# Runtime Inventory + Rebuild v0.6

## Scope

Phase 4.2 closes the first authoritative Garage 17 rebuild-state gap without claiming a playable Unreal slice.

Implemented source/runtime contracts:
- durable PostgreSQL inventory quantities,
- durable per-character blueprint unlocks,
- canonical vehicle-part catalog validation,
- MQ004 idempotent salvage grant,
- MQ005 starter-rebuild blueprint unlock,
- MQ009 idempotent recovered-part grant,
- inventory-authoritative vehicle rebuilds,
- immutable build revisions,
- atomic equip/unequip inventory accounting,
- operation-ID replay protection bound to the exact build hash,
- reconnect persistence for inventory, blueprints, and active build,
- Unreal snapshot parsing for inventory and blueprint state.

## Authoritative rebuild transaction

A rebuild request still contains only desired part IDs, expected revision, and an operation ID. The Go/PostgreSQL service treats all client input as untrusted.

Inside one transaction the service:
1. resolves the authenticated character,
2. rejects an operation ID previously used for a different build hash,
3. locks the owned starter vehicle and checks the expected active revision,
4. verifies the starter rebuild blueprint is unlocked,
5. loads the exact current build,
6. consumes newly equipped parts only when inventory quantity is positive,
7. returns removed parts to inventory,
8. appends a new immutable build revision,
9. advances the active revision,
10. commits atomically.

Any failure rolls back inventory and build changes together.

## Quest-linked inventory proof

- MQ004 / Scrap Rights: grants one part_brakes_track_i; retry of the same completion does not duplicate it.
- MQ005 / Borrowed Tools: unlocks bp_starter_rebuild; rebuilds are rejected before the unlock.
- MQ009 / Missing Fastener: grants one part_tires_street_i; retry does not duplicate it.

These concrete prototype rewards implement the vertical-slice acceptance requirements for idempotent salvage/recovery and first rebuild state. They are runtime proof assets, not final economy tuning.

## Catalog boundary

core.NormalizeParts now rejects identifiers absent from design/catalog/vehicle-parts.json. The v0.6 static validator verifies that the runtime allowlist and canonical design catalog have exactly the same part IDs.

## Evidence

The PostgreSQL integration test proves:
- rebuild rejection before blueprint unlock,
- MQ004 grant and retry idempotency,
- MQ005 blueprint persistence,
- authorized part consumption,
- build retry without double consumption or revision advance,
- operation-ID conflict on changed payload,
- rejection of unowned parts,
- MQ009 recovery retry idempotency,
- unequipped-part return,
- reconnect persistence.

GitHub Actions remain the authoritative execution evidence.

## Still open

v0.6 does not prove:
- successful Unreal 5.8 source compilation,
- live packaged Unreal-Go transport,
- Garage 17 environment/playability,
- MQ010 ignition-ready legality,
- MQ011 exact-revision ignition transition,
- race authority,
- anti-cheat,
- load/soak,
- HA/DR,
- production deployment.

Those remain evidence-gated.
