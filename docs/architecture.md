# Architecture — PROJECT: NEON DRIVE Online World

The authority model is technology-independent; the v0.3 implementation profile selects technologies without weakening those domain boundaries.

## System context

```text
Unreal Client(s)
   │ input / intent / prediction
   ▼
Unreal Dedicated Gameplay Servers
   │ validated gameplay outcomes / durable mutation requests
   ▼
Go Service Plane
   ├─ Identity & Session
   ├─ Character & Relationship
   ├─ Vehicle & Build Metadata
   ├─ Inventory / Blueprint
   ├─ Quest & Narrative State
   ├─ Faction & Reputation
   ├─ Crew / Social
   ├─ Event Coordination
   ├─ Entitlement
   └─ Audit / Live Operations
   │
   ├─ PostgreSQL — durable transactional state
   ├─ Redis — ephemeral coordination/cache/rate-limit state
   └─ NATS JetStream — only for asynchronous durable workflows that justify it
```

## Authority rules

The player client may predict and render, but cannot be trusted as the source of truth for:
- ownership,
- inventory quantities,
- vehicle part/build legality,
- money/XP/reputation,
- quest completion,
- faction standing,
- race start/finish/results,
- ranked timing,
- crew roles,
- VIP entitlements.

The dedicated gameplay server owns moment-to-moment session authority. The durable service plane owns long-lived value and progression.

## Selected implementation profile

### Client
**Unreal Engine 5.8**, C++ plus Blueprint.

C++ is preferred for reusable core gameplay, network-sensitive logic, authority boundaries, and performance-critical systems. Blueprint remains appropriate for designer-authored presentation and rapid iteration where state cannot violate server authority.

See ADR-0004.

### Gameplay plane
**Unreal dedicated servers** own:
- world/race instance simulation,
- vehicle session state,
- authoritative checkpoints/laps,
- event-local rules,
- gameplay validation before durable outcome submission.

A gameplay server is not the sole database of player ownership/economy.

### Durable service plane
**Go** services own long-lived mutations and domain APIs.

The first alpha may deploy as one modular Go binary with strict package/domain boundaries rather than premature microservices.

See ADR-0005.

### Persistence
**PostgreSQL** is the durable transactional source of truth for account/profile, vehicle ownership/build metadata, progression, quests, economy, relationships/factions, crews, entitlements, and auditable operations.

**Redis** is used for ephemeral state where loss can be recovered or rebuilt: cache, presence, short-lived coordination, rate-limit counters, and similar workloads.

### Asynchronous events
Use **NATS JetStream** only for concrete asynchronous workflows that need durable delivery. Do not add an event bus merely to appear “enterprise”.

### Deployment
Vertical slice starts simple. Containers and repeatable manifests are required; k3s/Kubernetes enters after measured operational need.

See ADR-0006.

## Domain boundaries

### Identity & Session
Account identity, character selection, session lifecycle, restrictions, and entitlement lookup.

### Character & Relationship
Primary character profile, NPC relationship states, story flags, and personal-world overlays.

### Vehicle & Build
Vehicle identity, garage slots, components, build revisions, condition, provenance, and ownership history.

### Inventory / Blueprint
Parts, materials, blueprint knowledge, storage capacity, grants, consumption, and transfer rules.

### Quest & Narrative
Quest definitions are content. Quest instances/state are runtime data. Completion and rewards must be idempotent.

### Race & Event
Registration, eligibility, accepted build revision, ruleset, authoritative checkpoints, penalties, results, and awards.

### World
District availability, global season state, public event state, and shared environmental state.

### Economy
Money, rewards, sinks, and item grants. Mutations are attributable to a reason/event and idempotency key.

### Faction & Reputation
Independent faction standing, public reputation, thresholds, consequences, and repair paths.

### Crew & Social
Crew membership, roles, crew reputation, shared progression, and social presence.

## Stable identifiers

Runtime systems consume IDs from `design/catalog/`. Display text may change; IDs should not.

Examples:
- `char_maya_voss`
- `faction_vantex`
- `district_foundry_9`
- `MQ001`
- `part_engine_ice_street_i`

## Persistence principles

- Reward/economy mutations are idempotent.
- Vehicle build revisions are immutable snapshots or reconstructable from append-only history.
- Race results point to a specific accepted build revision.
- Story decisions are auditable.
- Destructive operations require explicit ownership/authorization checks.
- Unique story vehicle provenance cannot be silently rewritten.
- Durable mutation requests carry actor, reason/context, and idempotency metadata.

## Synchronization classes

- **Immediate authoritative:** race state, ownership-sensitive gameplay validation, competitive session state.
- **Near-real-time:** crew/social presence, public events.
- **Transactional:** inventory, economy, build changes, garage operations, quest rewards.
- **Eventually consistent acceptable:** leaderboards, analytics, non-critical discovery feeds.

## Trust boundaries

Untrusted:
- player client,
- modified client,
- local files,
- client clocks,
- client-provided prices/rewards/results.

Validated trust only:
- gameplay-server outcome proposals,
- service-to-service messages,
- administration/live-ops actions,
- content deployments.

## Security requirements

- Never trust client-calculated rewards or race results.
- Authenticate gameplay servers to the service plane.
- Rate-limit state mutations.
- Use replay protection/idempotency for rewardable actions.
- Audit privileged/live-ops changes.
- Separate public profile data from private account/security data.
- Minimize personally identifiable information.
- Treat telemetry as adversarial input.
- Require anti-cheat evidence before ranked production release.

## Observability

Production candidates require:
- OpenTelemetry-compatible trace/log/metric correlation,
- gameplay server session and health metrics,
- request/error/latency metrics,
- race validation and impossible-state rejection metrics,
- economy/idempotency rejection metrics,
- reconnect outcome metrics,
- database and cache saturation metrics,
- structured privileged audit logs.

## Availability and recovery evidence

Production readiness requires:
- backup and restore of persistent player state,
- explicit RPO/RTO,
- race/event degradation behavior,
- reconnect semantics,
- queue/backpressure policy,
- duplicate-reward prevention after retry/failure,
- rollback of service and content deployments,
- tested recovery rather than architecture diagrams alone.

## Current constraint

The architecture direction is selected, but production implementation evidence is not yet complete. Unreal project code, Go service-plane persistence, real transport, anti-cheat, load/soak, HA/DR, and production deployment remain open roadmap work.
