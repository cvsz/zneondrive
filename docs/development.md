# Development

## Current development mode

zNeonDrive has moved from design-only work into **Phase 4 runtime prototype** while retaining strict evidence gates.

Selected direction:
- Unreal Engine 5.8.x client + dedicated gameplay server,
- Go 1.27.x durable service plane,
- PostgreSQL transactional persistence,
- Redis reserved for ephemeral coordination/presence/rate-limit work.

## Fast validation

Requirements: Python 3.x and standard shell tooling.

```bash
make validate-design
make validate-runtime
make ci
```

## Go service plane

```bash
cd services/game-api
go test ./...
go vet ./...
go build ./cmd/server
```

PostgreSQL integration tests are run in GitHub Actions. For local integration testing:

```bash
cp .env.example .env
docker compose up -d postgres redis
cd services/game-api
TEST_DATABASE_URL='postgres://zneondrive:zneondrive-local@127.0.0.1:55432/zneondrive?sslmode=disable' \
  go test -tags=integration ./...
```

Run the whole local service stack:

```bash
docker compose up --build
curl http://127.0.0.1:18080/healthz
```

## Unreal runtime

See [game/README.md](../game/README.md).

The repository contains Game, Editor, and Server targets plus a minimal authoritative replicated vehicle pawn. Full Unreal compile evidence requires a UE 5.8 source toolchain and is deliberately kept on a manual self-hosted workflow.

## Content authoring rules

- Never rename a published stable content ID merely to change display text.
- New main quests must preserve explicit prerequisite semantics.
- Character, faction, district, part, pet, and quest references must resolve.
- VIP design changes require fairness review.
- Pet changes must not add ranked competitive buffs.
- Story changes that affect authority/persistence assumptions require an ADR or architecture update.
- Do not mark runtime behavior complete because a design record or source file exists.

## Runtime engineering rules

- Client RPC input is untrusted and clamped/validated by authority.
- Durable reward/build mutations require operation IDs.
- PostgreSQL is the durable source of truth for current Phase 4 account/character/vehicle/quest state.
- Do not move economy or permanent ownership authority into only an Unreal match process.
- Secrets/resume keys/session tokens must never be logged or committed.
- A CI workflow definition is not equivalent to a successful build artifact.

## Quality expectations

- Keep changes reviewable.
- Add tests when durable mutation behavior changes.
- Add validator assertions for new invariants.
- Prefer deterministic/reproducible tooling.
- Never commit credentials or player/user data.
- Never weaken CI/security gates to make a change pass.
