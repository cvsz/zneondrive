# Runtime Reference Operation Parity v2.6

**Status:** reference-oracle parity hardening only. This does not close full Python↔Go parity, live Unreal integration, or production-readiness gates.

## Scope

The Python executable reference authority now binds every idempotent mutation receipt to a deterministic semantic fingerprint containing the mutation kind and canonical payload. Replaying the same operation ID with the same semantic payload returns the original receipt. Reusing that operation ID for a different mutation, changed reward, changed vehicle/build payload, or changed race timing is rejected with `OperationConflictError` before authoritative state can change.

This aligns the executable reference behavior with the existing Go/PostgreSQL runtime contract for payload-bound idempotency and removes a weak-oracle case where a changed request could previously receive an unrelated cached receipt.

Canonicalized vehicle part sets remain order- and duplicate-insensitive before fingerprinting, matching build-hash semantics.

## Evidence added

Reference-runtime unit coverage now verifies:

- build retries with reordered equivalent parts replay safely without advancing the immutable revision;
- build operation IDs reject changed part payloads;
- quest operation IDs reject changed reward payloads without changing money/XP/reputation;
- vehicle-grant operation IDs reject changed vehicle identity;
- race-result operation IDs reject changed checkpoint timing payloads;
- one operation ID cannot be replayed across different mutation kinds.

## Trust boundary

No production authority moved. Unreal Engine 5.8 remains gameplay-authoritative, Go 1.27 remains the service/auth boundary, PostgreSQL remains durable authority, and Redis remains ephemeral coordination. The Python model remains a dependency-free reference oracle and is not promoted into the production data path.

## Evidence boundary

This increment narrows one documented full-reference-parity gap but does **not** prove complete Python↔Go parity. Account/character uniqueness, entitlement/garage capacity, starter deletion protection, complete inventory/blueprint/rebuild behavior, the entire race lifecycle, and live Unreal transport remain to be compared or exercised. Real UE 5.8 Client/Server build/package artifacts, live Unreal↔Go E2E, deployment-scale load/soak, HA/DR, and production deployment evidence also remain open.
