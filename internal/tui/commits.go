package tui

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/benny123tw/bumpkin/internal/git"
)

// Commit display constants
const (
	conventionalCommitDescTruncate = 50 // Max length for conventional commit descriptions
	nonConventionalCommitTruncate  = 60 // Max length for non-conventional commit messages
	maxCommitsToDisplay            = 10 // Max commits to show before truncating
	noMessagePlaceholder           = "(no message)"
)

// CommitDisplay represents a formatted commit for TUI display
type CommitDisplay struct {
	Hash        string // Short hash (7 chars)
	Type        string // Conventional commit type (feat, fix, docs, etc.)
	Description string // Commit description (without type prefix)
	IsBreaking  bool   // Whether this is a breaking change (has !)
	RawMessage  string // Original full message (fallback)
}

// conventionalCommitRegex matches: type(!)?(\(scope\))?:description
var conventionalCommitRegex = regexp.MustCompile(`^(\w+)(!)?(?:\([^)]*\))?:\s*(.*)$`)

// ParseCommitForDisplay parses a git commit into display format
func ParseCommitForDisplay(hash, message string) CommitDisplay {
	shortHash := hash
	if len(hash) > 7 {
		shortHash = hash[:7]
	}

	display := CommitDisplay{
		Hash:       shortHash,
		RawMessage: message,
	}

	// Try to parse conventional commit format
	matches := conventionalCommitRegex.FindStringSubmatch(message)
	if len(matches) == 4 {
		display.Type = matches[1]
		display.IsBreaking = matches[2] == "!"
		display.Description = matches[3]
	}

	return display
}

// RenderCommitListWithBadges renders commits with colored type badges
func RenderCommitListWithBadges(commits []*git.Commit, maxDisplay int) string {
	if len(commits) == 0 {
		return WarningStyle.Render("No commits since last tag")
	}

	var sb strings.Builder

	displayCount := len(commits)
	if displayCount > maxDisplay {
		displayCount = maxDisplay
	}

	for i := 0; i < displayCount; i++ {
		commit := commits[i]
		display := ParseCommitForDisplay(commit.Hash, commit.Subject)

		// Hash
		sb.WriteString(CommitHashStyle.Render(display.Hash))
		sb.WriteString("  ")

		if display.Type != "" {
			// Conventional commit with type badge
			style := GetCommitTypeStyle(display.Type, display.IsBreaking)
			sb.WriteString(style.Render(display.Type))
			sb.WriteString(" : ")
			sb.WriteString(truncateString(display.Description, conventionalCommitDescTruncate))
		} else {
			// Non-conventional commit
			sb.WriteString(truncateString(display.RawMessage, nonConventionalCommitTruncate))
		}

		sb.WriteString("\n")
	}

	// Show "and X more commits..." if truncated
	if len(commits) > maxDisplay {
		remaining := len(commits) - maxDisplay
		sb.WriteString(HelpStyle.Render(
			fmt.Sprintf("...and %d more commit(s)", remaining),
		))
		sb.WriteString("\n")
	}

	return sb.String()
}

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

// RenderCommitListForViewport renders all commits without truncation for use in viewport
// selectedIndex indicates which commit should be highlighted (-1 for no selection)
func RenderCommitListForViewport(commits []*git.Commit, selectedIndex int) string {
	if len(commits) == 0 {
		return WarningStyle.Render("No new commits")
	}

	var sb strings.Builder

	for i, commit := range commits {
		subject := commit.Subject
		if strings.TrimSpace(subject) == "" {
			subject = noMessagePlaceholder
		}
		display := ParseCommitForDisplay(commit.Hash, subject)

		// Build the line content
		var line strings.Builder

		// Hash
		line.WriteString(CommitHashStyle.Render(display.Hash))
		line.WriteString("  ")

		if display.Type != "" {
			// Conventional commit with type badge
			style := GetCommitTypeStyle(display.Type, display.IsBreaking)
			line.WriteString(style.Render(display.Type))
			line.WriteString(" : ")
			desc := display.Description
			if strings.TrimSpace(desc) == "" {
				desc = noMessagePlaceholder
			}
			line.WriteString(desc)
		} else {
			// Non-conventional commit
			msg := display.RawMessage
			if strings.TrimSpace(msg) == "" {
				msg = noMessagePlaceholder
			}
			line.WriteString(msg)
		}

		// Apply selection indicator
		if i == selectedIndex {
			sb.WriteString(SelectedItemStyle.Render("▸ " + line.String()))
		} else {
			sb.WriteString("  " + line.String())
		}

		sb.WriteString("\n")
	}

	return strings.TrimSuffix(sb.String(), "\n")
}
