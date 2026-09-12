# PROJECT: NEON DRIVE — Game Design Bible v0.2

**Status:** canonical pre-production design contract  
**Supersedes:** v0.1 for implementation planning while preserving all v0.1 lore  
**Setting:** NOVA CITY, 2097  
**Rating target:** 18+  
**Runtime stance:** persistent online world, server-authoritative competitive/persistent state

## Locked pillars

1. **Vehicle as identity.** The starter prototype is persistent provenance, not a disposable mount.
2. **Earn the road.** Progress comes from play, craft, discovery, relationships, reputation, and competition.
3. **One account, one primary character.** Alternate-character systems are outside the launch baseline.
4. **Free baseline remains competitively complete.** VIP changes capacity/convenience only.
5. **Shared world, personal consequences.** Major choices use character/crew overlays without fragmenting the global seasonal world.
6. **Server-authoritative value.** Ownership, rewards, inventory, quest completion, builds, ranked eligibility, race results, crew roles, and entitlements are never trusted from the client.

## Canonical campaign spine

The 100-quest, seven-chapter graph in `design/catalog/quests-main.json` remains canonical:

1. The Broken Machine
2. Street Reputation
3. Underground
4. The Championship
5. VANTEX
6. Drive Zero
7. War for Nova

The first production target is not all 100 quests. It is a measurable vertical slice covering **MQ001–MQ012** in Foundry 9 and Garage 17, ending with a roadworthy starter vehicle and a durable save/reconnect boundary.

## Vertical-slice player promise

A new player can:

```text
create/resume account
→ enter NOVA CITY
→ follow Rowan Gray's message
→ reach Garage 17
→ meet Maya Voss
→ claim the broken DRIVE ZERO prototype
→ salvage/fit required parts
→ discover the Legacy ECU anomaly
→ perform Cold Crank
→ complete First Ignition
→ make the vehicle Roadworthy
→ disconnect/reconnect without duplicated rewards or lost ownership
```

Detailed acceptance criteria live in `docs/vertical-slice-v0.2.md` and machine-readable records in `design/catalog/vertical-slice.json`.

## Vehicle identity contract

Every persistent vehicle has:
- a stable vehicle ID,
- one authoritative owner character,
- immutable/reconstructable build revisions,
- provenance and starter/prototype lineage,
- condition and reputation,
- race/build history references.

Routine deletion of the starter-lineage vehicle is prohibited. Any future exceptional recovery/destruction mechanic must preserve provenance and provide a deterministic restore path.

## Build contract

A build revision is an immutable snapshot of equipped stable part IDs plus validated tuning parameters. Competitive registration binds to an exact accepted revision. Changing the build after registration does not rewrite the accepted snapshot.

## Economy and reward contract

Every rewardable mutation requires an operation/idempotency identifier. Retrying a completed quest, race award, inventory grant, or economy mutation must not duplicate value. The authoritative system records why value changed.

## Racing contract

A ranked result must be attributable to:
- registered event/ruleset,
- authoritative account/character/vehicle ownership,
- accepted vehicle build revision,
- ordered checkpoints/laps,
- timing provenance,
- penalties/disconnect policy,
- a single idempotent award operation.

The v0.2 executable reference proves the build-binding and ordered-checkpoint invariants only. It is not an anti-cheat or production race server.

## VIP fairness contract

Free baseline:
- one primary character,
- one vehicle/garage slot,
- baseline blueprint storage,
- full quest/faction/ranked progression.

VIP may add:
- garage capacity,
- blueprint/storage capacity,
- saved configurations,
- cosmetic/social customization.

VIP must not change the competitive build signature, hidden vehicle performance, race timing, matchmaking priority, reward multiplier, or access to required ranked progression.

## Relationship contract

Relationship state remains authored rather than grind-only:

```text
Stranger → Acquaintance → Friend → Trusted → Partner
                                   ↘ Rival
                                    ↘ Enemy
```

The vertical slice proves the first Maya transition through story-state persistence. Later branches may move Adrian and other NPCs toward ally/friend/rival/enemy paths.

## Faction and endgame continuity

The five launch factions remain NOVA Authority, Iron Wolves, Mechanist Guild, VANTEX, and Free Roads. Endgame governance choices remain:
- destroy DRIVE ZERO,
- transfer it to civic authority,
- return it to VANTEX,
- seize it,
- secret: decentralize it.

These outcomes affect personal/crew overlays, NPC disposition, quest access, and selected world presentation while preserving shared seasonal compatibility.

## Mature-world boundary

The 18+ target supports crime, violence, betrayal, nightlife, alcohol, gambling as story/environmental material, adult relationships, and moral ambiguity. Explicit sexual content is not required. Content still requires rating, localization, safety, and regional legal review before release.

## v0.2 executable evidence

`src/zneondrive/` is a dependency-free **reference contract**, not a selected production stack. Its tests prove:
- free account garage capacity,
- VIP capacity without performance advantage,
- append-only vehicle build revisions,
- optimistic build mutation checks,
- idempotent quest rewards,
- no duplicate rewards under a different retry operation,
- race results bound to accepted build revision and ordered checkpoints,
- starter vehicle routine-deletion protection.

## Definition of v0.2 complete

- v0.1 lore/catalog baseline preserved.
- MQ001–MQ012 vertical slice specified with acceptance criteria.
- first-session, reconnect/save, movement/interaction, accessibility, graybox, and race proof targets documented.
- executable authoritative reference model added.
- reference tests run in CI.
- no claim that engine, networking, persistence database, anti-cheat, HA/DR, or production deployment is complete.
