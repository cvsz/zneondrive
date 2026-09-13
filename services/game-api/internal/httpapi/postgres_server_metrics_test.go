package httpapi

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/cvsz/zneondrive/services/game-api/internal/store"
)

type fakePostgresServerStatsProvider struct {
	stats       store.PostgresServerStats
	activity    store.PostgresQueryActivityStats
	err         error
	activityErr error
}

func (f fakePostgresServerStatsProvider) PostgresServerStats(context.Context) (store.PostgresServerStats, error) {
	return f.stats, f.err
}

func (f fakePostgresServerStatsProvider) PostgresQueryActivityStats(context.Context) (store.PostgresQueryActivityStats, error) {
	return f.activity, f.activityErr
}

func TestPostgresServerMetricsExposeBoundedNumericStats(t *testing.T) {
	provider := fakePostgresServerStatsProvider{
		stats: store.PostgresServerStats{
			Backends: 4, TransactionsCommit: 101, TransactionsRollback: 3,
			BlocksRead: 11, BlocksHit: 211, TuplesReturned: 501, TuplesFetched: 401,
			TuplesInserted: 31, TuplesUpdated: 21, TuplesDeleted: 7,
			Deadlocks: 2, TempFiles: 5, TempBytes: 4096, DatabaseSizeBytes: 1048576,
		},
		activity: store.PostgresQueryActivityStats{
			ActiveQueries: 3, WaitingQueries: 1, IdleInTransaction: 2, LongRunningQueries: 1,
			OldestActiveQuerySeconds: 12.5, OldestTransactionSeconds: 21.25,
		},
	}
	h := NewPostgresServerMetricsHandler(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("base_metric 1\n"))
	}), provider)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/metrics", nil))
	body := rec.Body.String()
	for _, want := range []string{
		"base_metric 1", "zneondrive_postgres_up 1", "zneondrive_postgres_backends 4",
		`zneondrive_postgres_transactions_total{outcome="commit"} 101`,
		`zneondrive_postgres_transactions_total{outcome="rollback"} 3`,
		`zneondrive_postgres_blocks_total{source="read"} 11`,
		`zneondrive_postgres_blocks_total{source="hit"} 211`,
		`zneondrive_postgres_tuples_total{operation="inserted"} 31`,
		"zneondrive_postgres_deadlocks_total 2", "zneondrive_postgres_temp_bytes_total 4096",
		"zneondrive_postgres_database_size_bytes 1048576",
		"zneondrive_postgres_query_activity_up 1",
		`zneondrive_postgres_query_activity{state="active"} 3`,
		`zneondrive_postgres_query_activity{state="waiting"} 1`,
		`zneondrive_postgres_query_activity{state="idle_in_transaction"} 2`,
		`zneondrive_postgres_query_activity{state="long_running"} 1`,
		"zneondrive_postgres_oldest_active_query_seconds 12.500000",
		"zneondrive_postgres_oldest_transaction_seconds 21.250000",
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("missing PostgreSQL server metric %q: %s", want, body)
		}
	}
	for _, forbidden := range []string{"postgres://", "DATABASE_URL", "SELECT ", "zneondrive_test", "account", "vehicle", "race", "queryid", "client_addr", "usename", "application_name"} {
		if strings.Contains(body, forbidden) {
			t.Fatalf("metrics exposed forbidden value %q: %s", forbidden, body)
		}
	}
}

func TestPostgresServerMetricsFailureDoesNotExposeError(t *testing.T) {
	provider := fakePostgresServerStatsProvider{
		err:         errors.New("dial postgres://secret@db.internal/zneondrive_private"),
		activityErr: errors.New("query leaked-player@example.invalid from 10.0.0.8"),
	}
	h := NewPostgresServerMetricsHandler(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {}), provider)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/metrics", nil))
	body := rec.Body.String()
	if !strings.Contains(body, "zneondrive_postgres_up 0") || !strings.Contains(body, "zneondrive_postgres_query_activity_up 0") {
		t.Fatalf("expected failed PostgreSQL probe metrics: %s", body)
	}
	for _, forbidden := range []string{"secret", "db.internal", "zneondrive_private", "dial postgres", "leaked-player", "10.0.0.8"} {
		if strings.Contains(body, forbidden) {
			t.Fatalf("failed PostgreSQL probe exposed error detail %q: %s", forbidden, body)
		}
	}
}
