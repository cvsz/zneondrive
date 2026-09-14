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
	"sync/atomic"
	"time"
)

// redisRateLimiter provides a cross-process token bucket using one atomic Lua
// script per decision. Redis is coordination only: it never becomes a source
// of gameplay authority or durable account/race state.
type redisRateLimiter struct {
	addr          string
	sentinelAddrs []string
	sentinelMaster string
	dialTimeout   time.Duration
	ioTimeout     time.Duration
}

type redisRateLimitStats struct {
	Allowed  uint64
	Rejected uint64
	Errors   uint64
}

var redisRateLimitCounters struct {
	allowed  atomic.Uint64
	rejected atomic.Uint64
	errors   atomic.Uint64
}

const redisRateLimitScript = `
local key = KEYS[1]
local burst = tonumber(ARGV[1])
local refill_ms = tonumber(ARGV[2])
local now_ms = tonumber(ARGV[3])
local ttl_ms = tonumber(ARGV[4])
local values = redis.call('HMGET', key, 'tokens', 'updated_ms')
local tokens = tonumber(values[1])
local updated = tonumber(values[2])
if tokens == nil or updated == nil then
  tokens = burst
  updated = now_ms
end
local elapsed = now_ms - updated
if elapsed > 0 then
  tokens = math.min(burst, tokens + (elapsed / refill_ms))
end
local allowed = 0
local retry_ms = 0
if tokens >= 1 then
  tokens = tokens - 1
  allowed = 1
else
  retry_ms = math.ceil((1 - tokens) * refill_ms)
  if retry_ms < 1 then retry_ms = 1 end
end
redis.call('HSET', key, 'tokens', tokens, 'updated_ms', now_ms)
redis.call('PEXPIRE', key, ttl_ms)
return {allowed, retry_ms}
`

func newRedisRateLimiter(addr string) *redisRateLimiter {
	return &redisRateLimiter{
		addr:        strings.TrimSpace(addr),
		dialTimeout: 300 * time.Millisecond,
		ioTimeout:   500 * time.Millisecond,
	}
}

func newRedisSentinelRateLimiter(addrs []string, master string) *redisRateLimiter {
	clean := make([]string, 0, len(addrs))
	for _, addr := range addrs {
		if addr = strings.TrimSpace(addr); addr != "" {
			clean = append(clean, addr)
		}
	}
	return &redisRateLimiter{
		sentinelAddrs: clean,
		sentinelMaster: strings.TrimSpace(master),
		dialTimeout:   300 * time.Millisecond,
		ioTimeout:     500 * time.Millisecond,
	}
}

func redisRateLimitStatsSnapshot() redisRateLimitStats {
	return redisRateLimitStats{
		Allowed:  redisRateLimitCounters.allowed.Load(),
		Rejected: redisRateLimitCounters.rejected.Load(),
		Errors:   redisRateLimitCounters.errors.Load(),
	}
}

func (r *redisRateLimiter) allow(ctx context.Context, key string, policy rateLimitPolicy, now time.Time) (allowed bool, retry time.Duration, err error) {
	defer func() {
		switch {
		case err != nil:
			redisRateLimitCounters.errors.Add(1)
		case allowed:
			redisRateLimitCounters.allowed.Add(1)
		default:
			redisRateLimitCounters.rejected.Add(1)
		}
	}()

	if r == nil || (r.addr == "" && (len(r.sentinelAddrs) == 0 || r.sentinelMaster == "")) {
		return false, 0, errors.New("redis rate limiter is not configured")
	}
	if policy.Burst <= 0 || policy.RefillEvery <= 0 {
		return true, 0, nil
	}

	addr, err := r.resolveAddr(ctx)
	if err != nil {
		return false, 0, err
	}

	dialer := net.Dialer{Timeout: r.dialTimeout}
	conn, err := dialer.DialContext(ctx, "tcp", addr)
	if err != nil {
		return false, 0, fmt.Errorf("dial redis: %w", err)
	}
	defer conn.Close()

	deadline := time.Now().Add(r.ioTimeout)
	if ctxDeadline, ok := ctx.Deadline(); ok && ctxDeadline.Before(deadline) {
		deadline = ctxDeadline
	}
	if err = conn.SetDeadline(deadline); err != nil {
		return false, 0, fmt.Errorf("set redis deadline: %w", err)
	}

	refillMS := policy.RefillEvery.Milliseconds()
	if refillMS < 1 {
		refillMS = 1
	}
	ttlMS := refillMS * int64(policy.Burst+2)
	if ttlMS < 60_000 {
		ttlMS = 60_000
	}

	args := []string{
		"EVAL", redisRateLimitScript, "1", "zneondrive:ratelimit:" + key,
		strconv.Itoa(policy.Burst), strconv.FormatInt(refillMS, 10),
		strconv.FormatInt(now.UTC().UnixMilli(), 10), strconv.FormatInt(ttlMS, 10),
	}
	if err = writeRESPCommand(conn, args); err != nil {
		return false, 0, fmt.Errorf("write redis command: %w", err)
	}

	values, err := readRESPIntegerArray(bufio.NewReader(conn))
	if err != nil {
		return false, 0, fmt.Errorf("read redis response: %w", err)
	}
	if len(values) != 2 {
		return false, 0, fmt.Errorf("unexpected redis limiter response length %d", len(values))
	}
	if values[0] == 1 {
		return true, 0, nil
	}
	retry = time.Duration(values[1]) * time.Millisecond
	if retry < time.Second {
		retry = time.Second
	}
	return false, retry, nil
}

