package httpapi

import (
	"errors"
	"strings"
	"testing"

	"github.com/cvsz/zneondrive/services/game-api/internal/core"
	"github.com/cvsz/zneondrive/services/game-api/internal/store"
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

func TestRaceIntegrityReasonClassifiesOnlyIntegritySignals(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want string
	}{
		{name: "order", err: store.ErrRaceOrder, want: "sequence_or_time_order"},
		{name: "operation", err: store.ErrOperationKey, want: "operation_replay_conflict"},
		{name: "invalid", err: core.ErrInvalidRace, want: "invalid_race_domain"},
		{name: "wrapped", err: errors.Join(errors.New("outer"), store.ErrRaceOrder), want: "sequence_or_time_order"},
		{name: "not found", err: store.ErrNotFound, want: ""},
		{name: "backend", err: errors.New("database unavailable"), want: ""},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := raceIntegrityReason(tc.err); got != tc.want {
				t.Fatalf("raceIntegrityReason(%v)=%q want %q", tc.err, got, tc.want)
			}
		})
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
