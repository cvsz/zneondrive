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

## Phase 4 — Runtime prototype / integration v0.6
- [x] Select engine/client — Unreal Engine 5.8
- [x] Select gameplay/service-plane architecture — UE dedicated servers + Go
- [x] Select PostgreSQL/Redis persistence/ephemeral direction
- [x] Select phase-gated deployment strategy
- [x] Create Unreal C++ project / Game / Editor / Server target baseline
- [x] Add manual self-hosted UE 5.8 source-build workflow
- [ ] Produce successful archived UE source-build evidence
- [x] Implement prototype identity/session bootstrap against PostgreSQL
- [x] Implement one-primary-character persistence
- [x] Implement starter-vehicle persistence and immutable build revisions
- [x] Implement optimistic/idempotent durable mutation semantics
- [x] Implement inventory/blueprint prototype
- [x] Implement sequential quest-state persistence
- [ ] Implement authoritative race instance/result runtime
- [ ] Implement relationship/faction-state runtime
- [x] Add Go unit tests
- [x] Add PostgreSQL integration test suite
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

## Repository documentation readiness v0.7
- [x] GitHub governance/support/maintainers baseline
- [x] complete documentation index
- [x] product/API/data/networking contracts
- [x] security/threat/testing/performance documentation
- [x] deployment/observability/incident/DR plans
- [x] accessibility/localization/content/monetization/safety/privacy/live-ops policies
- [x] release-readiness/IP/glossary/brand guidance
- [x] bug/documentation issue templates
- [x] README banner asset

> Documentation readiness does not close runtime evidence gates such as live UE build, anti-cheat, load/soak, restore/DR or production deployment.

## Phase 5 — Multiplayer alpha
- [ ] Multi-player district instance/shard
- [ ] Presence/social
- [ ] Crew MVP
- [ ] Match/event registration
- [ ] Race anti-cheat telemetry
- [ ] Disconnect/rejoin behavior
- [ ] Load and soak testing
- [ ] Backup/restore drill

## Phase 6 — Content alpha
- [ ] Chapter 1 fully playable in selected runtime
- [ ] Street jobs and early side quests
- [ ] Luna companion questline
- [ ] Foundry 9 and Neon Mile content
- [ ] Early faction reputation
- [ ] First championship qualifier
- [ ] Localization pipeline

## Phase 7 — Production hardening
- [ ] Threat model
- [ ] Abuse/cheat model
- [ ] SLOs and observability
- [ ] Incident/rollback runbooks
- [ ] Capacity model
- [ ] Restore evidence
- [ ] Security review
- [ ] Privacy/data-retention review
- [ ] VIP fairness test suite in production runtime
- [ ] Ranked race integrity test suite

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

**Rule:** do not mark playable content, production, HA, DR, anti-cheat, or live-ops complete without executable evidence.
