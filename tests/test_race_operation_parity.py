import json
from pathlib import Path

from zneondrive.race_operation import (
    race_checkpoint_replay_matches,
    race_finish_replay_matches,
    race_start_replay_matches,
)


def test_race_operation_parity_vectors() -> None:
    vectors = json.loads((Path(__file__).with_name("race-operation-parity-v4.2.json")).read_text(encoding="utf-8"))
    groups = (
        ("start_cases", race_start_replay_matches),
        ("checkpoint_cases", race_checkpoint_replay_matches),
        ("finish_cases", race_finish_replay_matches),
    )
    for key, validator in groups:
        assert vectors[key], key
        for case in vectors[key]:
            assert validator(existing=case["existing"], request=case["request"]) is case["matches"], case["name"]
