# zNeonDrive — PROJECT: NEON DRIVE

<p align="center">
  <img src="./assets/zneondrive-banner.png" alt="PROJECT: NEON DRIVE — NOVA CITY 2097" width="100%" />
</p>

<p align="center">

[![CI](https://github.com/cvsz/zneondrive/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/cvsz/zneondrive/actions/workflows/ci.yml)
[![CodeQL](https://github.com/cvsz/zneondrive/actions/workflows/codeql.yml/badge.svg?branch=main)](https://github.com/cvsz/zneondrive/actions/workflows/codeql.yml)
[![Runtime Go](https://github.com/cvsz/zneondrive/actions/workflows/runtime-go.yml/badge.svg?branch=main)](https://github.com/cvsz/zneondrive/actions/workflows/runtime-go.yml)
[![Dependency Review](https://github.com/cvsz/zneondrive/actions/workflows/dependency-review.yml/badge.svg)](https://github.com/cvsz/zneondrive/actions/workflows/dependency-review.yml)

![Phase](https://img.shields.io/badge/Phase-4.19%20Unreal%20Rotation%20Envelope-00D8FF?style=flat-square)
![Runtime](https://img.shields.io/badge/Runtime-v2.3-8A2BE2?style=flat-square)
![Production Readiness](https://img.shields.io/badge/Production%20Readiness-Evidence%20Gated-F59E0B?style=flat-square)
![UE Source Build](https://img.shields.io/badge/UE%20Source%20Build-Evidence%20Pending-EF4444?style=flat-square)

![Unreal Engine](https://img.shields.io/badge/Unreal%20Engine-5.8-0E1128?style=flat-square&logo=unrealengine&logoColor=white)
![Go](https://img.shields.io/badge/Go-1.27-00ADD8?style=flat-square&logo=go&logoColor=white)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-Durable%20State-4169E1?style=flat-square&logo=postgresql&logoColor=white)
![Redis](https://img.shields.io/badge/Redis-Ephemeral%20Coordination-DC382D?style=flat-square&logo=redis&logoColor=white)
![Python](https://img.shields.io/badge/Python-Reference%20Oracle-3776AB?style=flat-square&logo=python&logoColor=white)

![Main Quests](https://img.shields.io/badge/Main%20Quests-100-FF4FD8?style=flat-square)
![Districts](https://img.shields.io/badge/Districts-10-00D8FF?style=flat-square)
![Characters](https://img.shields.io/badge/Characters-25-A7FF4F?style=flat-square)
![Factions](https://img.shields.io/badge/Factions-5-FFC857?style=flat-square)

[![License](https://img.shields.io/github/license/cvsz/zneondrive?style=flat-square)](./LICENSE)
[![Stars](https://img.shields.io/github/stars/cvsz/zneondrive?style=flat-square)](https://github.com/cvsz/zneondrive/stargazers)
[![Forks](https://img.shields.io/github/forks/cvsz/zneondrive?style=flat-square)](https://github.com/cvsz/zneondrive/forks)
[![Open Issues](https://img.shields.io/github/issues/cvsz/zneondrive?style=flat-square)](https://github.com/cvsz/zneondrive/issues)
[![Last Commit](https://img.shields.io/github/last-commit/cvsz/zneondrive/main?style=flat-square)](https://github.com/cvsz/zneondrive/commits/main)

</p>

> **You don't own the road. You earn it.**

zNeonDrive is the design and implementation repository for **PROJECT: NEON DRIVE**, an 18+ persistent online open-world RPG centered on vehicle building, racing, adventure, factions, relationships, crews, pets, and a long-lived social world.

## Product direction

- **World:** NOVA CITY, 2097
- **Genre:** Online Open World / RPG / Vehicle Building / Racing / Adventure / Social
- **World model:** persistent online world, server-authoritative by design
- **Player identity:** 1 account → 1 primary character → 1 starter vehicle
- **Vehicle philosophy:** a vehicle is a persistent identity object with ownership, builder, build, repair, race, and reputation history
- **VIP constraint:** garage/storage/convenience capacity only; no direct competitive performance advantage
- **Current phase:** **Phase 4.19 Unreal Rotation Envelope v2.3**
- **Selected implementation direction:** Unreal Engine 5.8 client/dedicated gameplay server + Go 1.27 service plane + PostgreSQL + Redis

## Runtime prototype v0.4

The repository now contains executable implementation scaffolding for both planes.

### Gameplay plane — Unreal

`game/` contains:
- UE 5.8 C++ project,
- Game / Editor / dedicated Server targets,
- default GameMode,
- prototype vehicle pawn,
- client input → server RPC → authoritative movement,
- replicated movement back to clients,
- bounded aggregate dedicated-server telemetry source for sessions, tick timing, input clamps, ticket redemption, authority movement, identity-gate blocks, collision blocks, explicit net-update requests, impossible-displacement envelope violations, and impossible-rotation envelope violations,
- source-build workflow baseline for a self-hosted UE 5.8 Linux runner.

See [game/README.md](./game/README.md).

### Durable service plane — Go

`services/game-api/` contains:
- Go 1.27 modular service binary,
- PostgreSQL account/session/character/vehicle/build/quest/inventory/blueprint/race persistence,
- hashed resume/session credentials,
- one-character + starter-vehicle bootstrap,
- immutable vehicle build revisions,
- optimistic revision checks,
- idempotent durable mutation operation IDs,
- MQ001–MQ100 sequential quest gate,
- MQ012 Roadworthy transition,
- authoritative race instance/checkpoint/result state,
- bounded local HTTP token-bucket abuse controls,
- Redis-coordinated distributed rate-limit state with bounded local fallback,
- explicit trusted-ingress client identity policy,
- integrated PostgreSQL + Redis trust-boundary abuse-path evidence,
- repeated two-replica Redis limiter stability evidence,
- credential-safe rate-limit and race-integrity/auth security events,
- low-cardinality Prometheus-format HTTP request/latency/in-flight metrics,
- bounded PostgreSQL pgx connection-pool metrics,
- bounded Redis server INFO metrics on the private metrics surface,
- bounded current-database PostgreSQL server metrics,
- isolated PostgreSQL 17 backup/restore CI evidence,
- unit, PostgreSQL/Redis integration, and concurrent two-replica HTTP limiter tests.

Local stack:

```bash
cp .env.example .env
docker compose up --build
curl http://127.0.0.1:18080/healthz
curl http://127.0.0.1:19090/metrics
```

Host defaults intentionally avoid common 5432/6379 conflicts:
- PostgreSQL: `55432`
- Redis: `56379`
- Game API: `18080`
- Internal metrics: `127.0.0.1:19090`

See [Runtime Prototype v0.4](./docs/runtime-prototype-v0.4.md).

### Authenticated runtime integration — v0.5

The gameplay/service boundary now has an implemented source contract:
- Unreal client bootstraps/resumes durable identity through the Go API,
- session credentials stay client-side; dedicated servers receive only short-lived one-time gameplay tickets,
- Go stores gameplay tickets only as hashes and atomically consumes them once,
- the internal redemption endpoint requires a server-only shared key,
- the dedicated server obtains the authoritative PostgreSQL snapshot directly from Go,
- durable VehicleID, build revision, active parts, and Roadworthy state are bound on the authority and replicated to clients,
- networked prototype movement remains blocked until durable identity is server-bound,
- HTTP/PostgreSQL E2E proves ticket issuance, single-use redemption, MQ001 persistence, and resume-key reconnect.

See [Runtime Integration v0.5](./docs/runtime-integration-v0.5.md).

### Inventory-authoritative rebuild — v0.6

The first Garage 17 rebuild state is now implemented in the durable service plane:
- PostgreSQL inventory quantities and blueprint unlocks,
- MQ004 idempotent salvage grant,
- MQ005 starter-rebuild blueprint unlock,
- MQ009 idempotent recovered-part grant,
- exact catalog validation against the canonical vehicle-part catalog,
- atomic consume/return inventory accounting when a build revision changes,
- operation-ID replay binding to the exact build hash,
- reconnect persistence for inventory, blueprints, and the active build,
- Unreal snapshot parsing for inventory and blueprint state.

See [Runtime Inventory + Rebuild v0.6](./docs/runtime-inventory-rebuild-v0.6.md).

### Server-authoritative race lifecycle — v0.7

The first durable race authority slice is implemented in Go/PostgreSQL:
- race start is restricted to the dedicated-server shared-key boundary,
- the service derives account/character/vehicle ownership from PostgreSQL,
- the vehicle must already be Roadworthy,
- every race instance binds the exact active build revision and validation hash at start,
- checkpoint indices must match the server-side cursor and elapsed time must increase monotonically,
- finish must match the recorded checkpoint count and exceed the last checkpoint time,
- start/checkpoint/finish writes are idempotent and reject operation-key payload reuse,
- final result hashes bind race identity, account/character/vehicle, immutable build evidence, checkpoint count, and finish time.

This is service-plane evidence only; live packaged Unreal race transport, physics anti-cheat, playable race content, load/soak, HA/DR, and production deployment remain open.

See [Runtime Authoritative Race v0.7](./docs/runtime-authoritative-race-v0.7.md).

### Runtime abuse resistance — v0.8

The Go service process wraps public and internal mutation traffic in a bounded token-bucket middleware:
- bootstrap/state/game-ticket/quest/build endpoints are limited by remote identity or hashed bearer identity,
- internal ticket and race mutations are rate-limited before reaching the existing server-only authorization boundary,
- raw bearer credentials are not retained in limiter keys,
- limiter bucket cardinality is bounded with stale/oldest eviction,
- rejected requests return HTTP `429` with `Retry-After`,
- `/healthz` remains exempt for orchestration probes.

See [Runtime Security Hardening v0.8](./docs/runtime-security-hardening-v0.8.md).

### Distributed abuse controls — v0.9

The v0.8 limiter has a Redis coordination layer for multi-process API deployments:
- each limited request uses an atomic Redis Lua token-bucket decision,
- independent API limiter instances share one Redis abuse budget,
- distributed bucket keys reuse the credential-safe hashed identity scheme,
- Redis limiter keys expire automatically,
- short Redis dial/I/O deadlines bound dependency impact,
- Redis failure falls back to the bounded local token bucket rather than removing rate limiting,
- rejection and fallback events are logged with scope + hashed bucket only,
- Compose wires the Go API to Redis,
- Runtime Go CI includes a Redis service and verifies shared budget/refill behavior across independent limiter instances.

Redis remains **ephemeral abuse-control coordination only**. PostgreSQL and the dedicated gameplay server remain authoritative for durable/gameplay state.

See [Runtime Distributed Abuse Controls v0.9](./docs/runtime-distributed-abuse-controls-v0.9.md).

### Trusted ingress identity — v1.0

The service plane now has an explicit proxy trust boundary for source-IP based abuse controls:
- `RemoteAddr` remains authoritative by default and forwarding headers are ignored for direct/untrusted peers,
- `TRUSTED_PROXY_CIDRS` explicitly allowlists ingress/proxy peers that may influence client identity,
- invalid CIDR configuration fails startup instead of silently broadening trust,
- trusted `X-Forwarded-For` chains are evaluated from right to left and skip only configured trusted proxy hops,
- malformed forwarding chains fail closed to the immediate socket peer,
- forwarded header contents are not written to the malformed-chain security event,
- unit tests cover spoof resistance, multi-hop chains, malformed values, exact-IP/CIDR configuration, and distinct client buckets behind one ingress.

This is source/unit evidence only. It does not prove a deployed ingress sanitizes forwarding headers correctly or blocks direct bypass traffic.

See [Runtime Trusted Ingress Identity v1.0](./docs/runtime-trusted-ingress-v1.0.md).

### Multi-replica HTTP limiter evidence — v1.1

Runtime Go integration now exercises the distributed limiter at the real HTTP middleware boundary:
- two independent HTTP server replicas share one Redis limiter backend,
- a trusted test ingress supplies one effective client identity,
- 128 concurrent bootstrap requests are split evenly across both replicas,
- the shared burst is fixed at 32 requests with refill disabled for the test window,
- exactly 32 requests may reach the downstream handler across both replicas,
- the remaining 96 requests must return HTTP 429,
- transport failures or unexpected HTTP statuses fail the integration test.

This is CI-scale concurrency evidence only. It does **not** close deployed ingress verification, deployment-scale load, long-duration soak, measured performance budgets, ranked anti-cheat, observability/SLO, HA/DR, backup/restore, or production deployment gates.

See [Runtime Multi-Replica HTTP Load Evidence v1.1](./docs/runtime-multi-replica-load-v1.1.md).

### Race integrity rejection telemetry — v1.2

The Go service now observes existing authoritative internal-route outcomes and emits credential-safe security telemetry:
- authoritative checkpoint/finish HTTP 400/409 outcomes emit `security_event=race_integrity_rejected`,
- internal game-server HTTP 401 outcomes emit `security_event=game_server_auth_rejected`,
- race-instance and socket-peer identifiers are hashed into deterministic correlation buckets,
- raw shared keys, raw race IDs, and raw peer addresses are not written to these events,
- static route scopes prevent dynamic path identifiers from leaking into logs,
- telemetry is observational only and cannot override PostgreSQL/dedicated-server race acceptance.

This is a rejection/correlation baseline, not physics anti-cheat. It does **not** prove live Unreal telemetry transport, acceleration/teleport envelope behavior under packaged physics, sanctions, deployed alerting, or production anti-cheat readiness.

See [Runtime Race Integrity Telemetry v1.2](./docs/runtime-race-integrity-telemetry-v1.2.md).

### Service observability metrics — v1.3

The Go service now exposes a source-level measurement baseline suitable for later load/soak and SLO evidence:
- HTTP request totals by bounded route scope and status class,
- fixed-bucket request-duration histogram by bounded route scope,
- current in-flight request gauge,
- static route labels prevent dynamic quest, vehicle, race, account, operation, peer, or credential values from entering metrics,
- metrics use a separate `METRICS_LISTEN_ADDR` listener rather than the player-facing API port,
- Docker Compose publishes the metrics listener only to host loopback by default (`127.0.0.1:19090`).

This is source/unit instrumentation evidence only. It does **not** prove deployed SLO compliance, dashboard/alert correctness, PostgreSQL/Redis/Unreal metrics coverage, deployment-scale load, long soak, HA/DR, or production readiness.

See [Runtime Observability Metrics v1.3](./docs/runtime-observability-metrics-v1.3.md).

### PostgreSQL restore evidence — v1.4

Runtime Go CI now exercises an isolated PostgreSQL 17 recovery path:
- every canonical migration is applied to a throwaway source database,
- representative account/progression, vehicle/build, quest, inventory, blueprint, and authoritative race state is seeded,
- a custom-format `pg_dump` is restored into a separate database with `pg_restore --exit-on-error`,
- restored operation IDs and authoritative state are asserted,
- every canonical migration is replayed after restore as a schema/application compatibility check,
- the dump and a credential-free report are retained as a 30-day workflow artifact.

This is CI/source-level restore evidence only. It does **not** prove production backup scheduling/custody, encryption, off-host retention, PITR, production-volume restore, accepted RPO/RTO, HA failover, or regional DR.

See [Runtime PostgreSQL Restore Drill v1.4](./docs/runtime-postgres-restore-drill-v1.4.md).

### PostgreSQL pool observability — v1.5

The private metrics surface now includes bounded pgx/PostgreSQL connection-pool telemetry:
- fixed-state connection gauges for max, total, idle, acquired, and constructing connections,
- acquire, empty-acquire, canceled-acquire, and new-connection counters,
- cumulative connection-acquire duration in seconds,
- no SQL text, connection strings, database names, credentials, player IDs, vehicle IDs, race IDs, or peer identifiers are exported,
- the metrics remain read-only and cannot affect PostgreSQL authority or gameplay acceptance.

This is source/unit instrumentation evidence only. It does **not** prove query-level tracing, PostgreSQL server-exporter coverage, deployed scrape/dashboard/alert correctness, production pool sizing, Redis/Unreal observability, SLO compliance, load/soak, HA/DR, or production readiness.

See [Runtime PostgreSQL Observability v1.5](./docs/runtime-postgres-observability-v1.5.md).

### Redis server observability — v1.6

The private metrics surface now includes a bounded Redis server `INFO` probe:
- Redis availability (`zneondrive_redis_up`), connected clients, used/peak memory, rejected connections, evictions, keyspace hit/miss counters and instantaneous operations/second,
- only a fixed allowlist of numeric Redis fields is parsed and exported,
- metric labels are fixed enums only; Redis addresses, keys, client names, replication identifiers, credentials, player/race identifiers and raw probe errors are not exported,
- TCP dial/I/O deadlines bound scrape impact and probe failure emits only `zneondrive_redis_up 0`,
- the probe is observational only and cannot alter rate-limiter decisions, PostgreSQL durable authority or Unreal gameplay authority.

This closes only the Redis server metrics source/unit baseline. It does **not** prove deployed scrape/dashboard/alert correctness, Redis HA/failover, production sizing, deployment-scale load/soak, SLO compliance, PostgreSQL server/query exporter coverage, Unreal runtime metrics, HA/DR, or production readiness.

See [Runtime Redis Server Observability v1.6](./docs/runtime-redis-observability-v1.6.md).

### PostgreSQL server observability — v1.7

The private metrics surface now includes bounded current-database PostgreSQL server statistics:
- probe health plus current backend count,
- committed/rolled-back transactions,
- block reads/cache hits,
- tuple returned/fetched/inserted/updated/deleted counters,
- deadlocks, temporary files/bytes, and current database size,
- only numeric `pg_stat_database` values for `current_database()` are read; database names, SQL text, query fingerprints, relations, connection URLs, credentials and raw query errors are not exported,
- PostgreSQL 17 integration coverage requires a live backend count, positive database size and non-negative counters.

The query is observational/read-only and cannot alter persistence, progression, rate limiting, authoritative race state or Unreal gameplay acceptance. This closes only the bounded PostgreSQL server-metrics source/integration baseline; query-level latency tracing, statement/fingerprint metrics, an external exporter, deployed dashboards/SLO evidence, production sizing/HA and production readiness remain open.

See [Runtime PostgreSQL Server Observability v1.7](./docs/runtime-postgres-server-observability-v1.7.md).

### Integrated trust-boundary evidence — v1.8

Runtime Go CI now exercises key security boundaries together against real PostgreSQL 17 and Redis 8 dependencies:
- direct/untrusted clients cannot evade the bootstrap limiter by rotating `X-Forwarded-For`,
- a configured trusted ingress preserves distinct client identities without weakening authentication,
- forged session credentials remain rejected after ingress identity and distributed rate-limit processing,
- the internal gameplay-ticket redemption route still requires the server-only shared key,
- a valid one-time ticket returns the PostgreSQL-authoritative account/vehicle snapshot,
- replay of the same ticket is rejected through the complete middleware stack.

This is CI-scale service-plane security evidence only. It does **not** prove deployed ingress sanitization/direct-bypass prevention, packaged Unreal↔Go transport, physics anti-cheat, deployment-scale load/soak, deployed SLOs, HA/DR, or production readiness.

See [Runtime Integrated Trust-Boundary Evidence v1.8](./docs/runtime-integrated-trust-boundary-v1.8.md).

### Multi-replica repeated stability evidence — v1.9

Runtime Go integration now repeats the real two-replica HTTP/Redis limiter scenario across multiple fresh identities:
- 12 rounds run against two independent HTTP replicas sharing the real Redis CI service,
- each round sends 64 concurrent bootstrap requests balanced evenly across replicas,
- the shared burst is 16 with refill disabled during the evidence window,
- every round must produce exactly 16 downstream successes and 48 HTTP 429 responses,
- every limited response must include `Retry-After`,
- aggregate evidence covers 768 requests: exactly 192 allowed and 576 rate-limited.

This strengthens CI-scale stability evidence and can detect intermittent coordination drift that a single burst may miss. It is **not** deployment-scale load evidence and is **not** a long-duration soak; those production gates remain open.

See [Runtime Multi-Replica Stability Evidence v1.9](./docs/runtime-multireplica-stability-v1.9.md).

### Unreal dedicated-server observability — v2.0

The UE 5.8 dedicated-server source now has a bounded aggregate telemetry baseline:
- active sessions plus join/leave counters,
- latest and peak server tick duration,
- authoritative driving-input clamp events,
- gameplay-ticket redemption attempts, successes and failures,
- one fixed-format numeric aggregate log record every 10 seconds,
- no player, vehicle, race, ticket, credential, endpoint or peer identifiers retained in the telemetry accumulator.

CI statically validates the required GameMode/PlayerController/VehiclePawn hooks and rejects dynamic/string fields or credential/endpoint reads in the telemetry accumulator. This is **source/static-validation evidence only** until the real UE 5.8 self-hosted runner successfully builds/packages the dedicated server and retained runtime logs prove emission.

Replication rate/bytes, live race-validation metrics, authority corrections from final physics, deployed log shipping/scraping, SLO evidence, load/soak, HA/DR and production readiness remain open.

See [Runtime Unreal Dedicated-Server Observability v2.0](./docs/runtime-unreal-server-observability-v2.0.md).

### Unreal authority activity observability — v2.1

The authoritative vehicle source now extends the bounded aggregate record with:
- authority movement ticks recorded only after the durable-identity gate passes,
- networked durable-identity gate blocks,
- blocking collisions derived from the authoritative swept movement `FHitResult`,
- explicit `ForceNetUpdate()` requests made when durable vehicle state becomes bound.

The static validator requires each counter to remain attached to its real authority-side source hook and continues rejecting dynamic/string identifiers from telemetry. This is **source/static-validation evidence only**. It is not NetDriver packet/byte measurement, live race-validation telemetry, final-physics correction-distance telemetry, or proof that the UE 5.8 server binary emitted these records at runtime.

See [Runtime Unreal Authority Activity Observability v2.1](./docs/runtime-unreal-authority-activity-v2.1.md).

### Unreal authority displacement envelope — v2.2

The authoritative prototype pawn now maintains a server-side position baseline after durable identity is bound. Before each authoritative movement step it compares observed displacement against `MaxSpeedCmPerSecond × DeltaSeconds + AuthorityDisplacementSlackCm`; an over-envelope observation increments a bounded aggregate `impossible_displacements_total` counter. The baseline is reset on authoritative durable-state binding and updated only after the server's swept movement step.

This is telemetry and detection evidence only: it does not accept a client transform, modify durable state, automatically ban a player, or override authoritative movement/race decisions. The static validator requires the envelope comparison to remain behind authority + durable-identity gates, before movement, with the baseline updated after movement. Real UE 5.8 build/runtime emission, final vehicle-physics acceleration/teleport envelopes, live race correlation and sanctions remain open.

See [Runtime Unreal Authority Displacement Envelope v2.2](./docs/runtime-unreal-displacement-envelope-v2.2.md).

### Unreal authority rotation envelope — v2.3

The authoritative prototype pawn now also maintains a server-side yaw baseline after durable identity is bound. Before each authoritative movement step it compares wrap-safe observed yaw delta against `TurnRateDegreesPerSecond × DeltaSeconds + AuthorityRotationSlackDegrees`; an over-envelope observation increments the bounded aggregate `impossible_rotations_total` counter. The baseline is reset during durable identity binding and updated only after authoritative rotation/movement.

This is source/static-validation telemetry evidence only. It does not accept client transforms, mutate durable state, sanction players, or override race/gameplay authority. Real UE 5.8 build/runtime emission, packaged physics calibration, acceleration/angular-velocity envelopes, live race correlation and ranked sanctions remain open.

See [Runtime Unreal Authority Rotation Envelope v2.3](./docs/runtime-unreal-rotation-envelope-v2.3.md).

## Present to a client now

Start with [client/README.md](./client/README.md).

Run the interactive offline presentation demo:

```bash
python3 -m http.server 8080 --directory demo
```

Open `http://localhost:8080`.

The browser demo is presentation-only; the Unreal/Go code is the implementation baseline. Neither is a claim that the full MMORPG is production-ready.

## Documentation

Start with the complete [PROJECT: NEON DRIVE documentation index](./docs/README.md). It covers product requirements, API/data/networking contracts, gameplay systems, security/threat modeling, testing, observability/SLO targets, deployment, backup/restore/DR, incident response, accessibility, localization, moderation, privacy, live ops, release readiness, governance and support.

## Canonical design

- [Original Story Bible / Master Plot](./zNeonDrive-concept.md)
- [Game Design Bible v0.2](./docs/game-design-bible-v0.2.md)
- [Game Design Bible v0.1](./docs/game-design-bible-v0.1.md)
- [Vertical Slice v0.2](./docs/vertical-slice-v0.2.md)
- [Runtime Prototype v0.4](./docs/runtime-prototype-v0.4.md)
- [Runtime Integration v0.5](./docs/runtime-integration-v0.5.md)
- [Runtime Inventory + Rebuild v0.6](./docs/runtime-inventory-rebuild-v0.6.md)
- [Runtime Authoritative Race v0.7](./docs/runtime-authoritative-race-v0.7.md)
- [Runtime Security Hardening v0.8](./docs/runtime-security-hardening-v0.8.md)
- [Runtime Distributed Abuse Controls v0.9](./docs/runtime-distributed-abuse-controls-v0.9.md)
- [Runtime Trusted Ingress Identity v1.0](./docs/runtime-trusted-ingress-v1.0.md)
- [Runtime Multi-Replica HTTP Load Evidence v1.1](./docs/runtime-multi-replica-load-v1.1.md)
- [Runtime Race Integrity Telemetry v1.2](./docs/runtime-race-integrity-telemetry-v1.2.md)
- [Runtime Observability Metrics v1.3](./docs/runtime-observability-metrics-v1.3.md)
- [Runtime PostgreSQL Restore Drill v1.4](./docs/runtime-postgres-restore-drill-v1.4.md)
- [Runtime PostgreSQL Observability v1.5](./docs/runtime-postgres-observability-v1.5.md)
- [Runtime Redis Server Observability v1.6](./docs/runtime-redis-observability-v1.6.md)
- [Runtime PostgreSQL Server Observability v1.7](./docs/runtime-postgres-server-observability-v1.7.md)
- [Runtime Integrated Trust-Boundary Evidence v1.8](./docs/runtime-integrated-trust-boundary-v1.8.md)
- [Runtime Multi-Replica Stability Evidence v1.9](./docs/runtime-multireplica-stability-v1.9.md)
- [Runtime Unreal Dedicated-Server Observability v2.0](./docs/runtime-unreal-server-observability-v2.0.md)
- [Runtime Unreal Authority Activity Observability v2.1](./docs/runtime-unreal-authority-activity-v2.1.md)
- [Runtime Unreal Authority Displacement Envelope v2.2](./docs/runtime-unreal-displacement-envelope-v2.2.md)
- [Runtime Unreal Authority Rotation Envelope v2.3](./docs/runtime-unreal-rotation-envelope-v2.3.md)
- [NOVA CITY World Bible](./docs/nova-city-world-bible.md)
- [Gameplay Systems](./docs/gameplay-systems.md)
- [Architecture](./docs/architecture.md)
- [Content Contracts](./docs/content-contracts.md)
- [Roadmap](./ROADMAP.md)
- [Implementation Checklist](./IMPLEMENTATION-CHECKLIST.md)
- [Complete Documentation Index](./docs/README.md)
- [Governance](./GOVERNANCE.md)
- [Support](./SUPPORT.md)

## Production-direction ADRs

- [ADR-0001 — Server-authoritative state](./docs/adr/0001-server-authoritative-state.md)
- [ADR-0002 — Starter vehicle identity](./docs/adr/0002-starter-vehicle-identity.md)
- [ADR-0003 — Executable reference runtime](./docs/adr/0003-executable-reference-runtime.md)
- [ADR-0004 — Unreal Engine 5.8 playable client + dedicated gameplay server](./docs/adr/0004-unreal-engine-client-and-gameplay-server.md)
- [ADR-0005 — Go/PostgreSQL/Redis durable service plane](./docs/adr/0005-durable-service-plane.md)
- [ADR-0006 — Phase-gated deployment](./docs/adr/0006-phase-gated-deployment.md)

## Machine-readable design catalogs

Content under `design/catalog/` provides stable IDs and implementation contracts:
- 25 named launch characters
- 5 primary factions
- 10 launch districts
- 100 main-story quest records across 7 chapters
- 30 side-content seeds
- vehicle part taxonomy
- functional pet archetypes
- Garage 17 / Foundry 9 vertical-slice proof catalog

Schemas live under `design/schemas/`.

## Evidence layers

### Python reference oracle
`src/zneondrive/` remains a dependency-free contract oracle for core authority invariants.

### Go durable runtime
`services/game-api/` implements the persistent service-plane slice, including durable rebuild, authoritative race lifecycle state, bounded local rate limiting, Redis-coordinated distributed rate-limit state, trusted-ingress identity resolution, concurrent and repeated two-replica HTTP limiter evidence, integrated PostgreSQL + Redis trust-boundary abuse-path evidence, credential-safe race-integrity/auth rejection telemetry, bounded HTTP observability metrics, bounded PostgreSQL pool metrics, bounded Redis server INFO metrics, bounded current-database PostgreSQL server metrics, and isolated PostgreSQL backup/restore verification.

### Unreal runtime
`game/` contains the gameplay-plane source baseline, including aggregate dedicated-server observability for sessions, tick timing, authoritative input clamps, ticket redemption outcomes, authority movement, identity gating, collision blocks, explicit net-update requests, and authority-side displacement/rotation envelope counters. A successful self-hosted UE 5.8 Client/Server build/package run is still required before claiming compiled/runtime Unreal evidence.

Run repository checks:

```bash
make ci
```

## Core loop

```text
Arrive in NOVA CITY
→ receive the broken prototype
→ rebuild the first vehicle
→ take jobs / discover districts / earn reputation
→ tune and reconfigure the same vehicle
→ build relationships and faction standing
→ enter underground and sanctioned competition
→ qualify for NOVA GRAND PRIX
→ expose PROJECT DRIVE ZERO
→ choose the future governance of mobility
→ continue into seasons, crews, world events, and expansions
```

## Fairness rules

1. Ranked outcomes are based on skill, earned/configured performance, event rules, and validated state.
2. VIP cannot add hidden horsepower, grip, durability, matchmaking priority, reward multipliers, or race-stat multipliers.
3. Competitive vehicle builds must be server-validatable and bind results to an exact accepted revision.
4. Rewardable mutations must be idempotent.
5. Reputation and relationship consequences are persistent gameplay state.
6. Endgame choices may alter a player's world state without breaking shared-world consistency.

## Status

**Implemented now:** pre-production design/content contracts, vertical-slice specification, client presentation package, Python authority oracle, Unreal C++ source integration layer, bounded Unreal dedicated-server aggregate telemetry source/static validation including authority movement/identity-gate/collision/net-update activity plus server-side displacement and rotation envelope counters, Go durable service plane, PostgreSQL persistence, one-time gameplay tickets, inventory/blueprint persistence, catalog-validated transactional rebuilds, authoritative PostgreSQL race instances/checkpoints/results, bounded local HTTP rate limiting, Redis-coordinated shared limiter state with local fallback, trusted-proxy source identity resolution with explicit CIDR allowlisting, integrated PostgreSQL + Redis trust-boundary abuse-path evidence, concurrent and repeated two-replica Redis limiter CI evidence, credential-safe rate-limit security events, race-integrity/auth rejection telemetry with hashed correlation buckets, bounded HTTP request/status/latency/in-flight metrics on a separate internal listener, bounded PostgreSQL pgx pool connection/acquire metrics, bounded Redis server INFO metrics, bounded current-database PostgreSQL server metrics, isolated PostgreSQL 17 pg_dump/pg_restore CI recovery evidence, committed Go module lock, HTTP/PostgreSQL reconnect-ticket E2E, Redis limiter integration evidence, local Compose stack, and CI/security validation.

**Still evidence-gated:** successful UE 5.8 source-build artifact, successful UE 5.8 Client/Server cook/package evidence, live packaged Unreal↔Go client/server and race E2E, live Unreal replication-rate/bytes/race-validation/final-physics-correction telemetry, live validation/calibration of physics-derived displacement/rotation/acceleration/teleport envelopes, Garage 17 playable content, final vehicle physics, playable Garage 17 rebuild interaction, relationships/factions, deployed ingress header-sanitization/direct-bypass evidence, deployment-scale multi-replica distributed-limiter load, long-duration soak, ranked anti-cheat and sanctions, matchmaking, deployed metrics/dashboard/SLO evidence, PostgreSQL query-level latency/fingerprint instrumentation and external exporter coverage, Redis production HA/failover evidence, production backup scheduling/custody/PITR/restore, accepted RPO/RTO, HA/regional DR, platform certification, and production deployment.

## License

MIT. See [LICENSE](./LICENSE).