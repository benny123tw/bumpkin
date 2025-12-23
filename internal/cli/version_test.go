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
	cmd := NewRootCmd(testBuildInfo())
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"version"})

	err := cmd.Execute()
	require.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, "bumpkin")
	assert.Contains(t, output, "built")
}

func TestVersionCommand_MatchesFlag(t *testing.T) {
	// Test that `bumpkin version` produces same output as `bumpkin --show-version`
	// Both use the same BuildInfo, so outputs should be identical

	// Run version subcommand
	bufVersion := new(bytes.Buffer)
	cmdVersion := NewRootCmd(testBuildInfo())
	cmdVersion.SetOut(bufVersion)
	cmdVersion.SetArgs([]string{"version"})
	err := cmdVersion.Execute()
	require.NoError(t, err)
	versionOutput := bufVersion.String()

	// Run --show-version flag using a fresh command to avoid flag state issues
	cmdFlag := NewRootCmd(testBuildInfo())
	bufFlag := new(bytes.Buffer)
	cmdFlag.SetOut(bufFlag)
	cmdFlag.SetArgs([]string{"--show-version"})
	err = cmdFlag.Execute()
	require.NoError(t, err)
	flagOutput := bufFlag.String()

	// Both should produce identical output (they use the same BuildInfo)
	assert.Equal(t, versionOutput, flagOutput)
}

func TestVersionCommand_Help(t *testing.T) {
	// Test that `bumpkin version --help` shows help text
	buf := new(bytes.Buffer)
	cmd := NewRootCmd(testBuildInfo())
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"version", "--help"})

	err := cmd.Execute()
	require.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, "version")
	assert.Contains(t, output, "Print")
}
