package httpapi

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"math"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

type rateLimitPolicy struct {
	Burst       int
	RefillEvery time.Duration
}

type rateBucket struct {
	tokens  float64
	last    time.Time
	touched time.Time
}

type rateLimiter struct {
	mu         sync.Mutex
	buckets    map[string]*rateBucket
	maxEntries int
	now        func() time.Time
}

// TrustedProxyPolicy defines the only peers allowed to influence client
// identity through X-Forwarded-For. A zero-value policy trusts no proxy.
// This deliberately keeps trust configuration explicit and fail-closed.
type TrustedProxyPolicy struct {
	trusted []*net.IPNet
}

type rateLimitMiddleware struct {
	next        http.Handler
	limiter     *rateLimiter
	redis       *redisRateLimiter
	proxyPolicy TrustedProxyPolicy
}

func NewRateLimitedHandler(next http.Handler) http.Handler {
	return NewRateLimitedHandlerWithTrustedProxies(next, TrustedProxyPolicy{})
}

func NewRateLimitedHandlerWithTrustedProxies(next http.Handler, policy TrustedProxyPolicy) http.Handler {
	return &rateLimitMiddleware{next: next, limiter: newRateLimiter(10000), proxyPolicy: policy}
}

// NewDistributedRateLimitedHandler coordinates rate-limit decisions through
// Redis so multiple API replicas share the same abuse budget. If Redis is
// temporarily unavailable, the bounded in-process limiter remains active and
// a credential-safe security event is logged. Redis failure never moves
// gameplay authority away from PostgreSQL/dedicated-server validation.
func NewDistributedRateLimitedHandler(next http.Handler, redisAddr string) http.Handler {
	return NewDistributedRateLimitedHandlerWithTrustedProxies(next, redisAddr, TrustedProxyPolicy{})
}

func NewDistributedRateLimitedHandlerWithTrustedProxies(next http.Handler, redisAddr string, policy TrustedProxyPolicy) http.Handler {
	middleware := &rateLimitMiddleware{next: next, limiter: newRateLimiter(10000), proxyPolicy: policy}
	if strings.TrimSpace(redisAddr) != "" {
		middleware.redis = newRedisRateLimiter(redisAddr)
	}
	return middleware
}

// ParseTrustedProxyCIDRs parses a comma-separated allowlist such as
// "10.0.0.0/8,192.168.0.0/16". Invalid entries fail startup rather than
// silently widening trust. Plain IP literals are accepted as exact /32 or /128
// networks for small deployments.
func ParseTrustedProxyCIDRs(raw string) (TrustedProxyPolicy, error) {
	var policy TrustedProxyPolicy
	for _, entry := range strings.Split(raw, ",") {
		entry = strings.TrimSpace(entry)
		if entry == "" {
			continue
		}

		if ip := net.ParseIP(entry); ip != nil {
			bits := 128
			if ip.To4() != nil {
				ip = ip.To4()
				bits = 32
			}
			policy.trusted = append(policy.trusted, &net.IPNet{IP: ip, Mask: net.CIDRMask(bits, bits)})
			continue
		}

		_, network, err := net.ParseCIDR(entry)
		if err != nil {
			return TrustedProxyPolicy{}, fmt.Errorf("invalid trusted proxy CIDR %q: %w", entry, err)
		}
		policy.trusted = append(policy.trusted, network)
	}
	return policy, nil
}

func (p TrustedProxyPolicy) configured() bool {
	return len(p.trusted) > 0
}

func (p TrustedProxyPolicy) isTrusted(ip net.IP) bool {
	if ip == nil {
		return false
	}
	for _, network := range p.trusted {
		if network.Contains(ip) {
			return true
		}
	}
	return false
}

func newRateLimiter(maxEntries int) *rateLimiter {
	if maxEntries < 1 {
		maxEntries = 1
	}
	return &rateLimiter{
		buckets:    make(map[string]*rateBucket),
		maxEntries: maxEntries,
		now:        time.Now,
	}
}

