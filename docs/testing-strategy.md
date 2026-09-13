# Testing Strategy

## Principle

Tests are evidence for specific claims. A workflow file, test stub or design document is not evidence that runtime behavior passed.

## Current automated layers

### Design validation
`tools/validate_design.py`
- stable counts/IDs,
- quest sequence/references,
- catalog integrity,
- pet competitive-buff constraints,
- schema presence/parseability.

### Client presentation validation
`tools/validate_client_demo.py`
- required demo sections,
- offline-first assets,
- no placeholder markers,
- MQ001–MQ012 presentation coverage.

### Python authority oracle
`tests/test_reference_runtime.py`
- ownership,
- vehicle build revision behavior,
- idempotent rewards,
- race/build binding,
- starter-vehicle protection,
- VIP fairness invariants.

### Runtime structure validators
- `validate_runtime_v0_4.py`
- `validate_runtime_v0_5.py`
- `validate_runtime_v0_6.py`

These verify expected source/contracts exist; they do not substitute for a real Unreal build.

### Go unit/integration
The Go service uses unit tests plus PostgreSQL-tagged integration tests for durable behavior.

## Required future layers

### Unreal automated tests
- C++ unit/automation tests,
- replication tests,
- gameplay ticket binding,
- build snapshot application,
- reconnect behavior.

### End-to-end
A packaged client + dedicated server + Go + PostgreSQL path should prove:
- fresh bootstrap,
- Garage 17 rebuild,
- disconnect/reconnect,
- First Ignition,
- race registration/checkpoints/result,
- reward idempotency.

### Security
- replay,
- forged ownership,
- stale revision,
- impossible build,
- ticket reuse,
- session abuse,
- rate-limit behavior.

### Performance
- dedicated server frame/time budget,
- service latency,
- database saturation,
- replication/bandwidth,
- soak stability.

### Recovery
- backup,
- restore,
- process/node loss,
- duplicate mutation safety.

## Merge policy

A PR should run all tests applicable to the files it changes. Failures should be fixed rather than bypassed. Evidence-gated checklist items remain open until the exact required test has passed.
