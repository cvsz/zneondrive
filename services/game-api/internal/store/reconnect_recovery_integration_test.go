//go:build integration

package store

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/cvsz/zneondrive/services/game-api/internal/core"
)

func TestPostgresServiceRestartReconnectRecovery(t *testing.T) {
	databaseURL := os.Getenv("TEST_DATABASE_URL")
	if databaseURL == "" {
		t.Skip("TEST_DATABASE_URL is not set")
	}

	ctx := context.Background()
	db, err := OpenPostgres(ctx, databaseURL)
	if err != nil {
		t.Fatal(err)
	}
	if err = db.EnsureSchema(ctx); err != nil {
		db.Close()
		t.Fatal(err)
	}

	resumeKey, err := core.NewSecret()
	if err != nil {
		db.Close()
		t.Fatal(err)
	}
	resumeHash := core.HashSecret(resumeKey)
	before, created, err := db.Bootstrap(ctx, resumeHash)
	if err != nil {
		db.Close()
		t.Fatal(err)
	}
	if !created {
		db.Close()
		t.Fatal("expected fresh bootstrap")
	}

	sessionOne, err := core.NewSecret()
	if err != nil {
		db.Close()
		t.Fatal(err)
	}
	sessionOneHash := core.HashSecret(sessionOne)
	if err = db.CreateSession(ctx, before.AccountID, sessionOneHash, time.Now().Add(time.Hour)); err != nil {
		db.Close()
		t.Fatal(err)
	}

	clientOperationID := "restart-recovery-mq001-" + before.AccountID
	durableOperationID, ok := core.ScopeQuestOperationID(before.CharacterID, clientOperationID)
	if !ok {
		db.Close()
		t.Fatal("expected scoped quest operation id")
	}
	mutated, receipt, err := db.CompleteQuest(ctx, sessionOneHash, "MQ001", durableOperationID)
	if err != nil {
		db.Close()
		t.Fatal(err)
	}
	if !receipt.Applied || mutated.Money <= before.Money || mutated.XP <= before.XP {
		db.Close()
		t.Fatalf("MQ001 did not persist before restart: before=%+v after=%+v receipt=%+v", before, mutated, receipt)
	}

	// Closing and reopening the PostgreSQL-backed store simulates service-process
	// loss/restart without relying on in-memory state from the first process.
	db.Close()

	db, err = OpenPostgres(ctx, databaseURL)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if err = db.EnsureSchema(ctx); err != nil {
		t.Fatal(err)
	}

	resumed, created, err := db.Bootstrap(ctx, resumeHash)
	if err != nil {
		t.Fatal(err)
	}
	if created {
		t.Fatal("restart recovery must resume the existing durable account")
	}
	if resumed.AccountID != before.AccountID || resumed.CharacterID != before.CharacterID || resumed.VehicleID != before.VehicleID {
		t.Fatalf("durable identity changed across service restart: before=%+v resumed=%+v", before, resumed)
	}
	if resumed.Money != mutated.Money || resumed.XP != mutated.XP || resumed.Reputation != mutated.Reputation {
		t.Fatalf("durable progression changed across service restart: mutated=%+v resumed=%+v", mutated, resumed)
	}

	// A renewed session after process restart must recover the same authoritative
	// snapshot and preserve operation-id replay semantics for the same character.
	sessionTwo, err := core.NewSecret()
	if err != nil {
		t.Fatal(err)
	}
	sessionTwoHash := core.HashSecret(sessionTwo)
	if err = db.CreateSession(ctx, resumed.AccountID, sessionTwoHash, time.Now().Add(time.Hour)); err != nil {
		t.Fatal(err)
	}

	bySession, err := db.SnapshotBySession(ctx, sessionTwoHash)
	if err != nil {
		t.Fatal(err)
	}
	if bySession.AccountID != before.AccountID || bySession.CharacterID != before.CharacterID {
		t.Fatalf("renewed session resolved wrong authority after restart: %+v", bySession)
	}

	replayed, replayReceipt, err := db.CompleteQuest(ctx, sessionTwoHash, "MQ001", durableOperationID)
	if err != nil {
		t.Fatal(err)
	}
	if replayReceipt.Applied {
		t.Fatal("quest operation replay after restart must not re-apply rewards")
	}
	if replayed.Money != mutated.Money || replayed.XP != mutated.XP || replayed.Reputation != mutated.Reputation {
		t.Fatalf("quest replay after restart changed progression: mutated=%+v replayed=%+v", mutated, replayed)
	}

	mq002Operation, ok := core.ScopeQuestOperationID(before.CharacterID, "restart-recovery-mq002-"+before.AccountID)
	if !ok {
		t.Fatal("expected scoped MQ002 operation id")
	}
	afterMQ002, mq002Receipt, err := db.CompleteQuest(ctx, sessionTwoHash, "MQ002", mq002Operation)
	if err != nil {
		t.Fatalf("complete MQ002 after restart: %v", err)
	}
	if !mq002Receipt.Applied || afterMQ002.Money <= replayed.Money {
		t.Fatalf("post-restart mutation did not advance authoritative state: replayed=%+v after=%+v receipt=%+v", replayed, afterMQ002, mq002Receipt)
	}

	_ = fmt.Sprintf("%s/%s", resumed.AccountID, resumed.CharacterID) // keep the evidence deliberately identity-local without logging secrets
}
