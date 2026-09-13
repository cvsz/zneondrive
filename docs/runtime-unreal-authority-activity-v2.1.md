# Runtime Unreal Authority Activity Observability v2.1

**Status:** implemented source/static-validation evidence only. A successful real UE 5.8 build/package and retained live runtime logs are still required before claiming runtime verification.

## Scope

This increment extends the bounded Unreal dedicated-server telemetry accumulator with authority-side activity counters that can be correlated with future live Unreal↔Go and load evidence without retaining dynamic player identifiers.

The aggregate record now includes:

- `authority_movement_ticks_total` — authoritative vehicle movement ticks that passed the durable-identity gate;
- `identity_gate_blocks_total` — networked authoritative ticks rejected because durable identity was not yet bound;
- `collision_blocks_total` — swept authoritative movement attempts that ended in a blocking collision;
- `net_update_requests_total` — explicit `ForceNetUpdate()` requests issued when a durable vehicle snapshot is bound.

These counters join the v2.0 session, tick, input-clamp and gameplay-ticket redemption aggregates in the same fixed-format log record.

## Trust boundary

The counters are observational only. They do not authorize movement, modify durable state, weaken the durable-identity gate, accept client transforms, bypass input clamps, or alter gameplay-ticket redemption.

The existing authority model remains unchanged:

- Unreal dedicated server accepts and applies gameplay movement;
- Go enforces service/authentication contracts;
- PostgreSQL owns durable account/character/vehicle/build/race state;
- Redis remains ephemeral coordination only.

No player ID, vehicle ID, race ID, session token, gameplay ticket, shared key, endpoint, IP address or peer address is retained in the telemetry snapshot or emitted as a metric field.

## Source evidence

`ANDVehiclePawn` now records:

1. durable-identity gate blocks immediately before rejecting networked movement;
2. authoritative movement ticks only after the gate passes;
3. collision blocks from the actual `FHitResult::bBlockingHit` produced by swept server movement;
4. explicit replication update requests immediately after the real `ForceNetUpdate()` call used by durable identity binding.

`tools/validate_unreal_server_observability.py` requires all four hooks, requires collision telemetry to be derived from `Hit.bBlockingHit`, requires the replication-update counter to remain paired with an actual `ForceNetUpdate()` call, and continues rejecting dynamic/string identifiers in the aggregate metric line.

## Explicit non-claims

This increment is **not** replication-byte or replication-rate measurement. It does not claim NetDriver packet/byte telemetry, live race-validation telemetry, client/server correction-distance telemetry from final physics, deployed log shipping, dashboards, SLO attainment, load/soak, anti-cheat readiness, HA/DR, or production readiness.

Those gates remain open until a real UE 5.8 Client/Server build/package succeeds and live runtime evidence is retained.
