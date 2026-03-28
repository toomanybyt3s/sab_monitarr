package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

// Config holds application configuration.
type Config struct {
	SabnzbdURL      string
	SabnzbdAPIKey   string
	RefreshInterval int  // in seconds
	Debug           bool // enable debug logging
	LogClientInfo   bool // log client IP and user agent
}

// Environment variable names.
const (
	EnvSabnzbdURL      = "SABMON_SABNZBD_URL"
	EnvSabnzbdAPIKey   = "SABMON_SABNZBD_API_KEY"
	EnvRefreshInterval = "SABMON_REFRESH_INTERVAL"
	EnvDebug           = "SABMON_DEBUG"
	EnvLogClientInfo   = "SABMON_LOG_CLIENT_INFO"
)

// Application constants.
const (
	AppPort            = "5959" // Fixed application port
	MinRefreshInterval = 2      // Minimum allowed refresh interval in seconds
)

// Load reads all configuration from environment variables.
func Load() (Config, error) {
	var cfg Config

	cfg.SabnzbdURL = os.Getenv(EnvSabnzbdURL)
	cfg.SabnzbdAPIKey = os.Getenv(EnvSabnzbdAPIKey)

	if envRefresh := os.Getenv(EnvRefreshInterval); envRefresh != "" {
		if val, err := strconv.Atoi(envRefresh); err == nil {
			cfg.RefreshInterval = val
		} else {
			log.Printf("Warning: Invalid %s value '%s', must be a number", EnvRefreshInterval, envRefresh)
		}
	}

	if envDebug := os.Getenv(EnvDebug); envDebug != "" {
		cfg.Debug = envDebug == "1" || strings.ToLower(envDebug) == "true"
	}
	if envLogClient := os.Getenv(EnvLogClientInfo); envLogClient != "" {
		cfg.LogClientInfo = envLogClient == "1" || strings.ToLower(envLogClient) == "true"
	}

	if err := Validate(&cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}

// Validate checks config values and sets defaults where needed.
func Validate(cfg *Config) error {
	if cfg.SabnzbdURL == "" {
		return fmt.Errorf("sabnzbd URL is required (set %s)", EnvSabnzbdURL)
	}

	if cfg.SabnzbdAPIKey == "" {
		return fmt.Errorf("sabnzbd API key is required (set %s)", EnvSabnzbdAPIKey)
	}

	if cfg.RefreshInterval <= 0 {
		cfg.RefreshInterval = 5
		log.Println("Invalid refresh interval, defaulting to 5 seconds")
	} else if cfg.RefreshInterval < MinRefreshInterval {
		cfg.RefreshInterval = MinRefreshInterval
		log.Printf("Refresh interval too low, raising to minimum %d seconds", MinRefreshInterval)
	}

	return nil
}
