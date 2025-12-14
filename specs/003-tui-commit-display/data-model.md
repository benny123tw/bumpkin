# Data Model: TUI Commit Display Enhancement

**Feature**: 003-tui-commit-display
**Date**: 2024-12-14

## Overview

This document describes the data model changes for the TUI commit display enhancement. Changes are minimal and primarily extend existing structures.

## Entity Changes

### 1. CommitDisplay (New)

**File**: `internal/tui/commits.go`

```go
// CommitDisplay represents a formatted commit for TUI display
type CommitDisplay struct {
    Hash        string // Short hash (7 chars)
    Type        string // Conventional commit type (feat, fix, docs, etc.)
    Description string // Commit description (without type prefix)
    IsBreaking  bool   // Whether this is a breaking change (has !)
    RawMessage  string // Original full message (fallback)
}
```

**Relationships**: Created from `git.Commit` via `ParseCommitForDisplay()`

---

### 2. CommitTypeStyle (New)

**File**: `internal/tui/styles.go`

```go
// CommitTypeStyles maps commit types to their lipgloss styles
var CommitTypeStyles = map[string]lipgloss.Style{
    "feat":     FeatStyle,
    "fix":      FixStyle,
    "docs":     DocsStyle,
    "chore":    ChoreStyle,
    "refactor": RefactorStyle,
    "test":     TestStyle,
    "style":    ChoreStyle,  // Maps to ChoreStyle (gray)
    "perf":     PerfStyle,
    "ci":       ChoreStyle,  // Maps to ChoreStyle (gray)
    "build":    ChoreStyle,  // Maps to ChoreStyle (gray)
}

// Individual styles
var (
    FeatStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("154")) // Lime
    FixStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("214")) // Yellow
    DocsStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("75"))  // Blue
    ChoreStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("245")) // Gray
    RefactorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("81"))  // Cyan
    TestStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("213")) // Magenta
    PerfStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("208")) // Orange
    
    // Breaking change style (red background, uses errorColor)
    BreakingStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("255")).
        Background(errorColor)
)
```

---

### 3. Model Extension

**File**: `internal/tui/model.go`

No new fields needed - the existing `commits` field already stores `[]*git.Commit`.

The change is in the `View()` method to call `RenderCommitListWithBadges()` instead of `RenderCommitList()`.

---

## New Functions

### ParseCommitForDisplay

**File**: `internal/tui/commits.go`

```go
// ParseCommitForDisplay parses a git commit into a display-friendly format
func ParseCommitForDisplay(commit *git.Commit) CommitDisplay {
    // Parse conventional commit format: type(!)?:description
    // Returns CommitDisplay with extracted type, description, isBreaking
}
```

**Behavior**:
- Extracts type from conventional commit format
- Detects `!` for breaking changes
- Falls back to raw message if not conventional format

---

### RenderCommitListWithBadges

**File**: `internal/tui/commits.go`

```go
// RenderCommitListWithBadges renders commits with colored type badges
func RenderCommitListWithBadges(commits []*git.Commit, maxDisplay int) string {
    // Format each commit as: <hash>  <type-badge> : <description>
    // Truncate to maxDisplay with "and X more commits..." indicator
}
```

**Parameters**:
- `commits`: List of git commits to display
- `maxDisplay`: Maximum commits to show (default: 10)

**Returns**: Formatted string ready for TUI display

---

### GetCommitTypeStyle

**File**: `internal/tui/styles.go`

```go
// GetCommitTypeStyle returns the appropriate style for a commit type
func GetCommitTypeStyle(commitType string, isBreaking bool) lipgloss.Style {
    if isBreaking {
        return BreakingStyle
    }
    if style, ok := CommitTypeStyles[commitType]; ok {
        return style
    }
    return lipgloss.NewStyle() // Default unstyled
}
```

---

## Display Format

### Standard Commit
```
<hash>  <type> : <description>
383d79a  feat : add power function
```

### Breaking Change Commit
```
<hash>  <type> : <description>
383d79a  feat : add power function
         ^^^^
         (red background)
```

### Non-Conventional Commit
```
<hash>  <message>
a1b2c3d  Updated readme
```

---

## View Layout

```
🎃 bumpkin

Current version: v1.0.0

5 Commits since the last version:

383d79a  feat : add power function
d677b41  docs : add package documentation
6bf3686  feat : add divide function
386500e  fix  : handle overflow
...and 1 more commit

? Current version 1.0.0 >
        major 2.0.0
        minor 1.1.0
      > patch 1.0.1
        custom ...
```

---

## Summary of Changes

| Entity | Change Type | Description |
|--------|-------------|-------------|
| `CommitDisplay` | New struct | Parsed commit for display |
| `CommitTypeStyles` | New map | Type to style mapping |
| `*Style` variables | New styles | Individual commit type styles |
| `ParseCommitForDisplay` | New function | Parse git commit to display format |
| `RenderCommitListWithBadges` | New function | Render commits with badges |
| `GetCommitTypeStyle` | New function | Get style for commit type |
| `Model.View()` | Modify | Show commits on version selection |
