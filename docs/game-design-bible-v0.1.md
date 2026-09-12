# PROJECT: NEON DRIVE — Game Design Bible v0.1

**Status:** Canonical design baseline  
**Setting:** NOVA CITY, 2097  
**Rating target:** 18+  
**Runtime stance:** Persistent online world / server-client; technology intentionally undecided

## Product pillars

### Vehicle as identity
The starter vehicle is not a disposable mount. It is a persistent identity object carrying ownership, builder, build, repair, race, championship, and reputation history. A vehicle may progress from scrap-tier to Legendary without being replaced.

### Earn the road
Progression should answer what the player earned, built, learned, risked, or changed.

### One city, many lives
NOVA CITY supports racing, delivery work, salvage, mechanical work, faction contracts, exploration, social venues, crews, pets, relationship stories, and seasonal world events.

### Consequences without breaking the shared world
Story decisions matter through personal/crew world-state overlays: NPC disposition, quest access, faction state, selected instanced areas, and ambience. Shared seasonal state remains compatible for all players.

### Fair competition
VIP is capacity/convenience, never hidden competitive power.

## Player fantasy

The player arrives with little money, one phone, and a message directing them to **Garage 17** before midnight. Maya Voss gives them a broken vehicle tied to their past. Rebuilding it opens the city; later the same machine becomes evidence in the conflict around PROJECT DRIVE ZERO.

```text
outsider → builder → street name → contender → faction asset
→ championship candidate → hunted witness → city-scale decision maker
```

## Main cast

The launch roster contains 25 named characters in `design/catalog/characters.json`. Core anchors:

- **Maya Voss** — engineer, former racer, ex-VANTEX DRIVE ZERO engineer.
- **Adrian Cross** — former champion; Rival → Ally → Friend/Trusted or Enemy.
- **Victor Kane** — VANTEX CEO and principal antagonist.
- **Luna** — early companion with search/discovery utility rather than race buffs.

Relationship state model:

```text
Stranger → Acquaintance → Friend → Trusted → Partner
                                   ↘ Rival
                                    ↘ Enemy
```

Major transitions require authored story moments, not only numerical grinding.

## Factions

Five launch factions:

1. **NOVA Authority** — law, continuity, surveillance, legitimacy.
2. **Iron Wolves** — underground racing, loyalty, territorial pride, reputation.
3. **Mechanist Guild** — builders, repair networks, invention, provenance.
4. **VANTEX Corporation** — mobility infrastructure, capital, DRIVE ZERO.
5. **Free Roads** — anti-monopoly mobility activists and decentralized communities.

Faction standing is independent per faction. Some content becomes mutually exclusive; damaged relationships may sometimes be repaired through restitution or service.

## Campaign

The campaign contains exactly 100 stable-ID main quests across seven chapters.

| Chapter | Theme | Transformation |
|---|---|---|
| 1 — The Broken Machine | survival and rebuilding | player becomes mobile |
| 2 — Street Reputation | work and identity | player earns a name |
| 3 — Underground | risk and belonging | player chooses whom to trust |
| 4 — The Championship | public legitimacy | player becomes visible |
| 5 — VANTEX | conspiracy | vehicle becomes evidence |
| 6 — Drive Zero | control vs autonomy | player becomes a target |
| 7 — War for Nova | governance | player determines the future |

The canonical catalog is `design/catalog/quests-main.json`.

## Plot twist

The starter vehicle is a surviving prototype from PROJECT DRIVE ZERO. Its ECU contains a missing control/identity artifact VANTEX has searched for since the project's collapse. The previous owner is tied to the player's family history.

Foreshadowing:
- anomalous ECU behavior,
- legacy hardware identifiers,
- Maya's incomplete explanations,
- VANTEX interest disproportionate to the vehicle's apparent value,
- archived telemetry,
- a hidden diagnostic partition,
- NPC recognition of the chassis lineage.

## Vehicle lifecycle

```text
Chassis
 ├─ Engine / Motor
 ├─ Transmission / Drive unit
 ├─ Suspension
 ├─ Brakes
 ├─ Wheels / Tires
 ├─ ECU / Control
 ├─ Body / Aero
 ├─ Interior / Driver interface
 ├─ Electronics / Sensors
 └─ Special Modules
```

A build is a versioned configuration, not merely a bag of stats. Race results reference the exact build revision used.

The starter vehicle may be reconfigured into Street, Drift, Circuit, Drag, Off-road, Delivery, or Hybrid roles.

## Competition

Launch disciplines:
- Street Race
- Circuit
- Drag
- Drift
- Off-road
- Time Attack
- Delivery Challenge
- Team Race
- Championship

**NOVA GRAND PRIX** is the seasonal pinnacle. Entry requires qualification, never purchase.

Competitive integrity requires validated eligibility, build revision, event ruleset, start state, checkpoints/laps, penalties, disconnect policy, and server-verifiable results.

## Jobs and non-race play

- courier / delivery
- passenger transport
- salvage / parts recovery
- mechanic assistance
- route discovery
- faction contracts
- recovery / towing
- telemetry capture
- convoy escort
- event support

These provide money, XP, parts, blueprints, reputation, faction standing, and relationship changes without forcing continuous racing.

## Pets

Pets offer bounded exploration/support utility and no race-stat boosts. Archetypes include dog, cat, hawk/scout, and cybernetic companion. Important progression must offer non-pet alternatives.

## Crews and social world

Social venues include Garage, Café, Car Meet, Night Club, Beach, Park, Apartment, Workshop, and Race Track.

Crew loop:

```text
meet players → form/join crew → shared jobs/events
→ crew reputation → crew garage/social identity
→ team qualifiers → crew championship → seasonal legacy
```

## Economy and VIP

v0.1 defines only in-world economy concepts: Money, XP, Reputation, Faction Standing, Relationship State, Blueprint Knowledge, and Vehicle Reputation. It intentionally does not define real-money currency.

Free baseline:
- 1 primary character
- 1 starter/active vehicle slot
- baseline blueprint storage
- full competitive progression

VIP envelope:
- 5–20 garage slots depending on future packaging
- additional blueprint/storage capacity
- saved loadouts/configurations
- cosmetic/social customization
- no direct performance multiplier

Any future monetization implementation requires a separate fairness and legal/compliance review.

## Mature themes

18+ supports crime, betrayal, violence, nightlife, alcohol, gambling as narrative/environmental material, adult relationships, and moral ambiguity. Explicit sexual content is not required by the design.

## Endgame choices

- **Destroy DRIVE ZERO** — mobility autonomy.
- **Transfer to civic authority** — regulated control.
- **Return to VANTEX** — corporate order/efficiency.
- **Seize it** — power centralized around the player.
- **Secret: decentralize it** — no single actor can dominate.

Outcomes change NPCs, faction standing, quest availability, access, and selected world-state presentation for that character.

## Post-campaign

The world continues through seasons, qualifiers, championships, crew competition, rotating faction contracts, world events, new districts, relationship chapters, vehicle systems, and expansions.

## Definition of design-complete v0.1

- 25 launch characters cataloged.
- 5 factions and 10 districts cataloged.
- 100 main quests have stable IDs.
- starter vehicle/part taxonomy defined.
- pets, economy constraints, VIP fairness, relationships, crews, and multiplayer loops defined.
- technology-neutral server/client authority boundaries documented.
- runtime systems can reference stable content IDs without renaming core lore.
