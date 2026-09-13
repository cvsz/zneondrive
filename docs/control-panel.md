# Full-Stack Install & Control Panel

This control surface makes the current zNeonDrive runtime easier to install and operate without weakening the evidence boundary.

## Profiles

### Server service plane

The local/server service plane is:

- PostgreSQL 17
- Redis 8
- Go game API
- optional packaged Unreal dedicated server

Install and start:

```bash
make deps-server
make env-init
make server-install
make server-status
make server-health
```

`make env-init` creates a local `.env` and replaces the example game-server key with a generated secret. The file remains git-ignored.

### Player client

A player client never receives `GAME_SERVER_SHARED_KEY` or `ZNEON_GAME_SERVER_KEY`.

Linux package install:

```bash
make client-install CLIENT_PACKAGE=/path/to/NeonDrive-client.zip
ZNEON_GAME_API_URL=https://api.example.com make client-play
```

Windows package install:

```powershell
powershell -ExecutionPolicy Bypass -File tools/install-client.ps1 `
  -Mode package `
  -Package C:\packages\NeonDrive-client.zip `
  -ApiUrl https://api.example.com `
  -CreateShortcut
```

Play later:

```powershell
powershell -ExecutionPolicy Bypass -File tools/install-client.ps1 `
  -Mode play `
  -ApiUrl https://api.example.com
```

### Unreal Linux development client/server

Linux supports both a full source tree and a precompiled installed-build layout. The latter may not contain a root-level `GenerateProjectFiles.sh`; in that case the repository invokes `UnrealBuildTool.dll -projectfiles` through Unreal's bundled `dotnet` runtime.

```bash
export UE_ROOT=/opt/UnrealEngine-5.8
make deps-client
make ue-detect
make client-generate
make client-build
make game-server-build
```

The helper requires `Engine/Build/BatchFiles/Linux/Build.sh`. Project generation uses this order:

1. root `GenerateProjectFiles.sh` when present and executable;
2. `Engine/Binaries/DotNET/UnrealBuildTool/UnrealBuildTool.dll` with a bundled `Engine/Binaries/ThirdParty/DotNet/*/linux/dotnet`;
3. system `dotnet` only when the bundled runtime is unavailable.

For non-standard mount points, override the detection roots without changing the build root itself:

```bash
UE_SEARCH_ROOTS=/opt:$HOME:/mnt/zworkforce-storage make ue-detect
```

Packaging uses the same `UE_ROOT` and requires `Engine/Build/BatchFiles/RunUAT.sh`:

```bash
make client-package-linux
make game-server-package-linux
# or both
make package-all-linux
```

Windows source development remains:

```powershell
$env:UE_ROOT = "D:\UnrealEngine-5.8"
powershell -ExecutionPolicy Bypass -File tools/install-client.ps1 -Mode source
```

The explicit `NeonDriveClient` target remains a Client target. Tooling compatibility and fixture tests are not successful UE 5.8 build/package evidence; that gate requires a real retained Client/Server build artifact from the selected engine installation.

## Dedicated gameplay server

Install a packaged server:

```bash
make game-server-install SERVER_PACKAGE=/path/to/NeonDrive-server.tar.gz
make game-server-start
make game-server-status
make game-server-logs
```

The launcher reads the internal shared key from local `.env` and passes it only through the dedicated-server environment.

## Interactive control panel

```bash
make control-panel
```

The menu provides:

- host diagnostics,
- environment initialization,
- server install/start/stop/status/logs,
- player-package installation,
- Linux UE client build,
- client launch,
- dedicated-server install/start/stop,
- full-stack status,
- repository CI,
- guarded local database/Redis reset.

The Makefile's UE build/package targets use `tools/ue-linux.sh` so installed builds without a root `GenerateProjectFiles.sh` are supported. The older interactive control-panel source-build menu remains operational for source-tree layouts; use the Make targets above for installed-build development until that menu is unified with the helper.

## Full-stack lifecycle

When both packaged client and packaged dedicated server artifacts exist:

```bash
CLIENT_PACKAGE=/path/client.zip \
SERVER_PACKAGE=/path/server.tar.gz \
make full-install

make full-up
make status
make full-down
```

`full-up` is strict: it starts the Compose service plane and then requires an installed packaged dedicated server. This avoids pretending that a Go/PostgreSQL-only stack is the complete game runtime.

## Runtime endpoint selection

The Unreal player client resolves its public API base URL in this precedence order:

1. `DefaultEngine.ini` / `[NeonDrive.Service] BaseUrl`
2. `ZNEON_GAME_API_URL`
3. `-ZNeonApi=https://api.example.com`

The command-line value wins. This is intentionally limited to the public API URL. It does not provide a path for a player client to receive the game-server shared key.

## Destructive operations

Local volume reset is guarded:

```bash
CONFIRM_RESET=YES make server-reset
```

It destroys the local PostgreSQL and Redis Compose volumes. It is not a production backup/restore or DR mechanism.

## Evidence boundary

These tools improve installation and operation only. They do **not** make the following claims green by themselves:

- successful archived UE 5.8 Client/Server build/package,
- packaged Unreal client/server ↔ Go E2E,
- Garage 17 playability,
- First Ignition playability,
- load/soak,
- HA/DR,
- production deployment.

Those gates still require executable retained evidence.
