# Runtime Reference Entitlement Parity v3.6

## Scope

This increment adds executable shared-vector parity for two existing reference-authority invariants that previously lacked an explicit Go-side parity contract:

- garage-slot capacity is an entitlement boundary only: a vehicle grant is accepted only while the character owns fewer vehicles than the account's garage-slot entitlement;
- starter-lineage vehicles cannot be removed by the routine deletion path.

The shared fixture is `tests/entitlement-parity-vectors.json`. Python exercises the existing `WorldState` reference oracle and Go exercises small pure policy helpers in `services/game-api/internal/core/entitlements.go` from the same vectors.

## Security and fairness boundary

This does **not** make VIP capacity part of race performance, build validation, or result hashing. VIP may expand garage capacity, but the authoritative competitive build/race contracts remain independent of that capacity. Starter-lineage protection is a deletion invariant only; it does not grant race authority or bypass ownership checks.

## Evidence

Automated evidence in this repository now covers:

- free-tier first-vehicle allowance and second-vehicle capacity rejection;
- larger entitled garage capacity while still rejecting grants at the configured limit;
- rejection of routine deletion for starter-lineage vehicles;
- allowance of routine deletion for non-starter vehicles;
- fixture-integrity checks on both Python and Go paths.

## Non-claims

This is partial reference-parity evidence. It does not close the full Python↔Go reference-oracle parity gate, does not add a production entitlement persistence/billing system, and does not prove live Unreal↔Go transport, real UE 5.8 build/package evidence, Garage 17 playability, authoritative live race execution, deployed anti-cheat, load/soak, HA/DR, or production deployment readiness.
