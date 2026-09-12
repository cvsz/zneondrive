#!/usr/bin/env bash
set -Eeuo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

ENV_FILE="${ENV_FILE:-$ROOT/.env}"
RUNTIME_DIR="${RUNTIME_DIR:-$ROOT/.runtime}"
DIST_DIR="${DIST_DIR:-$ROOT/dist}"
CLIENT_DIR="${CLIENT_INSTALL_DIR:-$DIST_DIR/client}"
SERVER_DIR="${SERVER_INSTALL_DIR:-$DIST_DIR/server}"
LOG_TAIL="${LOG_TAIL:-200}"

mkdir -p "$RUNTIME_DIR" "$DIST_DIR"

die() { printf 'ERROR: %s\n' "$*" >&2; exit 1; }
note() { printf '==> %s\n' "$*"; }
warn() { printf 'WARN: %s\n' "$*" >&2; }
have() { command -v "$1" >/dev/null 2>&1; }

env_value() {
  local key="$1"
  [[ -f "$ENV_FILE" ]] || return 1
  awk -F= -v key="$key" '$1 == key { sub(/^[^=]*=/, ""); print; exit }' "$ENV_FILE"
}

compose() {
  if docker compose version >/dev/null 2>&1; then docker compose "$@";
  elif have docker-compose; then docker-compose "$@";
  else die "Docker Compose is required."; fi
}

random_secret() {
  if have openssl; then openssl rand -hex 32;
  elif have python3; then python3 -c 'import secrets; print(secrets.token_hex(32))';
  else die "openssl or python3 is required."; fi
}

env_init() {
  if [[ ! -f "$ENV_FILE" ]]; then
    cp .env.example "$ENV_FILE"
    chmod 600 "$ENV_FILE"
    note "Created $ENV_FILE"
  fi
  local current secret
  current="$(env_value GAME_SERVER_SHARED_KEY || true)"
  if [[ -z "$current" || "$current" == "zneondrive-local-server-key-change-me-0001" ]]; then
    secret="$(random_secret)"
    python3 - "$ENV_FILE" "$secret" <<'PY'
from pathlib import Path
import sys
path = Path(sys.argv[1])
secret = sys.argv[2]
lines = path.read_text(encoding="utf-8").splitlines()
out, found = [], False
for line in lines:
    if line.startswith("GAME_SERVER_SHARED_KEY="):
        out.append("GAME_SERVER_SHARED_KEY=" + secret)
        found = True
    else:
        out.append(line)
if not found:
    out.append("GAME_SERVER_SHARED_KEY=" + secret)
path.write_text("\n".join(out) + "\n", encoding="utf-8")
PY
    chmod 600 "$ENV_FILE"
    note "Generated a unique local game-server key."
  fi
}

doctor() {
  local failed=0
  printf 'PROJECT: NEON DRIVE control-surface doctor\n'
  printf 'root: %s\n' "$ROOT"
  printf 'os: %s\n' "$(uname -srm)"
  for cmd in git python3 make curl; do
    if have "$cmd"; then printf '  [ok] %-10s %s\n' "$cmd" "$(command -v "$cmd")";
    else printf '  [--] %-10s missing\n' "$cmd"; failed=$((failed+1)); fi
  done
  if have docker; then
    printf '  [ok] docker     %s\n' "$(docker --version 2>/dev/null || true)"
    docker info >/dev/null 2>&1 || { printf '  [--] docker daemon not reachable\n'; failed=$((failed+1)); }
    if docker compose version >/dev/null 2>&1 || have docker-compose; then printf '  [ok] compose\n';
    else printf '  [--] compose missing\n'; failed=$((failed+1)); fi
  else
    printf '  [--] docker missing\n'; failed=$((failed+1))
  fi
  if have go; then printf '  [ok] go         %s\n' "$(go version)"; else printf '  [--] go optional locally; CI uses Go 1.27.x\n'; fi
  if [[ -n "${UE_ROOT:-}" && -x "$UE_ROOT/Engine/Build/BatchFiles/Linux/Build.sh" ]]; then
    printf '  [ok] UE_ROOT    %s\n' "$UE_ROOT"
  else
    printf '  [--] UE_ROOT not configured for Linux UE source builds\n'
  fi
  [[ -f "$ENV_FILE" ]] && printf '  [ok] env        %s\n' "$ENV_FILE" || printf '  [--] env        run make env-init\n'
  [[ "$failed" -eq 0 ]]
}

