# Release Readiness Checklist

Use this checklist for release/go-no-go reviews. Check only items backed by evidence.

## Source and CI

- [ ] target commit/tag identified
- [ ] repository CI green
- [ ] CodeQL/security scans reviewed
- [ ] dependency review clean or exceptions documented
- [ ] Go unit/integration tests pass
- [ ] Unreal target build passes on retained toolchain evidence
- [ ] packaged client/server artifacts retained
- [ ] UE 5.8.2 archive integrity and install evidence retained
- [ ] Windows 11 client package launch/API test passed
- [ ] Android APK/AAB install/launch/API test passed

## Gameplay/content

- [ ] milestone quests playable end-to-end
- [ ] save/reconnect verified
- [ ] vehicle build/rebuild verified
- [ ] race authority/result flow verified
- [ ] localization coverage reviewed
- [ ] accessibility review completed
- [ ] known content blockers documented

## Security

- [ ] threat-model tests executed
- [ ] TLS/service authentication deployed
- [ ] secrets externalized and rotation tested
- [ ] rate limits verified
- [ ] privileged actions audited
- [ ] ranked anti-cheat/integrity evidence reviewed

## Reliability

- [ ] load test passed
- [ ] soak test passed
- [ ] capacity limits understood
- [ ] backup current
- [ ] restore drill passed
- [ ] RPO/RTO accepted
- [ ] rollback tested

## Operations

- [ ] dashboards ready
- [ ] alerts tested
- [ ] incident owner/on-call defined
- [ ] support/moderation runbooks ready
- [ ] live-ops controls tested
- [ ] release communication prepared

## Legal/business

- [ ] IP/assets cleared
- [ ] privacy/legal review completed
- [ ] age-rating/platform requirements addressed
- [ ] monetization/store behavior approved

## Decision

- [ ] GO
- [ ] NO-GO
- [ ] CONDITIONAL GO with explicitly accepted risks

Record approver, date, commit/build IDs and evidence links with the decision.
