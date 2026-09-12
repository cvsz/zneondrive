package httpapi

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/cvsz/zneondrive/services/game-api/internal/core"
	"github.com/cvsz/zneondrive/services/game-api/internal/store"
)

func telemetryBucket(scope, value string) string {
	digest := sha256.Sum256([]byte(strings.TrimSpace(value)))
	return scope + ":" + hex.EncodeToString(digest[:12])
}

func raceIntegrityReason(err error) string {
	switch {
	case errors.Is(err, store.ErrRaceOrder):
		return "sequence_or_time_order"
	case errors.Is(err, store.ErrOperationKey):
		return "operation_replay_conflict"
	case errors.Is(err, core.ErrInvalidRace):
		return "invalid_race_domain"
	default:
		return ""
	}
}

func logRaceIntegrityRejection(action, raceInstanceID, operationID string, checkpointIndex int, elapsedMS int64, err error) {
	reason := raceIntegrityReason(err)
	if reason == "" {
		return
	}
	log.Printf(
		"security_event=race_integrity_rejected action=%s reason=%s race=%s operation=%s checkpoint_index=%d elapsed_ms=%d",
		action,
		reason,
		telemetryBucket("race", raceInstanceID),
		telemetryBucket("operation", operationID),
		checkpointIndex,
		elapsedMS,
	)
}

func internalRouteScope(path string) string {
	switch {
	case path == "/v1/internal/game-tickets/redeem":
		return "ticket-redeem"
	case path == "/v1/internal/races/start":
		return "race-start"
	case strings.HasPrefix(path, "/v1/internal/races/") && strings.HasSuffix(path, "/checkpoints"):
		return "race-checkpoint"
	case strings.HasPrefix(path, "/v1/internal/races/") && strings.HasSuffix(path, "/finish"):
		return "race-finish"
	default:
		return "internal"
	}
}

func logGameServerAuthRejection(r *http.Request) {
	peer := "unknown"
	if ip := remoteIP(r.RemoteAddr); ip != nil {
		peer = ip.String()
	}
	log.Printf(
		"security_event=game_server_auth_rejected route=%s peer=%s",
		internalRouteScope(r.URL.Path),
		telemetryBucket("peer", peer),
	)
}
