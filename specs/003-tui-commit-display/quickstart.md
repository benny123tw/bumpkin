# Quickstart: TUI Commit Display Enhancement

**Feature**: 003-tui-commit-display
**Date**: 2024-12-14

## Prerequisites

- Go 1.24+
- golangci-lint v2
- Existing bumpkin codebase checked out

## Development Setup

```bash
# Clone and enter repo
cd /path/to/bumpkin
git checkout 003-tui-commit-display

# Verify dependencies
go mod download

# Run existing tests
just test

# Run linter
just lint
```

## Implementation Steps

### Step 1: Add Commit Type Styles (10 min)

**File**: `internal/tui/styles.go`

```go
// Add after existing style definitions

// Commit type styles
var (
    FeatStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("154")) // Lime/Green
    
    FixStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("214")) // Yellow
    
    DocsStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("75")) // Blue
    
    ChoreStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("245")) // Gray
    
    RefactorStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("81")) // Cyan
    
    TestStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("213")) // Magenta
    
    PerfStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("214")) // Orange
    
    BreakingStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("255")).
        Background(lipgloss.Color("196")) // Red background
)

// CommitTypeStyles maps types to styles
var CommitTypeStyles = map[string]lipgloss.Style{
    "feat":     FeatStyle,
    "fix":      FixStyle,
    "docs":     DocsStyle,
    "chore":    ChoreStyle,
    "refactor": RefactorStyle,
    "test":     TestStyle,
    "style":    ChoreStyle,
    "perf":     PerfStyle,
    "ci":       ChoreStyle,
    "build":    ChoreStyle,
}
```

---

### Step 2: Add CommitDisplay Struct and Parser (15 min)

**File**: `internal/tui/commits.go`

```go
import (
    "regexp"
    "strings"
)

// CommitDisplay represents a formatted commit for TUI display
type CommitDisplay struct {
    Hash        string
    Type        string
    Description string
    IsBreaking  bool
    RawMessage  string
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
```

---

### Step 3: Add RenderCommitListWithBadges (15 min)

**File**: `internal/tui/commits.go`

```go
// RenderCommitListWithBadges renders commits with colored type badges
func RenderCommitListWithBadges(commits []*git.Commit, maxDisplay int) string {
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
            sb.WriteString(display.Description)
        } else {
            // Non-conventional commit
            sb.WriteString(display.RawMessage)
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

// GetCommitTypeStyle returns the style for a commit type
func GetCommitTypeStyle(commitType string, isBreaking bool) lipgloss.Style {
    if isBreaking {
        // Get base style and apply breaking background
        if baseStyle, ok := CommitTypeStyles[commitType]; ok {
            return baseStyle.Background(lipgloss.Color("196"))
        }
        return BreakingStyle
    }
    if style, ok := CommitTypeStyles[commitType]; ok {
        return style
    }
    return lipgloss.NewStyle()
}
```

---

### Step 4: Update Model View (10 min)

**File**: `internal/tui/model.go`

Modify `renderVersionSelectView()` to include commits:

```go
func (m Model) renderVersionSelectView() string {
    var sb strings.Builder

    // Current version
    sb.WriteString(fmt.Sprintf("Current version: %s\n\n",
        CurrentVersionStyle.Render(m.currentVersion.StringWithPrefix(m.config.Prefix)),
    ))

    // Commit history (NEW)
    if len(m.commits) > 0 {
        sb.WriteString(SubtitleStyle.Render(
            fmt.Sprintf("%d Commits since the last version:", len(m.commits)),
        ))
        sb.WriteString("\n\n")
        sb.WriteString(RenderCommitListWithBadges(m.commits, 10))
        sb.WriteString("\n")
    }

    // Version selector
    sb.WriteString(SubtitleStyle.Render("Select version bump:"))
    sb.WriteString("\n")
    sb.WriteString(RenderVersionSelector(m.versionOptions, m.selectedOption))

    return sb.String()
}
```

---

### Step 5: Test Manually

```bash
# Build
just build

# Create test commits
cd /tmp
mkdir test-repo && cd test-repo
git init
git commit --allow-empty -m "feat: initial commit"
git tag v1.0.0
git commit --allow-empty -m "feat: add new feature"
git commit --allow-empty -m "fix: fix bug"
git commit --allow-empty -m "feat!: breaking change"
git commit --allow-empty -m "docs: update readme"

# Run bumpkin
/path/to/bumpkin
```

Expected output:
```
🎃 bumpkin

Current version: v1.0.0

4 Commits since the last version:

abc1234  feat : add new feature
def5678  fix  : fix bug
ghi9012  feat : breaking change  (red background)
jkl3456  docs : update readme

Select version bump:
  > major 2.0.0 (recommended)
    minor 1.1.0
    patch 1.0.1
    custom ...
```

---

## Verification Checklist

```bash
# All tests pass
just test

# Linter passes
just lint

# Manual test with various commits
bumpkin  # Interactive mode shows commits with badges
```

## Files Changed Summary

| File | Lines Changed | Type |
|------|---------------|------|
| `internal/tui/styles.go` | +30 | Add type styles |
| `internal/tui/commits.go` | +60 | Add display functions |
| `internal/tui/model.go` | +10 | Update view rendering |

**Estimated Total**: ~100 lines of code
