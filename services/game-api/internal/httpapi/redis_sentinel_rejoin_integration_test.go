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

func TestRedisSentinelRejoinsOldPrimaryAsReadOnlyReplica(t *testing.T) {
	if _, err := exec.LookPath("docker"); err != nil {
		if os.Getenv("CI") == "true" {
			t.Fatalf("docker is required for the Redis Sentinel rejoin evidence gate: %v", err)
		}
		t.Skipf("docker unavailable: %v", err)
	}

	primaryPort := reserveTCPPort(t)
	replicaPort := reserveTCPPort(t)
	sentinelPorts := []int{reserveTCPPort(t), reserveTCPPort(t), reserveTCPPort(t)}
	suffix := fmt.Sprintf("%d", time.Now().UnixNano())
	network := "zneondrive-redis-sentinel-rejoin-" + suffix
	primary := "zneondrive-sentinel-rejoin-primary-" + suffix
	replica := "zneondrive-sentinel-rejoin-replica-" + suffix
	masterName := "zneondrive-primary"
	sentinels := []string{
		"zneondrive-sentinel-rejoin-a-" + suffix,
		"zneondrive-sentinel-rejoin-b-" + suffix,
		"zneondrive-sentinel-rejoin-c-" + suffix,
	}

	runDocker(t, "network", "create", network)
	t.Cleanup(func() {
		args := append([]string{"rm", "-f", primary, replica}, sentinels...)
		_ = exec.Command("docker", args...).Run()
		_ = exec.Command("docker", "network", "rm", network).Run()
	})

	runDocker(t,
		"run", "-d", "--name", primary, "--network", network,
		"-p", fmt.Sprintf("127.0.0.1:%d:6379", primaryPort),
		"redis:8-alpine", "redis-server", "--save", "", "--appendonly", "no",
	)
	waitRedisContainerReady(t, primary)
	primaryIP := strings.TrimSpace(runDocker(t, "inspect", "-f", "{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}", primary))
	if primaryIP == "" {
		t.Fatal("Redis primary did not expose a Docker network address")
	}

	runDocker(t,
		"run", "-d", "--name", replica, "--network", network,
		"-p", fmt.Sprintf("127.0.0.1:%d:6379", replicaPort),
		"redis:8-alpine", "redis-server", "--save", "", "--appendonly", "no",
		"--replicaof", primaryIP, "6379",
	)
	waitRedisContainerReady(t, replica)
	waitRedisReplicationLink(t, replica)
	replicaIP := strings.TrimSpace(runDocker(t, "inspect", "-f", "{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}", replica))
	if replicaIP == "" {
		t.Fatal("Redis replica did not expose a Docker network address")
	}

	for i, sentinel := range sentinels {
		config := fmt.Sprintf("port 6379\nsentinel monitor %s %s 6379 2\nsentinel down-after-milliseconds %s 1000\nsentinel failover-timeout %s 10000\nsentinel parallel-syncs %s 1\n", masterName, primaryIP, masterName, masterName, masterName)
		runDocker(t,
			"run", "-d", "--name", sentinel, "--network", network,
			"-p", fmt.Sprintf("127.0.0.1:%d:6379", sentinelPorts[i]),
			"redis:8-alpine", "sh", "-c",
			fmt.Sprintf("printf '%%b' %q > /tmp/sentinel.conf && exec redis-server /tmp/sentinel.conf --sentinel", config),
		)
		waitRedisContainerReady(t, sentinel)
	}

	for _, sentinel := range sentinels {
		waitRedisSentinelQuorum(t, sentinel, masterName)
		waitRedisSentinelReplica(t, sentinel, masterName, replicaIP)
	}

	sentinelAddrs := make([]string, 0, len(sentinelPorts))
	for _, port := range sentinelPorts {
		sentinelAddrs = append(sentinelAddrs, fmt.Sprintf("127.0.0.1:%d", port))
	}
	limiter := newRedisSentinelRateLimiter(sentinelAddrs, masterName)
	policy := rateLimitPolicy{Burst: 2, RefillEvery: 30 * time.Second}
	now := time.Now().UTC().Truncate(time.Millisecond)
	key := "sentinel-rejoin:" + suffix
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	waitForRedisCondition(t, 30*time.Second, func() bool {
		addr, err := limiter.resolveAddr(ctx)
		return err == nil && addr == primaryIP+":6379"
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

	// Fail the original authority and wait for Sentinel to promote the replica.
	runDocker(t, "stop", "-t", "1", primary)
	waitRedisRole(t, replica, "role:master")
	waitForRedisCondition(t, 30*time.Second, func() bool {
		addr, err := limiter.resolveAddr(ctx)
		return err == nil && addr == replicaIP+":6379"
	}, "Sentinel never converged on the promoted Redis replica")

	// Reintroduce the old primary without manual REPLICAOF commands. Sentinel must
	// reconfigure it as a replica of the promoted authority before this evidence
	// considers the rejoin safe.
	runDocker(t, "start", primary)
	waitRedisContainerReady(t, primary)
	waitForRedisCondition(t, 45*time.Second, func() bool {
		out, err := exec.Command("docker", "exec", primary, "redis-cli", "INFO", "replication").CombinedOutput()
		text := string(out)
		return err == nil &&
			strings.Contains(text, "role:slave") &&
			strings.Contains(text, "master_host:"+replicaIP) &&
			strings.Contains(text, "master_link_status:up")
	}, "old Redis primary never rejoined as a connected replica of the promoted authority")

	// A rejoined old primary must be fenced from direct writes by its replica role.
	waitForRedisCondition(t, 15*time.Second, func() bool {
		out, _ := exec.Command("docker", "exec", primary, "redis-cli", "SET", "zneondrive:sentinel-rejoin-write-probe", "1").CombinedOutput()
		return strings.Contains(strings.ToUpper(string(out)), "READONLY")
	}, "rejoined old Redis primary accepted a direct write instead of remaining read-only")

	// The application must continue discovering the promoted authority and the
	// replicated limiter state must remain coherent after the old node rejoins.
	waitForRedisCondition(t, 15*time.Second, func() bool {
		addr, err := limiter.resolveAddr(ctx)
		return err == nil && addr == replicaIP+":6379"
	}, "Sentinel master discovery regressed to the rejoined old primary")
	waitRedisHashField(t, primary, redisKey, "tokens", "0")
	if ok, retry, err := limiter.allow(ctx, key, policy, now); err != nil || ok || retry < time.Second {
		t.Fatalf("promoted Redis must preserve exhausted budget after old-primary rejoin: ok=%v retry=%s err=%v", ok, retry, err)
	}
}
