# Incident Response

## Severity

### SEV-1
Major security compromise, widespread inability to play, durable player-state corruption, uncontrolled economy duplication, or loss of administrative control.

### SEV-2
Significant degradation, regional/feature outage, elevated failure rate, recoverable progression impact.

### SEV-3
Limited defect with workaround and low player/business impact.

## Response flow

1. **Detect** — alert/report creates an incident record.
2. **Triage** — identify severity, scope, affected release/environment.
3. **Contain** — disable dangerous mutations/events or isolate compromised components.
4. **Preserve evidence** — logs, traces, deployment SHA, database/audit context.
5. **Recover** — rollback, restore, rotate credentials, or deploy fix.
6. **Validate** — prove player state/economy/race integrity before reopening.
7. **Communicate** — internal/client/player updates appropriate to impact.
8. **Review** — blameless post-incident review with concrete actions.

## Security-specific actions

When credentials may be exposed:
- rotate affected keys/tokens,
- invalidate sessions/tickets where feasible,
- audit privileged actions,
- preserve evidence,
- follow vulnerability disclosure obligations.

## Economy/ranked integrity

If an incident can create duplicated rewards or invalid ranked outcomes:
- stop the affected grant/result path,
- identify exact mutation/result set,
- prefer auditable corrective transactions over silent edits,
- document player-impact remediation.

## Minimum incident record

- start/end time,
- severity,
- detection source,
- affected services/builds,
- customer/player impact,
- root/contributing causes,
- actions taken,
- recovery evidence,
- follow-up owners.

No production incident-response claim should be made until the runbook has been exercised in at least a staging/game-day scenario.