deps_server() {
  [[ "$(uname -s)" == Linux ]] || die "Automatic dependency install currently supports Linux only."
  have apt-get || die "An apt-based host is required for automatic dependency installation."
  local sudo_cmd=""
  if [[ "$EUID" -ne 0 ]]; then have sudo || die "sudo is required."; sudo_cmd="sudo"; fi
  $sudo_cmd apt-get update
  $sudo_cmd apt-get install -y ca-certificates curl git make python3 python3-venv openssl docker.io
  if ! docker compose version >/dev/null 2>&1; then
    $sudo_cmd apt-get install -y docker-compose-v2 2>/dev/null ||
      $sudo_cmd apt-get install -y docker-compose-plugin 2>/dev/null ||
      $sudo_cmd apt-get install -y docker-compose
  fi
  have systemctl && $sudo_cmd systemctl enable --now docker || true
  note "Server host dependencies installed."
}

client_doctor() {
  printf 'PROJECT: NEON DRIVE player-client doctor\n'
  printf 'root: %s\n' "$ROOT"
  printf 'os: %s\n' "$(uname -srm)"
  have curl && printf '  [ok] curl       %s\n' "$(command -v curl)" || printf '  [--] curl missing\n'
  have python3 && printf '  [ok] python3    %s\n' "$(command -v python3)" || printf '  [--] python3 missing (needed for archive fallback/tooling)\n'
  local bin
  bin="$(find_client || true)"
  [[ -n "$bin" ]] && printf '  [ok] client     %s\n' "$bin" || printf '  [--] client package not installed\n'
  if [[ -n "${UE_ROOT:-}" && -x "$UE_ROOT/Engine/Build/BatchFiles/Linux/Build.sh" ]]; then
    printf '  [ok] UE_ROOT    %s\n' "$UE_ROOT"
  else
    printf '  [--] UE_ROOT not set (only required for source build/package)\n'
  fi
  printf '  [ok] secret boundary: player tooling does not require the game-server shared key\n'
}

deps_client() {
  [[ "$(uname -s)" == Linux ]] || die "Use tools/install-client.ps1 on Windows."
  have apt-get || die "An apt-based host is required for automatic dependency installation."
  local sudo_cmd=""
  if [[ "$EUID" -ne 0 ]]; then have sudo || die "sudo is required."; sudo_cmd="sudo"; fi
  $sudo_cmd apt-get update
  $sudo_cmd apt-get install -y ca-certificates curl git make python3 unzip tar build-essential
  note "Base Linux client/dev dependencies installed. UE 5.8 itself is not downloaded by this repository."
}

server_health() {
  local port url
  port="$(env_value GAME_API_PORT || true)"; port="${port:-18080}"
  url="http://127.0.0.1:$port/healthz"
  for _ in $(seq 1 30); do
    if curl --fail --silent "$url" >/dev/null 2>&1; then note "Game API healthy: $url"; return 0; fi
    sleep 2
  done
  compose ps >&2 || true
  die "Game API health check failed: $url"
}

server_install() {
  have docker || die "Docker missing. Run make deps-server."
  docker info >/dev/null 2>&1 || die "Docker daemon is not reachable."
  env_init
  compose pull postgres redis
  compose build game-api
  compose up -d
  server_health
}
server_up() { env_init; compose up -d; server_health; }
server_down() { compose down; }
server_restart() { compose restart; server_health; }
server_status() { compose ps; }
server_logs() { compose logs --tail "$LOG_TAIL" "$@"; }
server_reset() {
  [[ "${CONFIRM_RESET:-}" == YES ]] || die "Refusing reset. Re-run with CONFIRM_RESET=YES."
  compose down -v --remove-orphans
  rm -rf "$RUNTIME_DIR"; mkdir -p "$RUNTIME_DIR"
  note "Local PostgreSQL/Redis volumes removed."
}
db_shell() {
  local db user
  db="$(env_value POSTGRES_DB || true)"; db="${db:-zneondrive}"
  user="$(env_value POSTGRES_USER || true)"; user="${user:-zneondrive}"
  compose exec postgres psql -U "$user" -d "$db"
}
redis_cli() { compose exec redis redis-cli; }

require_ue() {
  [[ -n "${UE_ROOT:-}" ]] || die "UE_ROOT is required. Run 'make ue-detect' to search common locations."
  if [[ ! -x "$UE_ROOT/GenerateProjectFiles.sh" ]]; then
    printf 'UE_ROOT=%s\n' "$UE_ROOT" >&2
    printf 'Expected: %s/GenerateProjectFiles.sh\n' "$UE_ROOT" >&2
    printf 'Search with: make ue-detect\n' >&2
    die "Invalid UE_ROOT: GenerateProjectFiles.sh missing."
  fi
  [[ -x "$UE_ROOT/Engine/Build/BatchFiles/Linux/Build.sh" ]] || die "Invalid UE_ROOT: Linux Build.sh missing."
}

