## Summary

Describe what changed and why.

## Validation

- [ ] `make ci` passes.
- [ ] Content IDs/references validate.
- [ ] Tests added/updated where behavior exists.
- [ ] Security impact reviewed.
- [ ] Documentation and catalogs updated together.

## Game-design compatibility

- [ ] No published stable ID was renamed without a migration plan.
- [ ] VIP/fairness rules remain intact.
- [ ] Pet changes do not add ranked competitive buffs.
- [ ] Vehicle/race changes remain compatible with server-authoritative validation.
- [ ] Story changes identify their world-state scope and relationship/faction consequences.

## Runtime / operational risk

Describe persistence, migration, networking, security, compatibility, live-ops, and rollback impact when applicable.

## Rollback

Describe how the change can be reverted or mitigated.

## Checklist

- [ ] Focused and reviewable.
- [ ] No credentials, private keys, secrets, personal player data, or sensitive telemetry included.
- [ ] Security/quality gates were not weakened.
- [ ] `CHANGELOG.md` updated for material changes.
