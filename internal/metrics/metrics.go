// Package metrics registers and exposes Prometheus metrics for the demo app.
package metrics

import (
	"net/http"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// PageViewsTotal counts dashboard page loads.
	PageViewsTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "demo_page_views_total",
		Help: "Total number of Mission Status Dashboard page views.",
	})

	// ThreatLevel reflects the current threat level as a numeric gauge.
	// 0=GREEN, 1=AMBER, 2=RED, 3=BLACK
	ThreatLevel = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "demo_threat_level",
		Help: "Current threat level: 0=GREEN, 1=AMBER, 2=RED, 3=BLACK.",
	})

	// ActiveSessions tracks the approximate number of in-flight HTTP requests.
	ActiveSessions = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "demo_active_sessions",
		Help: "Approximate number of currently active HTTP sessions.",
	})

	// RequestDuration measures HTTP request latency per handler.
	RequestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "demo_request_duration_seconds",
		Help:    "HTTP request latency in seconds, labelled by handler.",
		Buckets: prometheus.DefBuckets,
	}, []string{"handler"})
)

// activeCount is an atomic counter used to track in-flight requests.
var activeCount int64

// SetThreatLevel sets the demo_threat_level gauge to the given numeric value.
func SetThreatLevel(v float64) {
	ThreatLevel.Set(v)
}

// InstrumentHandler wraps an http.Handler with Prometheus instrumentation:
// - records request duration in the demo_request_duration_seconds histogram
// - increments demo_active_sessions on entry, decrements on exit
// - increments demo_page_views_total for the "dashboard" handler
func InstrumentHandler(name string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Track active sessions
		atomic.AddInt64(&activeCount, 1)
		ActiveSessions.Set(float64(atomic.LoadInt64(&activeCount)))
		defer func() {
			atomic.AddInt64(&activeCount, -1)
			ActiveSessions.Set(float64(atomic.LoadInt64(&activeCount)))
		}()

		// Increment page view counter for the dashboard handler
		if name == "dashboard" {
			PageViewsTotal.Inc()
		}

		// Record request duration
		start := time.Now()
		rw := &statusResponseWriter{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rw, r)
		duration := time.Since(start).Seconds()

		RequestDuration.WithLabelValues(name).Observe(duration)
		_ = strconv.Itoa(rw.status) // suppress unused warning; status available for future use
	})
}

// statusResponseWriter wraps http.ResponseWriter to capture the status code.
type statusResponseWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusResponseWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}
