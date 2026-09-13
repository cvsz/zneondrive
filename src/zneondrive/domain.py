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


STARTER_REBUILD_BLUEPRINT = "bp_starter_rebuild"
QUEST_ITEM_GRANTS = {
    "MQ004": "part_brakes_track_i",
    "MQ009": "part_tires_street_i",
}
QUEST_BLUEPRINT_GRANTS = {"MQ005": STARTER_REBUILD_BLUEPRINT}


class DomainError(RuntimeError):
    """Base error for rejected authoritative state transitions."""


class AuthorizationError(DomainError):
    """Raised when an account attempts to mutate state it does not own."""


class CapacityError(DomainError):
    """Raised when an entitlement capacity would be exceeded."""


class BuildConflictError(DomainError):
    """Raised when optimistic build revision preconditions do not match."""


class BlueprintRequiredError(DomainError):
    """Raised when a rebuild requires a blueprint the character has not unlocked."""


class InsufficientInventoryError(DomainError):
    """Raised when a rebuild requires parts the character does not own."""


class OperationConflictError(DomainError):
    """Raised when an idempotency key is reused for a different mutation."""


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
    inventory: dict[str, int] = field(default_factory=dict)
    blueprints: set[str] = field(default_factory=set)


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


def _operation_fingerprint(kind: str, payload: dict[str, Any]) -> str:
    canonical = json.dumps(
        {"kind": kind, "payload": payload},
        sort_keys=True,
        separators=(",", ":"),
        ensure_ascii=True,
    )
    return sha256(canonical.encode("utf-8")).hexdigest()


