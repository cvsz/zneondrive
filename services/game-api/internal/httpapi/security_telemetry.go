package httpapi

import (
	"crypto/sha256"
	"encoding/hex"
	"log"
	"net/http"
	"strings"
)

func telemetryBucket(scope, value string) string {
	digest := sha256.Sum256([]byte(strings.TrimSpace(value)))
	return scope + ":" + hex.EncodeToString(digest[:12])
}

type statusCaptureWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusCaptureWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *statusCaptureWriter) Write(p []byte) (int, error) {
	if w.status == 0 {
		w.WriteHeader(http.StatusOK)
	}
	return w.ResponseWriter.Write(p)
}

func NewSecurityTelemetryHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capture := &statusCaptureWriter{ResponseWriter: w}
		next.ServeHTTP(capture, r)
		status := capture.status
		if status == 0 {
			status = http.StatusOK
		}
		emitSecurityTelemetry(r, status)
	})
}

func emitSecurityTelemetry(r *http.Request, status int) {
	scope := internalRouteScope(r.URL.Path)
	if scope == "internal" {
		return
	}
	peer := "unknown"
	if ip := remoteIP(r.RemoteAddr); ip != nil {
		peer = ip.String()
	}
	if status == http.StatusUnauthorized {
		log.Printf("security_event=game_server_auth_rejected route=%s peer=%s", scope, telemetryBucket("peer", peer))
		return
	}
	if scope != "race-checkpoint" && scope != "race-finish" {
		return
	}
	if status != http.StatusBadRequest && status != http.StatusConflict {
		return
	}
	rawRaceID := raceInstanceIDFromPath(r.URL.Path)
	reason := "invalid_request"
	if status == http.StatusConflict {
		reason = "authoritative_state_conflict"
	}
	log.Printf("security_event=race_integrity_rejected action=%s reason=%s race=%s peer=%s status=%d",
		scope, reason, telemetryBucket("race", rawRaceID), telemetryBucket("peer", peer), status)
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

func raceInstanceIDFromPath(path string) string {
	const prefix = "/v1/internal/races/"
	if !strings.HasPrefix(path, prefix) {
		return "unknown"
	}
	remainder := strings.TrimPrefix(path, prefix)
	if idx := strings.IndexByte(remainder, '/'); idx >= 0 {
		remainder = remainder[:idx]
	}
	if remainder == "" {
		return "unknown"
	}
	return remainder
}
