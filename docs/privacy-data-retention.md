# Privacy & Data Retention

**Status:** engineering baseline; not legal advice and not a substitute for jurisdiction-specific counsel.

## Data minimization

Collect only data required for:
- authentication/account continuity,
- gameplay/progression,
- security/anti-cheat,
- support/moderation,
- reliability/analytics with defined purpose.

Avoid collecting sensitive personal data unless a documented requirement exists.

## Credential handling

Never log:
- resume keys,
- session tokens,
- gameplay tickets,
- server shared secrets,
- private keys.

Durable credentials should remain hashed where the design supports one-way verification.

## Data classes

### Durable gameplay
Account ID, character, vehicles, builds, inventory, quests, progression.

### Security/audit
Authentication/security events, privileged operations, anti-cheat evidence.

### Telemetry
Performance, errors, gameplay events, device/platform metadata where justified.

### Player communications
Chat/report/moderation evidence when those systems are introduced.

## Retention

Before production, define per-class:
- purpose,
- retention period,
- deletion/anonymization behavior,
- access roles,
- backup retention interaction,
- legal/security hold behavior.

## Player rights workflow

If required by applicable law/product region, support verified:
- access/export,
- correction,
- deletion,
- objection/consent controls.

Deletion must account for fraud/security/audit/legal retention requirements and backups.

## Environment rule

Production player data must not be copied into development/test environments without an approved anonymization process.
