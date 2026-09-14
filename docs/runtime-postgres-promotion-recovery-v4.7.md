# Runtime PostgreSQL Promotion Recovery Evidence v4.7

## Scope

This evidence increment exercises an isolated PostgreSQL 17 physical-streaming-replication topology in Runtime Go CI and validates application behavior across a real database authority transition.

The selected architecture remains unchanged:

- Unreal Engine 5.8 dedicated server is the gameplay authority.
- Go 1.27 is the authenticated service plane.
- PostgreSQL is the durable source of truth.
- Redis remains ephemeral coordination and distributed abuse-control state.

## Executable evidence

`services/game-api/internal/store/postgres_promotion_recovery_integration_test.go` is compiled only with the `integration` build tag. During the existing Runtime Go PostgreSQL/Redis integration gate it:

1. creates an isolated Docker network, PostgreSQL 17 primary and persistent standby volume;
2. enables physical WAL streaming on the isolated primary;
3. creates a replication-only role and a test-network-only replication trust rule;
4. takes a physical `pg_basebackup -R` and starts a hot standby;
5. requires `pg_stat_replication` to report an active streaming standby;
6. runs the real Go/PostgreSQL store schema, bootstrap, session and MQ001 mutation path against the primary;
7. waits until the standby has replayed at least the primary WAL LSN containing the durable mutation;
8. stops the primary and promotes the standby with `pg_ctl promote`;
9. reconnects the real Go store to the promoted PostgreSQL instance;
10. resumes the same authoritative account/character/vehicle from the retained PostgreSQL state;
11. replays the same character-scoped MQ001 operation ID and requires `Applied=false`, zero replay reward and exactly one durable MQ001 completion;
12. completes MQ002 after promotion and requires a new authoritative mutation to persist successfully.

The test fails in CI if Docker is unavailable rather than silently skipping the promotion evidence gate.

## Security and trust-boundary notes

The physical-replication trust rule exists only inside the ephemeral isolated Docker network used by the CI test. Player-facing or production database authentication policy is not changed.

No client-provided reward, identity, build, inventory, blueprint or race state becomes authoritative. The post-promotion store resolves durable identity and progression from PostgreSQL and continues using the existing authoritative-character operation scoping and PostgreSQL transaction/idempotency rules.

## What this proves

Within an isolated single-primary/single-physical-standby PostgreSQL 17 topology, after WAL catch-up and an explicit primary stop followed by standby promotion:

- durable account/character/vehicle identity survives the authority transition;
- MQ001 durable progression survives the authority transition;
- replaying the same character-scoped operation cannot issue MQ001 rewards twice;
- the promoted PostgreSQL authority accepts the next sequential MQ002 mutation through the normal Go store path.

## Non-claims / gates that remain open

This is CI-scale single-standby promotion evidence. It does **not** prove production HA or DR readiness and does not close the broad failover/production gates. In particular it does not yet prove:

- automatic failover orchestration, leader election or fencing;
- split-brain prevention under network partition;
- synchronous-replication durability or zero-RPO guarantees;
- behavior for transactions interrupted exactly during primary loss;
- Redis failover/recovery;
- host/node loss plus workload rescheduling;
- production-volume performance during promotion;
- accepted production RPO/RTO;
- off-host backup custody or regional disaster recovery;
- live packaged Unreal reconnect across the database transition;
- production deployment readiness.

Those items remain evidence-gated and must stay unchecked until exercised against the selected deployment topology.
