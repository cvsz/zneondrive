# Runtime Reference Parity Harness v2.4

**Status:** implemented test-harness evidence only. Full Python-reference parity in Go remains open.

## Purpose

This increment adds one versioned fixture consumed by both the dependency-free Python reference oracle and the Go 1.27 runtime tests. The goal is to detect semantic drift in deterministic contracts without changing the selected production architecture or authority boundaries.

The shared vectors currently cover:

- canonical vehicle build hashing for the starter parts set;
- order/duplicate normalization producing the same build validation hash;
- the MQ001–MQ100 quest-ID domain and predecessor mapping on the Go runtime side;
- Python-side structural verification that the shared quest vectors themselves stay internally consistent.

## Evidence boundary

Passing these tests proves only that the covered deterministic vectors remain aligned. It does **not** prove full oracle parity. In particular, the following remain outside this harness and must stay open until equivalent executable coverage exists:

- account/character uniqueness semantics;
- garage-slot/VIP capacity semantics;
- starter-lineage deletion protection;
- complete quest reward/idempotency semantics across both implementations;
- build revision conflict behavior across both implementations;
- complete race registration/checkpoint/result parity;
- inventory/blueprint/rebuild parity;
- live Unreal transport or gameplay behavior.

## Trust boundaries

This test-only increment does not move authority. Unreal Engine 5.8 dedicated servers remain authoritative for moment-to-moment gameplay, Go remains the service/auth contract plane, PostgreSQL remains durable authority, and Redis remains ephemeral coordination.

## Validation

The repository CI must run both Python reference tests and Go unit tests. The shared fixture is intentionally free of credentials, player data, endpoints, environment-specific identifiers, or production secrets.
