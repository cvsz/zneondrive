# Runtime Redis Limiter Observability v1.5a

## Scope

This increment adds bounded, process-local Prometheus telemetry for the Redis-backed distributed rate limiter used by the Go service plane. It does not change gameplay authority, durable state ownership, or the existing fail-safe limiter behavior.

The private metrics listener now exports:

- `zneondrive_redis_rate_limit_decisions_total{outcome="allowed"}`
- `zneondrive_redis_rate_limit_decisions_total{outcome="rejected"}`
- `zneondrive_redis_rate_limit_decisions_total{outcome="error"}`

The outcome label is a fixed three-value enum. No limiter bucket, IP address, bearer/session token, race identifier, player identifier, Redis address, or request-controlled value is exported.

## Trust boundary

Redis remains ephemeral coordination only. PostgreSQL remains the durable service-plane authority and the Unreal dedicated gameplay server remains authoritative for gameplay acceptance. If Redis is unavailable, the existing bounded local token bucket remains active; Redis failures are counted as `outcome="error"` and are not misreported as Redis allows/rejections.

These metrics are observational only. They cannot authorize requests, mutate token buckets, change fallback behavior, or weaken dedicated-server/PostgreSQL checks.

## Evidence

Unit coverage verifies:

- unavailable Redis attempts increment only the Redis error outcome counter while local fallback continues enforcing the configured burst;
- Redis outcome metrics use only the bounded outcome label;
- metric output does not contain bucket identities, authorization values, bearer tokens, race identifiers, or peer addresses.

Existing Redis integration coverage continues to prove shared budget behavior across independent limiter instances.

## Explicit non-claims

This increment is **not** Redis server observability. It does not expose Redis CPU, memory, command latency, replication, persistence, eviction, connection, or availability metrics from Redis itself. The repository production-readiness gate for complete Redis server metrics therefore remains open.

It also does not close deployed scrape/dashboard/alert verification, SLO measurement, deployment-scale load, long-duration soak, live Unreal telemetry, HA/DR, or production deployment evidence.
