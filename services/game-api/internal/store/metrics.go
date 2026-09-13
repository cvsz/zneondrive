package store

import (
	"context"
	"time"
)

// PostgresPoolStats is a bounded, credential-free snapshot of the pgx pool.
// It intentionally excludes SQL text, connection strings, database names,
// account identifiers, and any player-controlled values so it is safe to
// expose on the existing private metrics listener.
type PostgresPoolStats struct {
	MaxConns             int32
	TotalConns           int32
	IdleConns            int32
	AcquiredConns        int32
	ConstructingConns    int32
	AcquireCount         int64
	EmptyAcquireCount    int64
	CanceledAcquireCount int64
	NewConnsCount        int64
	AcquireDuration      time.Duration
}

// PostgresServerStats is a bounded snapshot of current-database pg_stat_database
// counters plus the current database size. It deliberately excludes database
// names, SQL text, query fingerprints, credentials, relations and user data.
type PostgresServerStats struct {
	Backends          int64
	TransactionsCommit int64
	TransactionsRollback int64
	BlocksRead        int64
	BlocksHit         int64
	TuplesReturned    int64
	TuplesFetched     int64
	TuplesInserted    int64
	TuplesUpdated     int64
	TuplesDeleted     int64
	Deadlocks         int64
	TempFiles         int64
	TempBytes         int64
	DatabaseSizeBytes int64
}

// PostgresPoolStats returns a point-in-time pgx pool snapshot for operational
// observability. PostgreSQL remains the durable source of truth; this method is
// read-only and does not influence admission, persistence, or gameplay state.
func (p *Postgres) PostgresPoolStats() PostgresPoolStats {
	if p == nil || p.pool == nil {
		return PostgresPoolStats{}
	}
	stats := p.pool.Stat()
	return PostgresPoolStats{
		MaxConns:             stats.MaxConns(),
		TotalConns:           stats.TotalConns(),
		IdleConns:            stats.IdleConns(),
		AcquiredConns:        stats.AcquiredConns(),
		ConstructingConns:    stats.ConstructingConns(),
		AcquireCount:         stats.AcquireCount(),
		EmptyAcquireCount:    stats.EmptyAcquireCount(),
		CanceledAcquireCount: stats.CanceledAcquireCount(),
		NewConnsCount:        stats.NewConnsCount(),
		AcquireDuration:      stats.AcquireDuration(),
	}
}

// PostgresServerStats returns current-database PostgreSQL counters using one
// fixed read-only query. No database name or SQL text is returned to callers.
func (p *Postgres) PostgresServerStats(ctx context.Context) (PostgresServerStats, error) {
	if p == nil || p.pool == nil {
		return PostgresServerStats{}, nil
	}
	const query = `SELECT
		numbackends, xact_commit, xact_rollback, blks_read, blks_hit,
		tup_returned, tup_fetched, tup_inserted, tup_updated, tup_deleted,
		deadlocks, temp_files, temp_bytes, pg_database_size(current_database())
	FROM pg_stat_database WHERE datname = current_database()`
	var stats PostgresServerStats
	err := p.pool.QueryRow(ctx, query).Scan(
		&stats.Backends,
		&stats.TransactionsCommit,
		&stats.TransactionsRollback,
		&stats.BlocksRead,
		&stats.BlocksHit,
		&stats.TuplesReturned,
		&stats.TuplesFetched,
		&stats.TuplesInserted,
		&stats.TuplesUpdated,
		&stats.TuplesDeleted,
		&stats.Deadlocks,
		&stats.TempFiles,
		&stats.TempBytes,
		&stats.DatabaseSizeBytes,
	)
	return stats, err
}
