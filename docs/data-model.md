# Data Model — Durable Runtime v0.6

**Status:** implemented prototype schema in PostgreSQL migrations.

## Authority

PostgreSQL is the durable source of truth for the current service-plane prototype. Unreal gameplay servers consume/validate session state but do not own permanent player ownership or economy state.

## Current entities

### accounts

- `id` — stable account identifier
- `resume_key_hash` — unique credential hash
- `created_at`

One account currently maps to one primary character.

### characters

- `id`
- `account_id` — unique
- `money`
- `xp`
- `reputation`
- `created_at`

### vehicles

- `id`
- `owner_character_id`
- `active_revision`
- `starter_lineage`
- `roadworthy`
- `created_at`

A partial unique index enforces at most one starter-lineage vehicle per character.

### vehicle_builds

Immutable/reconstructable build history:
- `vehicle_id`
- `revision`
- `part_ids` JSONB
- `validation_hash`
- `operation_id`
- `created_at`

Primary key: `(vehicle_id, revision)`.

### sessions

- hashed session token
- account owner
- expiration
- creation time

### quest_completions

- character
- quest ID
- unique operation ID
- reward money/xp/reputation
- completion time

Primary key prevents duplicate completion of the same quest for a character.

### game_tickets

One-time server-redemption credentials:
- `ticket_hash`
- account
- expiry
- consumed timestamp
- creation time

### inventory_items

- character
- item ID
- non-negative quantity
- update timestamp

Primary key: `(character_id, item_id)`.

### character_blueprints

- character
- blueprint ID
- unlock timestamp

Primary key: `(character_id, blueprint_id)`.

## Invariants

- monetary and XP balances are non-negative,
- vehicle build revisions start at 1,
- starter vehicle lineage is preserved,
- build operation IDs are unique,
- quest operation IDs are unique,
- inventory quantities cannot be negative,
- ticket consumption is single-use,
- foreign keys cascade with account/character deletion where currently defined.

## Migration policy

Future migrations must:
1. be additive when possible,
2. preserve durable IDs,
3. define rollback/forward-recovery behavior,
4. avoid destructive production migrations without backup/restore evidence,
5. include tests for data compatibility.

## Planned domains not yet represented in schema

- faction standing,
- NPC relationships,
- crews/roles,
- race registrations/results,
- entitlement/VIP grants,
- social presence,
- moderation actions,
- live-ops event definitions,
- telemetry/audit event persistence.

These are planned, not implemented claims.
