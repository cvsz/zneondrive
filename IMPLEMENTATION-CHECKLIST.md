# Implementation Checklist — PROJECT: NEON DRIVE

## Design baseline
- [x] Story premise / NOVA CITY 2097
- [x] Garage 17 opening and starter prototype
- [x] Maya / Adrian / Victor / Luna anchors
- [x] PROJECT DRIVE ZERO conflict
- [x] seven campaign chapters / five endings
- [x] VIP fairness, social/crew, pets
- [x] 25 characters / 100 main quests / world catalogs
- [x] automated design validation
- [x] Game Design Bible v0.2

## Vertical-slice contract
- [x] MQ001–MQ012 acceptance matrix
- [x] Garage 17 / Foundry 9 graybox scope
- [x] first-session UX
- [x] starter handling target
- [x] movement/interaction minimum
- [x] legal + underground race proof definitions
- [x] save/reconnect semantics
- [x] accessibility baseline
- [x] machine-readable slice catalog/schema

## Executable reference authority
- [x] Python contract oracle
- [x] one account → one character
- [x] vehicle ownership
- [x] append-only build revisions
- [x] quest reward idempotency
- [x] race/build binding reference rule
- [x] ordered checkpoint reference rule
- [x] VIP separation reference rule
- [x] starter deletion protection reference rule

## Client meeting package
- [x] executive brief / talk track
- [x] interactive offline demo
- [x] runbook / commercial scope / Q&A
- [x] Thai/English toggle
- [x] client-demo CI / manual Pages
- [x] explicit evidence language

## Production technology decisions
- [x] Unreal Engine 5.8 gameplay ADR
- [x] Go + PostgreSQL + Redis service-plane ADR
- [x] phase-gated deployment ADR

## Documentation / governance baseline
- [x] complete documentation index
- [x] product requirements
- [x] API + durable data model docs through authoritative race v0.7
- [x] multiplayer/networking authority spec
- [x] threat-model baseline (production abuse testing still open)
- [x] testing strategy + performance-budget targets
- [x] observability/SLO target contract (measurement still open)
- [x] deployment + backup/restore/DR plans (verification still open)
- [x] incident-response runbook (exercise still open)
- [x] accessibility + localization baselines
- [x] content/quest/world/vehicle/faction/crew/companion authoring specs
- [x] economy/fairness + moderation + privacy + live-ops policies
- [x] governance + support + maintainers
- [x] release-readiness + asset/IP + brand guidance
- [x] documentation completeness + relative-link validator

## Phase 4 runtime source
- [x] Unreal .uproject baseline
- [x] Unreal Game / Editor / dedicated Server targets
- [x] server-authoritative replicated prototype pawn
- [x] manual self-hosted UE source-build workflow
- [ ] successful UE 5.8 source-build artifact
- [x] Go 1.27 service binary
- [x] PostgreSQL durable schema
- [x] local Docker Compose PostgreSQL + Redis + API
- [x] hashed resume/session credentials
- [x] durable one-character bootstrap
- [x] durable starter vehicle + immutable revisions
- [x] idempotent quest/build operations
- [x] MQ001–MQ100 prerequisite enforcement
- [x] MQ012 Roadworthy transition
- [x] Go unit tests
- [x] PostgreSQL integration tests
- [x] committed go.mod/go.sum module lock
- [x] one-time gameplay-ticket persistence and atomic redemption
- [x] server-only shared-key internal redemption endpoint
- [x] HTTP/PostgreSQL reconnect + ticket E2E
- [x] Unreal session/resume subsystem source
- [x] dedicated-server ticket redemption source
- [x] authority-only durable VehicleID/build/parts/Roadworthy binding
- [x] PostgreSQL inventory + blueprint persistence
- [x] canonical vehicle-part catalog validation in Go runtime
- [x] MQ004/MQ009 idempotent item grants
- [x] MQ005 durable starter-rebuild blueprint unlock
- [x] atomic inventory consume/return + immutable rebuild revision
- [x] build operation replay bound to exact validation hash
- [x] reconnect persistence for inventory/blueprints/rebuild state
- [x] Unreal snapshot source parses inventory/blueprints
- [x] PostgreSQL race instance/checkpoint/result persistence
- [x] race start bound to owned Roadworthy vehicle + exact active build revision/hash
- [x] race lifecycle endpoints restricted to dedicated-server shared-key boundary
- [x] ordered checkpoint cursor + monotonic elapsed-time enforcement
- [x] idempotent race start/checkpoint/finish operation semantics
- [x] deterministic final result hash bound to authoritative build evidence
- [ ] full reference-oracle parity in Go
- [ ] live packaged Unreal ↔ Go integration evidence
- [ ] live packaged Unreal ↔ Go race lifecycle evidence

## Content production
- [ ] Full dialogue/script MQ001–MQ100
- [ ] Side-quest narratives
- [ ] Cinematic/storyboard list
- [ ] Environment storytelling assets
- [ ] Voice/localization guide
- [ ] Rating/sensitivity review

## Playable runtime vertical slice
- [x] selected client/server ADRs
- [x] Unreal source project baseline
- [ ] Garage 17 environment playable
- [ ] first vehicle rebuild playable (durable service contract implemented; Unreal interaction/evidence still open)
- [ ] First Ignition end-to-end in Unreal
- [ ] Foundry 9 traversal/jobs
- [ ] legal + underground race playable (service-plane race authority exists; live Unreal gameplay/evidence still open)
- [ ] Maya relationship visible
- [ ] Luna utility playable
- [ ] durable Unreal save/reconnect verified

## Multiplayer and security
- [ ] Threat model complete and exercised
- [x] Threat-model baseline documented
- [x] Client cannot self-assert durable snapshot in v0.5 source contract
- [x] One-time ticket replay rejected by service E2E
- [x] Durable race result/build identity cannot be submitted directly by player-facing endpoints
- [x] Ordered/monotonic race checkpoint acceptance enforced in PostgreSQL runtime
- [ ] Client trust boundaries tested over live Unreal↔Go transport
- [ ] Rate limiting
- [ ] Redis replay/idempotency acceleration
- [ ] Cheat telemetry
- [ ] Ranked impossible-state detection
- [ ] Admin/live-ops audit trail
- [ ] Abuse/moderation runtime

## Reliability
- [ ] Load test executed against agreed target
- [ ] Soak test
- [ ] Backup verification
- [ ] Restore drill
- [ ] RPO/RTO accepted and demonstrated
- [ ] process/node reconnect recovery
- [ ] duplicate-reward prevention under failover
- [ ] content rollback

## Release evidence
- [ ] Stack-specific Unreal + Go CI all green
- [x] Go HTTP/PostgreSQL runtime integration/e2e
- [ ] Live Unreal/Go packaged integration/e2e
- [ ] Live Unreal/Go race lifecycle e2e
- [ ] Security scans pass for complete runtime surface
- [ ] Performance budgets measured/passed
- [ ] Accessibility review
- [ ] Privacy/data-retention review
- [ ] Production deployment evidence
- [ ] DR evidence
- [ ] Go/no-go approval

Unchecked production claims must not be represented as complete without evidence.
