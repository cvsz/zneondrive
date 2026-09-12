# Crew & Social Systems

## Goal

Crews turn individual vehicle identity into persistent social identity without allowing social privilege to bypass competitive fairness.

## Crew MVP target

A crew should eventually support:
- stable crew ID/name,
- owner/admin/member roles,
- membership lifecycle,
- crew reputation,
- shared presentation/garage identity,
- event/team registration,
- moderation/audit hooks.

## Authority

Crew roles, membership and competitive eligibility are durable service-plane state. Client UI may request changes but cannot authorize itself.

## Permission model target

At minimum:
- Owner — destructive/high-impact administration,
- Admin — membership/event administration,
- Member — normal participation.

Role changes must be authorized and audited.

## Competitive rules

Crew membership must not grant hidden vehicle performance. Team events may use event-specific composition/rules but eligibility remains server validated.

## Social safety

Before public social release:
- block/mute/report,
- name/content validation,
- spam/rate controls,
- moderator tooling,
- sanction/appeal flow.

## Evidence

Crew MVP is not complete until concurrent players can form/use a crew, reconnect safely, and enter a team event with persisted role/reputation state.
