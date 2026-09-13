# Runtime Redis Server Observability v1.6

## Scope

This increment adds bounded Redis server telemetry to the existing private Prometheus-format metrics surface. Redis remains ephemeral coordination only. PostgreSQL remains the durable service-plane authority and the Unreal dedicated gameplay server remains authoritative for gameplay acceptance.

The Go service performs a short-deadline Redis `INFO` probe only when the private metrics endpoint is scraped and exports a deliberately small allowlisted subset of numeric fields.

## Exported metrics

- `zneondrive_redis_up`
- `zneondrive_redis_connected_clients`
- `zneondrive_redis_memory_bytes{kind="used|peak"}`
- `zneondrive_redis_rejected_connections_total`
- `zneondrive_redis_evicted_keys_total`
- `zneondrive_redis_keyspace_total{outcome="hit|miss"}`
- `zneondrive_redis_instantaneous_ops_per_second`

The only labels are fixed enums controlled by source code. Redis addresses, client names, replication IDs, keys, player identifiers, session/bearer credentials, race identifiers, peer addresses and raw probe errors are never exported.

## Failure behavior

The probe uses bounded TCP dial and I/O deadlines. Probe failure emits only `zneondrive_redis_up 0`; the error string is not exposed through Prometheus output. A failed metrics probe does not change rate-limiter behavior, does not disable the existing bounded local fallback and cannot authorize gameplay or durable mutations.

## Evidence

Unit coverage verifies:

- allowlisted Redis INFO parsing;
- RESP bulk-string parsing through a real loopback TCP exchange;
- bounded successful metric output;
- Redis-down output without leaking an address/error;
- omission of dynamic Redis INFO fields such as replication identifiers.

Runtime Go CI remains responsible for Go tests, PostgreSQL/Redis integration coverage, `go vet`, module-lock cleanliness and service build.

## Explicit non-claims

This closes only the repository's Redis **server metric source/unit baseline**. It does not prove deployed scraping, dashboard or alert correctness, Redis HA/replication/failover, production sizing, long-duration soak, deployment-scale load, SLO attainment, Unreal runtime metrics, PostgreSQL server/query exporter coverage, production HA/DR or production readiness.
