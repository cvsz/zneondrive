package core

// BuildOperationReplay captures the durable fields that must match before an
// existing rebuild operation may be treated as an idempotent replay.
type BuildOperationReplay struct {
	OwnerCharacterID string
	VehicleID        string
	ResultRevision   int
	ValidationHash   string
}

// ValidateBuildOperationReplay rejects operation-key reuse unless the durable
// receipt belongs to the same character/vehicle and represents exactly the
// revision/hash implied by the caller's original request.
func ValidateBuildOperationReplay(existing BuildOperationReplay, ownerCharacterID, vehicleID string, expectedRevision int, validationHash string) bool {
	if ownerCharacterID == "" || vehicleID == "" || validationHash == "" || expectedRevision < 1 {
		return false
	}
	return existing.OwnerCharacterID == ownerCharacterID &&
		existing.VehicleID == vehicleID &&
		existing.ResultRevision == expectedRevision+1 &&
		existing.ValidationHash == validationHash
}
