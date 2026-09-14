from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
compose = (ROOT / "compose.yaml").read_text(encoding="utf-8")
env_example = (ROOT / ".env.example").read_text(encoding="utf-8")
runtime_workflow = (ROOT / ".github/workflows/runtime-go.yml").read_text(encoding="utf-8")
init_monitoring = (ROOT / "deploy/postgres/init-monitoring.sh").read_text(encoding="utf-8")

required_compose = (
    "postgres-exporter:",
    "quay.io/prometheuscommunity/postgres-exporter:v0.19.1",
    "DATA_SOURCE_URI: postgres:5432/${POSTGRES_DB:-zneondrive}?sslmode=disable",
    "DATA_SOURCE_USER: ${POSTGRES_EXPORTER_USER:-zneondrive_monitor}",
    "DATA_SOURCE_PASS: ${POSTGRES_EXPORTER_PASSWORD:-zneondrive-monitor-local}",
    "PG_EXPORTER_COLLECTION_TIMEOUT: 5s",
    '127.0.0.1:${POSTGRES_EXPORTER_PORT:-19187}:9187',
    "./deploy/postgres/init-monitoring.sh:/docker-entrypoint-initdb.d/20-monitoring.sh:ro",
    "shared_preload_libraries=pg_stat_statements",
    "compute_query_id=on",
    "pg_stat_statements.max=500",
    "pg_stat_statements.track=top",
    "--collector.stat_statements",
    "--collector.stat_statements.limit=25",
)
for marker in required_compose:
    assert marker in compose, f"missing postgres-exporter compose marker: {marker}"

for marker in (
    "POSTGRES_EXPORTER_PORT=19187",
    "POSTGRES_EXPORTER_USER=zneondrive_monitor",
    "POSTGRES_EXPORTER_PASSWORD=zneondrive-monitor-local",
):
    assert marker in env_example, f"missing exporter env example marker: {marker}"

assert '"${POSTGRES_EXPORTER_PORT:-19187}:9187"' not in compose, "exporter must not publish on all host interfaces"
assert "DATA_SOURCE_NAME: postgres://" not in compose, "do not embed credential-bearing DSN in exporter config"
assert "DATA_SOURCE_USER: ${POSTGRES_USER" not in compose, "exporter must not reuse application/admin user"
assert "DATA_SOURCE_PASS: ${POSTGRES_PASSWORD" not in compose, "exporter must not reuse application/admin password"
assert "--collector.stat_statements.include_query" not in compose, "statement collector must not export SQL/query text"
assert "latest" not in compose.split("postgres-exporter:", 1)[1].split("\n\n", 1)[0], "exporter image must be version pinned"

required_init = (
    "CREATE EXTENSION IF NOT EXISTS pg_stat_statements",
    "CREATE ROLE zneondrive_monitor LOGIN",
    "NOSUPERUSER NOCREATEDB NOCREATEROLE NOREPLICATION NOBYPASSRLS",
    "GRANT pg_monitor TO zneondrive_monitor",
    "POSTGRES_EXPORTER_PASSWORD must be set",
)
for marker in required_init:
    assert marker in init_monitoring, f"missing least-privilege monitoring marker: {marker}"

required_ci = (
    "Bounded pg_stat_statements exporter integration",
    "shared_preload_libraries=pg_stat_statements",
    "compute_query_id=on",
    "pg_stat_statements.max=500",
    "pg_stat_statements.track=top",
    "CREATE EXTENSION IF NOT EXISTS pg_stat_statements",
    "DATA_SOURCE_USER='zneondrive_monitor'",
    "pg_has_role('zneondrive_monitor', 'pg_monitor', 'member')",
    "role_privileges_ok",
    "NOT rolsuper AND NOT rolcreatedb AND NOT rolcreaterole AND NOT rolreplication AND NOT rolbypassrls",
    "quay.io/prometheuscommunity/postgres-exporter:v0.19.1",
    "PG_EXPORTER_COLLECTION_TIMEOUT=5s",
    "--collector.stat_statements",
    "--collector.stat_statements.limit=25",
    "pg_stat_statements_calls_total",
    "statement_series",
    "query_marker",
    "exporter metrics leaked SQL/query text",
    "zneondrive-test",
    "zneondrive-monitor-test",
)
for marker in required_ci:
    assert marker in runtime_workflow, f"missing postgres-exporter CI evidence marker: {marker}"

assert "--collector.stat_statements.include_query" not in runtime_workflow, "CI must not enable SQL/query text export"

print("postgres-exporter bounded statement evidence boundary valid")
