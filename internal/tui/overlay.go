package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/benny123tw/bumpkin/internal/git"
)

// OverlayStyle defines the style for the commit detail overlay
var OverlayStyle = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	BorderForeground(primaryColor).
	Padding(1, 2).
	Background(lipgloss.Color("235"))

// RenderCommitDetailOverlay renders a full commit detail overlay
func RenderCommitDetailOverlay(commit *git.Commit, width, height int) string {
	if commit == nil {
		return OverlayStyle.Render("No commit selected")
	}

	var sb strings.Builder

	// Header with commit hash
	sb.WriteString(CommitHashStyle.Render("Commit: " + commit.Hash))
	sb.WriteString("\n\n")

	// Author and date
	if commit.Author != "" {
		sb.WriteString(CommitAuthorStyle.Render("Author: " + commit.Author))
		sb.WriteString("\n")
	}
	if !commit.Timestamp.IsZero() {
		dateStr := commit.Timestamp.Format("2006-01-02 15:04:05")
		sb.WriteString(CommitAuthorStyle.Render("Date:   " + dateStr))
		sb.WriteString("\n")
	}
	sb.WriteString("\n")

	// Parse commit for type badge
	display := ParseCommitForDisplay(commit.Hash, commit.Subject)
	if display.Type != "" {
		style := GetCommitTypeStyle(display.Type, display.IsBreaking)
		sb.WriteString(style.Render(display.Type))
		if display.IsBreaking {
			sb.WriteString(" ")
			sb.WriteString(BreakingStyle.Render("BREAKING"))
		}
		sb.WriteString("\n\n")
	}

	// Subject (first line of message)
	sb.WriteString(CommitMessageStyle.Render(commit.Subject))
	sb.WriteString("\n")

	// Full message (may contain more details than subject)
	if commit.Message != "" && commit.Message != commit.Subject {
		sb.WriteString("\n")
		sb.WriteString(SubtitleStyle.Render(commit.Message))
	}

	// Calculate overlay dimensions
	overlayWidth := width - 10
	if overlayWidth < 40 {
		overlayWidth = 40
	}
	if overlayWidth > 80 {
		overlayWidth = 80
	}

	content := sb.String()

	// Apply overlay style with width
	styled := OverlayStyle.Width(overlayWidth).Render(content)

	// Center the overlay
	centered := lipgloss.Place(
		width,
		height,
		lipgloss.Center,
		lipgloss.Center,
		styled,
	)

	return centered
}

// RenderOverlayHeader renders a simple header for the overlay
func RenderOverlayHeader() string {
	return fmt.Sprintf("%s Commit Details %s", "─", "─")
}
