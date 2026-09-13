# UE 5.8 build evidence workflow v1.6

## Purpose

The manual `Unreal Build Evidence` workflow is the repository-controlled path for collecting retained evidence from a real self-hosted Linux runner with Unreal Engine 5.8 installed.

It preserves the selected architecture and authority model: Unreal Engine 5.8 remains the gameplay client/server runtime, Go 1.27 remains the service plane, PostgreSQL remains durable authority, and Redis remains ephemeral coordination.

## Execution contract

The workflow runs only through `workflow_dispatch` on a runner labeled:

```text
self-hosted, linux, unreal-5.8
```

The runner must provide `UE_ROOT`. The workflow then:

1. runs repository/static validators and the fake-engine UE tooling unit suite;
2. records repository SHA, run identity, runner metadata, and UTC timestamps;
3. validates the real engine through `tools/ue-linux.sh info` and retains `Engine/Build/Build.version`;
4. generates project files through the source/installed-build-aware helper;
5. builds `NeonDriveClient` in Linux Development;
6. builds `NeonDriveServer` in Linux Development;
7. inventories matching project binaries and records SHA-256 hashes;
8. uploads logs and metadata as a 30-day Actions artifact, including on failure.

The artifact is deliberately compact. It retains build provenance and diagnostics without attempting to upload multi-gigabyte cooked game packages.

## Evidence interpretation

A green repository PR that adds this workflow proves only the workflow/static contract. It does **not** prove Unreal builds succeeded.

The UE build-evidence gate may be advanced only after a real manual workflow run finishes successfully and retained evidence shows both Client and Server targets were produced from the expected UE 5.8 installation and repository commit.

Packaging is a separate gate. `make client-package-linux`, `make game-server-package-linux`, and `make package-all-linux` remain available, but this build-evidence workflow does not silently claim package evidence.

## Security boundary

The workflow does not require or embed player credentials, gameplay tickets, `GAME_SERVER_SHARED_KEY`, `ZNEON_GAME_SERVER_KEY`, PostgreSQL URLs, or Redis URLs. It exercises compilation only and does not weaken server-authoritative runtime boundaries.

## Still open

Even after a successful Client/Server build-evidence run, the following remain separate evidence gates:

- successful retained Client/Server package/cook evidence;
- live packaged Unreal ↔ Go bootstrap/resume/ticket redemption;
- live authoritative Unreal ↔ Go race lifecycle;
- Garage 17 / MQ001–MQ012 playable flow;
- physics-derived race anti-cheat;
- deployed ingress evidence;
- deployment-scale load/soak and measured SLOs;
- production backup/PITR/RPO/RTO, HA/DR, and deployment evidence.
