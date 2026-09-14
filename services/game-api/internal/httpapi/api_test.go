package httpapi

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/cvsz/zneondrive/services/game-api/internal/core"
)

type fakeStore struct {
	accountID        string
	tokenHash        string
	ticketHash       string
	characterID      string
	questOperationID string
}

func (f *fakeStore) Ping(context.Context) error { return nil }
func (f *fakeStore) Close()                     {}

func (f *fakeStore) Bootstrap(context.Context, string) (core.Snapshot, bool, error) {
	return core.Snapshot{
		AccountID:           "acct_test",
		CharacterID:         "char_test",
		VehicleID:           "veh_test",
		ActiveBuildRevision: 1,
		ActivePartIDs:       []string{"part_chassis_starter_prototype"},
		StarterLineage:      true,
	}, true, nil
}

func (f *fakeStore) CreateSession(_ context.Context, accountID, tokenHash string, _ time.Time) error {
	f.accountID = accountID
	f.tokenHash = tokenHash
	return nil
}

func (f *fakeStore) SnapshotBySession(_ context.Context, tokenHash string) (core.Snapshot, error) {
	characterID := f.characterID
	if characterID == "" {
		characterID = "char_test"
	}
	return core.Snapshot{AccountID: f.accountID, CharacterID: characterID, VehicleID: "veh_test"}, nil
}

func (f *fakeStore) IssueGameTicket(_ context.Context, _ string, ticketHash string, _ time.Time) error {
	f.ticketHash = ticketHash
	return nil
}

func (f *fakeStore) RedeemGameTicket(_ context.Context, ticketHash string) (core.Snapshot, error) {
	return core.Snapshot{AccountID: f.accountID, CharacterID: "char_test", VehicleID: "veh_test", ActiveBuildRevision: 1}, nil
}

func (f *fakeStore) CompleteQuest(_ context.Context, _ string, _ string, operationID string) (core.Snapshot, core.RewardReceipt, error) {
	f.questOperationID = operationID
	return core.Snapshot{}, core.RewardReceipt{}, nil
}

func (f *fakeStore) ReviseBuild(context.Context, string, string, int, []string, string) (core.Snapshot, error) {
	return core.Snapshot{}, nil
}

func (f *fakeStore) StartRace(_ context.Context, accountID, vehicleID, raceID, _ string) (core.RaceInstance, error) {
	return core.RaceInstance{
		RaceInstanceID: "raceinst_test",
		RaceID: raceID,
		AccountID: accountID,
		CharacterID: "char_test",
		VehicleID: vehicleID,
		BuildRevision: 1,
		BuildValidationHash: "build_hash",
		State: "active",
	}, nil
}

func (f *fakeStore) RecordRaceCheckpoint(_ context.Context, raceInstanceID string, checkpointIndex int, elapsedMS int64, _ string) (core.RaceInstance, error) {
	return core.RaceInstance{RaceInstanceID: raceInstanceID, State: "active", NextCheckpoint: checkpointIndex + 1, LastElapsedMS: elapsedMS}, nil
}

func (f *fakeStore) FinishRace(_ context.Context, raceInstanceID string, checkpointCount int, finishElapsedMS int64, _ string) (core.RaceResult, error) {
	return core.RaceResult{RaceInstanceID: raceInstanceID, CheckpointCount: checkpointCount, FinishElapsedMS: finishElapsedMS, ResultHash: "result_hash"}, nil
}

func TestBootstrapCreatesOpaqueSessionAndResumeKey(t *testing.T) {
	s := &fakeStore{}
	handler := New(s, "test-game-server-key-32-characters-minimum")
	req := httptest.NewRequest(http.MethodPost, "/v1/sessions/bootstrap", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", rec.Code, rec.Body.String())
	}

	var body bootstrapResponse
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}
	if body.SessionToken == "" || body.ResumeKey == "" {
		t.Fatal("expected generated session and resume key")
	}
	if s.tokenHash != core.HashSecret(body.SessionToken) {
		t.Fatal("store must receive only hashed session token")
	}
}

func TestStateRejectsMissingBearer(t *testing.T) {
	handler := New(&fakeStore{}, "test-game-server-key-32-characters-minimum")
	req := httptest.NewRequest(http.MethodGet, "/v1/state", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rec.Code)
	}
}

func TestQuestOperationIDIsScopedToAuthoritativeCharacter(t *testing.T) {
	s := &fakeStore{accountID: "acct_test", characterID: "char_test"}
	handler := New(s, "test-game-server-key-32-characters-minimum")
	req := httptest.NewRequest(http.MethodPost, "/v1/quests/MQ001/complete", strings.NewReader(`{"operation_id":"client-op-1"}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer session-token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	expected, ok := core.ScopeQuestOperationID("char_test", "client-op-1")
	if !ok {
		t.Fatal("expected valid scoped operation id")
	}
	if s.questOperationID != expected {
		t.Fatalf("store received %q, expected %q", s.questOperationID, expected)
	}
	if s.questOperationID == "client-op-1" {
		t.Fatal("raw caller operation id must not be used as durable quest operation key")
	}
}

func TestGameplayTicketRequiresServerKey(t *testing.T) {
	handler := New(&fakeStore{}, "test-game-server-key-32-characters-minimum")
	req := httptest.NewRequest(http.MethodPost, "/v1/internal/game-tickets/redeem", strings.NewReader(`{"ticket":"0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Game-Server-Key", "wrong-key")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rec.Code)
	}
}

func TestRaceStartRequiresServerKey(t *testing.T) {
	handler := New(&fakeStore{}, "test-game-server-key-32-characters-minimum")
	body := `{"account_id":"acct_test","vehicle_id":"veh_test","race_id":"race_first_ignition","operation_id":"op-1"}`
	req := httptest.NewRequest(http.MethodPost, "/v1/internal/races/start", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rec.Code)
	}
}

func TestRaceStartAcceptsDedicatedServerOnly(t *testing.T) {
	handler := New(&fakeStore{}, "test-game-server-key-32-characters-minimum")
	body := `{"account_id":"acct_test","vehicle_id":"veh_test","race_id":"race_first_ignition","operation_id":"op-1"}`
	req := httptest.NewRequest(http.MethodPost, "/v1/internal/races/start", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Game-Server-Key", "test-game-server-key-32-characters-minimum")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", rec.Code, rec.Body.String())
	}
}
