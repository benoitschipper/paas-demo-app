package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/ministerie-van-defensie/paas-demo-app/internal/config"
)

// DashboardData holds all data passed to the dashboard template.
type DashboardData struct {
	ThreatLevel config.Level
	Classes     config.TailwindClasses
	Timestamp   string
	Uptime      string
	Version     string
}

var startTime = time.Now()

// NewDashboardHandler returns an http.HandlerFunc that renders the Mission Status Dashboard.
// It takes the loaded config and the parsed template set.
func NewDashboardHandler(cfg *config.Config, tmpl *template.Template, version string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := DashboardData{
			ThreatLevel: cfg.ThreatLevel,
			Classes:     cfg.Classes(),
			Timestamp:   time.Now().UTC().Format("2006-01-02 15:04:05 UTC"),
			Uptime:      formatUptime(time.Since(startTime)),
			Version:     version,
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := tmpl.ExecuteTemplate(w, "dashboard.html", data); err != nil {
			http.Error(w, "template error: "+err.Error(), http.StatusInternalServerError)
		}
	}
}

// formatUptime formats a duration as "Xh Ym Zs".
func formatUptime(d time.Duration) string {
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	s := int(d.Seconds()) % 60
	if h > 0 {
		return fmt.Sprintf("%dh %dm %ds", h, m, s)
	}
	if m > 0 {
		return fmt.Sprintf("%dm %ds", m, s)
	}
	return fmt.Sprintf("%ds", s)
}
