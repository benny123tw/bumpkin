package cli

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/spf13/cobra"
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

	// Run version subcommand
	bufVersion := new(bytes.Buffer)
	rootCmd.SetOut(bufVersion)
	rootCmd.SetArgs([]string{"version"})
	err := rootCmd.Execute()
	require.NoError(t, err)
	versionOutput := bufVersion.String()

	// Run --show-version flag using a fresh command instance
	cmdFlag := NewRootCmd()
	// Add the version subcommand to the new instance for consistency
	cmdFlag.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Run: func(cmd *cobra.Command, _ []string) {
			commit := GitCommit
			if len(commit) > 7 {
				commit = commit[:7]
			}
			fmt.Fprintf(cmd.OutOrStdout(), "bumpkin %s (%s, built %s)\n", AppVersion, commit, BuildDate)
		},
	})
	bufFlag := new(bytes.Buffer)
	cmdFlag.SetOut(bufFlag)
	cmdFlag.SetArgs([]string{"--show-version"})
	err = cmdFlag.Execute()
	require.NoError(t, err)
	flagOutput := bufFlag.String()

	// Both should produce identical output
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
