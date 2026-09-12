# Multiplayer & Networking

## Authority model

The client is untrusted. Unreal dedicated servers own moment-to-moment gameplay authority; the Go service plane owns durable identity, ownership, inventory, quests and progression.

## Implemented source path

Current Unreal source proves:
- client input reaches a Server RPC,
- server clamps/validates prototype input,
- authority performs movement,
- movement/state replicates to clients,
- gameplay is blocked until a durable service-plane identity is bound,
- client sends a one-time gameplay ticket rather than a long-lived service credential,
- dedicated server redeems the ticket with the Go service and receives the durable snapshot.

This is source-level implementation evidence. Runtime v0.7 also implements a server-only PostgreSQL race start/checkpoint/finish lifecycle that binds results to the exact authoritative active build. A successful packaged live Unreal↔Go run and live race transport remain separate evidence gates.

## Target session lifecycle

1. Client bootstraps/resumes identity with Go.
2. Client obtains short-lived gameplay ticket.
3. Client connects to gameplay server.
4. Client supplies one-time ticket.
5. Gameplay server authenticates to service plane and redeems ticket.
6. Service returns authoritative snapshot.
7. Gameplay server binds vehicle/build/player state.
8. Gameplay begins.
9. Durable outcomes are submitted/committed through service-owned mutation rules.

## Race authority target

A future authoritative race instance must own:
- accepted vehicle build revision,
- start authorization,
- ordered checkpoints/laps,
- penalties,
- finish state,
- disconnect/rejoin policy,
- result signature/audit metadata.

The client may display predicted timing but cannot finalize ranked results.

## Reconnect requirements

Reconnect must not:
- duplicate quest/economy rewards,
- reuse consumed gameplay tickets,
- silently switch active build revision,
- reset penalties,
- bypass event eligibility.

## Transport security target

Production requires encrypted transport, service authentication, replay resistance, rate limiting and secret rotation. Local HTTP is not a production transport claim.

## Scale assumptions

No concurrency number is currently verified. Sharding/instance density, replication graph strategy, bandwidth budget and regional topology require load evidence before being marked complete.
