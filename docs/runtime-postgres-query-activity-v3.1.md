# Runtime PostgreSQL Query Activity Evidence v3.1

## Scope

This increment adds bounded, read-only PostgreSQL query-activity pressure metrics to the existing private Go metrics listener. PostgreSQL remains the durable source of truth and Unreal Engine 5.8 dedicated servers remain gameplay-authoritative.

The service reads aggregate numeric values from `pg_stat_activity` for `current_database()` only. The fixed query never selects SQL text, query identifiers/fingerprints, user names, application names, client addresses, backend PIDs, credentials, account/vehicle/race identifiers, or other player-controlled values.

## Exported metrics

- `zneondrive_postgres_query_activity_up`
- `zneondrive_postgres_query_activity{state="active"}`
- `zneondrive_postgres_query_activity{state="waiting"}`
- `zneondrive_postgres_query_activity{state="idle_in_transaction"}`
- `zneondrive_postgres_query_activity{state="long_running"}` where long-running means active for more than five seconds
- `zneondrive_postgres_oldest_active_query_seconds`
- `zneondrive_postgres_oldest_transaction_seconds`

The metrics query excludes its own backend when counting active/waiting/age signals so scraping does not create a permanent false active-query signal.

## Evidence

Unit coverage verifies numeric output, bounded fixed labels, and non-disclosure of SQL/connection/user/client details. PostgreSQL integration coverage executes the query against the Runtime Go CI PostgreSQL service and requires all aggregate counts and ages to be non-negative.

Probe errors are represented only by the `*_up` gauge; raw PostgreSQL errors are never written into Prometheus output.

## Trust boundary

The query is observational and read-only. It cannot authorize requests, mutate PostgreSQL durable state, alter Redis rate-limit coordination, accept race checkpoints/results, or override Unreal dedicated-server gameplay authority.

## Explicit non-claims

This is **partial query-level observability**, not completion of the PostgreSQL query-level/external-exporter production gate. It does not provide SQL text/fingerprint metrics, per-statement latency histograms, `pg_stat_statements` coverage, an external PostgreSQL exporter, deployed scrape/dashboard/alert evidence, production sizing, SLO attainment, load/soak evidence, HA/failover, PITR/RPO/RTO, or production readiness.

The production checklist item `PostgreSQL query-level/external exporter metrics` therefore remains open.
