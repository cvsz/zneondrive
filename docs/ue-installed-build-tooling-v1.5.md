# UE Linux installed-build tooling v1.5

## Scope

This increment removes a repository-tooling assumption that every usable Linux Unreal Engine 5.8 installation exposes a root-level `GenerateProjectFiles.sh`.

The selected architecture remains unchanged:

- Unreal Engine 5.8 client and dedicated gameplay server
- Go 1.27 service plane
- PostgreSQL as durable authority
- Redis as ephemeral coordination / rate-limit state

## Supported Linux Unreal layouts

`tools/ue-linux.sh` accepts either:

1. a source-tree layout with executable `$UE_ROOT/GenerateProjectFiles.sh`, or
2. a precompiled installed-build layout with:
   - `$UE_ROOT/Engine/Build/BatchFiles/Linux/Build.sh`,
   - `$UE_ROOT/Engine/Binaries/DotNET/UnrealBuildTool/UnrealBuildTool.dll`, and
   - Unreal's bundled `Engine/Binaries/ThirdParty/DotNet/*/linux/dotnet` or a system `dotnet` fallback.

Installed-build project generation invokes:

```text
dotnet UnrealBuildTool.dll -projectfiles -project=<repo>/game/NeonDrive.uproject -game
```

Build targets continue to use Unreal's Linux `Build.sh`; packaging continues to use `RunUAT.sh BuildCookRun`.

## CI/static evidence

`tools/test-ue-linux.sh` creates disposable fake Unreal layouts and verifies:

- installed-build UnrealBuildTool project-generation arguments,
- source-tree `GenerateProjectFiles.sh` compatibility,
- NeonDriveClient Linux Development build invocation,
- Linux dedicated-server packaging invocation,
- UE root detection through configurable `UE_SEARCH_ROOTS`.

`make validate-control` runs syntax checks plus this fixture test. This is tooling-contract evidence only; it does not execute Unreal Engine itself.

## Security and authority boundary

The helper does not handle session tokens, player credentials, `GAME_SERVER_SHARED_KEY`, `ZNEON_GAME_SERVER_KEY`, PostgreSQL state, Redis state, quest rewards, race results, or inventory authority. It changes only how local Linux Unreal tooling is discovered and invoked.

## Non-claims

This work does **not** prove:

- a successful real UE 5.8 Client build,
- a successful real UE 5.8 dedicated Server build,
- cooking or packaged artifacts from the selected UE installation,
- live Unreal ↔ Go ticket/session redemption,
- live authoritative race lifecycle,
- Garage 17 / MQ001–MQ012 playability,
- production readiness.

Those gates remain open until retained executable evidence exists.
