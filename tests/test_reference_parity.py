import json
from pathlib import Path
import unittest

from zneondrive.domain import (
    DomainError,
    QuestOutOfOrderError,
    RaceValidationError,
    WorldState,
    parse_quest_id,
    previous_quest_id,
)
from zneondrive.race_contract import (
    normalize_race_id,
    race_result_hash,
    validate_race_checkpoint,
)


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
        quest_prerequisite_cases = self.vectors.get("quest_prerequisite_cases", [])
        race_id_cases = self.vectors.get("race_id_cases", [])
        checkpoint_cases = self.vectors.get("race_checkpoint_cases", [])
        result_hash_cases = self.vectors.get("race_result_hash_cases", [])
        self.assertTrue(build_cases, "build_hash_cases must not be empty")
        self.assertTrue(quest_cases, "quest_id_cases must not be empty")
        self.assertTrue(
            quest_prerequisite_cases, "quest_prerequisite_cases must not be empty"
        )
        self.assertTrue(race_id_cases, "race_id_cases must not be empty")
        self.assertTrue(checkpoint_cases, "race_checkpoint_cases must not be empty")
        self.assertTrue(result_hash_cases, "race_result_hash_cases must not be empty")

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

        quest_prerequisite_names: set[str] = set()
        for vector in quest_prerequisite_cases:
            name = vector["name"]
            self.assertTrue(name, "quest prerequisite vector name must not be empty")
            self.assertNotIn(
                name,
                quest_prerequisite_names,
                f"duplicate quest prerequisite vector name: {name}",
            )
            quest_prerequisite_names.add(name)
            self.assertRegex(vector["quest_id"], r"^MQ\d{3}$")
            self.assertEqual(len(vector["completed"]), len(set(vector["completed"])))

        race_names: set[str] = set()
        for vector in race_id_cases:
            name = vector["name"]
            self.assertTrue(name, "race id vector name must not be empty")
            self.assertNotIn(name, race_names, f"duplicate race id vector name: {name}")
            race_names.add(name)
            if vector["valid"]:
                self.assertIsInstance(vector["normalized"], str)
                self.assertTrue(vector["normalized"])
            else:
                self.assertIsNone(vector["normalized"])

        checkpoint_names: set[str] = set()
        for vector in checkpoint_cases:
            name = vector["name"]
            self.assertTrue(name, "race checkpoint vector name must not be empty")
            self.assertNotIn(name, checkpoint_names, f"duplicate checkpoint vector name: {name}")
            checkpoint_names.add(name)

        result_names: set[str] = set()
        required_instance_fields = {
            "race_instance_id",
            "race_id",
            "account_id",
            "character_id",
            "vehicle_id",
            "build_revision",
            "build_validation_hash",
        }
        for vector in result_hash_cases:
            name = vector["name"]
            self.assertTrue(name, "race result vector name must not be empty")
            self.assertNotIn(name, result_names, f"duplicate race result vector name: {name}")
            result_names.add(name)
            self.assertEqual(set(vector["instance"]), required_instance_fields)
            self.assertRegex(vector["instance"]["build_validation_hash"], r"^[0-9a-f]{64}$")
            if vector["valid"]:
                self.assertRegex(vector["expected_hash"], r"^[0-9a-f]{64}$")
            else:
                self.assertIsNone(vector["expected_hash"])

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

    def test_python_quest_ids_match_shared_vectors(self) -> None:
        for vector in self.vectors["quest_id_cases"]:
            quest_id = vector["quest_id"]
            with self.subTest(quest_id=quest_id):
                if not vector["valid"]:
                    with self.assertRaises(DomainError):
                        parse_quest_id(quest_id)
                    with self.assertRaises(DomainError):
                        previous_quest_id(quest_id)
                    continue

                number = parse_quest_id(quest_id)
                self.assertGreaterEqual(number, 1)
                self.assertLessEqual(number, 100)
                self.assertEqual(previous_quest_id(quest_id), vector["previous"])

    def test_python_quest_prerequisites_match_shared_vectors(self) -> None:
        for index, vector in enumerate(
            self.vectors["quest_prerequisite_cases"], start=1
        ):
            with self.subTest(vector=vector["name"]):
                world = WorldState()
                account_id = f"acct-prereq-{index}"
                character_id = f"char-prereq-{index}"
                world.create_account(account_id, character_id)
                world.characters[character_id].completed_quests.update(vector["completed"])
                if vector["allowed"]:
                    world.complete_quest(
                        account_id,
                        vector["quest_id"],
                        operation_id=f"quest-prereq-{index}",
                    )
                    self.assertIn(
                        vector["quest_id"],
                        world.characters[character_id].completed_quests,
                    )
                else:
                    before = set(world.characters[character_id].completed_quests)
                    with self.assertRaises(QuestOutOfOrderError):
                        world.complete_quest(
                            account_id,
                            vector["quest_id"],
                            operation_id=f"quest-prereq-{index}",
                        )
                    self.assertEqual(
                        world.characters[character_id].completed_quests,
                        before,
                    )

    def test_python_race_ids_match_shared_vectors(self) -> None:
        for vector in self.vectors["race_id_cases"]:
            with self.subTest(vector=vector["name"]):
                if vector["valid"]:
                    self.assertEqual(
                        normalize_race_id(vector["race_id"]),
                        vector["normalized"],
                    )
                else:
                    with self.assertRaises(RaceValidationError):
                        normalize_race_id(vector["race_id"])

    def test_python_race_checkpoints_match_shared_vectors(self) -> None:
        for vector in self.vectors["race_checkpoint_cases"]:
            with self.subTest(vector=vector["name"]):
                if vector["valid"]:
                    validate_race_checkpoint(vector["index"], vector["elapsed_ms"])
                else:
                    with self.assertRaises(RaceValidationError):
                        validate_race_checkpoint(vector["index"], vector["elapsed_ms"])

    def test_python_race_result_hashes_match_shared_vectors(self) -> None:
        for vector in self.vectors["race_result_hash_cases"]:
            with self.subTest(vector=vector["name"]):
                if vector["valid"]:
                    self.assertEqual(
                        race_result_hash(
                            vector["instance"],
                            vector["checkpoint_count"],
                            vector["finish_elapsed_ms"],
                        ),
                        vector["expected_hash"],
                    )
                else:
                    with self.assertRaises(RaceValidationError):
                        race_result_hash(
                            vector["instance"],
                            vector["checkpoint_count"],
                            vector["finish_elapsed_ms"],
                        )


if __name__ == "__main__":
    unittest.main()
