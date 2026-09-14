# Runtime Service Restart Recovery Evidence v4.6

## Scope

This increment adds executable PostgreSQL-backed recovery evidence for a Go service-process restart and dual-service idempotency evidence for concurrent quest completion against the same durable PostgreSQL authority. It does not claim node failover, PostgreSQL HA, Redis HA, production RPO/RTO, or live Unreal recovery.

The selected production architecture remains unchanged:

- Unreal Engine 5.8 dedicated server is gameplay authority.
- Go 1.27 is the authenticated service plane.
- PostgreSQL is the durable source of truth for account, character, vehicle, progression, quest, inventory, blueprint, build and race state.
- Redis remains ephemeral/distributed coordination state and is not treated as the durable recovery source.

## Executable evidence

`services/game-api/internal/store/reconnect_recovery_integration_test.go` runs under the existing Runtime Go integration job with PostgreSQL 17. The test:

1. creates a durable account/character/starter vehicle through the canonical bootstrap path;
2. creates a hashed session and applies `MQ001` through a character-scoped durable operation key;
3. closes the PostgreSQL-backed store, removing process-local service state;
4. reopens the store against the same PostgreSQL database and replays canonical schema setup;
5. resumes with the same hashed resume credential and requires the same account, character and vehicle identities plus the previously committed progression;
6. creates a renewed session after restart and resolves the same authoritative snapshot;
7. replays the original `MQ001` durable operation and requires `Applied=false` with no duplicate reward;
8. completes `MQ002` after restart to prove the recovered service can continue the sequential authoritative progression path.

`services/game-api/internal/store/quest_concurrency_integration_test.go` adds a second PostgreSQL-backed evidence path representing two independent Go service instances sharing the same durable database authority. The test:

1. opens two independent PostgreSQL stores against the same PostgreSQL 17 database;
2. creates two separate authenticated sessions for the same authoritative character;
3. derives one character-scoped durable `MQ001` operation ID through `core.ScopeQuestOperationID`;
4. releases two `CompleteQuest` calls concurrently from separate store instances;
5. requires both calls to remain bound to the same authoritative character;
6. requires exactly one receipt with `Applied=true` and canonical `MQ001` reward values;
7. requires exactly one replay receipt with `Applied=false` and zero duplicate reward;
8. reloads the authoritative snapshot and requires exactly one durable `MQ001` completion and exactly one reward application.

These tests deliberately derive durable quest operation IDs through `core.ScopeQuestOperationID` so restart/concurrency evidence exercises the same character-scoped idempotency boundary as the player-facing service path.

## Security and authority boundary

The recovery path never reconstructs durable gameplay state from client assertions or Redis. Recovery is driven by the hashed resume/session identities and PostgreSQL state. Renewed/concurrent sessions can read and mutate only the account resolved by PostgreSQL-backed authentication.

The dual-service test strengthens evidence for duplicate-reward resistance under concurrent application instances sharing one PostgreSQL authority, but it is **not** a database failover test and must not be used as proof of duplicate-reward safety across PostgreSQL primary promotion or split-brain conditions.

No raw bearer/session secret is persisted by these tests as evidence, and no credential is emitted to logs intentionally.

## Non-claims / remaining gates

This is **service-process restart plus shared-database concurrent-service CI evidence**, not production HA/DR evidence. The following remain open:

- real Unreal Engine 5.8 Client and dedicated Server build/package evidence;
- live packaged Unreal ↔ Go reconnect and race lifecycle evidence;
- host/node loss and service rescheduling recovery;
- PostgreSQL primary failure/failover and replica promotion;
- Redis failover/recovery under the selected deployment topology;
- duplicate-reward prevention during actual database failover;
- production backup/off-host retention and restore drill;
- accepted/measured production RPO/RTO;
- deployment-scale load and long-duration soak;
- deployed observability/SLO and regional DR evidence.

Accordingly, this increment must not be used to mark the repository production-ready or to close the broad HA/DR gates in `IMPLEMENTATION-CHECKLIST.md`.
