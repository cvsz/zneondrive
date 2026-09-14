import json
from pathlib import Path

import pytest

from zneondrive.quest_side_effects import quest_side_effects


VECTORS = json.loads((Path(__file__).parent / "quest-side-effects-parity-v4.3.json").read_text())


def test_quest_side_effect_vectors_match_reference_oracle():
    assert VECTORS["schema_version"] == 1
    for case in VECTORS["cases"]:
        effect = quest_side_effects(case["quest_id"])
        assert effect.inventory_item_id == case["inventory_item_id"]
        assert effect.inventory_quantity == case["inventory_quantity"]
        assert effect.blueprint_id == case["blueprint_id"]
        assert effect.mark_starter_roadworthy is case["mark_starter_roadworthy"]


def test_quest_side_effects_reject_invalid_ids():
    for quest_id in VECTORS["invalid_quest_ids"]:
        with pytest.raises(RuntimeError):
            quest_side_effects(quest_id)
