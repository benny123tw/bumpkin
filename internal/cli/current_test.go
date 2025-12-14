package cli

import (
	"bytes"
	"testing"

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
