package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitCommand(t *testing.T) {
	// Create a temp directory for testing
	tmpDir, err := os.MkdirTemp("", "bumpkin-init-test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Change to temp directory
	originalDir, err := os.Getwd()
	require.NoError(t, err)
	defer func() {
		_ = os.Chdir(originalDir)
	}()

	err = os.Chdir(tmpDir)
	require.NoError(t, err)

	// Run init command
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"init"})

	err = rootCmd.Execute()
	require.NoError(t, err)

	// Verify .bumpkin.yaml was created
	configPath := filepath.Join(tmpDir, ".bumpkin.yaml")
	_, err = os.Stat(configPath)
	require.NoError(t, err, ".bumpkin.yaml should be created")

	// Verify content has expected fields
	content, err := os.ReadFile(configPath)
	require.NoError(t, err)
	assert.Contains(t, string(content), "prefix:")
	assert.Contains(t, string(content), "remote:")
}

func TestInitCommand_ConfigExists(t *testing.T) {
	// Create a temp directory for testing
	tmpDir, err := os.MkdirTemp("", "bumpkin-init-test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Change to temp directory
	originalDir, err := os.Getwd()
	require.NoError(t, err)
	defer func() {
		_ = os.Chdir(originalDir)
	}()

	err = os.Chdir(tmpDir)
	require.NoError(t, err)

	// Create existing config file
	configPath := filepath.Join(tmpDir, ".bumpkin.yaml")
	//nolint:gosec // Test file, permissions don't matter
	err = os.WriteFile(configPath, []byte("existing: true"), 0o644)
	require.NoError(t, err)

	// Run init command - should fail
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"init"})

	err = rootCmd.Execute()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "already exists")
}

func TestInitCommand_Help(t *testing.T) {
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"init", "--help"})

	err := rootCmd.Execute()
	require.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, "init")
	assert.Contains(t, output, ".bumpkin.yaml")
}
