CREATE TABLE IF NOT EXISTS inventory_items (
    character_id TEXT NOT NULL REFERENCES characters(id) ON DELETE CASCADE,
    item_id TEXT NOT NULL,
    quantity BIGINT NOT NULL DEFAULT 0 CHECK (quantity >= 0),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (character_id, item_id)
);

CREATE TABLE IF NOT EXISTS character_blueprints (
    character_id TEXT NOT NULL REFERENCES characters(id) ON DELETE CASCADE,
    blueprint_id TEXT NOT NULL,
    unlocked_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (character_id, blueprint_id)
);

CREATE INDEX IF NOT EXISTS ix_inventory_items_character
    ON inventory_items(character_id);

CREATE INDEX IF NOT EXISTS ix_character_blueprints_character
    ON character_blueprints(character_id);
