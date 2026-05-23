package config_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/cronlog/internal/config"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "cronlog.yaml")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("writeTemp: %v", err)
	}
	return p
}

func TestDefault_Values(t *testing.T) {
	cfg := config.Default()
	if cfg.JobName != "cron-job" {
		t.Errorf("expected job_name=cron-job, got %q", cfg.JobName)
	}
	if cfg.WebhookTimeout != 10*time.Second {
		t.Errorf("expected timeout=10s, got %v", cfg.WebhookTimeout)
	}
	if !cfg.NotifyOnFailure {
		t.Error("expected notify_on_failure=true by default")
	}
}

func TestLoad_MissingFile_ReturnsDefaults(t *testing.T) {
	cfg, err := config.Load("/nonexistent/path/cronlog.yaml")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.JobName != "cron-job" {
		t.Errorf("expected default job_name, got %q", cfg.JobName)
	}
}

func TestLoad_ValidFile(t *testing.T) {
	yaml := `
job_name: backup
webhook_url: https://example.com/hook
notify_on_success: true
notify_on_failure: false
`
	p := writeTemp(t, yaml)
	cfg, err := config.Load(p)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.JobName != "backup" {
		t.Errorf("expected job_name=backup, got %q", cfg.JobName)
	}
	if cfg.WebhookURL != "https://example.com/hook" {
		t.Errorf("unexpected webhook_url: %q", cfg.WebhookURL)
	}
	if !cfg.NotifyOnSuccess {
		t.Error("expected notify_on_success=true")
	}
}

func TestLoad_InvalidYAML(t *testing.T) {
	p := writeTemp(t, ":::invalid yaml:::")
	_, err := config.Load(p)
	if err == nil {
		t.Error("expected error for invalid YAML, got nil")
	}
}

func TestValidate_EmptyJobName(t *testing.T) {
	cfg := config.Default()
	cfg.JobName = ""
	if err := cfg.Validate(); err == nil {
		t.Error("expected validation error for empty job_name")
	}
}

func TestValidate_ZeroTimeout(t *testing.T) {
	cfg := config.Default()
	cfg.WebhookTimeout = 0
	if err := cfg.Validate(); err == nil {
		t.Error("expected validation error for zero webhook_timeout")
	}
}
