# NeonDrive Unreal Runtime

## Required engine

Use Unreal Engine 5.8.x. The repository targets the UE 5.8 API line; prefer the current 5.8 hotfix after regression testing.

The Linux tooling supports both:
- a UE 5.8 source tree with `GenerateProjectFiles.sh`, and
- a precompiled/installed UE 5.8 build that provides `UnrealBuildTool.dll`, `Build.sh`, `RunUAT.sh`, and a bundled/system `dotnet` runtime.

`tools/ue-linux.sh` validates `Engine/Build/Build.version` and rejects non-5.8 engine roots before build/package commands run.

## Generate project files

Recommended Linux path:

```bash
export UE_ROOT=/opt/UnrealEngine-5.8
make client-generate
```

For source trees the helper delegates to `$UE_ROOT/GenerateProjectFiles.sh`. For installed builds without that script, it invokes `UnrealBuildTool.dll -projectfiles` using the engine-bundled Linux `dotnet` when available. A local shim is not required.

Windows:

```powershell
& "$env:UE_ROOT\Engine\Build\BatchFiles\GenerateProjectFiles.bat" -project="$PWD\game\NeonDrive.uproject" -game
```

## Build

Linux:

```bash
export UE_ROOT=/opt/UnrealEngine-5.8
make editor-build
make client-build
make game-server-build
```

The Make targets route through `tools/ue-linux.sh`, which calls the engine's canonical `Engine/Build/BatchFiles/Linux/Build.sh` for the selected target.

To inspect a candidate engine root without compiling:

```bash
UE_ROOT=/opt/UnrealEngine-5.8 bash tools/ue-linux.sh info
make ue-detect
```

## Package

After Client and Server targets build successfully:

```bash
make client-package-linux
make game-server-package-linux
# or both, sequentially
make package-all-linux
```

Packages are archived separately under `dist/packages/client-linux` and `dist/packages/server-linux` through the engine's `RunUAT.sh BuildCookRun` path.

## Prototype driving

Default controls:
- W/S: throttle/reverse
- A/D: steering

The pawn is intentionally simple. Client input is clamped on the server and movement is performed only by authority, then replicated by Unreal. This proves the authority direction; it is not final vehicle physics or anti-cheat.

## Dedicated server CI

`.github/workflows/unreal-source-build.yml` is manual and requires a self-hosted Linux runner with:
- label `unreal-5.8`,
- `UE_ROOT` pointing at a usable UE 5.8 source or installed build,
- sufficient disk/RAM for Unreal C++ builds.

A workflow file, engine detection, or successful project-file generation is not build evidence by itself. Keep the roadmap build-evidence item open until Client/Server builds complete successfully and retained artifacts/logs exist.

## v0.5 service-plane binding

The client-side `UNDServiceSubsystem` bootstraps/resumes the local prototype identity and requests a short-lived gameplay ticket. The owning `ANDPlayerController` sends only that ticket through a Server RPC.

The dedicated server must have:

```bash
export ZNEON_GAME_SERVER_KEY='same-value-as-GAME_SERVER_SHARED_KEY'
export ZNEON_GAME_API_INTERNAL_URL='http://127.0.0.1:18080'
```

The server redeems the ticket directly with the Go service plane and binds the returned durable vehicle ID, build revision, active parts, and Roadworthy state to the replicated pawn.

Do not package `ZNEON_GAME_SERVER_KEY` in a client build. Plain HTTP is local-development only; production transport remains an explicit security gate.

## Player client target and runtime endpoint

The repository includes an explicit `NeonDriveClient` target for client-only builds.

Linux:

```bash
export UE_ROOT=/opt/UnrealEngine-5.8
make client-generate
make client-build
```

Windows:

```powershell
$env:UE_ROOT = "D:\\UnrealEngine-5.8"
powershell -ExecutionPolicy Bypass -File tools/install-client.ps1 -Mode source
```

A packaged player client may select the public Go API endpoint without rebuilding:

```text
ZNEON_GAME_API_URL=https://api.example.com
-ZNeonApi=https://api.example.com
```

The command-line override wins. Neither mechanism accepts the dedicated-server shared key.

For player package installation, dedicated-server lifecycle, Compose service-plane control, and the interactive menu, see [Full-Stack Install & Control Panel](../docs/control-panel.md).
