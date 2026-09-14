//go:build integration

package store

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/cvsz/zneondrive/services/game-api/internal/core"
)

func TestPostgresQuestReplayRejectsCrossCharacterOperationReuse(t *testing.T) {
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

	bootstrap := func() (core.Snapshot, string) {
		resumeKey, err := core.NewSecret()
		if err != nil {
			t.Fatal(err)
		}
		snapshot, created, err := db.Bootstrap(ctx, core.HashSecret(resumeKey))
		if err != nil {
			t.Fatal(err)
		}
		if !created {
			t.Fatal("expected isolated bootstrap to create a new account")
		}
		session, err := core.NewSecret()
		if err != nil {
			t.Fatal(err)
		}
		sessionHash := core.HashSecret(session)
		if err = db.CreateSession(ctx, snapshot.AccountID, sessionHash, time.Now().UTC().Add(time.Hour)); err != nil {
			t.Fatal(err)
		}
		return snapshot, sessionHash
	}

	first, firstSession := bootstrap()
	second, secondSession := bootstrap()
	if first.CharacterID == second.CharacterID {
		t.Fatal("expected distinct authoritative characters")
	}

	operationID := "store-boundary-shared-operation"
	firstAfter, firstReceipt, err := db.CompleteQuest(ctx, firstSession, "MQ001", operationID)
	if err != nil {
		t.Fatal(err)
	}
	if !firstReceipt.Applied || firstAfter.Money != 110 || firstAfter.XP != 55 || firstAfter.Reputation != 1 {
		t.Fatalf("unexpected first authoritative mutation: snapshot=%+v receipt=%+v", firstAfter, firstReceipt)
	}

	_, _, err = db.CompleteQuest(ctx, secondSession, "MQ001", operationID)
	if !errors.Is(err, ErrOperationKey) {
		t.Fatalf("cross-character operation reuse must be rejected with ErrOperationKey, got %v", err)
	}

	secondAfter, err := db.SnapshotBySession(ctx, secondSession)
	if err != nil {
		t.Fatal(err)
	}
	if secondAfter.Money != 0 || secondAfter.XP != 0 || secondAfter.Reputation != 0 || len(secondAfter.CompletedQuests) != 0 {
		t.Fatalf("rejected cross-character replay mutated the second character: %+v", secondAfter)
	}

	firstReplay, replayReceipt, err := db.CompleteQuest(ctx, firstSession, "MQ001", operationID)
	if err != nil {
		t.Fatal(err)
	}
	if replayReceipt.Applied || replayReceipt.Money != 0 || replayReceipt.XP != 0 || replayReceipt.Reputation != 0 {
		t.Fatalf("same-character replay must remain idempotent: %+v", replayReceipt)
	}
	if firstReplay.Money != 110 || firstReplay.XP != 55 || firstReplay.Reputation != 1 || len(firstReplay.CompletedQuests) != 1 {
		t.Fatalf("same-character replay changed authoritative state: %+v", firstReplay)
	}
}
