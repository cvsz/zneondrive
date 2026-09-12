# Changelog

All notable changes to zNeonDrive are documented here.

## [Unreleased]

### Added
- Canonical Game Design Bible v0.1
- NOVA CITY world bible and gameplay systems specification
- Technology-neutral authoritative online-world architecture
- Stable content-ID and runtime contract guidance
- 25-character launch roster
- Five-faction and ten-district catalogs
- 100-record, seven-chapter main campaign graph
- 30 side-content seed records
- Vehicle part taxonomy and companion catalog
- Quest and vehicle JSON schemas
- Automated design integrity validator
- Architecture Decision Records for authority and starter-vehicle identity

### Changed
- Replaced generic template README, roadmap, architecture, and implementation checklist with zNeonDrive-specific material
- Extended CI to validate JSON design catalogs and cross-reference invariants
- Reframed repository status around evidence-gated design, prototype, alpha, and production milestones

### Security
- Player clients, client clocks, rewards, build legality, and race results are explicitly untrusted until authoritative validation
