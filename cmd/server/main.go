package main

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/ministerie-van-defensie/paas-demo-app/internal/config"
	"github.com/ministerie-van-defensie/paas-demo-app/internal/handlers"
	appmetrics "github.com/ministerie-van-defensie/paas-demo-app/internal/metrics"
	"github.com/ministerie-van-defensie/paas-demo-app/templates"
)

// version is set at build time via -ldflags.
var version = "dev"

func main() {
	// Load configuration
	cfg := config.Load()
	log.Printf("Starting paas-demo-app version=%s threat_level=%s port=%s", version, cfg.ThreatLevel, cfg.Port)

	// Set threat level gauge at startup
	appmetrics.SetThreatLevel(cfg.NumericLevel())

	// Parse embedded templates
	tmpl, err := template.ParseFS(templates.FS, "*.html")
	if err != nil {
		log.Fatalf("failed to parse templates: %v", err)
	}

	// Register routes
	mux := http.NewServeMux()

	// Dashboard — instrumented with metrics middleware
	mux.Handle("/", appmetrics.InstrumentHandler("dashboard",
		http.HandlerFunc(handlers.NewDashboardHandler(cfg, tmpl, version)),
	))

	// Prometheus metrics
	mux.Handle("/metrics", appmetrics.InstrumentHandler("metrics", promhttp.Handler()))

	// Health and readiness probes
	mux.Handle("/health", appmetrics.InstrumentHandler("health", http.HandlerFunc(handlers.Health)))
	mux.Handle("/ready", appmetrics.InstrumentHandler("ready", http.HandlerFunc(handlers.Ready)))

	// HTTP server with timeouts
	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Listening on :%s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	// Graceful shutdown on SIGTERM / SIGINT
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("forced shutdown: %v", err)
	}
	log.Println("Server stopped.")
}
