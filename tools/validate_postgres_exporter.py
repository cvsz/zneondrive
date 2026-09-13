from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
compose = (ROOT / "compose.yaml").read_text(encoding="utf-8")
env_example = (ROOT / ".env.example").read_text(encoding="utf-8")
runtime_workflow = (ROOT / ".github/workflows/runtime-go.yml").read_text(encoding="utf-8")

required_compose = (
    "postgres-exporter:",
    "quay.io/prometheuscommunity/postgres-exporter:v0.19.1",
    "DATA_SOURCE_URI: postgres:5432/${POSTGRES_DB:-zneondrive}?sslmode=disable",
    "DATA_SOURCE_USER: ${POSTGRES_USER:-zneondrive}",
    "DATA_SOURCE_PASS: ${POSTGRES_PASSWORD:-zneondrive-local}",
    "PG_EXPORTER_COLLECTION_TIMEOUT: 5s",
    '127.0.0.1:${POSTGRES_EXPORTER_PORT:-19187}:9187',
)
for marker in required_compose:
    assert marker in compose, f"missing postgres-exporter compose marker: {marker}"

assert "POSTGRES_EXPORTER_PORT=19187" in env_example
assert '"${POSTGRES_EXPORTER_PORT:-19187}:9187"' not in compose, "exporter must not publish on all host interfaces"
assert "DATA_SOURCE_NAME: postgres://" not in compose, "do not embed credential-bearing DSN in exporter config"
assert "latest" not in compose.split("postgres-exporter:", 1)[1].split("\n\n", 1)[0], "exporter image must be version pinned"

required_ci = (
    "External PostgreSQL exporter integration",
    "quay.io/prometheuscommunity/postgres-exporter:v0.19.1",
    "PG_EXPORTER_COLLECTION_TIMEOUT=5s",
    "http://127.0.0.1:9187/metrics",
    "pg_up",
    "zneondrive-test",
)
for marker in required_ci:
    assert marker in runtime_workflow, f"missing postgres-exporter CI evidence marker: {marker}"

print("postgres-exporter evidence boundary valid")
