# Runtime PostgreSQL Statement Metrics Evidence v3.4

## Status

**Implemented / CI-exercised partial evidence.** This increment does not prove deployed production observability, capacity, SLO compliance, or production readiness.

## Scope

zNeonDrive keeps PostgreSQL as the durable source of truth and adds bounded statement-level telemetry through PostgreSQL 17 `pg_stat_statements` and the existing version-pinned Prometheus Community `postgres_exporter`.

The local evidence stack configures PostgreSQL with:

- `shared_preload_libraries=pg_stat_statements`
- `compute_query_id=on`
- `pg_stat_statements.max=500`
- `pg_stat_statements.track=top`
- `CREATE EXTENSION IF NOT EXISTS pg_stat_statements`

The exporter enables only the built-in statement collector and caps it to 25 statements per scrape:

- `--collector.stat_statements`
- `--collector.stat_statements.limit=25`

`--collector.stat_statements.include_query` is intentionally not enabled, so SQL/query text is not exported. The collector exposes query IDs and fixed database/user labels for statement statistics; the local PostgreSQL statistics table is itself bounded to 500 normalized statement entries and the exporter emits at most 25 statement rows per scrape.

## Trust and privacy boundary

The monitoring principal remains `zneondrive_monitor`, a dedicated login with `pg_monitor` membership and explicit `NOSUPERUSER NOCREATEDB NOCREATEROLE NOREPLICATION NOBYPASSRLS` restrictions. It is observational only and cannot alter gameplay authority, progression, race acceptance, ticket redemption, Redis limiter decisions, or durable mutations.

The CI exercise uses a distinctive SQL value marker and fails if that marker, the application database password, or the monitoring password appears in the scraped metrics payload. This verifies query-text/credential non-disclosure for the exercised collector configuration.

## CI evidence

Runtime Go CI starts an isolated PostgreSQL 17 container with `pg_stat_statements` preloaded, creates the extension and least-privilege monitoring role, generates representative statement activity, and then starts `postgres_exporter:v0.19.1` with the statement collector enabled.

The gate requires:

- PostgreSQL startup with the expected preload configuration;
- `pg_stat_statements.max = 500`;
- at least one recorded statement;
- `pg_monitor` membership with all privileged role flags false;
- `pg_up = 1`;
- normal `pg_stat_database_*` metrics;
- `pg_stat_statements_calls_total` metrics;
- no more than 25 statement call series in the scrape;
- no SQL marker or database credential leakage.

`tools/validate_postgres_exporter.py` statically pins those controls so the statement collector cannot silently become unbounded or start exporting query text.

## Non-claims / remaining gates

This does **not** prove:

- deployed Prometheus scraping or network isolation;
- production secret-manager delivery or credential rotation;
- production workload cardinality behavior over long retention periods;
- dashboards, alert rules, or measured SLOs;
- deployment-scale load or long-duration soak;
- PostgreSQL HA/failover behavior or exporter continuity during failover;
- production backup/PITR/RPO/RTO or regional DR;
- any Unreal Engine 5.8 build/package, live Unreal↔Go, Garage 17, MQ001–MQ012, or authoritative live-race gate.

The repository therefore remains evidence-gated and not production-ready.
