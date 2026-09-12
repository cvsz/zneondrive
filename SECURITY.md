# Security Policy

Security is a first-class requirement for zNeonDrive because the product includes persistent ownership, inventory/economy, ranked competition, social systems and privileged live operations.

## Reporting a vulnerability

Do not disclose exploitable vulnerabilities in public issues, pull requests, discussions, commit messages, screenshots or logs.

Use GitHub private vulnerability reporting/security advisories when available, or an agreed private channel with the repository owner.

Useful reports include affected commit/version, reproduction steps, prerequisites, impact, exploitability and suggested remediation.

## Current supported repository surface

The public repository currently contains:
- Unreal Engine 5.8 C++ gameplay/source integration code,
- Go 1.27 HTTP service-plane code,
- PostgreSQL migrations and integration tests,
- Docker Compose development runtime,
- Python reference oracle,
- design/content catalogs and schemas,
- client presentation demo,
- CI/CodeQL/dependency review workflows,
- validation and documentation tooling.

This is a prototype/development support boundary, not a statement that a public production game service is available.

## Implemented security controls

Current source/tests include:
- hashed resume/session credentials,
- one-time hashed gameplay tickets,
- gameplay-ticket TTL and atomic redemption,
- server-only shared-key redemption boundary,
- constant-time server-key comparison,
- operation-ID idempotency semantics,
- optimistic vehicle build revision checks,
- canonical part validation,
- inventory/blueprint authorization checks,
- non-negative durable inventory constraints,
- client input treated as untrusted in Unreal prototype authority flow.

## Mandatory security expectations

- Player clients, client clocks, client-computed rewards, build legality and race results are untrusted.
- Ownership, inventory/economy, quest completion, crew authorization, VIP entitlements and ranked results must be authoritative.
- Rewardable state transitions require replay/idempotency protection.
- Privileged administration/live-ops actions require auditability.
- Never commit credentials, tokens, private keys, production secrets, personal player data or sensitive telemetry.
- Validate untrusted input and authorize every state-changing action.
- Use least-privilege GitHub Actions and service permissions.
- Do not weaken security gates merely to obtain a passing build.

## Production gaps

Before production/ranked release, explicitly verify:
- TLS and service identity,
- rate limiting/backpressure,
- secret rotation,
- fabricated inventory/reward attempts,
- timing/checkpoint and impossible movement manipulation,
- reconnect abuse,
- crew/role privilege escalation,
- entitlement/VIP spoofing,
- leaderboard/result tampering,
- admin/live-ops misuse,
- DoS/capacity behavior,
- security incident exercises.

See:
- [Threat Model](./docs/threat-model.md)
- [Privacy & Data Retention](./docs/privacy-data-retention.md)
- [Incident Response](./docs/incident-response.md)
- [Backup, Restore & DR](./docs/backup-restore-dr.md)

## Incident handling

Future production operations must document and exercise containment, impact assessment, secret rotation, rollback, economy/result correction, evidence preservation, disclosure and recovery validation.
