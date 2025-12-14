package version

import "fmt"

// BumpType represents the type of version bump
type BumpType int

const (
	BumpPatch BumpType = iota + 1
	BumpMinor
	BumpMajor
	BumpCustom
	BumpPrereleaseAlpha
	BumpPrereleaseBeta
	BumpPrereleaseRC
	BumpRelease
)

// String returns the string representation of BumpType
func (b BumpType) String() string {
	switch b {
	case BumpPatch:
		return "patch"
	case BumpMinor:
		return "minor"
	case BumpMajor:
		return "major"
	case BumpCustom:
		return "custom"
	case BumpPrereleaseAlpha:
		return "prerelease-alpha"
	case BumpPrereleaseBeta:
		return "prerelease-beta"
	case BumpPrereleaseRC:
		return "prerelease-rc"
	case BumpRelease:
		return "release"
	default:
		return "unknown"
	}
}

// ParseBumpType parses a string into a BumpType
func ParseBumpType(s string) (BumpType, error) {
	switch s {
	case "patch":
		return BumpPatch, nil
	case "minor":
		return BumpMinor, nil
	case "major":
		return BumpMajor, nil
	case "custom":
		return BumpCustom, nil
	case "prerelease-alpha":
		return BumpPrereleaseAlpha, nil
	case "prerelease-beta":
		return BumpPrereleaseBeta, nil
	case "prerelease-rc":
		return BumpPrereleaseRC, nil
	case "release":
		return BumpRelease, nil
	default:
		return 0, fmt.Errorf("unknown bump type: %q", s)
	}
}

// Bump applies the specified bump type to a version and returns the new version
func Bump(v Version, bumpType BumpType) Version {
	// Clear prerelease and metadata for standard bumps
	result := Version{
		Major: v.Major,
		Minor: v.Minor,
		Patch: v.Patch,
	}

	switch bumpType {
	case BumpPatch:
		result.Patch++
	case BumpMinor:
		result.Minor++
		result.Patch = 0
	case BumpMajor:
		result.Major++
		result.Minor = 0
		result.Patch = 0
	case BumpRelease:
		// Just strip prerelease, keep version numbers
		// If already a release version, no change
	case BumpCustom, BumpPrereleaseAlpha, BumpPrereleaseBeta, BumpPrereleaseRC:
		// For custom and prerelease bumps, will be handled in prerelease.go
		// Return the original version for now
		return v
	}

	return result
}
