package server

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/toomanybyt3s/sab_monitarr/internal/config"
	"github.com/toomanybyt3s/sab_monitarr/internal/logger"
	"github.com/toomanybyt3s/sab_monitarr/internal/sabnzbd"
)

// New wires up all routes and middleware and returns the root handler.
func New(cfg config.Config, tmpl *template.Template) http.Handler {
	mux := http.NewServeMux()

	// Static assets
	fs := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	// Index page
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		logger.Log(cfg.Debug, "INFO", "Serving index page", r, cfg.LogClientInfo)

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := tmpl.ExecuteTemplate(w, "index.html", map[string]interface{}{
			"RefreshInterval": cfg.RefreshInterval,
			"Debug":           cfg.Debug,
		}); err != nil {
			logger.Log(cfg.Debug, "ERROR", fmt.Sprintf("Template execution error: %v", err), r, cfg.LogClientInfo)
		}
	})

	// SABnzbd status (GET/HEAD only — polled by HTMX)
	mux.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		logger.Log(cfg.Debug, "INFO", "Fetching SABnzbd status", r, cfg.LogClientInfo)

		status, err := sabnzbd.FetchStatus(cfg.SabnzbdURL, cfg.SabnzbdAPIKey, cfg.Debug)
		if err != nil {
			logger.Log(cfg.Debug, "ERROR", fmt.Sprintf("Failed to fetch status: %v", err), r, cfg.LogClientInfo)
			http.Error(w, "Failed to fetch status", http.StatusInternalServerError)
			return
		}

		logger.Log(cfg.Debug, "INFO", "SABnzbd status fetched successfully", r, cfg.LogClientInfo)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := tmpl.ExecuteTemplate(w, "status.html", status); err != nil {
			logger.Log(cfg.Debug, "ERROR", fmt.Sprintf("Template execution error: %v", err), r, cfg.LogClientInfo)
		}
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

	logger.Log(true, "INFO", "Application starting", nil, false)
	if cfg.Debug {
		safe := strings.Replace(fmt.Sprintf("%+v", cfg), cfg.SabnzbdAPIKey, "[REDACTED]", 1)
		logger.Log(true, "INFO", fmt.Sprintf("Configuration: %s", safe), nil, false)
	}

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
