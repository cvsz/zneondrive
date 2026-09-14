package core

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

type questRewardParityVectors struct {
	SchemaVersion int `json:"schema_version"`
	Cases []struct {
		QuestID    string `json:"quest_id"`
		Money      int64  `json:"money"`
		XP         int64  `json:"xp"`
		Reputation int64  `json:"reputation"`
	} `json:"cases"`
	InvalidQuestIDs []string `json:"invalid_quest_ids"`
}

func loadQuestRewardParityVectors(t *testing.T) questRewardParityVectors {
	t.Helper()
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("resolve quest reward parity test path")
	}
	path := filepath.Clean(filepath.Join(filepath.Dir(filename), "../../../../tests/quest-reward-parity-v3.8.json"))
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read quest reward parity vectors: %v", err)
	}
	var vectors questRewardParityVectors
	if err := json.Unmarshal(raw, &vectors); err != nil {
		t.Fatalf("decode quest reward parity vectors: %v", err)
	}
	if vectors.SchemaVersion != 1 {
		t.Fatalf("unsupported quest reward parity schema version: %d", vectors.SchemaVersion)
	}
	return vectors
}

func TestQuestRewardParity(t *testing.T) {
	vectors := loadQuestRewardParityVectors(t)
	if len(vectors.Cases) == 0 {
		t.Fatal("quest reward parity cases must not be empty")
	}
	for _, vector := range vectors.Cases {
		vector := vector
		t.Run(vector.QuestID, func(t *testing.T) {
			reward, err := QuestReward(vector.QuestID)
			if err != nil {
				t.Fatalf("QuestReward(%q): %v", vector.QuestID, err)
			}
			if reward.QuestID != vector.QuestID || !reward.Applied || reward.Money != vector.Money || reward.XP != vector.XP || reward.Reputation != vector.Reputation {
				t.Fatalf("quest reward drift for %s: got %+v want money=%d xp=%d reputation=%d", vector.QuestID, reward, vector.Money, vector.XP, vector.Reputation)
			}
		})
	}
}

func TestQuestRewardRejectsInvalidQuestIDs(t *testing.T) {
	vectors := loadQuestRewardParityVectors(t)
	for _, questID := range vectors.InvalidQuestIDs {
		if _, err := QuestReward(questID); err == nil {
			t.Fatalf("QuestReward(%q) unexpectedly succeeded", questID)
		}
	}
}
