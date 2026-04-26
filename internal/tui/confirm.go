package tui

import (
	"fmt"
	"strings"
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
	fmt.Fprintf(&sb, "  Version: %s %s %s\n",
		CurrentVersionStyle.Render(prevVersion),
		IconArrow,
		NewVersionStyle.Render(newVersion),
	)
	fmt.Fprintf(&sb, "  Tag:     %s\n", NewVersionStyle.Render(tagName))
	fmt.Fprintf(&sb, "  Commits: %d\n", commitCount)

	if noPush {
		fmt.Fprintf(&sb, "  Push:    %s\n", WarningStyle.Render("disabled (--no-push)"))
	} else {
		fmt.Fprintf(&sb, "  Remote:  %s\n", remoteName)
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
		fmt.Fprintf(&sb, "  %s %s\n", IconSelected, SelectedStyle.Render(confirmLabel))
		fmt.Fprintf(&sb, "    %s\n", UnselectedStyle.Render(cancelLabel))
	} else {
		fmt.Fprintf(&sb, "    %s\n", UnselectedStyle.Render(confirmLabel))
		fmt.Fprintf(&sb, "  %s %s\n", IconSelected, SelectedStyle.Render(cancelLabel))
	}

	return sb.String()
}

// RenderSuccess renders the success state
func RenderSuccess(result *ExecutionSummary) string {
	var sb strings.Builder

	sb.WriteString(SuccessStyle.Render(fmt.Sprintf("%s Success!", IconCheck)))
	sb.WriteString("\n\n")

	fmt.Fprintf(&sb, "  Created tag: %s\n", NewVersionStyle.Render(result.TagName))
	fmt.Fprintf(&sb, "  Commit:      %s\n", CommitHashStyle.Render(result.CommitHash[:7]))

	if result.Pushed {
		fmt.Fprintf(&sb, "  Pushed to:   %s\n", result.Remote)
	} else {
		fmt.Fprintf(&sb, "  Push:        %s\n", WarningStyle.Render("skipped"))
	}

	// Display post-push hook warnings if any
	if len(result.PostPushWarnings) > 0 {
		sb.WriteString("\n")
		sb.WriteString(WarningStyle.Render("  Post-push hook warnings:"))
		sb.WriteString("\n")
		for _, warning := range result.PostPushWarnings {
			fmt.Fprintf(&sb, "    %s %s\n", IconCross, warning)
		}
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
	fmt.Fprintf(&sb, "  %s\n", err.Error())
	sb.WriteString("\n")
	sb.WriteString(HelpStyle.Render("Press any key to exit"))

	return sb.String()
}

// ExecutionSummary contains summary info for display
type ExecutionSummary struct {
	TagName          string
	CommitHash       string
	Pushed           bool
	Remote           string
	PostPushWarnings []string
}
