# Reference Parity Evidence Index v2.8

**Status:** documentation/CI-governance hardening only. The promoted runtime baseline remains Phase 4.19 / Runtime v2.3 and full reference-oracle parity remains open.

## Purpose

Reference-parity evidence added after the promoted Unreal runtime baseline had become difficult to discover from the canonical documentation index. This document provides one explicit evidence map without promoting those test-harness increments into live Unreal, deployment, or production evidence.

## Covered reference-parity increments

| Increment | Evidence | Scope |
| --- | --- | --- |
| v2.4 | `runtime-reference-parity-v2.4.md` | shared Python/Go deterministic build-hash and quest-ID/predecessor vectors |
| v2.5 | `runtime-reference-parity-integrity-v2.5.md` | fixture integrity: non-empty sets, uniqueness, canonical hashes, valid predecessor declarations |
| v2.6 | `runtime-reference-operation-parity-v2.6.md` | payload-bound idempotency and conflicting operation-ID reuse rejection in the Python oracle |
| v2.7 | `runtime-reference-rebuild-parity-v2.7.md` | Garage 17 MQ004/MQ005/MQ009 inventory/blueprint effects and inventory-authoritative rebuild semantics |

## Authority boundary

None of these increments changes the selected production architecture:

- Unreal Engine 5.8 dedicated servers remain gameplay authority.
- Go 1.27 remains the authenticated service plane.
- PostgreSQL remains durable authority.
- Redis remains ephemeral coordination and distributed abuse-control state.
- Python remains a dependency-free executable reference oracle and test surface only.

Reference tests must not become an alternate production authority, accept client-asserted durable state, or bypass server-only gameplay-ticket redemption.

## Evidence boundary / non-claims

This index does **not** close any of the following gates:

- full Python-reference ↔ Go/PostgreSQL semantic parity;
- successful retained UE 5.8 Client or dedicated Server build evidence;
- retained UE 5.8 Client/Server cook/package evidence;
- live packaged Unreal ↔ Go bootstrap/session/ticket redemption;
- playable Garage 17 / First Ignition / MQ001–MQ012 flow;
- live Unreal authoritative race lifecycle;
- final-physics anti-cheat calibration or ranked sanctions;
- deployment-scale load, long-duration soak, deployed SLO evidence, HA/DR, or production deployment.

## Governance

`tools/validate_evidence_status_sync.py` must require the canonical documentation index to retain links to v2.4–v2.8 while separately requiring the promoted public status and live UE build/package/integration gates to remain evidence-gated. This prevents discoverability fixes from accidentally being interpreted as production-state promotion.
