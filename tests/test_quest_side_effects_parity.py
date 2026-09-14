import json
from pathlib import Path
import unittest

from zneondrive.quest_side_effects import quest_side_effects


VECTORS = json.loads((Path(__file__).parent / "quest-side-effects-parity-v4.3.json").read_text())


class QuestSideEffectsParityTests(unittest.TestCase):
    def test_vectors_match_reference_oracle(self):
        self.assertEqual(VECTORS["schema_version"], 1)
        for case in VECTORS["cases"]:
            effect = quest_side_effects(case["quest_id"])
            self.assertEqual(effect.inventory_item_id, case["inventory_item_id"])
            self.assertEqual(effect.inventory_quantity, case["inventory_quantity"])
            self.assertEqual(effect.blueprint_id, case["blueprint_id"])
            self.assertEqual(effect.mark_starter_roadworthy, case["mark_starter_roadworthy"])

    def test_rejects_invalid_ids(self):
        for quest_id in VECTORS["invalid_quest_ids"]:
            with self.assertRaises(RuntimeError):
                quest_side_effects(quest_id)


if __name__ == "__main__":
    unittest.main()
