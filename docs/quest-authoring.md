# Quest Authoring Guide

## Stable identity

Main quests use `MQ###`. Published IDs are compatibility contracts.

## Required quest fields

Every implementation-ready quest should define:
- title/display key,
- chapter,
- prerequisites,
- primary district/NPC references,
- objectives,
- success criteria,
- failure/retry behavior,
- choices,
- durable state changes,
- rewards,
- relationship/faction effects,
- replay policy,
- accessibility/localization notes where relevant.

## Durable mutation rule

Quest completion and reward grants must be authoritative and idempotent. Client presentation may show completion immediately, but durable reward state is committed by the service plane.

## Sequence

The current service prototype enforces sequential MQ progression. Any future branching relaxation requires an explicit quest-state design and migration.

## Narrative implementation

For each playable quest add:
- dialogue script,
- cinematic beats,
- environment/interaction requirements,
- telemetry events,
- QA acceptance cases.

## QA cases

At minimum test:
- normal completion,
- reconnect mid-objective,
- duplicate completion request,
- prerequisite violation,
- reward retry,
- stale content/version behavior where applicable.

## Source hierarchy

Canonical story intent comes from the Game Design Bible; stable IDs/references come from machine-readable catalogs.
