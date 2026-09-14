"""Reference-only race-start eligibility rule.

Mirrors deterministic Go validation after PostgreSQL resolves ownership and the exact
active build. It owns no production persistence, session, vehicle, or race authority.
"""


def validate_race_start_binding(*, race_id: str, roadworthy: bool, build_revision: int, build_validation_hash: str) -> bool:
    normalized = race_id.strip()
    if not (6 <= len(normalized) <= 128) or not normalized.startswith("race_"):
        return False
    if any(not (ch.islower() or ch.isdigit() or ch in "_-") for ch in normalized):
        return False
    return roadworthy and build_revision >= 1 and bool(build_validation_hash.strip())
