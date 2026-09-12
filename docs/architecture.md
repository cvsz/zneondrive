# Architecture — Technology-Neutral Online World Model

This document defines **domain authority and contracts**, not a mandated engine, programming language, database, cloud, or transport protocol.

## System context

```text
Player Client(s)
   │ intents / inputs / presentation acknowledgements
   ▼
Authoritative Online World
   ├─ Identity & Session
   ├─ Character & Relationship
   ├─ Vehicle & Build
   ├─ Inventory / Blueprint
   ├─ Quest & Narrative State
   ├─ World / District State
   ├─ Race / Event Authority
   ├─ Faction & Reputation
   ├─ Crew / Social
   ├─ Pet / Companion
   └─ Audit / Telemetry / Live Operations
```

## Authority rules

The client may predict and render, but cannot be trusted as the source of truth for:
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
Money, rewards, sinks, and item grants. Mutations should always be attributable to a reason/event.

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

## Synchronization classes

- **Immediate authoritative:** race state, ownership mutations, competitive validation.
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

Trusted only after validation:
- service-to-service messages,
- administration/live-ops actions,
- content deployments.

## Security requirements

- Never trust client-calculated rewards or race results.
- Rate-limit state mutations.
- Use replay protection/idempotency for rewardable actions.
- Audit privileged/live-ops changes.
- Separate public profile data from private account/security data.
- Minimize personally identifiable information.
- Treat telemetry as adversarial input.
- Require anti-cheat evidence before ranked production release.

## Availability and recovery evidence

Production readiness later requires:
- backup and restore of persistent player state,
- RPO/RTO,
- race/event degradation behavior,
- reconnect semantics,
- queue/backpressure policy,
- duplicate-reward prevention after retries,
- rollback of content deployments.

## Known constraint

The repository is currently design/contract ready, not a completed online runtime. Engine, protocol, database, and deployment decisions remain open and should be recorded as ADRs after prototypes provide evidence.
