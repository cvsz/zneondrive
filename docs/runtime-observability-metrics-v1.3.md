# Runtime Observability Metrics v1.3

**Status:** source/unit evidence implemented; deployed SLO measurement remains open.

## Goal

Add a low-cardinality, credential-safe measurement surface for the Go service plane so later deployment-scale load, soak, incident, and SLO exercises can collect comparable evidence without changing gameplay or durable-state authority.

## Implemented instrumentation

The Go HTTP service records:

- request totals by bounded route scope and HTTP status class,
- request-duration histogram by bounded route scope,
- current in-flight request gauge.

The route label set is intentionally static (`health`, `bootstrap`, `state`, `game-ticket`, `internal-ticket-redeem`, `race-start`, `race-checkpoint`, `race-finish`, `quest-complete`, `build-revise`, `other`). Dynamic quest, vehicle, race-instance, account, session, peer, token, and operation identifiers are never used as metric labels.

The duration histogram uses fixed buckets from 5 ms through 5 seconds and exports Prometheus text exposition without adding a runtime metrics dependency.

## Metrics network boundary

Metrics are served from a separate listener controlled by `METRICS_LISTEN_ADDR`; an empty value disables the listener. Docker Compose enables it on container port 9090 and publishes it only to host loopback (`127.0.0.1:${GAME_API_METRICS_PORT:-19090}`), keeping the metrics endpoint off the player-facing API port.

Production deployments must bind or expose this listener only on a private observability network. This source configuration is not evidence that a deployed firewall, ingress policy, or NetworkPolicy is correct.

## Security and authority invariants

- Metrics contain no Authorization header, session token, resume key, gameplay ticket, game-server key, raw peer IP, vehicle ID, quest ID, race-instance ID, or other unbounded player-controlled label.
- Metrics are observational only and cannot affect PostgreSQL transactions, race acceptance, reward decisions, Redis limiter decisions, or Unreal dedicated-server authority.
- The existing `/healthz` dependency probe remains on the service API; metrics scraping uses the separate listener.

## Automated evidence

Unit tests verify:

- normalized dynamic routes produce bounded labels,
- 2xx/4xx status classes are counted correctly,
- duration histogram count is emitted,
- in-flight requests return to zero after completion,
- raw authorization values and dynamic vehicle/race identifiers do not appear in exported metrics.

The standard Runtime Go workflow remains the merge gate for module-lock cleanliness, unit tests, PostgreSQL/Redis integration tests, `go vet ./...`, and service build.

## Explicit non-claims

This does **not** prove:

- production SLO compliance,
- deployed metrics reachability/isolation,
- PostgreSQL query/pool metrics,
- Redis server metrics,
- Unreal dedicated-server tick/replication metrics,
- alert rules or dashboard correctness,
- deployment-scale load capacity,
- long-duration soak stability,
- HA/DR, backup/restore, or production readiness.

Those require deployed measurements and retained operational evidence.