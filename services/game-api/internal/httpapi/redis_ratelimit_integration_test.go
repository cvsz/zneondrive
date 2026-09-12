//go:build integration

package httpapi

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"testing"
	"time"
)

func TestRedisRateLimiterSharesBudgetAcrossInstances(t *testing.T) {
	addr := os.Getenv("TEST_REDIS_ADDR")
	if addr == "" {
		t.Skip("TEST_REDIS_ADDR not configured")
	}

	first := newRedisRateLimiter(addr)
	second := newRedisRateLimiter(addr)
	policy := rateLimitPolicy{Burst: 2, RefillEvery: 2 * time.Second}
	now := time.Now().UTC().Truncate(time.Millisecond)
	key := fmt.Sprintf("integration:%d", now.UnixNano())
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if ok, _, err := first.allow(ctx, key, policy, now); err != nil || !ok {
		t.Fatalf("first replica first token: ok=%v err=%v", ok, err)
	}
	if ok, _, err := first.allow(ctx, key, policy, now); err != nil || !ok {
		t.Fatalf("first replica second token: ok=%v err=%v", ok, err)
	}
	if ok, retry, err := second.allow(ctx, key, policy, now); err != nil || ok || retry < time.Second {
		t.Fatalf("second replica must observe exhausted shared budget: ok=%v retry=%s err=%v", ok, retry, err)
	}

	if ok, _, err := second.allow(ctx, key, policy, now.Add(2*time.Second)); err != nil || !ok {
		t.Fatalf("shared bucket should refill after interval: ok=%v err=%v", ok, err)
	}
}

type httpLoadResult struct {
	replica    int
	status     int
	retryAfter string
	err        error
}

func TestDistributedRateLimiterConcurrentHTTPAcrossReplicas(t *testing.T) {
	addr := os.Getenv("TEST_REDIS_ADDR")
	if addr == "" {
		t.Skip("TEST_REDIS_ADDR not configured")
	}

	proxyPolicy, err := ParseTrustedProxyCIDRs("127.0.0.1/32")
	if err != nil {
		t.Fatalf("parse trusted proxy policy: %v", err)
	}

	originalPolicy := bootstrapRatePolicy
	bootstrapRatePolicy = rateLimitPolicy{Burst: 32, RefillEvery: time.Hour}
	defer func() { bootstrapRatePolicy = originalPolicy }()

	next := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
	first := httptest.NewServer(NewDistributedRateLimitedHandlerWithTrustedProxies(next, addr, proxyPolicy))
	defer first.Close()
	second := httptest.NewServer(NewDistributedRateLimitedHandlerWithTrustedProxies(next, addr, proxyPolicy))
	defer second.Close()

	const totalRequests = 128
	clientIP := fmt.Sprintf("2001:db8::%x", uint64(time.Now().UnixNano())&0xffff)
	client := &http.Client{Timeout: 5 * time.Second}
	start := make(chan struct{})
	results := make(chan httpLoadResult, totalRequests)
	var wg sync.WaitGroup

	startedAt := time.Now()
	for i := 0; i < totalRequests; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			<-start

			replica := index % 2
			target := first.URL
			if replica == 1 {
				target = second.URL
			}

			req, reqErr := http.NewRequest(http.MethodPost, target+"/v1/sessions/bootstrap", nil)
			if reqErr != nil {
				results <- httpLoadResult{replica: replica, err: reqErr}
				return
			}
			req.Header.Set("X-Forwarded-For", clientIP)
			resp, doErr := client.Do(req)
			if doErr != nil {
				results <- httpLoadResult{replica: replica, err: doErr}
				return
			}
			_ = resp.Body.Close()
			results <- httpLoadResult{replica: replica, status: resp.StatusCode, retryAfter: resp.Header.Get("Retry-After")}
		}(i)
	}
	close(start)
	wg.Wait()
	close(results)

	allowed := 0
	limited := 0
	replicaRequests := [2]int{}
	for result := range results {
		if result.err != nil {
			t.Fatalf("replica %d HTTP request failed: %v", result.replica, result.err)
		}
		replicaRequests[result.replica]++
		switch result.status {
		case http.StatusNoContent:
			allowed++
		case http.StatusTooManyRequests:
			if result.retryAfter == "" {
				t.Fatalf("rate-limited response from replica %d omitted Retry-After", result.replica)
			}
			limited++
		default:
			t.Fatalf("unexpected HTTP status from replica %d: %d", result.replica, result.status)
		}
	}

	if replicaRequests[0] != totalRequests/2 || replicaRequests[1] != totalRequests/2 {
		t.Fatalf("expected balanced traffic across replicas, got replica0=%d replica1=%d", replicaRequests[0], replicaRequests[1])
	}
	if allowed != bootstrapRatePolicy.Burst {
		t.Fatalf("shared Redis budget must allow exactly %d requests across both replicas, got %d", bootstrapRatePolicy.Burst, allowed)
	}
	if limited != totalRequests-bootstrapRatePolicy.Burst {
		t.Fatalf("expected %d shared-budget rejections, got %d", totalRequests-bootstrapRatePolicy.Burst, limited)
	}

	t.Logf("concurrent multi-replica HTTP evidence: requests=%d allowed=%d limited=%d duration=%s", totalRequests, allowed, limited, time.Since(startedAt))
}