func (m *rateLimitMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	policy, scope, identity, limited := rateLimitRuleWithTrustedProxies(r, m.proxyPolicy)
	if limited {
		key := rateLimitKey(scope, identity)
		backend := "local"
		allowed, retryAfter := m.limiter.allow(key, policy)
		if m.redis != nil {
			redisAllowed, redisRetry, err := m.redis.allow(r.Context(), key, policy, m.limiter.now())
			if err == nil {
				allowed, retryAfter, backend = redisAllowed, redisRetry, "redis"
			} else {
				log.Printf("security_event=rate_limit_backend_fallback scope=%s bucket=%s backend=local", scope, key)
			}
		}
		if !allowed {
			seconds := int(math.Ceil(retryAfter.Seconds()))
			if seconds < 1 {
				seconds = 1
			}
			log.Printf("security_event=rate_limit_rejected scope=%s bucket=%s backend=%s retry_after_seconds=%d", scope, key, backend, seconds)
			w.Header().Set("Retry-After", strconv.Itoa(seconds))
			writeError(w, http.StatusTooManyRequests, "rate_limited")
			return
		}
	}
	m.next.ServeHTTP(w, r)
}

func rateLimitRule(r *http.Request) (rateLimitPolicy, string, string, bool) {
	return rateLimitRuleWithTrustedProxies(r, TrustedProxyPolicy{})
}

func rateLimitRuleWithTrustedProxies(r *http.Request, proxyPolicy TrustedProxyPolicy) (rateLimitPolicy, string, string, bool) {
	path := r.URL.Path
	remote := remoteIdentityWithTrustedProxies(r, proxyPolicy)

	switch {
	case r.Method == http.MethodPost && path == "/v1/sessions/bootstrap":
		return bootstrapRatePolicy, "bootstrap", remote, true
	case r.Method == http.MethodGet && path == "/v1/state":
		return stateRatePolicy, "state", bearerIdentity(r, remote), true
	case r.Method == http.MethodPost && path == "/v1/game-tickets":
		return gameTicketRatePolicy, "game-ticket", bearerIdentity(r, remote), true
	case r.Method == http.MethodPost && strings.HasPrefix(path, "/v1/quests/") && strings.HasSuffix(path, "/complete"):
		return mutationRatePolicy, "quest", bearerIdentity(r, remote), true
	case r.Method == http.MethodPost && strings.HasPrefix(path, "/v1/vehicles/") && strings.HasSuffix(path, "/builds"):
		return mutationRatePolicy, "build", bearerIdentity(r, remote), true
	case r.Method == http.MethodPost && path == "/v1/internal/game-tickets/redeem":
		return internalAuthRatePolicy, "internal-ticket-redeem", remote, true
	case r.Method == http.MethodPost && path == "/v1/internal/races/start":
		return raceStartRatePolicy, "race-start", remote, true
	case r.Method == http.MethodPost && strings.HasPrefix(path, "/v1/internal/races/") && strings.HasSuffix(path, "/checkpoints"):
		return raceCheckpointRatePolicy, "race-checkpoint", remote, true
	case r.Method == http.MethodPost && strings.HasPrefix(path, "/v1/internal/races/") && strings.HasSuffix(path, "/finish"):
		return raceFinishRatePolicy, "race-finish", remote, true
	default:
		return rateLimitPolicy{}, "", "", false
	}
}

func bearerIdentity(r *http.Request, fallback string) string {
	if token, ok := bearerToken(r); ok {
		return token
	}
	return fallback
}

func (l *rateLimiter) allow(key string, policy rateLimitPolicy) (bool, time.Duration) {
	if l == nil || policy.Burst <= 0 || policy.RefillEvery <= 0 {
		return true, 0
	}

	now := l.now().UTC()
	l.mu.Lock()
	defer l.mu.Unlock()

	bucket, ok := l.buckets[key]
	if !ok {
		l.makeRoom(now)
		bucket = &rateBucket{tokens: float64(policy.Burst), last: now, touched: now}
		l.buckets[key] = bucket
	}

	elapsed := now.Sub(bucket.last)
	if elapsed > 0 {
		bucket.tokens = math.Min(float64(policy.Burst), bucket.tokens+float64(elapsed)/float64(policy.RefillEvery))
		bucket.last = now
	}
	bucket.touched = now

	if bucket.tokens >= 1 {
		bucket.tokens--
		return true, 0
	}

	missing := 1 - bucket.tokens
	retryAfter := time.Duration(math.Ceil(missing * float64(policy.RefillEvery)))
	if retryAfter < time.Second {
		retryAfter = time.Second
	}
	return false, retryAfter
}

