package cli

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// T052: Test for root command creation
func TestRootCommand(t *testing.T) {
	cmd := NewRootCmd(testBuildInfo())

	assert.Equal(t, "bumpkin", cmd.Use)
	assert.NotEmpty(t, cmd.Short)
	assert.NotEmpty(t, cmd.Long)
}

// T065: Test that flags trigger non-interactive mode
func TestRootCommand_NonInteractiveModeWithFlags(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{"patch flag", []string{"--patch"}},
		{"minor flag", []string{"--minor"}},
		{"major flag", []string{"--major"}},
		{"version flag with value", []string{"--version", "2.0.0"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := NewRootCmd(testBuildInfo())
			cmd.SetArgs(tt.args)

			// The command should recognize these as non-interactive mode
			// We can't fully test execution without a git repo, but we can verify flag parsing
			err := cmd.ParseFlags(tt.args)
			require.NoError(t, err)
		})
	}
}

// T066: Test that no flags triggers interactive mode
func TestRootCommand_InteractiveModeWithoutFlags(t *testing.T) {
	cmd := NewRootCmd(testBuildInfo())
	cmd.SetArgs([]string{})

	// Parse with no args should work
	err := cmd.ParseFlags([]string{})
	require.NoError(t, err)

	// Verify no bump flags are set
	patch, _ := cmd.Flags().GetBool("patch")
	minor, _ := cmd.Flags().GetBool("minor")
	major, _ := cmd.Flags().GetBool("major")

	assert.False(t, patch)
	assert.False(t, minor)
	assert.False(t, major)
}

// T121: Test exit codes
func TestRootCommand_VersionFlag(t *testing.T) {
	cmd := NewRootCmd(testBuildInfo())
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"--show-version"})

	err := cmd.Execute()
	require.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, "bumpkin")
}
