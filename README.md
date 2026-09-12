# zNeonDrive — PROJECT: NEON DRIVE

> **You don't own the road. You earn it.**

zNeonDrive is the design and implementation repository for **PROJECT: NEON DRIVE**, an 18+ persistent online open-world RPG centered on vehicle building, racing, adventure, factions, relationships, crews, pets, and a long-lived social world.

## Product direction

- **World:** NOVA CITY, 2097
- **Genre:** Online Open World / RPG / Vehicle Building / Racing / Adventure / Social
- **World model:** persistent online world, server-authoritative by design
- **Player identity:** 1 account → 1 primary character → 1 starter vehicle
- **Vehicle philosophy:** a vehicle is a persistent identity object with ownership, builder, build, repair, race, and reputation history
- **VIP constraint:** garage/storage/convenience capacity only; no direct competitive performance advantage
- **Current phase:** Game Design Bible v0.2 + vertical-slice contract + executable authority proof + **client-presentation package v0.3**
- **Selected implementation direction:** Unreal Engine 5.8 client/dedicated gameplay server + Go service plane + PostgreSQL + Redis

## Present to a client now

Start with [client/README.md](./client/README.md).

Run the interactive offline demo:

```bash
python3 -m http.server 8080 --directory demo
```

Open `http://localhost:8080`.

The presentation package includes:
- Thai/English interactive browser demo,
- executive brief,
- Thai talk track,
- demo runbook,
- commercial scope options,
- prepared client Q&A,
- explicit evidence/status language.

The browser demo is intentionally presentation-only and makes no claim that the production MMORPG client is finished.

## Canonical design

- [Original Story Bible / Master Plot](./zNeonDrive-concept.md)
- [Game Design Bible v0.2](./docs/game-design-bible-v0.2.md)
- [Game Design Bible v0.1](./docs/game-design-bible-v0.1.md)
- [Vertical Slice v0.2](./docs/vertical-slice-v0.2.md)
- [NOVA CITY World Bible](./docs/nova-city-world-bible.md)
- [Gameplay Systems](./docs/gameplay-systems.md)
- [Architecture](./docs/architecture.md)
- [Content Contracts](./docs/content-contracts.md)
- [Roadmap](./ROADMAP.md)
- [Implementation Checklist](./IMPLEMENTATION-CHECKLIST.md)

## Production-direction ADRs

- [ADR-0001 — Server-authoritative state](./docs/adr/0001-server-authoritative-state.md)
- [ADR-0002 — Starter vehicle identity](./docs/adr/0002-starter-vehicle-identity.md)
- [ADR-0003 — Executable reference runtime](./docs/adr/0003-executable-reference-runtime.md)
- [ADR-0004 — Unreal Engine 5.8 playable client + dedicated gameplay server](./docs/adr/0004-unreal-engine-client-and-gameplay-server.md)
- [ADR-0005 — Go/PostgreSQL/Redis durable service plane](./docs/adr/0005-durable-service-plane.md)
- [ADR-0006 — Phase-gated deployment](./docs/adr/0006-phase-gated-deployment.md)

## Machine-readable design catalogs

Content under `design/catalog/` provides stable IDs and implementation contracts:

- 25 named launch characters
- 5 primary factions
- 10 launch districts
- 100 main-story quest records across 7 chapters
- 30 side-content seeds
- vehicle part taxonomy
- functional pet archetypes
- Garage 17 / Foundry 9 vertical-slice proof catalog

Schemas live under `design/schemas/`.

## Executable reference contracts

`src/zneondrive/` is a dependency-free Python reference model used to make core authoritative invariants executable before the production Unreal/Go runtime is implemented.

It proves:
- one-character ownership boundary,
- free vs VIP garage capacity without competitive stat advantage,
- append-only vehicle build revisions,
- idempotent quest rewards,
- race results bound to an accepted build revision and ordered checkpoints,
- starter vehicle deletion protection.

Run locally:

```bash
make ci
```

CI also validates the presentation demo is offline-first and free of unfinished placeholder markers.

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
2. VIP cannot add hidden horsepower, grip, durability, matchmaking priority, reward multipliers, or race-stat multipliers.
3. Competitive vehicle builds must be server-validatable and bind results to an exact accepted revision.
4. Rewardable mutations must be idempotent.
5. Reputation and relationship consequences are persistent gameplay state.
6. Endgame choices may alter a player's world state without breaking shared-world consistency.

## Status

**Ready now:** pre-production story/world/content contracts, vertical-slice specification, authority reference tests, implementation architecture direction, CI/security gates, and client presentation package.

**Not yet claimed complete:** playable Unreal MMORPG client, durable Go/PostgreSQL runtime, production networking, anti-cheat, matchmaking, live operations, HA/DR, platform certification, or production deployment.

Those remain evidence-gated roadmap phases.

## License

MIT. See [LICENSE](./LICENSE).
