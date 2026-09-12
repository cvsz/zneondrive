# Contributing

Thanks for contributing to zNeonDrive.

## Development workflow

1. Create a focused branch from `main`.
2. Keep changes reviewable.
3. Run `make ci`.
4. Update canonical prose and machine-readable catalogs together.
5. Update `CHANGELOG.md` for material design/contract changes.
6. Open a pull request and describe content, compatibility, security/fairness, and rollback impact.

## Stable content IDs

Published IDs are compatibility contracts.

Do not rename IDs such as `char_maya_voss`, `district_foundry_9`, or `MQ001` solely because display text changes.

If an ID truly must be retired, document migration/alias behavior before removing it.

## Story/content changes

When changing story:
- preserve chapter prerequisite integrity,
- update referenced NPC/district IDs,
- document new branch consequences,
- keep world-state scope explicit: global, shard/community, or character,
- avoid making one player's ending incompatible with the shared world.

## Vehicle/fairness changes

- Do not introduce direct VIP competitive stat advantages.
- Pets must not grant ranked race-stat buffs.
- Competitive parts/build rules must remain authoritatively validatable.
- Unique story vehicle provenance must remain auditable.

## Architecture changes

Changes to ownership, persistence, competitive authority, monetization/fairness, or starter-vehicle identity should update architecture/ADRs.

## Branch naming

Use prefixes such as `feat/`, `design/`, `fix/`, `docs/`, `chore/`, `security/`, or `runtime/`.

## Security

Do not report exploitable vulnerabilities publicly. Follow `SECURITY.md`.
