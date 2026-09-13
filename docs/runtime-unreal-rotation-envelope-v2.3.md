# Runtime Unreal Authority Rotation Envelope v2.3

## Scope

This increment extends the UE 5.8 authoritative prototype pawn with a bounded yaw-change integrity signal. It does not change movement authority, durable persistence, gameplay-ticket authorization, race acceptance, or the selected Unreal + Go + PostgreSQL + Redis architecture.

Before each authoritative movement step, after the existing authority and durable-identity gates have passed, the server reads the current authoritative yaw and compares the wrap-safe change from the previous authoritative baseline against:

`TurnRateDegreesPerSecond × DeltaSeconds + AuthorityRotationSlackDegrees`

If the observed yaw delta exceeds that envelope, the dedicated-server aggregate telemetry increments `impossible_rotations_total`.

## Trust boundary

The detector uses only server-owned state and server configuration:

- `GetActorRotation().Yaw` from the authoritative pawn;
- `TurnRateDegreesPerSecond` from the server-side movement model;
- server `DeltaSeconds`;
- bounded `AuthorityRotationSlackDegrees`.

It does not accept a client transform or client clock. It cannot authorize gameplay, mutate PostgreSQL durable state, alter Redis coordination state, override race acceptance, kick/ban a player, or bypass the server-only gameplay-ticket redemption boundary.

## Baseline handling

The yaw baseline is reset when durable vehicle identity is authoritatively bound. This avoids treating spawn/bind rotation as suspicious movement. During normal authoritative ticks, the baseline is updated after the server applies authoritative rotation/movement.

`FMath::FindDeltaAngleDegrees` is used before taking the absolute value so wraparound at ±180 degrees does not create a false large rotation.

## Evidence

Repository CI statically requires:

- the detector to remain behind `HasAuthority()` and the durable-identity gate;
- the threshold to remain derived from authoritative turn rate and tick delta;
- explicit bounded server-side angular slack;
- wrap-safe yaw delta calculation;
- `RecordImpossibleRotation()` to run before authoritative rotation is applied;
- the yaw baseline to update only after authoritative rotation;
- durable identity binding to reset the yaw baseline before `ForceNetUpdate()`;
- the aggregate telemetry contract to remain numeric-only and free of credentials/dynamic identifiers.

## Explicit non-claims

This is source/static-validation evidence only. It does not prove:

- successful compilation or runtime emission on a real UE 5.8 runner;
- packaged Client/Dedicated Server integration;
- calibration against final vehicle physics;
- acceleration, angular-velocity, drift, teleport, rewind, or collision-exploit detection completeness;
- false-positive/false-negative rates;
- live race correlation or ranked sanctions;
- production anti-cheat readiness.

Those remain evidence-gated until real UE 5.8 build/package and live-runtime evidence exists.
