# SABnzbd Monitor

A lightweight web frontend for monitoring SABnzbd download queues in real time.

![SABnzbd Monitor Dashboard](/docs/screenshots/dashboard.png)

## What it does

- Polls the SABnzbd API and displays the download queue with speed, size, and time remaining
- Auto-refreshes the UI via HTMX at a configurable interval (minimum 2 s)
- Runs as a single static binary inside a distroless container — minimal attack surface
- All configuration via environment variables — no config files

## Tech stack

| | |
|---|---|
| **Language** | Go 1.26.1 (stdlib only, no external dependencies) |
| **Frontend** | HTMX 2.0.8 (fetched at image build time, not committed) |
| **Base image** | `gcr.io/distroless/static:nonroot` |
| **Requires** | A running SABnzbd instance |

## Project layout

```
main.go               — Entry point
internal/
  config/             — Env var loading and validation
  logger/             — Levelled logging, client IP, HTTP middleware
  sabnzbd/            — SABnzbd API client and response types
  server/             — HTTP routes and handlers
templates/            — index.html, status.html
static/               — CSS; htmx.min.js injected at build time
```

## Configuration

Environment variables only — no config file.

| Variable | Description | Default |
|---|---|---|
| `SABMON_SABNZBD_URL` | Full URL to SABnzbd (no trailing slash) | **required** |
| `SABMON_SABNZBD_API_KEY` | SABnzbd API key | **required** |
| `SABMON_REFRESH_INTERVAL` | Poll interval in seconds (min 2) | `5` |
| `SABMON_DEBUG` | Verbose debug logging | `false` |
| `SABMON_LOG_CLIENT_INFO` | Log client IP and User-Agent | `false` |
| `HOST_PORT` | Host port mapped to container port 5959 | `5959` |

## Running

```bash
cp .env.example .env   # fill in SABMON_SABNZBD_URL and SABMON_SABNZBD_API_KEY
make up                # build image and start container
make down              # stop and remove container
```

UI available at `http://localhost:5959` (or your `HOST_PORT`).

## Makefile

| Target | Action |
|---|---|
| `make build` | Build the Docker image |
| `make up` | Build and start via Docker Compose |
| `make down` | Stop and remove containers |
| `make test` | Run all tests (`go test ./... -v`) |
| `make coverage` | Run tests and print total coverage % |

## Development

**Requirements:** Go 1.26+, Docker with the Compose plugin.

```bash
make test      # run all tests
make coverage  # test + coverage summary
make build     # build Docker image
```