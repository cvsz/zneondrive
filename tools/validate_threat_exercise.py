#!/usr/bin/env python3
from __future__ import annotations

import json
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
MANIFEST = ROOT / "security" / "threat-exercise-v3.0.json"

REQUIRED_CATEGORIES = {
    "credential_replay",
    "economy_reward_duplication",
    "vehicle_build_forgery",
    "race_cheating",
    "availability_and_identity_abuse",
}
ALLOWED_STATUS = {"ci_evidenced", "open"}


def fail(message: str) -> None:
    raise SystemExit(f"threat exercise validation failed: {message}")


def main() -> None:
    data = json.loads(MANIFEST.read_text(encoding="utf-8"))
    if data.get("schema_version") != 1:
        fail("schema_version must be 1")
    if data.get("evidence_level") != "ci_exercised_partial":
        fail("evidence_level must remain ci_exercised_partial until deployed/live gates pass")

    architecture = data.get("architecture") or {}
    expected_architecture = {
        "gameplay_authority": "Unreal Engine 5.8 dedicated server",
        "service_plane": "Go 1.27",
        "durable_authority": "PostgreSQL",
        "ephemeral_coordination": "Redis",
    }
    if architecture != expected_architecture:
        fail("selected Unreal/Go/PostgreSQL/Redis trust boundary changed")

    scenarios = data.get("scenarios")
    if not isinstance(scenarios, list) or not scenarios:
        fail("scenarios must be a non-empty list")

    seen_ids: set[str] = set()
    categories: set[str] = set()
    evidenced = 0
    open_count = 0

    for scenario in scenarios:
        scenario_id = scenario.get("id")
        category = scenario.get("category")
        status = scenario.get("status")
        if not isinstance(scenario_id, str) or not scenario_id:
            fail("every scenario needs a non-empty id")
        if scenario_id in seen_ids:
            fail(f"duplicate scenario id: {scenario_id}")
        seen_ids.add(scenario_id)
        if not isinstance(category, str) or not category:
            fail(f"{scenario_id}: missing category")
        categories.add(category)
        if status not in ALLOWED_STATUS:
            fail(f"{scenario_id}: invalid status {status!r}")

        if status == "ci_evidenced":
            evidenced += 1
            evidence = scenario.get("evidence")
            if not isinstance(evidence, list) or not evidence:
                fail(f"{scenario_id}: ci_evidenced scenario needs evidence")
            if scenario.get("blockers"):
                fail(f"{scenario_id}: ci_evidenced scenario must not carry blockers")
            for item in evidence:
                path_value = item.get("path")
                markers = item.get("markers")
                if not isinstance(path_value, str) or not path_value:
                    fail(f"{scenario_id}: evidence path missing")
                path = ROOT / path_value
                if not path.is_file():
                    fail(f"{scenario_id}: missing evidence file {path_value}")
                if not isinstance(markers, list) or not markers:
                    fail(f"{scenario_id}: evidence markers missing for {path_value}")
                text = path.read_text(encoding="utf-8")
                for marker in markers:
                    if not isinstance(marker, str) or not marker:
                        fail(f"{scenario_id}: empty evidence marker")
                    if marker not in text:
                        fail(f"{scenario_id}: marker not found in {path_value}: {marker}")
        else:
            open_count += 1
            blockers = scenario.get("blockers")
            if not isinstance(blockers, list) or not blockers:
                fail(f"{scenario_id}: open scenario requires explicit blockers")
            if scenario.get("evidence"):
                fail(f"{scenario_id}: open scenario must not masquerade as evidenced")

    missing_categories = REQUIRED_CATEGORIES - categories
    if missing_categories:
        fail(f"missing required abuse categories: {sorted(missing_categories)}")
    if evidenced < 5:
        fail("expected at least five CI-evidenced threat scenarios")
    if open_count < 3:
        fail("open production/live blockers must remain explicit")

    print(
        "threat exercise evidence valid: "
        f"ci_evidenced={evidenced} open={open_count} categories={len(categories)}"
    )


if __name__ == "__main__":
    main()
