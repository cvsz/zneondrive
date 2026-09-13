# zNeonDrive Roadmap

## Phase 0 — Vision and repository identity
- [x] Original Story Bible / Master Plot captured
- [x] Repository reframed from generic template to zNeonDrive
- [x] Security/contribution/release repository baseline retained

## Phase 1 — Game Design Bible v0.1
- [x] Product pillars and fairness rules
- [x] 25 named launch characters
- [x] 5 primary factions
- [x] 10 launch districts
- [x] 7-chapter / 100-main-quest campaign catalog
- [x] Side-content seed catalog
- [x] Vehicle/parts taxonomy
- [x] Pet archetypes
- [x] Economy, relationships, crews, and VIP constraints
- [x] Five endgame governance outcomes

## Phase 2 — Content contracts
- [x] Stable ID conventions
- [x] Quest schema
- [x] Vehicle schema
- [x] Content integrity validator
- [x] CI design validation
- [x] Technology-neutral authority boundaries
- [x] ADR: server-authoritative persistent/competitive state
- [x] ADR: preserve starter-vehicle identity

## Phase 3 — Vertical-slice specification v0.2
- [x] Garage 17 + Foundry 9 spaces/race proof targets
- [x] MQ001–MQ012 implementation acceptance criteria
- [x] starter vehicle handling target envelope
- [x] movement/interaction specification
- [x] first-session UX through First Ignition/Roadworthy
- [x] reconnect/save semantics
- [x] accessibility baseline
- [x] graybox/content checklist
- [x] machine-readable vertical-slice proof catalog/schema
- [x] executable reference contract model/tests

## Client presentation readiness v0.3
- [x] executive brief / Thai talk track
- [x] offline interactive browser demo
- [x] Thai/English presentation mode
- [x] commercial scope + Q&A + runbook
- [x] demo validation + manual Pages workflow
- [x] production technology ADRs

## Phase 4 — Runtime prototype / integration v1.4
- [x] Select engine/client — Unreal Engine 5.8
- [x] Select gameplay/service-plane architecture — UE dedicated servers + Go
- [x] Select PostgreSQL/Redis persistence/ephemeral direction
- [x] Select phase-gated deployment strategy
- [x] Create Unreal C++ project / Game / Editor / Server target baseline
- [x] Add manual self-hosted UE 5.8 source-build workflow
- [x] Support UE 5.8 source-tree and installed-build Linux tooling
- [x] Prepare retained Client/Server build-evidence workflow with logs, engine metadata and checksums
- [ ] Produce successful retained UE 5.8 Client/Server build evidence from a real self-hosted runner
- [ ] Produce successful retained UE 5.8 Client/Server package/cook evidence
- [x] Implement prototype identity/session bootstrap against PostgreSQL
- [x] Implement one-primary-character persistence
- [x] Implement starter-vehicle persistence and immutable build revisions
- [x] Implement optimistic/idempotent durable mutation semantics
- [x] Implement inventory/blueprint prototype
- [x] Implement sequential quest-state persistence
- [x] Implement authoritative race instance/result runtime in Go/PostgreSQL
- [x] Bind race instances to the exact active build revision + validation hash
- [x] Enforce ordered/monotonic checkpoint acceptance and idempotent race writes
- [x] Add bounded per-process token-bucket abuse controls for public/internal HTTP mutations
- [x] Hash credential-derived limiter identities and bound limiter memory cardinality
- [x] Add Redis-backed distributed rate-limit coordination for multi-process deployment
- [x] Retain bounded per-process limiting as Redis-unavailable fallback
- [x] Add credential-safe rate-limit rejection/fallback security events
- [x] Implement trusted-proxy/ingress client identity policy and tests
- [x] Fail closed on invalid proxy CIDR configuration and malformed forwarding chains
- [x] Add concurrent two-replica HTTP integration evidence for shared Redis limiter budget
- [x] Add credential-safe race-integrity/auth rejection telemetry baseline
- [x] Add bounded Go HTTP request/status/latency/in-flight metrics with separate internal listener
- [x] Add isolated PostgreSQL 17 pg_dump/pg_restore CI restore drill with representative durable-state assertions
- [ ] Implement relationship/faction-state runtime
- [x] Add Go unit tests
- [x] Add PostgreSQL integration test suite
- [x] Add Redis integration coverage for shared distributed limiter budget
- [x] Commit deterministic Go module lock and enforce tidy-clean CI
- [x] Add one-time gameplay ticket issue/redeem protocol
- [x] Add server-only shared-key redemption boundary
- [x] Add HTTP/PostgreSQL reconnect + gameplay-ticket E2E
- [x] Add Unreal session subsystem + dedicated-server ticket redemption source
- [x] Bind durable vehicle/build/Roadworthy state only on authority
- [x] Enforce canonical part catalog on build mutations
- [x] Make MQ004/MQ009 item grants idempotent and MQ005 blueprint unlock durable
- [x] Make rebuild inventory consume/return atomic with immutable build revision
- [ ] Complete Python-reference parity tests in Go
- [ ] Produce live Unreal ↔ Go packaged/session E2E evidence
- [ ] Produce live Unreal ↔ Go race lifecycle evidence

