import json
from pathlib import Path
import unittest

from zneondrive.race_lifecycle import validate_checkpoint_advance, validate_finish


VECTORS = json.loads(
    (Path(__file__).with_name("race-lifecycle-parity-v4.0.json")).read_text(encoding="utf-8")
)


class RaceLifecycleParityTests(unittest.TestCase):
    def test_checkpoint_vectors(self) -> None:
        for case in VECTORS["checkpoint_cases"]:
            with self.subTest(case=case["name"]):
                got = validate_checkpoint_advance(
                    state=case["state"],
                    next_checkpoint=case["next_checkpoint"],
                    last_elapsed_ms=case["last_elapsed_ms"],
                    checkpoint_index=case["checkpoint_index"],
                    elapsed_ms=case["elapsed_ms"],
                )
                self.assertEqual(case["valid"], got)

    def test_finish_vectors(self) -> None:
        for case in VECTORS["finish_cases"]:
            with self.subTest(case=case["name"]):
                got = validate_finish(
                    state=case["state"],
                    next_checkpoint=case["next_checkpoint"],
                    last_elapsed_ms=case["last_elapsed_ms"],
                    checkpoint_count=case["checkpoint_count"],
                    finish_elapsed_ms=case["finish_elapsed_ms"],
                )
                self.assertEqual(case["valid"], got)


if __name__ == "__main__":
    unittest.main()
