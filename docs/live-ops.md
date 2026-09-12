# Live Operations

**Status:** target operational design; live-ops tooling is not yet implemented.

## Live-ops scope

Future tools may manage:
- seasonal events,
- NOVA GRAND PRIX schedules,
- world-event rotations,
- reward tables,
- announcements,
- feature flags,
- moderation/support actions,
- emergency disable switches.

## Safety rules

High-impact operations require:
- authenticated operator identity,
- RBAC,
- reason/ticket reference,
- audit record,
- preview/dry-run where practical,
- rollback path,
- separation of duties for economy-wide actions when scale justifies it.

## Event lifecycle

1. author,
2. validate,
3. stage,
4. approve,
5. schedule,
6. activate,
7. monitor,
8. close,
9. reconcile rewards/results,
10. review.

## Economy controls

No live-ops tool should bypass idempotency/ownership rules. Emergency grants/compensation must be attributable and auditable.

## Kill switches

Production should support narrowly scoped disable controls for risky subsystems such as:
- a broken reward path,
- a race/event queue,
- a problematic content rotation,
- a compromised commerce/entitlement path.

## Evidence

Live-ops readiness requires real operator tooling, audit records, staging exercises and rollback tests.
