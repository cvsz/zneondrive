# Runtime Reference Quest Sequence Parity v3.7

## Scope

This increment closes a concrete semantic-drift gap between the executable Python reference authority and the Go/PostgreSQL service plane for the main-quest sequence contract.

The Go runtime already treats main-quest identifiers as the canonical `MQ001` through `MQ100` range and derives an immediate predecessor for every quest after `MQ001`. The PostgreSQL mutation path rejects a completion when that immediate predecessor is absent. The Python reference oracle previously modeled quest rewards and Garage 17 side effects but did not enforce the same identifier or prerequisite boundary.

v3.7 makes that rule executable on both sides:

- Python now parses only the canonical `MQ001`..`MQ100` identifier range before accepting a quest mutation;
- Python rejects an out-of-order completion before mutating quest, reward, inventory, or blueprint state;
- the shared fixture `tests/reference-parity-vectors.json` now contains quest identifier and immediate-prerequisite cases;
- Python exercises those vectors through the reference authority rather than a regex-only structural check;
- Go exercises the same prerequisite vectors through `ParseQuestID` and `PreviousQuestID`;
- existing reference-runtime tests now establish prerequisite state before testing idempotency and Garage 17 quest side effects.

## Authority and security boundary

This change does not move quest authority to Python or to the player client. The production service plane remains Go 1.27 + PostgreSQL, and the durable PostgreSQL quest mutation remains the authoritative implementation. Unreal Engine 5.8 remains the gameplay plane, and Redis remains ephemeral coordination only.

A client-provided quest identifier is still untrusted. The production Go/PostgreSQL path validates the identifier, derives the prerequisite server-side, checks durable completion state, and only then applies durable rewards and quest side effects. The Python implementation is a test/reference oracle used to detect contract drift.

## Evidence

Automated evidence now covers:

- `MQ001` as the root quest with no predecessor;
- valid predecessor derivation at the vertical-slice boundary (`MQ012` → `MQ011`) and campaign upper boundary (`MQ100` → `MQ099`);
- rejection of out-of-range and structurally invalid main-quest identifiers;
- rejection of `MQ002` without `MQ001`;
- acceptance of `MQ002` after `MQ001`;
- rejection of `MQ012` when `MQ011` is absent even if other earlier quest IDs are present;
- acceptance of `MQ012` when `MQ011` is present;
- no quest-state mutation on an out-of-order Python reference transition.

The existing PostgreSQL integration suite remains the production-path evidence for durable prerequisite enforcement and the MQ012 Roadworthy transition.

## Non-claims

This remains partial reference-parity evidence. It does **not** close the full Python↔Go reference-oracle parity gate. In particular, the Python oracle still accepts explicit reward amounts while the Go/PostgreSQL production path derives main-quest rewards server-side, and other production-only persistence/authentication semantics remain outside this reference increment.

It also does not prove a successful real Unreal Engine 5.8 Client or dedicated Server build, package/cook evidence, live packaged Unreal↔Go transport, Garage 17 playability, MQ001–MQ012 playable execution, live authoritative race gameplay, final-physics anti-cheat calibration, deployed ingress/observability, deployment-scale load/soak, HA/DR, or production deployment readiness.
