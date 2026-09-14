//go:build integration

package httpapi

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"
)

func TestRedisReplicaPromotionPreservesDistributedLimiterBudget(t *testing.T) {
	if _, err := exec.LookPath("docker"); err != nil {
		if os.Getenv("CI") == "true" {
			t.Fatalf("docker is required for the Redis promotion evidence gate: %v", err)
		}
		t.Skipf("docker unavailable: %v", err)
	}

	primaryPort := reserveTCPPort(t)
	replicaPort := reserveTCPPort(t)
	suffix := fmt.Sprintf("%d", time.Now().UnixNano())
	network := "zneondrive-redis-promotion-" + suffix
	primary := "zneondrive-redis-primary-" + suffix
	replica := "zneondrive-redis-replica-" + suffix

	runDocker(t, "network", "create", network)
	t.Cleanup(func() {
		_ = exec.Command("docker", "rm", "-f", primary, replica).Run()
		_ = exec.Command("docker", "network", "rm", network).Run()
	})

	runDocker(t,
		"run", "-d", "--name", primary, "--network", network,
		"-p", fmt.Sprintf("127.0.0.1:%d:6379", primaryPort),
		"redis:8-alpine", "redis-server", "--save", "", "--appendonly", "no",
	)
	waitRedisContainerReady(t, primary)

	runDocker(t,
		"run", "-d", "--name", replica, "--network", network,
		"-p", fmt.Sprintf("127.0.0.1:%d:6379", replicaPort),
		"redis:8-alpine", "redis-server", "--save", "", "--appendonly", "no",
		"--replicaof", primary, "6379",
	)
	waitRedisContainerReady(t, replica)
	waitRedisReplicationLink(t, replica)

	primaryLimiter := newRedisRateLimiter(fmt.Sprintf("127.0.0.1:%d", primaryPort))
	policy := rateLimitPolicy{Burst: 2, RefillEvery: 30 * time.Second}
	now := time.Now().UTC().Truncate(time.Millisecond)
	key := "promotion:" + suffix
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	for i := 0; i < policy.Burst; i++ {
		ok, _, err := primaryLimiter.allow(ctx, key, policy, now)
		if err != nil || !ok {
			t.Fatalf("consume primary limiter token %d: ok=%v err=%v", i+1, ok, err)
		}
	}
	if ok, _, err := primaryLimiter.allow(ctx, key, policy, now); err != nil || ok {
		t.Fatalf("primary must observe exhausted budget before promotion: ok=%v err=%v", ok, err)
	}

	redisKey := "zneondrive:ratelimit:" + key
	waitRedisHashField(t, replica, redisKey, "tokens", "0")

	// Stop the old primary before promotion. This test proves state continuity
	// after explicit promotion, not automatic failover or split-brain fencing.
	runDocker(t, "stop", "-t", "1", primary)
	runDocker(t, "exec", replica, "redis-cli", "REPLICAOF", "NO", "ONE")
	waitRedisRole(t, replica, "role:master")

	promotedLimiter := newRedisRateLimiter(fmt.Sprintf("127.0.0.1:%d", replicaPort))
	if ok, retry, err := promotedLimiter.allow(ctx, key, policy, now); err != nil || ok || retry < time.Second {
		t.Fatalf("promoted Redis must preserve exhausted shared budget: ok=%v retry=%s err=%v", ok, retry, err)
	}

	if ok, _, err := promotedLimiter.allow(ctx, key, policy, now.Add(policy.RefillEvery)); err != nil || !ok {
		t.Fatalf("promoted Redis must continue canonical refill semantics: ok=%v err=%v", ok, err)
	}
}

func reserveTCPPort(t *testing.T) int {
	t.Helper()
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("reserve TCP port: %v", err)
	}
	defer listener.Close()
	return listener.Addr().(*net.TCPAddr).Port
}

func runDocker(t *testing.T, args ...string) string {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	out, err := exec.CommandContext(ctx, "docker", args...).CombinedOutput()
	if err != nil {
		t.Fatalf("docker %s failed: %v\n%s", strings.Join(args, " "), err, strings.TrimSpace(string(out)))
	}
	return strings.TrimSpace(string(out))
}

func waitRedisContainerReady(t *testing.T, container string) {
	t.Helper()
	waitForRedisCondition(t, 30*time.Second, func() bool {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		out, err := exec.CommandContext(ctx, "docker", "exec", container, "redis-cli", "PING").CombinedOutput()
		return err == nil && strings.TrimSpace(string(out)) == "PONG"
	}, container+" never became PING-ready")
}

func waitRedisReplicationLink(t *testing.T, replica string) {
	t.Helper()
	waitForRedisCondition(t, 30*time.Second, func() bool {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		out, err := exec.CommandContext(ctx, "docker", "exec", replica, "redis-cli", "INFO", "replication").CombinedOutput()
		text := string(out)
		return err == nil && strings.Contains(text, "role:slave") && strings.Contains(text, "master_link_status:up")
	}, replica+" never reached connected replica state")
}

func waitRedisHashField(t *testing.T, container, key, field, expected string) {
	t.Helper()
	waitForRedisCondition(t, 30*time.Second, func() bool {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		out, err := exec.CommandContext(ctx, "docker", "exec", container, "redis-cli", "--raw", "HGET", key, field).CombinedOutput()
		return err == nil && strings.TrimSpace(string(out)) == expected
	}, fmt.Sprintf("%s never replicated %s[%s]=%s", container, key, field, expected))
}

func waitRedisRole(t *testing.T, container, expected string) {
	t.Helper()
	waitForRedisCondition(t, 15*time.Second, func() bool {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		out, err := exec.CommandContext(ctx, "docker", "exec", container, "redis-cli", "INFO", "replication").CombinedOutput()
		return err == nil && strings.Contains(string(out), expected)
	}, container+" never reached "+expected)
}

func waitForRedisCondition(t *testing.T, timeout time.Duration, condition func() bool, failure string) {
	t.Helper()
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if condition() {
			return
		}
		time.Sleep(250 * time.Millisecond)
	}
	t.Fatal(failure)
}
