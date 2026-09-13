# Threat Model

**Status:** baseline threat model with a CI-exercised partial abuse matrix. Production security remains evidence-gated.

See [Security Threat Exercise v3.0](./security-threat-exercise-v3.0.md) for the machine-checked mapping between selected abuse scenarios and current CI evidence. The exercise deliberately leaves deployed/live blockers open and does not claim production security readiness.

## Assets to protect

- account/session credentials,
- resume keys,
- gameplay tickets,
- vehicle ownership/build history,
- inventory/blueprints,
- money/XP/reputation,
- quest completion/rewards,
- race eligibility/results,
- VIP/entitlements,
- crew privileges,
- live-ops/admin capability,
- player personal data.

## Trust boundaries

Untrusted:
- player client,
- local save/config files,
- client clock,
- user-provided telemetry,
- public network traffic.

Conditionally trusted:
- authenticated dedicated gameplay server,
- internal service calls,
- CI/deployment identities,
- live-ops operators.

Authoritative:
- validated gameplay-server session outcomes for moment-to-moment play,
- service-plane mutation rules,
- PostgreSQL durable state.

## Primary abuse cases

### Credential theft/replay
Risks:
- stolen resume key/session token,
- gameplay ticket replay,
- leaked game-server key.

Current controls:
- hashed durable credentials,
- ticket TTL,
- single-use ticket consumption,
- constant-time game-server key comparison.

Still required:
- deployed TLS evidence,
- credential/key rotation evidence,
- anomaly detection,
- scoped service identities.

### Economy/reward duplication
Risks:
- retry storms,
- replayed quest completion,
- duplicated reconnect grants.

Controls:
- operation IDs,
- unique constraints,
- quest completion uniqueness,
- payload-bound idempotency semantics.

### Vehicle/build forgery
Risks:
- impossible part IDs,
- stale revision overwrite,
- using parts not owned/unlocked.

Controls:
- canonical catalog validation,
- optimistic expected revision,
- blueprint/inventory checks,
- transactional consumption/return,
- immutable build revisions.

### Race cheating
Risks:
- teleport/speed manipulation,
- checkpoint skipping,
- clock tampering,
- result fabrication.

Current source/CI controls:
- authoritative checkpoint sequence,
- accepted build binding,
- server timing for accepted checkpoint progression,
- deterministic build-bound result records,
- source-level impossible displacement/rotation telemetry envelopes.

Required before ranked release:
- live final-physics calibration,
- false-positive/false-negative evaluation,
- ranked sanctions policy,
- live packaged Unreal↔Go race evidence.

### Privilege abuse
Risks:
- crew role escalation,
- admin/live-ops misuse,
- unauthorized reward grants.

Required:
- RBAC,
- least privilege,
- immutable/auditable privileged actions,
- approval for high-impact economy actions.

### Availability attacks
Risks:
- request floods,
- session/ticket abuse,
- expensive query patterns,
- gameplay-instance exhaustion.

Current source/CI controls:
- bounded local token buckets,
- Redis-coordinated distributed limiter budgets,
- bounded local fallback when Redis is unavailable,
- bounded limiter state/cardinality,
- health-probe exemption.

Still required:
- deployment-scale load evidence,
- quotas/backpressure on the target environment,
- capacity alerts,
- graceful degradation evidence under target failure modes.

## Security exit criteria for production

Production security cannot be claimed until:
- transport security is deployed and verified,
- threat tests cover live Unreal↔Go and deployed ingress boundaries,
- secrets are externally managed/rotatable,
- rate limits are verified under deployment-scale load,
- ranked abuse cases are tested against final physics,
- privileged/admin actions are auditable,
- security findings are triaged/remediated,
- incident and rollback procedures are exercised.
