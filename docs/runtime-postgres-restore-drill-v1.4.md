# Runtime PostgreSQL Restore Drill v1.4

**Evidence level:** CI/source-level isolated restore evidence. This is not production DR, HA, off-host retention, encryption-at-rest, or demonstrated production RPO/RTO evidence.

## Purpose

`tools/postgres-restore-drill.sh` exercises the durable PostgreSQL recovery path with the same PostgreSQL 17 major version selected by the runtime. The drill intentionally excludes Redis because Redis is ephemeral coordination state and must not become the only copy of durable player state.

## What the drill proves

The Runtime Go workflow creates two isolated throwaway databases. It applies every canonical migration to the source database, seeds representative authoritative state, creates a custom-format `pg_dump`, restores that archive into the second database with `pg_restore --exit-on-error`, and validates the restored state.

The retained assertions cover:

- account identity and character progression,
- starter vehicle and active immutable build evidence,
- build operation-ID/idempotency evidence,
- quest completion operation-ID/idempotency evidence,
- inventory quantity,
- blueprint unlock,
- finished authoritative race result and operation ID,
- safe replay of every canonical migration after restore.

The workflow uploads the dump plus `report.txt` as a 30-day GitHub Actions artifact. The report records dump size, elapsed drill time, and the validated state classes without storing credentials.

## Trust boundary

The drill does not alter the gameplay authority model. PostgreSQL remains the durable state authority, Unreal dedicated servers remain gameplay authority, and Redis remains disposable coordination state. The script operates only on uniquely named throwaway databases and drops them on exit.

## Non-claims

Passing this drill does **not** establish:

- production backup scheduling or retention,
- encrypted/off-host backup custody,
- point-in-time recovery or a 15-minute production RPO,
- a two-hour production RTO,
- restore behavior against production data volume,
- cross-node, multi-region, provider-loss, or HA failover recovery,
- production incident/rollback exercise completion.

Those remain separate production gates in `ROADMAP.md` and `IMPLEMENTATION-CHECKLIST.md`.

## Local execution

With a PostgreSQL 17 service reachable on localhost and Docker available:

```bash
DATABASE_ADMIN_URL='postgres://zneondrive:zneondrive-test@127.0.0.1:5432/postgres?sslmode=disable' \
  bash tools/postgres-restore-drill.sh
```

Evidence is written under `.runtime/backup-restore-drill/`.
