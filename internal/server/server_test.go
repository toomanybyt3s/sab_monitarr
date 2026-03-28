package server_test

import (
	"html/template"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/toomanybyt3s/sab_monitarr/internal/config"
	"github.com/toomanybyt3s/sab_monitarr/internal/server"
)

func TestStatusEndpointMethodGuard(t *testing.T) {
	cfg := config.Config{
		SabnzbdURL:      "http://localhost:8080",
		SabnzbdAPIKey:   "key",
		RefreshInterval: 5,
	}
	tmpl := template.Must(template.New("index.html").Parse(``))
	template.Must(tmpl.New("status.html").Parse(``))

	handler := server.New(cfg, tmpl)

	for _, method := range []string{http.MethodPost, http.MethodPut, http.MethodDelete} {
		req := httptest.NewRequest(method, "/status", nil)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
		if rr.Code != http.StatusMethodNotAllowed {
			t.Errorf("expected 405 for %s, got %d", method, rr.Code)
		}
	}
}

func TestIndexNotFoundForUnknownPath(t *testing.T) {
	cfg := config.Config{
		SabnzbdURL:      "http://localhost:8080",
		SabnzbdAPIKey:   "key",
		RefreshInterval: 5,
	}
	tmpl := template.Must(template.New("index.html").Parse(``))
	template.Must(tmpl.New("status.html").Parse(``))

	handler := server.New(cfg, tmpl)

	req := httptest.NewRequest(http.MethodGet, "/nonexistent", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("expected 404 for unknown path, got %d", rr.Code)
	}
}
