# NeonDrive Unreal Runtime

## Required engine

Use Unreal Engine 5.8.x. The repository targets the UE 5.8 API line; prefer the current 5.8 hotfix after regression testing.

## Generate project files

Linux source build:

```bash
"$UE_ROOT/GenerateProjectFiles.sh" -project="$PWD/game/NeonDrive.uproject" -game
```

Windows:

```powershell
& "$env:UE_ROOT\Engine\Build\BatchFiles\GenerateProjectFiles.bat" -project="$PWD\game\NeonDrive.uproject" -game
```

## Build

Linux:

```bash
"$UE_ROOT/Engine/Build/BatchFiles/Linux/Build.sh" NeonDriveEditor Linux Development "$PWD/game/NeonDrive.uproject" -WaitMutex
"$UE_ROOT/Engine/Build/BatchFiles/Linux/Build.sh" NeonDriveServer Linux Development "$PWD/game/NeonDrive.uproject" -WaitMutex
```

## Prototype driving

Default controls:
- W/S: throttle/reverse
- A/D: steering

The pawn is intentionally simple. Client input is clamped on the server and movement is performed only by authority, then replicated by Unreal. This proves the authority direction; it is not final vehicle physics or anti-cheat.

## Dedicated server CI

`.github/workflows/unreal-source-build.yml` is manual and requires a self-hosted Linux runner with:
- label `unreal-5.8`,
- `UE_ROOT` pointing at a built UE 5.8 source tree,
- sufficient disk/RAM for Unreal C++ builds.

A workflow file is not build evidence by itself. Keep the roadmap build-evidence item open until a run succeeds and artifacts/logs are retained.
