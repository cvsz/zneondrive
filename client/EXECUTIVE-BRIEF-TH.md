# PROJECT: NEON DRIVE — Executive Brief for Client Meeting

**Purpose:** client-facing summary for discussion and scope alignment  
**Prepared:** 12 September 2026  
**Status:** pre-production design + executable authority proof + interactive concept demo

## One-line pitch

**PROJECT: NEON DRIVE** is an 18+ persistent online open-world RPG where the player's car is not a disposable mount — it is a second identity that is built, raced, repaired, remembered, and eventually becomes the key to a city-scale conflict over who controls mobility.

## What is already complete

- canonical Story Bible and Game Design Bible v0.2,
- NOVA CITY world structure,
- 25 launch characters,
- 5 factions,
- 10 launch districts,
- 100 main quests across 7 chapters,
- 30 side-content seeds,
- vehicle/part taxonomy,
- companion/pet rules,
- VIP fairness contract,
- Garage 17 / Foundry 9 MQ001–MQ012 vertical-slice specification,
- technology-neutral server authority rules,
- executable Python reference model proving core ownership/reward/race invariants,
- automated CI validation,
- recommended production architecture ADRs,
- client-facing offline browser demo.

## Core differentiator

Most racing games treat cars as replaceable inventory. NEON DRIVE treats a vehicle like an RPG character:

- stable Vehicle ID,
- ownership history,
- builder/build revision history,
- repair and race history,
- vehicle reputation,
- story provenance.

The first broken car can remain the player's iconic Legendary vehicle months later.

## Customer value proposition

The project combines four engagement loops in one persistent world:

1. **RPG progression** — story, factions, relationships, choices;
2. **vehicle mastery** — building, tuning, provenance, role-specific configurations;
3. **competition** — street, circuit, drag, drift, off-road, team and championship;
4. **social retention** — crews, garages, car meets, venues, seasonal events.

## Monetization principle

VIP is designed around **capacity and convenience**, not Pay-to-Win.

VIP can add:
- garage slots,
- blueprint/storage capacity,
- saved configurations,
- cosmetic/social customization.

VIP does not add hidden horsepower, grip, durability, ranked reward multipliers, matchmaking priority, or required progression access.

## Story hook

The player's first broken vehicle is secretly a surviving prototype from **PROJECT DRIVE ZERO**. Its ECU contains a missing artifact sought by VANTEX Corporation. The machine the player spends the game rebuilding becomes evidence, identity, and eventually the key to deciding the future governance of NOVA CITY.

## What the demo proves tomorrow

The presentation demo is intentionally a **concept + systems demonstrator**, not a fake claim of a finished MMORPG.

It demonstrates:
- world/product positioning,
- 7-chapter story spine,
- 10-district city plan,
- character/faction structure,
- Garage 17 → First Ignition vertical-slice flow,
- interactive vehicle build roles,
- VIP fairness,
- server-authoritative architecture,
- delivery roadmap and readiness evidence.

## Recommended next commercial engagement

### Package A — Playable Vertical Slice
Deliver a polished Garage 17 / Foundry 9 client:
- player movement/interactions,
- starter vehicle rebuild,
- MQ001–MQ012,
- First Ignition cinematic/gameplay moment,
- one legal time trial,
- one underground race,
- reconnect/save proof,
- production-style telemetry.

Target planning window: **8–12 weeks** after art/fidelity/platform scope is locked.

### Package B — Online Alpha
Extend the slice into:
- real account/session backend,
- durable PostgreSQL persistence,
- multi-player district instance,
- event registration,
- crews/presence MVP,
- faction/relationship state,
- anti-cheat telemetry,
- load/soak and restore evidence.

Target planning window: **4–6 months** depending team size and asset scope.

### Package C — Live-Service Production Program
Full NOVA CITY launch program:
- 100 main quests implemented,
- launch districts,
- seasonal championship operations,
- moderation/support,
- live ops,
- security hardening,
- observability,
- DR,
- platform certification,
- launch operations.

This is a multi-phase production program, not an overnight or single-sprint deliverable.

## Decision requested from the client

The meeting should finish with four decisions:

1. target platforms: PC only first, or PC + console;
2. visual fidelity target: stylized AA, realistic AA/AAA, or staged progression;
3. online concurrency target for first commercial milestone;
4. preferred engagement: Vertical Slice, Online Alpha, or Full Production Program.
