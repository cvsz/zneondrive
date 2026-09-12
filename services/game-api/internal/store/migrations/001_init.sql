CREATE TABLE IF NOT EXISTS accounts (
    id TEXT PRIMARY KEY,
    resume_key_hash TEXT NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS characters (
    id TEXT PRIMARY KEY,
    account_id TEXT NOT NULL UNIQUE REFERENCES accounts(id) ON DELETE CASCADE,
    money BIGINT NOT NULL DEFAULT 0 CHECK (money >= 0),
    xp BIGINT NOT NULL DEFAULT 0 CHECK (xp >= 0),
    reputation BIGINT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS vehicles (
    id TEXT PRIMARY KEY,
    owner_character_id TEXT NOT NULL REFERENCES characters(id) ON DELETE CASCADE,
    active_revision INTEGER NOT NULL DEFAULT 1 CHECK (active_revision >= 1),
    starter_lineage BOOLEAN NOT NULL DEFAULT false,
    roadworthy BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX IF NOT EXISTS ux_vehicle_starter_per_character
    ON vehicles(owner_character_id)
    WHERE starter_lineage = true;

CREATE TABLE IF NOT EXISTS vehicle_builds (
    vehicle_id TEXT NOT NULL REFERENCES vehicles(id) ON DELETE CASCADE,
    revision INTEGER NOT NULL CHECK (revision >= 1),
    part_ids JSONB NOT NULL,
    validation_hash TEXT NOT NULL,
    operation_id TEXT UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (vehicle_id, revision)
);

CREATE TABLE IF NOT EXISTS sessions (
    token_hash TEXT PRIMARY KEY,
    account_id TEXT NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS ix_sessions_account ON sessions(account_id);
CREATE INDEX IF NOT EXISTS ix_sessions_expiry ON sessions(expires_at);

CREATE TABLE IF NOT EXISTS quest_completions (
    character_id TEXT NOT NULL REFERENCES characters(id) ON DELETE CASCADE,
    quest_id TEXT NOT NULL,
    operation_id TEXT NOT NULL UNIQUE,
    reward_money BIGINT NOT NULL DEFAULT 0,
    reward_xp BIGINT NOT NULL DEFAULT 0,
    reward_reputation BIGINT NOT NULL DEFAULT 0,
    completed_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (character_id, quest_id)
);
