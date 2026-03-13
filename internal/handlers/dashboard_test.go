package handlers

import (
	"html/template"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ministerie-van-defensie/paas-demo-app/internal/config"
	"github.com/ministerie-van-defensie/paas-demo-app/templates"
)

func loadTemplates(t *testing.T) *template.Template {
	t.Helper()
	tmpl, err := template.ParseFS(templates.FS, "*.html")
	if err != nil {
		t.Fatalf("failed to parse templates: %v", err)
	}
	return tmpl
}

func TestDashboardHandler_RendersWithGreenThreatLevel(t *testing.T) {
	t.Setenv("THREAT_LEVEL", "GREEN")
	cfg := config.Load()
	tmpl := loadTemplates(t)

	handler := NewDashboardHandler(cfg, tmpl, "test-build")
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	handler(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
	body := rr.Body.String()
	if !strings.Contains(body, "THREAT LEVEL: GREEN") {
		t.Error("expected 'THREAT LEVEL: GREEN' in body")
	}
	if !strings.Contains(body, "text-emerald-400") {
		t.Error("expected emerald CSS class for GREEN level")
	}
}

func TestDashboardHandler_RendersWithAmberThreatLevel(t *testing.T) {
	t.Setenv("THREAT_LEVEL", "AMBER")
	cfg := config.Load()
	tmpl := loadTemplates(t)

	handler := NewDashboardHandler(cfg, tmpl, "test-build")
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	handler(rr, req)

	body := rr.Body.String()
	if !strings.Contains(body, "THREAT LEVEL: AMBER") {
		t.Error("expected 'THREAT LEVEL: AMBER' in body")
	}
	if !strings.Contains(body, "text-amber-400") {
		t.Error("expected amber CSS class for AMBER level")
	}
}

func TestDashboardHandler_RendersWithRedThreatLevel(t *testing.T) {
	t.Setenv("THREAT_LEVEL", "RED")
	cfg := config.Load()
	tmpl := loadTemplates(t)

	handler := NewDashboardHandler(cfg, tmpl, "test-build")
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	handler(rr, req)

	body := rr.Body.String()
	if !strings.Contains(body, "THREAT LEVEL: RED") {
		t.Error("expected 'THREAT LEVEL: RED' in body")
	}
	if !strings.Contains(body, "text-red-400") {
		t.Error("expected red CSS class for RED level")
	}
	if !strings.Contains(body, "pulse-glow") {
		t.Error("expected pulse-glow animation for RED level")
	}
}

func TestDashboardHandler_RendersWithBlackThreatLevel(t *testing.T) {
	t.Setenv("THREAT_LEVEL", "BLACK")
	cfg := config.Load()
	tmpl := loadTemplates(t)

	handler := NewDashboardHandler(cfg, tmpl, "test-build")
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	handler(rr, req)

	body := rr.Body.String()
	if !strings.Contains(body, "THREAT LEVEL: BLACK") {
		t.Error("expected 'THREAT LEVEL: BLACK' in body")
	}
	if !strings.Contains(body, "text-white") {
		t.Error("expected white CSS class for BLACK level")
	}
	if !strings.Contains(body, "pulse-glow") {
		t.Error("expected pulse-glow animation for BLACK level")
	}
}

func TestDashboardHandler_ContentTypeIsHTML(t *testing.T) {
	t.Setenv("THREAT_LEVEL", "GREEN")
	cfg := config.Load()
	tmpl := loadTemplates(t)

	handler := NewDashboardHandler(cfg, tmpl, "test-build")
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	handler(rr, req)

	ct := rr.Header().Get("Content-Type")
	if !strings.HasPrefix(ct, "text/html") {
		t.Errorf("expected text/html content type, got %q", ct)
	}
}

func TestHealthHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()
	Health(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
	if !strings.Contains(rr.Body.String(), `"ok"`) {
		t.Error("expected 'ok' in health response body")
	}
}

func TestReadyHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/ready", nil)
	rr := httptest.NewRecorder()
	Ready(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
	if !strings.Contains(rr.Body.String(), `"ready"`) {
		t.Error("expected 'ready' in ready response body")
	}
}
