package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config holds the cronlog configuration.
type Config struct {
	JobName        string        `yaml:"job_name"`
	WebhookURL     string        `yaml:"webhook_url"`
	WebhookTimeout time.Duration `yaml:"webhook_timeout"`
	LogTimestamps  bool          `yaml:"log_timestamps"`
	NotifyOnSuccess bool         `yaml:"notify_on_success"`
	NotifyOnFailure bool         `yaml:"notify_on_failure"`
}

// Default returns a Config with sensible defaults.
func Default() *Config {
	return &Config{
		JobName:         "cron-job",
		WebhookTimeout:  10 * time.Second,
		LogTimestamps:   true,
		NotifyOnSuccess: false,
		NotifyOnFailure: true,
	}
}

// Load reads a YAML config file from the given path.
// Missing fields fall back to defaults.
func Load(path string) (*Config, error) {
	cfg := Default()

	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return nil, fmt.Errorf("config: open %q: %w", path, err)
	}
	defer f.Close()

	dec := yaml.NewDecoder(f)
	dec.KnownFields(true)
	if err := dec.Decode(cfg); err != nil {
		return nil, fmt.Errorf("config: decode %q: %w", path, err)
	}

	return cfg, nil
}

// Validate returns an error if the configuration is invalid.
func (c *Config) Validate() error {
	if c.JobName == "" {
		return fmt.Errorf("config: job_name must not be empty")
	}
	if c.WebhookTimeout <= 0 {
		return fmt.Errorf("config: webhook_timeout must be positive")
	}
	return nil
}
