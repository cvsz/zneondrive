"""Technology-neutral executable domain model for PROJECT: NEON DRIVE.

This module is intentionally small and dependency-free. It is not a production
MMORPG server. It exists to make the design invariants executable before a
runtime stack is selected.
"""

from __future__ import annotations

from dataclasses import dataclass, field
from hashlib import sha256
import json
from typing import Any, Iterable


class DomainError(RuntimeError):
    """Base error for rejected authoritative state transitions."""


class AuthorizationError(DomainError):
    """Raised when an account attempts to mutate state it does not own."""


class CapacityError(DomainError):
    """Raised when an entitlement capacity would be exceeded."""


class BuildConflictError(DomainError):
    """Raised when optimistic build revision preconditions do not match."""


class RaceValidationError(DomainError):
    """Raised when a race result cannot be proven against accepted state."""


@dataclass(frozen=True, slots=True)
class Entitlements:
    garage_slots: int = 1
    blueprint_slots: int = 32

    def __post_init__(self) -> None:
        if self.garage_slots < 1:
            raise ValueError("garage_slots must be >= 1")
        if self.blueprint_slots < 1:
            raise ValueError("blueprint_slots must be >= 1")


@dataclass(frozen=True, slots=True)
class BuildRevision:
    revision: int
    part_ids: tuple[str, ...]
    validation_hash: str


@dataclass(slots=True)
class Vehicle:
    vehicle_id: str
    owner_character_id: str
    starter_lineage: bool
    build_history: list[BuildRevision]
    reputation: int = 0

    @property
    def active_build(self) -> BuildRevision:
        return self.build_history[-1]


@dataclass(slots=True)
class Character:
    character_id: str
    vehicle_ids: list[str] = field(default_factory=list)
    money: int = 0
    xp: int = 0
    reputation: int = 0
    completed_quests: set[str] = field(default_factory=set)


@dataclass(slots=True)
class Account:
    account_id: str
    character_id: str
    entitlements: Entitlements = field(default_factory=Entitlements)


@dataclass(frozen=True, slots=True)
class RaceRegistration:
    race_id: str
    vehicle_id: str
    accepted_build_revision: int
    accepted_build_hash: str
    ordered_checkpoints: tuple[str, ...]


@dataclass(frozen=True, slots=True)
class RaceResult:
    race_id: str
    vehicle_id: str
    build_revision: int
    elapsed_ms: int
    checkpoint_times_ms: tuple[int, ...]


def _build_hash(part_ids: Iterable[str]) -> str:
    canonical = json.dumps(sorted(part_ids), separators=(",", ":"), ensure_ascii=True)
    return sha256(canonical.encode("utf-8")).hexdigest()


