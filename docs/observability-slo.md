# Observability & SLOs

**Status:** Go/data-service instrumentation and an Unreal dedicated-server aggregate source baseline are implemented; production SLOs are not yet measured.

## Telemetry principles

Every durable mutation should be traceable across:
- client/session context,
- gameplay server,
- Go service,
- database transaction,
- resulting audit/receipt.

Never log or export raw secrets, resume keys, session tokens, gameplay tickets, shared keys, peer addresses, or player-controlled identifiers as metric labels.

## Implemented service/data baselines

The Go service exports a low-cardinality Prometheus text surface from a separate internal listener:
- request totals by bounded route scope + status class,
- request-duration histogram by bounded route scope,
- current in-flight request gauge,
- bounded PostgreSQL pgx pool metrics,
- bounded current-database PostgreSQL server metrics,
- bounded Redis server INFO metrics,
- fixed labels that do not contain quest, vehicle, race, account, token, operation, peer, database URL, Redis address, or credential identifiers.

`METRICS_LISTEN_ADDR` controls the listener. Docker Compose binds the host side to loopback only by default (`127.0.0.1:${GAME_API_METRICS_PORT:-19090}`). This is source/configuration evidence, not proof of deployed firewall or NetworkPolicy isolation.

See:
- [Runtime Observability Metrics v1.3](./runtime-observability-metrics-v1.3.md)
- [Runtime PostgreSQL Observability v1.5](./runtime-postgres-observability-v1.5.md)
- [Runtime Redis Server Observability v1.6](./runtime-redis-observability-v1.6.md)
- [Runtime PostgreSQL Server Observability v1.7](./runtime-postgres-server-observability-v1.7.md)

## Implemented Unreal source baseline — v2.0

The Unreal dedicated-server source records bounded aggregate telemetry for:
- active sessions,
- player joins and leaves,
- latest and peak server tick duration,
- authoritative driving-input clamp events,
- gameplay-ticket redemption attempts, successes, and failures.

The server emits one aggregate structured log record every 10 seconds using a fixed metric prefix and numeric fields only. No player, vehicle, race, ticket, credential, endpoint, or peer identifier is retained in the telemetry accumulator.

This is source/static-validation evidence only until a real UE 5.8 dedicated-server build/package run succeeds and emits retained runtime telemetry.

See [Runtime Unreal Dedicated-Server Observability v2.0](./runtime-unreal-server-observability-v2.0.md).

## Required telemetry still open

### Go service
- [x] request count,
- [x] status class,
- [x] latency histogram,
- [ ] database query/transaction latency,
- [x] connection-pool usage,
- [x] mutation conflict/race rejection baseline through HTTP status and security telemetry,
- [x] game-server auth rejection telemetry baseline,
- [ ] ticket issue/redeem/reuse domain counters on the Go service side.

### Unreal server
- [x] active sessions source aggregate,
- [x] server tick/frame-time source aggregate,
- [x] player joins/leaves source aggregate,
- [ ] reconnect-specific counter,
- [ ] replication rate/bytes,
- [x] authoritative input-clamp source aggregate,
- [ ] race checkpoint/result validation failures from live Unreal race transport,
- [ ] authority corrections from final vehicle physics,
- [ ] live packaged server emission/scrape/log-shipping evidence.

### Data services
- [x] PostgreSQL connection/pool metrics,
- [ ] PostgreSQL query/transaction latency and lock/wait time,
- [x] current-database PostgreSQL size metric,
- [ ] backup age metric from production backup custody,
- [ ] PostgreSQL replication/failover metrics when introduced,
- [x] Redis memory/eviction baseline,
- [ ] Redis command latency and deployed HA/failover metrics.

## Initial SLO targets

For an online alpha target:
- service availability target: 99.5% monthly, excluding declared maintenance,
- durable API p95 latency target: < 250 ms in-region,
- state-read p95 target: < 150 ms in-region,
- successful reconnect target: ≥ 99% under tested non-outage scenarios,
- duplicate durable reward target: 0 accepted duplicates in tested retry/failover cases.

These are planning targets until an environment is deployed and measured.

## Alerting principles

Page/urgent alert only for actionable symptoms such as:
- sustained service unavailability,
- severe latency/error-rate breach,
- database exhaustion,
- backup failure/age breach,
- widespread gameplay-server join failure,
- confirmed economy corruption/security incident.

Avoid alerting on noisy single requests.

## Evidence

SLO readiness requires deployed metrics/log shipping, dashboards, alert rules, retained load/soak or incident evidence, and a review of false-positive/false-negative behavior. Source/static instrumentation alone does not close that gate.
