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
- Architecture Decision Records for authority, starter-vehicle identity, and executable reference contracts
- Dependency-free executable reference domain model for ownership, build revisions, idempotent quest rewards, VIP fairness, and race validation
- Unit tests covering seven authoritative domain invariants

### Changed
- Replaced generic template README, roadmap, architecture, and implementation checklist with zNeonDrive-specific material
- Advanced Phase 3 vertical-slice specification to complete while keeping playable/runtime claims evidence-gated
- Extended CI to validate JSON design catalogs, cross-reference invariants, reference-runtime tests, and Python compilation
- Reframed repository status around evidence-gated design, prototype, alpha, and production milestones

### Security
- Player clients, client clocks, rewards, build legality, and race results are explicitly untrusted until authoritative validation
- Rewardable reference mutations require idempotency semantics and races bind results to the accepted build revision
