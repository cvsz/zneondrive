package core

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

type questSequenceParityVectors struct {
	SchemaVersion int `json:"schema_version"`
	QuestPrerequisiteCases []struct {
		Name      string   `json:"name"`
		QuestID   string   `json:"quest_id"`
		Completed []string `json:"completed"`
		Allowed   bool     `json:"allowed"`
	} `json:"quest_prerequisite_cases"`
}

func loadQuestSequenceParityVectors(t *testing.T) questSequenceParityVectors {
	t.Helper()
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("resolve quest sequence parity test path")
	}
	path := filepath.Clean(filepath.Join(filepath.Dir(filename), "../../../../tests/reference-parity-vectors.json"))
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read quest sequence parity vectors: %v", err)
	}
	var vectors questSequenceParityVectors
	if err := json.Unmarshal(raw, &vectors); err != nil {
		t.Fatalf("decode quest sequence parity vectors: %v", err)
	}
	if vectors.SchemaVersion != 1 {
		t.Fatalf("unsupported parity schema version: %d", vectors.SchemaVersion)
	}
	return vectors
}

func TestReferenceParityQuestPrerequisites(t *testing.T) {
	vectors := loadQuestSequenceParityVectors(t)
	if len(vectors.QuestPrerequisiteCases) == 0 {
		t.Fatal("quest_prerequisite_cases must not be empty")
	}

	seen := make(map[string]struct{}, len(vectors.QuestPrerequisiteCases))
	for _, vector := range vectors.QuestPrerequisiteCases {
		vector := vector
		t.Run(vector.Name, func(t *testing.T) {
			if vector.Name == "" {
				t.Fatal("quest prerequisite vector name must not be empty")
			}
			if _, exists := seen[vector.Name]; exists {
				t.Fatalf("duplicate quest prerequisite vector name: %s", vector.Name)
			}
			seen[vector.Name] = struct{}{}

			previous, hasPrevious, err := PreviousQuestID(vector.QuestID)
			if err != nil {
				t.Fatalf("PreviousQuestID(%q): %v", vector.QuestID, err)
			}

			completed := make(map[string]struct{}, len(vector.Completed))
			for _, questID := range vector.Completed {
				if _, err := ParseQuestID(questID); err != nil {
					t.Fatalf("invalid completed quest %q: %v", questID, err)
				}
				if _, exists := completed[questID]; exists {
					t.Fatalf("duplicate completed quest in vector: %s", questID)
				}
				completed[questID] = struct{}{}
			}

			allowed := true
			if hasPrevious {
				_, allowed = completed[previous]
			}
			if allowed != vector.Allowed {
				t.Fatalf("quest prerequisite drift for %s: got allowed=%v want %v (previous=%q completed=%v)", vector.QuestID, allowed, vector.Allowed, previous, vector.Completed)
			}
		})
	}
}
