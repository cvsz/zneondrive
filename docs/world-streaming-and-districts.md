# World Streaming & District Implementation

## Canonical city

NOVA CITY launch design contains 10 districts defined in `design/catalog/districts.json` and the world bible.

## Runtime objective

District implementation should support:
- predictable traversal boundaries,
- streaming/level partitioning,
- gameplay server instance ownership,
- district-specific activities,
- social/event density,
- graceful reconnect/spawn.

## Unreal direction

Use UE world/level streaming capabilities appropriate to measured content scale. Do not lock partition sizes or server density before profiling.

## Persistent vs instance state

Persistent service state:
- ownership,
- quest progression,
- unlocks,
- reputation/relationships.

Gameplay-instance state:
- momentary actors/traffic,
- race instance/checkpoints,
- transient combat/chase/event state.

Global/seasonal state:
- public events,
- season phase,
- live-ops rotations.

## District production checklist

Each district needs:
- gameplay purpose,
- traversal routes,
- race/job routes,
- landmarks,
- social spaces,
- streaming/performance budget,
- audio/lighting identity,
- accessibility navigation review,
- spawn/reconnect rules,
- telemetry coverage.

## Evidence

Foundry 9 / Garage 17 remain the first production slice. Full-city streaming/concurrency is not yet verified.
