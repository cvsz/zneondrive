# zNeonDrive 3D Character Production Pipeline

## Goal

Provide one production contract for the 25-character roster while keeping gameplay identity, server authority, and presentation assets independent.

## Runtime layers

1. **Identity** — `design/catalog/characters.json` is the source of truth for stable IDs.
2. **Presentation profile** — `design/catalog/character-presentation.json` selects tier, locomotion, facial profile, and optimization policy.
3. **UE runtime definition** — `UNDCharacterDefinition` binds a character ID to soft-referenced skeletal and animation assets.
4. **Actor component** — `UNDCharacterComponent` applies the definition to an actor mesh.
5. **Content packs** — actual `.uasset`/`.umap`/texture/groom files stay outside the source-control contract unless redistribution is explicitly permitted.

## Character tiers

| Tier | Purpose | Facial | Groom | Motion matching |
|---|---|---:|---:|---:|
| Hero | story/cinematic/high-value NPC | yes | optional | yes |
| Player | local player avatar | yes | optional | yes |
| NPC | ordinary interactive characters | no | no | yes |
| Crowd | background population | no | no | no |

## Asset contract

Every production character package should provide:

- one approved humanoid skeleton profile or a documented retarget source;
- body skeletal mesh with collision and physics asset;
- material instances with texture budgets appropriate to the presentation tier;
- idle, walk, run, stop, turn, pivot and locomotion transition coverage;
- root-motion policy explicitly declared per animation set;
- IK-ready foot and hand bones where required;
- LOD chain and screen-size thresholds;
- animation budget and update-rate policy;
- optional facial Control Rig / MetaHuman-compatible profile for approved hero assets;
- optional groom with non-groom fallback;
- provenance and redistribution record.

## Animation pipeline

**Source animation → retarget → cleanup → locomotion database → Motion Matching/BlendSpace → IK → gameplay state machine → network presentation.**

The authoritative server replicates gameplay state, not animation assets. Clients resolve the appropriate presentation profile locally.

## Facial pipeline

Use Control Rig and the project's approved facial profile for ordinary characters. Use MetaHuman-compatible assets only where licensing and distribution rights are documented. Facial state is cosmetic and must not become an authorization primitive.

## Performance policy

- Keep dedicated servers free of client-only character presentation dependencies.
- Use soft references for optional character assets.
- Avoid synchronous asset loads during latency-sensitive gameplay paths; preload/stream presentation assets at spawn or scene-transition boundaries.
- Use lower-cost animation update rates and LODs for distant NPCs.
- Disable groom/facial evaluation for characters outside the configured relevance range.
- Treat crowd characters as a separate optimization class rather than scaling hero rigs to hundreds of instances.

## Networking policy

Replicate stable character identity and gameplay state. Cosmetic selections are validated server-side and replicated as bounded IDs. Never trust client-provided mesh paths, animation classes, asset URLs, or arbitrary object references.

## EOS / Pixel Streaming

EOS remains an optional platform adapter for identity, sessions, social, reports, sanctions, and anti-cheat signals. Pixel Streaming is a delivery plane and never an authorization boundary. Both integrations must preserve the same character identity and authoritative gameplay contract as native clients.

## Asset provenance

Do not commit Epic licensed/generated character content merely to make CI green. The repository stores contracts, schemas, validation, and integration boundaries; authorized build environments provide the actual licensed content.

## Completion gates

A character is production-ready only when:

- catalog ID is unique and validated;
- presentation profile exists;
- skeleton/retarget validation passes;
- mesh/material/LOD validation passes;
- animation coverage passes;
- facial/groom dependencies are licensed and optional fallbacks exist;
- dedicated-server packaging excludes client-only content;
- native packaged-client smoke test passes;
- network spawn/despawn/reconnect behavior passes;
- Pixel Streaming smoke test passes when streaming is enabled;
- asset provenance is recorded.
