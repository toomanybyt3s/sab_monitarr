package server

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/toomanybyt3s/sab_monitarr/internal/config"
	"github.com/toomanybyt3s/sab_monitarr/internal/logger"
	"github.com/toomanybyt3s/sab_monitarr/internal/sabnzbd"
)

// New wires up all routes and middleware and returns the root handler.
func New(cfg config.Config, tmpl *template.Template) http.Handler {
	mux := http.NewServeMux()

	// Static assets — served with a long-lived cache header since files are baked into the image.
	fs := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		http.StripPrefix("/static/", fs).ServeHTTP(w, r)
	}))

	// Index page
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		logger.Log(cfg.Debug, "INFO", "Serving index page", r, cfg.LogClientInfo)

		var buf bytes.Buffer
		if err := tmpl.ExecuteTemplate(&buf, "index.html", map[string]interface{}{
			"RefreshInterval": cfg.RefreshInterval,
			"Debug":           cfg.Debug,
		}); err != nil {
			logger.Log(cfg.Debug, "ERROR", fmt.Sprintf("Template execution error: %v", err), r, cfg.LogClientInfo)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		buf.WriteTo(w)
	})

	// SABnzbd status (GET/HEAD only — polled by HTMX)
	mux.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		logger.Log(cfg.Debug, "INFO", "Fetching SABnzbd status", r, cfg.LogClientInfo)

		status, err := sabnzbd.FetchStatus(r.Context(), cfg.SabnzbdURL, cfg.SabnzbdAPIKey, cfg.Debug)
		if err != nil {
			logger.Log(cfg.Debug, "ERROR", fmt.Sprintf("Failed to fetch status: %v", err), r, cfg.LogClientInfo)
			http.Error(w, "Failed to fetch status", http.StatusInternalServerError)
			return
		}

		logger.Log(cfg.Debug, "INFO", "SABnzbd status fetched successfully", r, cfg.LogClientInfo)
		var buf bytes.Buffer
		if err := tmpl.ExecuteTemplate(&buf, "status.html", status); err != nil {
			logger.Log(cfg.Debug, "ERROR", fmt.Sprintf("Template execution error: %v", err), r, cfg.LogClientInfo)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		buf.WriteTo(w)
	})

	return logger.Middleware(mux, cfg.Debug, cfg.LogClientInfo)
}

// Run loads config, parses templates and starts the HTTP server. It only
// returns if the server fails to start.
func Run() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("configuration error: %w", err)
	}

	logger.Info("Application starting")
	logger.Debug(cfg.Debug, fmt.Sprintf(
		"Configuration: URL=%s RefreshInterval=%d Debug=%v LogClientInfo=%v APIKey=[REDACTED]",
		cfg.SabnzbdURL, cfg.RefreshInterval, cfg.Debug, cfg.LogClientInfo,
	))

	tmpl, err := template.ParseFiles(
		"templates/index.html",
		"templates/status.html",
	)
	if err != nil {
		return fmt.Errorf("failed to parse templates: %w", err)
	}

	handler := New(cfg, tmpl)

	log.Printf("Server starting on http://localhost:%s", config.AppPort)
	return http.ListenAndServe(":"+config.AppPort, handler)
}
