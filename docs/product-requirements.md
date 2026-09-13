# Product Requirements — PROJECT: NEON DRIVE

## Product statement

PROJECT: NEON DRIVE is an 18+ persistent online open-world RPG set in NOVA CITY, 2097. The defining product rule is that the player's vehicle is a long-lived identity object rather than disposable inventory.

> You don't own the road. You earn it.

## Primary player fantasy

A new player arrives with little status, receives a broken prototype at Garage 17, rebuilds it, earns reputation through jobs and racing, forms relationships and faction ties, discovers PROJECT DRIVE ZERO, and eventually influences the governance of mobility in NOVA CITY.

## Core pillars

1. **Vehicle as identity** — stable vehicle identity, provenance, build history and race history.
2. **Earned mastery** — progression through play, skill, knowledge and world participation.
3. **Persistent consequences** — quests, factions, relationships and ownership survive sessions.
4. **Fair competition** — VIP may add capacity/convenience but not direct ranked performance.
5. **Shared world** — social venues, crews, events and seasonal championship structure.
6. **Server authority** — valuable, durable and competitive outcomes are not trusted to clients.

## Launch design scope

Canonical design currently defines:
- 25 named characters,
- 5 primary factions,
- 10 launch districts,
- 100 main quests across 7 chapters,
- 30 side-content seeds,
- vehicle part/build taxonomy,
- pets/companions,
- crew/social systems,
- five endgame governance paths.

## First commercial milestone

The recommended milestone remains a playable vertical slice:
- Garage 17 + Foundry 9,
- MQ001–MQ012,
- starter rebuild,
- First Ignition,
- one legal race proof,
- one underground race proof,
- save/reconnect,
- authoritative durable identity,
- telemetry and acceptance evidence.

## Non-functional requirements

### Security
- client is untrusted,
- reward/economy mutations are idempotent,
- session/game-server trust is explicit,
- secrets are never embedded in client builds.

### Reliability
- reconnect cannot duplicate rewards,
- durable player ownership survives gameplay-server failure,
- schema/content compatibility is migration-aware.

### Accessibility
At minimum the vertical slice targets remappable controls, readable UI scaling, subtitle support, color-independent critical signals and reduced-motion consideration.

### Operability
Production candidates require metrics, logs, traces, audit events, backup/restore evidence, rollback procedures and incident response.

## Explicit non-goals for the current phase

Not currently claimed:
- complete MMORPG content implementation,
- verified production concurrency,
- final anti-cheat,
- verified HA/DR,
- console certification,
- 24/7 live operations.

Those remain roadmap gates.
