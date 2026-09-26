#!/usr/bin/env bash
set -Eeuo pipefail
umask 077

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
RESOURCE_LABEL="$DOCKER_LABEL=$RUN_ID"
GO_TEST_TIMEOUT="${GO_INTEGRATION_TIMEOUT:-15m}"
RAW_DIR=""

fail() { printf 'ERROR: %s\n' "$*" >&2; exit 1; }

mkdir -p "$ARTIFACT_ROOT"
if ! mkdir "$EVIDENCE_DIR"; then
  fail "Evidence run directory already exists or cannot be created; refusing to overwrite it."
fi
RAW_DIR="$(mktemp -d "${TMPDIR:-/tmp}/zneondrive-go-integration.XXXXXXXX")" || fail "Could not create a private temporary directory for raw logs."
printf 'run_id=%s\nstarted_utc=%s\npostgres_image=postgres:17-alpine\nredis_image=redis:8-alpine\n' "$RUN_ID" "$(date -u +%Y-%m-%dT%H:%M:%SZ)" > "$EVIDENCE_DIR/run.txt"

sanitize_file() {
  local input="$1" output="$2"
  sed -E \
    -e 's#(://)[^:/[:space:]@]+:[^@/[:space:]]+@#\1[REDACTED]@#gI' \
    -e 's/("[[:alnum:]_.-]*(password|passwd|token|secret|authorization|api[_-]?key|access[_-]?key|private[_-]?key|server[_-]?key|resume[_-]?key|credential|jwt|cookie)[[:alnum:]_.-]*"[[:space:]]*:[[:space:]]*).*/\1[REDACTED]/Ig' \
    -e 's/((password|passwd|token|secret|authorization|api[_-]?key|access[_-]?key|private[_-]?key|server[_-]?key|resume[_-]?key|credential|jwt|cookie)[[:alnum:]_.-]*[[:space:]]*[:=][[:space:]]*).*/\1[REDACTED]/Ig' \
    -e 's/(Bearer[[:space:]]+)[A-Za-z0-9._~+\/-]+=*/\1[REDACTED]/Ig' \
    -e 's/gh[pousr]_[A-Za-z0-9_]{20,}/[REDACTED]/g' \
    -e 's/github_pat_[A-Za-z0-9_]{20,}/[REDACTED]/g' \
    "$input" > "$output"
}

verify_sanitizer() {
  local probe="$RAW_DIR/sanitize-probe.raw.log" sanitized="$RAW_DIR/sanitize-probe.log"
  cat > "$probe" <<'PROBE'
Authorization: Bearer bearer-example
resume_key=resume-example
server_key: server-example
{"private_key":"private-example","ok":true}
postgres://runner:password-example@127.0.0.1/database
Bearer standalone-example
ghp_ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890
github_pat_ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890
PROBE
  sanitize_file "$probe" "$sanitized" || fail "Could not run the log sanitizer self-check."
  for secret in bearer-example resume-example server-example private-example password-example standalone-example ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890; do
    if grep -Fq "$secret" "$sanitized"; then fail "The log sanitizer self-check found an unredacted credential pattern."; fi
  done
  rm -f -- "$probe" "$sanitized"
}

docker_call() {
  local duration="$1"
  shift
  timeout --signal=TERM --kill-after=5s "$duration" docker "$@"
}

labeled_resource_ids() {
  local resource_type="$1"
  case "$resource_type" in
    containers) docker_call 10s ps -aq --filter "label=$RESOURCE_LABEL" ;;
    networks) docker_call 10s network ls -q --filter "label=$RESOURCE_LABEL" ;;
    volumes) docker_call 10s volume ls -q --filter "label=$RESOURCE_LABEL" ;;
    *) return 2 ;;
  esac
}

capture_labeled_container_logs() {
  local ids container_id raw="$RAW_DIR/container.raw.log"
  if ! ids="$(labeled_resource_ids containers)"; then
    printf 'Could not enumerate labeled integration containers for failure logs.\n' >&2
    return 1
  fi
  [[ -n "$ids" ]] || return 0
  while IFS= read -r container_id; do
    [[ -n "$container_id" ]] || continue
    if docker_call 10s logs "$container_id" > "$raw" 2>&1 || [[ -s "$raw" ]]; then
      if ! sanitize_file "$raw" "$EVIDENCE_DIR/container-${container_id:0:12}.log"; then
        printf 'Could not sanitize logs for container %s.\n' "$container_id" >&2
        rm -f -- "$raw"
        return 1
      fi
    fi
    rm -f -- "$raw"
  done <<< "$ids"
}

