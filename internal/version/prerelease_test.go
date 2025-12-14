package version

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// T104: Test for parsing prerelease version
func TestParsePrerelease(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		wantType   string
		wantNumber int
	}{
		{"alpha.0", "alpha.0", "alpha", 0},
		{"alpha.1", "alpha.1", "alpha", 1},
		{"beta.0", "beta.0", "beta", 0},
		{"beta.5", "beta.5", "beta", 5},
		{"rc.0", "rc.0", "rc", 0},
		{"rc.10", "rc.10", "rc", 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			preType, preNum, err := ParsePrerelease(tt.input)
			require.NoError(t, err)
			assert.Equal(t, tt.wantType, preType)
			assert.Equal(t, tt.wantNumber, preNum)
		})
	}
}

// T105: Test for extracting prerelease type and number
func TestParsePrerelease_Invalid(t *testing.T) {
	tests := []string{
		"invalid",
		"alpha",
		"beta.",
		".0",
		"",
	}

	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			_, _, err := ParsePrerelease(input)
			assert.Error(t, err)
		})
	}
}

// T107: Test v1.0.0 → v1.0.1-alpha.0
func TestBumpPrerelease_NewAlpha(t *testing.T) {
	v := Version{Major: 1, Minor: 0, Patch: 0}
	result := BumpPrerelease(v, "alpha")

	assert.Equal(t, uint64(1), result.Major)
	assert.Equal(t, uint64(0), result.Minor)
	assert.Equal(t, uint64(1), result.Patch)
	assert.Equal(t, "alpha.0", result.Prerelease)
}

// T108: Test v1.0.1-alpha.0 → v1.0.1-alpha.1
func TestBumpPrerelease_IncrementAlpha(t *testing.T) {
	v := Version{Major: 1, Minor: 0, Patch: 1, Prerelease: "alpha.0"}
	result := BumpPrerelease(v, "alpha")

	assert.Equal(t, uint64(1), result.Major)
	assert.Equal(t, uint64(0), result.Minor)
	assert.Equal(t, uint64(1), result.Patch)
	assert.Equal(t, "alpha.1", result.Prerelease)
}

// T109: Test v1.0.1-alpha.1 → v1.0.1-beta.0
func TestBumpPrerelease_AlphaToBeta(t *testing.T) {
	v := Version{Major: 1, Minor: 0, Patch: 1, Prerelease: "alpha.1"}
	result := BumpPrerelease(v, "beta")

	assert.Equal(t, uint64(1), result.Major)
	assert.Equal(t, uint64(0), result.Minor)
	assert.Equal(t, uint64(1), result.Patch)
	assert.Equal(t, "beta.0", result.Prerelease)
}

// T111: Test v1.0.1-rc.0 → v1.0.1 (release)
func TestBumpRelease(t *testing.T) {
	v := Version{Major: 1, Minor: 0, Patch: 1, Prerelease: "rc.0"}
	result := BumpToRelease(v)

	assert.Equal(t, uint64(1), result.Major)
	assert.Equal(t, uint64(0), result.Minor)
	assert.Equal(t, uint64(1), result.Patch)
	assert.Empty(t, result.Prerelease)
}

func TestBumpPrerelease_BetaToRC(t *testing.T) {
	v := Version{Major: 1, Minor: 0, Patch: 1, Prerelease: "beta.2"}
	result := BumpPrerelease(v, "rc")

	assert.Equal(t, uint64(1), result.Major)
	assert.Equal(t, uint64(0), result.Minor)
	assert.Equal(t, uint64(1), result.Patch)
	assert.Equal(t, "rc.0", result.Prerelease)
}

func TestBumpPrerelease_IncrementRC(t *testing.T) {
	v := Version{Major: 1, Minor: 0, Patch: 1, Prerelease: "rc.0"}
	result := BumpPrerelease(v, "rc")

	assert.Equal(t, "rc.1", result.Prerelease)
}

func TestBumpPrerelease_FromZero(t *testing.T) {
	v := Zero()
	result := BumpPrerelease(v, "alpha")

	assert.Equal(t, uint64(0), result.Major)
	assert.Equal(t, uint64(0), result.Minor)
	assert.Equal(t, uint64(1), result.Patch)
	assert.Equal(t, "alpha.0", result.Prerelease)
}

func TestBumpRelease_AlreadyRelease(t *testing.T) {
	v := Version{Major: 1, Minor: 0, Patch: 1}
	result := BumpToRelease(v)

	// No change for already-release version
	assert.Equal(t, v, result)
}
