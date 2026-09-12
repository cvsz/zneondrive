# Development

## Current development mode

zNeonDrive is presently a design/contract repository. Runtime client/server technologies are deliberately not selected yet.

## Local validation

Requirements: Python 3.x and standard shell tooling.

```bash
make validate-design
make check-json
make ci
```

No secrets or external services are required to validate the current design layer.

## Content authoring rules

- Never rename a published stable content ID merely to change display text.
- New main quests must preserve explicit prerequisite semantics.
- Character, faction, district, part, pet, and quest references must resolve.
- VIP design changes require fairness review.
- Pet changes must not add ranked competitive buffs.
- Story changes that affect authority/persistence assumptions require an ADR or architecture update.
- Do not mark runtime behavior complete because a design record exists.

## Quality expectations

- Keep changes reviewable.
- Update machine-readable catalogs when canonical prose changes.
- Add validator assertions for new invariants.
- Prefer deterministic/reproducible tooling.
- Never commit credentials or player/user data.
- Never weaken CI/security gates to make a change pass.

## Runtime implementation

Before runtime code begins, accept ADRs for:
1. client/engine,
2. server/runtime,
3. persistence,
4. realtime/network model,
5. deployment topology,
6. observability/security baseline.

The accepted choices must preserve ADR-0001 server authority and ADR-0002 starter-vehicle identity.
