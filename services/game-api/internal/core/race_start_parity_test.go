package core

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

type raceStartVectors struct {
	Cases []struct {
		Name                string `json:"name"`
		RaceID              string `json:"race_id"`
		Roadworthy          bool   `json:"roadworthy"`
		BuildRevision       int    `json:"build_revision"`
		BuildValidationHash string `json:"build_validation_hash"`
		Valid               bool   `json:"valid"`
	} `json:"cases"`
}

func TestRaceStartBindingParity(t *testing.T) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("resolve race-start parity test path")
	}
	path := filepath.Clean(filepath.Join(filepath.Dir(filename), "../../../../tests/race-start-parity-v4.1.json"))
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read race-start parity vectors: %v", err)
	}
	var vectors raceStartVectors
	if err := json.Unmarshal(raw, &vectors); err != nil {
		t.Fatalf("decode race-start parity vectors: %v", err)
	}
	if len(vectors.Cases) == 0 {
		t.Fatal("race-start parity vectors must not be empty")
	}
	for _, vector := range vectors.Cases {
		vector := vector
		t.Run(vector.Name, func(t *testing.T) {
			_, err := ValidateRaceStartBinding(vector.RaceID, vector.Roadworthy, vector.BuildRevision, vector.BuildValidationHash)
			if vector.Valid && err != nil {
				t.Fatalf("expected race-start acceptance: %v", err)
			}
			if !vector.Valid && err == nil {
				t.Fatal("expected race-start rejection")
			}
		})
	}
}
