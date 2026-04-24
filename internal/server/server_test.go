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

func TestIndexRouteReturnsOK(t *testing.T) {
	cfg := config.Config{
		SabnzbdURL:      "http://localhost:8080",
		SabnzbdAPIKey:   "key",
		RefreshInterval: 10,
	}
	const indexTmpl = `interval={{.RefreshInterval}} debug={{.Debug}}`
	tmpl := template.Must(template.New("index.html").Parse(indexTmpl))
	template.Must(tmpl.New("status.html").Parse(``))

	handler := server.New(cfg, tmpl)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200 for GET /, got %d", rr.Code)
	}
	if ct := rr.Header().Get("Content-Type"); ct != "text/html; charset=utf-8" {
		t.Errorf("expected Content-Type text/html; charset=utf-8, got %q", ct)
	}
	body := rr.Body.String()
	if body != "interval=10 debug=false" {
		t.Errorf("unexpected index body: %q", body)
	}
}

func TestIndexRoutePassesRefreshInterval(t *testing.T) {
	cfg := config.Config{
		SabnzbdURL:      "http://localhost:8080",
		SabnzbdAPIKey:   "key",
		RefreshInterval: 30,
		Debug:           true,
	}
	const indexTmpl = `{{.RefreshInterval}}`
	tmpl := template.Must(template.New("index.html").Parse(indexTmpl))
	template.Must(tmpl.New("status.html").Parse(``))

	handler := server.New(cfg, tmpl)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if body := rr.Body.String(); body != "30" {
		t.Errorf("expected RefreshInterval '30' in body, got %q", body)
	}
}
