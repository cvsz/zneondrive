package store

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"time"

	"github.com/cvsz/zneondrive/services/game-api/internal/core"
	"github.com/jackc/pgx/v5"
)

//go:embed migrations/002_game_tickets.sql
var ticketMigrationSQL string

func (p *Postgres) IssueGameTicket(
	ctx context.Context,
	sessionTokenHash string,
	ticketHash string,
	expiresAt time.Time,
) error {
	var accountID string
	err := p.pool.QueryRow(ctx, `
		SELECT account_id
		FROM sessions
		WHERE token_hash=$1 AND expires_at > now()
	`, sessionTokenHash).Scan(&accountID)
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrUnauthorized
	}
	if err != nil {
		return fmt.Errorf("authorize gameplay ticket: %w", err)
	}

	if _, err = p.pool.Exec(ctx, `
		INSERT INTO game_tickets(ticket_hash,account_id,expires_at)
		VALUES($1,$2,$3)
	`, ticketHash, accountID, expiresAt.UTC()); err != nil {
		return fmt.Errorf("insert gameplay ticket: %w", err)
	}
	return nil
}

func (p *Postgres) RedeemGameTicket(ctx context.Context, ticketHash string) (core.Snapshot, error) {
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return core.Snapshot{}, fmt.Errorf("begin gameplay ticket redemption: %w", err)
	}
	defer tx.Rollback(ctx)

	var accountID string
	err = tx.QueryRow(ctx, `
		UPDATE game_tickets
		SET consumed_at=now()
		WHERE ticket_hash=$1
		  AND consumed_at IS NULL
		  AND expires_at > now()
		RETURNING account_id
	`, ticketHash).Scan(&accountID)
	if errors.Is(err, pgx.ErrNoRows) {
		return core.Snapshot{}, ErrUnauthorized
	}
	if err != nil {
		return core.Snapshot{}, fmt.Errorf("redeem gameplay ticket: %w", err)
	}
	if err = tx.Commit(ctx); err != nil {
		return core.Snapshot{}, fmt.Errorf("commit gameplay ticket redemption: %w", err)
	}
	return p.snapshotByAccount(ctx, accountID)
}
