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

### Changed
- Replaced generic template README, roadmap, architecture, and implementation checklist with zNeonDrive-specific material
- Advanced Phase 3 vertical-slice specification to complete while keeping playable/runtime claims evidence-gated
- Extended CI to validate JSON design catalogs, client-demo integrity, reference-runtime tests, and Python compilation
- Advanced Phase 4 technology selection while leaving implementation evidence explicitly open
- Reframed repository status around evidence-gated design, presentation, prototype, alpha, and production milestones

### Security
- Player clients, client clocks, rewards, build legality, and race results are explicitly untrusted until authoritative validation
- Rewardable reference mutations require idempotency semantics and races bind results to the accepted build revision
- Client presentation demo has no external CDN, analytics, API key, or network dependency
