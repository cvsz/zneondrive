package httpapi

import (
	"crypto/sha256"
	"encoding/hex"
	"math"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

type rateLimitPolicy struct {
	Burst        int
	RefillEvery  time.Duration
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
	retryAfter := time.Duration(math.Ceil(missing*float64(policy.RefillEvery)))
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

func (a *API) enforceRateLimit(w http.ResponseWriter, r *http.Request, scope, identity string, policy rateLimitPolicy) bool {
	allowed, retryAfter := a.limiter.allow(rateLimitKey(scope, identity), policy)
	if allowed {
		return true
	}
	seconds := int(math.Ceil(retryAfter.Seconds()))
	if seconds < 1 {
		seconds = 1
	}
	w.Header().Set("Retry-After", strconv.Itoa(seconds))
	writeError(w, http.StatusTooManyRequests, "rate_limited")
	return false
}

func rateLimitKey(scope, identity string) string {
	digest := sha256.Sum256([]byte(strings.TrimSpace(identity)))
	return scope + ":" + hex.EncodeToString(digest[:16])
}

func remoteIdentity(r *http.Request) string {
	remote := strings.TrimSpace(r.RemoteAddr)
	if host, _, err := net.SplitHostPort(remote); err == nil && host != "" {
		return host
	}
	if remote == "" {
		return "unknown"
	}
	return remote
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
