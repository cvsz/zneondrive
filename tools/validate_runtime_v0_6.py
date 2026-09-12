#!/usr/bin/env python3
"""Static contract checks for Phase 4.2 inventory/rebuild runtime v0.6."""

import json
import re
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]

required = [
    "docs/runtime-inventory-rebuild-v0.6.md",
    "services/game-api/internal/core/parts_catalog.go",
    "services/game-api/internal/store/migrations/003_inventory_blueprints.sql",
    "services/game-api/internal/store/postgres_integration_test.go",
    "game/Source/NeonDrive/NDServiceSubsystem.h",
    "game/Source/NeonDrive/NDServiceSubsystem.cpp",
]
for relative in required:
    assert (ROOT / relative).is_file(), f"missing v0.6 runtime file: {relative}"

catalog = json.loads((ROOT / "design/catalog/vehicle-parts.json").read_text(encoding="utf-8"))
design_ids = {part["id"] for part in catalog["part_families"]}
go_catalog = (ROOT / "services/game-api/internal/core/parts_catalog.go").read_text(encoding="utf-8")
runtime_ids = set(re.findall(r'"(part_[a-z0-9_]+)"\s*:\s*\{\}', go_catalog))
assert runtime_ids == design_ids, (
    f"runtime/catalog part mismatch: missing={sorted(design_ids-runtime_ids)} "
    f"extra={sorted(runtime_ids-design_ids)}"
)
assert 'StarterRebuildBlueprint = "bp_starter_rebuild"' in go_catalog

core = (ROOT / "services/game-api/internal/core/core.go").read_text(encoding="utf-8")
for token in ('json:"inventory"', 'json:"blueprints"', "!IsCatalogPart(partID)"):
    assert token in core, f"missing v0.6 core invariant: {token}"
assert re.search(r"Inventory\s+\[\]InventoryItem", core), "missing typed inventory snapshot field"
assert re.search(r"Blueprints\s+\[\]string", core), "missing typed blueprint snapshot field"

migration = (ROOT / "services/game-api/internal/store/migrations/003_inventory_blueprints.sql").read_text(encoding="utf-8")
for token in (
    "inventory_items",
    "quantity BIGINT NOT NULL DEFAULT 0 CHECK (quantity >= 0)",
    "character_blueprints",
    "PRIMARY KEY (character_id, item_id)",
    "PRIMARY KEY (character_id, blueprint_id)",
):
    assert token in migration, f"missing v0.6 migration invariant: {token}"

postgres = (ROOT / "services/game-api/internal/store/postgres.go").read_text(encoding="utf-8")
assert postgres.index('"002_game_tickets"') < postgres.index('"003_inventory_blueprints"')
for token in (
    "ErrBlueprintRequired",
    "ErrInsufficientInventory",
    "part_brakes_track_i",
    "part_tires_street_i",
    "quantity=quantity-1",
    "inventory_items.quantity+1",
    "StarterRebuildBlueprint",
    "validation_hash",
    "FOR UPDATE",
):
    assert token in postgres, f"missing v0.6 transactional invariant: {token}"

tests = (ROOT / "services/game-api/internal/store/postgres_integration_test.go").read_text(encoding="utf-8")
for token in (
    "MQ004 retry duplicated salvage inventory",
    "ErrBlueprintRequired",
    "ErrInsufficientInventory",
    "idempotent build retry changed revision or inventory",
    "resume lost inventory/blueprint state",
):
    assert token in tests, f"missing v0.6 integration evidence assertion: {token}"

unreal_h = (ROOT / "game/Source/NeonDrive/NDServiceSubsystem.h").read_text(encoding="utf-8")
for token in ("FNDInventoryItem", "TArray<FNDInventoryItem> Inventory", "TArray<FString> Blueprints"):
    assert token in unreal_h, f"missing Unreal inventory snapshot invariant: {token}"

unreal_cpp = (ROOT / "game/Source/NeonDrive/NDServiceSubsystem.cpp").read_text(encoding="utf-8")
for token in ('TEXT("inventory")', 'TEXT("blueprints")', "Parsed.Inventory", "Parsed.Blueprints"):
    assert token in unreal_cpp, f"missing Unreal inventory parsing invariant: {token}"

print("runtime inventory/rebuild v0.6 static validation OK")
