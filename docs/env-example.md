# Local environment template

The repository root `.env.example` is the canonical local-development template for PostgreSQL 17, Redis 8, the Go game API, internal metrics, Unreal Engine 5.8 build/package tooling, and the local dedicated-server trust boundary.

## Usage

```bash
cp .env.example .env
make env-init
```

`make env-init` replaces the known local `GAME_SERVER_SHARED_KEY` placeholder with a random value. Keep `.env` untracked and never reuse local example credentials in staging or production.

For Unreal dedicated-server runtime, `ZNEON_GAME_SERVER_KEY` must match the Go service-plane `GAME_SERVER_SHARED_KEY`, but it must be injected only into the server runtime and must never be packaged into the player client.

The template includes optional host-side `DATABASE_URL`, `TEST_DATABASE_URL`, `REDIS_ADDR`, and `TEST_REDIS_ADDR` values for direct local tooling/tests. Docker Compose continues to construct its own service-network database/Redis addresses internally.

## Evidence boundary

This template is configuration documentation only. It does not constitute production secret management, deployed network isolation, successful Unreal packaging, live Unreal↔Go evidence, HA/DR evidence, or production readiness.
