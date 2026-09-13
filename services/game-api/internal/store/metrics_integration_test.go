//go:build integration

package store

import (
	"context"
	"os"
	"testing"
)

func TestPostgresServerStatsIntegration(t *testing.T) {
	databaseURL := os.Getenv("TEST_DATABASE_URL")
	if databaseURL == "" {
		t.Skip("TEST_DATABASE_URL is not set")
	}
	ctx := context.Background()
	db, err := OpenPostgres(ctx, databaseURL)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	stats, err := db.PostgresServerStats(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if stats.Backends < 1 {
		t.Fatalf("expected at least one PostgreSQL backend, got %d", stats.Backends)
	}
	if stats.DatabaseSizeBytes <= 0 {
		t.Fatalf("expected positive PostgreSQL database size, got %d", stats.DatabaseSizeBytes)
	}
	if stats.TransactionsCommit < 0 || stats.TransactionsRollback < 0 || stats.BlocksRead < 0 || stats.BlocksHit < 0 || stats.Deadlocks < 0 {
		t.Fatalf("unexpected negative PostgreSQL counters: %+v", stats)
	}

	activity, err := db.PostgresQueryActivityStats(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if activity.ActiveQueries < 0 || activity.WaitingQueries < 0 || activity.IdleInTransaction < 0 || activity.LongRunningQueries < 0 {
		t.Fatalf("unexpected negative PostgreSQL query activity counters: %+v", activity)
	}
	if activity.OldestActiveQuerySeconds < 0 || activity.OldestTransactionSeconds < 0 {
		t.Fatalf("unexpected negative PostgreSQL query activity ages: %+v", activity)
	}
}
