package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// Config holds application configuration
type Config struct {
	SabnzbdURL      string `json:"sabnzbd_url"`
	SabnzbdAPIKey   string `json:"sabnzbd_api_key"`
	RefreshInterval int    `json:"refresh_interval"` // in seconds
	Debug           bool   `json:"debug"`            // enable debug logging
	LogClientInfo   bool   `json:"log_client_info"`  // log client IP and user agent
}

// Environment variable names
const (
	EnvSabnzbdURL      = "SABMON_SABNZBD_URL"
	EnvSabnzbdAPIKey   = "SABMON_SABNZBD_API_KEY"
	EnvRefreshInterval = "SABMON_REFRESH_INTERVAL"
	EnvDebug           = "SABMON_DEBUG"
	EnvLogClientInfo   = "SABMON_LOG_CLIENT_INFO"
)

// Application constants
const (
	AppPort = "5959" // Fixed application port
)

// LoadConfig loads configuration from config.json in the current working directory
// and then overlays any environment variables that are set
func LoadConfig() (Config, error) {
	var config Config
	var configLoaded bool
	var configErr error

	// Try to load config from current working directory
	cwd, err := os.Getwd()
	if err != nil {
		log.Println("Warning: Could not get current working directory:", err)
	} else {
		configPath := filepath.Join(cwd, "config.json")
		file, err := os.Open(configPath)
		if err != nil {
			configErr = fmt.Errorf("could not open config file: %v", err)
		} else {
			defer file.Close()
			decoder := json.NewDecoder(file)
			if err := decoder.Decode(&config); err != nil {
				configErr = fmt.Errorf("could not parse config file: %v", err)
			} else {
				configLoaded = true
				log.Printf("Config loaded from %s", configPath)
			}
		}
	}

	// Apply environment variables (they take precedence)
	if envURL := os.Getenv(EnvSabnzbdURL); envURL != "" {
		config.SabnzbdURL = envURL
	}
	if envAPIKey := os.Getenv(EnvSabnzbdAPIKey); envAPIKey != "" {
		config.SabnzbdAPIKey = envAPIKey
	}
	// Port is fixed and no longer configurable
	if envRefresh := os.Getenv(EnvRefreshInterval); envRefresh != "" {
		if val, err := strconv.Atoi(envRefresh); err == nil {
			config.RefreshInterval = val
		} else {
			log.Printf("Warning: Invalid %s value '%s', must be a number", EnvRefreshInterval, envRefresh)
		}
	}
	if envDebug := os.Getenv(EnvDebug); envDebug != "" {
		config.Debug = envDebug == "1" || strings.ToLower(envDebug) == "true"
	}
	if envLogClient := os.Getenv(EnvLogClientInfo); envLogClient != "" {
		config.LogClientInfo = envLogClient == "1" || strings.ToLower(envLogClient) == "true"
	}

	// Validate the configuration
	if err := validateConfig(&config); err != nil {
		return config, err
	}

	// If we had a config error but no environment variables were used, return the error
	if !configLoaded && configErr != nil {
		log.Println("No environment variables set and config file couldn't be loaded:", configErr)
	}

	return config, nil
}

// validateConfig checks if the config has valid values and sets defaults if needed
func validateConfig(config *Config) error {
	// SabnzbdURL is required
	if config.SabnzbdURL == "" {
		return fmt.Errorf("sabnzbd URL is required (set via config or %s)", EnvSabnzbdURL)
	}

	// SabnzbdAPIKey is required
	if config.SabnzbdAPIKey == "" {
		return fmt.Errorf("sabnzbd API key is required (set via config or %s)", EnvSabnzbdAPIKey)
	}

	// RefreshInterval must be greater than 0, default to 5 if invalid
	if config.RefreshInterval <= 0 {
		config.RefreshInterval = 5
		log.Println("Invalid refresh interval, defaulting to 5 seconds")
	}

	return nil
}

// SabnzbdStatus represents the response from SABnzbd API
type SabnzbdStatus struct {
	Status string `json:"status"`
	Queue  Queue  `json:"queue"`
}

// Queue contains download queue information
type Queue struct {
	Status     string      `json:"status"`
	Speed      string      `json:"speed"`
	SizeLeft   string      `json:"sizeleft"`
	TimeLeft   string      `json:"timeleft"`
	Percentage string      `json:"percentage"`
	Slots      []QueueItem `json:"slots"`
}

// QueueItem represents a single item in the download queue
type QueueItem struct {
	Filename   string `json:"filename"`
	Status     string `json:"status"`
	SizeLeft   string `json:"sizeleft"`
	Percentage string `json:"percentage"`
	TimeLeft   string `json:"timeleft"`
}

