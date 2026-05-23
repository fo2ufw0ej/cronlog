package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/example/cronlog/internal/config"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

func TestDefault_Values(t *testing.T) {
	cfg := config.Default()
	if cfg.WebhookURL != "" {
		t.Errorf("WebhookURL = %q, want empty", cfg.WebhookURL)
	}
	if cfg.Rotate.MaxFiles != 0 {
		t.Errorf("Rotate.MaxFiles = %d, want 0", cfg.Rotate.MaxFiles)
	}
}

func TestLoad_MissingFile_ReturnsDefaults(t *testing.T) {
	cfg, err := config.Load("/nonexistent/path.yaml")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.WebhookURL != "" {
		t.Errorf("unexpected WebhookURL: %q", cfg.WebhookURL)
	}
}

func TestLoad_ValidFile(t *testing.T) {
	cfg, err := config.Load(filepath.Join("testdata", "sample.yaml"))
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.WebhookURL == "" {
		t.Error("expected WebhookURL to be set")
	}
}

func TestLoad_InvalidYAML(t *testing.T) {
	path := writeTemp(t, ": bad: yaml: {")
	_, err := config.Load(path)
	if err == nil {
		t.Error("expected error for invalid YAML")
	}
}

func TestLoad_RotateConfig(t *testing.T) {
	cfg, err := config.Load(filepath.Join("testdata", "rotate.yaml"))
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.Rotate.Dir != "/var/log/cronlog" {
		t.Errorf("Rotate.Dir = %q, want /var/log/cronlog", cfg.Rotate.Dir)
	}
	if cfg.Rotate.Prefix != "myjob" {
		t.Errorf("Rotate.Prefix = %q, want myjob", cfg.Rotate.Prefix)
	}
	if cfg.Rotate.MaxFiles != 7 {
		t.Errorf("Rotate.MaxFiles = %d, want 7", cfg.Rotate.MaxFiles)
	}
}

func TestLoad_NotifyConfig(t *testing.T) {
	cfg, err := config.Load(filepath.Join("testdata", "notify.yaml"))
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.Notify.Condition == "" {
		t.Error("expected Notify.Condition to be set")
	}
}
