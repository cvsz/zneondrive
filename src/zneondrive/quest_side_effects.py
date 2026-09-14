"""Reference-only quest side-effect schedule for Python/Go parity tests."""

from __future__ import annotations

from dataclasses import dataclass

from .domain import STARTER_REBUILD_BLUEPRINT, parse_quest_id


@dataclass(frozen=True, slots=True)
class QuestSideEffects:
    inventory_item_id: str = ""
    inventory_quantity: int = 0
    blueprint_id: str = ""
    mark_starter_roadworthy: bool = False


def quest_side_effects(quest_id: str) -> QuestSideEffects:
    """Derive canonical quest side effects without granting production authority."""
    parse_quest_id(quest_id)
    if quest_id == "MQ004":
        return QuestSideEffects("part_brakes_track_i", 1)
    if quest_id == "MQ005":
        return QuestSideEffects(blueprint_id=STARTER_REBUILD_BLUEPRINT)
    if quest_id == "MQ009":
        return QuestSideEffects("part_tires_street_i", 1)
    if quest_id == "MQ012":
        return QuestSideEffects(mark_starter_roadworthy=True)
    return QuestSideEffects()
