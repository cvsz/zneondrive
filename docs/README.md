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
- [Authoritative Race Runtime v0.7](./runtime-authoritative-race-v0.7.md)
- [Runtime Security Hardening v0.8](./runtime-security-hardening-v0.8.md)
- [Runtime Distributed Abuse Controls v0.9](./runtime-distributed-abuse-controls-v0.9.md)
- [Runtime Trusted Ingress Identity v1.0](./runtime-trusted-ingress-v1.0.md)
- [Runtime Multi-Replica HTTP Evidence v1.1](./runtime-multi-replica-load-v1.1.md)
- [Runtime Race Integrity Telemetry v1.2](./runtime-race-integrity-telemetry-v1.2.md)
- [Runtime Observability Metrics v1.3](./runtime-observability-metrics-v1.3.md)
- [Runtime PostgreSQL Restore Drill v1.4](./runtime-postgres-restore-drill-v1.4.md)
- [Runtime PostgreSQL Observability v1.5](./runtime-postgres-observability-v1.5.md)
- [Runtime Redis Limiter Observability v1.5a](./runtime-redis-limiter-observability-v1.5a.md)
- [Runtime Redis Server Observability v1.6](./runtime-redis-observability-v1.6.md)
- [Runtime PostgreSQL Server Observability v1.7](./runtime-postgres-server-observability-v1.7.md)
- [Runtime Integrated Trust-Boundary Evidence v1.8](./runtime-integrated-trust-boundary-v1.8.md)
- [Runtime Multi-Replica Stability Evidence v1.9](./runtime-multireplica-stability-v1.9.md)
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
- [Full-Stack Install & Control Panel](./control-panel.md)
- [Roadmap](../ROADMAP.md)
- [Implementation Checklist](../IMPLEMENTATION-CHECKLIST.md)
- [Changelog](../CHANGELOG.md)
- [Governance](../GOVERNANCE.md)
- [Support](../SUPPORT.md)
- [Security](../SECURITY.md)
- [Glossary](./glossary.md)

## Evidence boundary

The repository currently has implementation evidence for the Python authority oracle, Go/PostgreSQL service-plane slices, bounded local HTTP abuse controls, Redis-coordinated shared limiter state, trusted-ingress client-identity source handling, CI-scale concurrent and repeated HTTP traffic across two independent limiter instances sharing Redis, integrated PostgreSQL + Redis trust-boundary abuse-path tests, credential-safe authoritative race/auth rejection telemetry, low-cardinality Go HTTP request/latency/in-flight metrics, bounded PostgreSQL connection-pool metrics, bounded process-local Redis limiter outcome telemetry, bounded Redis server INFO metrics, bounded current-database PostgreSQL server metrics, isolated PostgreSQL restore verification, content/design validators, and source-level Unreal integration. The repeated multi-replica stability and integrated trust-boundary evidence remain CI-scale service-plane evidence only. The repository does not yet claim deployed ingress/header-sanitization evidence, deployed metrics/dashboard/SLO evidence, PostgreSQL query-level/external-exporter coverage, Redis HA/failover verification, Unreal runtime metrics, deployment-scale or long-duration load/soak performance, packaged Unreal↔Go production E2E, final vehicle physics, physics-derived ranked anti-cheat, HA/DR, or production deployment.
