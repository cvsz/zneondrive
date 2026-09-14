package httpapi

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"strings"
	"testing"
	"time"
)

func startFakeSentinel(t *testing.T, masterHost, masterPort string) string {
	t.Helper()
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen fake Sentinel: %v", err)
	}
	t.Cleanup(func() { _ = listener.Close() })

	go func() {
		for {
			conn, acceptErr := listener.Accept()
			if acceptErr != nil {
				return
			}
			go func() {
				defer conn.Close()
				_ = conn.SetDeadline(time.Now().Add(2 * time.Second))
				reader := bufio.NewReader(conn)
				// SENTINEL get-master-addr-by-name <master> is encoded as one
				// RESP array header followed by three bulk-string header/value pairs.
				for i := 0; i < 7; i++ {
					if _, readErr := reader.ReadString('\n'); readErr != nil {
						return
					}
				}
				_, _ = fmt.Fprintf(conn, "*2\r\n$%d\r\n%s\r\n$%d\r\n%s\r\n", len(masterHost), masterHost, len(masterPort), masterPort)
			}()
		}
	}()

	return listener.Addr().String()
}

func TestRedisSentinelResolverRequiresMajorityAgreement(t *testing.T) {
	stale := startFakeSentinel(t, "10.0.0.10", "6379")
	majorityA := startFakeSentinel(t, "10.0.0.20", "6379")
	majorityB := startFakeSentinel(t, "10.0.0.20", "6379")

	limiter := newRedisSentinelRateLimiter([]string{stale, majorityA, majorityB}, "zneondrive-primary")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	addr, err := limiter.resolveAddr(ctx)
	if err != nil {
		t.Fatalf("resolve Sentinel majority: %v", err)
	}
	if want := "10.0.0.20:6379"; addr != want {
		t.Fatalf("resolved %q, want majority master %q", addr, want)
	}
}

func TestRedisSentinelResolverRejectsNoMajority(t *testing.T) {
	first := startFakeSentinel(t, "10.0.0.10", "6379")
	second := startFakeSentinel(t, "10.0.0.20", "6379")
	third := startFakeSentinel(t, "10.0.0.30", "6379")

	limiter := newRedisSentinelRateLimiter([]string{first, second, third}, "zneondrive-primary")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	_, err := limiter.resolveAddr(ctx)
	if err == nil {
		t.Fatal("expected Sentinel resolution to fail without majority agreement")
	}
	if !strings.Contains(err.Error(), "no majority agreement") {
		t.Fatalf("unexpected no-majority error: %v", err)
	}
}

func TestRedisSentinelResolverToleratesMinorityEndpointFailure(t *testing.T) {
	majorityA := startFakeSentinel(t, "10.0.0.20", "6379")
	majorityB := startFakeSentinel(t, "10.0.0.20", "6379")

	limiter := newRedisSentinelRateLimiter([]string{"127.0.0.1:1", majorityA, majorityB}, "zneondrive-primary")
	limiter.dialTimeout = 50 * time.Millisecond
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	addr, err := limiter.resolveAddr(ctx)
	if err != nil {
		t.Fatalf("resolve Sentinel majority with one failed endpoint: %v", err)
	}
	if want := "10.0.0.20:6379"; addr != want {
		t.Fatalf("resolved %q, want majority master %q", addr, want)
	}
}