func (l *rateLimiter) makeRoom(now time.Time) {
	if len(l.buckets) < l.maxEntries {
		return
	}

	staleBefore := now.Add(-15 * time.Minute)
	for key, bucket := range l.buckets {
		if bucket.touched.Before(staleBefore) {
			delete(l.buckets, key)
		}
	}
	if len(l.buckets) < l.maxEntries {
		return
	}

	var oldestKey string
	var oldest time.Time
	for key, bucket := range l.buckets {
		if oldestKey == "" || bucket.touched.Before(oldest) {
			oldestKey = key
			oldest = bucket.touched
		}
	}
	if oldestKey != "" {
		delete(l.buckets, oldestKey)
	}
}

func rateLimitKey(scope, identity string) string {
	digest := sha256.Sum256([]byte(strings.TrimSpace(identity)))
	return scope + ":" + hex.EncodeToString(digest[:16])
}

func remoteIdentity(r *http.Request) string {
	return remoteIdentityWithTrustedProxies(r, TrustedProxyPolicy{})
}

func remoteIdentityWithTrustedProxies(r *http.Request, policy TrustedProxyPolicy) string {
	peer := remoteIP(r.RemoteAddr)
	if peer == nil {
		return "unknown"
	}
	peerText := peer.String()

	// Forwarded identity is ignored unless the socket peer itself is trusted.
	// This prevents direct clients from spoofing X-Forwarded-For.
	if !policy.isTrusted(peer) {
		return peerText
	}

	xff := strings.TrimSpace(r.Header.Get("X-Forwarded-For"))
	if xff == "" {
		return peerText
	}

	parts := strings.Split(xff, ",")
	chain := make([]net.IP, 0, len(parts)+1)
	for _, part := range parts {
		ip := net.ParseIP(strings.TrimSpace(part))
		if ip == nil {
			log.Printf("security_event=forwarded_identity_rejected reason=malformed_x_forwarded_for")
			return peerText
		}
		chain = append(chain, ip)
	}
	chain = append(chain, peer)

	// Walk from the immediate peer toward the client. Trusted proxy hops are
	// skipped; the first untrusted address is the effective client identity.
	// This is robust when a trusted edge appends to a pre-existing XFF chain.
	for i := len(chain) - 1; i >= 0; i-- {
		if policy.isTrusted(chain[i]) {
			continue
		}
		return chain[i].String()
	}

	// An all-trusted chain can legitimately occur for internal callers. Use the
	// left-most forwarded address so distinct trusted clients remain distinct.
	return chain[0].String()
}

func remoteIP(remoteAddr string) net.IP {
	remote := strings.TrimSpace(remoteAddr)
	if host, _, err := net.SplitHostPort(remote); err == nil {
		if ip := net.ParseIP(strings.TrimSpace(host)); ip != nil {
			return ip
		}
	}
	return net.ParseIP(remote)
}

var (
	bootstrapRatePolicy      = rateLimitPolicy{Burst: 8, RefillEvery: 8 * time.Second}
	stateRatePolicy          = rateLimitPolicy{Burst: 30, RefillEvery: 2 * time.Second}
	gameTicketRatePolicy     = rateLimitPolicy{Burst: 12, RefillEvery: 5 * time.Second}
	mutationRatePolicy       = rateLimitPolicy{Burst: 30, RefillEvery: 2 * time.Second}
	internalAuthRatePolicy   = rateLimitPolicy{Burst: 30, RefillEvery: 2 * time.Second}
	raceStartRatePolicy      = rateLimitPolicy{Burst: 20, RefillEvery: 3 * time.Second}
	raceCheckpointRatePolicy = rateLimitPolicy{Burst: 120, RefillEvery: 500 * time.Millisecond}
	raceFinishRatePolicy     = rateLimitPolicy{Burst: 30, RefillEvery: 2 * time.Second}
)
