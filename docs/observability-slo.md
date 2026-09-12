# Observability & SLOs

**Status:** target operating contract; production SLOs are not yet measured.

## Telemetry principles

Every durable mutation should be traceable across:
- client/session context,
- gameplay server,
- Go service,
- database transaction,
- resulting audit/receipt.

Never log raw secrets, resume keys, session tokens or gameplay tickets.

## Required telemetry

### Go service
- request count,
- status/error code,
- latency histogram,
- database query/transaction latency,
- connection-pool usage,
- mutation conflict/idempotency rejects,
- ticket issue/redeem/reuse rejects.

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
- Redis memory/eviction/latency when introduced.

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

SLO readiness requires dashboards, alert rules, retained test/incident evidence and a review of false-positive/false-negative behavior.
