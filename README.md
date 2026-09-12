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
- **Current phase:** **Phase 4 Runtime Prototype v0.4**
- **Selected implementation direction:** Unreal Engine 5.8 client/dedicated gameplay server + Go 1.27 service plane + PostgreSQL + Redis

## Runtime prototype v0.4

The repository now contains executable implementation scaffolding for both planes.

### Gameplay plane — Unreal

`game/` contains:
- UE 5.8 C++ project,
- Game / Editor / dedicated Server targets,
- default GameMode,
- prototype vehicle pawn,
- client input → server RPC → authoritative movement,
- replicated movement back to clients,
- source-build workflow baseline for a self-hosted UE 5.8 Linux runner.

See [game/README.md](./game/README.md).

### Durable service plane — Go

`services/game-api/` contains:
- Go 1.27 modular service binary,
- PostgreSQL account/session/character/vehicle/build/quest persistence,
- hashed resume/session credentials,
- one-character + starter-vehicle bootstrap,
- immutable vehicle build revisions,
- optimistic revision checks,
- idempotent durable mutation operation IDs,
- MQ001–MQ100 sequential quest gate,
- MQ012 Roadworthy transition,
- unit and PostgreSQL integration tests.

Local stack:

```bash
cp .env.example .env
docker compose up --build
curl http://127.0.0.1:18080/healthz
```

Host defaults intentionally avoid common 5432/6379 conflicts:
- PostgreSQL: `55432`
- Redis: `56379`
- Game API: `18080`

See [Runtime Prototype v0.4](./docs/runtime-prototype-v0.4.md).

## Present to a client now

Start with [client/README.md](./client/README.md).

Run the interactive offline presentation demo:

```bash
python3 -m http.server 8080 --directory demo
```

Open `http://localhost:8080`.

The browser demo is presentation-only; the Unreal/Go code is the implementation baseline. Neither is a claim that the full MMORPG is production-ready.

## Canonical design

- [Original Story Bible / Master Plot](./zNeonDrive-concept.md)
- [Game Design Bible v0.2](./docs/game-design-bible-v0.2.md)
- [Game Design Bible v0.1](./docs/game-design-bible-v0.1.md)
- [Vertical Slice v0.2](./docs/vertical-slice-v0.2.md)
- [Runtime Prototype v0.4](./docs/runtime-prototype-v0.4.md)
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

## Evidence layers

### Python reference oracle
`src/zneondrive/` remains a dependency-free contract oracle for core authority invariants.

### Go durable runtime
`services/game-api/` implements the first persistent service-plane slice.

### Unreal runtime
`game/` contains the first gameplay-plane source baseline. A successful self-hosted UE source build is still required before claiming Unreal build evidence.

Run repository checks:

```bash
make ci
```

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

**Implemented now:** pre-production design/content contracts, vertical-slice specification, client presentation package, Python authority oracle, Unreal C++ source baseline, Go durable service-plane prototype, PostgreSQL schema/integration tests, local Compose stack, and CI validation.

**Still evidence-gated:** successful UE source-build artifact, Garage 17 playable content, final vehicle physics, Unreal↔Go authenticated integration, inventory/blueprints, relationships/factions, authoritative race instances, anti-cheat, matchmaking, live operations, HA/DR, platform certification, and production deployment.

## License

MIT. See [LICENSE](./LICENSE).
