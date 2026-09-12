package httpapi

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestTelemetryBucketDoesNotExposeRawIdentifier(t *testing.T) {
	raw := "raceinst_super-sensitive-example"
	bucket := telemetryBucket("race", raw)
	if strings.Contains(bucket, raw) {
		t.Fatal("telemetry bucket must not expose raw identifier")
	}
	if !strings.HasPrefix(bucket, "race:") {
		t.Fatalf("expected scoped telemetry bucket, got %q", bucket)
	}
	if bucket != telemetryBucket("race", raw) {
		t.Fatal("telemetry bucket must be deterministic for correlation")
	}
	if bucket == telemetryBucket("race", raw+"-other") {
		t.Fatal("different identifiers must not collapse into one telemetry bucket")
	}
}

func TestInternalRouteScopeDoesNotLeakRaceInstancePath(t *testing.T) {
	if got := internalRouteScope("/v1/internal/races/raceinst_secret/checkpoints"); got != "race-checkpoint" {
		t.Fatalf("unexpected checkpoint route scope %q", got)
	}
	if got := internalRouteScope("/v1/internal/races/raceinst_secret/finish"); got != "race-finish" {
		t.Fatalf("unexpected finish route scope %q", got)
	}
	if got := internalRouteScope("/v1/internal/game-tickets/redeem"); got != "ticket-redeem" {
		t.Fatalf("unexpected ticket route scope %q", got)
	}
}

func TestSecurityTelemetryLogsRaceConflictWithoutRawRaceID(t *testing.T) {
	var logs bytes.Buffer
	previous := log.Writer()
	log.SetOutput(&logs)
	defer log.SetOutput(previous)

	next := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		writeError(w, http.StatusConflict, "race_event_out_of_order")
	})
	handler := NewSecurityTelemetryHandler(next)
	req := httptest.NewRequest(http.MethodPost, "/v1/internal/races/raceinst_secret/checkpoints", nil)
	req.RemoteAddr = "203.0.113.44:5000"
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusConflict {
		t.Fatalf("expected downstream 409 to be preserved, got %d", rec.Code)
	}
	text := logs.String()
	if !strings.Contains(text, "security_event=race_integrity_rejected") || !strings.Contains(text, "reason=authoritative_state_conflict") {
		t.Fatalf("missing integrity telemetry: %s", text)
	}
	if strings.Contains(text, "raceinst_secret") || strings.Contains(text, "203.0.113.44") {
		t.Fatalf("telemetry leaked raw race/peer identifier: %s", text)
	}
}

func TestSecurityTelemetryLogsGameServerAuthReject(t *testing.T) {
	var logs bytes.Buffer
	previous := log.Writer()
	log.SetOutput(&logs)
	defer log.SetOutput(previous)

	next := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		writeError(w, http.StatusUnauthorized, "unauthorized_game_server")
	})
	handler := NewSecurityTelemetryHandler(next)
	req := httptest.NewRequest(http.MethodPost, "/v1/internal/races/raceinst_secret/finish", nil)
	req.RemoteAddr = "198.51.100.9:6000"
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	text := logs.String()
	if !strings.Contains(text, "security_event=game_server_auth_rejected") || !strings.Contains(text, "route=race-finish") {
		t.Fatalf("missing game-server auth telemetry: %s", text)
	}
	if strings.Contains(text, "198.51.100.9") || strings.Contains(text, "raceinst_secret") {
		t.Fatalf("auth telemetry leaked raw identifiers: %s", text)
	}
}
