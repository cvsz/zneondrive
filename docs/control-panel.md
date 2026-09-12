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

### Unreal source-development client

Linux:

```bash
export UE_ROOT=/opt/UnrealEngine-5.8
make deps-client
make client-generate
make client-build
```

Windows:

```powershell
$env:UE_ROOT = "D:\UnrealEngine-5.8"
powershell -ExecutionPolicy Bypass -File tools/install-client.ps1 -Mode source
```

The explicit `NeonDriveClient` target is a Client target. A successful source build is still not equivalent to archived packaged-play evidence.

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

- successful archived UE 5.8 source build,
- packaged Unreal client/server E2E,
- Garage 17 playability,
- First Ignition playability,
- load/soak,
- HA/DR,
- production deployment.

Those gates still require executable retained evidence.
