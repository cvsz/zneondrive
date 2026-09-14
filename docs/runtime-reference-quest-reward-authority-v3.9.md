# Runtime Reference Quest Reward Authority v3.9

## Scope

This increment removes caller-controlled quest reward amounts from the Python executable reference authority. `WorldState.complete_quest()` now accepts only account identity, canonical `MQ001..MQ100` quest identity, and an operation ID. Money, XP, and reputation are derived from the same deterministic schedule exercised by the Go/PostgreSQL production path:

- money = `100 + quest_number * 10`
- XP = `50 + quest_number * 5`
- reputation = `1`

The Python model remains a reference/test oracle only. Go 1.27 plus PostgreSQL remains the durable service-plane authority, Unreal Engine 5.8 remains gameplay authority, and Redis remains ephemeral/distributed coordination.

## Evidence

Automated Python coverage verifies that:

- `complete_quest()` exposes no `money`, `xp`, or `reputation` parameters;
- attempting to inject a reward amount is rejected by the Python call contract;
- MQ011 derives the canonical `210 money / 105 XP / 1 reputation` reward;
- replaying the same operation ID is idempotent;
- completing the same quest under a new operation ID does not grant a second reward;
- prerequisite, blueprint, inventory, rebuild, and cross-mutation idempotency tests continue to run through the same reference authority.

The existing v3.8 shared Python/Go reward vectors continue to pin formula parity against the production Go core.

## Security boundary

This change narrows the reference-oracle API so tests cannot accidentally model a client-authoritative reward path. It does not add a player-facing reward endpoint and does not move reward authority away from Go/PostgreSQL.

## Non-claims

This evidence does **not** prove:

- successful real Unreal Engine 5.8 Client or Dedicated Server builds/packages;
- live packaged Unreal↔Go quest completion;
- Garage 17 / MQ001–MQ012 playable completion;
- full Python↔Go reference-oracle parity;
- deployment-scale load/soak, HA/failover, regional DR, or production deployment readiness.

Those gates remain open until their required live/deployed evidence exists.
