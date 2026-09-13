# Unreal Engine 5.8 Linux installed-build tooling

## Scope

`tools/ue-linux.sh` makes the repository's Linux Unreal control surface work with either a UE 5.8 source tree or a precompiled/installed UE 5.8 build.

The helper validates `Engine/Build/Build.version` and requires the 5.8 API line before generation, build, or package operations continue.

For project-file generation it prefers a source-tree `GenerateProjectFiles.sh`. If that script is absent, it uses `Engine/Binaries/DotNET/UnrealBuildTool/UnrealBuildTool.dll -projectfiles` with an engine-bundled Linux `dotnet` runtime when available, falling back to a system `dotnet` only when necessary.

Build and package operations still delegate to Unreal's canonical `Build.sh` and `RunUAT.sh BuildCookRun` entry points.

## Security / authority boundary

This tooling does not receive or embed `GAME_SERVER_SHARED_KEY` / `ZNEON_GAME_SERVER_KEY`. It changes build orchestration only. PostgreSQL remains durable authority, Unreal dedicated servers remain gameplay authority, and Redis remains ephemeral coordination.

## Evidence boundary

Repository tests use fake engine fixtures to verify source-tree and installed-build routing, 5.8 version rejection, target delegation, and separate Client/Server archive paths. These tests do not prove that the real NeonDrive Client or Server compiles, cooks, packages, or connects to Go.

The following gates therefore remain open until retained real-engine evidence exists:

- successful UE 5.8 Client/Server build artifacts;
- packaged Unreal ↔ Go session E2E;
- packaged Unreal ↔ Go race lifecycle E2E;
- Garage 17 / MQ001–MQ012 playable runtime evidence;
- physics-derived anti-cheat evidence.
