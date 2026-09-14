import json
from pathlib import Path

from zneondrive.race_start import validate_race_start_binding


def test_race_start_parity_vectors() -> None:
    vectors = json.loads((Path(__file__).with_name("race-start-parity-v4.1.json")).read_text(encoding="utf-8"))
    assert vectors["cases"]
    for case in vectors["cases"]:
        actual = validate_race_start_binding(
            race_id=case["race_id"],
            roadworthy=case["roadworthy"],
            build_revision=case["build_revision"],
            build_validation_hash=case["build_validation_hash"],
        )
        assert actual is case["valid"], case["name"]
