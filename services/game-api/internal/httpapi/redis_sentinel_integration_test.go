//go:build integration

package httpapi

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"
)

func TestRedisSentinelFailoverPreservesDistributedLimiterBudget(t *testing.T) {
	if _, err := exec.LookPath("docker"); err != nil {
		if os.Getenv("CI") == "true" {
			t.Fatalf("docker is required for the Redis Sentinel evidence gate: %v", err)
		}
		t.Skipf("docker unavailable: %v", err)
	}

	primaryPort := reserveTCPPort(t)
	replicaPort := reserveTCPPort(t)
	sentinelPorts := []int{reserveTCPPort(t), reserveTCPPort(t), reserveTCPPort(t)}
	suffix := fmt.Sprintf("%d", time.Now().UnixNano())
	network := "zneondrive-redis-sentinel-" + suffix
	primary := "zneondrive-sentinel-primary-" + suffix
	replica := "zneondrive-sentinel-replica-" + suffix
	masterName := "zneondrive-primary"
	sentinels := []string{
		"zneondrive-sentinel-a-" + suffix,
		"zneondrive-sentinel-b-" + suffix,
		"zneondrive-sentinel-c-" + suffix,
	}

	runDocker(t, "network", "create", network)
	t.Cleanup(func() {
		args := append([]string{"rm", "-f", primary, replica}, sentinels...)
		_ = exec.Command("docker", args...).Run()
		_ = exec.Command("docker", "network", "rm", network).Run()
	})

	runDocker(t,
		"run", "-d", "--name", primary, "--network", network, "--network-alias", masterName,
		"-p", fmt.Sprintf("127.0.0.1:%d:6379", primaryPort),
		"redis:8-alpine", "redis-server", "--save", "", "--appendonly", "no",
	)
	waitRedisContainerReady(t, primary)

	runDocker(t,
		"run", "-d", "--name", replica, "--network", network,
		"-p", fmt.Sprintf("127.0.0.1:%d:6379", replicaPort),
		"redis:8-alpine", "redis-server", "--save", "", "--appendonly", "no",
		"--replicaof", masterName, "6379",
	)
	waitRedisContainerReady(t, replica)
	waitRedisReplicationLink(t, replica)

	for i, sentinel := range sentinels {
		config := fmt.Sprintf("port 26379\nsentinel monitor %s %s 6379 2\nsentinel down-after-milliseconds %s 1000\nsentinel failover-timeout %s 10000\nsentinel parallel-syncs %s 1\n", masterName, masterName, masterName, masterName, masterName)
		runDocker(t,
			"run", "-d", "--name", sentinel, "--network", network,
			"-p", fmt.Sprintf("127.0.0.1:%d:26379", sentinelPorts[i]),
			"redis:8-alpine", "sh", "-c",
			fmt.Sprintf("printf '%%s' %q > /tmp/sentinel.conf && exec redis-server /tmp/sentinel.conf --sentinel", config),
		)
		waitRedisContainerReady(t, sentinel)
	}

	sentinelAddrs := make([]string, 0, len(sentinelPorts))
	for _, port := range sentinelPorts {
		sentinelAddrs = append(sentinelAddrs, fmt.Sprintf("127.0.0.1:%d", port))
	}
	limiter := newRedisSentinelRateLimiter(sentinelAddrs, masterName)
	policy := rateLimitPolicy{Burst: 2, RefillEvery: 30 * time.Second}
	now := time.Now().UTC().Truncate(time.Millisecond)
	key := "sentinel:" + suffix
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	waitForRedisCondition(t, 30*time.Second, func() bool {
		addr, err := limiter.resolveAddr(ctx)
		return err == nil && strings.HasSuffix(addr, ":6379")
	}, "Sentinel never returned the initial Redis master")

	for i := 0; i < policy.Burst; i++ {
		ok, _, err := limiter.allow(ctx, key, policy, now)
		if err != nil || !ok {
			t.Fatalf("consume Sentinel-backed limiter token %d: ok=%v err=%v", i+1, ok, err)
		}
	}
	if ok, _, err := limiter.allow(ctx, key, policy, now); err != nil || ok {
		t.Fatalf("Sentinel-backed primary must observe exhausted budget: ok=%v err=%v", ok, err)
	}

	redisKey := "zneondrive:ratelimit:" + key
	waitRedisHashField(t, replica, redisKey, "tokens", "0")

	// Remove the old authority without changing the application's Sentinel
	// endpoints. A quorum of three Sentinels must discover and promote the replica.
	runDocker(t, "stop", "-t", "1", primary)
	waitRedisRole(t, replica, "role:master")
	waitForRedisCondition(t, 30*time.Second, func() bool {
		addr, err := limiter.resolveAddr(ctx)
		if err != nil {
			return false
		}
		containerIP := strings.TrimSpace(runDocker(t, "inspect", "-f", "{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}", replica))
		return addr == containerIP+":6379"
	}, "Sentinel never converged on the promoted Redis replica")

	if ok, retry, err := limiter.allow(ctx, key, policy, now); err != nil || ok || retry < time.Second {
		t.Fatalf("Sentinel-discovered promoted Redis must preserve exhausted shared budget: ok=%v retry=%s err=%v", ok, retry, err)
	}
	if ok, _, err := limiter.allow(ctx, key, policy, now.Add(policy.RefillEvery)); err != nil || !ok {
		t.Fatalf("Sentinel-discovered promoted Redis must continue canonical refill semantics: ok=%v err=%v", ok, err)
	}
}
