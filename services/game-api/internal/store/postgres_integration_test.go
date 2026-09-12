//go:build integration

package store

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/cvsz/zneondrive/services/game-api/internal/core"
)

func TestPostgresDurableVerticalSliceState(t *testing.T) {
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
	if err = db.EnsureSchema(ctx); err != nil {
		t.Fatal(err)
	}

	resumeKey, err := core.NewSecret()
	if err != nil {
		t.Fatal(err)
	}
	snapshot, created, err := db.Bootstrap(ctx, core.HashSecret(resumeKey))
	if err != nil {
		t.Fatal(err)
	}
	if !created || !snapshot.StarterLineage || snapshot.ActiveBuildRevision != 1 {
		t.Fatalf("unexpected bootstrap snapshot: %+v", snapshot)
	}

	sessionToken, err := core.NewSecret()
	if err != nil {
		t.Fatal(err)
	}
	tokenHash := core.HashSecret(sessionToken)
	if err = db.CreateSession(ctx, snapshot.AccountID, tokenHash, time.Now().Add(time.Hour)); err != nil {
		t.Fatal(err)
	}

	firstState, receipt, err := db.CompleteQuest(ctx, tokenHash, "MQ001", "quest-op-1-"+snapshot.AccountID)
	if err != nil {
		t.Fatal(err)
	}
	if !receipt.Applied || firstState.Money <= 0 {
		t.Fatalf("quest reward not applied: %+v %+v", receipt, firstState)
	}

	retryState, retryReceipt, err := db.CompleteQuest(ctx, tokenHash, "MQ001", "quest-op-1-"+snapshot.AccountID)
	if err != nil {
		t.Fatal(err)
	}
	if retryReceipt.Applied || retryState.Money != firstState.Money {
		t.Fatal("idempotent retry duplicated quest reward")
	}

	parts := append([]string{}, snapshot.ActivePartIDs...)
	parts = append(parts, "part_brakes_track_i")
	revised, err := db.ReviseBuild(
		ctx,
		tokenHash,
		snapshot.VehicleID,
		1,
		parts,
		"build-op-1-"+snapshot.AccountID,
	)
	if err != nil {
		t.Fatal(err)
	}
	if revised.ActiveBuildRevision != 2 {
		t.Fatalf("expected build revision 2, got %d", revised.ActiveBuildRevision)
	}

	retryBuild, err := db.ReviseBuild(
		ctx,
		tokenHash,
		snapshot.VehicleID,
		1,
		parts,
		"build-op-1-"+snapshot.AccountID,
	)
	if err != nil {
		t.Fatal(err)
	}
	if retryBuild.ActiveBuildRevision != 2 {
		t.Fatal("idempotent build retry changed revision")
	}

	if _, err = db.ReviseBuild(
		ctx,
		tokenHash,
		snapshot.VehicleID,
		1,
		parts,
		"build-op-stale-"+snapshot.AccountID,
	); err != ErrConflict {
		t.Fatalf("expected ErrConflict for stale build revision, got %v", err)
	}
}
