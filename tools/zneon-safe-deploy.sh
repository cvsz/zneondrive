#!/usr/bin/env bash
# Safe, non-destructive zNeonDrive web deployment and evidence collection.
set -Eeuo pipefail
cd "$(dirname "$0")/.."
PORT="${WEB_PORT:-18085}"
REPORT_DIR="${REPORT_DIR:-.runtime/deploy-evidence}"
mkdir -p "$REPORT_DIR"
REPORT="$REPORT_DIR/zneon-$(date -u +%Y%m%dT%H%M%SZ).log"
chmod 700 "$REPORT_DIR"
touch "$REPORT"; chmod 600 "$REPORT"
exec > >(tee -a "$REPORT") 2>&1
echo "=== zNeonDrive deployment evidence $(date -u -Is) ==="
echo "commit=$(git rev-parse HEAD)"
echo "branch=$(git branch --show-current)"
if ! git diff --quiet || ! git diff --cached --quiet; then
  echo "BLOCKED: tracked working tree changes. Preserve and review them before deployment." >&2
  exit 2
fi
if [[ ! -f .env ]]; then
  echo "BLOCKED: .env missing. Initialize secrets using make env-init, then inspect settings privately." >&2
  exit 2
fi
if ! grep -Eq '^GAME_SERVER_SHARED_KEY=.{16,}' .env; then
  echo "BLOCKED: GAME_SERVER_SHARED_KEY missing or too short (value not displayed)." >&2
  exit 2
fi
if grep -q '^WEB_PORT=' .env; then
  PORT="$(sed -n 's/^WEB_PORT=//p' .env | tail -1)"
fi
if [[ "$PORT" != 18085 ]]; then
  echo "BLOCKED: WEB_PORT must be 18085 for the reviewed Cloudflare tunnel origin; actual=$PORT" >&2
  exit 2
fi
command -v docker >/dev/null || { echo "BLOCKED: docker missing"; exit 2; }
docker compose config --quiet
if ss -ltn 2>/dev/null | grep -Eq '(:18085)[[:space:]]'; then
  echo "NOTE: 18085 already listening; checking whether existing Compose owns it."
  docker compose ps
fi
echo "=== build/start (preserve named volumes) ==="
docker compose up -d --build
echo "=== containers ==="
docker compose ps
echo "=== local health ==="
for url in "http://127.0.0.1:18085/healthz" "http://127.0.0.1:18085/api/healthz" "http://127.0.0.1:18085/client-v-0-0-1/manifest.json"; do
  echo "CHECK $url"
  curl --fail --silent --show-error --retry 12 --retry-connrefused --retry-delay 5 --max-time 10 "$url" >/dev/null
  echo "PASS $url"
done
echo "=== unpublished binary routes must fail closed ==="
for ext in zip apk exe bat; do
  code="$(curl --silent --output /dev/null --write-out '%{http_code}' "http://127.0.0.1:18085/client-v-0-0-1/zneondrive-client-v-0-0-1.$ext")"
  if [[ "$code" == 404 ]]; then echo "PENDING $ext (404)"; else echo "REVIEW $ext HTTP=$code (check manifest and actual SHA-256 before publishing)"; fi
done
echo "=== public probe (may require Cloudflare apply/reload) ==="
for url in "https://zneon.zeaz.dev/healthz" "https://zneon.zeaz.dev/client-v-0-0-1/manifest.json"; do
  if curl --fail --silent --show-error --max-time 12 "$url" >/dev/null; then
    echo "PUBLIC PASS $url"
  else
    echo "PUBLIC NOT VERIFIED $url"
  fi
done
echo "=== COMPLETE: inspect $REPORT ==="
