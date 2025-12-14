package version

import (
	"fmt"
	"strings"

	"github.com/Masterminds/semver/v3"
)

// Version represents a semantic version
type Version struct {
	Major      uint64
	Minor      uint64
	Patch      uint64
	Prerelease string
	Metadata   string
}

// String returns the version string without prefix
func (v Version) String() string {
	s := fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
	if v.Prerelease != "" {
		s += "-" + v.Prerelease
	}
	if v.Metadata != "" {
		s += "+" + v.Metadata
	}
	return s
}

// StringWithPrefix returns the version string with the given prefix
func (v Version) StringWithPrefix(prefix string) string {
	return prefix + v.String()
}

// Parse parses a version string into a Version struct
// It accepts versions with or without a "v" prefix
func Parse(s string) (Version, error) {
	if s == "" {
		return Version{}, fmt.Errorf("empty version string")
	}

	// Remove "v" prefix if present
	s = strings.TrimPrefix(s, "v")

	sv, err := semver.NewVersion(s)
	if err != nil {
		return Version{}, fmt.Errorf("invalid version %q: %w", s, err)
	}

	return Version{
		Major:      sv.Major(),
		Minor:      sv.Minor(),
		Patch:      sv.Patch(),
		Prerelease: sv.Prerelease(),
		Metadata:   sv.Metadata(),
	}, nil
}

// LessThan returns true if v is less than other
func (v Version) LessThan(other Version) bool {
	sv1 := v.toSemver()
	sv2 := other.toSemver()
	return sv1.LessThan(sv2)
}

// Equal returns true if v equals other
func (v Version) Equal(other Version) bool {
	return v.Major == other.Major &&
		v.Minor == other.Minor &&
		v.Patch == other.Patch &&
		v.Prerelease == other.Prerelease
}

// toSemver converts Version to Masterminds semver for comparison
func (v Version) toSemver() *semver.Version {
	s := v.String()
	sv, _ := semver.NewVersion(s)
	return sv
}

// IsPrerelease returns true if the version has a prerelease identifier
func (v Version) IsPrerelease() bool {
	return v.Prerelease != ""
}

// Zero returns the zero version (0.0.0)
func Zero() Version {
	return Version{Major: 0, Minor: 0, Patch: 0}
}
