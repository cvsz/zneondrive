# zNeonDrive — PROJECT: NEON DRIVE

> **You don't own the road. You earn it.**

zNeonDrive is the design and implementation repository for **PROJECT: NEON DRIVE**, an 18+ persistent online open-world RPG centered on vehicle building, racing, adventure, factions, relationships, crews, pets, and a long-lived social world.

## Product direction

- **World:** NOVA CITY, 2097
- **Genre:** Online Open World / RPG / Vehicle Building / Racing / Adventure / Social
- **World model:** Persistent online world, server-authoritative by design
- **Player identity:** 1 account → 1 primary character → 1 starter vehicle
- **Vehicle philosophy:** a vehicle is a persistent identity object with ownership, builder, build, repair, race, and reputation history
- **VIP constraint:** garage/storage/convenience capacity only; no direct competitive performance advantage
- **Current phase:** Game Design Bible v0.1 + implementation contracts. Runtime technology is intentionally not locked yet.

## Canonical design

- [Original Story Bible / Master Plot](./zNeonDrive-concept.md)
- [Game Design Bible v0.1](./docs/game-design-bible-v0.1.md)
- [NOVA CITY World Bible](./docs/nova-city-world-bible.md)
- [Gameplay Systems](./docs/gameplay-systems.md)
- [Architecture](./docs/architecture.md)
- [Content Contracts](./docs/content-contracts.md)
- [Roadmap](./ROADMAP.md)
- [Implementation Checklist](./IMPLEMENTATION-CHECKLIST.md)

## Machine-readable design catalogs

Content under `design/catalog/` is intended to become the stable source of content IDs for future clients/services:

- 25 named launch characters
- 5 primary factions
- 10 launch districts
- 100 main-story quest records across 7 chapters
- side-content seeds
- vehicle part taxonomy
- functional pet archetypes

Schemas live under `design/schemas/`.

## Core loop

```text
Arrive in NOVA CITY
→ receive the broken prototype
→ rebuild the first vehicle
→ take jobs / discover districts / earn reputation
→ tune and reconfigure the same vehicle
→ build relationships and faction standing
→ enter underground and sanctioned competition
→ qualify for NOVA GRAND PRIX
→ expose PROJECT DRIVE ZERO
→ choose the future governance of mobility
→ continue into seasons, crews, world events, and expansions
```

## Fairness rules

1. Ranked outcomes are based on skill, earned/configured performance, event rules, and validated state.
2. VIP cannot add hidden horsepower, grip, durability, matchmaking priority, or race-stat multipliers.
3. Competitive vehicle builds must be server-validatable.
4. Reputation and relationship consequences are persistent gameplay state.
5. Endgame choices may alter a player's world state without breaking shared-world consistency.

## Status

The repository contains an implementation-ready **design foundation**. It does **not** claim that the full MMORPG runtime, networking, persistence, anti-cheat, matchmaking, live operations, HA/DR, or production deployment are implemented. Those remain evidence-gated roadmap phases.

## License

MIT. See [LICENSE](./LICENSE).
