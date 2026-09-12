# Vehicle Physics & Race Integrity

**Status:** design/engineering target. Final vehicle physics are not yet implemented.

## Vehicle handling goals

The same persistent vehicle can support multiple builds while preserving a recognizable identity.

Target build roles:
- Street,
- Drift,
- Circuit,
- Drag,
- Off-road,
- Delivery,
- Hybrid.

Each role should create meaningful trade-offs rather than a single dominant stat stack.

## Authority split

Client:
- input capture,
- local presentation/prediction,
- camera/audio/FX.

Gameplay server:
- accepted input constraints,
- authoritative transform/vehicle state,
- race checkpoints/timing,
- eligibility/build revision.

Service plane:
- durable build ownership/revision,
- inventory/blueprints,
- durable race/reward result when implemented.

## Race integrity requirements

Before ranked play:
- bind entry to exact accepted build revision,
- server-owned start time,
- ordered checkpoint/lap validation,
- impossible speed/teleport telemetry,
- disconnect/rejoin policy,
- penalty authority,
- result audit record,
- idempotent reward/result commit.

## Physics evidence

Final handling claims require:
- chosen vehicle physics implementation,
- reference hardware/server tick,
- deterministic/acceptable replication behavior,
- latency/correction tests,
- competitive exploit tests,
- tuning acceptance by game design.

The current prototype pawn proves authority direction only.
