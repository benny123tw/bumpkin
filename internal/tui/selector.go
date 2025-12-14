package tui

import (
	"fmt"
	"strings"

	"github.com/benny123tw/bumpkin/internal/version"
)

// VersionOption represents a version bump option
type VersionOption struct {
	Label         string
	Description   string
	BumpType      version.BumpType
	NewVersion    string
	IsRecommended bool
}

// CreateVersionOptions creates version options based on current version
func CreateVersionOptions(current version.Version, prefix string) []VersionOption {
	options := []VersionOption{
		{
			Label:       "patch",
			Description: "Bug fixes, backwards compatible",
			BumpType:    version.BumpPatch,
			NewVersion:  version.Bump(current, version.BumpPatch).StringWithPrefix(prefix),
		},
		{
			Label:       "minor",
			Description: "New features, backwards compatible",
			BumpType:    version.BumpMinor,
			NewVersion:  version.Bump(current, version.BumpMinor).StringWithPrefix(prefix),
		},
		{
			Label:       "major",
			Description: "Breaking changes",
			BumpType:    version.BumpMajor,
			NewVersion:  version.Bump(current, version.BumpMajor).StringWithPrefix(prefix),
		},
	}

	// Add prerelease options
	options = append(options, createPrereleaseOptions(current, prefix)...)

	// Add custom option at the end
	options = append(options, VersionOption{
		Label:       "custom",
		Description: "Enter a custom version",
		BumpType:    version.BumpCustom,
		NewVersion:  "...",
	})

	return options
}

// createPrereleaseOptions creates prerelease version options
func createPrereleaseOptions(current version.Version, prefix string) []VersionOption {
	var options []VersionOption

	// If current version is a prerelease, show relevant options
	if current.IsPrerelease() {
		preType := current.PrereleaseType()

		// Show option to increment current prerelease type
		switch preType {
		case "alpha":
			options = append(options, VersionOption{
				Label:       "alpha",
				Description: "Increment alpha version",
				BumpType:    version.BumpPrereleaseAlpha,
				NewVersion: version.Bump(current, version.BumpPrereleaseAlpha).
					StringWithPrefix(prefix),
			})
			options = append(options, VersionOption{
				Label:       "beta",
				Description: "Promote to beta",
				BumpType:    version.BumpPrereleaseBeta,
				NewVersion: version.Bump(current, version.BumpPrereleaseBeta).
					StringWithPrefix(prefix),
			})
		case "beta":
			options = append(options, VersionOption{
				Label:       "beta",
				Description: "Increment beta version",
				BumpType:    version.BumpPrereleaseBeta,
				NewVersion: version.Bump(current, version.BumpPrereleaseBeta).
					StringWithPrefix(prefix),
			})
			options = append(options, VersionOption{
				Label:       "rc",
				Description: "Promote to release candidate",
				BumpType:    version.BumpPrereleaseRC,
				NewVersion: version.Bump(current, version.BumpPrereleaseRC).
					StringWithPrefix(prefix),
			})
		case "rc":
			options = append(options, VersionOption{
				Label:       "rc",
				Description: "Increment release candidate",
				BumpType:    version.BumpPrereleaseRC,
				NewVersion: version.Bump(current, version.BumpPrereleaseRC).
					StringWithPrefix(prefix),
			})
		}

		// Always show release option for prereleases
		options = append(options, VersionOption{
			Label:       "release",
			Description: "Promote to stable release",
			BumpType:    version.BumpRelease,
			NewVersion:  version.Bump(current, version.BumpRelease).StringWithPrefix(prefix),
		})
	} else {
		// For stable releases, show alpha option
		options = append(options, VersionOption{
			Label:       "alpha",
			Description: "Start new alpha prerelease",
			BumpType:    version.BumpPrereleaseAlpha,
			NewVersion:  version.Bump(current, version.BumpPrereleaseAlpha).StringWithPrefix(prefix),
		})
	}

	return options
}

// CreateVersionOptionsWithRecommendation creates options with a recommended bump highlighted
func CreateVersionOptionsWithRecommendation(
	current version.Version,
	prefix string,
	recommended version.BumpType,
) []VersionOption {
	options := CreateVersionOptions(current, prefix)

	// Mark the recommended option
	for i := range options {
		if options[i].BumpType == recommended {
			options[i].IsRecommended = true
		}
	}

	return options
}

// RenderVersionSelector renders the version selector
func RenderVersionSelector(options []VersionOption, selected int) string {
	var sb strings.Builder

	for i, opt := range options {
		cursor := "  "
		style := UnselectedStyle
		if i == selected {
			cursor = IconSelected + " "
			style = SelectedStyle
		}

		// Add recommendation indicator
		label := opt.Label
		if opt.IsRecommended {
			label = opt.Label + " " + RecommendedStyle.Render("(recommended)")
		}

		line := fmt.Sprintf(
			"%s%-24s %s %s",
			cursor,
			label,
			IconArrow,
			NewVersionStyle.Render(opt.NewVersion),
		)

		sb.WriteString(style.Render(line))
		sb.WriteString("\n")

		// Show description for selected item
		if i == selected {
			desc := fmt.Sprintf("   %s", SubtitleStyle.Render(opt.Description))
			sb.WriteString(desc)
			sb.WriteString("\n")
		}
	}

	return sb.String()
}
