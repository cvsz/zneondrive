CREATE TABLE IF NOT EXISTS game_tickets (
    ticket_hash TEXT PRIMARY KEY,
    account_id TEXT NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
    expires_at TIMESTAMPTZ NOT NULL,
    consumed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS ix_game_tickets_account ON game_tickets(account_id);
CREATE INDEX IF NOT EXISTS ix_game_tickets_expiry ON game_tickets(expires_at);
