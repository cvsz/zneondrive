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

func TestParseTrustedProxyCIDRsRejectsInvalidConfig(t *testing.T) {
	if _, err := ParseTrustedProxyCIDRs("10.0.0.0/8,not-a-network"); err == nil {
		t.Fatal("invalid trusted proxy CIDR must fail closed")
	}
}

func TestParseTrustedProxyCIDRsAcceptsCIDRsAndExactIPs(t *testing.T) {
	policy, err := ParseTrustedProxyCIDRs("10.0.0.0/8, 2001:db8::/32, 192.0.2.44")
	if err != nil {
		t.Fatalf("parse trusted proxies: %v", err)
	}
	if got := len(policy.trusted); got != 3 {
		t.Fatalf("expected 3 trusted networks, got %d", got)
	}
}

func TestUntrustedPeerCannotSpoofForwardedClientIdentity(t *testing.T) {
	policy, err := ParseTrustedProxyCIDRs("10.0.0.0/8")
	if err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest(http.MethodPost, "/v1/sessions/bootstrap", nil)
	req.RemoteAddr = "203.0.113.77:41234"
	req.Header.Set("X-Forwarded-For", "198.51.100.9")

	if got := remoteIdentityWithTrustedProxies(req, policy); got != "203.0.113.77" {
		t.Fatalf("untrusted peer must not control forwarded identity, got %q", got)
	}
}

func TestTrustedProxyUsesForwardedClientIdentity(t *testing.T) {
	policy, err := ParseTrustedProxyCIDRs("10.0.0.0/8")
	if err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest(http.MethodPost, "/v1/sessions/bootstrap", nil)
	req.RemoteAddr = "10.10.0.8:443"
	req.Header.Set("X-Forwarded-For", "198.51.100.23")

	if got := remoteIdentityWithTrustedProxies(req, policy); got != "198.51.100.23" {
		t.Fatalf("trusted proxy should preserve real client identity, got %q", got)
	}
}

func TestTrustedProxyWalksForwardedChainRightToLeft(t *testing.T) {
	policy, err := ParseTrustedProxyCIDRs("10.0.0.0/8,192.168.0.0/16")
	if err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest(http.MethodPost, "/v1/sessions/bootstrap", nil)
	req.RemoteAddr = "10.0.0.7:443"
	// A client-supplied left-most spoof must not win when the trusted edge has
	// appended the real client followed by an internal proxy hop.
	req.Header.Set("X-Forwarded-For", "192.0.2.200, 198.51.100.42, 192.168.20.10")

	if got := remoteIdentityWithTrustedProxies(req, policy); got != "198.51.100.42" {
		t.Fatalf("expected first untrusted hop from the right, got %q", got)
	}
}

func TestMalformedForwardedChainFallsBackToSocketPeer(t *testing.T) {
	policy, err := ParseTrustedProxyCIDRs("10.0.0.0/8")
	if err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest(http.MethodPost, "/v1/sessions/bootstrap", nil)
	req.RemoteAddr = "10.0.0.7:443"
	req.Header.Set("X-Forwarded-For", "198.51.100.42, definitely-not-an-ip")

	if got := remoteIdentityWithTrustedProxies(req, policy); got != "10.0.0.7" {
		t.Fatalf("malformed chain must fail closed to socket peer, got %q", got)
	}
}

func TestTrustedProxySeparatesBootstrapBucketsByClient(t *testing.T) {
	policy, err := ParseTrustedProxyCIDRs("10.0.0.0/8")
	if err != nil {
		t.Fatal(err)
	}

	reqA := httptest.NewRequest(http.MethodPost, "/v1/sessions/bootstrap", nil)
	reqA.RemoteAddr = "10.0.0.9:443"
	reqA.Header.Set("X-Forwarded-For", "198.51.100.10")
	reqB := httptest.NewRequest(http.MethodPost, "/v1/sessions/bootstrap", nil)
	reqB.RemoteAddr = "10.0.0.9:443"
	reqB.Header.Set("X-Forwarded-For", "198.51.100.11")

	_, scopeA, identityA, limitedA := rateLimitRuleWithTrustedProxies(reqA, policy)
	_, scopeB, identityB, limitedB := rateLimitRuleWithTrustedProxies(reqB, policy)
	if !limitedA || !limitedB || scopeA != "bootstrap" || scopeB != "bootstrap" {
		t.Fatal("bootstrap route should be rate limited")
	}
	if rateLimitKey(scopeA, identityA) == rateLimitKey(scopeB, identityB) {
		t.Fatal("distinct clients behind a trusted proxy must not collapse into one bucket")
	}
}
