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

func TestPostgresInventoryBlueprintAndRebuildState(t *testing.T) {
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
	resumeHash := core.HashSecret(resumeKey)
	snapshot, created, err := db.Bootstrap(ctx, resumeHash)
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

	preBlueprintParts := append([]string{}, snapshot.ActivePartIDs...)
	preBlueprintParts = append(preBlueprintParts, "part_brakes_track_i")
	if _, err = db.ReviseBuild(
		ctx, tokenHash, snapshot.VehicleID, 1, preBlueprintParts,
		"build-before-blueprint-"+snapshot.AccountID,
	); !errors.Is(err, ErrBlueprintRequired) {
		t.Fatalf("expected ErrBlueprintRequired before MQ005, got %v", err)
	}

	var state core.Snapshot
	for questNumber := 1; questNumber <= 4; questNumber++ {
		questID := fmt.Sprintf("MQ%03d", questNumber)
		state, _, err = db.CompleteQuest(
			ctx, tokenHash, questID, "quest-"+questID+"-"+snapshot.AccountID,
		)
		if err != nil {
			t.Fatalf("complete %s: %v", questID, err)
		}
	}
	if qty := inventoryQuantity(state, "part_brakes_track_i"); qty != 1 {
		t.Fatalf("MQ004 should grant one salvage part, got %d", qty)
	}

	retryMQ004, receipt, err := db.CompleteQuest(
		ctx, tokenHash, "MQ004", "quest-MQ004-"+snapshot.AccountID,
	)
	if err != nil {
		t.Fatal(err)
	}
	if receipt.Applied || inventoryQuantity(retryMQ004, "part_brakes_track_i") != 1 {
		t.Fatal("MQ004 retry duplicated salvage inventory")
	}

	state, _, err = db.CompleteQuest(
		ctx, tokenHash, "MQ005", "quest-MQ005-"+snapshot.AccountID,
	)
	if err != nil {
		t.Fatal(err)
	}
	if !containsString(state.Blueprints, core.StarterRebuildBlueprint) {
		t.Fatalf("MQ005 did not unlock %s: %+v", core.StarterRebuildBlueprint, state.Blueprints)
	}

	rebuildParts := append([]string{}, state.ActivePartIDs...)
	rebuildParts = append(rebuildParts, "part_brakes_track_i")
	revised, err := db.ReviseBuild(
		ctx,
		tokenHash,
		snapshot.VehicleID,
		1,
		rebuildParts,
		"build-op-1-"+snapshot.AccountID,
	)
	if err != nil {
		t.Fatal(err)
	}
	if revised.ActiveBuildRevision != 2 {
		t.Fatalf("expected build revision 2, got %d", revised.ActiveBuildRevision)
	}
	if qty := inventoryQuantity(revised, "part_brakes_track_i"); qty != 0 {
		t.Fatalf("equipped part should be consumed from inventory, got %d", qty)
	}

	retryBuild, err := db.ReviseBuild(
		ctx,
		tokenHash,
		snapshot.VehicleID,
		1,
		rebuildParts,
		"build-op-1-"+snapshot.AccountID,
	)
	if err != nil {
		t.Fatal(err)
	}
	if retryBuild.ActiveBuildRevision != 2 || inventoryQuantity(retryBuild, "part_brakes_track_i") != 0 {
		t.Fatal("idempotent build retry changed revision or inventory")
	}

	changedPayload := append([]string{}, rebuildParts...)
	changedPayload = append(changedPayload, "part_tires_street_i")
	if _, err = db.ReviseBuild(
		ctx,
		tokenHash,
		snapshot.VehicleID,
		1,
		changedPayload,
		"build-op-1-"+snapshot.AccountID,
	); !errors.Is(err, ErrOperationKey) {
		t.Fatalf("expected operation key conflict for changed payload, got %v", err)
	}

	unowned := append([]string{}, rebuildParts...)
	unowned = append(unowned, "part_aero_downforce_i")
	if _, err = db.ReviseBuild(
		ctx,
		tokenHash,
		snapshot.VehicleID,
		2,
		unowned,
		"build-unowned-"+snapshot.AccountID,
	); !errors.Is(err, ErrInsufficientInventory) {
		t.Fatalf("expected ErrInsufficientInventory, got %v", err)
	}

	for questNumber := 6; questNumber <= 9; questNumber++ {
		questID := fmt.Sprintf("MQ%03d", questNumber)
		state, _, err = db.CompleteQuest(
			ctx, tokenHash, questID, "quest-"+questID+"-"+snapshot.AccountID,
		)
		if err != nil {
			t.Fatalf("complete %s: %v", questID, err)
		}
	}
	if qty := inventoryQuantity(state, "part_tires_street_i"); qty != 1 {
		t.Fatalf("MQ009 should grant one recovered item, got %d", qty)
	}

	retryMQ009, retryReceipt, err := db.CompleteQuest(
		ctx, tokenHash, "MQ009", "quest-MQ009-"+snapshot.AccountID,
	)
	if err != nil {
		t.Fatal(err)
	}
	if retryReceipt.Applied || inventoryQuantity(retryMQ009, "part_tires_street_i") != 1 {
		t.Fatal("MQ009 retry duplicated recovered item")
	}

	swapParts := make([]string, 0, len(rebuildParts))
	for _, partID := range rebuildParts {
		if partID != "part_brakes_track_i" {
			swapParts = append(swapParts, partID)
		}
	}
	swapParts = append(swapParts, "part_tires_street_i")
	swapped, err := db.ReviseBuild(
		ctx,
		tokenHash,
		snapshot.VehicleID,
		2,
		swapParts,
		"build-op-2-"+snapshot.AccountID,
	)
	if err != nil {
		t.Fatal(err)
	}
	if swapped.ActiveBuildRevision != 3 {
		t.Fatalf("expected build revision 3, got %d", swapped.ActiveBuildRevision)
	}
	if inventoryQuantity(swapped, "part_tires_street_i") != 0 ||
		inventoryQuantity(swapped, "part_brakes_track_i") != 1 {
		t.Fatalf("inventory swap accounting incorrect: %+v", swapped.Inventory)
	}

	resumed, created, err := db.Bootstrap(ctx, resumeHash)
	if err != nil {
		t.Fatal(err)
	}
	if created || resumed.AccountID != snapshot.AccountID || resumed.ActiveBuildRevision != 3 {
		t.Fatalf("resume did not preserve durable rebuild state: %+v", resumed)
	}
	if inventoryQuantity(resumed, "part_brakes_track_i") != 1 ||
		!containsString(resumed.Blueprints, core.StarterRebuildBlueprint) {
		t.Fatalf("resume lost inventory/blueprint state: %+v", resumed)
	}
}

func inventoryQuantity(snapshot core.Snapshot, itemID string) int64 {
	for _, item := range snapshot.Inventory {
		if item.ItemID == itemID {
			return item.Quantity
		}
	}
	return 0
}

func containsString(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}
