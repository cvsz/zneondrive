# About zNeonDrive

**PROJECT: NEON DRIVE** is an 18+ persistent online open-world RPG set in NOVA CITY in 2097.

The project combines:
- persistent character and vehicle identity,
- deep vehicle building and build history,
- street and sanctioned racing,
- RPG quests and branching relationships,
- faction reputation,
- functional pets,
- crews and social spaces,
- seasonal championships and world events.

## Product position

The game's central rule is that the player's vehicle is a second identity, not a disposable mount. The starter vehicle begins as a broken machine at Garage 17 and can evolve into a server-famous Legendary vehicle while preserving provenance, build revisions and race history.

## Current repository phase

The repository is in **Phase 4.3 — Authoritative Race runtime v0.7**.

Implemented evidence includes:
- Game Design Bible/content catalogs and validators,
- Python authority reference oracle,
- Unreal Engine 5.8 C++ gameplay/source integration baseline,
- Go 1.27 durable service plane,
- PostgreSQL account/session/character/vehicle/build/quest/inventory/blueprint persistence,
- one-time gameplay-ticket trust boundary,
- idempotent quest and rebuild mutations,
- server-only authoritative race instance/checkpoint/result lifecycle,
- exact active-build binding for durable race results,
- PostgreSQL unit/integration evidence,
- local Docker Compose stack,
- client presentation package,
- comprehensive product/engineering/operations documentation baseline.

Still evidence-gated:
- successful retained UE 5.8 source-build artifact,
- packaged live Unreal↔Go E2E,
- live Unreal race lifecycle E2E,
- playable Garage 17/Foundry 9 content,
- final vehicle physics,
- relationships/factions,
- physics-derived anti-cheat/impossible-state detection,
- load/soak,
- backup/restore and DR verification,
- production deployment.

See [docs/README.md](./docs/README.md) for the complete documentation map.

## Project motto

> You don't own the road. You earn it.
