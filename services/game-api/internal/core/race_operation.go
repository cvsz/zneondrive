package core

// RaceStartReplayMatches reports whether an existing operation-id binding represents
// the same authoritative race-start mutation. The operation id lookup and durable
// uniqueness remain PostgreSQL responsibilities; this helper only centralizes the
// deterministic payload equality contract used for safe retries.
func RaceStartReplayMatches(existing RaceInstance, accountID, vehicleID, raceID string) bool {
	normalizedRaceID, err := NormalizeRaceID(raceID)
	if err != nil {
		return false
	}
	return existing.AccountID == accountID && existing.VehicleID == vehicleID && existing.RaceID == normalizedRaceID
}

// RaceCheckpointReplayMatches binds checkpoint retries to the exact race instance,
// checkpoint cursor and elapsed time originally associated with an operation id.
func RaceCheckpointReplayMatches(existingRaceInstanceID string, existingCheckpointIndex int, existingElapsedMS int64, raceInstanceID string, checkpointIndex int, elapsedMS int64) bool {
	return existingRaceInstanceID == raceInstanceID && existingCheckpointIndex == checkpointIndex && existingElapsedMS == elapsedMS
}

// RaceFinishReplayMatches binds finish retries to the exact race instance, checkpoint
// count and finish elapsed time originally associated with an operation id.
func RaceFinishReplayMatches(existing RaceResult, raceInstanceID string, checkpointCount int, finishElapsedMS int64) bool {
	return existing.RaceInstanceID == raceInstanceID && existing.CheckpointCount == checkpointCount && existing.FinishElapsedMS == finishElapsedMS
}
