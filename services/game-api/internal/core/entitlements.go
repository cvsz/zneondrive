package core

import "fmt"

// ValidateGarageGrantCapacity mirrors the reference authority rule used when a
// character is granted another vehicle. Garage capacity is an entitlement
// boundary only; it must not affect competitive build identity or race outcome.
func ValidateGarageGrantCapacity(currentVehicles, garageSlots int) error {
	if garageSlots < 1 {
		return fmt.Errorf("garage slots must be >= 1")
	}
	if currentVehicles < 0 {
		return fmt.Errorf("current vehicle count must be >= 0")
	}
	if currentVehicles >= garageSlots {
		return fmt.Errorf("garage slot capacity exceeded")
	}
	return nil
}

// ValidateRoutineVehicleDeletion preserves the starter-lineage vehicle from
// ordinary deletion paths. Any future exceptional recovery/admin path must be
// separately authorized and audited instead of weakening this invariant.
func ValidateRoutineVehicleDeletion(starterLineage bool) error {
	if starterLineage {
		return fmt.Errorf("starter-lineage vehicle cannot be routine-deleted")
	}
	return nil
}
