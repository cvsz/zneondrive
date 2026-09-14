package core

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

type raceOperationVectors struct {
	StartCases []struct {
		Name string `json:"name"`
		Existing struct {
			AccountID string `json:"account_id"`
			VehicleID string `json:"vehicle_id"`
			RaceID string `json:"race_id"`
		} `json:"existing"`
		Request struct {
			AccountID string `json:"account_id"`
			VehicleID string `json:"vehicle_id"`
			RaceID string `json:"race_id"`
		} `json:"request"`
		Matches bool `json:"matches"`
	} `json:"start_cases"`
	CheckpointCases []struct {
		Name string `json:"name"`
		Existing struct {
			RaceInstanceID string `json:"race_instance_id"`
			CheckpointIndex int `json:"checkpoint_index"`
			ElapsedMS int64 `json:"elapsed_ms"`
		} `json:"existing"`
		Request struct {
			RaceInstanceID string `json:"race_instance_id"`
			CheckpointIndex int `json:"checkpoint_index"`
			ElapsedMS int64 `json:"elapsed_ms"`
		} `json:"request"`
		Matches bool `json:"matches"`
	} `json:"checkpoint_cases"`
	FinishCases []struct {
		Name string `json:"name"`
		Existing struct {
			RaceInstanceID string `json:"race_instance_id"`
			CheckpointCount int `json:"checkpoint_count"`
			FinishElapsedMS int64 `json:"finish_elapsed_ms"`
		} `json:"existing"`
		Request struct {
			RaceInstanceID string `json:"race_instance_id"`
			CheckpointCount int `json:"checkpoint_count"`
			FinishElapsedMS int64 `json:"finish_elapsed_ms"`
		} `json:"request"`
		Matches bool `json:"matches"`
	} `json:"finish_cases"`
}

func loadRaceOperationVectors(t *testing.T) raceOperationVectors {
	t.Helper()
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("resolve race operation parity test path")
	}
	path := filepath.Clean(filepath.Join(filepath.Dir(filename), "../../../../tests/race-operation-parity-v4.2.json"))
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read race operation parity vectors: %v", err)
	}
	var vectors raceOperationVectors
	if err := json.Unmarshal(raw, &vectors); err != nil {
		t.Fatalf("decode race operation parity vectors: %v", err)
	}
	return vectors
}

func TestRaceOperationReplayParity(t *testing.T) {
	vectors := loadRaceOperationVectors(t)
	if len(vectors.StartCases) == 0 || len(vectors.CheckpointCases) == 0 || len(vectors.FinishCases) == 0 {
		t.Fatal("race operation parity vectors must cover start, checkpoint and finish")
	}
	for _, vector := range vectors.StartCases {
		vector := vector
		t.Run("start/"+vector.Name, func(t *testing.T) {
			existing := RaceInstance{AccountID: vector.Existing.AccountID, VehicleID: vector.Existing.VehicleID, RaceID: vector.Existing.RaceID}
			actual := RaceStartReplayMatches(existing, vector.Request.AccountID, vector.Request.VehicleID, vector.Request.RaceID)
			if actual != vector.Matches {
				t.Fatalf("start replay match=%v, want %v", actual, vector.Matches)
			}
		})
	}
	for _, vector := range vectors.CheckpointCases {
		vector := vector
		t.Run("checkpoint/"+vector.Name, func(t *testing.T) {
			actual := RaceCheckpointReplayMatches(vector.Existing.RaceInstanceID, vector.Existing.CheckpointIndex, vector.Existing.ElapsedMS, vector.Request.RaceInstanceID, vector.Request.CheckpointIndex, vector.Request.ElapsedMS)
			if actual != vector.Matches {
				t.Fatalf("checkpoint replay match=%v, want %v", actual, vector.Matches)
			}
		})
	}
	for _, vector := range vectors.FinishCases {
		vector := vector
		t.Run("finish/"+vector.Name, func(t *testing.T) {
			existing := RaceResult{RaceInstanceID: vector.Existing.RaceInstanceID, CheckpointCount: vector.Existing.CheckpointCount, FinishElapsedMS: vector.Existing.FinishElapsedMS}
			actual := RaceFinishReplayMatches(existing, vector.Request.RaceInstanceID, vector.Request.CheckpointCount, vector.Request.FinishElapsedMS)
			if actual != vector.Matches {
				t.Fatalf("finish replay match=%v, want %v", actual, vector.Matches)
			}
		})
	}
}
