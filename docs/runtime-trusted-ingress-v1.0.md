# Runtime Trusted Ingress Identity v1.0

## Scope

This increment hardens the HTTP abuse-control identity boundary when the Go service plane is deployed behind one or more reverse proxies or ingress hops. It does not change gameplay authority: Unreal dedicated servers remain authoritative for gameplay decisions, PostgreSQL remains authoritative for durable state, and Redis remains ephemeral rate-limit coordination.

## Trust model

`RemoteAddr` is the default client identity source. Forwarded headers are ignored unless the immediate socket peer belongs to the explicit `TRUSTED_PROXY_CIDRS` allowlist.

When the immediate peer is trusted, the runtime evaluates `X-Forwarded-For` from right to left, skipping configured trusted proxy hops. The first untrusted address is treated as the effective client IP. This prevents an arbitrary left-most value supplied by a direct client from becoming authoritative when a trusted edge appends the real source address.

If `X-Forwarded-For` is malformed, the resolver fails closed to the immediate socket peer and emits `security_event=forwarded_identity_rejected` without logging the supplied header value.

The zero-value/default policy trusts no proxy. Invalid `TRUSTED_PROXY_CIDRS` configuration fails service startup rather than silently widening trust.

## Configuration

Example for an ingress network and one exact proxy address:

```bash
TRUSTED_PROXY_CIDRS=10.42.0.0/16,192.0.2.44
```

Only configure networks that are actually controlled proxy/ingress hops. Do not use `0.0.0.0/0` or `::/0`; that would make direct client-supplied forwarding headers trusted input and defeat the boundary.

The ingress itself must append or replace `X-Forwarded-For` correctly and must not allow untrusted traffic to reach the API directly on a path that is considered trusted.

## Implemented evidence

Unit coverage proves:

- invalid proxy CIDR configuration is rejected;
- CIDR and exact-IP allowlist entries parse correctly;
- an untrusted direct peer cannot spoof `X-Forwarded-For`;
- a trusted proxy can preserve distinct client identities for rate-limit buckets;
- multi-hop chains are evaluated from the trusted edge toward the client;
- malformed forwarding chains fall back to the socket peer;
- the existing credential-safe hashed bucket keys and Redis coordination remain unchanged.

## Non-claims / remaining gates

This source-level and unit-test evidence does **not** prove:

- a deployed ingress is configured to sanitize/append forwarding headers correctly;
- the API is network-isolated so only trusted ingress hops can reach the trusted listener path;
- rate limiting meets targets under real multi-replica HTTP load or soak;
- live Unreal↔Go transport, race anti-cheat, observability/SLO, HA/DR, backup/restore, or production deployment readiness.

The next production-readiness gate is multi-replica HTTP load/soak evidence using the trusted-ingress identity policy, followed by ranked-race impossible-state / anti-cheat telemetry.
