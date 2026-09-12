# Runtime Security Hardening v0.8

Status: **implemented source-level abuse control; production/distributed evidence still gated**.

This slice hardens the Go service-plane HTTP boundary without changing the selected Unreal Engine 5.8 + Go 1.27 + PostgreSQL + Redis architecture or moving authority to clients.

## Implemented

The production server now wraps the API with a bounded in-process token-bucket middleware. The limiter covers:

- session bootstrap,
- authenticated state reads,
- gameplay-ticket issuance,
- quest and vehicle-build mutations,
- internal gameplay-ticket redemption,
- internal race start/checkpoint/finish mutations.

`/healthz` is intentionally excluded so orchestration probes are not starved by player traffic.

Authenticated player buckets are keyed from a one-way hash of the bearer token. Internal endpoints are keyed by remote network identity and continue to require the server-only `X-Game-Server-Key`; the limiter does not replace authentication. Raw bearer credentials are never embedded in rate-limit keys.

The limiter has bounded bucket cardinality and evicts stale/oldest entries to prevent attacker-controlled identity cardinality from growing memory without bound. Rejected requests return HTTP `429` with `Retry-After` and `{ "error": "rate_limited" }`.

## Trust boundary

Rate limiting is an abuse-resistance control only. It does not authorize gameplay state and it does not make client data authoritative. PostgreSQL remains the source of durable account, character, inventory, build, quest and race state. Dedicated-server race mutations remain behind the shared-key boundary and existing authoritative ordering/idempotency checks.

## Evidence

Unit tests cover:

- burst exhaustion and deterministic refill,
- bounded bucket cardinality,
- credential-safe hashed bucket keys,
- HTTP 429 + `Retry-After`,
- per-bearer bucket isolation,
- health-probe exemption.

Repository CI, Go unit/integration/vet/build, CodeQL and Dependency Review remain required before merge.

## Explicit non-claims / remaining production gates

This implementation is intentionally **per process**. It is useful immediately against local abuse and request storms, but it is not sufficient for a horizontally scaled production fleet because counters are not coordinated across API replicas.

Still required before production-readiness claims:

1. Redis-backed distributed rate-limit counters with fail-open/fail-closed policy defined per endpoint.
2. Trusted-proxy configuration and canonical client identity extraction at the ingress boundary; this middleware deliberately does not trust arbitrary `X-Forwarded-For` headers.
3. Runtime telemetry for rate-limit rejects and auth failures.
4. Ranked-race impossible-state / physics-derived anti-cheat telemetry.
5. Abuse/load testing proving limits under concurrent clients and multi-replica deployment.
6. Live Unreal↔Go transport evidence, production deployment, restore/DR and HA evidence.

The per-process limiter must remain as a defense-in-depth fallback even after Redis coordination is added.
