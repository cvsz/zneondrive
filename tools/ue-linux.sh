#!/usr/bin/env bash
set -Eeuo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
PROJECT="${ZNEON_UE_PROJECT:-$ROOT/game/NeonDrive.uproject}"
DIST_DIR="${DIST_DIR:-$ROOT/dist}"

fail() { printf 'ERROR: %s\n' "$*" >&2; exit 1; }
note() { printf '==> %s\n' "$*"; }

require_root() {
  [[ -n "${UE_ROOT:-}" ]] || fail "UE_ROOT is required. Run 'make ue-detect' to search common locations."
  [[ -f "$UE_ROOT/Engine/Build/Build.version" ]] || fail "Invalid UE_ROOT: Engine/Build/Build.version missing."
  [[ -x "$UE_ROOT/Engine/Build/BatchFiles/Linux/Build.sh" ]] || fail "Invalid UE_ROOT: Linux Build.sh missing."
  [[ -x "$UE_ROOT/Engine/Build/BatchFiles/RunUAT.sh" ]] || fail "Invalid UE_ROOT: RunUAT.sh missing."
  [[ -f "$PROJECT" ]] || fail "Unreal project missing: $PROJECT"
}

version_json() {
  require_root
  python3 - "$UE_ROOT/Engine/Build/Build.version" <<'PY'
import json, sys
p = sys.argv[1]
v = json.load(open(p, encoding='utf-8'))
major = int(v.get('MajorVersion', -1))
minor = int(v.get('MinorVersion', -1))
patch = int(v.get('PatchVersion', 0))
if (major, minor) != (5, 8):
    raise SystemExit(f"ERROR: Unreal Engine 5.8.x required; found {major}.{minor}.{patch}")
print(f"{major}.{minor}.{patch}")
PY
}

find_dotnet() {
  if [[ -n "${UE_DOTNET:-}" && -x "$UE_DOTNET" ]]; then
    printf '%s\n' "$UE_DOTNET"
    return
  fi
  local bundled=""
  bundled="$(find "$UE_ROOT/Engine/Binaries/ThirdParty/DotNet" -type f -path '*/linux/dotnet' -perm -u+x 2>/dev/null | sort -V | tail -n1 || true)"
  if [[ -n "$bundled" ]]; then
    printf '%s\n' "$bundled"
    return
  fi
  command -v dotnet 2>/dev/null || true
}

ubt_dll() {
  local dll="$UE_ROOT/Engine/Binaries/DotNET/UnrealBuildTool/UnrealBuildTool.dll"
  [[ -f "$dll" ]] && printf '%s\n' "$dll"
}

engine_info() {
  local version
  version="$(version_json)"
  printf 'UE_ROOT=%s\n' "$UE_ROOT"
  printf 'UE_VERSION=%s\n' "$version"
  if [[ -x "$UE_ROOT/GenerateProjectFiles.sh" ]]; then
    printf 'UE_PROJECTFILES_MODE=source-script\n'
  elif [[ -n "$(ubt_dll || true)" ]]; then
    printf 'UE_PROJECTFILES_MODE=installed-ubt\n'
  else
    printf 'UE_PROJECTFILES_MODE=unavailable\n'
    return 1
  fi
}

projectfiles() {
  version_json >/dev/null
  if [[ -x "$UE_ROOT/GenerateProjectFiles.sh" ]]; then
    note "Generating project files with source-tree GenerateProjectFiles.sh"
    "$UE_ROOT/GenerateProjectFiles.sh" -project="$PROJECT" -game
    return
  fi

  local dll dotnet_bin
  dll="$(ubt_dll || true)"
  [[ -n "$dll" ]] || fail "No GenerateProjectFiles.sh or UnrealBuildTool.dll found under UE_ROOT."
  dotnet_bin="$(find_dotnet)"
  [[ -n "$dotnet_bin" && -x "$dotnet_bin" ]] || fail "Installed UE build requires bundled/system dotnet to run UnrealBuildTool."
  note "Generating project files with installed-build UnrealBuildTool"
  "$dotnet_bin" "$dll" -projectfiles -project="$PROJECT" -game -engine
}

build_target() {
  local target="${1:-}"
  [[ -n "$target" ]] || fail "build-target requires a target name"
  version_json >/dev/null
  "$UE_ROOT/Engine/Build/BatchFiles/Linux/Build.sh" "$target" Linux Development "$PROJECT" -WaitMutex
}

package_client() {
  version_json >/dev/null
  local out="$DIST_DIR/packages/client-linux"
  rm -rf "$out"; mkdir -p "$out"
  "$UE_ROOT/Engine/Build/BatchFiles/RunUAT.sh" BuildCookRun \
    -project="$PROJECT" -noP4 -build -cook -stage -pak -archive \
    -archivedirectory="$out" -targetplatform=Linux -clientconfig=Development -client -utf8output
  note "Linux player package archived under $out"
}

package_server() {
  version_json >/dev/null
  local out="$DIST_DIR/packages/server-linux"
  rm -rf "$out"; mkdir -p "$out"
  "$UE_ROOT/Engine/Build/BatchFiles/RunUAT.sh" BuildCookRun \
    -project="$PROJECT" -noP4 -build -cook -stage -pak -archive \
    -archivedirectory="$out" -server -noclient -serverplatform=Linux -serverconfig=Development -utf8output
  note "Linux dedicated-server package archived under $out"
}

ue_detect() {
  local found=0 version_file candidate
  printf 'Searching for Linux Unreal Engine 5.8 roots...\n'
  while IFS= read -r version_file; do
    [[ -n "$version_file" ]] || continue
    candidate="$(dirname "$(dirname "$(dirname "$version_file")")")"
    if UE_ROOT="$candidate" ZNEON_UE_PROJECT="$PROJECT" bash "$0" info >/dev/null 2>&1; then
      found=1
      printf '  candidate: %s\n' "$candidate"
    fi
  done < <(find /opt "$HOME" /usr/local /mnt/c /mnt/d /mnt/zworkforce-storage \
    -maxdepth 7 -type f -path '*/Engine/Build/Build.version' -print 2>/dev/null | sort -u)
  [[ "$found" -eq 1 ]] || fail "No usable Linux Unreal Engine 5.8 installation found in common locations."
}

cmd="${1:-help}"; shift || true
case "$cmd" in
  info) engine_info ;;
  detect) ue_detect ;;
  generate) projectfiles ;;
  build-target) build_target "${1:-}" ;;
  package-client) package_client ;;
  package-server) package_server ;;
  package-all) package_client; package_server ;;
  help|-h|--help)
    printf '%s\n' 'usage: bash tools/ue-linux.sh {info|detect|generate|build-target TARGET|package-client|package-server|package-all}'
    ;;
  *) fail "Unknown command: $cmd" ;;
esac
