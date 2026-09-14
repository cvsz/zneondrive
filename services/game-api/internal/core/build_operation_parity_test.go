package core

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

type buildOperationReplayVector struct {
	Name                   string `json:"name"`
	ExistingOwner          string `json:"existing_owner"`
	ExistingVehicle        string `json:"existing_vehicle"`
	ExistingResultRevision int    `json:"existing_result_revision"`
	ExistingHash           string `json:"existing_hash"`
	Owner                  string `json:"owner"`
	Vehicle                string `json:"vehicle"`
	ExpectedRevision       int    `json:"expected_revision"`
	Hash                   string `json:"hash"`
	Allow                  bool   `json:"allow"`
}

func TestBuildOperationReplayParityVectors(t *testing.T) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("resolve test path")
	}
	path := filepath.Join(filepath.Dir(file), "..", "..", "..", "..", "tests", "build-operation-replay-parity-v4.4.json")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read parity vectors: %v", err)
	}
	var fixture struct {
		Cases []buildOperationReplayVector `json:"cases"`
	}
	if err := json.Unmarshal(data, &fixture); err != nil {
		t.Fatalf("decode parity vectors: %v", err)
	}
	for _, tc := range fixture.Cases {
		t.Run(tc.Name, func(t *testing.T) {
			got := ValidateBuildOperationReplay(BuildOperationReplay{
				OwnerCharacterID: tc.ExistingOwner,
				VehicleID: tc.ExistingVehicle,
				ResultRevision: tc.ExistingResultRevision,
				ValidationHash: tc.ExistingHash,
			}, tc.Owner, tc.Vehicle, tc.ExpectedRevision, tc.Hash)
			if got != tc.Allow {
				t.Fatalf("allow=%v want %v", got, tc.Allow)
			}
		})
	}
}
