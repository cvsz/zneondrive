# ADR-0003: Executable Reference Runtime Contracts

- **Status:** Accepted
- **Date:** 2026-09-12

## Context

The project needs executable evidence for server-authoritative design invariants before selecting a production client engine, server language, database, protocol, or hosting topology.

## Decision

Maintain a small dependency-free Python reference model under `src/zneondrive/` with unit tests. It defines behavior, not the production stack.

The reference must cover ownership, garage capacity, immutable build revisions, idempotent reward mutations, race/build binding, and starter-prototype protection.

## Consequences

- Design rules become testable early.
- Future implementations can port the behavior and run equivalent contract tests.
- Python is **not** selected as the production server by this ADR.
- Networking, persistence storage, concurrency control, anti-cheat, HA/DR, and deployment remain separate evidence-gated decisions.
