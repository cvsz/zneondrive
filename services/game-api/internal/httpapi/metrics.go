package httpapi

import (
	"context"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/cvsz/zneondrive/services/game-api/internal/store"
)

var httpDurationBuckets = []float64{0.005, 0.010, 0.025, 0.050, 0.100, 0.250, 0.500, 1.000, 2.500, 5.000}

type routeMetrics struct {
	requests        [6]uint64
	durationBuckets []uint64
	durationCount   uint64
	durationSum     float64
}

type postgresStatsProvider interface {
	PostgresPoolStats() store.PostgresPoolStats
}

type redisServerMetricSnapshot struct {
	configured bool
	stats      RedisServerStats
	err        error
}

type HTTPMetrics struct {
	mu            sync.RWMutex
	routes        map[string]*routeMetrics
	inFlight      int64
	postgresStats postgresStatsProvider
	redisStats    RedisServerStatsProvider
}

func NewHTTPMetrics(postgresStats ...postgresStatsProvider) *HTTPMetrics {
	metrics := &HTTPMetrics{routes: make(map[string]*routeMetrics)}
	if len(postgresStats) > 0 {
		metrics.postgresStats = postgresStats[0]
	}
	return metrics
}

func (m *HTTPMetrics) SetRedisServerStatsProvider(provider RedisServerStatsProvider) {
	if m == nil {
		return
	}
	m.redisStats = provider
}

func (m *HTTPMetrics) Wrap(next http.Handler) http.Handler {
	if m == nil {
		return next
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		route := metricsRoute(r)
		start := time.Now()
		m.changeInFlight(1)
		defer m.changeInFlight(-1)

		recorder := &metricsResponseWriter{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(recorder, r)
		m.observe(route, recorder.status, time.Since(start))
	})
}

func (m *HTTPMetrics) Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; version=0.0.4; charset=utf-8")
		w.Header().Set("Cache-Control", "no-store")
		_, _ = w.Write([]byte(m.render(r.Context())))
	})
}

type metricsResponseWriter struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func (w *metricsResponseWriter) WriteHeader(status int) {
	if w.wroteHeader {
		return
	}
	w.status = status
	w.wroteHeader = true
	w.ResponseWriter.WriteHeader(status)
}

func (w *metricsResponseWriter) Write(p []byte) (int, error) {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	return w.ResponseWriter.Write(p)
}

func (m *HTTPMetrics) changeInFlight(delta int64) {
	m.mu.Lock()
	m.inFlight += delta
	m.mu.Unlock()
}

func (m *HTTPMetrics) observe(route string, status int, elapsed time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()

	rm := m.routes[route]
	if rm == nil {
		rm = &routeMetrics{durationBuckets: make([]uint64, len(httpDurationBuckets))}
		m.routes[route] = rm
	}
	class := status / 100
	if class < 1 || class > 5 {
		class = 0
	}
	rm.requests[class]++
	seconds := elapsed.Seconds()
	for i, boundary := range httpDurationBuckets {
		if seconds <= boundary {
			rm.durationBuckets[i]++
		}
	}
	rm.durationCount++
	rm.durationSum += seconds
}

