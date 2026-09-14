# Runtime Redis Promotion Recovery Evidence v5.0

## Scope

This increment adds executable CI evidence for the Redis-backed distributed HTTP rate limiter across an **explicit Redis 8 primary/replica promotion**.

Redis remains ephemeral coordination only. It is not a source of durable account, progression, inventory, blueprint, build, quest, race, or gameplay authority. PostgreSQL remains the durable service-plane authority and Unreal dedicated servers remain gameplay authority.

## Executable evidence

`services/game-api/internal/httpapi/redis_promotion_integration_test.go` creates an isolated Docker network with one Redis 8 primary and one physical replica, then:

1. waits for the replica link to report connected;
2. consumes the complete canonical limiter burst against the primary;
3. verifies the primary rejects the next request;
4. waits until the replica has the exact exhausted token state;
5. stops the old primary before promotion;
6. promotes the replica with `REPLICAOF NO ONE`;
7. verifies the promoted Redis still rejects the same exhausted limiter key;
8. advances the deterministic limiter clock by one refill interval and verifies canonical refill semantics continue.

The test is part of the existing `go test -tags=integration ./...` Runtime Go gate and fails in CI if Docker is unavailable instead of silently skipping the recovery evidence.

## Security and authority boundary

The evidence strengthens availability/abuse-control confidence without changing trust boundaries:

- Redis contains only expiring limiter coordination state.
- Redis promotion cannot authorize gameplay mutations or replace PostgreSQL state.
- The limiter Lua decision remains atomic on the active Redis authority.
- The existing Redis-unavailable path still falls back to bounded process-local limiting rather than unbounded fail-open behavior.
- No credentials, player IDs, race IDs, or durable gameplay data are introduced into the Redis recovery fixture.

## Non-claims / remaining blockers

This CI-scale test does **not** prove:

- automatic Redis failover or Sentinel/Cluster orchestration;
- split-brain prevention or fencing if the old Redis primary remains writable/reappears;
- zero-loss replication for writes not received by the replica;
- production Redis persistence/custody guarantees (Redis is intentionally non-durable authority here);
- Kubernetes/node rescheduling;
- deployment-scale limiter load or long-duration soak;
- accepted production RPO/RTO;
- PostgreSQL automatic failover/fencing or regional DR;
- live packaged Unreal Engine reconnect/recovery;
- production readiness.

Accordingly the broad Redis HA, deployment HA/DR, live Unreal, load/soak, and production-readiness gates remain open.
