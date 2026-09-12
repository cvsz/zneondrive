# Factions & Relationships

## Factions

Canonical launch factions:
- NOVA Authority,
- Iron Wolves,
- Mechanist Guild,
- VANTEX,
- Free Roads.

Faction standing is independent persistent state; alignment with one faction must not be represented by a single global morality score.

## Relationship model

NPC relationship progression:
- Stranger,
- Acquaintance,
- Friend,
- Trusted,
- Partner / Rival / Enemy.

Relationships can affect dialogue, access, quest variants and story consequences.

## Authority requirements

Durable relationship/faction changes:
- originate from validated quest/story actions,
- are persisted server-side,
- are auditable enough to diagnose progression defects,
- cannot be set directly by a modified client.

## Shared-world consistency

Personal story consequences may alter a player's view/access while global shared-world state remains coherent for other players.

## Initial runtime scope

Relationship/faction persistence is designed but not yet implemented in the current v0.6 PostgreSQL schema. It remains a Phase 4/6 gate.
