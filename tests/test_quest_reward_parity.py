import json
from pathlib import Path
import unittest

from zneondrive.domain import DomainError
from zneondrive.quest_reward import quest_reward


ROOT = Path(__file__).resolve().parents[1]
VECTORS = ROOT / "tests" / "quest-reward-parity-v3.8.json"
POSTGRES = ROOT / "services" / "game-api" / "internal" / "store" / "postgres.go"


class QuestRewardParityTests(unittest.TestCase):
    @classmethod
    def setUpClass(cls) -> None:
        cls.vectors = json.loads(VECTORS.read_text(encoding="utf-8"))

    def test_schema_version(self) -> None:
        self.assertEqual(self.vectors["schema_version"], 1)

    def test_reference_reward_vectors(self) -> None:
        self.assertTrue(self.vectors["cases"])
        for case in self.vectors["cases"]:
            with self.subTest(quest_id=case["quest_id"]):
                reward = quest_reward(case["quest_id"])
                self.assertEqual(reward.quest_id, case["quest_id"])
                self.assertEqual(reward.money, case["money"])
                self.assertEqual(reward.xp, case["xp"])
                self.assertEqual(reward.reputation, case["reputation"])

    def test_invalid_quest_ids_are_rejected(self) -> None:
        for quest_id in self.vectors["invalid_quest_ids"]:
            with self.subTest(quest_id=quest_id):
                with self.assertRaises(DomainError):
                    quest_reward(quest_id)

    def test_postgres_authority_formula_remains_server_derived(self) -> None:
        source = POSTGRES.read_text(encoding="utf-8")
        required = (
            "questNumber, err := core.ParseQuestID(questID)",
            "money := int64(100 + questNumber*10)",
            "xp := int64(50 + questNumber*5)",
            "reputation := int64(1)",
        )
        for marker in required:
            self.assertIn(marker, source)

        # Public quest completion receives quest identity and operation identity,
        # not caller-controlled reward amounts.
        signature = "func (p *Postgres) CompleteQuest(ctx context.Context, tokenHash, questID, operationID string)"
        self.assertIn(signature, source)


if __name__ == "__main__":
    unittest.main()
