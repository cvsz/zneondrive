"""Reference-only race lifecycle acceptance rules.

This module mirrors the deterministic Go core validation used by the PostgreSQL
production authority. It does not own persistence, transport, idempotency, or
production gameplay authority.
"""

from __future__ import annotations


def validate_checkpoint_advance(
    *,
    state: str,
    next_checkpoint: int,
    last_elapsed_ms: int,
    checkpoint_index: int,
    elapsed_ms: int,
) -> bool:
    if checkpoint_index < 0 or checkpoint_index > 1024:
        return False
    if elapsed_ms <= 0 or elapsed_ms > 24 * 60 * 60 * 1000:
        return False
    return (
        state == "active"
        and checkpoint_index == next_checkpoint
        and elapsed_ms > last_elapsed_ms
    )


def validate_finish(
    *,
    state: str,
    next_checkpoint: int,
    last_elapsed_ms: int,
    checkpoint_count: int,
    finish_elapsed_ms: int,
) -> bool:
    return (
        state == "active"
        and checkpoint_count == next_checkpoint
        and checkpoint_count >= 1
        and finish_elapsed_ms > last_elapsed_ms
    )
