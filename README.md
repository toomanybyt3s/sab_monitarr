# SABnzbd Monitor

A simple, lightweight frontend for monitoring your SABnzbd downloads. Designed to give users an easy way to view recently requested media and track download progress.

![SABnzbd Monitor Dashboard](/docs/screenshots/dashboard.png)

## Features

- **Real-time updates**: Monitors download progress with configurable refresh intervals
- **Clean, responsive UI**: Works on desktop and mobile devices
- **Minimal resource usage**: Lightweight Go binary in a distroless container image
- **Zero external JavaScript dependencies**: HTMX (v2.0.8) is downloaded at image build time — no CDN required at runtime
- **Docker-only deployment**: Designed to run exclusively via Docker / Docker Compose
- **Environment-variable configuration**: No config files — all settings passed as env vars

## Technology

- **Backend**: Go (stdlib only, no external dependencies)
- **Frontend**: HTMX for dynamic content updates without JavaScript frameworks
- **Container**: Multi-stage Docker build; final image is `gcr.io/distroless/static:nonroot`
- **Dependencies**: Only requires an existing SABnzbd instance — nothing else

## Project Structure

```
internal/
  config/     — Config struct, env var loading and validation
  logger/     — Debug logging, client IP extraction, HTTP middleware
  sabnzbd/    — SABnzbd API client and response types
  server/     — HTTP routes, handlers, server startup
main.go       — Entry point (calls server.Run)
templates/    — Go HTML templates (index.html, status.html)
static/       — CSS and assets (htmx.min.js injected at build time)
```

## Configuration

Configuration is done **exclusively via environment variables**. There is no config file.

| Variable | Description | Default |
|----------|-------------|---------|
| `SABMON_SABNZBD_URL` | Full URL to your SABnzbd instance | *Required* |
| `SABMON_SABNZBD_API_KEY` | SABnzbd API key | *Required* |
| `SABMON_REFRESH_INTERVAL` | UI poll interval in seconds (minimum: 2) | `5` |
| `SABMON_DEBUG` | Enable verbose debug logging | `false` |
| `SABMON_LOG_CLIENT_INFO` | Log client IP and User-Agent per request | `false` |
| `HOST_PORT` | Host port to expose the web UI on | `5959` |

> The container always listens on port **5959**. Use `HOST_PORT` to remap it on the host.

## Running with Docker Compose

```bash
# Copy the example env file and fill in your values
cp .env.example .env
$EDITOR .env

# Build and start
make up

# Stop
make down
```

The UI will be available at `http://localhost:5959` (or whatever `HOST_PORT` is set to).

### docker-compose.yml environment variables

```yaml
environment:
  - SABMON_SABNZBD_URL=${SABMON_SABNZBD_URL}
  - SABMON_SABNZBD_API_KEY=${SABMON_SABNZBD_API_KEY}
  - SABMON_REFRESH_INTERVAL=${SABMON_REFRESH_INTERVAL:-5}
  - SABMON_DEBUG=${SABMON_DEBUG:-false}
  - SABMON_LOG_CLIENT_INFO=${SABMON_LOG_CLIENT_INFO:-false}
```

## Makefile Targets

| Target | Description |
|--------|-------------|
| `make build` | Build the Docker image |
| `make up` | Build and start the service with Docker Compose |
| `make down` | Stop and remove containers |
| `make test` | Run all Go tests with verbose output |
| `make coverage` | Run tests and print total coverage percentage |

## Development

### Prerequisites

- Go 1.26+
- Docker with the Compose plugin (`docker compose`)

### Running tests

```bash
make test
```

### Coverage report

```bash
make coverage
```

### Building the Docker image manually

```bash
make build
```

HTMX v2.0.8 is fetched via `wget` during the Docker build stage and embedded into the image — no internet access is required at runtime and the file is never committed to the repository.