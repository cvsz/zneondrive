# Runtime Reference Rebuild Parity v2.7

## Purpose

This increment narrows the remaining Python-reference ↔ Go/PostgreSQL parity gap for the Garage 17 rebuild slice. It does not replace the production service plane and it does not close the full reference-oracle parity, playable Garage 17, live Unreal↔Go, or production-readiness gates.

## Evidence added

The dependency-free Python reference authority now models the same durable rebuild invariants already exercised by the Go/PostgreSQL runtime:

- MQ004 grants exactly one `part_brakes_track_i` on first completion only.
- MQ005 unlocks `bp_starter_rebuild` once and respects blueprint capacity.
- MQ009 grants exactly one `part_tires_street_i` on first completion only.
- Garage 17 rebuilds require the starter-rebuild blueprint.
- Newly equipped parts must exist in authoritative inventory before any mutation occurs.
- Newly equipped parts are consumed and removed parts are returned.
- Inventory accounting and immutable build-revision advancement occur only after all preconditions pass.
- Operation IDs remain payload-bound: a semantic replay returns the original receipt while reuse with changed parts is rejected.
- Canonical part ordering keeps equivalent reordered payloads idempotent.

Regression tests cover duplicate quest-completion operations, blueprint gating, insufficient inventory with no partial mutation, consume/return swap accounting, semantic replay, and operation-key conflict behavior.

## Trust boundary

The architecture remains unchanged:

- Unreal Engine 5.8 dedicated server: gameplay authority.
- Go 1.27: authenticated service/runtime contracts.
- PostgreSQL: durable source of truth.
- Redis: ephemeral coordination and distributed abuse-control state.
- Python: executable reference oracle only.

No player-facing client is allowed to assert inventory, blueprint, build, reward, or race authority through this reference model.

## Non-claims

This increment does **not** prove:

- successful Unreal Engine 5.8 Client or dedicated-server compilation;
- Client/Server cook or package artifacts;
- live packaged Unreal↔Go transport;
- playable Garage 17 interactions or First Ignition;
- complete MQ001–MQ012 Unreal gameplay;
- full Python↔Go reference parity;
- deployed load/soak, observability, HA/DR, or production readiness.

The top-level evidence status therefore remains Phase 4.19 / Runtime v2.3 / Evidence Gated with UE build evidence pending.
