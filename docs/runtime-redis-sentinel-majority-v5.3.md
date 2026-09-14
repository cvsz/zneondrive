# Runtime Redis Sentinel Majority Resolution Evidence v5.3

## Scope

This increment hardens the Go distributed rate-limit coordination client used with Redis Sentinel. Redis remains ephemeral coordination only; PostgreSQL remains the durable authority for account, character, inventory, blueprint, quest, build, race, and reward state.

## Implemented behavior

When `REDIS_SENTINEL_ADDRS` is configured, the limiter no longer trusts the first Sentinel endpoint that returns a master address. Every limiter mutation resolves Sentinel state and requires a strict majority of the configured Sentinel endpoints to report the same Redis master address before a Redis mutation is attempted.

For three configured Sentinels, at least two must agree. A single stale or isolated Sentinel therefore cannot redirect the application limiter to a former primary. If no address reaches majority, Redis coordination returns an error and the existing bounded process-local limiter fallback remains in force rather than failing open.

Direct `REDIS_ADDR` mode is unchanged for non-Sentinel deployments.

## Automated evidence

`redis_sentinel_quorum_test.go` covers:

- a stale first Sentinel while two peers agree on the promoted authority;
- three mutually inconsistent Sentinel answers, which must fail resolution;
- one unavailable Sentinel while the remaining majority agrees.

Existing Redis Sentinel integration tests continue to exercise real Redis 8 primary/replica + three-Sentinel failover and old-primary rejoin paths. Because those tests use the production resolver, they now also require majority agreement before limiter traffic follows a promoted Redis authority.

## Security / trust boundary

This is defense in depth against stale/minority Sentinel observations. It does not make Redis authoritative for gameplay or durable state and does not permit clients to assert reward, progression, inventory, blueprint, Roadworthy, build, or race state.

## Non-claims

This evidence does **not** prove full split-brain prevention. In particular it does not provide an external fencing or STONITH mechanism for a former Redis primary that remains independently reachable and writable during an asymmetric network partition. It also does not prove Redis Cluster behavior, zero-loss asynchronous replication, deployment-scale failover/load, long-duration soak, Kubernetes node rescheduling, accepted production RPO/RTO, regional DR, or live packaged Unreal recovery.

The repository must therefore remain **Evidence Gated / not Production Ready**, and the broad Redis asymmetric-partition split-brain/Cluster/deployment-HA gate remains open.
