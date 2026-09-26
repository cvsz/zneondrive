#!/usr/bin/env bash
set -Eeuo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
RUN_ID="$(date -u +%Y%m%dT%H%M%SZ)-$$"
ARTIFACT_ROOT="${GO_INTEGRATION_ARTIFACT_DIR:-$ROOT/artifacts/go-integration}"
EVIDENCE_DIR="$ARTIFACT_ROOT/$RUN_ID"
PG_NAME="zneondrive-it-pg-$RUN_ID"
REDIS_NAME="zneondrive-it-redis-$RUN_ID"
PG_DATABASE="zneondrive_test"
PG_USER="zneondrive"
PG_PASSWORD="zneondrive-integration-test"
DOCKER_LABEL="com.zneondrive.integration-run"
GO_TEST_TIMEOUT="${GO_INTEGRATION_TIMEOUT:-15m}"

mkdir -p "$EVIDENCE_DIR"
printf 'run_id=%s\nstarted_utc=%s\npostgres_image=postgres:17-alpine\nredis_image=redis:8-alpine\n' "$RUN_ID" "$(date -u +%Y-%m-%dT%H:%M:%SZ)" > "$EVIDENCE_DIR/run.txt"

fail() { printf 'ERROR: %s\n' "$*" >&2; exit 1; }

sanitize_file() {
  local input="$1" output="$2"
  sed -E -e 's#(://[^:/[:space:]]+):[^@/[:space:]]+@#\1:[REDACTED]@#gI' -e 's/(password|passwd|token|secret|authorization|api[_-]?key)(=|:)[[:space:]]*[^[:space:],;]+/\1\2[REDACTED]/Ig' -e 's/(Bearer[[:space:]]+)[A-Za-z0-9._~+\/-]+=*/\1[REDACTED]/Ig' "$input" > "$output"
}

