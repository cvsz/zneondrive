//go:build integration

package httpapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/cvsz/zneondrive/services/game-api/internal/core"
	"github.com/cvsz/zneondrive/services/game-api/internal/store"
)

func TestIntegratedTrustBoundaryAbusePaths(t *testing.T) {
	databaseURL := os.Getenv("TEST_DATABASE_URL")
	redisAddr := os.Getenv("TEST_REDIS_ADDR")
	if databaseURL == "" || redisAddr == "" {
		t.Skip("TEST_DATABASE_URL and TEST_REDIS_ADDR are required")
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

	originalBootstrapPolicy := bootstrapRatePolicy
	bootstrapRatePolicy = rateLimitPolicy{Burst: 3, RefillEvery: time.Hour}
	t.Cleanup(func() { bootstrapRatePolicy = originalBootstrapPolicy })

	api := New(db, integrationServerKey)
	telemetry := NewSecurityTelemetryHandler(api)

	// Direct clients are not trusted to supply forwarding headers. Distinct spoofed
	// X-Forwarded-For values must therefore consume one shared socket-peer budget.
	direct := httptest.NewServer(NewDistributedRateLimitedHandlerWithTrustedProxies(telemetry, redisAddr, TrustedProxyPolicy{}))
	defer direct.Close()
	for i := 0; i < bootstrapRatePolicy.Burst+1; i++ {
		want := http.StatusCreated
		if i == bootstrapRatePolicy.Burst {
			want = http.StatusTooManyRequests
		}
		var out integrationBootstrap
		status, body := integratedJSONRequest(
			t,
			direct.URL+"/v1/sessions/bootstrap",
			http.MethodPost,
			"",
			"",
			fmt.Sprintf("198.51.100.%d", i+10),
			map[string]any{},
			&out,
		)
		if status != want {
			t.Fatalf("direct spoof attempt %d: expected %d, got %d: %s", i, want, status, body)
		}
	}

	proxyPolicy, err := ParseTrustedProxyCIDRs("127.0.0.1/32")
	if err != nil {
		t.Fatal(err)
	}
	trusted := httptest.NewServer(NewDistributedRateLimitedHandlerWithTrustedProxies(telemetry, redisAddr, proxyPolicy))
	defer trusted.Close()

	// A configured trusted ingress is allowed to preserve distinct real client
	// identities. Each client receives an independent bootstrap abuse budget.
	var bootstrap integrationBootstrap
	for i := 0; i < bootstrapRatePolicy.Burst+1; i++ {
		var out integrationBootstrap
		status, body := integratedJSONRequest(
			t,
			trusted.URL+"/v1/sessions/bootstrap",
			http.MethodPost,
			"",
			"",
			fmt.Sprintf("2001:db8:1::%x", i+1),
			map[string]any{},
			&out,
		)
		if status != http.StatusCreated {
			t.Fatalf("trusted ingress client %d: expected 201, got %d: %s", i, status, body)
		}
		if i == 0 {
			bootstrap = out
		}
	}
	if bootstrap.SessionToken == "" || bootstrap.Snapshot.AccountID == "" {
		t.Fatalf("invalid trusted-ingress bootstrap response: %+v", bootstrap)
	}

	// A forged bearer credential cannot become authoritative merely because the
	// request traverses the trusted ingress and distributed limiter.
	status, body := integratedJSONRequest(
		t,
		trusted.URL+"/v1/state",
		http.MethodGet,
		"forged-session-token",
		"",
		"2001:db8:1::1",
		nil,
		&core.Snapshot{},
	)
	if status != http.StatusUnauthorized {
		t.Fatalf("forged session token: expected 401, got %d: %s", status, body)
	}

	var ticket integrationTicket
	status, body = integratedJSONRequest(
		t,
		trusted.URL+"/v1/game-tickets",
		http.MethodPost,
		bootstrap.SessionToken,
		"",
		"2001:db8:1::1",
		nil,
		&ticket,
	)
	if status != http.StatusCreated || len(ticket.Ticket) != 64 {
		t.Fatalf("gameplay ticket issuance failed: status=%d ticket_len=%d body=%s", status, len(ticket.Ticket), body)
	}

	// The internal redemption route remains server-only even when all outer
	// middleware succeeds.
	status, body = integratedJSONRequest(
		t,
		trusted.URL+"/v1/internal/game-tickets/redeem",
		http.MethodPost,
		"",
		"wrong-game-server-key-that-is-long-enough",
		"2001:db8:ffff::10",
		map[string]any{"ticket": ticket.Ticket},
		&core.Snapshot{},
	)
	if status != http.StatusUnauthorized {
		t.Fatalf("wrong game-server key: expected 401, got %d: %s", status, body)
	}

	var redeemed core.Snapshot
	status, body = integratedJSONRequest(
		t,
		trusted.URL+"/v1/internal/game-tickets/redeem",
		http.MethodPost,
		"",
		integrationServerKey,
		"2001:db8:ffff::10",
		map[string]any{"ticket": ticket.Ticket},
		&redeemed,
	)
	if status != http.StatusOK {
		t.Fatalf("valid game-server redemption: expected 200, got %d: %s", status, body)
	}
	if redeemed.AccountID != bootstrap.Snapshot.AccountID || redeemed.VehicleID != bootstrap.Snapshot.VehicleID {
		t.Fatalf("redeemed snapshot changed authoritative identity: %+v vs %+v", redeemed, bootstrap.Snapshot)
	}

	// One-time ticket replay remains rejected through the complete middleware stack.
	status, body = integratedJSONRequest(
		t,
		trusted.URL+"/v1/internal/game-tickets/redeem",
		http.MethodPost,
		"",
		integrationServerKey,
		"2001:db8:ffff::10",
		map[string]any{"ticket": ticket.Ticket},
		&core.Snapshot{},
	)
	if status != http.StatusUnauthorized {
		t.Fatalf("ticket replay: expected 401, got %d: %s", status, body)
	}

	t.Logf("integrated trust-boundary evidence: direct_xff_spoof_limited=true trusted_ingress_identity=true forged_session_rejected=true wrong_server_key_rejected=true ticket_replay_rejected=true")
}

func integratedJSONRequest(
	t *testing.T,
	url string,
	method string,
	bearer string,
	serverKey string,
	xForwardedFor string,
	body any,
	target any,
) (int, string) {
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
	if xForwardedFor != "" {
		req.Header.Set("X-Forwarded-For", xForwardedFor)
	}

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	payload, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if target != nil && len(payload) > 0 && resp.StatusCode >= 200 && resp.StatusCode < 300 {
		if err = json.Unmarshal(payload, target); err != nil {
			t.Fatalf("decode %s: %v (%s)", url, err, string(payload))
		}
	}
	return resp.StatusCode, string(payload)
}
