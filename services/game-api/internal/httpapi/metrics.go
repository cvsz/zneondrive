package httpapi

import (
	"fmt"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"
)

var httpDurationBuckets = []float64{0.005, 0.010, 0.025, 0.050, 0.100, 0.250, 0.500, 1.000, 2.500, 5.000}

type routeMetrics struct {
	requests        [6]uint64
	durationBuckets []uint64
	durationCount   uint64
	durationSum     float64
}

type HTTPMetrics struct {
	mu       sync.RWMutex
	routes   map[string]*routeMetrics
	inFlight int64
}

func NewHTTPMetrics() *HTTPMetrics {
	return &HTTPMetrics{routes: make(map[string]*routeMetrics)}
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
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/plain; version=0.0.4; charset=utf-8")
		w.Header().Set("Cache-Control", "no-store")
		_, _ = w.Write([]byte(m.render()))
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

func (m *HTTPMetrics) render() string {
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
	return b.String()
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