class WorldState:
    """In-memory reference authority with payload-bound idempotent receipts."""

    def __init__(self) -> None:
        self.accounts: dict[str, Account] = {}
        self.characters: dict[str, Character] = {}
        self.vehicles: dict[str, Vehicle] = {}
        self.race_registrations: dict[str, RaceRegistration] = {}
        self.race_results: dict[str, RaceResult] = {}
        self._operation_receipts: dict[str, Any] = {}
        self._operation_fingerprints: dict[str, str] = {}

    def _replay_operation(
        self, operation_id: str, kind: str, payload: dict[str, Any]
    ) -> tuple[bool, Any]:
        fingerprint = _operation_fingerprint(kind, payload)
        if operation_id not in self._operation_receipts:
            return False, fingerprint
        if self._operation_fingerprints.get(operation_id) != fingerprint:
            raise OperationConflictError(
                f"operation_id already used for different mutation: {operation_id}"
            )
        return True, self._operation_receipts[operation_id]

    def _record_operation(
        self, operation_id: str, fingerprint: str, receipt: Any
    ) -> Any:
        self._operation_fingerprints[operation_id] = fingerprint
        self._operation_receipts[operation_id] = receipt
        return receipt

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
        parts = tuple(sorted(set(part_ids)))
        replayed, replay_or_fingerprint = self._replay_operation(
            operation_id,
            "grant_vehicle",
            {
                "account_id": account_id,
                "vehicle_id": vehicle_id,
                "part_ids": parts,
                "starter_lineage": starter_lineage,
            },
        )
        if replayed:
            return replay_or_fingerprint
        account, character = self._owned_character(account_id)
        if len(character.vehicle_ids) >= account.entitlements.garage_slots:
            raise CapacityError("garage slot capacity exceeded")
        if vehicle_id in self.vehicles:
            raise DomainError(f"vehicle already exists: {vehicle_id}")
        initial = BuildRevision(1, parts, _build_hash(parts))
        vehicle = Vehicle(vehicle_id, character.character_id, starter_lineage, [initial])
        self.vehicles[vehicle_id] = vehicle
        character.vehicle_ids.append(vehicle_id)
        return self._record_operation(operation_id, replay_or_fingerprint, vehicle)

    def revise_build(
        self,
        account_id: str,
        vehicle_id: str,
        part_ids: Iterable[str],
        *,
        expected_revision: int,
        operation_id: str,
    ) -> BuildRevision:
        parts = tuple(sorted(set(part_ids)))
        replayed, replay_or_fingerprint = self._replay_operation(
            operation_id,
            "revise_build",
            {
                "account_id": account_id,
                "vehicle_id": vehicle_id,
                "part_ids": parts,
                "expected_revision": expected_revision,
            },
        )
        if replayed:
            return replay_or_fingerprint
        vehicle = self._owned_vehicle(account_id, vehicle_id)
        if vehicle.active_build.revision != expected_revision:
            raise BuildConflictError(
                f"expected revision {expected_revision}, current is {vehicle.active_build.revision}"
            )
        revision = BuildRevision(expected_revision + 1, parts, _build_hash(parts))
        vehicle.build_history.append(revision)
        return self._record_operation(operation_id, replay_or_fingerprint, revision)

    def rebuild_vehicle(
        self,
        account_id: str,
        vehicle_id: str,
        part_ids: Iterable[str],
        *,
        expected_revision: int,
        operation_id: str,
    ) -> BuildRevision:
        """Apply the inventory-authoritative Garage 17 rebuild reference contract.

        This mirrors the Go/PostgreSQL service-plane semantics without replacing them:
        the starter rebuild blueprint is mandatory, newly equipped parts are consumed,
        removed parts are returned, and the immutable build revision advances atomically.
        """
        parts = tuple(sorted(set(part_ids)))
        replayed, replay_or_fingerprint = self._replay_operation(
            operation_id,
            "rebuild_vehicle",
            {
                "account_id": account_id,
                "vehicle_id": vehicle_id,
                "part_ids": parts,
                "expected_revision": expected_revision,
            },
        )
        if replayed:
            return replay_or_fingerprint

        _, character = self._owned_character(account_id)
        vehicle = self._owned_vehicle(account_id, vehicle_id)
        if STARTER_REBUILD_BLUEPRINT not in character.blueprints:
            raise BlueprintRequiredError("starter rebuild blueprint is required")
        if vehicle.active_build.revision != expected_revision:
            raise BuildConflictError(
                f"expected revision {expected_revision}, current is {vehicle.active_build.revision}"
            )

        previous_parts = set(vehicle.active_build.part_ids)
        next_parts = set(parts)
        added_parts = sorted(next_parts - previous_parts)
        removed_parts = sorted(previous_parts - next_parts)
        missing_parts = [part_id for part_id in added_parts if character.inventory.get(part_id, 0) < 1]
        if missing_parts:
            raise InsufficientInventoryError(
                f"insufficient inventory for rebuild: {','.join(missing_parts)}"
            )

        # Validate every precondition before mutating inventory or build history.
        for part_id in added_parts:
            remaining = character.inventory.get(part_id, 0) - 1
            if remaining:
                character.inventory[part_id] = remaining
            else:
                character.inventory.pop(part_id, None)
        for part_id in removed_parts:
            character.inventory[part_id] = character.inventory.get(part_id, 0) + 1

        revision = BuildRevision(expected_revision + 1, parts, _build_hash(parts))
        vehicle.build_history.append(revision)
        return self._record_operation(operation_id, replay_or_fingerprint, revision)

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
        replayed, replay_or_fingerprint = self._replay_operation(
            operation_id,
            "complete_quest",
            {
                "account_id": account_id,
                "quest_id": quest_id,
                "money": money,
                "xp": xp,
                "reputation": reputation,
            },
        )
        if replayed:
            return replay_or_fingerprint
        account, character = self._owned_character(account_id)
        if quest_id in character.completed_quests:
            receipt = {"quest_id": quest_id, "money": 0, "xp": 0, "reputation": 0}
            return self._record_operation(operation_id, replay_or_fingerprint, receipt)
        if min(money, xp, reputation) < 0:
            raise DomainError("quest rewards cannot be negative")

        blueprint_id = QUEST_BLUEPRINT_GRANTS.get(quest_id)
        if blueprint_id and blueprint_id not in character.blueprints:
            if len(character.blueprints) >= account.entitlements.blueprint_slots:
                raise CapacityError("blueprint slot capacity exceeded")

        character.completed_quests.add(quest_id)
        character.money += money
        character.xp += xp
        character.reputation += reputation
        item_id = QUEST_ITEM_GRANTS.get(quest_id)
        if item_id:
            character.inventory[item_id] = character.inventory.get(item_id, 0) + 1
        if blueprint_id:
            character.blueprints.add(blueprint_id)

        receipt = {
            "quest_id": quest_id,
            "money": money,
            "xp": xp,
            "reputation": reputation,
        }
        return self._record_operation(operation_id, replay_or_fingerprint, receipt)

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
        times = tuple(checkpoint_times_ms)
        replayed, replay_or_fingerprint = self._replay_operation(
            operation_id,
            "submit_race_result",
            {
                "account_id": account_id,
                "race_id": race_id,
                "build_revision": build_revision,
                "checkpoint_times_ms": times,
            },
        )
        if replayed:
            return replay_or_fingerprint
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
        return self._record_operation(operation_id, replay_or_fingerprint, result)

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
