#!/usr/bin/env bash
set -Eeuo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
PROJECT="$ROOT/game/NeonDrive.uproject"
DIST_DIR="${DIST_DIR:-$ROOT/dist}"

fail() { printf 'ERROR: %s\n' "$*" >&2; exit 1; }
note() { printf '==> %s\n' "$*"; }

require_root() {
  [[ -n "${UE_ROOT:-}" ]] || fail "UE_ROOT is required. Run 'make ue-detect' to search common locations."
  [[ -x "$UE_ROOT/Engine/Build/BatchFiles/Linux/Build.sh" ]] || fail "Invalid UE_ROOT: Linux Build.sh missing."
}

find_bundled_dotnet() {
  local candidate
  while IFS= read -r candidate; do
    [[ -x "$candidate" ]] && { printf '%s\n' "$candidate"; return 0; }
  done < <(find "$UE_ROOT/Engine/Binaries/ThirdParty/DotNet" -maxdepth 4 -type f -path '*/linux/dotnet' -print 2>/dev/null | sort -V -r)
  command -v dotnet 2>/dev/null || true
}

project_generate() {
  require_root
  if [[ -x "$UE_ROOT/GenerateProjectFiles.sh" ]]; then
    note "Generating project files with source-tree GenerateProjectFiles.sh"
    "$UE_ROOT/GenerateProjectFiles.sh" -project="$PROJECT" -game
    return
  fi

  local ubt dotnet
  ubt="$UE_ROOT/Engine/Binaries/DotNET/UnrealBuildTool/UnrealBuildTool.dll"
  [[ -f "$ubt" ]] || fail "UE installed build is missing UnrealBuildTool.dll and no GenerateProjectFiles.sh exists."
  dotnet="$(find_bundled_dotnet)"
  [[ -n "$dotnet" ]] || fail "UE installed build requires bundled or system dotnet to run UnrealBuildTool."

  note "Generating project files with installed-build UnrealBuildTool"
  "$dotnet" "$ubt" -projectfiles -project="$PROJECT" -game
}

ue_build() {
  local target="${1:-}"
  [[ -n "$target" ]] || fail "build requires a target name."
  require_root
  "$UE_ROOT/Engine/Build/BatchFiles/Linux/Build.sh" "$target" Linux Development "$PROJECT" -WaitMutex
}

require_uat() {
  require_root
  [[ -x "$UE_ROOT/Engine/Build/BatchFiles/RunUAT.sh" ]] || fail "RunUAT.sh is missing under UE_ROOT."
}

package_client() {
  require_uat
  local out="$DIST_DIR/packages/client-linux"
  rm -rf "$out"; mkdir -p "$out"
  "$UE_ROOT/Engine/Build/BatchFiles/RunUAT.sh" BuildCookRun \
    -project="$PROJECT" -noP4 -build -cook -stage -pak -archive \
    -archivedirectory="$out" -targetplatform=Linux -clientconfig=Development -client -utf8output
  note "Linux player package archived under $out"
}

package_server() {
  require_uat
  local out="$DIST_DIR/packages/server-linux"
  rm -rf "$out"; mkdir -p "$out"
  "$UE_ROOT/Engine/Build/BatchFiles/RunUAT.sh" BuildCookRun \
    -project="$PROJECT" -noP4 -build -cook -stage -pak -archive \
    -archivedirectory="$out" -server -noclient -serverplatform=Linux -serverconfig=Development -utf8output
  note "Linux dedicated-server package archived under $out"
}

ue_detect() {
  local roots_raw="${UE_SEARCH_ROOTS:-/opt:$HOME:/usr/local:/mnt/c:/mnt/d:/mnt/zworkforce-storage}"
  local found=0 root candidate
  local -a roots=()
  IFS=':' read -r -a roots <<< "$roots_raw"
  printf 'Searching for Linux Unreal Engine roots (source or installed build)...\n'
  for root in "${roots[@]}"; do
    [[ -d "$root" ]] || continue
    while IFS= read -r candidate; do
      [[ -n "$candidate" ]] || continue
      found=1
      printf '  candidate: %s\n' "${candidate%/Engine/Build/BatchFiles/Linux/Build.sh}"
    done < <(find "$root" -maxdepth 9 -type f -path '*/Engine/Build/BatchFiles/Linux/Build.sh' -print 2>/dev/null | sort -u)
  done
  [[ "$found" -eq 1 ]] || { printf 'No usable Linux Unreal Engine root found.\n' >&2; return 1; }
}

usage() {
  cat <<'EOF'
Usage: bash tools/ue-linux.sh <command> [args]

Commands:
  detect
  generate
  build <NeonDriveClient|NeonDriveEditor|NeonDriveServer>
  package-client
  package-server
  package-all
EOF
}

cmd="${1:-help}"; shift || true
case "$cmd" in
  detect) ue_detect ;;
  generate) project_generate ;;
  build) ue_build "${1:-}" ;;
  package-client) package_client ;;
  package-server) package_server ;;
  package-all) package_client; package_server ;;
  help|-h|--help) usage ;;
  *) usage; fail "Unknown command: $cmd" ;;
esac
