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

// T009: Test for config parsing hooks.post-push array
func TestLoad_WithPostPushHooks(t *testing.T) {
	tmpDir := t.TempDir()

	configContent := `
hooks:
  pre-tag:
    - echo "before tag"
  post-tag:
    - echo "after tag"
  post-push:
    - curl -X POST $SLACK_WEBHOOK
    - ./scripts/notify-team.sh
`
	configPath := filepath.Join(tmpDir, ".bumpkin.yml")
	//nolint:gosec // test file
	err := os.WriteFile(configPath, []byte(configContent), 0o644)
	require.NoError(t, err)

	cfg, err := Load(tmpDir)
	require.NoError(t, err)

	assert.Len(t, cfg.Hooks.PostPush, 2)
	assert.Equal(t, "curl -X POST $SLACK_WEBHOOK", cfg.Hooks.PostPush[0])
	assert.Equal(t, "./scripts/notify-team.sh", cfg.Hooks.PostPush[1])
}

// T010: Test that empty post-push returns empty slice
func TestLoad_EmptyPostPush(t *testing.T) {
	tmpDir := t.TempDir()

	configContent := `
hooks:
  pre-tag:
    - echo "before"
`
	configPath := filepath.Join(tmpDir, ".bumpkin.yml")
	//nolint:gosec // test file
	err := os.WriteFile(configPath, []byte(configContent), 0o644)
	require.NoError(t, err)

	cfg, err := Load(tmpDir)
	require.NoError(t, err)

	assert.Empty(t, cfg.Hooks.PostPush)
}

// T011: Test that post-push hooks preserve order
func TestLoad_PostPushPreservesOrder(t *testing.T) {
	tmpDir := t.TempDir()

	configContent := `
hooks:
  post-push:
    - echo "first"
    - echo "second"
    - echo "third"
    - echo "fourth"
`
	configPath := filepath.Join(tmpDir, ".bumpkin.yml")
	//nolint:gosec // test file
	err := os.WriteFile(configPath, []byte(configContent), 0o644)
	require.NoError(t, err)

	cfg, err := Load(tmpDir)
	require.NoError(t, err)

	require.Len(t, cfg.Hooks.PostPush, 4)
	assert.Equal(t, "echo \"first\"", cfg.Hooks.PostPush[0])
	assert.Equal(t, "echo \"second\"", cfg.Hooks.PostPush[1])
	assert.Equal(t, "echo \"third\"", cfg.Hooks.PostPush[2])
	assert.Equal(t, "echo \"fourth\"", cfg.Hooks.PostPush[3])
}

func TestConfig_MergeWithPostPush(t *testing.T) {
	base := &Config{
		Prefix: "v",
		Remote: "origin",
		Hooks: Hooks{
			PostPush: []string{"echo base"},
		},
	}

	overrides := &Config{
		Hooks: Hooks{
			PostPush: []string{"echo override"},
		},
	}

	merged := base.Merge(overrides)

	assert.Len(t, merged.Hooks.PostPush, 1)
	assert.Equal(t, "echo override", merged.Hooks.PostPush[0])
}
