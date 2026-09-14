# Runtime Reference Quest Side-Effect Parity v4.3

## Scope

This increment pins deterministic side effects for the canonical main-quest IDs that currently mutate inventory, blueprint knowledge, or starter-vehicle Roadworthy state:

- `MQ004` grants one `part_brakes_track_i` inventory item.
- `MQ005` unlocks `bp_starter_rebuild`.
- `MQ009` grants one `part_tires_street_i` inventory item.
- `MQ012` marks starter-lineage vehicles Roadworthy.
- other canonical `MQ001..MQ100` IDs have no side effect in this contract.

The shared fixture `tests/quest-side-effects-parity-v4.3.json` is consumed by executable Python and Go tests. Both implementations reject malformed/out-of-range quest IDs before deriving effects.

## Authority boundary

This is parity evidence, not a new production mutation path. Go/PostgreSQL remains the durable service-plane authority that applies quest completion, rewards, inventory grants, blueprint unlocks, and the MQ012 Roadworthy transition transactionally. Python remains a reference/test oracle only. Unreal Engine 5.8 remains the gameplay authority and does not gain a client-authoritative progression path from this work.

The Go helper accepts only a canonical quest ID and derives fixed effects; it accepts no client-supplied item ID, quantity, blueprint ID, or Roadworthy flag.

## Evidence

Automated evidence consists of:

- Go unit/parity tests consuming the shared v4.3 vectors;
- Python parity tests consuming the same vectors;
- existing PostgreSQL integration coverage for duplicate-grant prevention, durable inventory/blueprint persistence, rebuild accounting, and MQ012 progression;
- existing static/runtime validators that keep live Unreal and production deployment claims evidence-gated.

## Non-claims / open gates

This increment does **not** claim:

- full Python↔Go reference-oracle parity;
- successful real UE 5.8 Client or Dedicated Server build/package evidence;
- live packaged Unreal↔Go quest/inventory integration;
- Garage 17 or MQ001–MQ012 playable end-to-end completion;
- final-physics anti-cheat calibration;
- deployment-scale load/soak, deployed observability/SLOs, HA/failover, PITR/RPO/RTO, regional DR, or production go-live readiness.
