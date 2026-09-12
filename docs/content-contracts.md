# Content Contracts

Content definitions are separate from runtime player state. Catalogs under `design/catalog/` provide stable IDs for prototypes and future services.

## Identifier conventions

| Type | Pattern | Example |
|---|---|---|
| Character | `char_<slug>` | `char_maya_voss` |
| Faction | `faction_<slug>` | `faction_free_roads` |
| District | `district_<slug>` | `district_old_grid` |
| Main quest | `MQ###` | `MQ042` |
| Side quest | `SQ###` | `SQ014` |
| Part | `part_<category>_<slug>` | `part_ecu_legacy_zero` |
| Pet | `pet_<slug>` | `pet_luna` |

IDs are permanent once published. Display names may change; IDs should not.

## Quest definition

Minimum authored quest data:
- id
- chapter where applicable
- title and summary
- district
- primary NPC
- objectives
- prerequisites
- rewards
- relationship/faction consequences
- choice tags
- replay policy

Quest runtime state is separate and may include acceptance/completion timestamps, objective progress, selected choices, and reward grant IDs.

## Vehicle runtime contract

A persistent vehicle should expose:
- `vehicle_id`
- `owner_character_id`
- provenance and prototype lineage
- garage slot
- active build revision
- condition
- reputation
- history references

A build revision should include:
- immutable revision ID
- equipped part IDs
- tuning parameters within accepted bounds
- derived role/classification
- authoritative validation hash/signature or equivalent server record

Race results must reference the exact accepted build revision.

## Design validation

`tools/validate_design.py` verifies:
- expected catalog counts
- uniqueness of IDs
- MQ001–MQ100 continuity
- sequential campaign prerequisite graph
- character/district references
- allowed character faction references
- pet competitive-buff prohibition
- schema parseability

CI runs the validator on pull requests.
