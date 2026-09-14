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
- [Runtime Unreal Dedicated-Server Observability v2.0](./runtime-unreal-server-observability-v2.0.md)
- [Runtime Unreal Authority Activity Observability v2.1](./runtime-unreal-authority-activity-v2.1.md)
- [Runtime Unreal Authority Displacement Envelope v2.2](./runtime-unreal-displacement-envelope-v2.2.md)
- [Runtime Unreal Authority Rotation Envelope v2.3](./runtime-unreal-rotation-envelope-v2.3.md)
- [Runtime Reference Parity v2.4](./runtime-reference-parity-v2.4.md)
- [Runtime Reference Parity Integrity v2.5](./runtime-reference-parity-integrity-v2.5.md)
- [Runtime Reference Operation Parity v2.6](./runtime-reference-operation-parity-v2.6.md)
- [Runtime Reference Rebuild Parity v2.7](./runtime-reference-rebuild-parity-v2.7.md)
- [Reference Parity Evidence Index v2.8](./runtime-reference-parity-index-v2.8.md)
- [Runtime Reference Race Input Parity v2.9](./runtime-reference-race-parity-v2.9.md)
- [Runtime PostgreSQL Query Activity Evidence v3.1](./runtime-postgres-query-activity-v3.1.md)
- [Runtime PostgreSQL External Exporter Evidence v3.2](./runtime-postgres-exporter-v3.2.md)
- [Runtime PostgreSQL Exporter Least-Privilege Evidence v3.3](./runtime-postgres-exporter-least-privilege-v3.3.md)
- [Runtime PostgreSQL Statement Metrics Evidence v3.4](./runtime-postgres-statements-v3.4.md)
- [Runtime Reference Race Result Parity v3.5](./runtime-reference-race-result-parity-v3.5.md)
- [Runtime Reference Entitlement Parity v3.6](./runtime-reference-entitlement-parity-v3.6.md)
- [Runtime Reference Quest Sequence Parity v3.7](./runtime-reference-quest-sequence-parity-v3.7.md)
- [Runtime Reference Quest Reward Parity v3.8](./runtime-reference-quest-reward-parity-v3.8.md)
- [Runtime Reference Quest Reward Authority v3.9](./runtime-reference-quest-reward-authority-v3.9.md)
- [Runtime Reference Race Lifecycle Parity v4.0](./runtime-reference-race-lifecycle-parity-v4.0.md)
- [Runtime Reference Race Start Parity v4.1](./runtime-reference-race-start-parity-v4.1.md)
- [Runtime Reference Race Operation Parity v4.2](./runtime-reference-race-operation-parity-v4.2.md)
- [Runtime Reference Quest Side Effects Parity v4.3](./runtime-reference-quest-side-effects-parity-v4.3.md)
- [Runtime Reference Build Operation Replay v4.4](./runtime-reference-build-operation-replay-v4.4.md)
- [Runtime Reference Quest Operation Scope v4.5](./runtime-reference-quest-operation-scope-v4.5.md)
- [Runtime Service Restart Recovery Evidence v4.6](./runtime-service-restart-recovery-v4.6.md)
- [Content Contracts](./content-contracts.md)
- [ADR Index](./adr/README.md)

## Quality, security and operations

- [Threat Model](./threat-model.md)
- [Security Threat Exercise v3.0](./security-threat-exercise-v3.0.md)
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

The repository currently has implementation evidence for the Python authority oracle, Go/PostgreSQL service-plane slices, bounded local HTTP abuse controls, Redis-coordinated shared limiter state, trusted-ingress client-identity source handling, CI-scale concurrent and repeated HTTP traffic across two independent limiter instances sharing Redis, integrated PostgreSQL + Redis trust-boundary abuse-path tests, a machine-checked partial threat-exercise manifest tying selected abuse scenarios to concrete CI controls/tests, credential-safe authoritative race/auth rejection telemetry, low-cardinality Go HTTP request/latency/in-flight metrics, bounded PostgreSQL connection-pool metrics, bounded process-local Redis limiter outcome telemetry, bounded Redis server INFO metrics, bounded current-database PostgreSQL server metrics, bounded query-text-free PostgreSQL `pg_stat_activity` pressure metrics, CI-exercised external PostgreSQL exporter coverage using a dedicated least-privilege `pg_monitor` principal, bounded query-text-free PostgreSQL `pg_stat_statements` exporter coverage capped at 25 statement rows per scrape against a 500-entry PostgreSQL statement table, isolated PostgreSQL restore verification, service-process restart/reconnect recovery against the same PostgreSQL durable state, content/design validators, source-level Unreal integration, bounded Unreal dedicated-server aggregate telemetry for sessions/tick/input clamps/ticket redemption, bounded authority-side Unreal movement/identity-gate/collision/net-update activity counters, source-level authority displacement + rotation envelope counters, and reference-oracle parity/hardening evidence through v4.5 for canonical build/quest vectors, fixture integrity, payload-bound idempotency, Garage 17 inventory/blueprint rebuild semantics, bounded race-ID/checkpoint input validation, deterministic authoritative race-result hashing bound to durable identity/build evidence, entitlement-bounded garage capacity/starter-lineage deletion protection, canonical MQ001–MQ100 identifier/immediate-prerequisite sequencing, server-derived quest reward formula parity, a Python oracle quest-completion API that no longer accepts caller-controlled reward amounts, shared Python/Go checkpoint-order + finish-order race lifecycle acceptance vectors, race-start Roadworthy/active-build binding acceptance vectors, start/checkpoint/finish operation-id replay payload-binding vectors, deterministic MQ004/MQ005/MQ009/MQ012 inventory/blueprint/Roadworthy side-effect vectors, rebuild operation replay owner/vehicle/source-result revision/hash binding, and public quest-operation durable-key scoping to authoritative character identity. PostgreSQL observability v3.1-v3.4 remains partial: CI proves the least-privilege monitoring principal, PostgreSQL 17 `pg_stat_statements` preload/extension, bounded statement collector, statement-series cap, and exercised query-text/credential non-disclosure, but production secret-manager delivery/rotation, deployed scrape/dashboard/SLO evidence, long-retention cardinality behavior, deployment-scale load, and HA/failover evidence remain open. Runtime Service Restart Recovery v4.6 is CI evidence for Go service-process loss/restart only; it does not prove node rescheduling, PostgreSQL/Redis failover, production RPO/RTO, regional DR, or live Unreal recovery. The Unreal telemetry and integrity envelopes remain source/static-validation evidence only until a real UE 5.8 build/package run succeeds and retained live logs prove emission. Reference parity/hardening through v4.5 remains partial evidence and does not close the full reference-oracle parity gate. Threat Exercise v3.0 remains CI-exercised partial evidence and does not close deployed ingress, live Unreal transport, final-physics anti-cheat, privileged-admin audit, or production-security gates. The repository does not yet claim deployed ingress/header-sanitization evidence, deployed metrics/dashboard/SLO evidence, Redis HA/failover verification, live Unreal replication-byte/rate or race-validation metrics, deployment-scale or long-duration load/soak performance, packaged Unreal↔Go production E2E, final vehicle physics, live calibration of physics-derived displacement/rotation/acceleration envelopes, ranked anti-cheat, HA/DR, or production deployment.
