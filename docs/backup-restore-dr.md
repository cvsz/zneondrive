# Backup, Restore & Disaster Recovery

**Status:** required plan; production restore/DR is not yet verified.

## Scope

Durable assets requiring protection include:
- PostgreSQL player/account/progression state,
- migration history,
- canonical content/configuration required to interpret durable IDs,
- live-ops configuration when introduced,
- audit records required by policy.

Redis must not become the only copy of irreplaceable player state.

## Backup target

For production design:
- automated PostgreSQL backups,
- encrypted storage,
- retention policy,
- off-host/off-node copy,
- periodic restore validation,
- access controls and auditability.

## Initial recovery objectives

Planning targets only:
- **RPO:** ≤ 15 minutes for durable player state where architecture supports it,
- **RTO:** ≤ 2 hours for service restoration.

Final values require business approval and architecture evidence.

## Restore drill

A restore drill is successful only if it:
1. restores to an isolated environment,
2. runs schema/application compatibility checks,
3. verifies representative accounts, vehicles, builds, inventory and quests,
4. verifies operation-ID/idempotency behavior,
5. records duration and data-loss window,
6. captures failures and remediation.

## Disaster scenarios to test

- accidental data deletion,
- failed migration,
- corrupted primary database,
- lost service node,
- gameplay-server fleet loss,
- region/provider outage when multi-region becomes a requirement,
- leaked credential requiring rotation.

## Evidence rule

Backup existence is not restore evidence. Production DR remains open until restore drills and declared RPO/RTO have been demonstrated.
