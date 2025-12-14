package tui

import (
	"fmt"
	"strings"

	"github.com/benny123tw/bumpkin/internal/version"
)

// RenderConfirmation renders the confirmation view
func RenderConfirmation(
	prevVersion string,
	newVersion string,
	tagName string,
	commitCount int,
	remoteName string,
	noPush bool,
	dryRun bool,
	selected int,
) string {
	var sb strings.Builder

	// Header
	sb.WriteString(TitleStyle.Render("Confirm Version Bump"))
	sb.WriteString("\n\n")

	// Summary
	sb.WriteString(fmt.Sprintf("  Version: %s %s %s\n",
		CurrentVersionStyle.Render(prevVersion),
		IconArrow,
		NewVersionStyle.Render(newVersion),
	))
	sb.WriteString(fmt.Sprintf("  Tag:     %s\n", NewVersionStyle.Render(tagName)))
	sb.WriteString(fmt.Sprintf("  Commits: %d\n", commitCount))

	if noPush {
		sb.WriteString(fmt.Sprintf("  Push:    %s\n", WarningStyle.Render("disabled (--no-push)")))
	} else {
		sb.WriteString(fmt.Sprintf("  Remote:  %s\n", remoteName))
	}

	if dryRun {
		sb.WriteString("\n")
		sb.WriteString(DryRunStyle.Render(" DRY RUN - No changes will be made "))
		sb.WriteString("\n")
	}

	sb.WriteString("\n")

	// Confirmation buttons
	confirmLabel := "Yes, create tag"
	cancelLabel := "No, cancel"

	if selected == 0 {
		sb.WriteString(fmt.Sprintf("  %s %s\n", IconSelected, SelectedStyle.Render(confirmLabel)))
		sb.WriteString(fmt.Sprintf("    %s\n", UnselectedStyle.Render(cancelLabel)))
	} else {
		sb.WriteString(fmt.Sprintf("    %s\n", UnselectedStyle.Render(confirmLabel)))
		sb.WriteString(fmt.Sprintf("  %s %s\n", IconSelected, SelectedStyle.Render(cancelLabel)))
	}

	return sb.String()
}

// RenderExecuting renders the executing state
func RenderExecuting(step string) string {
	return fmt.Sprintf("%s %s...", SpinnerStyle.Render("⠋"), step)
}

// RenderSuccess renders the success state
func RenderSuccess(result *ExecutionSummary) string {
	var sb strings.Builder

	sb.WriteString(SuccessStyle.Render(fmt.Sprintf("%s Success!", IconCheck)))
	sb.WriteString("\n\n")

	sb.WriteString(fmt.Sprintf("  Created tag: %s\n", NewVersionStyle.Render(result.TagName)))
	sb.WriteString(
		fmt.Sprintf("  Commit:      %s\n", CommitHashStyle.Render(result.CommitHash[:7])),
	)

	if result.Pushed {
		sb.WriteString(fmt.Sprintf("  Pushed to:   %s\n", result.Remote))
	} else {
		sb.WriteString(fmt.Sprintf("  Push:        %s\n", WarningStyle.Render("skipped")))
	}

	sb.WriteString("\n")
	sb.WriteString(HelpStyle.Render("Press any key to exit"))

	return sb.String()
}

// RenderError renders an error state
func RenderError(err error) string {
	var sb strings.Builder

	sb.WriteString(ErrorStyle.Render(fmt.Sprintf("%s Error", IconCross)))
	sb.WriteString("\n\n")
	sb.WriteString(fmt.Sprintf("  %s\n", err.Error()))
	sb.WriteString("\n")
	sb.WriteString(HelpStyle.Render("Press any key to exit"))

	return sb.String()
}

// ExecutionSummary contains summary info for display
type ExecutionSummary struct {
	TagName    string
	CommitHash string
	Pushed     bool
	Remote     string
}

// BumpTypeLabel returns a human-readable label for a bump type
func BumpTypeLabel(bt version.BumpType) string {
	switch bt {
	case version.BumpPatch:
		return "patch"
	case version.BumpMinor:
		return "minor"
	case version.BumpMajor:
		return "major"
	case version.BumpCustom:
		return "custom"
	case version.BumpPrereleaseAlpha:
		return "alpha"
	case version.BumpPrereleaseBeta:
		return "beta"
	case version.BumpPrereleaseRC:
		return "rc"
	case version.BumpRelease:
		return "release"
	default:
		return "unknown"
	}
}
