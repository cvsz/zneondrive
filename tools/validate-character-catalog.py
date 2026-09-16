#!/usr/bin/env python3
"""Validate the source-of-truth character roster and 3D presentation contract."""
from __future__ import annotations

import json
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
ROSTER = ROOT / "design/catalog/characters.json"
PRESENTATION = ROOT / "design/catalog/character-presentation.json"


def load(path: Path) -> dict:
    with path.open(encoding="utf-8") as handle:
        return json.load(handle)


def main() -> int:
    roster = load(ROSTER)
    presentation = load(PRESENTATION)
    characters = roster.get("characters", [])
    visual = presentation.get("characters", {})
    expected = roster.get("count")

    errors: list[str] = []
    ids = [item.get("id") for item in characters]
    if expected != len(characters):
        errors.append(f"roster count={expected!r} but contains {len(characters)} entries")
    if len(ids) != len(set(ids)):
        errors.append("roster contains duplicate character IDs")
    if any(not value or not isinstance(value, str) for value in ids):
        errors.append("every character must have a non-empty string id")

    missing = sorted(set(ids) - set(visual))
    extra = sorted(set(visual) - set(ids))
    if missing:
        errors.append("missing presentation profiles: " + ", ".join(missing))
    if extra:
        errors.append("presentation profiles without roster entries: " + ", ".join(extra))

    tiers = set(presentation.get("tiers", {}))
    for character_id, profile in visual.items():
        if profile.get("tier") not in tiers:
            errors.append(f"{character_id}: invalid presentation tier {profile.get('tier')!r}")
        for field in ("motion", "facial"):
            if not isinstance(profile.get(field), str) or not profile[field]:
                errors.append(f"{character_id}: {field} must be a non-empty string")

    if errors:
        for error in errors:
            print(f"ERROR: {error}")
        return 1

    print(f"OK: validated {len(characters)} character identities and presentation profiles")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
