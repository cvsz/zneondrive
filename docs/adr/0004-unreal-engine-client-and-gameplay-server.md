# ADR-0004: Unreal Engine 5.8 for the Playable Client and Dedicated Gameplay Server

- Status: Accepted for vertical-slice and alpha implementation
- Date: 2026-09-12

## Context

PROJECT: NEON DRIVE requires a high-fidelity vehicle-focused 3D client, large streamed environments, multiplayer replication, dedicated-server support, tooling for designers, cinematic presentation, and a path toward competitive authoritative gameplay.

Epic's Unreal Engine documentation explicitly describes a client-server model where the server moderates the true game state, and dedicated servers are recommended for large-scale or competitive multiplayer. That directly matches ADR-0001.

## Decision

Use **Unreal Engine 5.8** as the primary playable client technology for the next implementation phase.

Use:
- C++ for authoritative/network-sensitive gameplay systems and reusable core components,
- Blueprint for designer-authored presentation, encounter flow, UI glue, and rapid iteration where safe,
- Unreal dedicated-server builds for moment-to-moment authoritative gameplay such as movement, race state, checkpoints, vehicle session rules, and event instances.

The existing Python reference package remains a contract/test oracle only. It is not the production game server.

## Consequences

### Benefits
- native client/server multiplayer model aligned with project authority rules,
- dedicated-server build support,
- strong 3D environment, vehicle, animation, cinematic, UI, and world tooling,
- one gameplay technology for client prediction plus authoritative session logic.

### Costs
- larger build/toolchain footprint,
- dedicated-server source-build/CI complexity,
- C++ engineering requirements,
- content-production pipeline and asset-performance discipline become critical.

## Guardrails

- Competitive ownership, economy, inventory, quest rewards, entitlements, and long-lived progression do not live only inside a single Unreal match server.
- Dedicated gameplay servers request/commit durable mutations through the service plane defined in ADR-0005.
- Client RPCs remain untrusted.
- Any engine upgrade requires multiplayer, physics, replay, save compatibility, and performance regression testing.
