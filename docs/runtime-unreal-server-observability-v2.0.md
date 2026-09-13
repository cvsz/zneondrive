# Runtime Unreal Dedicated-Server Observability v2.0

## Scope

This increment adds a bounded, aggregate telemetry baseline inside the Unreal Engine 5.8 dedicated gameplay server source. It is intentionally observational: it does not change gameplay authority, durable state ownership, ticket authorization, vehicle legality, race acceptance, or Redis/PostgreSQL responsibilities.

The authoritative gameplay server now records aggregate counters/timing for:

- active server sessions;
- player joins and leaves;
- latest and peak server tick duration in milliseconds;
- authoritative driving-input clamp events;
- gameplay-ticket redemption attempts, successes, and failures.

Every 10 seconds the GameMode emits one structured aggregate log record with the fixed prefix `metric=zneondrive_unreal_server`.

## Security and cardinality

The telemetry accumulator stores only integers and floating-point timing values. It does not retain or export:

- player/account/character identifiers;
- vehicle/build/race identifiers;
- gameplay tickets, session tokens, resume keys, or shared keys;
- service URLs, Redis/PostgreSQL addresses, peer addresses, or client IPs;
- arbitrary request-controlled metric labels.

`tools/validate_unreal_server_observability.py` statically verifies the required hooks and rejects string-valued/dynamic identifiers in the aggregate metric contract. CI requires this validator.

## Authority boundary

The existing trust model is unchanged:

- Unreal dedicated servers remain authoritative for gameplay acceptance;
- Go remains the durable service-plane API and authorization boundary;
- PostgreSQL remains the durable state authority;
- Redis remains ephemeral coordination for abuse controls only;
- clients cannot write telemetry values to authorize gameplay or mutate durable state.

Input-clamp telemetry observes the existing server-side `[-1, 1]` clamp and does not weaken it. Ticket telemetry observes the existing server-only redemption path and never records the ticket or server shared key.

## Evidence in this increment

Repository/CI evidence covers:

- source wiring in `ANeonDriveGameModeBase` for tick and player lifecycle aggregates;
- authoritative input-clamp accounting in `ANDVehiclePawn`;
- gameplay-ticket redemption aggregate accounting in `ANDPlayerController`;
- numeric-only aggregate snapshot/log format;
- static validation that required hooks exist and dynamic/secret identifiers are absent from the metric record.

## Explicit non-claims

This is **source/static-validation evidence only** until a real UE 5.8 Client/Server build executes successfully on the self-hosted Unreal runner.

It does not prove:

- successful UE 5.8 compilation or package/cook output;
- live telemetry emission from a packaged dedicated server;
- replication rate/bytes metrics;
- race checkpoint/result validation metrics from live Unreal race transport;
- authority-correction metrics derived from final vehicle physics;
- deployed log shipping, Prometheus scraping, dashboards, alerts, or SLO compliance;
- deployment-scale load/soak, HA/DR, or production readiness.

The next evidence step is a successful retained UE 5.8 Client/Server build/package run followed by live Unreal↔Go session/ticket/binding E2E with retained server telemetry.
