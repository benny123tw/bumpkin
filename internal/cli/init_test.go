package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// withTempDir creates a temp directory, changes to it, and returns a cleanup function.
func withTempDir(t *testing.T) (tmpDir string, cleanup func()) {
	t.Helper()

	tmpDir, err := os.MkdirTemp("", "bumpkin-init-test")
	require.NoError(t, err)

	originalDir, err := os.Getwd()
	require.NoError(t, err)

	err = os.Chdir(tmpDir)
	require.NoError(t, err)

	cleanup = func() {
		_ = os.Chdir(originalDir)
		os.RemoveAll(tmpDir)
	}

	return tmpDir, cleanup
}

func TestInitCommand(t *testing.T) {
	tmpDir, cleanup := withTempDir(t)
	defer cleanup()

	// Run init command
	buf := new(bytes.Buffer)
	cmd := NewRootCmd(testBuildInfo())
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"init"})

	err := cmd.Execute()
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
	tests := []struct {
		name           string
		existingFile   string
		expectedErrMsg string
	}{
		{
			name:           "yaml extension",
			existingFile:   ".bumpkin.yaml",
			expectedErrMsg: "already exists",
		},
		{
			name:           "yml extension",
			existingFile:   ".bumpkin.yml",
			expectedErrMsg: ".bumpkin.yml already exists",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir, cleanup := withTempDir(t)
			defer cleanup()

			// Create existing config file
			configPath := filepath.Join(tmpDir, tt.existingFile)
			//nolint:gosec // Test file, permissions don't matter
			err := os.WriteFile(configPath, []byte("existing: true"), 0o644)
			require.NoError(t, err)

			// Run init command - should fail
			buf := new(bytes.Buffer)
			cmd := NewRootCmd(testBuildInfo())
			cmd.SetOut(buf)
			cmd.SetArgs([]string{"init"})

			err = cmd.Execute()
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.expectedErrMsg)
		})
	}
}

func TestInitCommand_Help(t *testing.T) {
	buf := new(bytes.Buffer)
	cmd := NewRootCmd(testBuildInfo())
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"init", "--help"})

	err := cmd.Execute()
	require.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, "init")
	assert.Contains(t, output, ".bumpkin.yaml")
}