ue_detect() {
  local found=0 candidate
  printf 'Searching for Unreal Engine source roots...\n'
  while IFS= read -r candidate; do
    [[ -n "$candidate" ]] || continue
    found=1
    printf '  candidate: %s\n' "$candidate"
  done < <(
    find /opt "$HOME" /usr/local /mnt/c /mnt/d       -maxdepth 5 -type f -name GenerateProjectFiles.sh -print 2>/dev/null |
      sed 's#/GenerateProjectFiles.sh$##' |
      sort -u
  )
  if [[ "$found" -eq 0 ]]; then
    printf 'No Linux Unreal source tree found in common locations.\n'
    printf 'A placeholder path such as /opt/UnrealEngine-5.8 is not enough; UE 5.8 source must actually be present and built.\n'
    return 1
  fi
}
client_generate() { require_ue; "$UE_ROOT/GenerateProjectFiles.sh" -project="$ROOT/game/NeonDrive.uproject" -game; }
ue_build() { local target="$1"; require_ue; "$UE_ROOT/Engine/Build/BatchFiles/Linux/Build.sh" "$target" Linux Development "$ROOT/game/NeonDrive.uproject" -WaitMutex; }
client_build() { ue_build NeonDriveClient; }
editor_build() { ue_build NeonDriveEditor; }
game_server_build() { ue_build NeonDriveServer; }

require_uat() {
  require_ue
  [[ -x "$UE_ROOT/Engine/Build/BatchFiles/RunUAT.sh" ]] || die "RunUAT.sh is missing under UE_ROOT."
}

client_package_linux() {
  require_uat
  local out="$DIST_DIR/packages/client-linux"
  rm -rf "$out"; mkdir -p "$out"
  "$UE_ROOT/Engine/Build/BatchFiles/RunUAT.sh" BuildCookRun \
    -project="$ROOT/game/NeonDrive.uproject" -noP4 -build -cook -stage -pak -archive \
    -archivedirectory="$out" -targetplatform=Linux -clientconfig=Development -client -utf8output
  note "Linux player package archived under $out"
}

game_server_package_linux() {
  require_uat
  local out="$DIST_DIR/packages/server-linux"
  rm -rf "$out"; mkdir -p "$out"
  "$UE_ROOT/Engine/Build/BatchFiles/RunUAT.sh" BuildCookRun \
    -project="$ROOT/game/NeonDrive.uproject" -noP4 -build -cook -stage -pak -archive \
    -archivedirectory="$out" -server -noclient -serverplatform=Linux -serverconfig=Development -utf8output
  note "Linux dedicated-server package archived under $out"
}

package_all_linux() { client_package_linux; game_server_package_linux; }

install_package() {
  local src="$1" dst="$2"
  [[ -e "$src" ]] || die "Package not found: $src"
  rm -rf "$dst"; mkdir -p "$dst"
  if [[ -d "$src" ]]; then cp -a "$src"/. "$dst"/;
  elif [[ "$src" == *.zip ]]; then
    if have unzip; then unzip -q "$src" -d "$dst"; else python3 -m zipfile -e "$src" "$dst"; fi
  elif [[ "$src" == *.tar.gz || "$src" == *.tgz ]]; then tar -xzf "$src" -C "$dst";
  elif [[ "$src" == *.tar ]]; then tar -xf "$src" -C "$dst";
  else die "Supported packages: directory, .zip, .tar, .tar.gz, .tgz"; fi
}
find_client() { find "$CLIENT_DIR" -type f \( -name NeonDrive -o -name NeonDrive.sh -o -name NeonDrive.exe \) -print 2>/dev/null | head -n1; }
find_game_server() { find "$SERVER_DIR" -type f \( -name NeonDriveServer -o -name NeonDriveServer.sh -o -name NeonDriveServer.exe \) -print 2>/dev/null | head -n1; }

client_install() {
  local package="${1:-${CLIENT_PACKAGE:-}}"
  if [[ -n "$package" ]]; then
    install_package "$package" "$CLIENT_DIR"
    local bin; bin="$(find_client || true)"
    [[ -n "$bin" ]] || die "No NeonDrive client executable found in package."
    [[ "$bin" == *.exe ]] || chmod +x "$bin" 2>/dev/null || true
    note "Player client installed: $bin"
    return
  fi
  [[ -n "${UE_ROOT:-}" ]] || die "Set CLIENT_PACKAGE or UE_ROOT."
  client_generate
  client_build
  note "NeonDriveClient source target built. Packaged-play evidence is still a separate gate."
}

