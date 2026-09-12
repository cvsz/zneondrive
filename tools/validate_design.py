#!/usr/bin/env python3
"""Validate zNeonDrive design catalogs using only Python stdlib."""

from __future__ import annotations

import json
import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
CAT = ROOT / "design" / "catalog"


def load(name: str):
    with (CAT / name).open("r", encoding="utf-8") as fh:
        return json.load(fh)


def unique(items, label: str) -> set[str]:
    ids = [item["id"] for item in items]
    if len(ids) != len(set(ids)):
        raise AssertionError(f"duplicate {label} ids")
    return set(ids)


def main() -> int:
    characters_doc = load("characters.json")
    factions_doc = load("factions.json")
    districts_doc = load("districts.json")
    main_doc = load("quests-main.json")
    side_doc = load("quests-side.json")
    parts_doc = load("vehicle-parts.json")
    pets_doc = load("pets.json")
    slice_doc = load("vertical-slice.json")

    characters = characters_doc["characters"]
    factions = factions_doc["factions"]
    districts = districts_doc["districts"]
    quests = main_doc["quests"]
    side = side_doc["quests"]
    parts = parts_doc["part_families"]
    pets = pets_doc["pets"]

    assert characters_doc["count"] == len(characters) == 25
    assert factions_doc["count"] == len(factions) == 5
    assert districts_doc["count"] == len(districts) == 10
    assert main_doc["count"] == len(quests) == 100
    assert side_doc["count"] == len(side)
    assert pets_doc["count"] == len(pets)

    char_ids = unique(characters, "character")
    faction_ids = unique(factions, "faction")
    district_ids = unique(districts, "district")
    quest_ids = unique(quests, "main quest")
    unique(side, "side quest")
    unique(parts, "part")
    unique(pets, "pet")

    expected = {f"MQ{i:03d}" for i in range(1, 101)}
    assert quest_ids == expected, "main quest IDs must be exactly MQ001..MQ100"
    assert {q["chapter"] for q in quests} == set(range(1, 8))

    for idx, quest in enumerate(quests, start=1):
        assert quest["district"] in district_ids, f"{quest['id']}: unknown district"
        assert quest["primary_npc"] in char_ids, f"{quest['id']}: unknown primary NPC"
        if idx == 1:
            assert quest["prerequisites"] == []
        else:
            assert quest["prerequisites"] == [f"MQ{idx - 1:03d}"]

    for quest in side:
        assert quest["district"] in district_ids, f"{quest['id']}: unknown district"
        assert quest["primary_npc"] in char_ids, f"{quest['id']}: unknown primary NPC"

    allowed_character_factions = faction_ids | {"independent", "unknown"}
    for character in characters:
        assert character["faction"] in allowed_character_factions

    for pet in pets:
        assert pet["competitive_buff"] is False, f"{pet['id']}: pets may not grant competitive buffs"

    slice_quest_ids = [entry["quest_id"] for entry in slice_doc["quests"]]
    assert slice_doc["version"] == "0.2"
    assert slice_doc["district"] == "district_foundry_9"
    assert slice_quest_ids == [f"MQ{i:03d}" for i in range(1, 13)], (
        "vertical slice must cover MQ001..MQ012 in order"
    )
    assert len(slice_quest_ids) == len(set(slice_quest_ids))
    assert all(entry["proof"] for entry in slice_doc["quests"])
    assert len(slice_doc["race_proofs"]) >= 2
    assert all(race["requires_build_binding"] for race in slice_doc["race_proofs"])
    assert all(race["ordered_checkpoints"] for race in slice_doc["race_proofs"])
    assert slice_doc["vip_competitive_boost"] is False

    required_schemas = [
        ROOT / "design" / "schemas" / "quest.schema.json",
        ROOT / "design" / "schemas" / "vehicle.schema.json",
        ROOT / "design" / "schemas" / "vertical-slice.schema.json",
    ]
    for path in required_schemas:
        assert path.is_file(), f"missing schema: {path.relative_to(ROOT)}"
        with path.open("r", encoding="utf-8") as fh:
            json.load(fh)

    print(
        "design validation OK:",
        f"{len(characters)} characters,",
        f"{len(factions)} factions,",
        f"{len(districts)} districts,",
        f"{len(quests)} main quests,",
        f"{len(side)} side quests,",
        f"{len(parts)} part families,",
        f"{len(pets)} pets,",
        f"{len(slice_quest_ids)} vertical-slice quests",
    )
    return 0


if __name__ == "__main__":
    try:
        raise SystemExit(main())
    except (AssertionError, KeyError, json.JSONDecodeError) as exc:
        print(f"design validation FAILED: {exc}", file=sys.stderr)
        raise SystemExit(1)
