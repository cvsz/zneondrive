package httpapi

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/cvsz/zneondrive/services/game-api/internal/store"
)

type fakePostgresStatsProvider struct {
	stats store.PostgresPoolStats
}

func (f fakePostgresStatsProvider) PostgresPoolStats() store.PostgresPoolStats {
	return f.stats
}

func TestHTTPMetricsRecordsBoundedRouteAndStatus(t *testing.T) {
	metrics := NewHTTPMetrics()
	next := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusConflict)
	})
	handler := metrics.Wrap(next)

	req := httptest.NewRequest(http.MethodPost, "/v1/internal/races/race-secret-123/checkpoints", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	metricsRec := httptest.NewRecorder()
	metrics.Handler().ServeHTTP(metricsRec, httptest.NewRequest(http.MethodGet, "/metrics", nil))
	body := metricsRec.Body.String()

	if !strings.Contains(body, `zneondrive_http_requests_total{route="race-checkpoint",status_class="4xx"} 1`) {
		t.Fatalf("missing bounded request metric: %s", body)
	}
	if !strings.Contains(body, `zneondrive_http_request_duration_seconds_count{route="race-checkpoint"} 1`) {
		t.Fatalf("missing duration count: %s", body)
	}
	if strings.Contains(body, "race-secret-123") {
		t.Fatal("metrics must not expose dynamic race identifiers")
	}
}

func TestHTTPMetricsDoNotExposeAuthorizationOrDynamicIDs(t *testing.T) {
	metrics := NewHTTPMetrics()
	handler := metrics.Wrap(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))

	req := httptest.NewRequest(http.MethodPost, "/v1/vehicles/vehicle-private-456/builds", nil)
	req.Header.Set("Authorization", "Bearer super-secret-session-token")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	metricsRec := httptest.NewRecorder()
	metrics.Handler().ServeHTTP(metricsRec, httptest.NewRequest(http.MethodGet, "/metrics", nil))
	body := metricsRec.Body.String()

	for _, forbidden := range []string{"super-secret-session-token", "vehicle-private-456"} {
		if strings.Contains(body, forbidden) {
			t.Fatalf("metrics exposed sensitive/dynamic value %q", forbidden)
		}
	}
	if !strings.Contains(body, `route="build-revise"`) {
		t.Fatalf("expected normalized route label, got: %s", body)
	}
}

func TestHTTPMetricsTracksDefault200AndInflightReturnsToZero(t *testing.T) {
	metrics := NewHTTPMetrics()
	handler := metrics.Wrap(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("ok"))
	}))

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/healthz", nil))

	metricsRec := httptest.NewRecorder()
	metrics.Handler().ServeHTTP(metricsRec, httptest.NewRequest(http.MethodGet, "/metrics", nil))
	body := metricsRec.Body.String()
	if !strings.Contains(body, `zneondrive_http_requests_total{route="health",status_class="2xx"} 1`) {
		t.Fatalf("expected default 200 request metric, got: %s", body)
	}
	if !strings.Contains(body, "zneondrive_http_in_flight_requests 0") {
		t.Fatalf("in-flight gauge must return to zero, got: %s", body)
	}
}

func TestHTTPMetricsExposeBoundedPostgresPoolStats(t *testing.T) {
	provider := fakePostgresStatsProvider{stats: store.PostgresPoolStats{
		MaxConns:             20,
		TotalConns:           7,
		IdleConns:            4,
		AcquiredConns:        2,
		ConstructingConns:    1,
		AcquireCount:         123,
		EmptyAcquireCount:    9,
		CanceledAcquireCount: 3,
		NewConnsCount:        11,
		AcquireDuration:      1500 * time.Millisecond,
	}}
	metrics := NewHTTPMetrics(provider)
	metricsRec := httptest.NewRecorder()
	metrics.Handler().ServeHTTP(metricsRec, httptest.NewRequest(http.MethodGet, "/metrics", nil))
	body := metricsRec.Body.String()

	for _, want := range []string{
		`zneondrive_postgres_pool_connections{state="max"} 20`,
		`zneondrive_postgres_pool_connections{state="total"} 7`,
		`zneondrive_postgres_pool_connections{state="idle"} 4`,
		`zneondrive_postgres_pool_connections{state="acquired"} 2`,
		`zneondrive_postgres_pool_connections{state="constructing"} 1`,
		`zneondrive_postgres_pool_acquires_total 123`,
		`zneondrive_postgres_pool_empty_acquires_total 9`,
		`zneondrive_postgres_pool_canceled_acquires_total 3`,
		`zneondrive_postgres_pool_new_connections_total 11`,
		`zneondrive_postgres_pool_acquire_duration_seconds 1.500000000`,
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("missing PostgreSQL pool metric %q: %s", want, body)
		}
	}
	for _, forbidden := range []string{"DATABASE_URL", "postgres://", "session", "vehicle", "race"} {
		if strings.Contains(body, forbidden) {
			t.Fatalf("PostgreSQL pool metrics exposed forbidden dynamic/sensitive token %q: %s", forbidden, body)
		}
	}
}

func TestMetricsRouteCardinalityIsStatic(t *testing.T) {
	cases := []struct {
		method string
		path   string
		want   string
	}{
		{http.MethodPost, "/v1/quests/MQ012/complete", "quest-complete"},
		{http.MethodPost, "/v1/vehicles/V-123/builds", "build-revise"},
		{http.MethodPost, "/v1/internal/races/R-123/checkpoints", "race-checkpoint"},
		{http.MethodPost, "/v1/internal/races/R-123/finish", "race-finish"},
		{http.MethodGet, "/anything/else", "other"},
	}
	for _, tc := range cases {
		req := httptest.NewRequest(tc.method, tc.path, nil)
		if got := metricsRoute(req); got != tc.want {
			t.Fatalf("%s %s: got %q want %q", tc.method, tc.path, got, tc.want)
		}
	}
}
