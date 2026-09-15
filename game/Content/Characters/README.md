# 3D Character Content

This directory is the controlled import boundary for zNeonDrive character assets.

## Rules

- Do not commit licensed/generated character packages unless redistribution rights are explicitly recorded.
- Keep source-control contracts (`design/catalog/characters.json`) separate from binary presentation assets.
- Use `UNDCharacterDefinition` assets as the runtime binding between stable character IDs and visual assets.
- Keep dedicated-server builds free of client-only presentation dependencies.
- Prefer shared skeletons, animation sets and LODs for NPC/crowd characters.
- Reserve high-fidelity facial rigs and dense groom assets for approved hero characters.

## Expected layout

```text
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

The repository intentionally contains no fake `.uasset` files. Real assets must be imported through an authorized Unreal/Fab/content pipeline and then validated against the character content contract.
