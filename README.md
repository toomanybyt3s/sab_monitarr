# SABnzbd Monitor

A simple, lightweight frontend for monitoring your SABnzbd downloads. Designed to give users an easy way to view recently requested media and track download progress.

![SABnzbd Monitor Dashboard](/docs/screenshots/dashboard.png)

## Features

- **Real-time updates**: Monitors download progress with configurable refresh intervals
- **Clean, responsive UI**: Works on desktop and mobile devices
- **Minimal resource usage**: Lightweight implementation with small footprint
- **Zero JavaScript dependencies**: Built with HTMX for dynamic content without bulky frameworks
- **Simple deployment**: Run as a standalone binary or with Docker

## Technology

- **Backend**: Written in Go (Golang) for high performance and low resource usage
- **Frontend**: Uses HTMX for dynamic content updates without JavaScript frameworks
- **Dependencies**: Only requires an existing SABnzbd service - nothing else!

## Configuration

SABnzbd Monitor can be configured in two ways:

1. Using a `config.json` file in the current working directory
2. Using environment variables (these take precedence over the config file)

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `SABMON_SABNZBD_URL` | URL to SABnzbd instance | *Required* |
| `SABMON_SABNZBD_API_KEY` | SABnzbd API key | *Required* |
| `SABMON_REFRESH_INTERVAL` | Refresh interval in seconds | `5` |
| `SABMON_DEBUG` | Enable debug logging | `false` |
| `SABMON_LOG_CLIENT_INFO` | Log client IP and user agent | `false` |

> Note: The application always runs on port 5959. Use Docker port mapping to change this if needed.

### Config File Format

```json
{
  "sabnzbd_url": "http://127.0.0.1:8080/sabnzbd",
  "sabnzbd_api_key": "your-api-key-here",
  "refresh_interval": 5,
  "debug": false,
  "log_client_info": false
}
```

## Running with Docker

```bash
docker run -p 8081:5959 \
  -e SABMON_SABNZBD_URL=http://sabnzbd:8080/sabnzbd \
  -e SABMON_SABNZBD_API_KEY=your-api-key-here \
  -e SABMON_DEBUG=true \
  ghcr.io/toomanybyt3s/sab_monitarr:latest
```

Or using docker-compose:

```bash
# Copy .env.example to .env and update values
cp .env.example .env
nano .env

# Run the application
docker-compose up -d
```

## Building from Source

```bash
go build -o sab_monitarr
./sab_monitarr
```

The application will be available at http://localhost:5959