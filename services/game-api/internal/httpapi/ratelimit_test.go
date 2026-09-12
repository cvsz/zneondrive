package httpapi

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestRateLimiterBurstAndRefill(t *testing.T) {
	limiter := newRateLimiter(10)
	now := time.Date(2026, 9, 13, 0, 0, 0, 0, time.UTC)
	limiter.now = func() time.Time { return now }
	policy := rateLimitPolicy{Burst: 2, RefillEvery: 5 * time.Second}

	if ok, _ := limiter.allow("key", policy); !ok {
		t.Fatal("first token should be allowed")
	}
	if ok, _ := limiter.allow("key", policy); !ok {
		t.Fatal("second token should be allowed")
	}
	if ok, retry := limiter.allow("key", policy); ok || retry != 5*time.Second {
		t.Fatalf("third token should be limited for 5s, got ok=%v retry=%s", ok, retry)
	}

	now = now.Add(5 * time.Second)
	if ok, _ := limiter.allow("key", policy); !ok {
		t.Fatal("one token should refill after refill interval")
	}
}

func TestRateLimiterBoundsBucketCardinality(t *testing.T) {
	limiter := newRateLimiter(2)
	policy := rateLimitPolicy{Burst: 1, RefillEvery: time.Minute}
	for _, key := range []string{"a", "b", "c"} {
		if ok, _ := limiter.allow(key, policy); !ok {
			t.Fatalf("new key %q should be allowed", key)
		}
	}
	if got := len(limiter.buckets); got != 2 {
		t.Fatalf("expected bounded bucket cardinality 2, got %d", got)
	}
}

func TestRateLimitKeyDoesNotExposeCredential(t *testing.T) {
	secret := "super-secret-session-token"
	key := rateLimitKey("state", secret)
	if strings.Contains(key, secret) {
		t.Fatal("rate-limit key must not embed raw credentials")
	}
	if !strings.HasPrefix(key, "state:") {
		t.Fatalf("expected scoped key, got %q", key)
	}
}

func TestRateLimitedHandlerReturns429AndRetryAfter(t *testing.T) {
	nextCalls := 0
	next := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		nextCalls++
		w.WriteHeader(http.StatusNoContent)
	})
	middleware := &rateLimitMiddleware{next: next, limiter: newRateLimiter(10)}
	middleware.limiter.now = func() time.Time { return time.Date(2026, 9, 13, 0, 0, 0, 0, time.UTC) }

	for i := 0; i < bootstrapRatePolicy.Burst; i++ {
		req := httptest.NewRequest(http.MethodPost, "/v1/sessions/bootstrap", nil)
		req.RemoteAddr = "203.0.113.10:40000"
		rec := httptest.NewRecorder()
		middleware.ServeHTTP(rec, req)
		if rec.Code != http.StatusNoContent {
			t.Fatalf("request %d expected 204, got %d", i+1, rec.Code)
		}
	}

	req := httptest.NewRequest(http.MethodPost, "/v1/sessions/bootstrap", nil)
	req.RemoteAddr = "203.0.113.10:40001"
	rec := httptest.NewRecorder()
	middleware.ServeHTTP(rec, req)
	if rec.Code != http.StatusTooManyRequests {
		t.Fatalf("expected 429, got %d", rec.Code)
	}
	if rec.Header().Get("Retry-After") == "" {
		t.Fatal("rate-limited response must include Retry-After")
	}
	if !strings.Contains(rec.Body.String(), `"error":"rate_limited"`) {
		t.Fatalf("unexpected response body: %s", rec.Body.String())
	}
	if nextCalls != bootstrapRatePolicy.Burst {
		t.Fatalf("limited request reached next handler: calls=%d", nextCalls)
	}
}

func TestBearerIdentitySeparatesAuthenticatedBuckets(t *testing.T) {
	reqA := httptest.NewRequest(http.MethodGet, "/v1/state", nil)
	reqA.Header.Set("Authorization", "Bearer token-a")
	reqB := httptest.NewRequest(http.MethodGet, "/v1/state", nil)
	reqB.Header.Set("Authorization", "Bearer token-b")

	_, scopeA, identityA, limitedA := rateLimitRule(reqA)
	_, scopeB, identityB, limitedB := rateLimitRule(reqB)
	if !limitedA || !limitedB || scopeA != "state" || scopeB != "state" {
		t.Fatal("state route should be rate limited")
	}
	if rateLimitKey(scopeA, identityA) == rateLimitKey(scopeB, identityB) {
		t.Fatal("different bearer credentials must not share a rate-limit bucket")
	}
}

func TestHealthEndpointIsNotRateLimited(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	_, _, _, limited := rateLimitRule(req)
	if limited {
		t.Fatal("health endpoint must remain available to probes")
	}
}
