# Runtime Quest Store Replay Binding v4.8

## Scope

This increment hardens the PostgreSQL quest idempotency boundary below the public HTTP operation-key scoping layer.

`CompleteQuest` already authenticates the session and resolves the authoritative `character_id` from PostgreSQL. A durable `operation_id` replay is now accepted only when the persisted quest-completion row matches both that authoritative character and the requested canonical quest ID.

## Executable evidence

`services/game-api/internal/store/quest_operation_binding_integration_test.go` creates two independent accounts/characters in PostgreSQL and deliberately reuses the same store-level operation ID for `MQ001`.

The test proves that:

1. the first character can apply canonical MQ001 exactly once;
2. the same operation ID presented by a different authenticated character is rejected with `ErrOperationKey` even when the requested quest ID is identical;
3. the rejected attempt does not mutate the second character's money, XP, reputation, or quest history; and
4. the original character can still replay the exact operation idempotently without receiving rewards twice.

## Trust boundary

- Unreal Engine 5.8 dedicated servers remain gameplay authority.
- Go 1.27 remains the authenticated service plane.
- PostgreSQL remains durable identity, progression, reward, and idempotency authority.
- Redis remains ephemeral coordination/rate-limit state only.
- Player-supplied operation IDs never establish character ownership.

This is defense-in-depth for internal/direct store callers. The public quest route continues to derive durable operation IDs from the PostgreSQL-resolved authoritative character identity via `core.ScopeQuestOperationID`.

## Non-claims

This evidence does not prove live packaged Unreal↔Go integration, Garage 17/MQ001–MQ012 playability, automatic PostgreSQL failover/fencing, interrupted-transaction safety, Redis failover, deployment-scale load/soak, production RPO/RTO, regional DR, final-physics anti-cheat calibration, or production readiness.
