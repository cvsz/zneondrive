# Runtime PostgreSQL Server Observability v1.7

## Scope

This increment adds a bounded, read-only PostgreSQL server snapshot to the existing private metrics surface. It preserves PostgreSQL as durable service authority, Redis as ephemeral coordination, and the Unreal dedicated server as gameplay authority.

The Go service reads only numeric values for the current application database from `pg_stat_database` plus `pg_database_size(current_database())`.

## Exported metrics

The private metrics listener now exposes:

- `zneondrive_postgres_up`;
- current backend count;
- committed and rolled-back transaction totals;
- block reads and cache hits;
- tuple returned/fetched/inserted/updated/deleted totals;
- deadlock count;
- temporary file and byte totals;
- current database size in bytes.

All labels are fixed enums. Database names, SQL text, query fingerprints, relation names, connection URLs, credentials, player identifiers, vehicle identifiers, race identifiers and raw query errors are not exported.

If the bounded stats query fails, the metrics surface reports only `zneondrive_postgres_up 0` for this probe and does not render the raw error.

## Evidence

Automated evidence includes:

- unit coverage for successful bounded Prometheus output;
- unit coverage proving raw PostgreSQL error/address/database material is not exposed on failure;
- PostgreSQL 17 integration coverage using `TEST_DATABASE_URL` that requires a live backend count, positive database size and non-negative cumulative counters;
- existing Runtime Go CI module-lock, unit, PostgreSQL/Redis integration, restore-drill, vet and build gates.

## Trust boundary

The stats query is observational and read-only. It cannot authorize a request, alter progression, mutate race state, change rate-limit decisions or override dedicated-server gameplay acceptance.

## Explicit non-claims

This closes only the bounded PostgreSQL server-metrics source/integration baseline. It does **not** prove query-level latency tracing, statement/fingerprint metrics, an external PostgreSQL exporter, deployed dashboards/alerts, production sizing, SLO attainment, HA/failover, PITR/RPO/RTO, deployment-scale load/soak, Unreal runtime metrics or production readiness.
