# Runtime PostgreSQL External Exporter Evidence v3.2

## Scope

This increment adds a version-pinned `prometheus-community/postgres_exporter` sidecar-style service to the local Docker Compose stack and exercises the same exporter image against PostgreSQL 17 in Runtime Go CI.

The selected runtime architecture is unchanged: Unreal Engine 5.8 dedicated servers remain gameplay-authoritative, Go 1.27 remains the authenticated service plane, PostgreSQL remains durable authority, and Redis remains ephemeral coordination.

## Implementation

- Image: `quay.io/prometheuscommunity/postgres-exporter:v0.19.1`.
- Local exporter endpoint: `127.0.0.1:${POSTGRES_EXPORTER_PORT:-19187}:9187`.
- Collection timeout: 5 seconds.
- The Compose endpoint is loopback-only and is not a player-facing or public service.
- The exporter receives host/database separately from user/password instead of embedding credentials in a URI.
- CI starts the pinned exporter against the PostgreSQL 17 service, waits for `/metrics`, requires `pg_up = 1`, requires PostgreSQL database metrics, and fails if the CI database password appears in the metrics payload.
- `tools/validate_postgres_exporter.py` statically enforces the pinned image, loopback publication, bounded timeout, non-credential URI form, and CI evidence markers.

The upstream exporter supports PostgreSQL 17 and recommends `DATA_SOURCE_PASS_FILE` for production secret handling. The Compose configuration here is a local-development evidence surface that intentionally reuses local example credentials; staging/production deployment must use a dedicated monitoring principal and file/secret-backed credentials.

## Evidence boundary

This closes a source/CI gap for an **external PostgreSQL exporter process**, but it does not complete the production observability gate.

Still open:

- dedicated production `pg_monitor`/`pg_read_all_stats` principal and secret-file delivery;
- `pg_stat_statements` activation and bounded statement-level policy;
- deployed Prometheus scrape evidence;
- dashboards, alerts and SLO measurement;
- production sizing/cardinality review;
- deployment-scale load and long soak;
- PostgreSQL HA/failover exporter verification.

Accordingly, `PostgreSQL query-level/external exporter metrics` remains open in the production checklist until deployed and statement-level evidence exists.

## Security / trust boundary

The exporter is read-only observability infrastructure. It is not allowed to authorize player requests, mutate durable gameplay state, redeem gameplay tickets, accept race checkpoints/results, alter Redis rate-limit decisions, or override Unreal dedicated-server authority.

Production deployments must not expose port 9187 publicly and must not reuse application-owner credentials when a least-privilege monitoring role can be provisioned.
