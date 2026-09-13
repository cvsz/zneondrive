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
assert "latest" not in compose.split("postgres-exporter:", 1)[1].split("\n\n", 1)[0], "exporter image must be version pinned"

required_init = (
    "CREATE ROLE zneondrive_monitor LOGIN",
    "NOSUPERUSER NOCREATEDB NOCREATEROLE NOREPLICATION NOBYPASSRLS",
    "GRANT pg_monitor TO zneondrive_monitor",
    "POSTGRES_EXPORTER_PASSWORD must be set",
)
for marker in required_init:
    assert marker in init_monitoring, f"missing least-privilege monitoring marker: {marker}"

required_ci = (
    "Least-privilege PostgreSQL exporter integration",
    "DATA_SOURCE_USER='zneondrive_monitor'",
    "pg_has_role('zneondrive_monitor', 'pg_monitor', 'member')",
    "0:0:0:0:0",
    "quay.io/prometheuscommunity/postgres-exporter:v0.19.1",
    "PG_EXPORTER_COLLECTION_TIMEOUT=5s",
    "http://127.0.0.1:9187/metrics",
    "pg_up",
    "zneondrive-test",
    "zneondrive-monitor-test",
)
for marker in required_ci:
    assert marker in runtime_workflow, f"missing postgres-exporter CI evidence marker: {marker}"

print("postgres-exporter least-privilege evidence boundary valid")
