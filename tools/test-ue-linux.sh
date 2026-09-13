#!/usr/bin/env bash
set -Eeuo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
TMP="$(mktemp -d)"
trap 'rm -rf "$TMP"' EXIT

assert_contains() {
  local file="$1" needle="$2"
  grep -F -- "$needle" "$file" >/dev/null || { echo "missing '$needle' in $file" >&2; cat "$file" >&2; exit 1; }
}

make_common() {
  local ue="$1"
  mkdir -p "$ue/Engine/Build/BatchFiles/Linux"
  cat >"$ue/Engine/Build/BatchFiles/Linux/Build.sh" <<'SH'
#!/usr/bin/env bash
printf '%s\n' "$@" >"${UE_TEST_BUILD_CAPTURE:?}"
SH
  chmod +x "$ue/Engine/Build/BatchFiles/Linux/Build.sh"
  mkdir -p "$ue/Engine/Build/BatchFiles"
  cat >"$ue/Engine/Build/BatchFiles/RunUAT.sh" <<'SH'
#!/usr/bin/env bash
printf '%s\n' "$@" >"${UE_TEST_UAT_CAPTURE:?}"
SH
  chmod +x "$ue/Engine/Build/BatchFiles/RunUAT.sh"
}

installed="$TMP/InstalledUE"
make_common "$installed"
mkdir -p "$installed/Engine/Binaries/DotNET/UnrealBuildTool"
touch "$installed/Engine/Binaries/DotNET/UnrealBuildTool/UnrealBuildTool.dll"
mkdir -p "$installed/Engine/Binaries/ThirdParty/DotNet/10.0.0/linux"
cat >"$installed/Engine/Binaries/ThirdParty/DotNet/10.0.0/linux/dotnet" <<'SH'
#!/usr/bin/env bash
printf '%s\n' "$@" >"${UE_TEST_GENERATE_CAPTURE:?}"
SH
chmod +x "$installed/Engine/Binaries/ThirdParty/DotNet/10.0.0/linux/dotnet"

UE_TEST_GENERATE_CAPTURE="$TMP/installed-generate.txt" \
UE_TEST_BUILD_CAPTURE="$TMP/installed-build.txt" \
UE_TEST_UAT_CAPTURE="$TMP/installed-uat.txt" \
UE_ROOT="$installed" bash "$ROOT/tools/ue-linux.sh" generate
assert_contains "$TMP/installed-generate.txt" "$installed/Engine/Binaries/DotNET/UnrealBuildTool/UnrealBuildTool.dll"
assert_contains "$TMP/installed-generate.txt" "-projectfiles"
assert_contains "$TMP/installed-generate.txt" "-project=$ROOT/game/NeonDrive.uproject"
assert_contains "$TMP/installed-generate.txt" "-game"

UE_TEST_BUILD_CAPTURE="$TMP/installed-build.txt" UE_ROOT="$installed" bash "$ROOT/tools/ue-linux.sh" build NeonDriveClient
assert_contains "$TMP/installed-build.txt" "NeonDriveClient"
assert_contains "$TMP/installed-build.txt" "Linux"
assert_contains "$TMP/installed-build.txt" "Development"
assert_contains "$TMP/installed-build.txt" "$ROOT/game/NeonDrive.uproject"
assert_contains "$TMP/installed-build.txt" "-WaitMutex"

UE_TEST_UAT_CAPTURE="$TMP/installed-uat.txt" UE_ROOT="$installed" DIST_DIR="$TMP/dist" bash "$ROOT/tools/ue-linux.sh" package-server
assert_contains "$TMP/installed-uat.txt" "BuildCookRun"
assert_contains "$TMP/installed-uat.txt" "-server"
assert_contains "$TMP/installed-uat.txt" "-noclient"
assert_contains "$TMP/installed-uat.txt" "-serverplatform=Linux"

source="$TMP/SourceUE"
make_common "$source"
cat >"$source/GenerateProjectFiles.sh" <<'SH'
#!/usr/bin/env bash
printf '%s\n' "$@" >"${UE_TEST_GENERATE_CAPTURE:?}"
SH
chmod +x "$source/GenerateProjectFiles.sh"
UE_TEST_GENERATE_CAPTURE="$TMP/source-generate.txt" UE_ROOT="$source" bash "$ROOT/tools/ue-linux.sh" generate
assert_contains "$TMP/source-generate.txt" "-project=$ROOT/game/NeonDrive.uproject"
assert_contains "$TMP/source-generate.txt" "-game"

UE_SEARCH_ROOTS="$TMP" bash "$ROOT/tools/ue-linux.sh" detect >"$TMP/detect.txt"
assert_contains "$TMP/detect.txt" "$installed"
assert_contains "$TMP/detect.txt" "$source"

printf 'UE Linux tooling validation OK: installed-build UBT, source-tree generation, build, package, and detection paths covered\n'
