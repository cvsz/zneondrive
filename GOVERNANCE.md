# Governance

## Project ownership

PROJECT: NEON DRIVE is maintained under the `cvsz/zneondrive` repository. The repository owner is the final approver for product, architecture, security, release, and IP decisions unless a later governance agreement replaces this document.

## Decision classes

### Product/content decisions
Examples: story canon, faction behavior, economy rules, VIP policy, quest IDs.

Required evidence:
- design document/catalog update,
- compatibility review for stable IDs,
- fairness review when competitive impact exists.

### Architecture decisions
Material architecture changes require an ADR under `docs/adr/`.

Examples:
- engine/runtime selection,
- data authority,
- persistence,
- networking,
- deployment topology,
- security boundary changes.

### Security decisions
Security gates may not be bypassed solely to make CI or a release pass. Vulnerabilities are handled through private disclosure as described in `SECURITY.md`.

### Release decisions
A release status is based on evidence, not intent. A green build alone is not production-readiness evidence.

## Change path

1. Open an issue for non-trivial work.
2. Implement on a focused branch.
3. Run applicable validation and tests.
4. Open a PR describing scope, evidence, risks, and remaining gaps.
5. Merge only when required checks pass.
6. Synchronize README, ROADMAP, checklist, changelog, and docs when status changes.

## Stable-contract rule

Published content IDs and durable API/database semantics are compatibility contracts. Renaming or destructive changes require migration/compatibility planning.

## Conflict resolution

When documents disagree, use this priority:
1. accepted ADRs for architecture/security,
2. machine-readable catalogs/schemas for stable IDs/contracts,
3. current runtime behavior with passing tests,
4. canonical Game Design Bible,
5. roadmap/planning documents.

Conflicts should be corrected rather than silently worked around.
