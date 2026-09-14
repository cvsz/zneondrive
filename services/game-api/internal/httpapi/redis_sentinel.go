package httpapi

import (
	"net/http"
	"strings"
)

// NewSentinelDistributedRateLimitedHandler resolves the active Redis primary
// through Sentinel before each distributed limiter decision. Redis remains
// coordination-only; PostgreSQL and the dedicated server retain durable/gameplay
// authority. If Sentinel/Redis is unavailable, the middleware's existing local
// bounded limiter remains the fail-safe fallback.
func NewSentinelDistributedRateLimitedHandler(next http.Handler, sentinelAddrs []string, masterName string) http.Handler {
	return NewSentinelDistributedRateLimitedHandlerWithTrustedProxies(next, sentinelAddrs, masterName, TrustedProxyPolicy{})
}

func NewSentinelDistributedRateLimitedHandlerWithTrustedProxies(next http.Handler, sentinelAddrs []string, masterName string, policy TrustedProxyPolicy) http.Handler {
	middleware := &rateLimitMiddleware{next: next, limiter: newRateLimiter(10000), proxyPolicy: policy}
	if strings.TrimSpace(masterName) != "" && len(sentinelAddrs) > 0 {
		middleware.redis = newRedisSentinelRateLimiter(sentinelAddrs, masterName)
	}
	return middleware
}
