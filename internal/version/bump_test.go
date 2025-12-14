package version

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// T011: Test for BumpType enum and string representation
func TestBumpType_String(t *testing.T) {
	tests := []struct {
		bumpType BumpType
		expected string
	}{
		{BumpPatch, "patch"},
		{BumpMinor, "minor"},
		{BumpMajor, "major"},
		{BumpCustom, "custom"},
		{BumpPrereleaseAlpha, "prerelease-alpha"},
		{BumpPrereleaseBeta, "prerelease-beta"},
		{BumpPrereleaseRC, "prerelease-rc"},
		{BumpRelease, "release"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.bumpType.String())
		})
	}
}

func TestParseBumpType(t *testing.T) {
	tests := []struct {
		input       string
		expected    BumpType
		expectError bool
	}{
		{"patch", BumpPatch, false},
		{"minor", BumpMinor, false},
		{"major", BumpMajor, false},
		{"custom", BumpCustom, false},
		{"prerelease-alpha", BumpPrereleaseAlpha, false},
		{"prerelease-beta", BumpPrereleaseBeta, false},
		{"prerelease-rc", BumpPrereleaseRC, false},
		{"release", BumpRelease, false},
		{"invalid", BumpType(0), true},
		{"", BumpType(0), true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			bt, err := ParseBumpType(tt.input)
			if tt.expectError {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, bt)
		})
	}
}

// T017-T019: Test for Bump operations (Patch, Minor, Major)
func TestBump(t *testing.T) {
	tests := []struct {
		name     string
		input    Version
		bumpType BumpType
		expected Version
	}{
		// Patch bumps (T017)
		{
			name:     "patch: simple bump",
			input:    Version{Major: 1, Minor: 2, Patch: 3},
			bumpType: BumpPatch,
			expected: Version{Major: 1, Minor: 2, Patch: 4},
		},
		{
			name:     "patch: from zero",
			input:    Version{Major: 0, Minor: 0, Patch: 0},
			bumpType: BumpPatch,
			expected: Version{Major: 0, Minor: 0, Patch: 1},
		},
		{
			name:     "patch: clears prerelease",
			input:    Version{Major: 1, Minor: 0, Patch: 0, Prerelease: "alpha.0"},
			bumpType: BumpPatch,
			expected: Version{Major: 1, Minor: 0, Patch: 1},
		},
		// Minor bumps (T018)
		{
			name:     "minor: simple bump",
			input:    Version{Major: 1, Minor: 2, Patch: 3},
			bumpType: BumpMinor,
			expected: Version{Major: 1, Minor: 3, Patch: 0},
		},
		{
			name:     "minor: from zero",
			input:    Version{Major: 0, Minor: 0, Patch: 0},
			bumpType: BumpMinor,
			expected: Version{Major: 0, Minor: 1, Patch: 0},
		},
		{
			name:     "minor: clears prerelease",
			input:    Version{Major: 1, Minor: 0, Patch: 0, Prerelease: "beta.1"},
			bumpType: BumpMinor,
			expected: Version{Major: 1, Minor: 1, Patch: 0},
		},
		// Major bumps (T019)
		{
			name:     "major: simple bump",
			input:    Version{Major: 1, Minor: 2, Patch: 3},
			bumpType: BumpMajor,
			expected: Version{Major: 2, Minor: 0, Patch: 0},
		},
		{
			name:     "major: from zero",
			input:    Version{Major: 0, Minor: 0, Patch: 0},
			bumpType: BumpMajor,
			expected: Version{Major: 1, Minor: 0, Patch: 0},
		},
		{
			name:     "major: clears prerelease",
			input:    Version{Major: 1, Minor: 0, Patch: 0, Prerelease: "rc.2"},
			bumpType: BumpMajor,
			expected: Version{Major: 2, Minor: 0, Patch: 0},
		},
		// Release bump
		{
			name:     "release: strips prerelease",
			input:    Version{Major: 1, Minor: 0, Patch: 1, Prerelease: "rc.0"},
			bumpType: BumpRelease,
			expected: Version{Major: 1, Minor: 0, Patch: 1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Bump(tt.input, tt.bumpType)
			assert.Equal(t, tt.expected, result)
		})
	}
}
