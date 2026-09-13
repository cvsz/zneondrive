# Runtime Multi-Replica HTTP Load Evidence v1.1

**Status:** CI-scale integration evidence only. This is not production-scale load/soak certification.

## Scope

This evidence slice verifies that the Go service-plane distributed limiter behaves correctly when concurrent HTTP traffic is split across two independent HTTP server instances backed by the same Redis service.

It preserves the selected architecture and trust boundaries:
- Unreal Engine 5.8 dedicated servers remain gameplay-authoritative.
- Go 1.27 remains the service plane.
- PostgreSQL remains durable authority.
- Redis remains ephemeral coordination for abuse-control state.
- trusted proxy identity is resolved only through the explicit `TrustedProxyPolicy` / `TRUSTED_PROXY_CIDRS` boundary.

## Executable evidence

`services/game-api/internal/httpapi/redis_ratelimit_integration_test.go` now includes `TestDistributedRateLimiterConcurrentHTTPAcrossReplicas`.

The test:
1. requires the real Redis service supplied by Runtime Go CI,
2. creates two independent `httptest.Server` HTTP replicas,
3. configures loopback as the trusted test ingress peer,
4. sends 128 concurrent bootstrap requests split evenly across both replicas,
5. presents one effective client identity through `X-Forwarded-For`,
6. configures a 32-request shared burst with no meaningful refill during the run,
7. requires exactly 32 requests to reach the downstream handler across both replicas,
8. requires the remaining 96 requests to receive HTTP 429,
9. rejects transport errors or unexpected status codes.

This verifies the limiter at the HTTP middleware boundary rather than only calling the Redis token-bucket primitive directly.

## What this closes

This provides repository/CI evidence that:
- two independent HTTP replicas consume one Redis-coordinated limiter budget,
- trusted-ingress-derived identity remains consistent across replicas,
- concurrent requests do not create a per-process burst multiplier while Redis is available,
- excess requests are rejected at the HTTP boundary with the existing limiter contract.

## What remains open

This does **not** prove:
- deployed ingress header sanitization or direct-bypass prevention,
- target-region latency budgets,
- alpha concurrency capacity,
- long-duration soak stability,
- Redis failover behavior under sustained production traffic,
- CPU/RAM/network saturation points,
- production observability/SLO attainment,
- HA/DR readiness,
- Unreal packaged-client or race load behavior.

Those require deployed benchmark evidence with recorded commit/build, infrastructure, scenario, duration, latency percentiles, resource metrics, and saturation/failure points as defined in `docs/performance-budget.md`.

## Merge gate

This evidence is valid only when the focused PR passes the repository CI, Runtime Go unit/integration/vet/build workflow, CodeQL, and Dependency Review. Production-readiness status remains evidence-gated.
