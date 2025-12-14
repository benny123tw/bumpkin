package version

import (
	"fmt"
	"strconv"
	"strings"
)

// ParsePrerelease parses a prerelease string (e.g., "alpha.0") into type and number
func ParsePrerelease(prerelease string) (string, int, error) {
	if prerelease == "" {
		return "", 0, fmt.Errorf("empty prerelease string")
	}

	parts := strings.SplitN(prerelease, ".", 2)
	if len(parts) != 2 {
		return "", 0, fmt.Errorf("invalid prerelease format: %q", prerelease)
	}

	preType := parts[0]
	if preType == "" {
		return "", 0, fmt.Errorf("empty prerelease type")
	}

	num, err := strconv.Atoi(parts[1])
	if err != nil {
		return "", 0, fmt.Errorf("invalid prerelease number: %w", err)
	}

	return preType, num, nil
}

// BumpPrerelease bumps to the specified prerelease type
// If already at that type, increments the number
// If changing type (e.g., alpha -> beta), starts at 0
// If not a prerelease, bumps patch and adds prerelease
func BumpPrerelease(v Version, preType string) Version {
	result := Version{
		Major: v.Major,
		Minor: v.Minor,
		Patch: v.Patch,
	}

	// If no existing prerelease, bump patch and start at 0
	if v.Prerelease == "" {
		result.Patch++
		result.Prerelease = fmt.Sprintf("%s.0", preType)
		return result
	}

	// Parse existing prerelease
	existingType, existingNum, err := ParsePrerelease(v.Prerelease)
	if err != nil {
		// Invalid existing prerelease, treat as new
		result.Patch++
		result.Prerelease = fmt.Sprintf("%s.0", preType)
		return result
	}

	// Same type: increment number
	if existingType == preType {
		result.Prerelease = fmt.Sprintf("%s.%d", preType, existingNum+1)
		return result
	}

	// Different type: start at 0
	result.Prerelease = fmt.Sprintf("%s.0", preType)
	return result
}

// BumpToRelease strips the prerelease identifier, making it a release version
func BumpToRelease(v Version) Version {
	return Version{
		Major: v.Major,
		Minor: v.Minor,
		Patch: v.Patch,
		// No Prerelease or Metadata
	}
}

// PrereleaseType returns the type of prerelease (alpha, beta, rc, or empty)
func (v Version) PrereleaseType() string {
	if v.Prerelease == "" {
		return ""
	}

	parts := strings.SplitN(v.Prerelease, ".", 2)
	return parts[0]
}
