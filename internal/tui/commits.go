package tui

import (
	"fmt"
	"strings"

	"github.com/benny123tw/bumpkin/internal/git"
)

// RenderCommitList renders a list of commits
func RenderCommitList(commits []*git.Commit, maxHeight int) string {
	if len(commits) == 0 {
		return WarningStyle.Render("No commits since last tag")
	}

	var sb strings.Builder

	// Limit commits to display
	displayCommits := commits
	truncated := false
	if len(commits) > maxHeight {
		displayCommits = commits[:maxHeight]
		truncated = true
	}

	for _, commit := range displayCommits {
		line := fmt.Sprintf(
			"%s %s",
			CommitHashStyle.Render(commit.ShortHash),
			CommitMessageStyle.Render(truncateString(commit.Subject, 60)),
		)
		sb.WriteString(line)
		sb.WriteString("\n")
	}

	if truncated {
		remaining := len(commits) - maxHeight
		sb.WriteString(
			SubtitleStyle.Render(fmt.Sprintf("  ... and %d more commits", remaining)),
		)
		sb.WriteString("\n")
	}

	return sb.String()
}

// RenderCommitSummary renders a summary of commits
func RenderCommitSummary(commits []*git.Commit) string {
	if len(commits) == 0 {
		return "No commits since last tag"
	}

	return fmt.Sprintf("%d commit(s) since last tag", len(commits))
}

// truncateString truncates a string to maxLen and adds ellipsis if needed
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
