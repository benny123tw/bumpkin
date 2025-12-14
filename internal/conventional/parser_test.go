package conventional

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// T070: Test for parsing "feat:" commit
func TestParseCommit_Feat(t *testing.T) {
	msg := "feat: add new login feature"
	cc, err := ParseCommit(msg)

	require.NoError(t, err)
	assert.Equal(t, "feat", cc.Type)
	assert.Equal(t, "add new login feature", cc.Description)
	assert.False(t, cc.IsBreaking)
	assert.Empty(t, cc.Scope)
}

// T071: Test for parsing "fix:" commit
func TestParseCommit_Fix(t *testing.T) {
	msg := "fix: resolve memory leak"
	cc, err := ParseCommit(msg)

	require.NoError(t, err)
	assert.Equal(t, "fix", cc.Type)
	assert.Equal(t, "resolve memory leak", cc.Description)
	assert.False(t, cc.IsBreaking)
}

// T072: Test for parsing "feat!:" breaking change
func TestParseCommit_BreakingWithBang(t *testing.T) {
	msg := "feat!: remove deprecated API"
	cc, err := ParseCommit(msg)

	require.NoError(t, err)
	assert.Equal(t, "feat", cc.Type)
	assert.True(t, cc.IsBreaking)
	assert.Equal(t, "remove deprecated API", cc.Description)
}

// T073: Test for parsing "BREAKING CHANGE:" footer
func TestParseCommit_BreakingChangeFooter(t *testing.T) {
	msg := `feat: update authentication

BREAKING CHANGE: The auth API now requires a token`
	cc, err := ParseCommit(msg)

	require.NoError(t, err)
	assert.Equal(t, "feat", cc.Type)
	assert.True(t, cc.IsBreaking)
}

// T075: Test for parsing commit with scope
func TestParseCommit_WithScope(t *testing.T) {
	msg := "feat(auth): add OAuth support"
	cc, err := ParseCommit(msg)

	require.NoError(t, err)
	assert.Equal(t, "feat", cc.Type)
	assert.Equal(t, "auth", cc.Scope)
	assert.Equal(t, "add OAuth support", cc.Description)
}

func TestParseCommit_BreakingWithScopeAndBang(t *testing.T) {
	msg := "refactor(api)!: change response format"
	cc, err := ParseCommit(msg)

	require.NoError(t, err)
	assert.Equal(t, "refactor", cc.Type)
	assert.Equal(t, "api", cc.Scope)
	assert.True(t, cc.IsBreaking)
}

func TestParseCommit_NonConventional(t *testing.T) {
	tests := []struct {
		name string
		msg  string
	}{
		{"simple message", "Update readme"},
		{"merge commit", "Merge branch 'main' into feature"},
		{"random message", "WIP: working on something"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cc, err := ParseCommit(tt.msg)
			require.NoError(t, err)
			assert.Equal(t, "other", cc.Type)
		})
	}
}

func TestParseCommit_AllTypes(t *testing.T) {
	types := []string{
		"feat", "fix", "docs", "style", "refactor",
		"perf", "test", "build", "ci", "chore", "revert",
	}

	for _, typ := range types {
		t.Run(typ, func(t *testing.T) {
			msg := typ + ": some description"
			cc, err := ParseCommit(msg)
			require.NoError(t, err)
			assert.Equal(t, typ, cc.Type)
		})
	}
}

func TestParseCommit_BreakingFooterVariants(t *testing.T) {
	tests := []struct {
		name string
		msg  string
	}{
		{"BREAKING CHANGE", "fix: update\n\nBREAKING CHANGE: old API removed"},
		{"BREAKING-CHANGE", "fix: update\n\nBREAKING-CHANGE: old API removed"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cc, err := ParseCommit(tt.msg)
			require.NoError(t, err)
			assert.True(t, cc.IsBreaking)
		})
	}
}

func TestParseCommit_WithBody(t *testing.T) {
	msg := `feat(core): add caching layer

This commit adds a Redis-based caching layer to improve performance.

The cache has a default TTL of 1 hour.`

	cc, err := ParseCommit(msg)
	require.NoError(t, err)
	assert.Equal(t, "feat", cc.Type)
	assert.Equal(t, "core", cc.Scope)
	assert.Equal(t, "add caching layer", cc.Description)
	assert.Contains(t, cc.Body, "Redis-based caching")
}
