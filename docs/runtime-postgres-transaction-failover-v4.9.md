# PostgreSQL Transaction-Boundary Failover Evidence v4.9

## Purpose

This evidence closes a narrower reliability question left open by v4.7: what happens when the PostgreSQL primary is lost while an authoritative quest/reward transaction is still in progress, after the quest-completion row has been inserted but before the transaction can commit.

The production trust boundary is unchanged: Unreal Engine 5.8 remains gameplay authority, Go 1.27 remains the authenticated service plane, PostgreSQL remains the durable source of truth, Redis remains ephemeral/distributed coordination state, and clients remain untrusted for durable progression or rewards.

## Executable evidence

`services/game-api/internal/store/postgres_transaction_failover_integration_test.go` runs under the existing `integration` build tag and uses PostgreSQL 17 physical streaming replication through the real Go store path. It:

1. starts an isolated PostgreSQL primary plus physical hot standby;
2. creates the canonical schema and bootstraps an authoritative account/character/starter vehicle;
3. creates a valid session and character-scoped MQ001 operation ID;
4. installs a test-only `AFTER INSERT` trigger on `quest_completions` that pauses the backend after the row insert while the transaction is still uncommitted;
5. waits until `pg_stat_activity` proves the authoritative MQ001 mutation is blocked inside that trigger;
6. terminates the primary immediately, forcing the in-flight transaction to fail before commit;
7. promotes the physical standby;
8. proves the promoted authority exposes neither the uncommitted MQ001 completion nor any reward/progression mutation;
9. proves there is no committed `quest_completions` row for the interrupted operation ID;
10. retries the same authoritative MQ001 operation against the promoted primary and requires exactly one canonical reward application;
11. replays that operation again and requires `Applied=false` with zero additional rewards.

This is materially stronger than post-commit promotion evidence because the primary is lost at a deterministic transaction boundary after an authoritative completion row has been written into the transaction but before commit.

## Security and authority invariants

The test does not introduce a production backdoor or client-controlled failover behavior. The pause trigger exists only inside the isolated integration-test database and is removed from the promoted test primary before retry. Player-facing clients still cannot choose reward values, durable character identity, database authority, or failover state.

The operation ID remains scoped to authoritative character identity and PostgreSQL transaction atomicity remains the durability boundary. Uncommitted state is not accepted as progression after promotion.

## Evidence boundary / non-claims

v4.9 demonstrates **rollback-safe retry across explicit physical-standby promotion when the old primary is terminated during an in-flight MQ001 transaction**.

It does **not** prove:

- automatic leader election or failover orchestration;
- fencing or split-brain prevention if the old primary remains writable/reappears;
- synchronous replication or zero-RPO durability for commits not replayed to standby;
- every possible transaction interruption point or every mutation type;
- Redis failover;
- Kubernetes/node rescheduling;
- production-volume latency/throughput during failover;
- production RPO/RTO targets;
- off-host/regional disaster recovery;
- live packaged Unreal Engine reconnect/recovery.

Therefore broad HA/DR, duplicate-reward-under-production-failover, node recovery, and production-readiness gates remain open until deployment-representative evidence exists.
