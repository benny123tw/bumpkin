package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// T085: Test for loading .bumpkin.yml
func TestLoad(t *testing.T) {
	tmpDir := t.TempDir()

	configContent := `
prefix: "v"
remote: "origin"
hooks:
  pre-tag:
    - echo "before tag"
  post-tag:
    - echo "after tag"
`
	configPath := filepath.Join(tmpDir, ".bumpkin.yml")
	//nolint:gosec // test file
	err := os.WriteFile(configPath, []byte(configContent), 0o644)
	require.NoError(t, err)

	cfg, err := Load(tmpDir)
	require.NoError(t, err)

	assert.Equal(t, "v", cfg.Prefix)
	assert.Equal(t, "origin", cfg.Remote)
	assert.Len(t, cfg.Hooks.PreTag, 1)
	assert.Len(t, cfg.Hooks.PostTag, 1)
	assert.Equal(t, "echo \"before tag\"", cfg.Hooks.PreTag[0])
	assert.Equal(t, "echo \"after tag\"", cfg.Hooks.PostTag[0])
}

// T086: Test for default config when file missing
func TestLoad_DefaultWhenMissing(t *testing.T) {
	tmpDir := t.TempDir()

	cfg, err := Load(tmpDir)
	require.NoError(t, err)

	// Should return default values
	assert.Equal(t, "v", cfg.Prefix)
	assert.Equal(t, "origin", cfg.Remote)
	assert.Empty(t, cfg.Hooks.PreTag)
	assert.Empty(t, cfg.Hooks.PostTag)
}

// T088: Test for config with hooks defined
func TestLoad_WithHooks(t *testing.T) {
	tmpDir := t.TempDir()

	configContent := `
hooks:
  pre-tag:
    - npm version ${VERSION} --no-git-tag-version
    - git add package.json
  post-tag:
    - npm publish
    - git push
    - git push --tags
`
	configPath := filepath.Join(tmpDir, ".bumpkin.yml")
	//nolint:gosec // test file
	err := os.WriteFile(configPath, []byte(configContent), 0o644)
	require.NoError(t, err)

	cfg, err := Load(tmpDir)
	require.NoError(t, err)

	assert.Len(t, cfg.Hooks.PreTag, 2)
	assert.Len(t, cfg.Hooks.PostTag, 3)
	assert.Contains(t, cfg.Hooks.PreTag[0], "npm version")
	assert.Contains(t, cfg.Hooks.PostTag[0], "npm publish")
}

func TestLoad_YAMLExtension(t *testing.T) {
	tmpDir := t.TempDir()

	configContent := `
prefix: "release-"
`
	// Test .yaml extension
	configPath := filepath.Join(tmpDir, ".bumpkin.yaml")
	//nolint:gosec // test file
	err := os.WriteFile(configPath, []byte(configContent), 0o644)
	require.NoError(t, err)

	cfg, err := Load(tmpDir)
	require.NoError(t, err)

	assert.Equal(t, "release-", cfg.Prefix)
}

func TestLoad_InvalidYAML(t *testing.T) {
	tmpDir := t.TempDir()

	configContent := `
prefix: [invalid yaml
`
	configPath := filepath.Join(tmpDir, ".bumpkin.yml")
	//nolint:gosec // test file
	err := os.WriteFile(configPath, []byte(configContent), 0o644)
	require.NoError(t, err)

	_, err = Load(tmpDir)
	assert.Error(t, err)
}

func TestConfig_Merge(t *testing.T) {
	base := &Config{
		Prefix: "v",
		Remote: "origin",
	}

	overrides := &Config{
		Prefix: "release-",
	}

	merged := base.Merge(overrides)

	assert.Equal(t, "release-", merged.Prefix)
	assert.Equal(t, "origin", merged.Remote)
}
