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
	accountID string
	tokenHash string
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
	if tokenHash != f.tokenHash {
		return core.Snapshot{}, errUnauthorized
	}
	return core.Snapshot{AccountID: f.accountID, CharacterID: "char_test", VehicleID: "veh_test"}, nil
}

func (f *fakeStore) CompleteQuest(context.Context, string, string, string) (core.Snapshot, core.RewardReceipt, error) {
	return core.Snapshot{}, core.RewardReceipt{}, nil
}

func (f *fakeStore) ReviseBuild(context.Context, string, string, int, []string, string) (core.Snapshot, error) {
	return core.Snapshot{}, nil
}

var errUnauthorized = &fakeUnauthorized{}

type fakeUnauthorized struct{}

func (*fakeUnauthorized) Error() string { return "unauthorized" }

func TestBootstrapCreatesOpaqueSessionAndResumeKey(t *testing.T) {
	s := &fakeStore{}
	handler := New(s)
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
	handler := New(&fakeStore{})
	req := httptest.NewRequest(http.MethodGet, "/v1/state", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rec.Code)
	}
}
