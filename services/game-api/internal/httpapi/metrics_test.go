package httpapi

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

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
