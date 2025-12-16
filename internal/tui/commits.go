package tui

import (
	"regexp"
	"strings"

	"github.com/benny123tw/bumpkin/internal/git"
)

// Commit display constants
const (
	noMessagePlaceholder = "(no message)"
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

// stringOrDefault returns the default value if the string is empty or whitespace
func stringOrDefault(s, def string) string {
	if strings.TrimSpace(s) == "" {
		return def
	}
	return s
}

// RenderCommitListForViewport renders all commits without truncation for use in viewport
// selectedIndex indicates which commit should be highlighted (-1 for no selection)
func RenderCommitListForViewport(commits []*git.Commit, selectedIndex int) string {
	if len(commits) == 0 {
		return WarningStyle.Render("No new commits")
	}

	var sb strings.Builder

	for i, commit := range commits {
		subject := stringOrDefault(commit.Subject, noMessagePlaceholder)
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
			line.WriteString(stringOrDefault(display.Description, noMessagePlaceholder))
		} else {
			// Non-conventional commit
			line.WriteString(stringOrDefault(display.RawMessage, noMessagePlaceholder))
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
