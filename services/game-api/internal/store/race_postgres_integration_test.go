//go:build integration

package store

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/cvsz/zneondrive/services/game-api/internal/core"
)

func TestPostgresAuthoritativeRaceRuntime(t *testing.T) {
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
	if err = db.EnsureRaceSchema(ctx); err != nil {
		t.Fatal(err)
	}

	resumeKey, err := core.NewSecret()
	if err != nil {
		t.Fatal(err)
	}
	snapshot, _, err := db.Bootstrap(ctx, core.HashSecret(resumeKey))
	if err != nil {
		t.Fatal(err)
	}
	sessionToken, err := core.NewSecret()
	if err != nil {
		t.Fatal(err)
	}
	tokenHash := core.HashSecret(sessionToken)
	if err = db.CreateSession(ctx, snapshot.AccountID, tokenHash, time.Now().Add(time.Hour)); err != nil {
		t.Fatal(err)
	}

	for questNumber := 1; questNumber <= 12; questNumber++ {
		questID := fmt.Sprintf("MQ%03d", questNumber)
		snapshot, _, err = db.CompleteQuest(ctx, tokenHash, questID, "race-prereq-"+questID+"-"+snapshot.AccountID)
		if err != nil {
			t.Fatalf("complete %s: %v", questID, err)
		}
	}
	if !snapshot.Roadworthy {
		t.Fatal("MQ012 must make starter vehicle roadworthy before race start")
	}

	instance, err := db.StartRace(ctx, snapshot.AccountID, snapshot.VehicleID, "race_first_ignition", "race-start-"+snapshot.AccountID)
	if err != nil {
		t.Fatal(err)
	}
	if instance.BuildRevision != snapshot.ActiveBuildRevision || instance.BuildValidationHash == "" {
		t.Fatalf("race did not bind immutable build evidence: %+v", instance)
	}

	replay, err := db.StartRace(ctx, snapshot.AccountID, snapshot.VehicleID, "race_first_ignition", "race-start-"+snapshot.AccountID)
	if err != nil || replay.RaceInstanceID != instance.RaceInstanceID {
		t.Fatalf("race start must be idempotent: %+v %v", replay, err)
	}

	checkpoint, err := db.RecordRaceCheckpoint(ctx, instance.RaceInstanceID, 0, 1000, "race-cp-0-"+snapshot.AccountID)
	if err != nil {
		t.Fatal(err)
	}
	if checkpoint.NextCheckpoint != 1 || checkpoint.LastElapsedMS != 1000 {
		t.Fatalf("unexpected checkpoint cursor: %+v", checkpoint)
	}
	if _, err = db.RecordRaceCheckpoint(ctx, instance.RaceInstanceID, 2, 2000, "race-cp-out-of-order-"+snapshot.AccountID); !errors.Is(err, ErrRaceOrder) {
		t.Fatalf("expected out-of-order rejection, got %v", err)
	}
	if _, err = db.RecordRaceCheckpoint(ctx, instance.RaceInstanceID, 1, 900, "race-cp-regress-"+snapshot.AccountID); !errors.Is(err, ErrRaceOrder) {
		t.Fatalf("expected elapsed regression rejection, got %v", err)
	}
	if _, err = db.RecordRaceCheckpoint(ctx, instance.RaceInstanceID, 1, 2500, "race-cp-1-"+snapshot.AccountID); err != nil {
		t.Fatal(err)
	}

	result, err := db.FinishRace(ctx, instance.RaceInstanceID, 2, 3000, "race-finish-"+snapshot.AccountID)
	if err != nil {
		t.Fatal(err)
	}
	if result.ResultHash == "" || result.BuildRevision != instance.BuildRevision || result.CheckpointCount != 2 {
		t.Fatalf("unexpected race result: %+v", result)
	}
	replayResult, err := db.FinishRace(ctx, instance.RaceInstanceID, 2, 3000, "race-finish-"+snapshot.AccountID)
	if err != nil || replayResult.ResultHash != result.ResultHash {
		t.Fatalf("race finish must be idempotent: %+v %v", replayResult, err)
	}
	if _, err = db.RecordRaceCheckpoint(ctx, instance.RaceInstanceID, 2, 3500, "race-cp-after-finish-"+snapshot.AccountID); !errors.Is(err, ErrRaceOrder) {
		t.Fatalf("finished race must reject new checkpoints, got %v", err)
	}
}
