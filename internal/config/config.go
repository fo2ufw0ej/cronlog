package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

// NotifyCondition mirrors notify.Condition as a plain string for unmarshalling.
type NotifyCondition string

// RotateConfig holds log rotation settings.
type RotateConfig struct {
	Dir      string `yaml:"dir"`
	Prefix   string `yaml:"prefix"`
	MaxFiles int    `yaml:"max_files"`
}

// Config holds all cronlog configuration.
type Config struct {
	WebhookURL string `yaml:"webhook_url"`
	Notify     struct {
		Condition NotifyCondition `yaml:"condition"`
	} `yaml:"notify"`
	Redact struct {
		Patterns []string `yaml:"patterns"`
	} `yaml:"redact"`
	Filter struct {
		Stream  string `yaml:"stream"`
		Pattern string `yaml:"pattern"`
	} `yaml:"filter"`
	Truncate struct {
		MaxLines int `yaml:"max_lines"`
		MaxBytes int `yaml:"max_bytes"`
	} `yaml:"truncate"`
	Rotate RotateConfig `yaml:"rotate"`
}

// Default returns a Config populated with sensible defaults.
func Default() Config {
	return Config{}
}

// Load reads a YAML config file from path. If the file does not exist the
// default config is returned without error.
func Load(path string) (Config, error) {
	cfg := Default()
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return cfg, err
	}
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}
