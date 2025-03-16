# SABnzbd Monitor

A simple web application to monitor your SABnzbd downloads.

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
