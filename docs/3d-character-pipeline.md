# zNeonDrive 3D Character Pipeline

## Goal

Provide a production-ready, data-driven character pipeline without making the core game dependent on a specific character vendor or licensed asset pack.

## Runtime contract

Each playable/NPC character is represented by a `UNDCharacterDefinition` primary data asset. The asset owns stable game identity and references presentation assets:

- `CharacterId` — immutable content identifier.
- `FactionId` / `ArchetypeId` — game-domain classification; these remain independent of the mesh.
- `BodyMesh` — soft reference to the skeletal mesh.
- `AnimationClass` — soft reference to the Animation Blueprint class.
- `LocomotionProfileId` — selects the movement/animation profile.
- `CharacterScale` — bounded presentation scale.
- `bUseHighFidelityFacialRig` — presentation capability flag.
- `bReplicateCosmeticState` — explicit network policy for cosmetic state.

`UNDCharacterComponent` applies the definition to an actor's skeletal mesh at runtime. Asset loading is soft and therefore keeps the gameplay module decoupled from optional character packs.

## Recommended character stack

### Tier A — gameplay skeleton

- Standard UE skeletal mesh + Animation Blueprint.
- Locomotion state machine.
- Motion Matching where the installed UE 5.8 content/plugins support it.
- Stride/turn-in-place handling.
- Foot IK and ground alignment.
- Network-safe authoritative movement state.

### Tier B — facial presentation

- Face Control Rig.
- Facial animation layers separated from locomotion.
- Optional Live Link / performance capture.
- Lip-sync and dialogue-driven facial state as cosmetic presentation only.

### Tier C — high-fidelity characters

MetaHuman can be used as an optional asset source for selected hero characters. Do not copy Epic-licensed character assets, DNA, textures, grooms, or generated packages into the public repository unless their applicable license explicitly permits it.

Epic's OpenRigLogic project provides open-source DNA/RigLogic libraries, but the character assets and surrounding Unreal-licensed technology remain separate concerns. Treat the MetaHuman asset pipeline as an external content dependency and keep generated assets out of source control unless licensing has been verified.

## Character architecture

```text
Go authoritative profile
        |
        | stable CharacterId / faction / archetype
        v
UNDCharacterDefinition
        |
        +---- BodyMesh ---------> Skeletal Mesh
        +---- AnimationClass ----> Animation Blueprint
        +---- LocomotionProfile -> locomotion/motion-matching profile
        +---- Facial capability -> optional Control Rig / facial layer
        |
        v
UNDCharacterComponent
        |
        v
UE character actor / pawn
```

## Network boundary

Character appearance is not authoritative gameplay state. The server owns identity, progression, faction, inventory and competitive state. Clients receive validated character identity/presentation references and may render optional high-fidelity cosmetics locally.

Do not send skeletal transforms every frame as a replacement for UE character movement replication. Use the existing authoritative movement model and replicate only the minimum cosmetic state needed to reproduce the presentation.

## Content layout

The intended content layout is:

```text
Content/
  Characters/
    Definitions/
    Hero/
    NPC/
    Shared/
    Animation/
      Locomotion/
      Facial/
      Additive/
      IK/
      MotionMatching/
    Materials/
    Groom/
    Retargeting/
```

Binary assets should be acquired/generated in an authorized UE/Fab/content environment and are intentionally not represented as fake placeholders in Git.

## Character production targets

### Launch character contract

The existing design bible contains 25 named launch characters. Each should eventually have:

1. Stable character ID.
2. Faction and archetype metadata.
3. Body/skeletal asset.
4. Locomotion profile.
5. Facial profile where applicable.
6. Clothing/material variants.
7. LOD policy.
8. Animation coverage for idle/walk/run/turn/interaction/race-adjacent states.
9. Dialogue/emote presentation hooks.
10. Content validation evidence.

### Hero characters

Use high-fidelity characters only where they materially improve the scene. Hero characters can use MetaHuman-compatible assets, facial rigs and performance capture; ordinary crowd/NPC characters should use cheaper shared rigs, LODs and animation sets.

### Crowd characters

For districts, use shared skeletal/animation assets and aggressive LOD/instance strategies. Do not create 25 unique high-cost rigs merely because there are 25 named characters.

## Validation gates

A character is not considered production-ready until:

- the definition asset resolves its required mesh;
- the Animation Blueprint loads without warnings;
- skeleton/retarget compatibility is validated;
- LODs exist for the intended camera ranges;
- animation state coverage is exercised;
- dedicated-server builds do not require client-only character presentation assets;
- packaged-client load is verified;
- networked spawn/reconnect preserves stable character identity;
- content licensing is recorded outside the public source tree.
