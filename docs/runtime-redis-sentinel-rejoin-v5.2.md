# Runtime Redis Sentinel Rejoin Fencing Evidence v5.2

## Scope

This increment extends the existing Redis Sentinel CI topology with explicit **old-primary rejoin** evidence for the Go 1.27 distributed HTTP abuse-control limiter while preserving the selected Unreal Engine 5.8 + Go 1.27 + PostgreSQL + Redis authority model.

Redis remains **ephemeral coordination only**. PostgreSQL remains the durable account/progression/inventory/build/quest/race authority, and Unreal dedicated servers remain gameplay authority. This increment does not move durable state or gameplay authority into Redis or the Go process.

## Executable evidence

`services/game-api/internal/httpapi/redis_sentinel_rejoin_integration_test.go` creates an isolated Redis 8 topology with one primary, one replica and three Sentinel processes using quorum 2. The test:

1. waits for Redis replication plus Sentinel quorum/replica discovery;
2. exhausts a canonical shared limiter bucket through the Sentinel-backed limiter;
3. verifies the exhausted state has replicated to the replica;
4. stops the old primary and waits for Sentinel to promote the replica;
5. verifies the existing application limiter resolves only the promoted authority;
6. restarts the old primary without issuing an application/test-side `REPLICAOF` command;
7. waits for Sentinel to reconfigure the old primary as a connected replica of the promoted authority;
8. verifies a direct `SET` against the rejoined node is rejected with Redis `READONLY` behavior;
9. verifies Sentinel discovery does not regress to the old node;
10. verifies the limiter state is replicated back to the rejoined node and the promoted authority still rejects the already-exhausted budget.

This is stronger than v5.1 because it exercises the failed node's return path and proves that, in this isolated Sentinel topology, the returned node is demoted to a read-only replica instead of becoming a competing writer.

## Security and trust boundary

- Sentinel addresses/master name are operator-controlled configuration, never player input.
- Redis coordinates only ephemeral abuse-control state; durable gameplay data remains PostgreSQL-authoritative.
- The rejoined old node is considered safe in this evidence only after Redis reports replica role, the replication link points to the promoted authority, and direct writes are rejected.
- Application master discovery must remain pinned to the Sentinel-reported promoted authority.
- Limiter identities remain hashed; the test stores no bearer/session credentials or player-controlled durable state in Redis.

## Explicit non-claims

This CI-scale evidence does **not** prove:

- safety during an asymmetric network partition where the old primary remains reachable and writable while quorum promotes another node;
- external fencing/STONITH or independent lease/consensus fencing;
- Redis Cluster behavior;
- zero-loss asynchronous replication for limiter writes not received before failure;
- deployment-scale multi-node failover under sustained production traffic;
- long-duration soak stability;
- Kubernetes node loss/rescheduling;
- production Redis authentication/TLS/ACL/secret-rotation posture;
- accepted production RPO/RTO or regional DR;
- live packaged Unreal reconnect/race behavior;
- production readiness.

Accordingly, the broader **Redis split-brain prevention / deployment HA** gate remains open. Repository status remains **Evidence Gated / Not Production Ready** until the real-Unreal, deployment, load/soak, observability, security and HA/DR gates are demonstrated with retained evidence.
