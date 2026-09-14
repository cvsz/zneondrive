# Runtime Reference Quest Reward Parity v3.8

Status: **CI evidence only; production readiness remains evidence-gated.**

## Purpose

The Go/PostgreSQL service plane already derives MQ quest currency/XP/reputation from the canonical quest identifier instead of accepting reward amounts from a player or client. The Python reference surface did not have an independently executable copy of that schedule, so drift could go undetected.

v3.8 adds a shared reward vector corpus plus executable Python and Go contracts for the current authoritative schedule:

- money = `100 + quest_number * 10`
- XP = `50 + quest_number * 5`
- reputation = `1`
- valid quest range = `MQ001` through `MQ100`

The Python test also pins the production PostgreSQL mutation boundary: `CompleteQuest` receives authenticated session identity, quest ID, and operation ID; it does not receive caller-supplied reward amounts.

## Evidence

- `tests/quest-reward-parity-v3.8.json` contains boundary and Garage 17/MQ001-MQ012 representative vectors.
- `src/zneondrive/quest_reward.py` executes the reference-only reward schedule.
- `services/game-api/internal/core/quest_reward.go` executes the Go-side reward schedule.
- `tests/test_quest_reward_parity.py` validates the Python vectors, invalid IDs, and the existing PostgreSQL server-derived formula/boundary.
- `services/game-api/internal/core/quest_reward_parity_test.go` validates the same vector corpus in Go.

## Trust boundary

This increment does **not** move authority to Python. Unreal Engine 5.8 dedicated servers remain gameplay-authoritative, Go 1.27 remains the authenticated service plane, PostgreSQL remains durable authority, and Redis remains ephemeral/distributed coordination.

The production PostgreSQL mutation remains the actual durable reward authority. Python is an executable oracle only.

## Explicit non-claims

v3.8 does **not** prove or close:

- full Python/Go reference-oracle parity;
- real UE 5.8 Client or Dedicated Server build/package evidence;
- packaged Unreal-to-Go quest completion;
- Garage 17 or MQ001-MQ012 playable E2E evidence;
- live authoritative race lifecycle;
- production load/soak, ingress isolation, HA, PITR/RPO/RTO, regional DR, or deployment readiness.

The current `README.md`, `ROADMAP.md`, and `IMPLEMENTATION-CHECKLIST.md` production status therefore must remain unpromoted until those gates have matching evidence.
