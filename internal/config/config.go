package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

// NotifyConfig holds webhook notification settings.
type NotifyConfig struct {
	URL       string `yaml:"url"`
	Condition string `yaml:"condition"`
}

// RotateConfig holds log rotation settings.
type RotateConfig struct {
	Dir      string `yaml:"dir"`
	MaxFiles int    `yaml:"max_files"`
}

// RedactConfig holds log redaction settings.
type RedactConfig struct {
	Patterns []string `yaml:"patterns"`
}

// SummaryConfig controls whether a stats summary line is appended.
type SummaryConfig struct {
	Enabled         bool `yaml:"enabled"`
	AppendToOutput  bool `yaml:"append_to_output"`
}

// Config is the top-level configuration structure.
type Config struct {
	Notify  NotifyConfig  `yaml:"notify"`
	Rotate  RotateConfig  `yaml:"rotate"`
	Redact  RedactConfig  `yaml:"redact"`
	Summary SummaryConfig `yaml:"summary"`
	MaxLines int          `yaml:"max_lines"`
	MaxBytes int          `yaml:"max_bytes"`
}

// Default returns a Config populated with sensible defaults.
func Default() Config {
	return Config{
		Notify: NotifyConfig{
			Condition: "failure",
		},
		Summary: SummaryConfig{
			Enabled:        true,
			AppendToOutput: false,
		},
	}
}

// Load reads a YAML config file and merges it over the defaults.
// If the file does not exist, defaults are returned without error.
func Load(path string) (Config, error) {
	cfg := Default()
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return cfg, nil
	}
	if err != nil {
		return cfg, err
	}
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}
