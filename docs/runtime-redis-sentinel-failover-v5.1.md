# Runtime Redis Sentinel Failover Evidence v5.1

## Scope

This increment adds Redis Sentinel-aware master discovery to the Go 1.27 distributed HTTP abuse-control limiter while preserving the selected Unreal Engine 5.8 + Go 1.27 + PostgreSQL + Redis authority model.

Redis remains **ephemeral coordination only**. PostgreSQL remains the durable account/progression/inventory/build/quest/race authority, and Unreal dedicated servers remain gameplay authority. Sentinel cannot grant rewards, mutate inventory/blueprints, mark vehicles Roadworthy, or author race results.

## Runtime contract

The game API may be configured with:

- `REDIS_SENTINEL_ADDRS` — comma-separated Sentinel TCP endpoints;
- `REDIS_SENTINEL_MASTER` — monitored Redis master name.

Both settings are required together. Partial Sentinel configuration fails startup rather than silently pretending HA is enabled. When both are configured, Sentinel discovery takes precedence over the legacy direct `REDIS_ADDR` limiter path.

Before a distributed limiter decision, the Go limiter queries the configured Sentinel endpoints with `SENTINEL get-master-addr-by-name`, selects the first valid response, and sends the existing atomic Lua token-bucket decision to that Redis primary. If Sentinel or Redis is temporarily unavailable, the existing bounded in-process limiter remains the fail-safe fallback; this does not move durable/gameplay authority into the application process.

## CI evidence

`services/game-api/internal/httpapi/redis_sentinel_integration_test.go` creates an isolated Redis 8 topology containing:

- one Redis primary;
- one Redis replica;
- three Redis Sentinel processes with quorum 2.

The executable test:

1. waits for the primary/replica replication link and Sentinel discovery;
2. exhausts a canonical distributed limiter bucket through the Sentinel-backed limiter;
3. waits until the exhausted bucket state is present on the replica;
4. stops the old primary without changing the application's Sentinel endpoints;
5. waits for Sentinel quorum to promote the replica and report the new primary;
6. reuses the same Sentinel-backed limiter instance and verifies the exhausted bucket remains rejected;
7. advances the deterministic limiter clock one refill interval and verifies canonical refill semantics continue on the promoted Redis.

This closes the repository gap between **explicit manual Redis promotion evidence (v5.0)** and **application-side automatic master rediscovery through Sentinel**.

## Security and trust boundary

- Sentinel addresses and master name are operator configuration, never player input.
- Sentinel discovery influences only where ephemeral rate-limit coordination is performed.
- A Redis/Sentinel outage falls back to the existing bounded local limiter; it never causes unbounded fail-open behavior.
- No raw bearer/session credential is written to Redis; limiter identities remain hashed before bucket storage.
- No durable gameplay state is stored in Redis.

## Explicit non-claims

This CI-scale evidence does **not** prove:

- fencing or split-brain prevention if an old Redis primary remains writable or rejoins incorrectly;
- zero-loss asynchronous replication for limiter writes not received by the replica before failure;
- Redis Cluster behavior;
- deployment-scale failover under sustained production traffic;
- long-duration soak stability;
- Kubernetes node loss/rescheduling;
- production Redis authentication/TLS/ACL/secret-rotation posture;
- accepted production RPO/RTO or regional DR;
- live packaged Unreal reconnect/race behavior;
- production readiness.

The repository must remain **Evidence Gated / Not Production Ready** until the wider real-Unreal, deployment, HA/DR, load/soak, observability, and security gates are proven.
