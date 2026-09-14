"""Cross-language reference contract for authoritative race validation.

This module mirrors bounded, dependency-free helpers used by the Go service
plane. It is reference/test code only; the Go/PostgreSQL runtime and Unreal
dedicated server remain authoritative for production race state.
"""

from __future__ import annotations

from hashlib import sha256
import json
from typing import Mapping

from .domain import RaceValidationError

MAX_RACE_ID_LENGTH = 128
MAX_CHECKPOINT_INDEX = 1024
MAX_CHECKPOINT_ELAPSED_MS = 24 * 60 * 60 * 1000
MAX_RESULT_CHECKPOINT_COUNT = 1025


def normalize_race_id(race_id: str) -> str:
    """Normalize and validate the canonical service-plane race identifier."""
    normalized = race_id.strip()
    if (
        len(normalized) < 6
        or len(normalized) > MAX_RACE_ID_LENGTH
        or not normalized.startswith("race_")
    ):
        raise RaceValidationError("invalid race id")
    if any(
        not ("a" <= char <= "z" or "0" <= char <= "9" or char in "_-")
        for char in normalized
    ):
        raise RaceValidationError("invalid race id")
    return normalized


def validate_race_checkpoint(index: int, elapsed_ms: int) -> None:
    """Validate the bounded checkpoint cursor/time accepted by the Go core."""
    if (
        index < 0
        or index > MAX_CHECKPOINT_INDEX
        or elapsed_ms <= 0
        or elapsed_ms > MAX_CHECKPOINT_ELAPSED_MS
    ):
        raise RaceValidationError("invalid race checkpoint")


def race_result_hash(
    instance: Mapping[str, object], checkpoint_count: int, finish_elapsed_ms: int
) -> str:
    """Mirror Go's deterministic authoritative race-result digest contract.

    The digest binds the durable race instance identity, account/character/vehicle
    ownership, immutable build revision/hash, checkpoint count, and finish time.
    It deliberately does not accept player-computed outcome data beyond the
    already-authoritative fields supplied by the service plane.
    """
    if (
        checkpoint_count < 1
        or checkpoint_count > MAX_RESULT_CHECKPOINT_COUNT
        or finish_elapsed_ms <= 0
    ):
        raise RaceValidationError("invalid race result")

    required = (
        "race_instance_id",
        "race_id",
        "account_id",
        "character_id",
        "vehicle_id",
        "build_revision",
        "build_validation_hash",
    )
    if any(key not in instance for key in required):
        raise RaceValidationError("incomplete race instance")

    payload = {
        "race_instance_id": instance["race_instance_id"],
        "race_id": instance["race_id"],
        "account_id": instance["account_id"],
        "character_id": instance["character_id"],
        "vehicle_id": instance["vehicle_id"],
        "build_revision": instance["build_revision"],
        "build_validation_hash": instance["build_validation_hash"],
        "checkpoint_count": checkpoint_count,
        "finish_elapsed_ms": finish_elapsed_ms,
    }
    canonical = json.dumps(payload, separators=(",", ":"), ensure_ascii=True)
    return sha256(canonical.encode("utf-8")).hexdigest()
