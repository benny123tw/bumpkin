package version

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// T006: Test for Version struct creation and string formatting
func TestVersion_String(t *testing.T) {
	tests := []struct {
		name     string
		version  Version
		expected string
	}{
		{
			name:     "simple version",
			version:  Version{Major: 1, Minor: 2, Patch: 3},
			expected: "1.2.3",
		},
		{
			name:     "zero version",
			version:  Version{Major: 0, Minor: 0, Patch: 0},
			expected: "0.0.0",
		},
		{
			name:     "version with prerelease",
			version:  Version{Major: 1, Minor: 0, Patch: 0, Prerelease: "alpha.0"},
			expected: "1.0.0-alpha.0",
		},
		{
			name:     "version with metadata",
			version:  Version{Major: 1, Minor: 0, Patch: 0, Metadata: "build.123"},
			expected: "1.0.0+build.123",
		},
		{
			name: "version with prerelease and metadata",
			version: Version{
				Major:      1,
				Minor:      0,
				Patch:      0,
				Prerelease: "beta.1",
				Metadata:   "sha.abc123",
			},
			expected: "1.0.0-beta.1+sha.abc123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.version.String())
		})
	}
}

func TestVersion_StringWithPrefix(t *testing.T) {
	v := Version{Major: 1, Minor: 2, Patch: 3}
	assert.Equal(t, "v1.2.3", v.StringWithPrefix("v"))
	assert.Equal(t, "ver1.2.3", v.StringWithPrefix("ver"))
	assert.Equal(t, "1.2.3", v.StringWithPrefix(""))
}

// T007: Test for parsing version strings (with/without v prefix)
func TestParse(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    Version
		expectError bool
	}{
		{
			name:     "simple version without prefix",
			input:    "1.2.3",
			expected: Version{Major: 1, Minor: 2, Patch: 3},
		},
		{
			name:     "simple version with v prefix",
			input:    "v1.2.3",
			expected: Version{Major: 1, Minor: 2, Patch: 3},
		},
		{
			name:     "version with prerelease",
			input:    "v1.0.0-alpha.0",
			expected: Version{Major: 1, Minor: 0, Patch: 0, Prerelease: "alpha.0"},
		},
		{
			name:     "version with metadata",
			input:    "1.0.0+build.123",
			expected: Version{Major: 1, Minor: 0, Patch: 0, Metadata: "build.123"},
		},
		{
			name:  "version with prerelease and metadata",
			input: "v2.1.0-rc.1+sha.def456",
			expected: Version{
				Major:      2,
				Minor:      1,
				Patch:      0,
				Prerelease: "rc.1",
				Metadata:   "sha.def456",
			},
		},
		{
			name:     "zero version",
			input:    "0.0.0",
			expected: Version{Major: 0, Minor: 0, Patch: 0},
		},
		{
			name:        "invalid version - negative",
			input:       "-1.0.0",
			expectError: true,
		},
		{
			name:        "invalid version - text",
			input:       "invalid",
			expectError: true,
		},
		{
			name:        "empty string",
			input:       "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v, err := Parse(tt.input)
			if tt.expectError {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.expected.Major, v.Major)
			assert.Equal(t, tt.expected.Minor, v.Minor)
			assert.Equal(t, tt.expected.Patch, v.Patch)
			assert.Equal(t, tt.expected.Prerelease, v.Prerelease)
			assert.Equal(t, tt.expected.Metadata, v.Metadata)
		})
	}
}

// T009: Test for version comparison (LessThan, Equal)
func TestVersion_Comparison(t *testing.T) {
	tests := []struct {
		name     string
		v1       Version
		v2       Version
		lessThan bool
		equal    bool
	}{
		{
			name:     "equal versions",
			v1:       Version{Major: 1, Minor: 2, Patch: 3},
			v2:       Version{Major: 1, Minor: 2, Patch: 3},
			lessThan: false,
			equal:    true,
		},
		{
			name:     "major less than",
			v1:       Version{Major: 1, Minor: 0, Patch: 0},
			v2:       Version{Major: 2, Minor: 0, Patch: 0},
			lessThan: true,
			equal:    false,
		},
		{
			name:     "minor less than",
			v1:       Version{Major: 1, Minor: 1, Patch: 0},
			v2:       Version{Major: 1, Minor: 2, Patch: 0},
			lessThan: true,
			equal:    false,
		},
		{
			name:     "patch less than",
			v1:       Version{Major: 1, Minor: 2, Patch: 3},
			v2:       Version{Major: 1, Minor: 2, Patch: 4},
			lessThan: true,
			equal:    false,
		},
		{
			name:     "greater major overrides minor",
			v1:       Version{Major: 2, Minor: 0, Patch: 0},
			v2:       Version{Major: 1, Minor: 9, Patch: 9},
			lessThan: false,
			equal:    false,
		},
		{
			name:     "prerelease less than release",
			v1:       Version{Major: 1, Minor: 0, Patch: 0, Prerelease: "alpha.0"},
			v2:       Version{Major: 1, Minor: 0, Patch: 0},
			lessThan: true,
			equal:    false,
		},
		{
			name:     "alpha less than beta",
			v1:       Version{Major: 1, Minor: 0, Patch: 0, Prerelease: "alpha.0"},
			v2:       Version{Major: 1, Minor: 0, Patch: 0, Prerelease: "beta.0"},
			lessThan: true,
			equal:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.lessThan, tt.v1.LessThan(tt.v2), "LessThan")
			assert.Equal(t, tt.equal, tt.v1.Equal(tt.v2), "Equal")
		})
	}
}
