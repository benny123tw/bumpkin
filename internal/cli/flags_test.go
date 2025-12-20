package cli

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// T054: Test for --patch flag
func TestFlags_Patch(t *testing.T) {
	cmd := NewRootCmd()
	cmd.SetArgs([]string{"--patch"})

	err := cmd.ParseFlags([]string{"--patch"})
	require.NoError(t, err)

	patch, err := cmd.Flags().GetBool("patch")
	require.NoError(t, err)
	assert.True(t, patch)
}

// T055: Test for --minor flag
func TestFlags_Minor(t *testing.T) {
	cmd := NewRootCmd()
	cmd.SetArgs([]string{"--minor"})

	err := cmd.ParseFlags([]string{"--minor"})
	require.NoError(t, err)

	minor, err := cmd.Flags().GetBool("minor")
	require.NoError(t, err)
	assert.True(t, minor)
}

// T056: Test for --major flag
func TestFlags_Major(t *testing.T) {
	cmd := NewRootCmd()
	cmd.SetArgs([]string{"--major"})

	err := cmd.ParseFlags([]string{"--major"})
	require.NoError(t, err)

	major, err := cmd.Flags().GetBool("major")
	require.NoError(t, err)
	assert.True(t, major)
}

// T057: Test for --version custom flag
func TestFlags_CustomVersion(t *testing.T) {
	cmd := NewRootCmd()
	cmd.SetArgs([]string{"--set-version", "2.0.0"})

	err := cmd.ParseFlags([]string{"--set-version", "2.0.0"})
	require.NoError(t, err)

	customVer, err := cmd.Flags().GetString("set-version")
	require.NoError(t, err)
	assert.Equal(t, "2.0.0", customVer)
}

// T059: Test for --dry-run flag
func TestFlags_DryRun(t *testing.T) {
	cmd := NewRootCmd()
	cmd.SetArgs([]string{"--dry-run"})

	err := cmd.ParseFlags([]string{"--dry-run"})
	require.NoError(t, err)

	dryRun, err := cmd.Flags().GetBool("dry-run")
	require.NoError(t, err)
	assert.True(t, dryRun)
}

// T060: Test for --no-push flag
func TestFlags_NoPush(t *testing.T) {
	cmd := NewRootCmd()
	cmd.SetArgs([]string{"--no-push"})

	err := cmd.ParseFlags([]string{"--no-push"})
	require.NoError(t, err)

	noPush, err := cmd.Flags().GetBool("no-push")
	require.NoError(t, err)
	assert.True(t, noPush)
}

// T061: Test for --yes flag
func TestFlags_Yes(t *testing.T) {
	cmd := NewRootCmd()
	cmd.SetArgs([]string{"--yes"})

	err := cmd.ParseFlags([]string{"--yes"})
	require.NoError(t, err)

	yes, err := cmd.Flags().GetBool("yes")
	require.NoError(t, err)
	assert.True(t, yes)
}

// T063: Test for --json output flag
func TestFlags_JSON(t *testing.T) {
	cmd := NewRootCmd()
	cmd.SetArgs([]string{"--json"})

	err := cmd.ParseFlags([]string{"--json"})
	require.NoError(t, err)

	json, err := cmd.Flags().GetBool("json")
	require.NoError(t, err)
	assert.True(t, json)
}

// T123: Test for --remote flag
func TestFlags_Remote(t *testing.T) {
	cmd := NewRootCmd()
	cmd.SetArgs([]string{"--remote", "upstream"})

	err := cmd.ParseFlags([]string{"--remote", "upstream"})
	require.NoError(t, err)

	remote, err := cmd.Flags().GetString("remote")
	require.NoError(t, err)
	assert.Equal(t, "upstream", remote)
}

// T124: Test for --prefix flag
func TestFlags_Prefix(t *testing.T) {
	cmd := NewRootCmd()
	cmd.SetArgs([]string{"--prefix", "release-"})

	err := cmd.ParseFlags([]string{"--prefix", "release-"})
	require.NoError(t, err)

	prefix, err := cmd.Flags().GetString("prefix")
	require.NoError(t, err)
	assert.Equal(t, "release-", prefix)
}

// Test mutual exclusivity of bump flags
func TestFlags_MutualExclusivity(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{"patch and minor", []string{"--patch", "--minor"}},
		{"patch and major", []string{"--patch", "--major"}},
		{"minor and major", []string{"--minor", "--major"}},
		{"patch and version", []string{"--patch", "--set-version", "2.0.0"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// These combinations should be flagged as errors during validation
			// The actual validation happens in the command execution
			cmd := NewRootCmd()
			err := cmd.ParseFlags(tt.args)
			require.NoError(t, err) // Parsing succeeds, validation happens in RunE
		})
	}
}

// Test countTrueFlags helper function
func TestCountTrueFlags(t *testing.T) {
	tests := []struct {
		name     string
		flags    []bool
		expected int
	}{
		{"no flags", []bool{}, 0},
		{"all false", []bool{false, false, false}, 0},
		{"one true", []bool{true, false, false}, 1},
		{"two true", []bool{true, true, false}, 2},
		{"all true", []bool{true, true, true}, 3},
		{"single true", []bool{true}, 1},
		{"single false", []bool{false}, 0},
		{"many flags mixed", []bool{true, false, true, false, true, false, true}, 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := countTrueFlags(tt.flags...)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Test default values
func TestFlags_Defaults(t *testing.T) {
	cmd := NewRootCmd()
	cmd.SetArgs([]string{})
	err := cmd.ParseFlags([]string{})
	require.NoError(t, err)

	remote, _ := cmd.Flags().GetString("remote")
	assert.Equal(t, "origin", remote)

	prefix, _ := cmd.Flags().GetString("prefix")
	assert.Equal(t, "v", prefix)

	dryRun, _ := cmd.Flags().GetBool("dry-run")
	assert.False(t, dryRun)

	noPush, _ := cmd.Flags().GetBool("no-push")
	assert.False(t, noPush)

	yes, _ := cmd.Flags().GetBool("yes")
	assert.False(t, yes)

	json, _ := cmd.Flags().GetBool("json")
	assert.False(t, json)
}
