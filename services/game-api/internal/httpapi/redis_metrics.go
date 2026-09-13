package httpapi

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"time"
)

// RedisServerStats is a bounded subset of Redis INFO values suitable for low-cardinality metrics.
// Redis remains ephemeral coordination only; these values never participate in authorization or gameplay decisions.
type RedisServerStats struct {
	ConnectedClients          uint64
	UsedMemoryBytes           uint64
	UsedMemoryPeakBytes       uint64
	RejectedConnectionsTotal  uint64
	EvictedKeysTotal          uint64
	KeyspaceHitsTotal         uint64
	KeyspaceMissesTotal       uint64
	InstantaneousOpsPerSecond uint64
}

// RedisServerStatsProvider supplies a bounded Redis server health snapshot.
type RedisServerStatsProvider interface {
	RedisServerStats(ctx context.Context) (RedisServerStats, error)
}

type redisServerStatsProbe struct {
	addr        string
	dialTimeout time.Duration
	ioTimeout   time.Duration
}

// NewRedisServerStatsProvider creates a small RESP/INFO probe without adding Redis as a durable authority.
func NewRedisServerStatsProvider(addr string) RedisServerStatsProvider {
	return &redisServerStatsProbe{
		addr:        strings.TrimSpace(addr),
		dialTimeout: 300 * time.Millisecond,
		ioTimeout:   500 * time.Millisecond,
	}
}

func (p *redisServerStatsProbe) RedisServerStats(ctx context.Context) (RedisServerStats, error) {
	if p == nil || p.addr == "" {
		return RedisServerStats{}, errors.New("redis server metrics probe is not configured")
	}
	dialer := net.Dialer{Timeout: p.dialTimeout}
	conn, err := dialer.DialContext(ctx, "tcp", p.addr)
	if err != nil {
		return RedisServerStats{}, fmt.Errorf("dial redis metrics probe: %w", err)
	}
	defer conn.Close()

	deadline := time.Now().Add(p.ioTimeout)
	if ctxDeadline, ok := ctx.Deadline(); ok && ctxDeadline.Before(deadline) {
		deadline = ctxDeadline
	}
	if err := conn.SetDeadline(deadline); err != nil {
		return RedisServerStats{}, fmt.Errorf("set redis metrics deadline: %w", err)
	}
	if err := writeRESPCommand(conn, []string{"INFO"}); err != nil {
		return RedisServerStats{}, fmt.Errorf("write redis INFO: %w", err)
	}
	payload, err := readRESPBulkString(bufio.NewReader(conn))
	if err != nil {
		return RedisServerStats{}, fmt.Errorf("read redis INFO: %w", err)
	}
	return parseRedisInfoStats(payload)
}

func readRESPBulkString(r *bufio.Reader) (string, error) {
	line, err := readRESPLine(r)
	if err != nil {
		return "", err
	}
	if strings.HasPrefix(line, "-") {
		return "", errors.New(strings.TrimPrefix(line, "-"))
	}
	if !strings.HasPrefix(line, "$") {
		return "", fmt.Errorf("expected RESP bulk string, got %q", line)
	}
	length, err := strconv.Atoi(strings.TrimPrefix(line, "$"))
	if err != nil || length < 0 {
		return "", fmt.Errorf("invalid RESP bulk length %q", line)
	}
	buf := make([]byte, length+2)
	if _, err := io.ReadFull(r, buf); err != nil {
		return "", err
	}
	if string(buf[length:]) != "\r\n" {
		return "", errors.New("malformed RESP bulk string terminator")
	}
	return string(buf[:length]), nil
}

func parseRedisInfoStats(payload string) (RedisServerStats, error) {
	values := make(map[string]uint64)
	wanted := map[string]struct{}{
		"connected_clients": {}, "used_memory": {}, "used_memory_peak": {},
		"rejected_connections": {}, "evicted_keys": {}, "keyspace_hits": {},
		"keyspace_misses": {}, "instantaneous_ops_per_sec": {},
	}
	for _, raw := range strings.Split(payload, "\n") {
		line := strings.TrimSpace(raw)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, value, ok := strings.Cut(line, ":")
		if !ok {
			continue
		}
		if _, keep := wanted[key]; !keep {
			continue
		}
		n, err := strconv.ParseUint(strings.TrimSpace(value), 10, 64)
		if err != nil {
			return RedisServerStats{}, fmt.Errorf("invalid Redis INFO value for %s: %w", key, err)
		}
		values[key] = n
	}
	return RedisServerStats{
		ConnectedClients: values["connected_clients"],
		UsedMemoryBytes: values["used_memory"],
		UsedMemoryPeakBytes: values["used_memory_peak"],
		RejectedConnectionsTotal: values["rejected_connections"],
		EvictedKeysTotal: values["evicted_keys"],
		KeyspaceHitsTotal: values["keyspace_hits"],
		KeyspaceMissesTotal: values["keyspace_misses"],
		InstantaneousOpsPerSecond: values["instantaneous_ops_per_sec"],
	}, nil
}
