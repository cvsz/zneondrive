"""Reference-only race operation replay equality contract.

These helpers mirror the deterministic Go checks used after PostgreSQL resolves an
existing operation-id binding. They do not own operation-id uniqueness, persistence,
or production race authority.
"""


def _normalize_race_id(race_id: str) -> str | None:
    normalized = race_id.strip()
    if not (6 <= len(normalized) <= 128) or not normalized.startswith("race_"):
        return None
    if any(not ("a" <= ch <= "z" or "0" <= ch <= "9" or ch in "_-") for ch in normalized):
        return None
    return normalized


def race_start_replay_matches(*, existing: dict, request: dict) -> bool:
    normalized = _normalize_race_id(request["race_id"])
    return normalized is not None and (
        existing["account_id"] == request["account_id"]
        and existing["vehicle_id"] == request["vehicle_id"]
        and existing["race_id"] == normalized
    )


def race_checkpoint_replay_matches(*, existing: dict, request: dict) -> bool:
    return (
        existing["race_instance_id"] == request["race_instance_id"]
        and existing["checkpoint_index"] == request["checkpoint_index"]
        and existing["elapsed_ms"] == request["elapsed_ms"]
    )


def race_finish_replay_matches(*, existing: dict, request: dict) -> bool:
    return (
        existing["race_instance_id"] == request["race_instance_id"]
        and existing["checkpoint_count"] == request["checkpoint_count"]
        and existing["finish_elapsed_ms"] == request["finish_elapsed_ms"]
    )
