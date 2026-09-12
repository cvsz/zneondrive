# Implementation Checklist — PROJECT: NEON DRIVE

## Design baseline
- [x] Story premise and NOVA CITY 2097
- [x] Garage 17 opening and starter prototype
- [x] Maya Voss / Adrian Cross / Victor Kane / Luna anchors
- [x] PROJECT DRIVE ZERO conflict
- [x] Seven campaign chapters
- [x] Five endgame governance paths
- [x] VIP fairness principle
- [x] Social/crew and pet design
- [x] 25-character launch catalog
- [x] 100-main-quest graph
- [x] World, faction, district, and part catalogs
- [x] Automated catalog integrity checks
- [x] Game Design Bible v0.2 canonical implementation rules

## Vertical-slice contract
- [x] MQ001–MQ012 acceptance matrix
- [x] Garage 17 / Foundry 9 graybox scope
- [x] First-session UX path
- [x] Starter handling target envelope
- [x] Movement/interaction minimum
- [x] Legal + underground race proof definitions
- [x] Save/reconnect semantics
- [x] Accessibility baseline
- [x] Machine-readable slice catalog and schema

## Executable reference authority
- [x] One account → one primary character reference invariant
- [x] Vehicle ownership authoritative reference invariant
- [x] Append-only build revisions with optimistic conflict checks
- [x] Quest reward idempotency
- [x] Race result bound to accepted build revision
- [x] Ordered checkpoint validation
- [x] VIP capacity separated from competitive build signature
- [x] Starter vehicle routine-deletion protection
- [x] Unit tests for reference invariants

## Content production
- [ ] Full dialogue/script for MQ001–MQ100
- [ ] Side-quest narratives beyond seed records
- [ ] Cinematic list and storyboard requirements
- [ ] Environment storytelling asset list
- [ ] Voice/localization style guides
- [ ] Content sensitivity and age-rating review

## Playable runtime vertical slice
- [ ] Selected client engine ADR
- [ ] Selected production server/runtime ADR
- [ ] Garage 17 environment playable
- [ ] First vehicle rebuild playable from broken state
- [ ] First Ignition end-to-end in selected runtime
- [ ] Foundry 9 traversal and jobs playable
- [ ] One legal and one underground race playable
- [ ] Maya relationship state transition visible in client
- [ ] Luna discovery utility playable
- [ ] Durable save/reconnect verified against selected persistence

## Multiplayer and security
- [ ] Threat model complete
- [ ] Client trust boundaries tested over real transport
- [ ] Rate limiting
- [ ] Replay/idempotency protection in selected persistence/runtime
- [ ] Cheat telemetry
- [ ] Ranked impossible-state detection
- [ ] Privileged admin/live-ops audit trail
- [ ] Abuse/moderation model

## Reliability
- [ ] Load test target defined
- [ ] Soak test passed
- [ ] Backup verified
- [ ] Restore drill passed
- [ ] RPO/RTO defined
- [ ] Reconnect/recovery tested under process/node failure
- [ ] Duplicate rewards prevented on retry/failover
- [ ] Content rollback tested

## Release evidence
- [ ] Stack-specific CI green
- [ ] Unit/integration/e2e suites for selected runtime
- [ ] Security scans pass
- [ ] Performance budgets pass
- [ ] Accessibility review
- [ ] Privacy/data-retention review
- [ ] Production deployment evidence
- [ ] DR evidence
- [ ] Go/no-go approval

Unchecked production claims must not be represented as complete without evidence.
