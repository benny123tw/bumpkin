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
	cmd := NewRootCmd(testBuildInfo())
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"help"})

	err := cmd.Execute()
	require.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, "bumpkin")
	assert.Contains(t, output, "Usage:")
}

func TestHelpCommand_MatchesFlag(t *testing.T) {
	// Test that `bumpkin help` produces same output as `bumpkin --help`

	// Run help subcommand
	bufHelp := new(bytes.Buffer)
	cmdHelp := NewRootCmd(testBuildInfo())
	cmdHelp.SetOut(bufHelp)
	cmdHelp.SetArgs([]string{"help"})
	err := cmdHelp.Execute()
	require.NoError(t, err)
	helpOutput := bufHelp.String()

	// Run --help flag
	bufFlag := new(bytes.Buffer)
	cmdFlag := NewRootCmd(testBuildInfo())
	cmdFlag.SetOut(bufFlag)
	cmdFlag.SetArgs([]string{"--help"})
	err = cmdFlag.Execute()
	require.NoError(t, err)
	flagOutput := bufFlag.String()

	// Both should produce identical output
	assert.Equal(t, helpOutput, flagOutput)
}

func TestHelpCommand_ForVersion(t *testing.T) {
	// Test that `bumpkin help version` shows version command help
	buf := new(bytes.Buffer)
	cmd := NewRootCmd(testBuildInfo())
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"help", "version"})

	err := cmd.Execute()
	require.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, "version")
	assert.Contains(t, output, "Print")
}

func TestHelpCommand_UnknownCommand(t *testing.T) {
	// Test that `bumpkin help nonexistent` shows unknown topic message
	buf := new(bytes.Buffer)
	cmd := NewRootCmd(testBuildInfo())
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"help", "nonexistent"})

	err := cmd.Execute()
	// Cobra doesn't return an error, but shows "Unknown help topic"
	require.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, "Unknown help topic")
	assert.Contains(t, output, "Available Commands")
}
