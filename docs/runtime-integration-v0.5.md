# Runtime Integration v0.5

## Goal

Connect the Phase 4 Unreal gameplay plane to the Go/PostgreSQL service plane without allowing a player client to self-assert durable vehicle identity.

## Trust flow

```text
Unreal client
  -> POST /v1/sessions/bootstrap
  -> session token + durable snapshot
  -> POST /v1/game-tickets
  -> one-time 60-second gameplay ticket

Unreal client
  -> Server RPC: gameplay ticket only

Unreal dedicated server
  -> POST /v1/internal/game-tickets/redeem
     X-Game-Server-Key: server-only secret
  -> Go atomically consumes ticket
  -> authoritative PostgreSQL snapshot returned once
  -> server binds VehicleID / build revision / part IDs / Roadworthy
     to the replicated vehicle pawn
```

The durable snapshot is never accepted from a client RPC.

## Client session subsystem

`UNDServiceSubsystem`:
- loads the prototype resume key from `Saved/NeonDrive/resume_key.txt`,
- bootstraps or resumes a durable account,
- keeps the 24-hour session token in memory only,
- refreshes `/v1/state`,
- performs quest/build mutations with new operation IDs,
- requests a short-lived gameplay ticket,
- passes only that one-time ticket to `ANDPlayerController`.

The resume key file is prototype local identity, not final platform authentication.

## Dedicated-server binding

`ANDPlayerController` accepts only a 64-character hex gameplay ticket from the owning client.

On authority it:
1. reads `ZNEON_GAME_SERVER_KEY` from the server environment,
2. calls the internal ticket redemption endpoint,
3. receives the durable snapshot directly from Go,
4. applies durable identity to the possessed `ANDVehiclePawn`.

`ANDVehiclePawn` replicates:
- durable vehicle ID,
- active build revision,
- active part IDs,
- Roadworthy status,
- durable-binding status.

Networked prototype movement is blocked until the server has bound a valid durable identity. Standalone editor play remains usable without service binding.

## Service-plane ticket rules

Gameplay tickets:
- are cryptographically random,
- are stored only as SHA-256 hashes,
- expire after 60 seconds,
- may be redeemed exactly once,
- are linked to the account behind a live session,
- require a server-only shared key at the internal redemption endpoint.

The API refuses to start unless `GAME_SERVER_SHARED_KEY` has at least 32 characters.

## Local run

Start the service plane:

```bash
cp .env.example .env
docker compose up --build
```

For a local Unreal dedicated server, use the same local secret as the API:

```bash
export ZNEON_GAME_SERVER_KEY="$(grep '^GAME_SERVER_SHARED_KEY=' .env | cut -d= -f2-)"
export ZNEON_GAME_API_INTERNAL_URL=http://127.0.0.1:18080
```

Do not place a production server key in client config, source code, packaged client assets, logs, or GitHub variables visible to untrusted builds.

## Automated evidence

The Go integration suite proves:
- HTTP bootstrap creates durable account/character/starter vehicle state,
- session issues a gameplay ticket,
- internal server key is required,
- ticket redemption succeeds once,
- second redemption is rejected,
- MQ001 persists,
- resume key reconnect returns the same account/vehicle/progression,
- a fresh ticket after reconnect redeems to the same durable state.

Static v0.5 validation checks the Unreal trust boundary source and migration ordering.

## Still open

v0.5 source/e2e evidence does **not** prove:
- Unreal 5.8 source compilation on the self-hosted runner,
- a live packaged client/dedicated-server network session,
- TLS/mTLS for production transport,
- final platform identity/authentication,
- inventory/blueprints,
- race instance/result authority,
- relationship/faction persistence,
- Garage 17 playable content,
- anti-cheat, load/soak, backup/restore, HA/DR.