remove_labeled_resources() {
  local resource_type="$1" ids output
  if ! ids="$(labeled_resource_ids "$resource_type")"; then
    printf 'Could not enumerate labeled integration %s during cleanup.\n' "$resource_type" >&2
    return 1
  fi
  [[ -n "$ids" ]] || return 0
  mapfile -t resource_ids <<< "$ids"
  case "$resource_type" in
    containers)
      if docker_call 30s rm -f -v "${resource_ids[@]}" >/dev/null 2>&1; then return 0; fi
      ;;
    networks)
      if docker_call 30s network rm "${resource_ids[@]}" >/dev/null 2>&1; then return 0; fi
      ;;
    volumes)
      if docker_call 30s volume rm -f "${resource_ids[@]}" >/dev/null 2>&1; then return 0; fi
      ;;
  esac
  if output="$(labeled_resource_ids "$resource_type")" && [[ -z "$output" ]]; then
    return 0
  fi
  printf 'Failed to remove labeled integration %s: %s\n' "$resource_type" "${resource_ids[*]}" >&2
  return 1
}

cleanup() {
  local exit_code="$?" cleanup_failed=0
  trap - EXIT INT TERM
  set +e
  if [[ "$exit_code" -ne 0 ]]; then
    if ! capture_labeled_container_logs; then cleanup_failed=1; fi
    if [[ -f "$RAW_DIR/go-test.raw.log" ]]; then
      if ! sanitize_file "$RAW_DIR/go-test.raw.log" "$EVIDENCE_DIR/go-test.log"; then
        printf 'Could not sanitize the integration test log.\n' >&2
        cleanup_failed=1
      fi
    fi
  fi
  if ! remove_labeled_resources containers; then cleanup_failed=1; fi
  if ! remove_labeled_resources networks; then cleanup_failed=1; fi
  if ! remove_labeled_resources volumes; then cleanup_failed=1; fi
  if [[ "$cleanup_failed" -eq 1 && "$exit_code" -eq 0 ]]; then exit_code=1; fi
  printf 'cleanup_exit_code=%s\nfinished_utc=%s\nexit_code=%s\n' "$cleanup_failed" "$(date -u +%Y-%m-%dT%H:%M:%SZ)" "$exit_code" >> "$EVIDENCE_DIR/run.txt"
  if [[ -n "$RAW_DIR" ]]; then rm -rf -- "$RAW_DIR"; fi
  if [[ "$exit_code" -eq 0 ]]; then
    printf 'GO_INTEGRATION_EXIT_CODE=0\nEVIDENCE_DIR=%s\n' "$EVIDENCE_DIR"
  else
    printf 'GO_INTEGRATION_EXIT_CODE=%s\nEVIDENCE_DIR=%s\n' "$exit_code" "$EVIDENCE_DIR" >&2
  fi
  exit "$exit_code"
}
trap cleanup EXIT
trap 'exit 130' INT
trap 'exit 143' TERM

[[ "$GO_TEST_TIMEOUT" =~ ^[1-9][0-9]*(s|m|h)$ ]] || fail "GO_INTEGRATION_TIMEOUT must be a positive duration such as 15m."
command -v docker >/dev/null 2>&1 || fail "docker is required for isolated Go integration tests."
command -v go >/dev/null 2>&1 || fail "go is required for isolated Go integration tests."
command -v timeout >/dev/null 2>&1 || fail "timeout is required for bounded Docker startup and cleanup."
command -v mktemp >/dev/null 2>&1 || fail "mktemp is required to keep raw logs outside the evidence directory."
command -v grep >/dev/null 2>&1 || fail "grep is required for the log sanitizer self-check."
verify_sanitizer
docker_call 10s info >/dev/null 2>&1 || fail "Docker daemon is unavailable."

ensure_name_available() {
  local name="$1" status
  if docker_call 10s inspect "$name" >/dev/null 2>&1; then
    fail "Generated integration container name is already in use; no existing container was changed."
  else
    status="$?"
    [[ "$status" -eq 1 ]] || fail "Could not safely inspect generated integration container name."
  fi
}
ensure_name_available "$PG_NAME"
ensure_name_available "$REDIS_NAME"

