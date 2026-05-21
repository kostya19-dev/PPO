package appconfig

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad_ConfigFileAndEnvOverride(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "app.json")

	content := `{
		"app": {"name": "Test App"},
		"api": {"address": ":9090"},
		"database": {"url": "postgres://from-file"},
		"logging": {"level": "warn", "file": "./logs/test.log"},
		"business": {"appointment_min_lead_time_minutes": 45}
	}`

	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write config file: %v", err)
	}

	t.Setenv("DATABASE_URL", "postgres://from-env")
	t.Setenv("VITALINE_APP_NAME", "Env App")

	cfg, resolved, err := Load(path)
	if err != nil {
		t.Fatalf("load config: %v", err)
	}

	if resolved != path {
		t.Fatalf("unexpected config path: %s", resolved)
	}
	if cfg.Database.URL != "postgres://from-env" {
		t.Fatalf("expected env override, got %s", cfg.Database.URL)
	}
	if cfg.App.Name != "Env App" {
		t.Fatalf("expected app env override, got %s", cfg.App.Name)
	}
	if cfg.Business.AppointmentMinLeadTimeMinutes != 45 {
		t.Fatalf("unexpected business parameter: %d", cfg.Business.AppointmentMinLeadTimeMinutes)
	}
}

func TestLoad_AppliesDefaults(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "app.json")

	content := `{
		"database": {"url": "postgres://localhost"},
		"logging": {"file": "./logs/test.log"}
	}`

	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write config file: %v", err)
	}

	cfg, _, err := Load(path)
	if err != nil {
		t.Fatalf("load config: %v", err)
	}

	if cfg.App.Name != "VitaLine Clinic" {
		t.Fatalf("unexpected default app name: %s", cfg.App.Name)
	}
	if cfg.API.Address != ":8080" {
		t.Fatalf("unexpected default api address: %s", cfg.API.Address)
	}
	if cfg.Logging.Level != "info" {
		t.Fatalf("unexpected default log level: %s", cfg.Logging.Level)
	}
	if cfg.Business.AppointmentMinLeadTimeMinutes != 30 {
		t.Fatalf("unexpected default lead time: %d", cfg.Business.AppointmentMinLeadTimeMinutes)
	}
}
