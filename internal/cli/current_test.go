package cli

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCurrentCommand_Help(t *testing.T) {
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"current", "--help"})

	err := rootCmd.Execute()
	require.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, "current")
	assert.Contains(t, output, "version")
}

func TestCurrentCommand_HasPrefixFlag(t *testing.T) {
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"current", "--help"})

	err := rootCmd.Execute()
	require.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, "--prefix")
	assert.Contains(t, output, "-p")
}

func TestCurrentCommand_InGitRepo(t *testing.T) {
	// This test runs in the bumpkin repo which has tags
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"current"})

	err := rootCmd.Execute()
	require.NoError(t, err)

	output := buf.String()
	// Should output a version tag (v0.1.0 or similar)
	assert.Contains(t, output, "v")
}

func TestCurrentCommand_NotGitRepo(t *testing.T) {
	// Create a temp directory that is not a git repo
	tmpDir, err := os.MkdirTemp("", "bumpkin-current-test")
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

	// Test runCurrent directly to avoid rootCmd state issues
	err = runCurrent(rootCmd, nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not a git repository")
}

func TestCurrentCommand_NoTags(t *testing.T) {
	// Create a temp directory with a git repo but no tags
	tmpDir, err := os.MkdirTemp("", "bumpkin-current-notags-test")
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

	// Initialize a git repo with a commit but no tags
	ctx := context.Background()
	cmd := exec.CommandContext(ctx, "git", "init")
	require.NoError(t, cmd.Run())

	cmd = exec.CommandContext(ctx, "git", "config", "user.email", "test@test.com")
	require.NoError(t, cmd.Run())

	cmd = exec.CommandContext(ctx, "git", "config", "user.name", "Test")
	require.NoError(t, cmd.Run())

	// Create a file and commit it
	require.NoError(t, os.WriteFile("test.txt", []byte("test"), 0o600))

	cmd = exec.CommandContext(ctx, "git", "add", ".")
	require.NoError(t, cmd.Run())

	cmd = exec.CommandContext(ctx, "git", "commit", "-m", "initial")
	require.NoError(t, cmd.Run())

	// Test runCurrent directly - should report no tags found
	buf := new(bytes.Buffer)
	testCmd := &cobra.Command{}
	testCmd.SetOut(buf)

	err = runCurrent(testCmd, nil)
	require.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, "No version tags found")
}