class WorldState:
    """In-memory reference authority with idempotent mutation receipts."""

    def __init__(self) -> None:
        self.accounts: dict[str, Account] = {}
        self.characters: dict[str, Character] = {}
        self.vehicles: dict[str, Vehicle] = {}
        self.race_registrations: dict[str, RaceRegistration] = {}
        self.race_results: dict[str, RaceResult] = {}
        self._operation_receipts: dict[str, Any] = {}

    def create_account(
        self,
        account_id: str,
        character_id: str,
        *,
        entitlements: Entitlements | None = None,
    ) -> Account:
        if account_id in self.accounts:
            raise DomainError(f"account already exists: {account_id}")
        if character_id in self.characters:
            raise DomainError(f"character already exists: {character_id}")
        account = Account(account_id, character_id, entitlements or Entitlements())
        self.accounts[account_id] = account
        self.characters[character_id] = Character(character_id)
        return account

    def grant_vehicle(
        self,
        account_id: str,
        vehicle_id: str,
        part_ids: Iterable[str],
        *,
        starter_lineage: bool = False,
        operation_id: str,
    ) -> Vehicle:
        if operation_id in self._operation_receipts:
            return self._operation_receipts[operation_id]
        account, character = self._owned_character(account_id)
        if len(character.vehicle_ids) >= account.entitlements.garage_slots:
            raise CapacityError("garage slot capacity exceeded")
        if vehicle_id in self.vehicles:
            raise DomainError(f"vehicle already exists: {vehicle_id}")
        parts = tuple(sorted(set(part_ids)))
        initial = BuildRevision(1, parts, _build_hash(parts))
        vehicle = Vehicle(vehicle_id, character.character_id, starter_lineage, [initial])
        self.vehicles[vehicle_id] = vehicle
        character.vehicle_ids.append(vehicle_id)
        self._operation_receipts[operation_id] = vehicle
        return vehicle

    def revise_build(
        self,
        account_id: str,
        vehicle_id: str,
        part_ids: Iterable[str],
        *,
        expected_revision: int,
        operation_id: str,
    ) -> BuildRevision:
        if operation_id in self._operation_receipts:
            return self._operation_receipts[operation_id]
        vehicle = self._owned_vehicle(account_id, vehicle_id)
        if vehicle.active_build.revision != expected_revision:
            raise BuildConflictError(
                f"expected revision {expected_revision}, current is {vehicle.active_build.revision}"
            )
        parts = tuple(sorted(set(part_ids)))
        revision = BuildRevision(expected_revision + 1, parts, _build_hash(parts))
        vehicle.build_history.append(revision)
        self._operation_receipts[operation_id] = revision
        return revision

    def complete_quest(
        self,
        account_id: str,
        quest_id: str,
        *,
        money: int = 0,
        xp: int = 0,
        reputation: int = 0,
        operation_id: str,
    ) -> dict[str, int | str]:
        if operation_id in self._operation_receipts:
            return self._operation_receipts[operation_id]
        _, character = self._owned_character(account_id)
        if quest_id in character.completed_quests:
            receipt = {"quest_id": quest_id, "money": 0, "xp": 0, "reputation": 0}
            self._operation_receipts[operation_id] = receipt
            return receipt
        if min(money, xp, reputation) < 0:
            raise DomainError("quest rewards cannot be negative")
        character.completed_quests.add(quest_id)
        character.money += money
        character.xp += xp
        character.reputation += reputation
        receipt = {
            "quest_id": quest_id,
            "money": money,
            "xp": xp,
            "reputation": reputation,
        }
        self._operation_receipts[operation_id] = receipt
        return receipt

    def register_race(
        self,
        account_id: str,
        race_id: str,
        vehicle_id: str,
        ordered_checkpoints: Iterable[str],
    ) -> RaceRegistration:
        vehicle = self._owned_vehicle(account_id, vehicle_id)
        checkpoints = tuple(ordered_checkpoints)
        if len(checkpoints) < 2 or len(checkpoints) != len(set(checkpoints)):
            raise RaceValidationError("race requires at least two unique ordered checkpoints")
        registration = RaceRegistration(
            race_id=race_id,
            vehicle_id=vehicle_id,
            accepted_build_revision=vehicle.active_build.revision,
            accepted_build_hash=vehicle.active_build.validation_hash,
            ordered_checkpoints=checkpoints,
        )
        self.race_registrations[race_id] = registration
        return registration

    def submit_race_result(
        self,
        account_id: str,
        race_id: str,
        *,
        build_revision: int,
        checkpoint_times_ms: Iterable[int],
        operation_id: str,
    ) -> RaceResult:
        if operation_id in self._operation_receipts:
            return self._operation_receipts[operation_id]
        registration = self.race_registrations.get(race_id)
        if registration is None:
            raise RaceValidationError("unknown race registration")
        vehicle = self._owned_vehicle(account_id, registration.vehicle_id)
        if build_revision != registration.accepted_build_revision:
            raise RaceValidationError("result build revision differs from accepted revision")
        accepted = next(
            (b for b in vehicle.build_history if b.revision == build_revision), None
        )
        if accepted is None or accepted.validation_hash != registration.accepted_build_hash:
            raise RaceValidationError("accepted build provenance cannot be verified")
        times = tuple(checkpoint_times_ms)
        if len(times) != len(registration.ordered_checkpoints):
            raise RaceValidationError("checkpoint count does not match race ruleset")
        if any(t <= 0 for t in times) or any(b <= a for a, b in zip(times, times[1:])):
            raise RaceValidationError("checkpoint times must be positive and strictly increasing")
        result = RaceResult(
            race_id=race_id,
            vehicle_id=registration.vehicle_id,
            build_revision=build_revision,
            elapsed_ms=times[-1],
            checkpoint_times_ms=times,
        )
        self.race_results[race_id] = result
        self._operation_receipts[operation_id] = result
        return result

    def delete_vehicle(self, account_id: str, vehicle_id: str) -> None:
        vehicle = self._owned_vehicle(account_id, vehicle_id)
        if vehicle.starter_lineage:
            raise DomainError("starter-lineage vehicle cannot be routine-deleted")
        _, character = self._owned_character(account_id)
        character.vehicle_ids.remove(vehicle_id)
        del self.vehicles[vehicle_id]

    def competitive_build_signature(self, vehicle_id: str) -> str:
        """Return a signature determined only by build state, never VIP capacity."""
        return self.vehicles[vehicle_id].active_build.validation_hash

    def _owned_character(self, account_id: str) -> tuple[Account, Character]:
        account = self.accounts.get(account_id)
        if account is None:
            raise AuthorizationError("unknown account")
        return account, self.characters[account.character_id]

    def _owned_vehicle(self, account_id: str, vehicle_id: str) -> Vehicle:
        _, character = self._owned_character(account_id)
        vehicle = self.vehicles.get(vehicle_id)
        if vehicle is None or vehicle.owner_character_id != character.character_id:
            raise AuthorizationError("vehicle is not owned by account character")
        return vehicle
