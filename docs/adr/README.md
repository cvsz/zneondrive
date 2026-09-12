# Architecture Decision Records

ADRs capture durable architecture decisions and their consequences.

## Accepted decisions

- [ADR-0001 — Server-authoritative state](./0001-server-authoritative-state.md)
- [ADR-0002 — Starter vehicle identity](./0002-starter-vehicle-identity.md)
- [ADR-0003 — Executable reference runtime](./0003-executable-reference-runtime.md)
- [ADR-0004 — Unreal Engine client and dedicated gameplay server](./0004-unreal-engine-client-and-gameplay-server.md)
- [ADR-0005 — Durable service plane](./0005-durable-service-plane.md)
- [ADR-0006 — Phase-gated deployment](./0006-phase-gated-deployment.md)

## Creating a new ADR

Copy [0000-template.md](./0000-template.md), assign the next sequence number, and include:
- context/problem,
- decision,
- alternatives considered,
- consequences,
- security/operational impact,
- migration/rollback implications where relevant.

Accepted ADRs should not be silently rewritten to reverse a decision. Create a superseding ADR instead.
