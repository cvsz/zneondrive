# Changelog

All notable changes to zNeonDrive are documented here.

## [Unreleased]

### Added
- Canonical Game Design Bible v0.1
- Game Design Bible v0.2 implementation contract layer
- NOVA CITY world bible and gameplay systems specification
- Garage 17 / Foundry 9 MQ001–MQ012 vertical-slice specification
- Machine-readable vertical-slice proof catalog and schema
- Technology-neutral authoritative online-world architecture
- Stable content-ID and runtime contract guidance
- 25-character launch roster
- Five-faction and ten-district catalogs
- 100-record, seven-chapter main campaign graph
- 30 side-content seed records
- Vehicle part taxonomy and companion catalog
- Quest and vehicle JSON schemas
- Automated design integrity validator
- Dependency-free executable reference domain model and unit tests
- ADR-0004 selecting Unreal Engine 5.8 for playable client and dedicated gameplay servers
- ADR-0005 selecting Go/PostgreSQL/Redis durable service-plane direction
- ADR-0006 defining phase-gated deployment and observability
- Offline interactive client presentation demo
- Thai/English presentation mode
- Executive brief, talk track, demo runbook, commercial scope options, and prepared client Q&A
- Automated client-demo validation
- Manual GitHub Pages deployment workflow
- Runtime Prototype v0.4 documentation
- Unreal Engine 5.8 C++ project with Game, Editor, and dedicated Server targets
- Server-authoritative replicated prototype vehicle pawn
- Manual self-hosted Unreal source-build workflow
- Go 1.27 service-plane module and HTTP API
- PostgreSQL durable schema for accounts, characters, vehicles, immutable builds, sessions, and quest completions
- Hashed resume/session credential handling
- Idempotent quest/build mutation semantics
- Sequential MQ001–MQ100 durable quest gate and MQ012 Roadworthy transition
- PostgreSQL integration tests and Go runtime CI workflow
- Docker Compose local PostgreSQL 17 + Redis 8 + game API stack
- Phase 4 static runtime validator
- Runtime Integration v0.5 trust-boundary specification
- Unreal GameInstance service subsystem for bootstrap/resume/state/quest/build flows
- One-time 60-second gameplay ticket issue/redeem protocol
- Server-only shared-key internal ticket redemption endpoint
- Atomic single-use gameplay-ticket persistence and migration
- Dedicated-server PlayerController ticket redemption and authority-only durable pawn binding
- Replicated durable VehicleID, build revision, active parts, Roadworthy, and binding state
- HTTP/PostgreSQL reconnect + single-use-ticket E2E test
- Committed go.mod/go.sum dependency lock with tidy-clean CI enforcement
- PostgreSQL advisory lock around concurrent schema setup
- Runtime Inventory + Rebuild v0.6 contract and evidence documentation
- PostgreSQL inventory item and character blueprint persistence
- Canonical runtime vehicle-part allowlist checked against the design catalog
- MQ004 idempotent salvage grant, MQ005 rebuild-blueprint unlock, and MQ009 idempotent recovered-part grant
- Inventory-authoritative rebuild transaction with atomic consume/return accounting
- Unreal inventory/blueprint snapshot parsing
- PostgreSQL integration proof for duplicate-grant prevention, rebuild replay safety, unowned-part rejection, swap accounting, and reconnect persistence

- Repository documentation index and status-language contract
- Product requirements and brand guide
- API contract grounded in current Go HTTP routes
- Durable PostgreSQL data-model documentation
- Multiplayer/networking authority specification
- Threat model and security-production exit criteria
- Testing strategy and performance-budget targets
- Observability/SLO target contract
- Deployment, backup/restore/DR, and incident-response runbooks
- Accessibility and localization baselines
- Content pipeline, quest authoring, world-streaming, vehicle/race-integrity, faction/relationship, crew/social, and companion specs
- Economy/monetization fairness, moderation/player-safety, privacy/data-retention, and live-ops policies
- Release-readiness checklist, asset/IP policy, glossary, governance, support, and maintainers documents
- GitHub bug-report and documentation issue templates
- Repository 1280×640 SVG README banner asset

### Changed
- Replaced generic template README, roadmap, architecture, and implementation checklist with zNeonDrive-specific material
- Advanced Phase 3 vertical-slice specification to complete while keeping playable/runtime claims evidence-gated
- Extended CI to validate JSON design catalogs, client-demo integrity, reference-runtime tests, runtime source structure, and Python compilation
- Advanced Phase 4 from technology selection into executable Unreal/Go/PostgreSQL source
- Reframed repository status around evidence-gated design, presentation, runtime prototype, alpha, and production milestones
- Local runtime ports default to 55432/56379/18080 to reduce collisions with common PostgreSQL/Redis/dev ports
- Advanced Phase 4.1 to source-level Unreal↔Go integration while keeping live Unreal build/session evidence open
- Advanced Phase 4.2 to durable inventory/blueprint and rebuild-state integration without claiming playable Garage 17 evidence

### Security
- Player clients, client clocks, rewards, build legality, and race results are explicitly untrusted until authoritative validation
- Rewardable mutations require idempotency semantics
- Resume keys and session tokens are stored as hashes in PostgreSQL
- Vehicle build mutations use optimistic revision checks and row locking
- Unreal driving inputs are clamped on the authority before movement
- Clients never receive or read the gameplay-server shared key
- Durable snapshots are redeemed by the dedicated server rather than accepted from a client RPC
- Gameplay tickets are hashed, short-lived, and single-use
- Vehicle builds reject unknown catalog IDs, locked blueprints, and parts not owned in inventory
- Inventory consumption/returns and build revision activation commit atomically
- Replayed build operation IDs must match the original vehicle and build validation hash
- Client presentation demo has no external CDN, analytics, API key, or network dependency
