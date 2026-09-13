# Runtime Reference Race Input Parity v2.9

## Status

**Implemented reference/test evidence. Not live Unreal runtime evidence and not production-readiness proof.**

This increment extends the shared Python/Go reference-parity fixture to the bounded race-input validation contract already used by the Go 1.27 service plane.

## Evidence added

The shared `tests/reference-parity-vectors.json` fixture now covers:

- canonical `race_` identifier normalization;
- outer-whitespace trimming;
- allowed lower-case ASCII letters, digits, `_`, and `-`;
- rejection of missing prefixes, uppercase characters, undersized IDs, and path separators;
- checkpoint index boundaries `0..1024`;
- checkpoint elapsed-time boundaries `1..86400000` milliseconds.

Both `tests/test_reference_parity.py` and `services/game-api/internal/core/reference_parity_test.go` consume the same vectors. The Python side uses `zneondrive.race_contract`; the Go side uses the production-core `NormalizeRaceID` and `ValidateRaceCheckpoint` helpers.

## Trust boundary

This work does not move authority into Python and does not accept race outcomes from clients. The selected architecture remains:

- Unreal Engine 5.8 dedicated server: gameplay authority;
- Go 1.27: authenticated service/race contract plane;
- PostgreSQL: durable authority;
- Redis: ephemeral coordination and shared abuse-control state;
- Python: executable reference oracle and parity tests only.

## Non-claims / still open

This evidence does **not** prove:

- successful real UE 5.8 Client or dedicated Server build/package artifacts;
- live Unreal↔Go race start/checkpoint/finish transport;
- authoritative physics/checkpoint emission from a packaged dedicated server;
- full Python↔Go reference-oracle parity;
- Garage 17 / MQ001–MQ012 playable flow;
- deployment-scale load, soak, HA/DR, or production deployment.

The promoted repository runtime baseline therefore remains Phase 4.19 / Runtime v2.3 and production readiness remains evidence-gated.
