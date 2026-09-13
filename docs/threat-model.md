# Threat Model

**Status:** baseline threat model. Mitigations are a mix of implemented prototype controls and production requirements.

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
- TLS,
- rotation,
- rate limits,
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
- quest completion uniqueness.

### Vehicle/build forgery
Risks:
- impossible part IDs,
- stale revision overwrite,
- using parts not owned/unlocked.

Controls:
- canonical catalog validation,
- optimistic expected revision,
- blueprint/inventory checks,
- transactional consumption/return.

### Race cheating
Risks:
- teleport/speed manipulation,
- checkpoint skipping,
- clock tampering,
- result fabrication.

Required before ranked release:
- authoritative checkpoint sequence,
- accepted build binding,
- impossible-state telemetry,
- server timing,
- result audit record,
- cheat response policy.

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

Required:
- rate limits,
- bounded request bodies,
- quotas/backpressure,
- capacity alerts,
- graceful degradation.

## Security exit criteria for production

Production security cannot be claimed until:
- transport security is deployed,
- threat tests pass,
- secrets are externally managed/rotatable,
- rate limits are verified,
- ranked abuse cases are tested,
- security findings are triaged/remediated,
- incident and rollback procedures are exercised.
