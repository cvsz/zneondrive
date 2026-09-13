# Security Threat Exercise v3.0

**Status:** CI-exercised partial evidence. This does **not** establish production security readiness.

This evidence slice turns the existing zNeonDrive threat-model baseline into a machine-checked exercise manifest without changing the selected architecture or authority boundaries:

- Unreal Engine 5.8 Dedicated Server remains gameplay authority.
- Go 1.27 remains the authenticated service plane.
- PostgreSQL remains durable authority.
- Redis remains ephemeral coordination / distributed abuse-control state.
- Player clients, client clocks, client-computed race outcomes, forwarding headers from untrusted peers, and client-supplied durable state remain untrusted.

## What is exercised in CI

`security/threat-exercise-v3.0.json` records concrete abuse scenarios and binds each `ci_evidenced` scenario to repository evidence markers. `tools/validate_threat_exercise.py` fails CI if those evidence files or markers disappear, if the architecture declaration changes, if core abuse categories disappear, or if open production/live blockers are silently promoted.

The current CI-evidenced set covers:

1. forged session credentials, incorrect game-server keys, and one-time gameplay-ticket replay rejection through the integrated HTTP/PostgreSQL/Redis trust boundary;
2. untrusted `X-Forwarded-For` spoof resistance and trusted-ingress identity separation;
3. Redis-coordinated distributed rate-limit budget plus bounded local fallback semantics;
4. server-only authoritative race start/checkpoint/result integrity rules;
5. PostgreSQL-authoritative inventory/blueprint/rebuild validation and immutable build revisions;
6. idempotent quest/item grants and replay-resistant durable mutation semantics.

These are regression/evidence bindings around controls already implemented and tested. The manifest does not create a second authorization plane and is not consulted by production request handling.

## Explicitly open

The manifest intentionally keeps the following scenarios `open`:

- deployed ingress header sanitization and direct service-bypass prevention;
- live packaged Unreal Engine 5.8 ↔ Go ticket/session trust-boundary evidence;
- final-physics anti-cheat calibration, false-positive/false-negative evaluation, and ranked sanctions policy;
- privileged admin/live-ops RBAC and immutable audit trail.

Additional production gates remain open in the repository, including real UE Client/Server build/package artifacts, Garage 17/MQ001–MQ012 playable evidence, deployed observability/SLOs, deployment-scale load and soak, production backup custody/PITR/RPO/RTO, HA/failover, regional DR, and production deployment approval.

## Evidence boundary

Passing this validator means only that the documented CI-evidenced abuse scenarios remain tied to concrete repository controls/tests and that known live/deployment blockers remain explicit. It does **not** prove:

- TLS or secret rotation on a deployed environment;
- live Unreal transport security;
- ranked anti-cheat effectiveness;
- deployed ingress isolation;
- deployment-scale capacity or DDoS resistance;
- production incident response, HA, DR, or recovery objectives.

Production security must remain evidence-gated until those environment- and runtime-dependent controls are exercised on the agreed deployment target.
