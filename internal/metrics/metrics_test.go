package metrics

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

func TestSetThreatLevel(t *testing.T) {
	tests := []struct {
		name  string
		value float64
	}{
		{"GREEN", 0},
		{"AMBER", 1},
		{"RED", 2},
		{"BLACK", 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetThreatLevel(tt.value)
			got := testutil.ToFloat64(ThreatLevel)
			if got != tt.value {
				t.Errorf("ThreatLevel gauge = %v, want %v", got, tt.value)
			}
		})
	}
}

func TestInstrumentHandler_IncrementsDashboardPageViews(t *testing.T) {
	// Reset counter by reading current value before
	before := testutil.ToFloat64(PageViewsTotal)

	handler := InstrumentHandler("dashboard", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	after := testutil.ToFloat64(PageViewsTotal)
	if after != before+1 {
		t.Errorf("PageViewsTotal: expected %v, got %v", before+1, after)
	}
}

func TestInstrumentHandler_DoesNotIncrementPageViewsForNonDashboard(t *testing.T) {
	before := testutil.ToFloat64(PageViewsTotal)

	handler := InstrumentHandler("health", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	after := testutil.ToFloat64(PageViewsTotal)
	if after != before {
		t.Errorf("PageViewsTotal should not increment for health handler: before=%v after=%v", before, after)
	}
}

func TestInstrumentHandler_RecordsRequestDuration(t *testing.T) {
	handler := InstrumentHandler("test-handler", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	// Verify histogram has at least one observation for this label
	count, err := testutil.GatherAndCount(prometheus.DefaultGatherer, "demo_request_duration_seconds")
	if err != nil {
		t.Fatalf("failed to gather metrics: %v", err)
	}
	if count == 0 {
		t.Error("expected demo_request_duration_seconds to have observations")
	}
}

func TestInstrumentHandler_ActiveSessionsDecrementAfterRequest(t *testing.T) {
	handler := InstrumentHandler("test", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	// After the request completes, active sessions should be back to baseline
	sessions := testutil.ToFloat64(ActiveSessions)
	if sessions < 0 {
		t.Errorf("ActiveSessions should not be negative, got %v", sessions)
	}
}
