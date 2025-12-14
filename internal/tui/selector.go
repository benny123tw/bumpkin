package tui

import (
	"fmt"
	"strings"

	"github.com/benny123tw/bumpkin/internal/version"
)

// VersionOption represents a version bump option
type VersionOption struct {
	Label       string
	Description string
	BumpType    version.BumpType
	NewVersion  string
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
		{
			Label:       "custom",
			Description: "Enter a custom version",
			BumpType:    version.BumpCustom,
			NewVersion:  "...",
		},
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

		line := fmt.Sprintf(
			"%s%-8s %s %s",
			cursor,
			opt.Label,
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
