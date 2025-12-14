package cli

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVersionCommand(t *testing.T) {
	// Test that `bumpkin version` subcommand works
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"version"})

	err := rootCmd.Execute()
	require.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, "bumpkin")
	assert.Contains(t, output, "built")
}

func TestVersionCommand_MatchesFlag(t *testing.T) {
	// Test that `bumpkin version` produces same output as `bumpkin --show-version`
	// Both use the shared PrintVersionInfo function, so outputs should be identical

	// Run version subcommand
	bufVersion := new(bytes.Buffer)
	rootCmd.SetOut(bufVersion)
	rootCmd.SetArgs([]string{"version"})
	err := rootCmd.Execute()
	require.NoError(t, err)
	versionOutput := bufVersion.String()

	// Run --show-version flag using a fresh command to avoid flag state issues
	cmdFlag := NewRootCmd()
	bufFlag := new(bytes.Buffer)
	cmdFlag.SetOut(bufFlag)
	cmdFlag.SetArgs([]string{"--show-version"})
	err = cmdFlag.Execute()
	require.NoError(t, err)
	flagOutput := bufFlag.String()

	// Both should produce identical output (they use the same PrintVersionInfo function)
	assert.Equal(t, versionOutput, flagOutput)
}

func TestVersionCommand_Help(t *testing.T) {
	// Test that `bumpkin version --help` shows help text
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"version", "--help"})

	err := rootCmd.Execute()
	require.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, "version")
	assert.Contains(t, output, "Print")
}
