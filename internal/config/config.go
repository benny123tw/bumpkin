package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config represents the bumpkin configuration
type Config struct {
	Prefix string `yaml:"prefix"`
	Remote string `yaml:"remote"`
	Hooks  Hooks  `yaml:"hooks"`
}

// Hooks contains pre and post tag hooks
type Hooks struct {
	PreTag  []string `yaml:"pre-tag"`
	PostTag []string `yaml:"post-tag"`
}

// Default returns a config with default values
func Default() *Config {
	return &Config{
		Prefix: "v",
		Remote: "origin",
		Hooks:  Hooks{},
	}
}

// Load loads configuration from the given directory
// It looks for .bumpkin.yml or .bumpkin.yaml
func Load(dir string) (*Config, error) {
	cfg := Default()

	// Try .bumpkin.yml first
	configPath := filepath.Join(dir, ".bumpkin.yml")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Try .bumpkin.yaml
		configPath = filepath.Join(dir, ".bumpkin.yaml")
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			// No config file, return defaults
			return cfg, nil
		}
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Ensure defaults for unset values
	if cfg.Prefix == "" {
		cfg.Prefix = "v"
	}
	if cfg.Remote == "" {
		cfg.Remote = "origin"
	}

	return cfg, nil
}

// Merge merges another config into this one, with the other config taking precedence
func (c *Config) Merge(other *Config) *Config {
	result := &Config{
		Prefix: c.Prefix,
		Remote: c.Remote,
		Hooks:  c.Hooks,
	}

	if other.Prefix != "" {
		result.Prefix = other.Prefix
	}
	if other.Remote != "" {
		result.Remote = other.Remote
	}
	if len(other.Hooks.PreTag) > 0 {
		result.Hooks.PreTag = other.Hooks.PreTag
	}
	if len(other.Hooks.PostTag) > 0 {
		result.Hooks.PostTag = other.Hooks.PostTag
	}

	return result
}
