# PROJECT: NEON DRIVE Documentation Index

This directory is the canonical documentation map for zNeonDrive.

## Status language

Every engineering document should distinguish:

- **Implemented** — source exists and applicable automated evidence passes.
- **Verified** — the implemented behavior has the required runtime/integration/operational evidence.
- **Target** — agreed design objective, not yet proven.
- **Planned** — backlog direction, not an implementation claim.

## Product and world

- [Game Design Bible v0.2](./game-design-bible-v0.2.md)
- [Original concept / Story Bible](../zNeonDrive-concept.md)
- [NOVA CITY World Bible](./nova-city-world-bible.md)
- [Gameplay Systems](./gameplay-systems.md)
- [Product Requirements](./product-requirements.md)
- [Vertical Slice v0.2](./vertical-slice-v0.2.md)
- [Economy & Monetization](./economy-monetization.md)
- [Content Pipeline](./content-pipeline.md)
- [Localization](./localization.md)
- [Accessibility](./accessibility.md)
- [Moderation & Player Safety](./moderation-safety.md)
- [Asset & IP Policy](./asset-ip-policy.md)
- [Brand Guide](./brand-guide.md)
- [Quest Authoring](./quest-authoring.md)
- [Vehicle Physics & Race Integrity](./vehicle-physics-and-race-integrity.md)
- [World Streaming & Districts](./world-streaming-and-districts.md)
- [Factions & Relationships](./factions-relationships.md)
- [Crew & Social Systems](./crew-social.md)
- [Pets & Companions](./pets-companions.md)

## Architecture and runtime

- [Architecture](./architecture.md)
- [API Contract](./api-contract.md)
- [Data Model](./data-model.md)
- [Multiplayer & Networking](./multiplayer-networking.md)
- [Runtime Prototype v0.4](./runtime-prototype-v0.4.md)
- [Runtime Integration v0.5](./runtime-integration-v0.5.md)
- [Inventory + Rebuild v0.6](./runtime-inventory-rebuild-v0.6.md)
- [Content Contracts](./content-contracts.md)
- [ADR Index](./adr/README.md)

## Quality, security and operations

- [Threat Model](./threat-model.md)
- [Testing Strategy](./testing-strategy.md)
- [Performance Budget](./performance-budget.md)
- [Observability & SLOs](./observability-slo.md)
- [Deployment](./deployment.md)
- [Backup, Restore & DR](./backup-restore-dr.md)
- [Incident Response](./incident-response.md)
- [Privacy & Data Retention](./privacy-data-retention.md)
- [Live Operations](./live-ops.md)
- [Release](./release.md)
- [Release Readiness Checklist](./release-readiness-checklist.md)

## Project management and contribution

- [Development](./development.md)
- [Roadmap](../ROADMAP.md)
- [Implementation Checklist](../IMPLEMENTATION-CHECKLIST.md)
- [Changelog](../CHANGELOG.md)
- [Governance](../GOVERNANCE.md)
- [Support](../SUPPORT.md)
- [Security](../SECURITY.md)
- [Glossary](./glossary.md)

## Evidence boundary

The repository currently has implementation evidence for the Python authority oracle, Go/PostgreSQL service-plane slices, content/design validators, and source-level Unreal integration. It does not yet claim verified packaged Unreal↔Go production E2E, final vehicle physics, ranked anti-cheat, load/soak, HA/DR, or production deployment.