docker_call 180s run -d --name "$PG_NAME" --label "$RESOURCE_LABEL" -p 127.0.0.1::5432 -e "POSTGRES_DB=$PG_DATABASE" -e "POSTGRES_USER=$PG_USER" -e "POSTGRES_PASSWORD=$PG_PASSWORD" postgres:17-alpine >/dev/null
docker_call 180s run -d --name "$REDIS_NAME" --label "$RESOURCE_LABEL" -p 127.0.0.1::6379 redis:8-alpine redis-server --save '' --appendonly no >/dev/null

pg_port="$(docker_call 10s port "$PG_NAME" 5432/tcp | sed -nE 's/.*:([0-9]+)$/\1/p' | head -n1)"
redis_port="$(docker_call 10s port "$REDIS_NAME" 6379/tcp | sed -nE 's/.*:([0-9]+)$/\1/p' | head -n1)"
[[ "$pg_port" =~ ^[0-9]+$ ]] || fail "Could not resolve the isolated PostgreSQL loopback port."
[[ "$redis_port" =~ ^[0-9]+$ ]] || fail "Could not resolve the isolated Redis loopback port."

pg_ready=0
for ((attempt = 0; attempt < 30; attempt++)); do
  if docker_call 3s exec -e "PGPASSWORD=$PG_PASSWORD" "$PG_NAME" psql -At -U "$PG_USER" -d "$PG_DATABASE" -c 'SELECT 1' 2>/dev/null | grep -qx '1'; then
    pg_ready=1
    break
  fi
  sleep 1
done
[[ "$pg_ready" -eq 1 ]] || fail "Isolated PostgreSQL did not pass SELECT 1 within 30 bounded attempts."

redis_ready=0
for ((attempt = 0; attempt < 30; attempt++)); do
  if [[ "$(docker_call 3s exec "$REDIS_NAME" redis-cli PING 2>/dev/null || true)" == "PONG" ]]; then
    redis_ready=1
    break
  fi
  sleep 1
done
[[ "$redis_ready" -eq 1 ]] || fail "Isolated Redis did not pass PING within 30 bounded attempts."

wait_host_port() {
  local port="$1" service="$2"
  for ((attempt = 0; attempt < 30; attempt++)); do
    if timeout --signal=TERM --kill-after=2s 2s bash -c 'exec 3<>"/dev/tcp/127.0.0.1/$1"' _ "$port" 2>/dev/null; then return 0; fi
    sleep 1
  done
  fail "The mapped $service host port did not accept connections within 30 bounded attempts."
}
wait_host_port "$pg_port" PostgreSQL
wait_host_port "$redis_port" Redis

case "${GO_TEST_TIMEOUT: -1}" in
  s) test_timeout_seconds="${GO_TEST_TIMEOUT%s}" ;;
  m) test_timeout_seconds="$(( ${GO_TEST_TIMEOUT%m} * 60 ))" ;;
  h) test_timeout_seconds="$(( ${GO_TEST_TIMEOUT%h} * 3600 ))" ;;
esac
hard_timeout="$((test_timeout_seconds + 30))s"
database_url="postgres://$PG_USER:$PG_PASSWORD@127.0.0.1:$pg_port/$PG_DATABASE?sslmode=disable"
printf 'postgres_host=127.0.0.1\npostgres_port=%s\nredis_host=127.0.0.1\nredis_port=%s\ngo_test_timeout=%s\ngo_test_package_parallelism=1\n' "$pg_port" "$redis_port" "$GO_TEST_TIMEOUT" >> "$EVIDENCE_DIR/run.txt"

set +e
(
  cd "$ROOT/services/game-api"
  TEST_DATABASE_URL="$database_url" \
    TEST_REDIS_ADDR="127.0.0.1:$redis_port" \
    ZNEONDRIVE_INTEGRATION_RESOURCE_LABEL="$RESOURCE_LABEL" \
    timeout --signal=TERM --kill-after=10s "$hard_timeout" go test -p 1 -timeout "$GO_TEST_TIMEOUT" -tags=integration ./...
) > "$RAW_DIR/go-test.raw.log" 2>&1
test_exit_code="$?"
set -e
printf 'go_test_exit_code=%s\n' "$test_exit_code" >> "$EVIDENCE_DIR/run.txt"
if ! sanitize_file "$RAW_DIR/go-test.raw.log" "$EVIDENCE_DIR/go-test.log"; then
  fail "Could not sanitize the integration test log."
fi
if [[ "$test_exit_code" -ne 0 ]]; then
  tail -n 80 "$EVIDENCE_DIR/go-test.log" >&2
  exit "$test_exit_code"
fi
tail -n 30 "$EVIDENCE_DIR/go-test.log"
