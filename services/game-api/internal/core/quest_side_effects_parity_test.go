package core

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

type questSideEffectVectors struct {
	SchemaVersion int `json:"schema_version"`
	Cases []struct {
		QuestID string `json:"quest_id"`
		InventoryItemID string `json:"inventory_item_id"`
		InventoryQuantity int64 `json:"inventory_quantity"`
		BlueprintID string `json:"blueprint_id"`
		MarkStarterRoadworthy bool `json:"mark_starter_roadworthy"`
	} `json:"cases"`
	InvalidQuestIDs []string `json:"invalid_quest_ids"`
}

func loadQuestSideEffectVectors(t *testing.T) questSideEffectVectors {
	t.Helper()
	_, filename, _, ok := runtime.Caller(0)
	if !ok { t.Fatal("resolve quest side-effect parity path") }
	path := filepath.Clean(filepath.Join(filepath.Dir(filename), "../../../../tests/quest-side-effects-parity-v4.3.json"))
	raw, err := os.ReadFile(path)
	if err != nil { t.Fatalf("read quest side-effect vectors: %v", err) }
	var vectors questSideEffectVectors
	if err := json.Unmarshal(raw, &vectors); err != nil { t.Fatalf("decode quest side-effect vectors: %v", err) }
	if vectors.SchemaVersion != 1 { t.Fatalf("unsupported schema version: %d", vectors.SchemaVersion) }
	return vectors
}

func TestQuestSideEffectsParity(t *testing.T) {
	for _, vector := range loadQuestSideEffectVectors(t).Cases {
		effect, err := QuestSideEffectsFor(vector.QuestID)
		if err != nil { t.Fatalf("QuestSideEffectsFor(%q): %v", vector.QuestID, err) }
		if effect.InventoryItemID != vector.InventoryItemID || effect.InventoryQuantity != vector.InventoryQuantity || effect.BlueprintID != vector.BlueprintID || effect.MarkStarterRoadworthy != vector.MarkStarterRoadworthy {
			t.Fatalf("quest side-effect drift for %s: got %+v", vector.QuestID, effect)
		}
	}
}

func TestQuestSideEffectsRejectInvalidQuestIDs(t *testing.T) {
	for _, questID := range loadQuestSideEffectVectors(t).InvalidQuestIDs {
		if _, err := QuestSideEffectsFor(questID); err == nil { t.Fatalf("QuestSideEffectsFor(%q) unexpectedly succeeded", questID) }
	}
}
