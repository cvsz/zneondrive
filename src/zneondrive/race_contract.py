"""Cross-language reference contract for authoritative race input validation.

This module mirrors the bounded, dependency-free validation helpers used by the
Go service plane. It is reference/test code only; the Go/PostgreSQL runtime and
Unreal dedicated server remain authoritative for production race state.
"""

from __future__ import annotations

from .domain import RaceValidationError

MAX_RACE_ID_LENGTH = 128
MAX_CHECKPOINT_INDEX = 1024
MAX_CHECKPOINT_ELAPSED_MS = 24 * 60 * 60 * 1000


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
