package httpapi

import (
	"bufio"
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestWriteRESPCommand(t *testing.T) {
	var buf bytes.Buffer
	if err := writeRESPCommand(&buf, []string{"PING", "hello"}); err != nil {
		t.Fatalf("write command: %v", err)
	}
	want := "*2\r\n$4\r\nPING\r\n$5\r\nhello\r\n"
	if got := buf.String(); got != want {
		t.Fatalf("unexpected RESP encoding: %q", got)
	}
}

func TestReadRESPIntegerArray(t *testing.T) {
	values, err := readRESPIntegerArray(bufio.NewReader(strings.NewReader("*2\r\n:1\r\n:2500\r\n")))
	if err != nil {
		t.Fatalf("read response: %v", err)
	}
	if len(values) != 2 || values[0] != 1 || values[1] != 2500 {
		t.Fatalf("unexpected values: %#v", values)
	}
}

func TestDistributedLimiterFallsBackToLocalWhenRedisUnavailable(t *testing.T) {
	nextCalls := 0
	next := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		nextCalls++
		w.WriteHeader(http.StatusNoContent)
	})
	middleware := NewDistributedRateLimitedHandler(next, "127.0.0.1:1").(*rateLimitMiddleware)
	middleware.redis.dialTimeout = 10 * time.Millisecond
	middleware.redis.ioTimeout = 10 * time.Millisecond
	middleware.limiter.now = func() time.Time { return time.Date(2026, 9, 13, 0, 0, 0, 0, time.UTC) }

	for i := 0; i < bootstrapRatePolicy.Burst; i++ {
		req := httptest.NewRequest(http.MethodPost, "/v1/sessions/bootstrap", nil)
		req.RemoteAddr = "203.0.113.50:40000"
		rec := httptest.NewRecorder()
		middleware.ServeHTTP(rec, req)
		if rec.Code != http.StatusNoContent {
			t.Fatalf("fallback request %d expected 204, got %d", i+1, rec.Code)
		}
	}

	req := httptest.NewRequest(http.MethodPost, "/v1/sessions/bootstrap", nil)
	req.RemoteAddr = "203.0.113.50:40000"
	rec := httptest.NewRecorder()
	middleware.ServeHTTP(rec, req)
	if rec.Code != http.StatusTooManyRequests {
		t.Fatalf("local fallback must still enforce limit, got %d", rec.Code)
	}
	if nextCalls != bootstrapRatePolicy.Burst {
		t.Fatalf("limited fallback request reached handler: %d", nextCalls)
	}
}
