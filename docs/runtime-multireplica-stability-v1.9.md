# Runtime Multi-Replica Stability Evidence v1.9

## Scope

This increment strengthens the existing CI-scale distributed-rate-limit evidence without changing the selected Unreal Engine 5.8 + Go 1.27 + PostgreSQL + Redis architecture or any authority boundary.

Runtime Go integration now runs a repeated two-replica HTTP/Redis stability scenario in addition to the existing single concurrent burst test.

## Scenario

The integration test:

- starts two independent HTTP server replicas using the real distributed rate-limit middleware;
- uses the Redis 8 service supplied by Runtime Go CI as the shared coordination backend;
- configures loopback as the trusted test ingress peer;
- executes 12 fresh-client rounds;
- sends 64 concurrent bootstrap requests per round, balanced evenly across both replicas;
- configures a 16-request shared burst with refill disabled during the evidence window;
- requires exactly 16 successful downstream requests and exactly 48 HTTP 429 responses in every round;
- requires every HTTP 429 response to include `Retry-After`;
- fails on transport errors, unexpected statuses, replica imbalance, per-round budget drift, or aggregate-count mismatch.

The complete scenario therefore exercises 768 HTTP requests and requires exactly 192 allows plus 576 shared-budget rejections across the repeated cycles.

## Trust boundary

Redis remains ephemeral abuse-control coordination only. PostgreSQL remains authoritative for durable account/progression/vehicle/race state, and the Unreal dedicated gameplay server remains authoritative for gameplay acceptance.

The repeated-load test cannot grant authentication, mutate durable progression, assert vehicle/race state, or weaken the existing server-only gameplay-ticket redemption boundary.

## Evidence value

The scenario is intended to detect failures that a single burst can miss, including:

- cross-replica budget drift between repeated fresh identities;
- intermittent Redis coordination failures;
- inconsistent 429/`Retry-After` behavior;
- HTTP transport instability under repeated concurrent bursts;
- accidental per-process burst multiplication.

## Explicit non-claims

This is **CI-scale repeated stability evidence**, not deployment-scale load testing and not a long-duration soak test.

It does not establish:

- production throughput or latency percentiles;
- production saturation or capacity limits;
- multi-node Redis HA/failover behavior;
- deployed ingress correctness;
- hours/days of soak stability;
- packaged Unreal↔Go runtime behavior;
- production SLO attainment;
- HA/DR or production readiness.

Deployment-scale load and long-duration soak gates therefore remain open.
