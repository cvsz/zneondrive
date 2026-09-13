# Runtime Reference Parity Integrity v2.5

This increment hardens the shared Python/Go reference-parity fixture without expanding the production-readiness claim.

## Evidence added

Both the Python reference-oracle tests and Go 1.27 core tests now reject malformed shared parity fixtures before semantic assertions run. The fixture must contain non-empty build-hash and quest-ID case sets, unique case identifiers, non-empty part lists, canonical 64-character lowercase SHA-256 build hashes, unique quest IDs, and no predecessor metadata on invalid quest vectors.

This reduces the risk of a parity test silently becoming weak because of duplicated, empty, or malformed vectors. The existing semantic checks remain unchanged: Python and Go must still agree on canonical starter build hashes, duplicate/order normalization, valid MQ001–MQ100 bounds, and deterministic predecessor mapping.

## Trust boundary

No runtime authority moved. Unreal Engine 5.8 remains gameplay-authoritative, Go remains the service/auth boundary, PostgreSQL remains durable authority, and Redis remains ephemeral coordination. These tests operate only on repository reference vectors.

## Evidence boundary

This does **not** close full reference-oracle parity. Account/character uniqueness, entitlement capacity, starter deletion protection, complete quest reward/idempotency semantics, build-conflict semantics, inventory/blueprint/rebuild behavior, race lifecycle semantics, and live Unreal transport still require additional cross-language or runtime evidence.

It also does not provide real UE 5.8 build/package evidence, live Unreal↔Go integration evidence, deployment-scale load/soak, HA/DR, or production deployment evidence.
