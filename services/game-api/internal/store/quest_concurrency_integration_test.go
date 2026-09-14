//go:build integration

package store

import (
	"context"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/cvsz/zneondrive/services/game-api/internal/core"
)

func TestPostgresQuestReplayAcrossServiceInstancesIsSingleApply(t *testing.T) {
	databaseURL := os.Getenv("TEST_DATABASE_URL")
	if databaseURL == "" {
		t.Skip("TEST_DATABASE_URL is not set")
	}

	ctx := context.Background()
	primary, err := OpenPostgres(ctx, databaseURL)
	if err != nil {
		t.Fatal(err)
	}
	defer primary.Close()
	if err = primary.EnsureSchema(ctx); err != nil {
		t.Fatal(err)
	}

	secondary, err := OpenPostgres(ctx, databaseURL)
	if err != nil {
		t.Fatal(err)
	}
	defer secondary.Close()

	resumeKey, err := core.NewSecret()
	if err != nil {
		t.Fatal(err)
	}
	snapshot, _, err := primary.Bootstrap(ctx, core.HashSecret(resumeKey))
	if err != nil {
		t.Fatal(err)
	}

	sessionOne, err := core.NewSecret()
	if err != nil {
		t.Fatal(err)
	}
	sessionTwo, err := core.NewSecret()
	if err != nil {
		t.Fatal(err)
	}
	sessionOneHash := core.HashSecret(sessionOne)
	sessionTwoHash := core.HashSecret(sessionTwo)
	expiresAt := time.Now().UTC().Add(time.Hour)
	if err = primary.CreateSession(ctx, snapshot.AccountID, sessionOneHash, expiresAt); err != nil {
		t.Fatal(err)
	}
	if err = secondary.CreateSession(ctx, snapshot.AccountID, sessionTwoHash, expiresAt); err != nil {
		t.Fatal(err)
	}

	operationID, ok := core.ScopeQuestOperationID(snapshot.CharacterID, "dual-service-mq001-"+snapshot.AccountID)
	if !ok {
		t.Fatal("failed to derive durable quest operation id")
	}

	type result struct {
		snapshot core.Snapshot
		receipt  core.RewardReceipt
		err      error
	}
	results := make(chan result, 2)
	start := make(chan struct{})
	var wg sync.WaitGroup
	invoke := func(db *Postgres, tokenHash string) {
		defer wg.Done()
		<-start
		mutated, receipt, callErr := db.CompleteQuest(ctx, tokenHash, "MQ001", operationID)
		results <- result{snapshot: mutated, receipt: receipt, err: callErr}
	}

	wg.Add(2)
	go invoke(primary, sessionOneHash)
	go invoke(secondary, sessionTwoHash)
	close(start)
	wg.Wait()
	close(results)

	applied := 0
	replayed := 0
	for outcome := range results {
		if outcome.err != nil {
			t.Fatalf("concurrent quest completion failed: %v", outcome.err)
		}
		if outcome.snapshot.CharacterID != snapshot.CharacterID {
			t.Fatalf("quest result crossed authoritative character boundary: got %q want %q", outcome.snapshot.CharacterID, snapshot.CharacterID)
		}
		if outcome.receipt.Applied {
			applied++
			if outcome.receipt.Money != 110 || outcome.receipt.XP != 55 || outcome.receipt.Reputation != 1 {
				t.Fatalf("unexpected canonical MQ001 reward: %+v", outcome.receipt)
			}
		} else {
			replayed++
			if outcome.receipt.Money != 0 || outcome.receipt.XP != 0 || outcome.receipt.Reputation != 0 {
				t.Fatalf("idempotent replay must not reissue rewards: %+v", outcome.receipt)
			}
		}
	}
	if applied != 1 || replayed != 1 {
		t.Fatalf("expected exactly one applied mutation and one replay, got applied=%d replayed=%d", applied, replayed)
	}

	finalSnapshot, err := primary.SnapshotBySession(ctx, sessionOneHash)
	if err != nil {
		t.Fatal(err)
	}
	if finalSnapshot.Money != 110 || finalSnapshot.XP != 55 || finalSnapshot.Reputation != 1 {
		t.Fatalf("durable progression was not exactly-once after concurrent service calls: %+v", finalSnapshot)
	}
	if len(finalSnapshot.CompletedQuests) != 1 || finalSnapshot.CompletedQuests[0] != "MQ001" {
		t.Fatalf("expected exactly one durable MQ001 completion, got %v", finalSnapshot.CompletedQuests)
	}
}
