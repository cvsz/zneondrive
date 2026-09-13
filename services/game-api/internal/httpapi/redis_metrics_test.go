package httpapi

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type fakeRedisServerStatsProvider struct {
	stats RedisServerStats
	err   error
}

func (f fakeRedisServerStatsProvider) RedisServerStats(context.Context) (RedisServerStats, error) {
	return f.stats, f.err
}

func TestParseRedisInfoStatsUsesBoundedNumericFields(t *testing.T) {
	payload := "# Clients\r\nconnected_clients:4\r\n" +
		"# Memory\r\nused_memory:1024\r\nused_memory_peak:2048\r\n" +
		"# Stats\r\nrejected_connections:2\r\nevicted_keys:3\r\nkeyspace_hits:40\r\nkeyspace_misses:5\r\ninstantaneous_ops_per_sec:17\r\n" +
		"master_replid:secret-dynamic-value\r\n"
	stats, err := parseRedisInfoStats(payload)
	if err != nil {
		t.Fatalf("parse Redis INFO: %v", err)
	}
	if stats.ConnectedClients != 4 || stats.UsedMemoryBytes != 1024 || stats.UsedMemoryPeakBytes != 2048 ||
		stats.RejectedConnectionsTotal != 2 || stats.EvictedKeysTotal != 3 || stats.KeyspaceHitsTotal != 40 ||
		stats.KeyspaceMissesTotal != 5 || stats.InstantaneousOpsPerSecond != 17 {
		t.Fatalf("unexpected stats: %+v", stats)
	}
}

func TestRedisServerStatsProbeReadsInfoOverRESP(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()

	done := make(chan error, 1)
	go func() {
		conn, acceptErr := listener.Accept()
		if acceptErr != nil {
			done <- acceptErr
			return
		}
		defer conn.Close()
		reader := bufio.NewReader(conn)
		first, readErr := readRESPLine(reader)
		if readErr != nil {
			done <- readErr
			return
		}
		if first != "*1" {
			done <- fmt.Errorf("unexpected command array %q", first)
			return
		}
		if _, readErr = readRESPLine(reader); readErr != nil {
			done <- readErr
			return
		}
		cmd, readErr := readRESPLine(reader)
		if readErr != nil {
			done <- readErr
			return
		}
		if cmd != "INFO" {
			done <- fmt.Errorf("unexpected command %q", cmd)
			return
		}
		payload := "connected_clients:2\r\nused_memory:4096\r\nused_memory_peak:8192\r\nrejected_connections:1\r\nevicted_keys:0\r\nkeyspace_hits:9\r\nkeyspace_misses:2\r\ninstantaneous_ops_per_sec:11\r\n"
		_, writeErr := fmt.Fprintf(conn, "$%d\r\n%s\r\n", len(payload), payload)
		done <- writeErr
	}()

	provider := NewRedisServerStatsProvider(listener.Addr().String())
	stats, err := provider.RedisServerStats(context.Background())
	if err != nil {
		t.Fatalf("probe Redis INFO: %v", err)
	}
	if err := <-done; err != nil {
		t.Fatalf("fake Redis server: %v", err)
	}
	if stats.ConnectedClients != 2 || stats.UsedMemoryBytes != 4096 || stats.KeyspaceHitsTotal != 9 {
		t.Fatalf("unexpected probe stats: %+v", stats)
	}
}

func TestHTTPMetricsExposeBoundedRedisServerStats(t *testing.T) {
	metrics := NewHTTPMetrics()
	metrics.SetRedisServerStatsProvider(fakeRedisServerStatsProvider{stats: RedisServerStats{
		ConnectedClients: 3, UsedMemoryBytes: 1000, UsedMemoryPeakBytes: 2000,
		RejectedConnectionsTotal: 4, EvictedKeysTotal: 5, KeyspaceHitsTotal: 60,
		KeyspaceMissesTotal: 7, InstantaneousOpsPerSecond: 8,
	}})
	rec := httptest.NewRecorder()
	metrics.Handler().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/metrics", nil))
	body := rec.Body.String()
	for _, want := range []string{
		"zneondrive_redis_up 1",
		"zneondrive_redis_connected_clients 3",
		`zneondrive_redis_memory_bytes{kind="used"} 1000`,
		`zneondrive_redis_memory_bytes{kind="peak"} 2000`,
		"zneondrive_redis_rejected_connections_total 4",
		"zneondrive_redis_evicted_keys_total 5",
		`zneondrive_redis_keyspace_total{outcome="hit"} 60`,
		`zneondrive_redis_keyspace_total{outcome="miss"} 7`,
		"zneondrive_redis_instantaneous_ops_per_second 8",
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("missing Redis server metric %q: %s", want, body)
		}
	}
	for _, forbidden := range []string{"REDIS_ADDR", "redis://", "master_replid", "client_name", "session-token"} {
		if strings.Contains(body, forbidden) {
			t.Fatalf("Redis metrics exposed forbidden token %q: %s", forbidden, body)
		}
	}
}

func TestHTTPMetricsExposeRedisDownWithoutLeakingError(t *testing.T) {
	metrics := NewHTTPMetrics()
	metrics.SetRedisServerStatsProvider(fakeRedisServerStatsProvider{err: fmt.Errorf("dial redis://secret@host:6379 failed")})
	rec := httptest.NewRecorder()
	metrics.Handler().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/metrics", nil))
	body := rec.Body.String()
	if !strings.Contains(body, "zneondrive_redis_up 0") {
		t.Fatalf("expected Redis down metric: %s", body)
	}
	if strings.Contains(body, "secret@host") || strings.Contains(body, "dial redis") {
		t.Fatalf("Redis probe error must not be exported: %s", body)
	}
}
