# Runtime Reference Build Operation Replay v4.4

**Status:** source/integration-contract hardening only. This does not close live Unreal, full reference-parity, deployment, or production-readiness gates.

## Change

The Go/PostgreSQL rebuild replay path now binds an existing operation receipt to the same authoritative character owner, vehicle, expected source revision, resulting revision, and canonical build validation hash before treating a repeated operation ID as an idempotent replay.

Previously the replay lookup compared only vehicle ID and validation hash. A repeated key with the same vehicle/hash but a different `expectedRevision` could therefore be accepted as a replay even though the Python reference oracle correctly treated `expected_revision` as part of the operation payload. The hardened contract requires `existing.result_revision == expectedRevision + 1` and rejects owner, vehicle, revision, or hash drift.

## Evidence

- `services/game-api/internal/core/build_operation.go` defines the deterministic replay validator.
- `services/game-api/internal/store/postgres.go` resolves the durable receipt together with the authoritative vehicle owner and invokes the validator before returning a replay.
- `tests/build-operation-replay-parity-v4.4.json` is consumed by both Go and Python tests.
- negative vectors cover changed owner, vehicle, expected revision, and validation hash.

## Trust boundary

Unreal Engine 5.8 dedicated servers remain gameplay authority; Go 1.27 remains the authenticated service plane; PostgreSQL remains durable authority; Redis remains ephemeral coordination; Python remains a test/reference oracle only.

## Non-claims

This increment does not prove successful UE Client/Server builds, cook/package artifacts, live packaged Unreal↔Go transport, Garage 17/MQ001–MQ012 playability, live race physics/anti-cheat calibration, deployment-scale load/soak, HA/DR, or overall production readiness. Full reference-oracle parity remains open.
