# Economy & Monetization

## Fairness contract

PROJECT: NEON DRIVE uses a non-pay-to-win design constraint.

VIP may provide:
- garage capacity,
- storage/blueprint capacity,
- saved build convenience,
- cosmetics/social presentation,
- other convenience that does not change ranked performance.

VIP must not provide hidden:
- horsepower,
- grip,
- durability,
- ranked reward multiplier,
- matchmaking priority,
- checkpoint/timing advantage,
- exclusive required progression power.

## Core economy categories

Planned:
- money/credits,
- parts/materials,
- blueprints/knowledge,
- repair/maintenance sinks,
- event/quest rewards,
- reputation-gated opportunities.

Current runtime implements only a subset: money/xp/reputation fields, inventory quantities, blueprint unlocks and quest-linked grants.

## Economy mutation rules

- server/service authoritative,
- idempotent for retryable grants,
- reason/context attributable,
- auditable for privileged/manual adjustments,
- no negative inventory/balance where disallowed,
- rollback/correction through traceable transactions where practical.

## Anti-exploit requirements

Test:
- duplicate quest completion,
- reconnect reward replay,
- stale build mutation,
- item duplication via swap/return,
- VIP entitlement spoofing,
- race reward replay,
- admin grant misuse.

## Commercial integrity

Pricing, store SKU design, tax/payment processing, refund handling and platform-commerce integration are not yet implemented and should not be inferred from VIP design documentation.