owns_container() {
  local name="$1" label_value
  label_value="$(docker inspect -f "{{ index .Config.Labels \"$DOCKER_LABEL\" }}" "$name" 2>/dev/null || true)"
  [[ "$label_value" == "$RUN_ID" ]]
}

capture_container_log() {
  local name="$1" output="$2" raw="$EVIDENCE_DIR/.container.raw.log"
  if owns_container "$name"; then
    docker logs "$name" > "$raw" 2>&1 || true
    sanitize_file "$raw" "$output"
    rm -f -- "$raw"
  fi
}

cleanup() {
  local exit_code="$?"
  trap - EXIT
  set +e
  if [[ "$exit_code" -ne 0 ]]; then
    capture_container_log "$PG_NAME" "$EVIDENCE_DIR/postgres.log"
    capture_container_log "$REDIS_NAME" "$EVIDENCE_DIR/redis.log"
    if [[ -f "$EVIDENCE_DIR/go-test.raw.log" ]]; then
      sanitize_file "$EVIDENCE_DIR/go-test.raw.log" "$EVIDENCE_DIR/go-test.log"
      rm -f -- "$EVIDENCE_DIR/go-test.raw.log"
    fi
  fi
  printf 'finished_utc=%s\nexit_code=%s\n' "$(date -u +%Y-%m-%dT%H:%M:%SZ)" "$exit_code" >> "$EVIDENCE_DIR/run.txt"
  if owns_container "$PG_NAME"; then docker rm -f -v "$PG_NAME" >/dev/null 2>&1; fi
  if owns_container "$REDIS_NAME"; then docker rm -f -v "$REDIS_NAME" >/dev/null 2>&1; fi
  if [[ "$exit_code" -eq 0 ]]; then
    printf 'GO_INTEGRATION_EXIT_CODE=0\nEVIDENCE_DIR=%s\n' "$EVIDENCE_DIR"
  else
    printf 'GO_INTEGRATION_EXIT_CODE=%s\nEVIDENCE_DIR=%s\n' "$exit_code" "$EVIDENCE_DIR" >&2
  fi
  exit "$exit_code"
}
trap cleanup EXIT

[[ "$GO_TEST_TIMEOUT" =~ ^[1-9][0-9]*(s|m|h)$ ]] || fail "GO_INTEGRATION_TIMEOUT must be a positive duration such as 15m."
command -v docker >/dev/null 2>&1 || fail "docker is required for isolated Go integration tests."
command -v go >/dev/null 2>&1 || fail "go is required for isolated Go integration tests."
command -v timeout >/dev/null 2>&1 || fail "timeout is required for bounded Docker startup."
timeout 10s docker info >/dev/null 2>&1 || fail "Docker daemon is unavailable."

if docker inspect "$PG_NAME" >/dev/null 2>&1 || docker inspect "$REDIS_NAME" >/dev/null 2>&1; then
  fail "Generated integration container name is already in use; no existing container was changed."
fi

timeout 180s docker run -d --name "$PG_NAME" --label "$DOCKER_LABEL=$RUN_ID" -p 127.0.0.1::5432 -e "POSTGRES_DB=$PG_DATABASE" -e "POSTGRES_USER=$PG_USER" -e "POSTGRES_PASSWORD=$PG_PASSWORD" postgres:17-alpine >/dev/null
timeout 180s docker run -d --name "$REDIS_NAME" --label "$DOCKER_LABEL=$RUN_ID" -p 127.0.0.1::6379 redis:8-alpine redis-server --save '' --appendonly no >/dev/null

pg_port="$(docker port "$PG_NAME" 5432/tcp | sed -nE 's/.*:([0-9]+)$/\1/p' | head -n1)"
redis_port="$(docker port "$REDIS_NAME" 6379/tcp | sed -nE 's/.*:([0-9]+)$/\1/p' | head -n1)"
[[ "$pg_port" =~ ^[0-9]+$ ]] || fail "Could not resolve the isolated PostgreSQL loopback port."
[[ "$redis_port" =~ ^[0-9]+$ ]] || fail "Could not resolve the isolated Redis loopback port."

pg_ready=0
for ((attempt = 0; attempt < 30; attempt++)); do
  if docker exec -e "PGPASSWORD=$PG_PASSWORD" "$PG_NAME" psql -At -U "$PG_USER" -d "$PG_DATABASE" -c 'SELECT 1' 2>/dev/null | grep -qx '1'; then
    pg_ready=1
    break
  fi
  sleep 1
done
[[ "$pg_ready" -eq 1 ]] || fail "Isolated PostgreSQL was not query-ready within 30 seconds."

redis_ready=0
for ((attempt = 0; attempt < 30; attempt++)); do
  if [[ "$(docker exec "$REDIS_NAME" redis-cli PING 2>/dev/null || true)" == "PONG" ]]; then
    redis_ready=1
    break
  fi
  sleep 1
done
[[ "$redis_ready" -eq 1 ]] || fail "Isolated Redis was not ready within 30 seconds."

database_url="postgres://$PG_USER:$PG_PASSWORD@127.0.0.1:$pg_port/$PG_DATABASE?sslmode=disable"
printf 'postgres_host=127.0.0.1\npostgres_port=%s\nredis_host=127.0.0.1\nredis_port=%s\ngo_test_timeout=%s\n' "$pg_port" "$redis_port" "$GO_TEST_TIMEOUT" >> "$EVIDENCE_DIR/run.txt"

set +e
(
  cd "$ROOT/services/game-api"
  TEST_DATABASE_URL="$database_url" TEST_REDIS_ADDR="127.0.0.1:$redis_port" timeout "$GO_TEST_TIMEOUT" go test -timeout "$GO_TEST_TIMEOUT" -tags=integration ./...
) > "$EVIDENCE_DIR/go-test.raw.log" 2>&1
test_exit_code="$?"
set -e
printf 'go_test_exit_code=%s\n' "$test_exit_code" >> "$EVIDENCE_DIR/run.txt"
sanitize_file "$EVIDENCE_DIR/go-test.raw.log" "$EVIDENCE_DIR/go-test.log"
rm -f -- "$EVIDENCE_DIR/go-test.raw.log"
if [[ "$test_exit_code" -ne 0 ]]; then
  tail -n 80 "$EVIDENCE_DIR/go-test.log" >&2
  exit "$test_exit_code"
fi
tail -n 30 "$EVIDENCE_DIR/go-test.log"
