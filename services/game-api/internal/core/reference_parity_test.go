package core

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

type parityVectors struct {
	SchemaVersion   int               `json:"schema_version"`
	BuildHashCases  []buildHashVector `json:"build_hash_cases"`
	QuestIDCases    []questIDVector   `json:"quest_id_cases"`
}

type buildHashVector struct {
	Name         string   `json:"name"`
	PartIDs      []string `json:"part_ids"`
	ExpectedHash string   `json:"expected_hash"`
}

type questIDVector struct {
	QuestID  string  `json:"quest_id"`
	Valid    bool    `json:"valid"`
	Previous *string `json:"previous"`
}

func loadReferenceParityVectors(t *testing.T) parityVectors {
	t.Helper()
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("resolve parity test path")
	}
	path := filepath.Clean(filepath.Join(filepath.Dir(filename), "../../../../tests/reference-parity-vectors.json"))
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read parity vectors: %v", err)
	}
	var vectors parityVectors
	if err := json.Unmarshal(raw, &vectors); err != nil {
		t.Fatalf("decode parity vectors: %v", err)
	}
	if vectors.SchemaVersion != 1 {
		t.Fatalf("unsupported parity schema version: %d", vectors.SchemaVersion)
	}
	return vectors
}

func TestReferenceParityBuildHashes(t *testing.T) {
	vectors := loadReferenceParityVectors(t)
	for _, vector := range vectors.BuildHashCases {
		vector := vector
		t.Run(vector.Name, func(t *testing.T) {
			got, err := BuildHash(vector.PartIDs)
			if err != nil {
				t.Fatalf("BuildHash: %v", err)
			}
			if got != vector.ExpectedHash {
				t.Fatalf("build hash drift: got %s want %s", got, vector.ExpectedHash)
			}
		})
	}
}

func TestReferenceParityQuestIDs(t *testing.T) {
	vectors := loadReferenceParityVectors(t)
	for _, vector := range vectors.QuestIDCases {
		vector := vector
		t.Run(vector.QuestID, func(t *testing.T) {
			_, err := ParseQuestID(vector.QuestID)
			if !vector.Valid {
				if err == nil {
					t.Fatalf("expected %s to be rejected", vector.QuestID)
				}
				return
			}
			if err != nil {
				t.Fatalf("expected %s to be valid: %v", vector.QuestID, err)
			}
			previous, ok, err := PreviousQuestID(vector.QuestID)
			if err != nil {
				t.Fatalf("PreviousQuestID: %v", err)
			}
			if vector.Previous == nil {
				if ok || previous != "" {
					t.Fatalf("expected no predecessor, got %q", previous)
				}
				return
			}
			if !ok || previous != *vector.Previous {
				t.Fatalf("predecessor drift: got %q want %q", previous, *vector.Previous)
			}
		})
	}
}
