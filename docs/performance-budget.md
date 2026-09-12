# Performance Budget

**Status:** target budget; not yet verified on a packaged production build.

## Client target

Initial PC vertical-slice target:
- 60 FPS presentation target on agreed reference hardware,
- frame-time spikes measured separately from average FPS,
- scalable graphics settings,
- no gameplay-critical dependency on cinematic-only effects.

Final hardware tiers must be agreed with the client/publisher before claiming compliance.

## Dedicated gameplay server target

Before multiplayer alpha:
- stable server tick appropriate to the selected vehicle/race model,
- no unbounded per-player allocations,
- bounded replication payload,
- measurable instance CPU/RAM per concurrent player,
- overload behavior defined rather than silent degradation.

Exact tick rate and concurrency are intentionally not fixed until real Unreal profiling exists.

## Service-plane targets

Pre-production SLO targets for load testing:
- health/state reads: p95 < 150 ms inside target region,
- durable mutation APIs: p95 < 250 ms under expected alpha load,
- error rate excluding client/domain rejects: < 1%,
- database pool saturation and queueing must be observable.

These are targets, not measured production claims.

## Network budget

Track:
- bytes/sec per player,
- RPC frequency,
- replicated actor count,
- correction rate,
- gameplay-ticket/auth round trips,
- reconnect time.

## Content budgets

Each district/content milestone should define:
- draw/actor budget,
- texture/memory budget,
- streaming-cell budget,
- animation/Niagara budget,
- audio voice count,
- shader/PSO warm-up plan.

## Evidence required

A performance item is complete only when the benchmark records:
- commit/build,
- hardware/environment,
- scenario/player count,
- duration,
- p50/p95/p99 where relevant,
- CPU/RAM/GPU/network metrics,
- failure/saturation point.
