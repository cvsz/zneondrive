//go:build integration

package httpapi

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/cvsz/zneondrive/services/game-api/internal/core"
	"github.com/cvsz/zneondrive/services/game-api/internal/store"
)

const integrationServerKey = "integration-game-server-key-32-characters"

type integrationBootstrap struct {
	SessionToken string        `json:"session_token"`
	ResumeKey    string        `json:"resume_key"`
	Snapshot     core.Snapshot `json:"snapshot"`
}

type integrationTicket struct {
	Ticket string `json:"ticket"`
}

func TestHTTPReconnectAndOneTimeGameplayTicket(t *testing.T) {
	databaseURL := os.Getenv("TEST_DATABASE_URL")
	if databaseURL == "" {
		t.Skip("TEST_DATABASE_URL is not set")
	}

	ctx := context.Background()
	db, err := store.OpenPostgres(ctx, databaseURL)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if err = db.EnsureSchema(ctx); err != nil {
		t.Fatal(err)
	}

	server := httptest.NewServer(New(db, integrationServerKey))
	defer server.Close()

	first := integrationBootstrap{}
	doJSON(t, http.MethodPost, server.URL+"/v1/sessions/bootstrap", "", "", map[string]any{}, http.StatusCreated, &first)
	if first.ResumeKey == "" || first.SessionToken == "" || first.Snapshot.AccountID == "" {
		t.Fatalf("invalid bootstrap response: %+v", first)
	}

	ticket := integrationTicket{}
	doJSON(t, http.MethodPost, server.URL+"/v1/game-tickets", first.SessionToken, "", nil, http.StatusCreated, &ticket)
	if len(ticket.Ticket) != 64 {
		t.Fatalf("expected 64-character gameplay ticket, got %d", len(ticket.Ticket))
	}

	doJSON(
		t,
		http.MethodPost,
		server.URL+"/v1/internal/game-tickets/redeem",
		"",
		integrationServerKey,
		map[string]any{"ticket": ticket.Ticket},
		http.StatusOK,
		&core.Snapshot{},
	)

	doJSON(
		t,
		http.MethodPost,
		server.URL+"/v1/internal/game-tickets/redeem",
		"",
		integrationServerKey,
		map[string]any{"ticket": ticket.Ticket},
		http.StatusUnauthorized,
		nil,
	)

	var quest struct {
		Snapshot core.Snapshot `json:"snapshot"`
	}
	doJSON(
		t,
		http.MethodPost,
		server.URL+"/v1/quests/MQ001/complete",
		first.SessionToken,
		"",
		map[string]any{"operation_id": "http-e2e-mq001-" + first.Snapshot.AccountID},
		http.StatusOK,
		&quest,
	)
	if len(quest.Snapshot.CompletedQuests) != 1 || quest.Snapshot.CompletedQuests[0] != "MQ001" {
		t.Fatalf("expected MQ001 completion, got %+v", quest.Snapshot.CompletedQuests)
	}

	second := integrationBootstrap{}
	doJSON(
		t,
		http.MethodPost,
		server.URL+"/v1/sessions/bootstrap",
		"",
		"",
		map[string]any{"resume_key": first.ResumeKey},
		http.StatusCreated,
		&second,
	)
	if second.Snapshot.AccountID != first.Snapshot.AccountID {
		t.Fatalf("resume changed account: %s != %s", second.Snapshot.AccountID, first.Snapshot.AccountID)
	}
	if second.SessionToken == "" || second.SessionToken == first.SessionToken {
		t.Fatal("resume should issue a fresh session token")
	}
	if len(second.Snapshot.CompletedQuests) != 1 || second.Snapshot.CompletedQuests[0] != "MQ001" {
		t.Fatalf("quest state did not survive reconnect: %+v", second.Snapshot.CompletedQuests)
	}

	var state core.Snapshot
	doJSON(t, http.MethodGet, server.URL+"/v1/state", second.SessionToken, "", nil, http.StatusOK, &state)
	if state.VehicleID != first.Snapshot.VehicleID || state.Money != quest.Snapshot.Money {
		t.Fatalf("durable state mismatch after reconnect: %+v", state)
	}

	secondTicket := integrationTicket{}
	doJSON(t, http.MethodPost, server.URL+"/v1/game-tickets", second.SessionToken, "", nil, http.StatusCreated, &secondTicket)

	var redeemed core.Snapshot
	doJSON(
		t,
		http.MethodPost,
		server.URL+"/v1/internal/game-tickets/redeem",
		"",
		integrationServerKey,
		map[string]any{"ticket": secondTicket.Ticket},
		http.StatusOK,
		&redeemed,
	)
	if redeemed.VehicleID != state.VehicleID || redeemed.ActiveBuildRevision != state.ActiveBuildRevision {
		t.Fatalf("redeemed gameplay snapshot differs from session state: %+v vs %+v", redeemed, state)
	}
}

func doJSON(
	t *testing.T,
	method string,
	url string,
	bearer string,
	serverKey string,
	body any,
	wantStatus int,
	target any,
) {
	t.Helper()

	var reader io.Reader
	if body != nil {
		payload, err := json.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}
		reader = bytes.NewReader(payload)
	}

	req, err := http.NewRequest(method, url, reader)
	if err != nil {
		t.Fatal(err)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if bearer != "" {
		req.Header.Set("Authorization", "Bearer "+bearer)
	}
	if serverKey != "" {
		req.Header.Set("X-Game-Server-Key", serverKey)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	payload, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != wantStatus {
		t.Fatalf("%s %s: expected %d, got %d: %s", method, url, wantStatus, resp.StatusCode, string(payload))
	}
	if target != nil && len(payload) > 0 {
		if err = json.Unmarshal(payload, target); err != nil {
			t.Fatalf("decode %s: %v (%s)", url, err, string(payload))
		}
	}
}