## Repository documentation readiness
- [x] GitHub governance / support / maintainers baseline
- [x] canonical documentation index + evidence status language
- [x] product requirements and brand guide
- [x] API / data / networking contracts synchronized through race runtime v0.7
- [x] gameplay authoring specs for quests, vehicles/races, world, factions, crews and companions
- [x] threat model, testing strategy and performance-budget targets
- [x] runtime security hardening v0.8 evidence/non-claim contract
- [x] distributed abuse-controls v0.9 evidence/non-claim contract
- [x] trusted-ingress identity v1.0 evidence/non-claim contract
- [x] multi-replica HTTP load v1.1 CI evidence/non-claim contract
- [x] race-integrity telemetry v1.2 evidence/non-claim contract
- [x] observability metrics v1.3 source/unit evidence contract
- [x] PostgreSQL restore drill v1.4 CI evidence/non-claim contract
- [x] UE 5.8 retained build-evidence workflow contract
- [x] observability/SLO, deployment, backup/restore/DR and incident-response plans
- [x] accessibility, localization, content pipeline and economy/fairness policy
- [x] moderation/player safety, privacy/data retention and live-ops policy
- [x] release-readiness, asset/IP and glossary docs
- [x] bug/documentation issue templates
- [x] README 1280×640 banner asset
- [x] documentation completeness + relative-link CI validator

> Documentation readiness does not close runtime evidence gates such as successful real UE 5.8 Client/Server build/package, live Unreal↔Go E2E, deployed ingress correctness, physics-derived anti-cheat, deployment-scale load/soak, deployed SLO measurement, production restore/DR, or production deployment.

## Phase 5 — Multiplayer alpha
- [ ] Multi-player district instance/shard
- [ ] Presence/social
- [ ] Crew MVP
- [ ] Match/event registration
- [x] Race rejection/auth telemetry baseline
- [ ] Physics-derived race anti-cheat telemetry
- [ ] Disconnect/rejoin behavior
- [ ] Deployment-scale load and long-duration soak testing
- [x] Isolated PostgreSQL backup/restore CI drill
- [ ] Production backup/restore drill against an agreed deployment target

## Phase 6 — Content alpha
- [ ] Chapter 1 fully playable in selected runtime
- [ ] Street jobs and early side quests
- [ ] Luna companion questline
- [ ] Foundry 9 and Neon Mile content
- [ ] Early faction reputation
- [ ] First championship qualifier
- [ ] Localization pipeline

## Phase 7 — Production hardening
- [x] Threat-model baseline documented
- [x] Per-process HTTP abuse limiter implemented and unit tested
- [x] Redis-backed limiter coordination exercised across independent limiter instances
- [x] Credential-safe rate-limit rejection/fallback events emitted
- [x] Trusted-ingress identity handling implemented and unit exercised
- [ ] Deployed ingress forwarding/sanitization configuration verified
- [x] Distributed abuse controls verified under concurrent two-replica HTTP CI load
- [ ] Deployment-scale distributed limiter load evidence
- [ ] Long-duration soak evidence
- [x] Authoritative race rejection + game-server auth telemetry emitted with hashed correlation buckets
- [ ] Physics-derived impossible-state telemetry from live Unreal race samples
- [x] Low-cardinality Go HTTP service metrics source/unit baseline
- [ ] PostgreSQL/Redis/Unreal runtime metrics completed
- [ ] Threat/abuse cases exercised against integrated runtime
- [ ] SLOs/observability measured in deployed environment
- [x] Incident/rollback runbook baseline documented
- [ ] Incident/rollback exercise passed
- [ ] Capacity model backed by load evidence
- [x] Backup/restore/DR plan documented
- [x] PostgreSQL 17 isolated CI backup/restore drill with retained artifact/report
- [ ] Production restore/DR evidence
- [ ] Security review of complete runtime
- [x] Privacy/data-retention baseline documented
- [ ] Privacy/legal review for target launch regions
- [ ] VIP fairness test suite in production runtime
- [ ] Ranked race integrity / anti-cheat suite

## Phase 8 — Launch readiness
- [ ] 100 main quests implemented and verified in runtime
- [ ] All launch districts implemented
- [ ] NOVA GRAND PRIX seasonal operations
- [ ] Crew championship
- [ ] Live-ops tools and audit controls
- [ ] Customer support/moderation workflows
- [ ] Production deployment evidence
- [ ] Disaster-recovery evidence
- [ ] Go/no-go release review

**Rule:** do not mark playable content, production, HA, DR, anti-cheat, deployment-scale distributed abuse protection, production load/soak, measured SLOs, or live-ops complete without executable evidence.
