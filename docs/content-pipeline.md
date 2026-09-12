# Content Pipeline

## Principle

Narrative prose, machine-readable content IDs and runtime implementation must remain synchronized.

## Content layers

1. **Canon** — Story Bible / Game Design Bible.
2. **Catalogs** — stable machine-readable IDs under `design/catalog/`.
3. **Schemas** — structural contracts under `design/schemas/`.
4. **Runtime content** — Unreal/Go representations derived from canonical contracts.
5. **Localization** — keyed display text/dialogue.
6. **Release content** — versioned and validation-passing package.

## Stable IDs

Published IDs are compatibility contracts. Display text may change; IDs should not be renamed casually.

Examples:
- `MQ001`
- `char_maya_voss`
- `district_foundry_9`
- `faction_vantex`

## Quest authoring lifecycle

1. narrative intent,
2. objectives/choices/failure criteria,
3. stable references,
4. rewards/state transitions,
5. graybox implementation,
6. dialogue/cinematic implementation,
7. localization,
8. QA,
9. telemetry hooks,
10. release validation.

## Content review gates

Check:
- reference integrity,
- progression order,
- server-authoritative mutation requirements,
- VIP/fairness impact,
- relationship/faction consequences,
- age-rating/sensitivity impact,
- localization context,
- rollback compatibility.

## Hotfix policy target

Data-driven content should be separable from executable code where safe, but hotfix capability must not bypass validation, signature/audit controls or stable-ID compatibility.
