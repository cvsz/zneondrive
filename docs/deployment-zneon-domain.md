# zneon.zeaz.dev deployment

This document describes the public marketing/web surface for **PROJECT: NEON DRIVE** by **ZEAZDEV COMPANY LIMITED**.

## Public architecture

`https://zneon.zeaz.dev` should terminate at Cloudflare and forward through a Cloudflare Tunnel to the local web container on `127.0.0.1:18081`.

The public web container serves the static frontend and exposes only `/api/healthz` to the Go service. Player/session/quest/build/race mutation APIs are intentionally **not** reverse-proxied through the marketing site.

Runtime topology:

```text
Internet
  -> Cloudflare TLS/WAF
  -> Cloudflare Tunnel
  -> 127.0.0.1:18081
  -> nginx web container
      -> static frontend
      -> /api/healthz -> game-api:8080/healthz

private compose network
  game-api -> PostgreSQL 17
           -> Redis 8
```

## Start locally

```bash
cp .env.example .env
make env-init

docker compose up -d --build

docker compose ps
curl -fsS http://127.0.0.1:18081/healthz
curl -fsS http://127.0.0.1:18081/api/healthz
```

The frontend is then available at `http://127.0.0.1:18081`.

## Cloudflare Tunnel

If a tunnel already exists for `zeaz.dev`, add this ingress rule before the catch-all rule:

```yaml
ingress:
  - hostname: zneon.zeaz.dev
    service: http://127.0.0.1:18081
  - service: http_status:404
```

Then route DNS to the tunnel and restart/reload `cloudflared` using the installation method already used on the host.

Typical verification:

```bash
curl -I https://zneon.zeaz.dev/
curl -fsS https://zneon.zeaz.dev/healthz
curl -fsS https://zneon.zeaz.dev/api/healthz
```

## Security boundary

- The public site is read-only marketing content plus a health probe.
- `GAME_SERVER_SHARED_KEY` remains server-only and is never sent to or embedded in the frontend.
- PostgreSQL remains durable authority.
- Unreal dedicated servers remain gameplay authority.
- Redis remains ephemeral coordination.
- Metrics stay bound to loopback by default and are not exposed by the web proxy.
- The site does not mark the game production-ready. Live Unreal↔Go E2E, load/soak, deployed SLO, HA/DR and release evidence remain separate gates.

## Website content

The site presents:

- PROJECT: NEON DRIVE / NOVA CITY 2097 hero experience;
- game scope and launch-content counts;
- Unreal Engine 5.8 + Go 1.27 + PostgreSQL 17 + Redis 8 architecture;
- ZEAZDEV COMPANY LIMITED studio profile;
- live service-plane health state;
- evidence-gated development status;
- links to the project repository and corporate site.

## Operations

Recommended routine checks:

```bash
docker compose ps
docker compose logs --tail=200 web game-api
curl -fsS http://127.0.0.1:18081/healthz
curl -fsS http://127.0.0.1:18081/api/healthz
```

For updates:

```bash
git pull --ff-only
docker compose up -d --build
```

Do not expose PostgreSQL, Redis, metrics, or the dedicated-server shared key through the public domain.
