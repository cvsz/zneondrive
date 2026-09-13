# Runtime PostgreSQL Observability v1.5

## Scope

This increment extends the existing private Prometheus-format metrics endpoint with bounded PostgreSQL connection-pool telemetry from the Go service plane.

It is **source/unit observability evidence only**. It does not claim deployed dashboards, alerting, SLO compliance, query-level tracing, production capacity validation, HA, DR, or production readiness.

## Implemented metrics

The metrics listener now exports a credential-free snapshot of the pgx connection pool:

- `zneondrive_postgres_pool_connections{state="max"}`
- `zneondrive_postgres_pool_connections{state="total"}`
- `zneondrive_postgres_pool_connections{state="idle"}`
- `zneondrive_postgres_pool_connections{state="acquired"}`
- `zneondrive_postgres_pool_connections{state="constructing"}`
- `zneondrive_postgres_pool_acquires_total`
- `zneondrive_postgres_pool_empty_acquires_total`
- `zneondrive_postgres_pool_canceled_acquires_total`
- `zneondrive_postgres_pool_new_connections_total`
- `zneondrive_postgres_pool_acquire_duration_seconds`

The labels are fixed and bounded. The pool snapshot does not include SQL text, connection strings, database names, player IDs, account IDs, session credentials, vehicle IDs, race IDs, peer addresses, or other player-controlled values.

## Trust boundary

PostgreSQL remains the durable source of truth. Metrics are read-only observations and cannot modify admission, persistence, quest progression, inventory, vehicle builds, or race acceptance.

The metrics remain on the separate `METRICS_LISTEN_ADDR` listener. Docker Compose already publishes that listener only on host loopback by default. This change does not make the metrics endpoint player-facing.

## Automated evidence

Unit coverage verifies that:

- the expected bounded pool gauges/counters are rendered;
- acquire duration is emitted in seconds;
- credential/database URL material is not part of the pool metric surface;
- the existing HTTP metric route-cardinality and credential non-disclosure tests remain intact.

Runtime Go CI remains responsible for compiling the pgx-backed `PostgresPoolStats` adapter and running Go unit/integration tests.

## Explicit non-claims

This increment does **not** prove:

- deployed Prometheus scraping or dashboard/alert correctness;
- query latency/error/cardinality instrumentation;
- PostgreSQL server exporter metrics;
- Redis server metrics;
- Unreal dedicated-server tick/replication metrics;
- production pool sizing or saturation thresholds;
- deployment-scale load or long-duration soak;
- SLO attainment;
- production HA/DR or production readiness.

Those remain evidence-gated in `ROADMAP.md` and `IMPLEMENTATION-CHECKLIST.md`.
