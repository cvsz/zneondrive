# Unreal Package Evidence v1.7

## Purpose

This gate prepares retained evidence for real Unreal Engine 5.8 Linux Client and dedicated Server cook/stage/package runs on the self-hosted `unreal-5.8` runner.

The workflow uses the repository's installed-build-aware `tools/ue-linux.sh` helper and preserves separate Client and Server archive roots.

## Retained evidence

A workflow run retains for 30 days:

- repository commit and workflow/run identity;
- runner identity and UTC timestamps;
- Unreal `Build.version` and helper-reported engine mode;
- project-generation log;
- Client package log;
- dedicated Server package log;
- Client and Server file manifests with byte sizes;
- Client and Server SHA-256 checksum inventories;
- total byte counts for each package root.

The workflow intentionally uploads compact evidence metadata rather than multi-gigabyte packaged game payloads.

On a successful run the Client and Server manifests and checksum inventories must all be non-empty. Failure runs still retain whatever evidence was produced before failure.

## Security and authority boundary

This is build/cook/package evidence only. It does not receive or embed player credentials, `GAME_SERVER_SHARED_KEY`, `ZNEON_GAME_SERVER_KEY`, PostgreSQL URLs, or Redis URLs.

Unreal dedicated servers remain gameplay authority, PostgreSQL remains durable authority, Go 1.27 remains the service plane, and Redis remains ephemeral coordination.

## Non-claims

Adding this workflow does not prove that a real package run succeeded. The package evidence gate remains open until a real self-hosted UE 5.8 workflow run succeeds and retained manifests/checksums exist.

It also does not prove:

- packaged Unreal ↔ Go session/bootstrap E2E;
- live authoritative race transport;
- Garage 17 / MQ001–MQ012 playability;
- physics-derived anti-cheat;
- deployment-scale load/soak;
- deployed SLOs/observability;
- production backup/restore, HA, DR, or production deployment.
