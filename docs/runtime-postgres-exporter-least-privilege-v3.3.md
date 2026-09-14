# Runtime PostgreSQL Exporter Least-Privilege Evidence v3.3

## Scope

This increment removes reuse of the application/admin PostgreSQL credential from the external `postgres_exporter` path and introduces a dedicated `zneondrive_monitor` login with only the built-in `pg_monitor` role plus explicit negative role attributes.

The selected architecture is unchanged: Unreal Engine 5.8 dedicated servers remain gameplay-authoritative, Go 1.27 remains the authenticated service plane, PostgreSQL remains durable authority, and Redis remains ephemeral coordination.

## Implementation

- Local Compose initializes `zneondrive_monitor` through a read-only `/docker-entrypoint-initdb.d` script on new PostgreSQL data volumes.
- The monitoring login is explicitly `NOSUPERUSER NOCREATEDB NOCREATEROLE NOREPLICATION NOBYPASSRLS`.
- The exporter uses separate `POSTGRES_EXPORTER_USER` / `POSTGRES_EXPORTER_PASSWORD` configuration and remains loopback-only on the host.
- The exporter image remains pinned to `quay.io/prometheuscommunity/postgres-exporter:v0.19.1` with a five-second collection timeout.

## Runtime CI evidence

Runtime Go CI creates the monitoring login against PostgreSQL 17, asserts all privileged role flags are false, asserts membership in `pg_monitor`, starts the pinned exporter with only the monitoring credential, requires `pg_up = 1` and `pg_stat_database_*`, and rejects leakage of either the application/admin or monitoring password in the metrics payload.

## Security boundary

The monitoring principal is read-only observability infrastructure. It cannot create databases or roles, bypass row-level security, replicate, or gain superuser authority. The exporter cannot accept gameplay mutations, issue/redeem gameplay tickets, mutate progression, or determine authoritative race outcomes.

## Operational note

PostgreSQL entrypoint initialization scripts only run when a data directory is first initialized. Existing local volumes must provision/rotate the monitoring role explicitly before switching exporter credentials; operators must not assume changing Compose alone retrofits an existing database.

## Non-claims

This evidence does **not** close production observability or production-readiness gates. Still open are `pg_stat_statements`/bounded statement-level policy, production secret-manager delivery and rotation, deployed Prometheus scrape evidence, dashboards/alerts/SLO measurement, deployment-scale cardinality/load/soak evidence, PostgreSQL/Redis HA/failover, regional DR, and real UE 5.8 build/package/live Unreal↔Go evidence.
