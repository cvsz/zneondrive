# zneon.zeaz.dev deployment

This document defines the public web and versioned client-release surface for **PROJECT: NEON DRIVE** by **ZEAZDEV COMPANY LIMITED**.

## Public architecture

`https://zneon.zeaz.dev` terminates at Cloudflare and forwards through the shared Cloudflare Tunnel to the dedicated zNeonDrive loopback origin `127.0.0.1:18085`.

Port `18085` is reserved for zNeonDrive so it does not collide with other ZEAZDEV services. The matching Terraform/Cloudflare configuration is maintained in `cvsz/zworkforce/infrastructure/terraform/cloudflare`.

The public nginx container serves static content and exposes only `/api/healthz` to the Go service. Player/session/quest/build/race mutation APIs are not reverse-proxied through this website.

Runtime topology:

```text
Internet
  -> Cloudflare
  -> Cloudflare Tunnel
  -> 127.0.0.1:18085
  -> nginx web container
      -> project website
      -> /client-v-0-0-1/
      -> /api/healthz -> game-api:8080/healthz

private compose network
  game-api -> PostgreSQL 17
           -> Redis 8
```

## Local start

Create `.env` from `.env.example`, initialize the local server credential with the repository-supported `make env-init`, then start the Compose stack. `WEB_PORT=18085` is the reviewed default.

Expected local checks:

- `http://127.0.0.1:18085/healthz`
- `http://127.0.0.1:18085/api/healthz`
- `http://127.0.0.1:18085/client-v-0-0-1/`
- `http://127.0.0.1:18085/client-v-0-0-1/manifest.json`

## Cloudflare route

The tunnel route is:

```yaml
- hostname: zneon.zeaz.dev
  service: http://127.0.0.1:18085
```

It must appear before the terminal `http_status:404` rule. Terraform keeps full tunnel configuration management opt-in until the live ingress has been reconciled and reviewed.

## Public release URLs

Release metadata is safe to publish before platform binaries exist:

- `https://zneon.zeaz.dev/client-v-0-0-1/`
- `https://zneon.zeaz.dev/client-v-0-0-1/manifest.json`
- `https://zneon.zeaz.dev/client-v-0-0-1/SHA256SUMS.txt`

The following files remain absent until real platform builds succeed:

- `zneondrive-client-v-0-0-1.zip`
- `zneondrive-client-v-0-0-1.apk`
- `zneondrive-client-v-0-0-1.exe`
- `zneondrive-client-v-0-0-1.bat`

The manifest must keep each unavailable artifact in `pending` state with a null checksum. Never create placeholder binaries or rename unrelated files to a release extension.

## Security boundary

- PostgreSQL remains durable authority.
- Unreal dedicated servers remain gameplay authority.
- Redis remains ephemeral coordination.
- Internal metrics stay on loopback.
- Server-only credentials are not embedded in frontend or downloadable client files.
- The release directory inherits the same CSP, frame, referrer and content-type protections as the website.
- Publishing the website or release metadata does not claim that the game is production-ready.
