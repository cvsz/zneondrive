# Runtime Unreal Authority Displacement Envelope v2.2

## Scope

This increment adds a bounded source-level integrity signal to the Unreal Engine 5.8 authoritative prototype vehicle path. It does not change gameplay authority, durable persistence ownership, race acceptance, or the selected Unreal + Go + PostgreSQL + Redis architecture.

## Authority-side envelope

After durable identity is bound, the authoritative pawn retains the location produced by the previous server movement step. Before the next movement step it compares the current authoritative location against that baseline.

The source-level threshold is:

`MaxSpeedCmPerSecond * DeltaSeconds + AuthorityDisplacementSlackCm`

The slack is a fixed server-side prototype setting. The comparison is performed only after `HasAuthority()` and the durable-identity gate have passed.

If the observed displacement exceeds the threshold, the server increments the aggregate `impossible_displacements_total` telemetry counter. No player, vehicle, race, ticket, address, credential, or request identifier is stored in the telemetry snapshot or log line.

## Baseline lifecycle

The displacement baseline is initialized/reset when authoritative durable vehicle identity is successfully bound. It is then updated only after the server completes its swept movement step. This avoids treating the initial authoritative binding position as an integrity event and keeps the detector tied to server-produced movement history.

## Trust boundary

This signal is observational and defense-in-depth only:

- clients do not submit an accepted transform through this path;
- the detector does not authorize or reject durable state;
- PostgreSQL remains the durable service-plane authority;
- Redis remains ephemeral coordination only;
- the Unreal dedicated server remains gameplay authority;
- the detector does not automatically ban, kick, sanction, or alter race results.

## Static evidence

`tools/validate_unreal_server_observability.py` verifies that:

- the counter remains numeric and identifier-free;
- the comparison occurs after authority and durable-identity gates;
- the envelope derives from authoritative max speed, server tick delta, and explicit bounded slack;
- the integrity hook occurs before the authoritative movement step;
- the baseline is updated after authoritative movement;
- durable identity binding resets the baseline before `ForceNetUpdate()`.

## Explicit non-claims

This is source/static-validation evidence only. It does **not** prove:

- a successful Unreal Engine 5.8 Client or dedicated Server build;
- packaged runtime emission of this telemetry;
- final vehicle-physics speed, acceleration, teleport, impulse, collision-recovery, or rollback envelopes;
- acceptable false-positive/false-negative rates under real network/physics conditions;
- live correlation with authoritative race checkpoints;
- ranked anti-cheat readiness or sanctions policy;
- deployed log shipping, alerting, SLOs, load/soak, HA/DR, or production readiness.

The next evidence step is a successful retained UE 5.8 build/package run followed by live packaged Unreal↔Go testing and calibration using real dedicated-server samples.
