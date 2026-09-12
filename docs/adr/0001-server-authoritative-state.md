# ADR-0001: Server-Authoritative Persistent and Competitive State

- Status: Accepted
- Date: 2026-09-12

## Context

PROJECT: NEON DRIVE depends on persistent ownership, vehicle provenance, economy, quests, crews, faction/relationship state, and ranked racing. Trusting client-calculated state would make duplication, impossible builds, fabricated rewards, and result manipulation structurally easy.

## Decision

Ownership, inventory/economy mutations, quest completion, faction/relationship state, vehicle build legality, and competitive race results are authoritative online-world state.

Clients may predict presentation and movement where appropriate but cannot be the final source of truth for rewardable or competitive outcomes.

## Consequences

The future runtime must support validation, idempotency, auditability, replay resistance, reconnect semantics, and authoritative race/event state. Engine, language, protocol, database, and hosting technology remain open.
