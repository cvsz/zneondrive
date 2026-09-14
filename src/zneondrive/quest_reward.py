"""Reference-only quest reward schedule mirroring the Go/PostgreSQL authority.

This module does not grant rewards. It makes the server-derived MQ reward formula
executable on the Python oracle side so CI can detect cross-language drift.
"""

from __future__ import annotations

from dataclasses import dataclass

from .domain import parse_quest_id


@dataclass(frozen=True, slots=True)
class QuestReward:
    quest_id: str
    money: int
    xp: int
    reputation: int


def quest_reward(quest_id: str) -> QuestReward:
    """Return the canonical server-derived reward for MQ001..MQ100."""
    number = parse_quest_id(quest_id)
    return QuestReward(
        quest_id=quest_id,
        money=100 + number * 10,
        xp=50 + number * 5,
        reputation=1,
    )
