# Runtime Reference Race Start Parity v4.1

## Scope

This increment centralizes the deterministic race-start eligibility rule that is applied after PostgreSQL resolves authoritative account/character/vehicle ownership and the exact active vehicle build.

Shared Python/Go vectors cover:

- canonical `race_*` identifier acceptance and normalization,
- Roadworthy gating,
- positive active build revision,
- non-empty authoritative build validation hash,
- rejection of malformed/uppercase race IDs and invalid build bindings.

The PostgreSQL transaction remains authoritative for ownership resolution, active-revision selection, row locking, idempotent operation binding, race-instance persistence, and subsequent checkpoint/result state. The Python implementation is reference-only.

## Evidence

- `services/game-api/internal/core/race.go` — deterministic `ValidateRaceStartBinding` contract.
- `services/game-api/internal/store/race_postgres.go` — PostgreSQL start transaction invokes the contract after authoritative build lookup.
- `tests/race-start-parity-v4.1.json` — shared vectors.
- `services/game-api/internal/core/race_start_parity_test.go` — Go executable parity.
- `src/zneondrive/race_start.py` and `tests/test_race_start_parity.py` — Python reference parity.

## Non-claims

This does **not** prove a packaged Unreal Engine 5.8 client or dedicated server can complete a live race start, does not prove gameplay-ticket transport over packaged Unreal↔Go, and does not close playable Garage 17/MQ001–MQ012, final-physics anti-cheat, load/soak, deployed observability, HA/DR, or production deployment gates.

Full reference-oracle parity remains open; v4.1 covers only the race-start eligibility/build-binding predicate.
