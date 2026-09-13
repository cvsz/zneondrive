package core

import (
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

type parityVectors struct {
	SchemaVersion       int                    `json:"schema_version"`
	BuildHashCases      []buildHashVector      `json:"build_hash_cases"`
	QuestIDCases        []questIDVector        `json:"quest_id_cases"`
	RaceIDCases         []raceIDVector         `json:"race_id_cases"`
	RaceCheckpointCases []raceCheckpointVector `json:"race_checkpoint_cases"`
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

type raceIDVector struct {
	Name       string  `json:"name"`
	RaceID     string  `json:"race_id"`
	Valid      bool    `json:"valid"`
	Normalized *string `json:"normalized"`
}

type raceCheckpointVector struct {
	Name      string `json:"name"`
	Index     int    `json:"index"`
	ElapsedMS int64  `json:"elapsed_ms"`
	Valid     bool   `json:"valid"`
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

func TestReferenceParityVectorIntegrity(t *testing.T) {
	vectors := loadReferenceParityVectors(t)
	if len(vectors.BuildHashCases) == 0 {
		t.Fatal("build_hash_cases must not be empty")
	}
	if len(vectors.QuestIDCases) == 0 {
		t.Fatal("quest_id_cases must not be empty")
	}
	if len(vectors.RaceIDCases) == 0 {
		t.Fatal("race_id_cases must not be empty")
	}
	if len(vectors.RaceCheckpointCases) == 0 {
		t.Fatal("race_checkpoint_cases must not be empty")
	}

	buildNames := make(map[string]struct{}, len(vectors.BuildHashCases))
	for _, vector := range vectors.BuildHashCases {
		if vector.Name == "" {
			t.Fatal("build hash vector name must not be empty")
		}
		if _, exists := buildNames[vector.Name]; exists {
			t.Fatalf("duplicate build hash vector name: %s", vector.Name)
		}
		buildNames[vector.Name] = struct{}{}
		if len(vector.PartIDs) == 0 {
			t.Fatalf("build hash vector %s has no parts", vector.Name)
		}
		hashBytes, err := hex.DecodeString(vector.ExpectedHash)
		if err != nil || len(hashBytes) != 32 {
			t.Fatalf("build hash vector %s has invalid SHA-256: %q", vector.Name, vector.ExpectedHash)
		}
	}

	questIDs := make(map[string]struct{}, len(vectors.QuestIDCases))
	for _, vector := range vectors.QuestIDCases {
		if _, exists := questIDs[vector.QuestID]; exists {
			t.Fatalf("duplicate quest vector: %s", vector.QuestID)
		}
		questIDs[vector.QuestID] = struct{}{}
		if !vector.Valid && vector.Previous != nil {
			t.Fatalf("invalid quest vector %s must not declare a predecessor", vector.QuestID)
		}
	}

	raceNames := make(map[string]struct{}, len(vectors.RaceIDCases))
	for _, vector := range vectors.RaceIDCases {
		if vector.Name == "" {
			t.Fatal("race id vector name must not be empty")
		}
		if _, exists := raceNames[vector.Name]; exists {
			t.Fatalf("duplicate race id vector name: %s", vector.Name)
		}
		raceNames[vector.Name] = struct{}{}
		if vector.Valid && (vector.Normalized == nil || *vector.Normalized == "") {
			t.Fatalf("valid race id vector %s must declare normalized value", vector.Name)
		}
		if !vector.Valid && vector.Normalized != nil {
			t.Fatalf("invalid race id vector %s must not declare normalized value", vector.Name)
		}
	}

	checkpointNames := make(map[string]struct{}, len(vectors.RaceCheckpointCases))
	for _, vector := range vectors.RaceCheckpointCases {
		if vector.Name == "" {
			t.Fatal("race checkpoint vector name must not be empty")
		}
		if _, exists := checkpointNames[vector.Name]; exists {
			t.Fatalf("duplicate race checkpoint vector name: %s", vector.Name)
		}
		checkpointNames[vector.Name] = struct{}{}
	}
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

func TestReferenceParityRaceIDs(t *testing.T) {
	vectors := loadReferenceParityVectors(t)
	for _, vector := range vectors.RaceIDCases {
		vector := vector
		t.Run(vector.Name, func(t *testing.T) {
			got, err := NormalizeRaceID(vector.RaceID)
			if !vector.Valid {
				if err == nil {
					t.Fatalf("expected race id %q to be rejected", vector.RaceID)
				}
				return
			}
			if err != nil {
				t.Fatalf("expected race id %q to be valid: %v", vector.RaceID, err)
			}
			if vector.Normalized == nil || got != *vector.Normalized {
				t.Fatalf("race id normalization drift: got %q want %v", got, vector.Normalized)
			}
		})
	}
}

func TestReferenceParityRaceCheckpoints(t *testing.T) {
	vectors := loadReferenceParityVectors(t)
	for _, vector := range vectors.RaceCheckpointCases {
		vector := vector
		t.Run(vector.Name, func(t *testing.T) {
			err := ValidateRaceCheckpoint(vector.Index, vector.ElapsedMS)
			if vector.Valid && err != nil {
				t.Fatalf("expected checkpoint to be valid: %v", err)
			}
			if !vector.Valid && err == nil {
				t.Fatalf("expected checkpoint index=%d elapsed_ms=%d to be rejected", vector.Index, vector.ElapsedMS)
			}
		})
	}
}
