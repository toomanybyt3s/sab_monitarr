# TODO

## 🐛 Bugs

- [x] **Silent `ExecuteTemplate` errors** — `main.go`
- [x] **Double `resp.Body.Close()` in `fetchSabnzbdStatus`** — `main.go`
- [x] **`X-Forwarded-For` header blindly trusted in `getClientIP`** — `main.go`
- [x] **`go.mod` declares non-existent Go version `1.26.1`** — `go.mod`
- [x] **`getClientIP` breaks on IPv6 `RemoteAddr`** — `main.go`
- [x] **Config file errors silently swallowed** — `main.go` *(config file support removed)*

---

## ⚠️ General Issues & Best Practices

- [x] **No HTTP method guard on `/status` endpoint** — `main.go`
- [x] **`http.Client` created on every request** — `main.go`
- [x] **Hard-coded relative paths for templates and static files** — `main.go` *(Docker-only; paths are always relative to `/app`)*
- [x] **API key sent as a URL query parameter** — `main.go`
- [x] **No minimum refresh interval enforced** — `main.go`
- [x] **No `Content-Type` header set on template responses** — `main.go`
- [x] **HTMX loaded from an external CDN** — `templates/index.html` *(now downloaded at image build time via `wget`)*
- [x] **Missing `.env.example` file**
- [x] **Go version mismatch between `go.mod` and `Dockerfile`** — `go.mod`, `Dockerfile`
- [x] **`config.json` volume in `docker-compose.yml`** — `docker-compose.yml` *(volumes block removed)*
- [x] **`VOLUME ["/app/config.json"]` in `Dockerfile`** — `Dockerfile` *(removed)*
- [x] **`TestLoadConfigFromEnvironment` lacks file-based and edge-case coverage** — `main_test.go` *(file-based tests removed; env-var coverage retained)*
- [x] **`TestGetClientIP` missing IPv6 test case** — `main_test.go`

---

## ✅ Architectural Changes

- [x] **Remove config file support** — `main.go`, `main_test.go`, `docker-compose.yml`
  Configuration is now exclusively via environment variables.

- [x] **HTMX fetched at Docker build time** — `Dockerfile`
  HTMX `2.0.8` is downloaded with `wget` during the builder stage and placed into `static/`. It is not committed to the repository.

- [x] **Update HTMX version to 2.0.8** — `Dockerfile`
  Previous vendored version was `2.0.4`.
