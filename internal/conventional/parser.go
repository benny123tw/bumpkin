package conventional

import (
	"regexp"
	"strings"
)

// ConventionalCommit represents a parsed conventional commit
type ConventionalCommit struct {
	Type        string
	Scope       string
	Description string
	Body        string
	IsBreaking  bool
	Footers     map[string]string
}

// Standard commit types per Conventional Commits spec
var standardTypes = map[string]bool{
	"feat":     true,
	"fix":      true,
	"docs":     true,
	"style":    true,
	"refactor": true,
	"perf":     true,
	"test":     true,
	"build":    true,
	"ci":       true,
	"chore":    true,
	"revert":   true,
}

// Regex patterns for parsing conventional commits
var (
	// Pattern: type(scope)!: description
	// Groups: 1=type, 2=scope (optional), 3=! (optional), 4=description
	headerPattern = regexp.MustCompile(
		`^([a-zA-Z]+)(?:\(([^)]+)\))?(!)?\s*:\s*(.+)$`,
	)

	// Pattern for BREAKING CHANGE footer
	breakingFooterPattern = regexp.MustCompile(
		`(?m)^BREAKING[ -]CHANGE\s*:\s*`,
	)
)

// ParseCommit parses a commit message according to the Conventional Commits spec
func ParseCommit(message string) (*ConventionalCommit, error) {
	cc := &ConventionalCommit{
		Type:    "other",
		Footers: make(map[string]string),
	}

	// Split into lines
	lines := strings.Split(message, "\n")
	if len(lines) == 0 {
		return cc, nil
	}

	// Parse the header (first line)
	header := strings.TrimSpace(lines[0])
	matches := headerPattern.FindStringSubmatch(header)

	if matches != nil {
		commitType := strings.ToLower(matches[1])

		// Only accept known types
		if standardTypes[commitType] {
			cc.Type = commitType
			cc.Scope = matches[2]
			cc.IsBreaking = matches[3] == "!"
			cc.Description = strings.TrimSpace(matches[4])
		}
	}

	// Parse body and footers
	if len(lines) > 1 {
		// Join remaining lines
		remaining := strings.Join(lines[1:], "\n")
		remaining = strings.TrimSpace(remaining)

		// Check for breaking change in footer
		if breakingFooterPattern.MatchString(remaining) {
			cc.IsBreaking = true
		}

		// Extract body (content between header and footers)
		cc.Body = extractBody(remaining)
	}

	return cc, nil
}

// extractBody extracts the body portion of the commit message
func extractBody(content string) string {
	// Simple extraction: take content up to any footer-like patterns
	// Footers typically start with a token followed by ": " or " #"

	lines := strings.Split(content, "\n")
	var bodyLines []string
	inBody := true

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Check if this looks like a footer
		if isFooterLine(trimmed) {
			inBody = false
		}

		if inBody && trimmed != "" {
			bodyLines = append(bodyLines, line)
		}
	}

	return strings.TrimSpace(strings.Join(bodyLines, "\n"))
}

// isFooterLine checks if a line looks like a conventional commit footer
func isFooterLine(line string) bool {
	// Footers match pattern: token: value or token #value
	// BREAKING CHANGE is a special case
	if strings.HasPrefix(line, "BREAKING CHANGE:") ||
		strings.HasPrefix(line, "BREAKING-CHANGE:") {
		return true
	}

	// Check for other footer patterns (e.g., "Reviewed-by: name", "Fixes #123")
	footerPattern := regexp.MustCompile(`^[A-Za-z-]+\s*[:#]\s*`)
	return footerPattern.MatchString(line)
}
