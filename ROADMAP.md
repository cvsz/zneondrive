# zNeonDrive Roadmap

## Phase 0 — Vision and repository identity
- [x] Original Story Bible / Master Plot captured
- [x] Repository reframed from generic template to zNeonDrive
- [x] Security/contribution/release repository baseline retained

## Phase 1 — Game Design Bible v0.1
- [x] Product pillars and fairness rules
- [x] 25 named launch characters
- [x] 5 primary factions
- [x] 10 launch districts
- [x] 7-chapter / 100-main-quest campaign catalog
- [x] Side-content seed catalog
- [x] Vehicle/parts taxonomy
- [x] Pet archetypes
- [x] Economy, relationships, crews, and VIP constraints
- [x] Five endgame governance outcomes

## Phase 2 — Content contracts
- [x] Stable ID conventions
- [x] Quest schema
- [x] Vehicle schema
- [x] Content integrity validator
- [x] CI design validation
- [x] Technology-neutral authority boundaries
- [x] ADR: server-authoritative persistent/competitive state
- [x] ADR: preserve starter-vehicle identity

## Phase 3 — Vertical-slice specification
- [ ] Lock Garage 17 + Foundry 9 + first race route
- [ ] Expand MQ001–MQ012 to full acceptance criteria/dialogue beats
- [ ] Define starter vehicle handling/physical target
- [ ] Define player movement and interaction specification
- [ ] Define first-session UX from login to First Ignition
- [ ] Define reconnect/save semantics
- [ ] Define accessibility baseline
- [ ] Create graybox/content production checklist

## Phase 4 — Runtime prototype
- [ ] Select engine/client technology by ADR
- [ ] Select server/runtime technology by ADR
- [ ] Implement identity/session prototype
- [ ] Implement one-character persistence
- [ ] Implement one-vehicle persistence and build revisions
- [ ] Implement inventory/blueprint prototype
- [ ] Implement quest-state prototype
- [ ] Implement authoritative race prototype
- [ ] Implement relationship/faction-state prototype
- [ ] Automated contract/integration tests

## Phase 5 — Multiplayer alpha
- [ ] Multi-player district instance/shard
- [ ] Presence/social
- [ ] Crew MVP
- [ ] Match/event registration
- [ ] Race anti-cheat telemetry
- [ ] Disconnect/rejoin behavior
- [ ] Load and soak testing
- [ ] Backup/restore drill

## Phase 6 — Content alpha
- [ ] Chapter 1 fully playable
- [ ] Street jobs and early side quests
- [ ] Luna companion questline
- [ ] Foundry 9 and Neon Mile content
- [ ] Early faction reputation
- [ ] First championship qualifier
- [ ] Localization pipeline

## Phase 7 — Production hardening
- [ ] Threat model
- [ ] Abuse/cheat model
- [ ] SLOs and observability
- [ ] Incident/rollback runbooks
- [ ] Capacity model
- [ ] Restore evidence
- [ ] Security review
- [ ] Privacy/data-retention review
- [ ] VIP fairness test suite
- [ ] Ranked race integrity test suite

## Phase 8 — Launch readiness
- [ ] 100 main quests implemented and verified in runtime
- [ ] All launch districts implemented
- [ ] NOVA GRAND PRIX seasonal operations
- [ ] Crew championship
- [ ] Live-ops tools and audit controls
- [ ] Customer support/moderation workflows
- [ ] Production deployment evidence
- [ ] Disaster-recovery evidence
- [ ] Go/no-go release review

**Rule:** do not mark runtime, production, HA, DR, anti-cheat, or live-ops items complete without executable evidence.