client_play() {
  local api port bin
  api="${ZNEON_GAME_API_URL:-}"
  if [[ -z "$api" ]]; then port="$(env_value GAME_API_PORT || true)"; port="${port:-18080}"; api="http://127.0.0.1:$port"; fi
  [[ "$api" =~ ^https?://[^[:space:]]+$ ]] || die "ZNEON_GAME_API_URL must be an http(s) URL."
  bin="$(find_client || true)"
  if [[ -n "$bin" ]]; then
    [[ "$bin" != *.exe ]] || die "Windows client detected; launch it with tools/install-client.ps1 on Windows."
    note "Launching client against $api"
    ZNEON_GAME_API_URL="$api" "$bin" -ZNeonApi="$api" ${CLIENT_ARGS:-}
  elif [[ -n "${UE_ROOT:-}" && -x "$UE_ROOT/Engine/Binaries/Linux/UnrealEditor" ]]; then
    note "Launching UE development client against $api"
    ZNEON_GAME_API_URL="$api" "$UE_ROOT/Engine/Binaries/Linux/UnrealEditor" "$ROOT/game/NeonDrive.uproject" -game -ZNeonApi="$api" ${CLIENT_ARGS:-}
  else
    die "No packaged client found. Install a package or configure UE_ROOT."
  fi
}

game_server_install() {
  local package="${1:-${SERVER_PACKAGE:-}}"
  [[ -n "$package" ]] || die "Set SERVER_PACKAGE or pass a package path."
  install_package "$package" "$SERVER_DIR"
  local bin; bin="$(find_game_server || true)"
  [[ -n "$bin" ]] || die "No NeonDriveServer executable found in package."
  [[ "$bin" == *.exe ]] || chmod +x "$bin" 2>/dev/null || true
  note "Dedicated game server installed: $bin"
}
game_server_start() {
  env_init
  local bin key port internal pidfile
  bin="$(find_game_server || true)"; [[ -n "$bin" ]] || die "Install a dedicated server package first."
  [[ "$bin" != *.exe ]] || die "Manage Windows server packages from Windows."
  key="$(env_value GAME_SERVER_SHARED_KEY || true)"; [[ "${#key}" -ge 32 ]] || die "GAME_SERVER_SHARED_KEY missing/short."
  port="$(env_value GAME_API_PORT || true)"; port="${port:-18080}"
  internal="${ZNEON_GAME_API_INTERNAL_URL:-http://127.0.0.1:$port}"
  pidfile="$RUNTIME_DIR/game-server.pid"
  if [[ -f "$pidfile" ]] && kill -0 "$(cat "$pidfile")" 2>/dev/null; then die "Game server already running."; fi
  ZNEON_GAME_SERVER_KEY="$key" ZNEON_GAME_API_INTERNAL_URL="$internal" nohup "$bin" -log ${SERVER_ARGS:-} >"$RUNTIME_DIR/game-server.log" 2>&1 &
  echo $! >"$pidfile"; sleep 1
  kill -0 "$(cat "$pidfile")" 2>/dev/null || { tail -n100 "$RUNTIME_DIR/game-server.log" >&2 || true; rm -f "$pidfile"; die "Game server exited during startup."; }
  note "Dedicated game server started: PID $(cat "$pidfile")"
}
game_server_stop() {
  local f="$RUNTIME_DIR/game-server.pid"; [[ -f "$f" ]] || { note "Game server not tracked as running."; return; }
  local pid; pid="$(cat "$f")"
  kill "$pid" 2>/dev/null || true
  for _ in $(seq 1 10); do kill -0 "$pid" 2>/dev/null || break; sleep 1; done
  kill -0 "$pid" 2>/dev/null && kill -9 "$pid" || true
  rm -f "$f"; note "Dedicated game server stopped."
}
game_server_status() {
  local f="$RUNTIME_DIR/game-server.pid"
  if [[ -f "$f" ]] && kill -0 "$(cat "$f")" 2>/dev/null; then
    printf 'running pid=%s\n' "$(cat "$f")"
  else
    printf 'stopped\n'
  fi
}
game_server_logs() { touch "$RUNTIME_DIR/game-server.log"; tail -n "$LOG_TAIL" -f "$RUNTIME_DIR/game-server.log"; }

full_install() {
  local client_package="${1:-${CLIENT_PACKAGE:-}}"
  local server_package="${SERVER_PACKAGE:-}"
  [[ -n "$client_package" ]] || die "CLIENT_PACKAGE is required for full-install."
  [[ -e "$client_package" ]] || die "Client package not found: $client_package"
  [[ -n "$server_package" ]] || die "SERVER_PACKAGE is required for full-install."
  [[ -e "$server_package" ]] || die "Server package not found: $server_package"

  note "Full-install preflight passed."
  server_install
  client_install "$client_package"
  game_server_install "$server_package"
}

full_up() {
  local bin
  bin="$(find_game_server || true)"
  [[ -n "$bin" ]] || die "Dedicated game-server package is not installed. Run make game-server-install first."
  server_up
  game_server_start
}
full_down() { game_server_stop || true; server_down; }
status_all() {
  printf '\n-- service plane --\n'; server_status || true
  printf '\n-- dedicated server --\n'; game_server_status || true
  printf '\n-- player client --\n'; local b; b="$(find_client || true)"; [[ -n "$b" ]] && printf 'installed: %s\n' "$b" || printf 'not installed\n'
}

control_panel() {
  while true; do
    cat <<'MENU'

PROJECT: NEON DRIVE — Control Panel
 1 Doctor                       9 Build Linux player client
 2 Initialize .env             10 Play client
 3 Install server stack        11 Install game-server package
 4 Start server stack          12 Start dedicated game server
 5 Stop server stack           13 Stop dedicated game server
 6 Server status               14 Full-stack status
 7 Server logs                 15 Run make ci
 8 Install player package      16 Reset LOCAL DB/Redis volumes
17 Package Linux player        18 Package Linux game server
19 Package both Linux builds
 0 Exit
MENU
    read -r -p "Select: " choice
    case "$choice" in
      1) doctor || true;; 2) env_init;; 3) server_install;; 4) server_up;; 5) server_down;;
      6) server_status;; 7) server_logs;; 8) read -r -p "Client package: " p; client_install "$p";;
      9) client_generate; client_build;; 10) client_play;;
      11) read -r -p "Server package: " p; game_server_install "$p";; 12) game_server_start;; 13) game_server_stop;;
      14) status_all;; 15) make ci;;
      16) read -r -p "Type YES to destroy LOCAL DB/Redis volumes: " a; [[ "$a" == YES ]] && CONFIRM_RESET=YES server_reset || warn "Reset cancelled.";;
      17) client_package_linux;; 18) game_server_package_linux;; 19) package_all_linux;;
      0) return;; *) warn "Unknown selection.";;
    esac
  done
}

