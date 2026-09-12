//go:build integration

package httpapi

import (
	"context"
	"fmt"
	"os"
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
