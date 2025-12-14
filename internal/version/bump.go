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
	switch bumpType {
	case BumpPatch:
		return Version{
			Major: v.Major,
			Minor: v.Minor,
			Patch: v.Patch + 1,
		}
	case BumpMinor:
		return Version{
			Major: v.Major,
			Minor: v.Minor + 1,
			Patch: 0,
		}
	case BumpMajor:
		return Version{
			Major: v.Major + 1,
			Minor: 0,
			Patch: 0,
		}
	case BumpRelease:
		return BumpToRelease(v)
	case BumpPrereleaseAlpha:
		return BumpPrerelease(v, "alpha")
	case BumpPrereleaseBeta:
		return BumpPrerelease(v, "beta")
	case BumpPrereleaseRC:
		return BumpPrerelease(v, "rc")
	case BumpCustom:
		// Custom versions are handled separately
		return v
	default:
		return v
	}
}
