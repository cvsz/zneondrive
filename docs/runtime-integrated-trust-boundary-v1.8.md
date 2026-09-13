# Runtime Integrated Trust-Boundary Evidence v1.8

## Scope

This increment exercises the existing service-plane security boundaries together against real PostgreSQL and Redis dependencies in Runtime Go CI. It does not change the selected Unreal Engine 5.8 + Go 1.27 + PostgreSQL + Redis architecture and does not move authority into Redis or the player client.

The integration test constructs the same security layers used by the service plane:

1. PostgreSQL-backed API authority,
2. security telemetry,
3. trusted-ingress identity resolution,
4. Redis-coordinated distributed rate limiting,
5. one-time gameplay-ticket issue/redeem flow,
6. server-only internal authorization.

## Evidence

`TestIntegratedTrustBoundaryAbusePaths` requires both `TEST_DATABASE_URL` and `TEST_REDIS_ADDR` and verifies that:

- an untrusted direct client cannot evade the bootstrap limiter by rotating `X-Forwarded-For`; spoofed values still consume one socket-peer Redis budget;
- a configured trusted loopback ingress may preserve distinct client identities, so independent clients receive independent abuse budgets;
- a forged bearer/session credential is rejected even after trusted-ingress and distributed-limiter processing;
- a gameplay ticket can be issued only from an authenticated session;
- the internal redemption endpoint rejects an incorrect game-server key;
- a valid server-only redemption returns the PostgreSQL-authoritative account/vehicle snapshot;
- replay of the same one-time gameplay ticket is rejected through the complete middleware stack.

The test runs inside the existing Runtime Go integration job using PostgreSQL 17 and Redis 8 service containers.

## Trust boundary

The test deliberately preserves the existing authority model:

- the client may request but cannot assert durable identity or race state;
- PostgreSQL remains authoritative for durable account, character, vehicle, progression and ticket state;
- Redis coordinates abuse-control budgets only;
- forwarding headers affect client identity only when the immediate peer is explicitly trusted;
- internal gameplay-server routes still require the server-only shared key;
- gameplay tickets remain hashed, short-lived and single-use.

## Explicit non-claims

This is **CI-scale integrated service-plane security evidence**, not production ingress or live Unreal evidence.

It does not prove:

- deployed proxy/header sanitization or direct-network bypass prevention;
- live packaged Unreal client or dedicated-server transport;
- physics-derived impossible-state detection or ranked anti-cheat;
- deployment-scale load or long-duration soak;
- deployed dashboard/SLO/alert correctness;
- HA/failover, regional DR, production backup custody/PITR, or production readiness.

Those gates remain open until their own executable/deployed evidence exists.
