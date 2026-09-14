#!/bin/sh

: "${POSTGRES_EXPORTER_PASSWORD:?POSTGRES_EXPORTER_PASSWORD must be set}"

psql -v ON_ERROR_STOP=1 \
  --username "$POSTGRES_USER" \
  --dbname "$POSTGRES_DB" \
  --set=monitor_password="$POSTGRES_EXPORTER_PASSWORD" <<'SQL'
CREATE EXTENSION IF NOT EXISTS pg_stat_statements;
DO $$
BEGIN
  IF NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'zneondrive_monitor') THEN
    CREATE ROLE zneondrive_monitor LOGIN;
  END IF;
END
$$;
ALTER ROLE zneondrive_monitor NOSUPERUSER NOCREATEDB NOCREATEROLE NOREPLICATION NOBYPASSRLS;
ALTER ROLE zneondrive_monitor PASSWORD :'monitor_password';
GRANT pg_monitor TO zneondrive_monitor;
SQL
