package cli

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHelpCommand(t *testing.T) {
	// Test that `bumpkin help` works (Cobra built-in)
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"help"})

	err := rootCmd.Execute()
	require.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, "bumpkin")
	assert.Contains(t, output, "Usage:")
}

func TestHelpCommand_MatchesFlag(t *testing.T) {
	// Test that `bumpkin help` produces same output as `bumpkin --help`

	// Run help subcommand
	bufHelp := new(bytes.Buffer)
	rootCmd.SetOut(bufHelp)
	rootCmd.SetArgs([]string{"help"})
	err := rootCmd.Execute()
	require.NoError(t, err)
	helpOutput := bufHelp.String()

	// Run --help flag
	bufFlag := new(bytes.Buffer)
	rootCmd.SetOut(bufFlag)
	rootCmd.SetArgs([]string{"--help"})
	err = rootCmd.Execute()
	require.NoError(t, err)
	flagOutput := bufFlag.String()

	// Both should produce identical output
	assert.Equal(t, helpOutput, flagOutput)
}

func TestHelpCommand_ForVersion(t *testing.T) {
	// Test that `bumpkin help version` shows version command help
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"help", "version"})

	err := rootCmd.Execute()
	require.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, "version")
	assert.Contains(t, output, "Print")
}

func TestHelpCommand_UnknownCommand(t *testing.T) {
	// Test that `bumpkin help nonexistent` shows unknown topic message
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"help", "nonexistent"})

	err := rootCmd.Execute()
	// Cobra doesn't return an error, but shows "Unknown help topic"
	require.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, "Unknown help topic")
	assert.Contains(t, output, "Available Commands")
}
