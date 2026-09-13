# Runtime Distributed Abuse Controls v0.9

## Scope

This slice extends the v0.8 HTTP abuse limiter without changing the selected Unreal Engine 5.8 + Go 1.27 + PostgreSQL + Redis architecture or any gameplay trust boundary.

Redis is used only to coordinate ephemeral rate-limit state across Go API processes. PostgreSQL remains the durable authority for accounts, characters, inventory, builds, quests and race state. Dedicated-server-only race mutations still require the server shared-key boundary and the existing server-authoritative ordering/idempotency rules.

## Implemented

- atomic Redis token-bucket decisions via one Lua `EVAL` per limited request,
- shared bucket keys derived from the existing credential-safe hashed identities,
- per-bucket expiry so abandoned distributed limiter state is reclaimed,
- bounded local token-bucket fallback if Redis is unavailable,
- short Redis dial/I/O deadlines so abuse-control coordination cannot indefinitely stall API requests,
- structured security log events for Redis fallback and HTTP 429 rejection without raw credentials,
- Docker Compose wiring from the Go API to Redis,
- CI Redis service and integration coverage proving two independent limiter instances consume one shared budget,
- existing HTTP 429 + `Retry-After` behavior retained.

## Failure policy

Redis coordination is defense in depth, not authorization. If Redis cannot be reached, the service does **not** fail open without limits: each API process continues enforcing the bounded v0.8 in-process limiter and emits `security_event=rate_limit_backend_fallback`.

This fallback intentionally trades perfect cross-replica coordination for availability while retaining a local abuse ceiling. Durable/gameplay authority is unaffected because Redis never authorizes quest rewards, builds, race progress, race results, or player identity.

## Security telemetry

The middleware emits:

- `security_event=rate_limit_rejected` with route scope, hashed bucket key, selected backend and retry interval,
- `security_event=rate_limit_backend_fallback` with route scope and hashed bucket key.

Raw bearer/session credentials are never written into limiter keys or these log events.

## Automated evidence

The Runtime Go workflow runs:

- ordinary Go unit tests,
- PostgreSQL + Redis integration tests,
- `go vet ./...`,
- service compilation.

The Redis integration test creates two separate limiter instances against the same Redis service. After the first instance consumes a two-token budget, the second instance must observe the shared bucket as exhausted; it then observes refill after the configured interval.

## Explicit non-claims

This slice does **not** establish production readiness. The following remain open evidence gates:

- trusted reverse-proxy/client-IP identity policy,
- multi-process HTTP load/soak evidence against multiple API replicas,
- deployed metrics/SLO dashboards and alert validation,
- authentication rejection metrics/correlation,
- ranked race impossible-state/physics anti-cheat,
- live packaged Unreal ↔ Go integration and race evidence,
- backup/restore and HA/DR drills,
- production deployment evidence.
