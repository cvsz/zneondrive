package core

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"strings"
)

var ErrInvalidRace = errors.New("invalid race input")

type RaceInstance struct {
	RaceInstanceID     string `json:"race_instance_id"`
	RaceID             string `json:"race_id"`
	AccountID          string `json:"account_id"`
	CharacterID        string `json:"character_id"`
	VehicleID          string `json:"vehicle_id"`
	BuildRevision      int    `json:"build_revision"`
	BuildValidationHash string `json:"build_validation_hash"`
	State              string `json:"state"`
	NextCheckpoint     int    `json:"next_checkpoint"`
	LastElapsedMS      int64  `json:"last_elapsed_ms"`
}

type RaceResult struct {
	RaceInstanceID string `json:"race_instance_id"`
	RaceID         string `json:"race_id"`
	VehicleID      string `json:"vehicle_id"`
	BuildRevision  int    `json:"build_revision"`
	CheckpointCount int   `json:"checkpoint_count"`
	FinishElapsedMS int64 `json:"finish_elapsed_ms"`
	ResultHash     string `json:"result_hash"`
}

func NormalizeRaceID(raceID string) (string, error) {
	raceID = strings.TrimSpace(raceID)
	if len(raceID) < 6 || len(raceID) > 128 || !strings.HasPrefix(raceID, "race_") {
		return "", ErrInvalidRace
	}
	for _, r := range raceID {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '_' || r == '-' {
			continue
		}
		return "", ErrInvalidRace
	}
	return raceID, nil
}

func ValidateRaceCheckpoint(index int, elapsedMS int64) error {
	if index < 0 || index > 1024 || elapsedMS <= 0 || elapsedMS > 24*60*60*1000 {
		return ErrInvalidRace
	}
	return nil
}

func RaceResultHash(instance RaceInstance, checkpointCount int, finishElapsedMS int64) (string, error) {
	if checkpointCount < 1 || checkpointCount > 1025 || finishElapsedMS <= 0 {
		return "", ErrInvalidRace
	}
	payload, err := json.Marshal(struct {
		RaceInstanceID      string `json:"race_instance_id"`
		RaceID              string `json:"race_id"`
		AccountID           string `json:"account_id"`
		CharacterID         string `json:"character_id"`
		VehicleID           string `json:"vehicle_id"`
		BuildRevision       int    `json:"build_revision"`
		BuildValidationHash string `json:"build_validation_hash"`
		CheckpointCount     int    `json:"checkpoint_count"`
		FinishElapsedMS     int64  `json:"finish_elapsed_ms"`
	}{
		RaceInstanceID: instance.RaceInstanceID,
		RaceID: instance.RaceID,
		AccountID: instance.AccountID,
		CharacterID: instance.CharacterID,
		VehicleID: instance.VehicleID,
		BuildRevision: instance.BuildRevision,
		BuildValidationHash: instance.BuildValidationHash,
		CheckpointCount: checkpointCount,
		FinishElapsedMS: finishElapsedMS,
	})
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(payload)
	return hex.EncodeToString(sum[:]), nil
}
