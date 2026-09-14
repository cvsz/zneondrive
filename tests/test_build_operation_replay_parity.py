import json
from pathlib import Path
import unittest

from zneondrive.build_operation import BuildOperationReplay, validate_build_operation_replay


class BuildOperationReplayParityTests(unittest.TestCase):
    def test_shared_vectors(self) -> None:
        fixture = json.loads(
            (Path(__file__).parent / "build-operation-replay-parity-v4.4.json").read_text()
        )
        for case in fixture["cases"]:
            with self.subTest(case=case["name"]):
                got = validate_build_operation_replay(
                    BuildOperationReplay(
                        owner_character_id=case["existing_owner"],
                        vehicle_id=case["existing_vehicle"],
                        result_revision=case["existing_result_revision"],
                        validation_hash=case["existing_hash"],
                    ),
                    case["owner"],
                    case["vehicle"],
                    case["expected_revision"],
                    case["hash"],
                )
                self.assertEqual(case["allow"], got)


if __name__ == "__main__":
    unittest.main()
