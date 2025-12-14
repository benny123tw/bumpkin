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

// Hooks contains pre-tag, post-tag, and post-push hooks
type Hooks struct {
	PreTag   []string `yaml:"pre-tag"`
	PostTag  []string `yaml:"post-tag"`
	PostPush []string `yaml:"post-push"`
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
// It looks for .bumpkin.yaml or .bumpkin.yml
func Load(dir string) (*Config, error) {
	// Try .bumpkin.yaml first (preferred)
	configPath := filepath.Join(dir, ".bumpkin.yaml")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Try .bumpkin.yml as fallback
		configPath = filepath.Join(dir, ".bumpkin.yml")
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			// No config file, return defaults
			return Default(), nil
		}
	}

	return LoadFile(configPath)
}

// LoadFile loads configuration from a specific file path
func LoadFile(path string) (*Config, error) {
	cfg := Default()

	// If path doesn't exist, return defaults
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return cfg, nil
	}

	data, err := os.ReadFile(path)
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
	if len(other.Hooks.PostPush) > 0 {
		result.Hooks.PostPush = other.Hooks.PostPush
	}

	return result
}
