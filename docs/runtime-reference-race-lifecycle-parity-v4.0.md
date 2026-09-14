# Runtime Reference Race Lifecycle Parity v4.0

## Scope

This increment centralizes deterministic authoritative race-lifecycle acceptance rules in the Go core and exercises the same contract through a Python reference oracle and shared JSON vectors.

Covered rules:

- checkpoints are accepted only while the race instance is `active`;
- the submitted checkpoint index must equal the authoritative `next_checkpoint` cursor;
- checkpoint elapsed time must be positive, no more than 24 hours, and strictly greater than the last accepted elapsed time;
- finish is accepted only while the instance is `active`;
- finish checkpoint count must exactly equal the authoritative cursor and be at least one;
- finish elapsed time must be strictly greater than the last accepted checkpoint time;
- the PostgreSQL store calls the centralized Go validators before durable checkpoint/result mutation.

The shared vectors live in `tests/race-lifecycle-parity-v4.0.json` and are executed by both `tests/test_race_lifecycle_parity.py` and `services/game-api/internal/core/race_lifecycle_parity_test.go`.

## Authority boundary

This is a parity/refactoring increment, not an authority migration. Unreal Engine 5.8 Dedicated Server remains gameplay authority, Go 1.27 remains the authenticated service plane, PostgreSQL remains durable authority, and Redis remains ephemeral/distributed coordination. Python is reference/test code only.

PostgreSQL still owns transaction locking, operation-id idempotency, checkpoint persistence, state transition to `finished`, and final result persistence. The extracted Go functions only centralize deterministic acceptance predicates that were previously embedded directly in the store.

## Evidence claims

This increment can claim source/unit/integration-compatible parity evidence for checkpoint-order and finish-order semantics once the applicable CI gates pass on the PR head.

It does **not** claim:

- a live packaged Unreal ↔ Go race lifecycle;
- real UE 5.8 Client/Server build or cook/package evidence;
- live physics-derived anti-cheat calibration;
- deployed ingress, load/soak, metrics/SLO, HA/DR, or production deployment evidence;
- completion of the full reference-oracle parity gate.

Those gates remain open until their own retained runtime/deployment evidence exists.
