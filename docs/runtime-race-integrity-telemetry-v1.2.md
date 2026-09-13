# Runtime Race Integrity Telemetry v1.2

**Status:** implemented source/CI evidence; physics-derived anti-cheat and deployed observability remain open.

## Purpose

This increment turns authoritative race rejections into credential-safe security telemetry without changing the trust model. Unreal dedicated servers remain gameplay-authoritative, PostgreSQL remains durable authority, and Redis remains ephemeral rate-limit coordination.

## Implemented signals

The Go HTTP service now emits structured log events for:

- `security_event=race_integrity_rejected` when authoritative checkpoint/finish requests are rejected with HTTP 400 or 409,
- `security_event=game_server_auth_rejected` when an internal game-server route returns HTTP 401.

Race-instance IDs and socket-peer addresses are SHA-256-derived correlation buckets before logging. The game-server shared key and raw identifiers are never included in these events. Route labels are bounded static scopes instead of raw request paths.

The telemetry middleware observes the result of the existing API handler. It does not override, weaken, or independently decide authoritative race acceptance.

## Test evidence

Unit coverage verifies:

- deterministic correlation buckets do not expose raw identifiers,
- dynamic race IDs are reduced to static route scopes,
- authoritative race conflicts preserve the downstream HTTP response while emitting `race_integrity_rejected`,
- unauthorized dedicated-server calls emit `game_server_auth_rejected`,
- raw race IDs and peer addresses are absent from captured logs.

Repository CI, Runtime Go unit/integration/vet/build, CodeQL, and Dependency Review remain merge gates.

## Explicit non-claims

This does **not** prove or implement:

- vehicle-physics envelope validation,
- teleport/speed/acceleration impossible-state detection,
- client-memory or process anti-tamper,
- automatic bans or sanctions,
- live Unreal race telemetry transport,
- deployed log aggregation/alerting/SLO compliance,
- production anti-cheat readiness.

The next anti-cheat step should consume authoritative Unreal dedicated-server physics/race samples and evaluate conservative impossible-state rules using measured vehicle envelopes. Until live Unreal evidence exists, this telemetry remains a rejection/correlation baseline only.