func (m *HTTPMetrics) render(ctx context.Context) string {
	redisSnapshot := redisServerMetricSnapshot{}
	if m.redisStats != nil {
		redisSnapshot.configured = true
		redisSnapshot.stats, redisSnapshot.err = m.redisStats.RedisServerStats(ctx)
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	var b strings.Builder
	b.WriteString("# HELP zneondrive_http_requests_total Total HTTP requests by bounded route and status class.\n")
	b.WriteString("# TYPE zneondrive_http_requests_total counter\n")

	routes := make([]string, 0, len(m.routes))
	for route := range m.routes {
		routes = append(routes, route)
	}
	sort.Strings(routes)
	for _, route := range routes {
		rm := m.routes[route]
		for class := 0; class <= 5; class++ {
			if rm.requests[class] == 0 {
				continue
			}
			label := "other"
			if class >= 1 && class <= 5 {
				label = fmt.Sprintf("%dxx", class)
			}
			fmt.Fprintf(&b, "zneondrive_http_requests_total{route=%q,status_class=%q} %d\n", route, label, rm.requests[class])
		}
	}

	b.WriteString("# HELP zneondrive_http_request_duration_seconds HTTP request latency by bounded route.\n")
	b.WriteString("# TYPE zneondrive_http_request_duration_seconds histogram\n")
	for _, route := range routes {
		rm := m.routes[route]
		for i, boundary := range httpDurationBuckets {
			fmt.Fprintf(&b, "zneondrive_http_request_duration_seconds_bucket{route=%q,le=\"%g\"} %d\n", route, boundary, rm.durationBuckets[i])
		}
		fmt.Fprintf(&b, "zneondrive_http_request_duration_seconds_bucket{route=%q,le=\"+Inf\"} %d\n", route, rm.durationCount)
		fmt.Fprintf(&b, "zneondrive_http_request_duration_seconds_sum{route=%q} %.9f\n", route, rm.durationSum)
		fmt.Fprintf(&b, "zneondrive_http_request_duration_seconds_count{route=%q} %d\n", route, rm.durationCount)
	}

	b.WriteString("# HELP zneondrive_http_in_flight_requests Current in-flight HTTP requests.\n")
	b.WriteString("# TYPE zneondrive_http_in_flight_requests gauge\n")
	fmt.Fprintf(&b, "zneondrive_http_in_flight_requests %d\n", m.inFlight)

	m.renderPostgresMetrics(&b)
	renderRedisRateLimitMetrics(&b)
	renderRedisServerMetrics(redisSnapshot, &b)
	return b.String()
}

func (m *HTTPMetrics) renderPostgresMetrics(b *strings.Builder) {
	if m.postgresStats == nil {
		return
	}
	stats := m.postgresStats.PostgresPoolStats()

	b.WriteString("# HELP zneondrive_postgres_pool_connections PostgreSQL pool connections by bounded state.\n")
	b.WriteString("# TYPE zneondrive_postgres_pool_connections gauge\n")
	fmt.Fprintf(b, "zneondrive_postgres_pool_connections{state=\"max\"} %d\n", stats.MaxConns)
	fmt.Fprintf(b, "zneondrive_postgres_pool_connections{state=\"total\"} %d\n", stats.TotalConns)
	fmt.Fprintf(b, "zneondrive_postgres_pool_connections{state=\"idle\"} %d\n", stats.IdleConns)
	fmt.Fprintf(b, "zneondrive_postgres_pool_connections{state=\"acquired\"} %d\n", stats.AcquiredConns)
	fmt.Fprintf(b, "zneondrive_postgres_pool_connections{state=\"constructing\"} %d\n", stats.ConstructingConns)

	b.WriteString("# HELP zneondrive_postgres_pool_acquires_total PostgreSQL pool acquire attempts.\n")
	b.WriteString("# TYPE zneondrive_postgres_pool_acquires_total counter\n")
	fmt.Fprintf(b, "zneondrive_postgres_pool_acquires_total %d\n", stats.AcquireCount)
	b.WriteString("# HELP zneondrive_postgres_pool_empty_acquires_total Acquires that found no immediately idle connection.\n")
	b.WriteString("# TYPE zneondrive_postgres_pool_empty_acquires_total counter\n")
	fmt.Fprintf(b, "zneondrive_postgres_pool_empty_acquires_total %d\n", stats.EmptyAcquireCount)
	b.WriteString("# HELP zneondrive_postgres_pool_canceled_acquires_total Canceled PostgreSQL pool acquire attempts.\n")
	b.WriteString("# TYPE zneondrive_postgres_pool_canceled_acquires_total counter\n")
	fmt.Fprintf(b, "zneondrive_postgres_pool_canceled_acquires_total %d\n", stats.CanceledAcquireCount)
	b.WriteString("# HELP zneondrive_postgres_pool_new_connections_total PostgreSQL connections created by the pool.\n")
	b.WriteString("# TYPE zneondrive_postgres_pool_new_connections_total counter\n")
	fmt.Fprintf(b, "zneondrive_postgres_pool_new_connections_total %d\n", stats.NewConnsCount)
	b.WriteString("# HELP zneondrive_postgres_pool_acquire_duration_seconds Cumulative time spent acquiring PostgreSQL pool connections.\n")
	b.WriteString("# TYPE zneondrive_postgres_pool_acquire_duration_seconds counter\n")
	fmt.Fprintf(b, "zneondrive_postgres_pool_acquire_duration_seconds %.9f\n", stats.AcquireDuration.Seconds())
}

func renderRedisRateLimitMetrics(b *strings.Builder) {
	stats := redisRateLimitStatsSnapshot()
	b.WriteString("# HELP zneondrive_redis_rate_limit_decisions_total Redis-backed rate-limit decisions by bounded outcome.\n")
	b.WriteString("# TYPE zneondrive_redis_rate_limit_decisions_total counter\n")
	fmt.Fprintf(b, "zneondrive_redis_rate_limit_decisions_total{outcome=\"allowed\"} %d\n", stats.Allowed)
	fmt.Fprintf(b, "zneondrive_redis_rate_limit_decisions_total{outcome=\"rejected\"} %d\n", stats.Rejected)
	fmt.Fprintf(b, "zneondrive_redis_rate_limit_decisions_total{outcome=\"error\"} %d\n", stats.Errors)
}

func renderRedisServerMetrics(snapshot redisServerMetricSnapshot, b *strings.Builder) {
	if !snapshot.configured {
		return
	}
	b.WriteString("# HELP zneondrive_redis_up Whether the configured Redis server metrics probe succeeded.\n")
	b.WriteString("# TYPE zneondrive_redis_up gauge\n")
	if snapshot.err != nil {
		b.WriteString("zneondrive_redis_up 0\n")
		return
	}
	stats := snapshot.stats
	b.WriteString("zneondrive_redis_up 1\n")
	b.WriteString("# HELP zneondrive_redis_connected_clients Connected Redis clients.\n# TYPE zneondrive_redis_connected_clients gauge\n")
	fmt.Fprintf(b, "zneondrive_redis_connected_clients %d\n", stats.ConnectedClients)
	b.WriteString("# HELP zneondrive_redis_memory_bytes Redis memory bytes by bounded kind.\n# TYPE zneondrive_redis_memory_bytes gauge\n")
	fmt.Fprintf(b, "zneondrive_redis_memory_bytes{kind=\"used\"} %d\n", stats.UsedMemoryBytes)
	fmt.Fprintf(b, "zneondrive_redis_memory_bytes{kind=\"peak\"} %d\n", stats.UsedMemoryPeakBytes)
	b.WriteString("# HELP zneondrive_redis_rejected_connections_total Redis rejected connections.\n# TYPE zneondrive_redis_rejected_connections_total counter\n")
	fmt.Fprintf(b, "zneondrive_redis_rejected_connections_total %d\n", stats.RejectedConnectionsTotal)
	b.WriteString("# HELP zneondrive_redis_evicted_keys_total Redis evicted keys.\n# TYPE zneondrive_redis_evicted_keys_total counter\n")
	fmt.Fprintf(b, "zneondrive_redis_evicted_keys_total %d\n", stats.EvictedKeysTotal)
	b.WriteString("# HELP zneondrive_redis_keyspace_total Redis keyspace hits and misses by bounded outcome.\n# TYPE zneondrive_redis_keyspace_total counter\n")
	fmt.Fprintf(b, "zneondrive_redis_keyspace_total{outcome=\"hit\"} %d\n", stats.KeyspaceHitsTotal)
	fmt.Fprintf(b, "zneondrive_redis_keyspace_total{outcome=\"miss\"} %d\n", stats.KeyspaceMissesTotal)
	b.WriteString("# HELP zneondrive_redis_instantaneous_ops_per_second Redis instantaneous operations per second.\n# TYPE zneondrive_redis_instantaneous_ops_per_second gauge\n")
	fmt.Fprintf(b, "zneondrive_redis_instantaneous_ops_per_second %d\n", stats.InstantaneousOpsPerSecond)
}

func metricsRoute(r *http.Request) string {
	path := r.URL.Path
	switch {
	case r.Method == http.MethodGet && path == "/healthz":
		return "health"
	case r.Method == http.MethodPost && path == "/v1/sessions/bootstrap":
		return "bootstrap"
	case r.Method == http.MethodGet && path == "/v1/state":
		return "state"
	case r.Method == http.MethodPost && path == "/v1/game-tickets":
		return "game-ticket"
	case r.Method == http.MethodPost && path == "/v1/internal/game-tickets/redeem":
		return "internal-ticket-redeem"
	case r.Method == http.MethodPost && path == "/v1/internal/races/start":
		return "race-start"
	case r.Method == http.MethodPost && strings.HasPrefix(path, "/v1/internal/races/") && strings.HasSuffix(path, "/checkpoints"):
		return "race-checkpoint"
	case r.Method == http.MethodPost && strings.HasPrefix(path, "/v1/internal/races/") && strings.HasSuffix(path, "/finish"):
		return "race-finish"
	case r.Method == http.MethodPost && strings.HasPrefix(path, "/v1/quests/") && strings.HasSuffix(path, "/complete"):
		return "quest-complete"
	case r.Method == http.MethodPost && strings.HasPrefix(path, "/v1/vehicles/") && strings.HasSuffix(path, "/builds"):
		return "build-revise"
	default:
		return "other"
	}
}
