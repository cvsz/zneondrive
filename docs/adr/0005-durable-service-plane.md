# ADR-0005: Split Gameplay Runtime and Durable Service Plane

- Status: Accepted for vertical-slice and alpha implementation
- Date: 2026-09-12

## Context

A persistent open-world racing RPG has two different workloads:

1. low-latency session authority for movement, vehicles, races, checkpoints, and world-instance simulation;
2. durable transactional authority for accounts, ownership, quests, economy, builds, factions, crews, entitlements, audit, and live operations.

Keeping both inside one monolith would make scaling, recovery, testing, and operational ownership harder.

## Decision

Adopt a split architecture.

### Gameplay plane
- Unreal Engine dedicated servers
- authoritative session/race/world-instance state
- short-lived session state
- reconnect handoff and event outcome proposals

### Service plane
- **Go** services for identity/session coordination, profile, vehicle ownership/build metadata, inventory/economy, quest progression, relationship/faction state, crews, entitlements, matchmaking/event coordination, and live-ops APIs
- **PostgreSQL** as the durable transactional store
- **Redis** for ephemeral cache, presence, rate-limit counters, short-lived coordination, and idempotency acceleration where appropriate
- **NATS JetStream** only where asynchronous durable events are justified by an implemented workflow; do not introduce it merely for architecture aesthetics
- S3-compatible object storage for large non-transactional artifacts where needed

## API stance

- external/player-facing service APIs use explicit versioned HTTP/JSON or binary game-service contracts;
- internal latency-sensitive service calls may use gRPC/Connect-style contracts;
- every rewardable/durable mutation carries an idempotency key and actor/context metadata;
- live-ops mutations are auditable.

## Data ownership

Each domain owns its mutation rules even if the first alpha deploys as a modular monolith. Physical microservice separation is a scaling decision, not a design requirement.

## Consequences

The first implementation can ship as:
- one Go service binary with strict internal modules,
- one PostgreSQL database,
- one Redis instance,
- one or more Unreal dedicated servers.

This keeps the alpha operable while preserving boundaries for later scale-out.
