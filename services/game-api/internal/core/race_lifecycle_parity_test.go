package core

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

type raceLifecycleVectors struct {
	CheckpointCases []struct {
		Name           string `json:"name"`
		State          string `json:"state"`
		NextCheckpoint int    `json:"next_checkpoint"`
		LastElapsedMS  int64  `json:"last_elapsed_ms"`
		CheckpointIndex int   `json:"checkpoint_index"`
		ElapsedMS      int64  `json:"elapsed_ms"`
		Valid          bool   `json:"valid"`
	} `json:"checkpoint_cases"`
	FinishCases []struct {
		Name            string `json:"name"`
		State           string `json:"state"`
		NextCheckpoint  int    `json:"next_checkpoint"`
		LastElapsedMS   int64  `json:"last_elapsed_ms"`
		CheckpointCount int    `json:"checkpoint_count"`
		FinishElapsedMS int64  `json:"finish_elapsed_ms"`
		Valid           bool   `json:"valid"`
	} `json:"finish_cases"`
}

func loadRaceLifecycleVectors(t *testing.T) raceLifecycleVectors {
	t.Helper()
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("resolve lifecycle parity test path")
	}
	path := filepath.Clean(filepath.Join(filepath.Dir(filename), "../../../../tests/race-lifecycle-parity-v4.0.json"))
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read lifecycle parity vectors: %v", err)
	}
	var vectors raceLifecycleVectors
	if err := json.Unmarshal(raw, &vectors); err != nil {
		t.Fatalf("decode lifecycle parity vectors: %v", err)
	}
	if len(vectors.CheckpointCases) == 0 || len(vectors.FinishCases) == 0 {
		t.Fatal("race lifecycle parity vectors must cover checkpoint and finish cases")
	}
	return vectors
}

func TestRaceLifecycleCheckpointParity(t *testing.T) {
	for _, vector := range loadRaceLifecycleVectors(t).CheckpointCases {
		vector := vector
		t.Run(vector.Name, func(t *testing.T) {
			instance := RaceInstance{State: vector.State, NextCheckpoint: vector.NextCheckpoint, LastElapsedMS: vector.LastElapsedMS}
			err := ValidateRaceCheckpointAdvance(instance, vector.CheckpointIndex, vector.ElapsedMS)
			if vector.Valid && err != nil {
				t.Fatalf("expected checkpoint acceptance: %v", err)
			}
			if !vector.Valid && err == nil {
				t.Fatal("expected checkpoint rejection")
			}
		})
	}
}

func TestRaceLifecycleFinishParity(t *testing.T) {
	for _, vector := range loadRaceLifecycleVectors(t).FinishCases {
		vector := vector
		t.Run(vector.Name, func(t *testing.T) {
			instance := RaceInstance{State: vector.State, NextCheckpoint: vector.NextCheckpoint, LastElapsedMS: vector.LastElapsedMS}
			err := ValidateRaceFinish(instance, vector.CheckpointCount, vector.FinishElapsedMS)
			if vector.Valid && err != nil {
				t.Fatalf("expected finish acceptance: %v", err)
			}
			if !vector.Valid && err == nil {
				t.Fatal("expected finish rejection")
			}
		})
	}
}
