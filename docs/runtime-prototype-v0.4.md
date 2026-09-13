# Runtime Prototype v0.4

## Scope

This milestone turns the accepted Unreal/Go architecture into executable repository structure without claiming alpha or production readiness.

## Version baseline

- Unreal Engine: 5.8 family; validate source builds against the current 5.8 hotfix used by the team.
- Go: 1.27.x.
- PostgreSQL: 17 for the local vertical-slice stack.
- Redis: 8 for local ephemeral coordination readiness.

## Gameplay plane

`game/` is an Unreal C++ project baseline containing:
- Game, Editor, and dedicated Server targets,
- a default GameMode,
- an authoritative prototype vehicle pawn,
- client driving input sent to a server RPC,
- movement executed only by authority and replicated to clients,
- default input mappings for throttle/steering.

This is a source/runtime skeleton. Garage 17 art, vehicle physics fidelity, World Partition, gameplay ability architecture, race routes, and vertical-slice assets are still pending.

## Service plane

`services/game-api/` is the first modular Go service binary.

Implemented:
- local device/resume-key bootstrap,
- opaque 24-hour session token issuance,
- session tokens and resume keys stored only as hashes,
- exactly one primary character on bootstrap,
- exactly one starter-lineage vehicle on bootstrap,
- PostgreSQL-backed durable character/economy/vehicle/quest state,
- immutable vehicle build revision rows,
- optimistic build revision conflict checks,
- idempotent operation IDs,
- ordered MQ001–MQ100 prerequisite enforcement,
- MQ012 marks the starter vehicle Roadworthy,
- versioned HTTP endpoints.

Prototype endpoints:

```text
GET  /healthz
POST /v1/sessions/bootstrap
GET  /v1/state
POST /v1/quests/{questID}/complete
POST /v1/vehicles/{vehicleID}/builds
```

The bootstrap resume key is a local prototype credential. It is not the final account/authentication model.

## Local stack

```bash
cp .env.example .env
docker compose up --build
curl http://127.0.0.1:18080/healthz
```

Host ports deliberately avoid the project's historically common 5432/6379 conflicts:
- PostgreSQL: 55432
- Redis: 56379
- Game API: 18080

## CI evidence

The Go runtime workflow must:
- compile on Go 1.27.x,
- run unit tests,
- run PostgreSQL integration tests,
- run `go vet`,
- build the service binary.

The Unreal source-build workflow is manual and targets a self-hosted Linux runner labelled `unreal-5.8`, because GitHub-hosted runners do not contain the licensed/source Unreal Engine toolchain.

## Explicit non-claims

v0.4 does not yet prove:
- Garage 17 is playable,
- production account authentication,
- Redis-backed presence/rate limiting,
- Unreal-to-Go authenticated service calls,
- inventory/blueprint persistence,
- relationships/factions,
- authoritative race instances,
- anti-cheat,
- load/soak behavior,
- backup/restore,
- HA/DR,
- production deployment.
