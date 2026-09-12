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
- [ ] full reference-oracle parity in Go
- [ ] live packaged Unreal ↔ Go integration evidence

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
- [ ] first vehicle rebuild playable
- [ ] First Ignition end-to-end in Unreal
- [ ] Foundry 9 traversal/jobs
- [ ] legal + underground race playable
- [ ] Maya relationship visible
- [ ] Luna utility playable
- [ ] durable Unreal save/reconnect verified

## Multiplayer and security
- [ ] Threat model complete
- [x] Client cannot self-assert durable snapshot in v0.5 source contract
- [x] One-time ticket replay rejected by service E2E
- [ ] Client trust boundaries tested over live Unreal↔Go transport
- [ ] Rate limiting
- [ ] Redis replay/idempotency acceleration
- [ ] Cheat telemetry
- [ ] Ranked impossible-state detection
- [ ] Admin/live-ops audit trail
- [ ] Abuse/moderation model

## Reliability
- [ ] Load test target
- [ ] Soak test
- [ ] Backup verification
- [ ] Restore drill
- [ ] RPO/RTO
- [ ] process/node reconnect recovery
- [ ] duplicate-reward prevention under failover
- [ ] content rollback

## Release evidence
- [ ] Stack-specific Unreal + Go CI all green
- [x] Go HTTP/PostgreSQL runtime integration/e2e
- [ ] Live Unreal/Go packaged integration/e2e
- [ ] Security scans pass for runtime
- [ ] Performance budgets
- [ ] Accessibility review
- [ ] Privacy/data retention
- [ ] Production deployment evidence
- [ ] DR evidence
- [ ] Go/no-go approval

Unchecked production claims must not be represented as complete without evidence.
