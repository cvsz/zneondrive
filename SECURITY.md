# Security Policy

Security is a first-class requirement for zNeonDrive, especially because the future product is a persistent online game with ownership, inventory/economy, ranked competition, crews, social systems, and privileged live operations.

## Reporting a vulnerability

Do not disclose exploitable vulnerabilities in public issues, pull requests, discussions, or commit messages.

Use GitHub private vulnerability reporting/security advisories when available, or an agreed private channel with the repository owner.

Useful reports include affected commit/version, reproduction steps, prerequisites, impact, exploitability, and suggested remediation.

## Current supported surface

The current public repository primarily contains design documents, content catalogs, schemas, CI, and design-validation tooling.

When runtime code is added, the supported-version table must be expanded to cover client/server and deployed release channels.

## Mandatory security expectations

- Player clients, client clocks, client-computed rewards, build legality, and race results are untrusted.
- Ownership, inventory/economy, quest completion, crew authorization, VIP entitlements, and ranked results must be authoritative.
- Rewardable state transitions require replay/idempotency protection.
- Privileged administration/live-ops actions require auditability.
- Never commit credentials, tokens, private keys, production secrets, personal player data, or sensitive telemetry.
- Validate untrusted input and authorize every state-changing action.
- Use least-privilege GitHub Actions permissions and service permissions.
- Keep dependency/security scanning enabled when relevant.
- Do not weaken security gates merely to obtain a passing build.

## Online-game threat model backlog

Before multiplayer production release, explicitly test:
- fabricated inventory/reward requests,
- replay/duplicate grants,
- impossible vehicle builds,
- timing/checkpoint manipulation,
- speed/teleport impossible states,
- disconnect/reconnect abuse,
- crew/role privilege escalation,
- entitlement/VIP spoofing,
- leaderboard/result tampering,
- admin/live-ops misuse,
- denial-of-service and rate-limit behavior.

## Incident handling

Future production operations must document containment, player-impact assessment, credential/secret rotation where applicable, rollback, economy/result correction, evidence preservation, disclosure, and recovery validation.
