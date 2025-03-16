package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestValidateConfig(t *testing.T) {
	// Test valid configuration
	validConfig := Config{
		SabnzbdURL:      "http://localhost:8080",
		SabnzbdAPIKey:   "test-api-key",
		RefreshInterval: 5,
	}

	if err := validateConfig(&validConfig); err != nil {
		t.Errorf("Expected valid config to pass validation, got error: %v", err)
	}

	// Test missing URL
	invalidConfig := Config{
		SabnzbdAPIKey:   "test-api-key",
		RefreshInterval: 5,
	}

	if err := validateConfig(&invalidConfig); err == nil {
		t.Error("Expected error for missing URL, got nil")
	}

	// Test missing API key
	invalidConfig = Config{
		SabnzbdURL:      "http://localhost:8080",
		RefreshInterval: 5,
	}

	if err := validateConfig(&invalidConfig); err == nil {
		t.Error("Expected error for missing API key, got nil")
	}

	// Test invalid refresh interval
	fixableConfig := Config{
		SabnzbdURL:      "http://localhost:8080",
		SabnzbdAPIKey:   "test-api-key",
		RefreshInterval: 0, // Invalid, should be set to default
	}

	if err := validateConfig(&fixableConfig); err != nil {
		t.Errorf("Expected fixable config to pass validation, got error: %v", err)
	}

	if fixableConfig.RefreshInterval != 5 {
		t.Errorf("Expected refresh interval to be set to default 5, got %d", fixableConfig.RefreshInterval)
	}
}

func TestLoadConfigFromEnvironment(t *testing.T) {
	// Set environment variables for testing
	os.Setenv(EnvSabnzbdURL, "http://test-env:8080")
	os.Setenv(EnvSabnzbdAPIKey, "env-api-key")
	os.Setenv(EnvRefreshInterval, "10")
	os.Setenv(EnvDebug, "true")

	// Clean up when test finishes
	defer func() {
		os.Unsetenv(EnvSabnzbdURL)
		os.Unsetenv(EnvSabnzbdAPIKey)
		os.Unsetenv(EnvRefreshInterval)
		os.Unsetenv(EnvDebug)
	}()

	config, err := LoadConfig()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Verify environment variables were loaded
	if config.SabnzbdURL != "http://test-env:8080" {
		t.Errorf("Expected URL from env var, got %s", config.SabnzbdURL)
	}

	if config.SabnzbdAPIKey != "env-api-key" {
		t.Errorf("Expected API key from env var, got %s", config.SabnzbdAPIKey)
	}

	if config.RefreshInterval != 10 {
		t.Errorf("Expected refresh interval 10 from env var, got %d", config.RefreshInterval)
	}

	if !config.Debug {
		t.Error("Expected debug to be true from env var")
	}
}

func TestGetClientIP(t *testing.T) {
	// Test with remote address only
	req1 := httptest.NewRequest("GET", "/", nil)
	req1.RemoteAddr = "192.168.1.100:12345"

	ip := getClientIP(req1)
	if ip != "192.168.1.100" {
		t.Errorf("Expected 192.168.1.100, got %s", ip)
	}

	// Test with X-Forwarded-For header
	req2 := httptest.NewRequest("GET", "/", nil)
	req2.RemoteAddr = "10.0.0.1:12345"
	req2.Header.Set("X-Forwarded-For", "192.168.1.100")

	ip = getClientIP(req2)
	if ip != "192.168.1.100" {
		t.Errorf("Expected 192.168.1.100 from X-Forwarded-For, got %s", ip)
	}
}

// Setup a mock SABnzbd API server for testing
func mockSabnzbdAPI() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{
			"status": "ok", 
			"queue": {
				"status": "Downloading", 
				"speed": "2.5 MB/s", 
				"sizeleft": "500 MB", 
				"timeleft": "00:03:20",
				"slots": [
					{
						"filename": "test_file.mkv",
						"status": "Downloading",
						"sizeleft": "500 MB",
						"percentage": "75",
						"timeleft": "00:03:20"
					}
				]
			}
		}`))
	}))
}

func TestFetchSabnzbdStatus(t *testing.T) {
	// Create a mock SABnzbd API server
	mockServer := mockSabnzbdAPI()
	defer mockServer.Close()

	// Configure to use the mock server
	config := Config{
		SabnzbdURL:      mockServer.URL,
		SabnzbdAPIKey:   "test-api-key",
		RefreshInterval: 5,
	}

	// Test fetching status
	status, err := fetchSabnzbdStatus(config)
	if err != nil {
		t.Fatalf("Failed to fetch SABnzbd status: %v", err)
	}

	// Verify the response was parsed correctly
	if status.Status != "ok" {
		t.Errorf("Expected status 'ok', got '%s'", status.Status)
	}

	if status.Queue.Status != "Downloading" {
		t.Errorf("Expected queue status 'Downloading', got '%s'", status.Queue.Status)
	}

	if status.Queue.Speed != "2.5 MB/s" {
		t.Errorf("Expected speed '2.5 MB/s', got '%s'", status.Queue.Speed)
	}

	if len(status.Queue.Slots) != 1 {
		t.Fatalf("Expected 1 slot, got %d", len(status.Queue.Slots))
	}

	if status.Queue.Slots[0].Filename != "test_file.mkv" {
		t.Errorf("Expected filename 'test_file.mkv', got '%s'", status.Queue.Slots[0].Filename)
	}

	if status.Queue.Slots[0].Percentage != "75" {
		t.Errorf("Expected percentage '75', got '%s'", status.Queue.Slots[0].Percentage)
	}
}
