#!/usr/bin/env bash
set -Eeuo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
CLIENT_IMAGE="${POSTGRES_CLIENT_IMAGE:-postgres:17-alpine}"
ADMIN_URL="${DATABASE_ADMIN_URL:-postgres://zneondrive:zneondrive-test@127.0.0.1:5432/postgres?sslmode=disable}"
SOURCE_DB="${BACKUP_DRILL_SOURCE_DB:-zneondrive_backup_drill}"
RESTORE_DB="${BACKUP_DRILL_RESTORE_DB:-zneondrive_restore_drill}"
ARTIFACT_DIR="${BACKUP_DRILL_ARTIFACT_DIR:-$ROOT/.runtime/backup-restore-drill}"
DUMP_FILE="$ARTIFACT_DIR/zneondrive-backup.dump"
REPORT_FILE="$ARTIFACT_DIR/report.txt"
MIGRATIONS_DIR="$ROOT/services/game-api/internal/store/migrations"

command -v docker >/dev/null 2>&1 || { echo 'ERROR: docker is required' >&2; exit 1; }
docker info >/dev/null 2>&1 || { echo 'ERROR: docker daemon is not reachable' >&2; exit 1; }
[[ -d "$MIGRATIONS_DIR" ]] || { echo "ERROR: migrations directory missing: $MIGRATIONS_DIR" >&2; exit 1; }
mkdir -p "$ARTIFACT_DIR"
rm -f "$DUMP_FILE" "$REPORT_FILE"

url_for_db() {
  python3 - "$ADMIN_URL" "$1" <<'PY'
from urllib.parse import urlsplit, urlunsplit
import sys
u = urlsplit(sys.argv[1])
db = sys.argv[2]
print(urlunsplit((u.scheme, u.netloc, '/' + db, u.query, u.fragment)))
PY
}

SOURCE_URL="$(url_for_db "$SOURCE_DB")"
RESTORE_URL="$(url_for_db "$RESTORE_DB")"

pg() {
  docker run --rm --network host \
    -v "$ROOT:/repo:ro" \
    -v "$ARTIFACT_DIR:/evidence" \
    "$CLIENT_IMAGE" "$@"
}

psql_admin() {
  pg psql "$ADMIN_URL" -v ON_ERROR_STOP=1 "$@"
}

cleanup() {
  set +e
  psql_admin -c "SELECT pg_terminate_backend(pid) FROM pg_stat_activity WHERE datname IN ('$SOURCE_DB', '$RESTORE_DB') AND pid <> pg_backend_pid();" >/dev/null 2>&1
  psql_admin -c "DROP DATABASE IF EXISTS \"$SOURCE_DB\";" >/dev/null 2>&1
  psql_admin -c "DROP DATABASE IF EXISTS \"$RESTORE_DB\";" >/dev/null 2>&1
}
trap cleanup EXIT

start_epoch="$(date +%s)"

# The drill always uses isolated throwaway databases. It never dumps or restores the
# CI application's working database, and it never touches Redis because Redis is not
# authoritative durable state in zNeonDrive.
cleanup
psql_admin -c "CREATE DATABASE \"$SOURCE_DB\";"

