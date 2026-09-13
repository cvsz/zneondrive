package httpapi

import (
	"context"
	"fmt"
	"net/http"

	"github.com/cvsz/zneondrive/services/game-api/internal/store"
)

type postgresServerStatsProvider interface {
	PostgresServerStats(context.Context) (store.PostgresServerStats, error)
	PostgresQueryActivityStats(context.Context) (store.PostgresQueryActivityStats, error)
}

// NewPostgresServerMetricsHandler appends bounded current-database PostgreSQL
// server and query-activity snapshots to the existing private metrics response.
// The queries are read-only and do not affect durable or gameplay authority.
func NewPostgresServerMetricsHandler(next http.Handler, provider postgresServerStatsProvider) http.Handler {
	if provider == nil {
		return next
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
		stats, err := provider.PostgresServerStats(r.Context())
		_, _ = fmt.Fprint(w, "# HELP zneondrive_postgres_up Whether the bounded PostgreSQL server stats query succeeded.\n")
		_, _ = fmt.Fprint(w, "# TYPE zneondrive_postgres_up gauge\n")
		if err != nil {
			_, _ = fmt.Fprint(w, "zneondrive_postgres_up 0\n")
	} else {
			_, _ = fmt.Fprint(w, "zneondrive_postgres_up 1\n")
			_, _ = fmt.Fprint(w, "# HELP zneondrive_postgres_backends Current backends for the application database.\n# TYPE zneondrive_postgres_backends gauge\n")
			_, _ = fmt.Fprintf(w, "zneondrive_postgres_backends %d\n", stats.Backends)
			_, _ = fmt.Fprint(w, "# HELP zneondrive_postgres_transactions_total PostgreSQL transactions by bounded outcome.\n# TYPE zneondrive_postgres_transactions_total counter\n")
			_, _ = fmt.Fprintf(w, "zneondrive_postgres_transactions_total{outcome=\"commit\"} %d\n", stats.TransactionsCommit)
			_, _ = fmt.Fprintf(w, "zneondrive_postgres_transactions_total{outcome=\"rollback\"} %d\n", stats.TransactionsRollback)
			_, _ = fmt.Fprint(w, "# HELP zneondrive_postgres_blocks_total PostgreSQL blocks by bounded source.\n# TYPE zneondrive_postgres_blocks_total counter\n")
			_, _ = fmt.Fprintf(w, "zneondrive_postgres_blocks_total{source=\"read\"} %d\n", stats.BlocksRead)
			_, _ = fmt.Fprintf(w, "zneondrive_postgres_blocks_total{source=\"hit\"} %d\n", stats.BlocksHit)
			_, _ = fmt.Fprint(w, "# HELP zneondrive_postgres_tuples_total PostgreSQL tuples by bounded operation.\n# TYPE zneondrive_postgres_tuples_total counter\n")
			_, _ = fmt.Fprintf(w, "zneondrive_postgres_tuples_total{operation=\"returned\"} %d\n", stats.TuplesReturned)
			_, _ = fmt.Fprintf(w, "zneondrive_postgres_tuples_total{operation=\"fetched\"} %d\n", stats.TuplesFetched)
			_, _ = fmt.Fprintf(w, "zneondrive_postgres_tuples_total{operation=\"inserted\"} %d\n", stats.TuplesInserted)
			_, _ = fmt.Fprintf(w, "zneondrive_postgres_tuples_total{operation=\"updated\"} %d\n", stats.TuplesUpdated)
			_, _ = fmt.Fprintf(w, "zneondrive_postgres_tuples_total{operation=\"deleted\"} %d\n", stats.TuplesDeleted)
			_, _ = fmt.Fprint(w, "# HELP zneondrive_postgres_deadlocks_total PostgreSQL deadlocks.\n# TYPE zneondrive_postgres_deadlocks_total counter\n")
			_, _ = fmt.Fprintf(w, "zneondrive_postgres_deadlocks_total %d\n", stats.Deadlocks)
			_, _ = fmt.Fprint(w, "# HELP zneondrive_postgres_temp_files_total PostgreSQL temporary files.\n# TYPE zneondrive_postgres_temp_files_total counter\n")
			_, _ = fmt.Fprintf(w, "zneondrive_postgres_temp_files_total %d\n", stats.TempFiles)
			_, _ = fmt.Fprint(w, "# HELP zneondrive_postgres_temp_bytes_total PostgreSQL temporary bytes.\n# TYPE zneondrive_postgres_temp_bytes_total counter\n")
			_, _ = fmt.Fprintf(w, "zneondrive_postgres_temp_bytes_total %d\n", stats.TempBytes)
			_, _ = fmt.Fprint(w, "# HELP zneondrive_postgres_database_size_bytes Current application database size in bytes.\n# TYPE zneondrive_postgres_database_size_bytes gauge\n")
			_, _ = fmt.Fprintf(w, "zneondrive_postgres_database_size_bytes %d\n", stats.DatabaseSizeBytes)
		}

		activity, activityErr := provider.PostgresQueryActivityStats(r.Context())
		_, _ = fmt.Fprint(w, "# HELP zneondrive_postgres_query_activity_up Whether the bounded pg_stat_activity query succeeded.\n")
		_, _ = fmt.Fprint(w, "# TYPE zneondrive_postgres_query_activity_up gauge\n")
		if activityErr != nil {
			_, _ = fmt.Fprint(w, "zneondrive_postgres_query_activity_up 0\n")
			return
		}
		_, _ = fmt.Fprint(w, "zneondrive_postgres_query_activity_up 1\n")
		_, _ = fmt.Fprint(w, "# HELP zneondrive_postgres_query_activity Current queries by bounded state.\n# TYPE zneondrive_postgres_query_activity gauge\n")
		_, _ = fmt.Fprintf(w, "zneondrive_postgres_query_activity{state=\"active\"} %d\n", activity.ActiveQueries)
		_, _ = fmt.Fprintf(w, "zneondrive_postgres_query_activity{state=\"waiting\"} %d\n", activity.WaitingQueries)
		_, _ = fmt.Fprintf(w, "zneondrive_postgres_query_activity{state=\"idle_in_transaction\"} %d\n", activity.IdleInTransaction)
		_, _ = fmt.Fprintf(w, "zneondrive_postgres_query_activity{state=\"long_running\"} %d\n", activity.LongRunningQueries)
		_, _ = fmt.Fprint(w, "# HELP zneondrive_postgres_oldest_active_query_seconds Age of the oldest active query, excluding the metrics query itself.\n# TYPE zneondrive_postgres_oldest_active_query_seconds gauge\n")
		_, _ = fmt.Fprintf(w, "zneondrive_postgres_oldest_active_query_seconds %.6f\n", activity.OldestActiveQuerySeconds)
		_, _ = fmt.Fprint(w, "# HELP zneondrive_postgres_oldest_transaction_seconds Age of the oldest open transaction, excluding the metrics query itself.\n# TYPE zneondrive_postgres_oldest_transaction_seconds gauge\n")
		_, _ = fmt.Fprintf(w, "zneondrive_postgres_oldest_transaction_seconds %.6f\n", activity.OldestTransactionSeconds)
	})
}
