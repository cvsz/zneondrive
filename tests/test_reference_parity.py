import json
from pathlib import Path
import unittest

from zneondrive.domain import WorldState


VECTORS_PATH = Path(__file__).with_name("reference-parity-vectors.json")


class ReferenceParityTests(unittest.TestCase):
    @classmethod
    def setUpClass(cls) -> None:
        cls.vectors = json.loads(VECTORS_PATH.read_text(encoding="utf-8"))
        if cls.vectors.get("schema_version") != 1:
            raise AssertionError("unsupported parity schema version")

    def test_shared_fixture_integrity(self) -> None:
        build_cases = self.vectors.get("build_hash_cases", [])
        quest_cases = self.vectors.get("quest_id_cases", [])
        self.assertTrue(build_cases, "build_hash_cases must not be empty")
        self.assertTrue(quest_cases, "quest_id_cases must not be empty")

        build_names: set[str] = set()
        for vector in build_cases:
            name = vector["name"]
            self.assertTrue(name, "build hash vector name must not be empty")
            self.assertNotIn(name, build_names, f"duplicate build hash vector name: {name}")
            build_names.add(name)
            self.assertTrue(vector["part_ids"], f"build hash vector {name} has no parts")
            expected_hash = vector["expected_hash"]
            self.assertRegex(expected_hash, r"^[0-9a-f]{64}$")

        quest_ids: set[str] = set()
        for vector in quest_cases:
            quest_id = vector["quest_id"]
            self.assertNotIn(quest_id, quest_ids, f"duplicate quest vector: {quest_id}")
            quest_ids.add(quest_id)
            if not vector["valid"]:
                self.assertIsNone(
                    vector["previous"],
                    f"invalid quest vector {quest_id} must not declare a predecessor",
                )

    def test_python_oracle_build_hashes_match_shared_vectors(self) -> None:
        for index, vector in enumerate(self.vectors["build_hash_cases"], start=1):
            with self.subTest(vector=vector["name"]):
                world = WorldState()
                account_id = f"acct-{index}"
                character_id = f"char-{index}"
                vehicle_id = f"veh-{index}"
                world.create_account(account_id, character_id)
                vehicle = world.grant_vehicle(
                    account_id,
                    vehicle_id,
                    vector["part_ids"],
                    operation_id=f"grant-{index}",
                )
                self.assertEqual(
                    vehicle.active_build.validation_hash,
                    vector["expected_hash"],
                )

    def test_shared_quest_vectors_remain_well_formed(self) -> None:
        for vector in self.vectors["quest_id_cases"]:
            quest_id = vector["quest_id"]
            with self.subTest(quest_id=quest_id):
                if vector["valid"]:
                    self.assertRegex(quest_id, r"^MQ\d{3}$")
                    number = int(quest_id[2:])
                    self.assertGreaterEqual(number, 1)
                    self.assertLessEqual(number, 100)
                    expected_previous = None if number == 1 else f"MQ{number - 1:03d}"
                    self.assertEqual(vector["previous"], expected_previous)
                else:
                    structurally_valid = (
                        len(quest_id) == 5
                        and quest_id.startswith("MQ")
                        and quest_id[2:].isdigit()
                        and 1 <= int(quest_id[2:]) <= 100
                    )
                    self.assertFalse(structurally_valid)


if __name__ == "__main__":
    unittest.main()
