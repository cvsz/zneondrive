CREATE TABLE IF NOT EXISTS race_instances (
    id TEXT PRIMARY KEY,
    race_id TEXT NOT NULL,
    account_id TEXT NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
    character_id TEXT NOT NULL REFERENCES characters(id) ON DELETE CASCADE,
    vehicle_id TEXT NOT NULL REFERENCES vehicles(id) ON DELETE CASCADE,
    build_revision INTEGER NOT NULL CHECK (build_revision >= 1),
    build_validation_hash TEXT NOT NULL,
    state TEXT NOT NULL DEFAULT 'active' CHECK (state IN ('active','finished','abandoned')),
    operation_id TEXT NOT NULL UNIQUE,
    next_checkpoint INTEGER NOT NULL DEFAULT 0 CHECK (next_checkpoint >= 0),
    last_elapsed_ms BIGINT NOT NULL DEFAULT 0 CHECK (last_elapsed_ms >= 0),
    started_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    finished_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS ix_race_instances_account_started
    ON race_instances(account_id, started_at DESC);
CREATE INDEX IF NOT EXISTS ix_race_instances_vehicle_started
    ON race_instances(vehicle_id, started_at DESC);

CREATE TABLE IF NOT EXISTS race_checkpoints (
    race_instance_id TEXT NOT NULL REFERENCES race_instances(id) ON DELETE CASCADE,
    checkpoint_index INTEGER NOT NULL CHECK (checkpoint_index >= 0),
    elapsed_ms BIGINT NOT NULL CHECK (elapsed_ms > 0),
    operation_id TEXT NOT NULL UNIQUE,
    recorded_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (race_instance_id, checkpoint_index)
);

CREATE TABLE IF NOT EXISTS race_results (
    race_instance_id TEXT PRIMARY KEY REFERENCES race_instances(id) ON DELETE CASCADE,
    checkpoint_count INTEGER NOT NULL CHECK (checkpoint_count > 0),
    finish_elapsed_ms BIGINT NOT NULL CHECK (finish_elapsed_ms > 0),
    result_hash TEXT NOT NULL,
    operation_id TEXT NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
