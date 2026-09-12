# Release

## Version classes

Until runtime implementation begins, releases represent **design/content contract versions**.

Suggested tags:
- `design-v0.1.0` — Game Design Bible/content contracts
- future runtime tags should adopt an explicit version policy after stack selection

## Design release checklist

1. CI and security checks pass.
2. `python3 tools/validate_design.py` passes.
3. Stable IDs were not accidentally renamed.
4. `CHANGELOG.md` is updated.
5. Canonical prose and machine-readable catalogs agree.
6. New authority/fairness decisions have ADR coverage.

## Runtime release evidence

A future playable/production release additionally requires:
- executable unit/integration/end-to-end tests,
- authoritative persistence/race validation evidence,
- security and anti-cheat review,
- load/soak evidence,
- backup and restore drill,
- rollback procedure,
- production deployment verification,
- privacy/data-retention review.

A Git tag or successful build alone is not production-readiness evidence.

## Rollback

Design releases can be reverted through version control while preserving published ID compatibility where possible.

Future runtime rollback documentation must cover code, content, schema/migration state, player-state compatibility, event/reward idempotency, and live-ops communication.