func (r *redisRateLimiter) resolveAddr(ctx context.Context) (string, error) {
	if r.addr != "" {
		return r.addr, nil
	}
	var lastErr error
	for _, sentinel := range r.sentinelAddrs {
		addr, err := r.resolveSentinelMaster(ctx, sentinel)
		if err == nil {
			return addr, nil
		}
		lastErr = err
	}
	if lastErr == nil {
		lastErr = errors.New("no Redis Sentinel endpoints configured")
	}
	return "", fmt.Errorf("resolve redis master via sentinel: %w", lastErr)
}

func (r *redisRateLimiter) resolveSentinelMaster(ctx context.Context, sentinel string) (string, error) {
	dialer := net.Dialer{Timeout: r.dialTimeout}
	conn, err := dialer.DialContext(ctx, "tcp", sentinel)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	deadline := time.Now().Add(r.ioTimeout)
	if ctxDeadline, ok := ctx.Deadline(); ok && ctxDeadline.Before(deadline) {
		deadline = ctxDeadline
	}
	if err := conn.SetDeadline(deadline); err != nil {
		return "", err
	}
	if err := writeRESPCommand(conn, []string{"SENTINEL", "get-master-addr-by-name", r.sentinelMaster}); err != nil {
		return "", err
	}
	values, err := readRESPBulkStringArray(bufio.NewReader(conn))
	if err != nil {
		return "", err
	}
	if len(values) != 2 || strings.TrimSpace(values[0]) == "" || strings.TrimSpace(values[1]) == "" {
		return "", fmt.Errorf("unexpected Sentinel master response")
	}
	if _, err := strconv.Atoi(values[1]); err != nil {
		return "", fmt.Errorf("invalid Sentinel master port: %w", err)
	}
	return net.JoinHostPort(values[0], values[1]), nil
}

func writeRESPCommand(w io.Writer, args []string) error {
	if _, err := fmt.Fprintf(w, "*%d\r\n", len(args)); err != nil {
		return err
	}
	for _, arg := range args {
		if _, err := fmt.Fprintf(w, "$%d\r\n%s\r\n", len(arg), arg); err != nil {
			return err
		}
	}
	return nil
}

func readRESPIntegerArray(r *bufio.Reader) ([]int64, error) {
	line, err := readRESPLine(r)
	if err != nil {
		return nil, err
	}
	if strings.HasPrefix(line, "-") {
		return nil, errors.New(strings.TrimPrefix(line, "-"))
	}
	if !strings.HasPrefix(line, "*") {
		return nil, fmt.Errorf("expected RESP array, got %q", line)
	}
	count, err := strconv.Atoi(strings.TrimPrefix(line, "*"))
	if err != nil || count < 0 {
		return nil, fmt.Errorf("invalid RESP array length %q", line)
	}
	values := make([]int64, 0, count)
	for i := 0; i < count; i++ {
		item, readErr := readRESPLine(r)
		if readErr != nil {
			return nil, readErr
		}
		if strings.HasPrefix(item, "-") {
			return nil, errors.New(strings.TrimPrefix(item, "-"))
		}
		if !strings.HasPrefix(item, ":") {
			return nil, fmt.Errorf("expected RESP integer, got %q", item)
		}
		value, parseErr := strconv.ParseInt(strings.TrimPrefix(item, ":"), 10, 64)
		if parseErr != nil {
			return nil, fmt.Errorf("invalid RESP integer %q: %w", item, parseErr)
		}
		values = append(values, value)
	}
	return values, nil
}

func readRESPBulkStringArray(r *bufio.Reader) ([]string, error) {
	line, err := readRESPLine(r)
	if err != nil {
		return nil, err
	}
	if strings.HasPrefix(line, "-") {
		return nil, errors.New(strings.TrimPrefix(line, "-"))
	}
	if !strings.HasPrefix(line, "*") {
		return nil, fmt.Errorf("expected RESP array, got %q", line)
	}
	count, err := strconv.Atoi(strings.TrimPrefix(line, "*"))
	if err != nil || count < 0 {
		return nil, fmt.Errorf("invalid RESP array length %q", line)
	}
	values := make([]string, 0, count)
	for i := 0; i < count; i++ {
		header, readErr := readRESPLine(r)
		if readErr != nil {
			return nil, readErr
		}
		if header == "$-1" {
			return nil, errors.New("RESP null bulk string")
		}
		if !strings.HasPrefix(header, "$") {
			return nil, fmt.Errorf("expected RESP bulk string, got %q", header)
		}
		length, parseErr := strconv.Atoi(strings.TrimPrefix(header, "$"))
		if parseErr != nil || length < 0 {
			return nil, fmt.Errorf("invalid RESP bulk string length %q", header)
		}
		payload := make([]byte, length+2)
		if _, readErr := io.ReadFull(r, payload); readErr != nil {
			return nil, readErr
		}
		if string(payload[length:]) != "\r\n" {
			return nil, errors.New("malformed RESP bulk string")
		}
		values = append(values, string(payload[:length]))
	}
	return values, nil
}

func readRESPLine(r *bufio.Reader) (string, error) {
	line, err := r.ReadString('\n')
	if err != nil {
		return "", err
	}
	if !strings.HasSuffix(line, "\r\n") {
		return "", errors.New("malformed RESP line")
	}
	return strings.TrimSuffix(line, "\r\n"), nil
}