// Helper function to get client IP, handling proxies
func getClientIP(r *http.Request) string {
	// Check for X-Forwarded-For header first
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		// The first IP in the list is the client
		return strings.Split(forwarded, ",")[0]
	}

	// Try to get IP from RemoteAddr
	ip := r.RemoteAddr
	// Strip port if present
	if strings.Contains(ip, ":") {
		ip = strings.Split(ip, ":")[0]
	}
	return ip
}

// Simple debug logger that outputs to stdout
func debugLog(debug bool, level, message string, r *http.Request, logClientInfo bool) {
	if !debug && level == "DEBUG" {
		return
	}

	log.Printf("[%s] %s", level, message)

	if logClientInfo && r != nil {
		clientIP := getClientIP(r)
		log.Printf("  Client: %s, User-Agent: %s, URI: %s",
			clientIP, r.UserAgent(), r.RequestURI)
	}
}

// Middleware to log HTTP requests
func loggingMiddleware(next http.Handler, debug bool, logClientInfo bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		debugLog(debug, "INFO", fmt.Sprintf("Request: %s %s", r.Method, r.URL.Path), r, logClientInfo)
		next.ServeHTTP(w, r)
	})
}

func main() {
	// Load and validate configuration
	config, err := LoadConfig()
	if err != nil {
		log.Fatalf("Configuration error: %v", err)
	}

	// Log startup information
	debugLog(true, "INFO", "Application starting", nil, false)
	if config.Debug {
		// Redact sensitive information
		safeConfig := fmt.Sprintf("%+v", config)
		safeConfig = strings.Replace(safeConfig, config.SabnzbdAPIKey, "[REDACTED]", 1)
		debugLog(true, "INFO", fmt.Sprintf("Configuration: %s", safeConfig), nil, false)
	}

	// Parse templates
	tmpl := template.Must(template.ParseFiles("templates/index.html", "templates/status.html"))

	// Create a mux for easier middleware use
	mux := http.NewServeMux()

	// Serve static files
	fs := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	// Main page handler
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		debugLog(config.Debug, "INFO", "Serving index page", r, config.LogClientInfo)

		tmpl.ExecuteTemplate(w, "index.html", map[string]interface{}{
			"RefreshInterval": config.RefreshInterval,
			"Debug":           config.Debug,
		})
	})

	// SABnzbd status handler
	mux.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		debugLog(config.Debug, "INFO", "Fetching SABnzbd status", r, config.LogClientInfo)

		status, err := fetchSabnzbdStatus(config)
		if err != nil {
			debugLog(config.Debug, "ERROR", fmt.Sprintf("Failed to fetch status: %v", err), r, config.LogClientInfo)
			http.Error(w, "Failed to fetch status", http.StatusInternalServerError)
			return
		}

		debugLog(config.Debug, "INFO", "SABnzbd status fetched successfully", r, config.LogClientInfo)
		tmpl.ExecuteTemplate(w, "status.html", status)
	})

	// Apply middleware
	handler := loggingMiddleware(mux, config.Debug, config.LogClientInfo)

	// Start server
	log.Printf("Server starting on http://localhost:%s", AppPort)
	log.Fatal(http.ListenAndServe(":"+AppPort, handler))
}

func fetchSabnzbdStatus(config Config) (*SabnzbdStatus, error) {
	url := fmt.Sprintf("%s/api?output=json&apikey=%s&mode=queue",
		config.SabnzbdURL, config.SabnzbdAPIKey)

	if config.Debug {
		// Don't log the full URL with API key for security reasons
		safeUrl := strings.Replace(url, config.SabnzbdAPIKey, "[REDACTED]", 1)
		debugLog(true, "DEBUG", fmt.Sprintf("Requesting SABnzbd API: %s", safeUrl), nil, false)
	}

	client := http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		debugLog(config.Debug, "ERROR", fmt.Sprintf("API request failed: %v", err), nil, false)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		debugLog(config.Debug, "ERROR", fmt.Sprintf("API returned non-OK status: %d", resp.StatusCode), nil, false)
		return nil, fmt.Errorf("API returned non-OK status: %d", resp.StatusCode)
	}

	if config.Debug {
		bodyBytes, _ := io.ReadAll(resp.Body)
		// Don't log full response potentially containing sensitive data
		debugLog(true, "DEBUG", fmt.Sprintf("API response received, status: %s, length: %d bytes",
			resp.Status, len(bodyBytes)), nil, false)

		// We need to recreate the response body as we've read it
		resp.Body.Close()
		resp.Body = io.NopCloser(strings.NewReader(string(bodyBytes)))
	}

	var status SabnzbdStatus
	err = json.NewDecoder(resp.Body).Decode(&status)
	if err != nil {
		debugLog(config.Debug, "ERROR", fmt.Sprintf("Failed to decode API response: %v", err), nil, false)
		return nil, err
	}

	return &status, nil
}
