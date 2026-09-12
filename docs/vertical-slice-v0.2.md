# Vertical Slice v0.2 — Garage 17 / Foundry 9

## Objective

Prove the emotional and authoritative core of PROJECT: NEON DRIVE before selecting a production engine/server stack. The slice covers MQ001–MQ012 and must make the player care about one broken vehicle while demonstrating safe persistence boundaries.

## Scope

**Included:** account resume, primary character, Garage 17, Foundry 9 traversal, starter vehicle ownership, basic salvage/inventory hooks, part fitting, build revision, Maya relationship state, Legacy ECU reveal, First Ignition, Roadworthy, one legal time trial, one underground race proof, save/reconnect.

**Excluded:** full city streaming, 100 quests in runtime, production matchmaking, monetization, real-money currency, complete anti-cheat, crews, housing, live ops, HA/DR claims.

## First-session UX

1. **Resume/Create identity** — player reaches one authoritative primary character.
2. **Arrival After Midnight (MQ001)** — phone message from Rowan Gray establishes urgency and distrust of police.
3. **Garage 17 (MQ002)** — enter the garage and meet Maya.
4. **A Chassis With No Name (MQ003)** — inspect the prototype and claim ownership.
5. **Scrap Rights (MQ004)** — learn salvage interaction and legal ownership boundaries.
6. **Borrowed Tools (MQ005)** — fit first part/blueprint and persist a build revision.
7. **Dead Grid (MQ006)** — traverse Foundry 9 and restore a local power/diagnostic route.
8. **First Fuel (MQ007)** — make a small authored choice with persistent consequence.
9. **Legacy ECU (MQ008)** — discover the anomalous DRIVE ZERO hardware signature.
10. **Missing Fastener (MQ009)** — short recovery/search beat; optional Luna foreshadowing.
11. **Cold Crank (MQ010)** — perform pre-ignition checks and validate build legality.
12. **First Ignition (MQ011)** — start the vehicle; authoritative ownership/build state becomes drive-capable.
13. **Roadworthy (MQ012)** — complete a controlled route and commit the chapter boundary.

## Quest acceptance matrix

| Quest | Required proof |
|---|---|
| MQ001 | message acknowledged; Garage 17 objective activated exactly once |
| MQ002 | Maya introduced; garage checkpoint persisted |
| MQ003 | unique starter vehicle bound to player character |
| MQ004 | salvage grant is idempotent |
| MQ005 | first build revision can be created; stale expected revision is rejected |
| MQ006 | traversal objective survives reconnect |
| MQ007 | authored choice stored once and visible after reconnect |
| MQ008 | Legacy ECU discovery flag cannot be client-forged into reward state |
| MQ009 | recovered item cannot be duplicated by retry |
| MQ010 | illegal/incomplete build cannot enter ignition-ready state |
| MQ011 | ignition transition references exact active build revision |
| MQ012 | chapter completion, rewards, reputation, and Roadworthy state are idempotent |

## Player movement and interaction target

Technology-neutral minimum:
- move, sprint, camera/aim/look,
- context interaction,
- inspect/pickup/use,
- enter/exit vehicle,
- garage/workbench interaction,
- clear interaction priority when multiple targets overlap,
- no progression-critical interaction requires pixel-perfect pointing.

The production engine may add prediction, animation root motion, physics, and streaming later; authoritative ownership/reward state remains server-side.

## Starter vehicle feel target

The slice should communicate “broken but worth saving” rather than final supercar performance.

Target properties:
- slow-to-moderate acceleration,
- readable weight transfer,
- forgiving low-speed steering,
- visible/audio mechanical imperfection before repair,
- improved response after MQ010/MQ011 without becoming high-tier,
- setup differences must be attributable to part/build state, never VIP.

Exact physics constants are deferred until engine selection and instrumented handling tests.

## Legal race proof

A short Foundry 9 sanctioned time trial validates:
- event registration,
- accepted build revision,
- ordered checkpoint completion,
- finish timing,
- no reward duplication on result retry.

## Underground race proof

A short invitation-only route proves a second ruleset and different presentation while reusing the same authoritative result contract. It is not required to implement faction-wide matchmaking in this slice.

## Reconnect/save semantics

On disconnect/reconnect:
- account resumes the same primary character,
- starter vehicle ownership is unchanged,
- latest committed build revision is restored,
- completed quest rewards are not re-granted,
- in-progress objectives resume from the last committed objective boundary,
- transient movement position may recover to a safe checkpoint,
- race participation follows event-specific reconnect policy and never trusts the client finish state.

## Accessibility baseline

The slice must be implementable with:
- remappable controls,
- subtitle/caption support with speaker identification,
- UI text scaling target,
- non-color-only objective/race cues,
- reduced camera shake option,
- hold/toggle alternatives for repeated actions where practical,
- timing-sensitive accessibility review before ranked launch.

## Graybox/content checklist

Required spaces:
- NOVA CITY arrival pocket,
- approach road to Garage 17,
- Garage 17 interior/workbench,
- Foundry 9 salvage yard,
- diagnostic/power route,
- Cold Crank bay,
- controlled Roadworthy loop,
- short legal time-trial route,
- short underground route.

Required content beats:
- Rowan message,
- Maya first meeting,
- prototype reveal,
- first meaningful part fit,
- Legacy ECU anomaly,
- ignition moment,
- first drive out of Garage 17.

## Telemetry contract for prototype evaluation

Measure, without collecting unnecessary personal data:
- time to Garage 17,
- time to claim starter vehicle,
- time per rebuild beat,
- failed/stale build mutations,
- reconnect recovery success,
- duplicate mutation attempts blocked,
- time to First Ignition,
- time to Roadworthy,
- legal/underground race completion/failure reasons.

## Exit criteria

The slice is accepted when all machine-testable authoritative invariants pass and a playable client can demonstrate MQ001–MQ012 without losing or duplicating persistent state. Playability itself cannot be claimed by this repository until a selected client/runtime implementation provides evidence.
