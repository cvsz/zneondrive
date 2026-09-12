# Runtime Authoritative Race v0.7

## Scope

v0.7 implements the first durable authoritative race instance/result slice in the Go/PostgreSQL service plane. It does not claim that a race is playable in Unreal yet.

## Trust boundary

Player clients cannot submit durable race results directly. Race lifecycle mutations are exposed only under `/v1/internal/...` and require the same `X-Game-Server-Key` secret used by dedicated-server gameplay-ticket redemption.

A race start request supplies account, vehicle, race ID, and an idempotency operation ID. PostgreSQL resolves the account's primary character and owned vehicle, verifies that the vehicle is Roadworthy, then binds the race instance to the vehicle's exact active build revision and stored validation hash. Build identity is therefore derived from durable authoritative state rather than trusted from a client payload.

## Lifecycle

1. Dedicated server redeems a one-time gameplay ticket and obtains the authoritative snapshot.
2. Dedicated server starts a race through `POST /v1/internal/races/start`.
3. Service creates an immutable race binding to account, character, vehicle, build revision, and validation hash.
4. Dedicated server records ordered checkpoints through `POST /v1/internal/races/{raceInstanceID}/checkpoints`.
5. Checkpoint index must equal the server-side cursor and elapsed time must increase monotonically.
6. Dedicated server finishes through `POST /v1/internal/races/{raceInstanceID}/finish`.
7. Finish checkpoint count must equal the recorded cursor and finish time must exceed the last checkpoint.
8. The service derives a deterministic result hash including race instance, race ID, account, character, vehicle, exact build revision/hash, checkpoint count, and finish elapsed time.

## Idempotency

Race start, checkpoint, and finish writes each require an operation ID. Replaying the same operation with the same semantic payload returns the existing state/result. Reusing an operation ID for a different mutation is rejected.

## Persistence

Migration `004_race_runtime.sql` adds:

- `race_instances`
- `race_checkpoints`
- `race_results`

The runtime stores the immutable build binding on `race_instances`, ordered checkpoint evidence, and the final result hash.

## Evidence in this revision

Automated evidence includes:

- core validation tests for race IDs/checkpoints and deterministic build-bound result hashes;
- HTTP tests proving race start is rejected without the dedicated-server key;
- PostgreSQL integration coverage for MQ012 Roadworthy prerequisite, idempotent race start, immutable build binding, ordered/monotonic checkpoints, out-of-order rejection, finish/result persistence, finish replay, and post-finish checkpoint rejection.

## Explicitly not proven

This revision does not prove:

- packaged Unreal Engine 5.8 ↔ Go live race transport;
- Garage 17 / First Ignition / Foundry 9 playable race content;
- physics-derived anti-cheat or impossible-state detection;
- distributed race orchestration, matchmaking, load/soak capacity, HA, failover, or DR;
- production deployment readiness.

Those remain evidence-gated roadmap items.
