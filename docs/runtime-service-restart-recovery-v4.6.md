# Runtime Service Restart Recovery Evidence v4.6

## Scope

This increment adds executable PostgreSQL-backed recovery evidence for a Go service-process restart. It does not claim node failover, PostgreSQL HA, Redis HA, production RPO/RTO, or live Unreal recovery.

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

The test deliberately derives durable quest operation IDs through `core.ScopeQuestOperationID` so restart evidence exercises the same character-scoped idempotency boundary as the player-facing service path.

## Security and authority boundary

The recovery path never reconstructs durable gameplay state from client assertions or Redis. Recovery is driven by the hashed resume/session identities and PostgreSQL state. The renewed session can read and mutate only the account resolved by PostgreSQL-backed authentication.

No raw bearer/session secret is persisted by this test as evidence, and no credential is emitted to logs intentionally.

## Non-claims / remaining gates

This is **service-process restart CI evidence**, not production HA/DR evidence. The following remain open:

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
