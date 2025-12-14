package cli

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCompletionCommand_Help(t *testing.T) {
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"completion", "--help"})

	err := rootCmd.Execute()
	require.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, "completion")
	assert.Contains(t, output, "bash")
	assert.Contains(t, output, "zsh")
	assert.Contains(t, output, "fish")
	assert.Contains(t, output, "powershell")
}

func TestCompletionCommand_BashHelp(t *testing.T) {
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"completion", "bash", "--help"})

	err := rootCmd.Execute()
	require.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, "bash")
	assert.Contains(t, output, "autocompletion")
}

func TestCompletionCommand_ZshHelp(t *testing.T) {
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"completion", "zsh", "--help"})

	err := rootCmd.Execute()
	require.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, "zsh")
}

func TestCompletionCommand_FishHelp(t *testing.T) {
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"completion", "fish", "--help"})

	err := rootCmd.Execute()
	require.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, "fish")
}

func TestCompletionCommand_PowershellHelp(t *testing.T) {
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"completion", "powershell", "--help"})

	err := rootCmd.Execute()
	require.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, "powershell")
}

// Note: Actual completion output goes to os.Stdout directly (Cobra built-in behavior)
// The commands are tested via help flags and execution without error
func TestCompletionCommand_BashExecutes(t *testing.T) {
	rootCmd.SetArgs([]string{"completion", "bash"})
	err := rootCmd.Execute()
	require.NoError(t, err)
}

func TestCompletionCommand_ZshExecutes(t *testing.T) {
	rootCmd.SetArgs([]string{"completion", "zsh"})
	err := rootCmd.Execute()
	require.NoError(t, err)
}

func TestCompletionCommand_FishExecutes(t *testing.T) {
	rootCmd.SetArgs([]string{"completion", "fish"})
	err := rootCmd.Execute()
	require.NoError(t, err)
}

func TestCompletionCommand_PowershellExecutes(t *testing.T) {
	rootCmd.SetArgs([]string{"completion", "powershell"})
	err := rootCmd.Execute()
	require.NoError(t, err)
}
