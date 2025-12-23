package cli

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// testBuildInfo returns a BuildInfo for testing purposes
func testBuildInfo() BuildInfo {
	return BuildInfo{
		Version:   "test",
		Commit:    "abc1234",
		Date:      "2024-01-01",
		GoVersion: "go1.21",
	}
}

func TestCompletionCommand_Help(t *testing.T) {
	buf := new(bytes.Buffer)
	cmd := NewRootCmd(testBuildInfo())
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"completion", "--help"})

	err := cmd.Execute()
	require.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, "completion")
	assert.Contains(t, output, "bash")
	assert.Contains(t, output, "zsh")
	assert.Contains(t, output, "fish")
	assert.Contains(t, output, "powershell")
}

func TestCompletionCommand_Shells(t *testing.T) {
	testCases := []struct {
		shell        string
		helpContains []string
	}{
		{"bash", []string{"bash", "autocompletion"}},
		{"zsh", []string{"zsh"}},
		{"fish", []string{"fish"}},
		{"powershell", []string{"powershell"}},
	}

	for _, tc := range testCases {
		t.Run(tc.shell+" help", func(t *testing.T) {
			buf := new(bytes.Buffer)
			cmd := NewRootCmd(testBuildInfo())
			cmd.SetOut(buf)
			cmd.SetArgs([]string{"completion", tc.shell, "--help"})

			err := cmd.Execute()
			require.NoError(t, err)

			output := buf.String()
			for _, s := range tc.helpContains {
				assert.Contains(t, output, s)
			}
		})

		t.Run(tc.shell+" executes", func(t *testing.T) {
			// Note: Actual completion output goes to os.Stdout directly (Cobra built-in behavior)
			// The commands are tested via execution without error
			cmd := NewRootCmd(testBuildInfo())
			cmd.SetArgs([]string{"completion", tc.shell})
			err := cmd.Execute()
			require.NoError(t, err)
		})
	}
}
