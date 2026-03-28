package config_test

import (
	"os"
	"testing"

	"github.com/toomanybyt3s/sab_monitarr/internal/config"
)

func TestValidate(t *testing.T) {
	valid := config.Config{
		SabnzbdURL:      "http://localhost:8080",
		SabnzbdAPIKey:   "test-api-key",
		RefreshInterval: 5,
	}
	if err := config.Validate(&valid); err != nil {
		t.Errorf("expected valid config to pass, got: %v", err)
	}

	noURL := config.Config{SabnzbdAPIKey: "key", RefreshInterval: 5}
	if err := config.Validate(&noURL); err == nil {
		t.Error("expected error for missing URL, got nil")
	}

	noKey := config.Config{SabnzbdURL: "http://localhost:8080", RefreshInterval: 5}
	if err := config.Validate(&noKey); err == nil {
		t.Error("expected error for missing API key, got nil")
	}

	zeroCfg := config.Config{
		SabnzbdURL:      "http://localhost:8080",
		SabnzbdAPIKey:   "key",
		RefreshInterval: 0,
	}
	if err := config.Validate(&zeroCfg); err != nil {
		t.Fatalf("expected fixable config to pass, got: %v", err)
	}
	if zeroCfg.RefreshInterval != 5 {
		t.Errorf("expected default 5, got %d", zeroCfg.RefreshInterval)
	}

	lowCfg := config.Config{
		SabnzbdURL:      "http://localhost:8080",
		SabnzbdAPIKey:   "key",
		RefreshInterval: 1,
	}
	if err := config.Validate(&lowCfg); err != nil {
		t.Fatalf("expected low-interval config to pass, got: %v", err)
	}
	if lowCfg.RefreshInterval != config.MinRefreshInterval {
		t.Errorf("expected %d, got %d", config.MinRefreshInterval, lowCfg.RefreshInterval)
	}
}

func TestLoad(t *testing.T) {
	os.Setenv(config.EnvSabnzbdURL, "http://test-env:8080")
	os.Setenv(config.EnvSabnzbdAPIKey, "env-api-key")
	os.Setenv(config.EnvRefreshInterval, "10")
	os.Setenv(config.EnvDebug, "true")
	defer func() {
		os.Unsetenv(config.EnvSabnzbdURL)
		os.Unsetenv(config.EnvSabnzbdAPIKey)
		os.Unsetenv(config.EnvRefreshInterval)
		os.Unsetenv(config.EnvDebug)
	}()

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if cfg.SabnzbdURL != "http://test-env:8080" {
		t.Errorf("expected URL from env, got %s", cfg.SabnzbdURL)
	}
	if cfg.SabnzbdAPIKey != "env-api-key" {
		t.Errorf("expected API key from env, got %s", cfg.SabnzbdAPIKey)
	}
	if cfg.RefreshInterval != 10 {
		t.Errorf("expected interval 10, got %d", cfg.RefreshInterval)
	}
	if !cfg.Debug {
		t.Error("expected debug true")
	}
}

func TestLoadEnvVarsOnly(t *testing.T) {
	os.Setenv(config.EnvSabnzbdURL, "http://env-override:9090")
	os.Setenv(config.EnvSabnzbdAPIKey, "override-key")
	defer func() {
		os.Unsetenv(config.EnvSabnzbdURL)
		os.Unsetenv(config.EnvSabnzbdAPIKey)
	}()

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if cfg.SabnzbdURL != "http://env-override:9090" {
		t.Errorf("expected env URL, got %s", cfg.SabnzbdURL)
	}
	if cfg.SabnzbdAPIKey != "override-key" {
		t.Errorf("expected env API key, got %s", cfg.SabnzbdAPIKey)
	}
}
