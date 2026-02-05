package metrics

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

// Metrics collects and exposes application metrics.
// It tracks request counts, error rates, latencies, and system resources.
type Metrics struct {
	startTime time.Time

	// Request counters
	totalRequests   atomic.Int64
	totalErrors     atomic.Int64
	activeRequests  atomic.Int64
	statusCounts    sync.Map // map[int]*atomic.Int64

	// Latency tracking (using histogram buckets)
	latencyMu      sync.Mutex
	latencyBuckets map[string]*atomic.Int64 // "<=10ms", "<=50ms", "<=100ms", "<=500ms", "<=1s", ">1s"
	latencySum      atomic.Int64 // total latency in microseconds

	// GraphQL specific
	queryCount       atomic.Int64
	mutationCount    atomic.Int64
	subscriptionCount atomic.Int64
	rateLimitHits    atomic.Int64
	authFailures     atomic.Int64
}

// New creates a new Metrics instance.
func New() *Metrics {
	m := &Metrics{
		startTime: time.Now(),
		latencyBuckets: map[string]*atomic.Int64{
			"le_10ms":  {},
			"le_50ms":  {},
			"le_100ms": {},
			"le_500ms": {},
			"le_1s":    {},
			"gt_1s":    {},
		},
	}
	return m
}

// RecordRequest records a completed request with its status code and duration.
func (m *Metrics) RecordRequest(statusCode int, duration time.Duration) {
	m.totalRequests.Add(1)

	if statusCode >= 400 {
		m.totalErrors.Add(1)
	}

	// Record status code count
	key := statusCode
	val, _ := m.statusCounts.LoadOrStore(key, &atomic.Int64{})
	val.(*atomic.Int64).Add(1)

	// Record latency bucket
	m.latencySum.Add(int64(duration.Microseconds()))
	switch {
	case duration <= 10*time.Millisecond:
		m.latencyBuckets["le_10ms"].Add(1)
	case duration <= 50*time.Millisecond:
		m.latencyBuckets["le_50ms"].Add(1)
	case duration <= 100*time.Millisecond:
		m.latencyBuckets["le_100ms"].Add(1)
	case duration <= 500*time.Millisecond:
		m.latencyBuckets["le_500ms"].Add(1)
	case duration <= 1*time.Second:
		m.latencyBuckets["le_1s"].Add(1)
	default:
		m.latencyBuckets["gt_1s"].Add(1)
	}
}

func (m *Metrics) RecordQuery()        { m.queryCount.Add(1) }
func (m *Metrics) RecordMutation()     { m.mutationCount.Add(1) }
func (m *Metrics) RecordSubscription() { m.subscriptionCount.Add(1) }
func (m *Metrics) RecordRateLimitHit() { m.rateLimitHits.Add(1) }
func (m *Metrics) RecordAuthFailure()  { m.authFailures.Add(1) }

func (m *Metrics) IncrementActive()  { m.activeRequests.Add(1) }
func (m *Metrics) DecrementActive()  { m.activeRequests.Add(-1) }

// MetricsMiddleware returns an HTTP middleware that records request metrics.
func (m *Metrics) MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		m.IncrementActive()
		defer m.DecrementActive()

		recorder := &statusRecorder{ResponseWriter: w, statusCode: 200}
		next.ServeHTTP(recorder, r)

		m.RecordRequest(recorder.statusCode, time.Since(start))
	})
}

type statusRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.statusCode = code
	r.ResponseWriter.WriteHeader(code)
}

// Handler returns an HTTP handler that exposes metrics as JSON.
// This endpoint can be scraped by monitoring systems.
func (m *Metrics) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var memStats runtime.MemStats
		runtime.ReadMemStats(&memStats)

		totalReqs := m.totalRequests.Load()
		avgLatencyUs := int64(0)
		if totalReqs > 0 {
			avgLatencyUs = m.latencySum.Load() / totalReqs
		}

		statusCounts := make(map[string]int64)
		m.statusCounts.Range(func(key, value interface{}) bool {
			statusCounts[fmt.Sprintf("%d", key.(int))] = value.(*atomic.Int64).Load()
			return true
		})

		latencyBuckets := make(map[string]int64)
		for k, v := range m.latencyBuckets {
			latencyBuckets[k] = v.Load()
		}

		data := map[string]interface{}{
			"uptime_seconds":   time.Since(m.startTime).Seconds(),
			"total_requests":   totalReqs,
			"total_errors":     m.totalErrors.Load(),
			"active_requests":  m.activeRequests.Load(),
			"error_rate":       errorRate(m.totalErrors.Load(), totalReqs),
			"avg_latency_ms":   float64(avgLatencyUs) / 1000.0,
			"status_codes":     statusCounts,
			"latency_buckets":  latencyBuckets,
			"graphql": map[string]int64{
				"queries":       m.queryCount.Load(),
				"mutations":     m.mutationCount.Load(),
				"subscriptions": m.subscriptionCount.Load(),
			},
			"security": map[string]int64{
				"rate_limit_hits": m.rateLimitHits.Load(),
				"auth_failures":   m.authFailures.Load(),
			},
			"system": map[string]interface{}{
				"goroutines":      runtime.NumGoroutine(),
				"heap_alloc_mb":   float64(memStats.HeapAlloc) / 1024 / 1024,
				"heap_sys_mb":     float64(memStats.HeapSys) / 1024 / 1024,
				"gc_cycles":       memStats.NumGC,
				"gc_pause_total_ms": float64(memStats.PauseTotalNs) / 1e6,
				"num_cpu":         runtime.NumCPU(),
			},
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		json.NewEncoder(w).Encode(data)
	}
}

func errorRate(errors, total int64) float64 {
	if total == 0 {
		return 0
	}
	return float64(errors) / float64(total) * 100
}
