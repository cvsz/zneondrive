# Runtime Reference Race Result Parity v3.5

This increment extends the shared Python/Go reference-parity contract to the deterministic authoritative race-result digest already used by the Go 1.27 service plane.

## Evidence added

`tests/reference-parity-vectors.json` now carries shared race-result-hash vectors covering:

- a canonical authoritative result bound to race instance, account, character, vehicle, immutable build revision/hash, checkpoint count, and finish elapsed time;
- the maximum accepted checkpoint-count boundary (`1025`);
- rejection of zero checkpoint count;
- rejection of checkpoint-count overflow;
- rejection of non-positive finish elapsed time.

The Python reference oracle implements the same canonical JSON/SHA-256 digest contract in `src/zneondrive/race_contract.py`. Both `tests/test_reference_parity.py` and `services/game-api/internal/core/reference_parity_test.go` consume the same vectors, so drift in field ordering, bound identity/build evidence, acceptance limits, or hash output fails CI on either side.

## Trust boundary

No production authority moved. Unreal Engine 5.8 remains gameplay-authoritative, Go 1.27 remains the authenticated service/race plane, PostgreSQL remains durable authority, and Redis remains ephemeral coordination. Python remains a dependency-free reference/test oracle only.

The digest does not make player-submitted race outcomes authoritative. Production race start/checkpoint/finish writes remain restricted to the dedicated-server shared-key boundary and are persisted by the Go/PostgreSQL runtime.

## Non-claims

This evidence does **not** prove:

- successful real Unreal Engine 5.8 Client or Dedicated Server builds;
- packaged Unreal↔Go transport or gameplay-ticket redemption;
- live Unreal authoritative race start/checkpoint/finish execution;
- live Garage 17 / MQ001–MQ012 playable flow;
- final-physics anti-cheat calibration;
- full Python↔Go reference-oracle parity;
- deployed load/soak, ingress isolation, SLOs, HA/DR, or production readiness.

Those gates remain open until their required runtime/deployment evidence exists.
