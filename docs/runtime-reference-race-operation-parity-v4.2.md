# Runtime Reference Race Operation Parity v4.2

## Scope

This increment centralizes the deterministic replay-equality rules applied after PostgreSQL finds an existing race operation-id binding.

Shared Python/Go vectors cover all three authoritative race mutations:

- race start: retry must preserve account, vehicle and canonical race ID,
- checkpoint: retry must preserve race instance, checkpoint index and elapsed time,
- finish: retry must preserve race instance, checkpoint count and finish elapsed time,
- reuse of an operation ID with any changed bound payload is rejected as an operation-id conflict.

PostgreSQL remains authoritative for operation-id uniqueness, transaction locking, durable race/checkpoint/result persistence and replay lookup. The extracted Go helpers only centralize deterministic payload equality; the Python implementation is reference-only.

## Evidence

- `services/game-api/internal/core/race_operation.go` — deterministic replay-equality contracts.
- `services/game-api/internal/store/race_postgres.go` — durable start/checkpoint/finish replay paths call the centralized contracts.
- `tests/race-operation-parity-v4.2.json` — shared start/checkpoint/finish vectors.
- `services/game-api/internal/core/race_operation_parity_test.go` — Go executable parity.
- `src/zneondrive/race_operation.py` and `tests/test_race_operation_parity.py` — Python reference parity.

## Non-claims

This does **not** prove operation replay over packaged Unreal↔Go transport, does not prove live dedicated-server retry behavior under network partitions, and does not close the full reference-oracle parity gate. It also does not close real UE 5.8 build/package, playable Garage 17/MQ001–MQ012, deployment-scale load/soak, final-physics anti-cheat, deployed observability, HA/DR or production deployment gates.
