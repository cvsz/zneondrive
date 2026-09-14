package core

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

type entitlementParityVectors struct {
	SchemaVersion int `json:"schema_version"`
	GarageCapacityCases []struct {
		Name string `json:"name"`
		CurrentVehicles int `json:"current_vehicles"`
		GarageSlots int `json:"garage_slots"`
		GrantAllowed bool `json:"grant_allowed"`
	} `json:"garage_capacity_cases"`
	RoutineDeleteCases []struct {
		Name string `json:"name"`
		StarterLineage bool `json:"starter_lineage"`
		DeleteAllowed bool `json:"delete_allowed"`
	} `json:"routine_delete_cases"`
}

func loadEntitlementParityVectors(t *testing.T) entitlementParityVectors {
	t.Helper()
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("resolve entitlement parity test path")
	}
	path := filepath.Clean(filepath.Join(filepath.Dir(filename), "../../../../tests/entitlement-parity-vectors.json"))
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read entitlement parity vectors: %v", err)
	}
	var vectors entitlementParityVectors
	if err := json.Unmarshal(raw, &vectors); err != nil {
		t.Fatalf("decode entitlement parity vectors: %v", err)
	}
	if vectors.SchemaVersion != 1 {
		t.Fatalf("unsupported entitlement parity schema version: %d", vectors.SchemaVersion)
	}
	return vectors
}

func TestEntitlementParityVectorIntegrity(t *testing.T) {
	vectors := loadEntitlementParityVectors(t)
	if len(vectors.GarageCapacityCases) == 0 {
		t.Fatal("garage_capacity_cases must not be empty")
	}
	if len(vectors.RoutineDeleteCases) == 0 {
		t.Fatal("routine_delete_cases must not be empty")
	}
	seen := map[string]struct{}{}
	for _, vector := range vectors.GarageCapacityCases {
		if vector.Name == "" {
			t.Fatal("garage capacity vector name must not be empty")
		}
		if _, exists := seen[vector.Name]; exists {
			t.Fatalf("duplicate entitlement parity vector name: %s", vector.Name)
		}
		seen[vector.Name] = struct{}{}
		if vector.CurrentVehicles < 0 || vector.GarageSlots < 1 {
			t.Fatalf("invalid garage capacity vector: %+v", vector)
		}
	}
	for _, vector := range vectors.RoutineDeleteCases {
		if vector.Name == "" {
			t.Fatal("routine delete vector name must not be empty")
		}
		if _, exists := seen[vector.Name]; exists {
			t.Fatalf("duplicate entitlement parity vector name: %s", vector.Name)
		}
		seen[vector.Name] = struct{}{}
	}
}

func TestReferenceParityGarageCapacity(t *testing.T) {
	vectors := loadEntitlementParityVectors(t)
	for _, vector := range vectors.GarageCapacityCases {
		vector := vector
		t.Run(vector.Name, func(t *testing.T) {
			err := ValidateGarageGrantCapacity(vector.CurrentVehicles, vector.GarageSlots)
			if vector.GrantAllowed && err != nil {
				t.Fatalf("expected grant to be allowed: %v", err)
			}
			if !vector.GrantAllowed && err == nil {
				t.Fatal("expected garage capacity to reject vehicle grant")
			}
		})
	}
}

func TestReferenceParityStarterDeletionProtection(t *testing.T) {
	vectors := loadEntitlementParityVectors(t)
	for _, vector := range vectors.RoutineDeleteCases {
		vector := vector
		t.Run(vector.Name, func(t *testing.T) {
			err := ValidateRoutineVehicleDeletion(vector.StarterLineage)
			if vector.DeleteAllowed && err != nil {
				t.Fatalf("expected routine deletion to be allowed: %v", err)
			}
			if !vector.DeleteAllowed && err == nil {
				t.Fatal("expected starter-lineage deletion to be rejected")
			}
		})
	}
}
