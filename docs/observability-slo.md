# Observability & SLOs

**Status:** Go HTTP instrumentation source implemented; production SLOs are not yet measured.

## Telemetry principles

Every durable mutation should be traceable across:
- client/session context,
- gameplay server,
- Go service,
- database transaction,
- resulting audit/receipt.

Never log or export raw secrets, resume keys, session tokens, gameplay tickets, shared keys, peer addresses, or player-controlled identifiers as metric labels.

## Implemented Go HTTP baseline — v1.3

The service plane now exports a low-cardinality Prometheus text surface from a separate internal listener:
- request totals by bounded route scope + status class,
- request-duration histogram by bounded route scope,
- current in-flight request gauge,
- fixed route labels that do not contain quest, vehicle, race, account, token, operation, or peer identifiers.

`METRICS_LISTEN_ADDR` controls the listener. Docker Compose binds the host side to loopback only by default (`127.0.0.1:${GAME_API_METRICS_PORT:-19090}`). This is source/configuration evidence, not proof of deployed firewall or NetworkPolicy isolation.

See [Runtime Observability Metrics v1.3](./runtime-observability-metrics-v1.3.md).

## Required telemetry still open

### Go service
- [x] request count,
- [x] status class,
- [x] latency histogram,
- [ ] database query/transaction latency,
- [ ] connection-pool usage,
- [x] mutation conflict/race rejection baseline through HTTP status and security telemetry,
- [x] game-server auth rejection telemetry baseline,
- [ ] ticket issue/redeem/reuse domain counters.

### Unreal server
- active sessions,
- server tick/frame time,
- player joins/leaves/reconnects,
- replication rate/bytes,
- input validation rejects,
- race checkpoint/result validation failures,
- authority corrections.

### Data services
- PostgreSQL connections,
- lock/wait time,
- storage growth,
- backup age,
- replication/failover metrics when introduced,
- Redis memory/eviction/latency.

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

SLO readiness requires deployed metrics, dashboards, alert rules, retained load/soak or incident evidence, and a review of false-positive/false-negative behavior. The v1.3 source metrics baseline alone does not close that gate.