for migration in "$MIGRATIONS_DIR"/*.sql; do
  name="$(basename "$migration")"
  pg psql "$SOURCE_URL" -v ON_ERROR_STOP=1 -f "/repo/services/game-api/internal/store/migrations/$name"
done

seed_sql=$(cat <<'SQL'
INSERT INTO accounts (id, resume_key_hash) VALUES ('drill-account', 'drill-resume-hash');
INSERT INTO characters (id, account_id, money, xp, reputation)
VALUES ('drill-character', 'drill-account', 1250, 420, 7);
INSERT INTO vehicles (id, owner_character_id, active_revision, starter_lineage, roadworthy)
VALUES ('drill-vehicle', 'drill-character', 1, true, true);
INSERT INTO vehicle_builds (vehicle_id, revision, part_ids, validation_hash, operation_id)
VALUES ('drill-vehicle', 1, '["part_chassis_starter_prototype","part_engine_ice_street_i"]'::jsonb, 'drill-build-hash', 'drill-build-op');
INSERT INTO quest_completions (character_id, quest_id, operation_id, reward_money, reward_xp, reward_reputation)
VALUES ('drill-character', 'MQ001', 'drill-quest-op', 250, 100, 2);
INSERT INTO inventory_items (character_id, item_id, quantity)
VALUES ('drill-character', 'part_brakes_track_i', 2);
INSERT INTO character_blueprints (character_id, blueprint_id)
VALUES ('drill-character', 'blueprint_starter_rebuild');
INSERT INTO race_instances (
  id, race_id, account_id, character_id, vehicle_id, build_revision,
  build_validation_hash, state, operation_id, next_checkpoint, last_elapsed_ms
) VALUES (
  'drill-race-instance', 'race-garage17-shakedown', 'drill-account',
  'drill-character', 'drill-vehicle', 1, 'drill-build-hash', 'finished',
  'drill-race-start-op', 2, 9000
);
INSERT INTO race_checkpoints (race_instance_id, checkpoint_index, elapsed_ms, operation_id)
VALUES
  ('drill-race-instance', 0, 4000, 'drill-race-cp0-op'),
  ('drill-race-instance', 1, 9000, 'drill-race-cp1-op');
INSERT INTO race_results (race_instance_id, checkpoint_count, finish_elapsed_ms, result_hash, operation_id)
VALUES ('drill-race-instance', 2, 12000, 'drill-result-hash', 'drill-race-finish-op');
SQL
)
printf '%s\n' "$seed_sql" | pg psql "$SOURCE_URL" -v ON_ERROR_STOP=1

# Custom-format backup gives pg_restore a real archive to validate rather than merely
# replaying SQL text.
pg pg_dump --format=custom --no-owner --no-acl --file=/evidence/zneondrive-backup.dump "$SOURCE_URL"
[[ -s "$DUMP_FILE" ]] || { echo 'ERROR: pg_dump produced no backup artifact' >&2; exit 1; }

psql_admin -c "CREATE DATABASE \"$RESTORE_DB\";"
pg pg_restore --exit-on-error --no-owner --no-acl --dbname="$RESTORE_URL" /evidence/zneondrive-backup.dump

validation_sql=$(cat <<'SQL'
DO $$
DECLARE
  failures text[] := ARRAY[]::text[];
BEGIN
  IF (SELECT count(*) FROM accounts WHERE id='drill-account') <> 1 THEN failures := array_append(failures, 'account'); END IF;
  IF (SELECT count(*) FROM characters WHERE id='drill-character' AND money=1250 AND xp=420 AND reputation=7) <> 1 THEN failures := array_append(failures, 'character progression'); END IF;
  IF (SELECT count(*) FROM vehicles WHERE id='drill-vehicle' AND active_revision=1 AND starter_lineage AND roadworthy) <> 1 THEN failures := array_append(failures, 'vehicle'); END IF;
  IF (SELECT count(*) FROM vehicle_builds WHERE vehicle_id='drill-vehicle' AND revision=1 AND validation_hash='drill-build-hash' AND operation_id='drill-build-op') <> 1 THEN failures := array_append(failures, 'build/idempotency'); END IF;
  IF (SELECT count(*) FROM quest_completions WHERE character_id='drill-character' AND quest_id='MQ001' AND operation_id='drill-quest-op') <> 1 THEN failures := array_append(failures, 'quest/idempotency'); END IF;
  IF (SELECT count(*) FROM inventory_items WHERE character_id='drill-character' AND item_id='part_brakes_track_i' AND quantity=2) <> 1 THEN failures := array_append(failures, 'inventory'); END IF;
  IF (SELECT count(*) FROM character_blueprints WHERE character_id='drill-character' AND blueprint_id='blueprint_starter_rebuild') <> 1 THEN failures := array_append(failures, 'blueprint'); END IF;
  IF (SELECT count(*) FROM race_results r JOIN race_instances i ON i.id=r.race_instance_id WHERE i.id='drill-race-instance' AND i.state='finished' AND r.checkpoint_count=2 AND r.finish_elapsed_ms=12000 AND r.operation_id='drill-race-finish-op') <> 1 THEN failures := array_append(failures, 'race result/idempotency'); END IF;
  IF coalesce(array_length(failures, 1), 0) <> 0 THEN
    RAISE EXCEPTION 'restore validation failed: %', array_to_string(failures, ', ');
  END IF;
END $$;
SQL
)
printf '%s\n' "$validation_sql" | pg psql "$RESTORE_URL" -v ON_ERROR_STOP=1

# Schema/application compatibility: every canonical migration must remain safely
# replayable after restore.
for migration in "$MIGRATIONS_DIR"/*.sql; do
  name="$(basename "$migration")"
  pg psql "$RESTORE_URL" -v ON_ERROR_STOP=1 -f "/repo/services/game-api/internal/store/migrations/$name" >/dev/null
done

end_epoch="$(date +%s)"
duration="$((end_epoch - start_epoch))"
dump_bytes="$(wc -c < "$DUMP_FILE" | tr -d ' ')"
{
  echo 'zNeonDrive PostgreSQL isolated restore drill: PASS'
  echo "postgres_client_image=$CLIENT_IMAGE"
  echo "source_database=$SOURCE_DB"
  echo "restore_database=$RESTORE_DB"
  echo "dump_bytes=$dump_bytes"
  echo "duration_seconds=$duration"
  echo 'validated=account,character_progression,vehicle,build_operation_id,quest_operation_id,inventory,blueprint,race_result_operation_id,migration_replay'
  echo 'scope=CI/source-level restore evidence; not production RPO/RTO, off-host retention, encryption, HA, or regional DR evidence'
} | tee "$REPORT_FILE"
