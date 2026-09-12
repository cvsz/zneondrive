# Release

## Version classes

The repository contains both design/content contracts and runtime prototype implementation.

Recommended version classes:
- `design-vX.Y.Z` — design/content contract milestones,
- `runtime-vX.Y.Z` — prototype/runtime implementation milestones,
- future playable/alpha/beta/release tags only when their defined evidence gates are satisfied.

A tag name must not imply production readiness by itself.

## Repository release checks

1. CI and applicable runtime workflows pass.
2. Design/runtime/documentation validators pass.
3. Go tests and integration tests pass where affected.
4. Unreal source-build evidence exists when claiming an Unreal build.
5. Stable IDs were not accidentally renamed.
6. `CHANGELOG.md`, README, ROADMAP and checklist reflect reality.
7. Canonical prose and machine-readable catalogs agree.
8. New authority/fairness decisions have ADR coverage.

## Runtime release evidence

A playable/production release additionally requires:
- packaged Unreal client/server build evidence,
- live Unreal↔Go integration/E2E,
- live authoritative race transport/result evidence,
- security and anti-cheat review,
- performance/load/soak evidence,
- backup and restore drill,
- accepted RPO/RTO,
- rollback procedure and exercise,
- observability/alert validation,
- privacy/data-retention review,
- support/moderation/live-ops readiness where applicable,
- production deployment verification.

See [Release Readiness Checklist](./release-readiness-checklist.md).

## Rollback

Rollback planning must cover:
- executable code,
- content/configuration,
- schema/migrations,
- player-state compatibility,
- race/result/economy idempotency,
- live-ops communication.

Destructive database rollback should not be assumed safe. Prefer forward fixes or restore strategies supported by tested backups and migration evidence.

A successful build alone is not production-readiness evidence.
