package store

import "time"

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
