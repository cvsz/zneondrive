"""Reference-only rebuild operation replay contract."""

from dataclasses import dataclass


@dataclass(frozen=True, slots=True)
class BuildOperationReplay:
    owner_character_id: str
    vehicle_id: str
    result_revision: int
    validation_hash: str


def validate_build_operation_replay(
    existing: BuildOperationReplay,
    owner_character_id: str,
    vehicle_id: str,
    expected_revision: int,
    validation_hash: str,
) -> bool:
    """Return True only for an exact payload-bound durable replay."""
    if not owner_character_id or not vehicle_id or not validation_hash or expected_revision < 1:
        return False
    return (
        existing.owner_character_id == owner_character_id
        and existing.vehicle_id == vehicle_id
        and existing.result_revision == expected_revision + 1
        and existing.validation_hash == validation_hash
    )
