# Runtime Reference Quest Operation Scope v4.5

## Status

Implemented source/unit evidence. This document does **not** claim full reference-oracle parity, live Unreal integration, deployment-scale verification, or production readiness.

## Problem closed by this increment

The durable `quest_completions.operation_id` key is globally unique. Before v4.5, the public quest-completion API forwarded the caller's opaque operation ID directly to PostgreSQL. The store replay path validated the existing quest ID, but the durable key itself was not namespaced by authoritative character identity. A reused caller key could therefore collide with another character's durable quest operation.

## Implemented boundary

The player-facing quest API now resolves the authenticated session snapshot first and derives the PostgreSQL operation key from:

- a fixed domain separator (`zneondrive:quest-operation:v1`);
- the authoritative `CharacterID` returned by the session-backed store; and
- the caller's opaque idempotency key.

The derived value is SHA-256 through the existing Go core hashing primitive. The same character plus the same caller operation ID produces the same durable key across renewed sessions; a different character produces a different key. The raw caller operation ID is no longer the durable quest-operation key used by the public API.

## Trust boundary

- Unreal/player clients still supply only an opaque idempotency key and canonical quest identity.
- Go resolves authenticated durable identity before deriving the operation key.
- PostgreSQL remains the durable quest/progression/reward authority and still owns uniqueness, prerequisite, reward, inventory, blueprint, and Roadworthy transactions.
- Redis is unchanged and remains ephemeral/distributed coordination infrastructure.
- This change does not make the client authoritative for identity, rewards, inventory, blueprints, vehicle state, or quest sequencing.

## Automated evidence

Go core unit coverage verifies deterministic replay for the same character/key, separation across characters and caller keys, and rejection of empty inputs. HTTP API coverage verifies that quest completion resolves authoritative character identity and passes only the scoped derived key to the store rather than the raw caller key.

## Non-claims / remaining gates

This increment does not close:

- successful real UE 5.8 Client or dedicated Server build/package evidence;
- packaged Unreal↔Go integration or race lifecycle evidence;
- Garage 17 / First Ignition / MQ001–MQ012 playable evidence;
- full reference-oracle parity in Go;
- deployed ingress isolation, deployment-scale load, or long-duration soak;
- live final-physics anti-cheat calibration;
- production backup/RPO/RTO, HA/failover, regional DR, or go/no-go approval.