usage() {
  cat <<'EOF'
Usage: bash tools/zneondrive-control.sh <command> [args]

doctor | client-doctor | ue-detect | env-init | deps-server | deps-client | control-panel
server-install | server-up | server-down | server-restart | server-status
server-health | server-logs [service] | server-reset | db-shell | redis-cli
client-generate | client-build | client-package-linux | editor-build | client-install [package] | client-play
game-server-build | game-server-package-linux | package-all-linux | game-server-install [package] | game-server-start | game-server-stop
game-server-status | game-server-logs | full-install [client-package] | full-up | full-down | status
EOF
}

cmd="${1:-help}"; shift || true
case "$cmd" in
  doctor) doctor;; client-doctor) client_doctor;; ue-detect) ue_detect;; env-init) env_init;; deps-server) deps_server;; deps-client) deps_client;;
  server-install) server_install;; server-up) server_up;; server-down) server_down;; server-restart) server_restart;;
  server-status) server_status;; server-health) server_health;; server-logs) server_logs "$@";; server-reset) server_reset;;
  db-shell) db_shell;; redis-cli) redis_cli;; client-generate) client_generate;; client-build) client_build;; editor-build) editor_build;;
  client-install) client_install "${1:-}";; client-play) client_play;; client-package-linux) client_package_linux;; game-server-build) game_server_build;;
  game-server-package-linux) game_server_package_linux;; package-all-linux) package_all_linux;; game-server-install) game_server_install "${1:-}";; game-server-start) game_server_start;; game-server-stop) game_server_stop;;
  game-server-status) game_server_status;; game-server-logs) game_server_logs;; full-install) full_install "${1:-}";;
  full-up) full_up;; full-down) full_down;; status) status_all;; control-panel) control_panel;; help|-h|--help) usage;;
  *) usage; die "Unknown command: $cmd";;
esac
