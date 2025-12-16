# Quickstart: Expandable Commit History

**Feature**: 005-expandable-commit-history  
**Date**: 2025-12-16

## Overview

This guide provides a quick reference for implementing the dual-pane TUI with scrollable commit history.

---

## Key Files to Modify

| File | Change |
|------|--------|
| `internal/tui/model.go` | Add pane fields, focus management, update routing |
| `internal/tui/styles.go` | Add focused/unfocused border styles |
| `internal/tui/commits.go` | Add full commit list rendering (no truncation) |

---

## Implementation Steps

### Step 1: Add Viewport Import

```go
import "github.com/charmbracelet/bubbles/viewport"
```

### Step 2: Define Pane Type

```go
type PaneType int

const (
    PaneVersion PaneType = iota
    PaneCommits
)
```

### Step 3: Extend Model

```go
type Model struct {
    // ... existing fields ...
    
    // Pane management
    commitsPane         viewport.Model
    focusedPane         PaneType
    showingDetail       bool
    selectedCommitIndex int
}
```

### Step 4: Initialize Viewport

In `New()`:
```go
m.focusedPane = PaneVersion  // Default focus on version pane
m.commitsPane = viewport.New(0, 0)  // Sized later on WindowSizeMsg
```

### Step 5: Handle Window Resize

In `Update()`:
```go
case tea.WindowSizeMsg:
    m.width = msg.Width
    m.height = msg.Height
    
    // Calculate pane heights (~30% commits, ~70% version)
    headerHeight := 4  // title + dry run indicator
    footerHeight := 2  // help text
    availableHeight := m.height - headerHeight - footerHeight
    
    commitsPaneHeight := availableHeight * 30 / 100
    m.commitsPane.Width = m.width - 4  // account for borders
    m.commitsPane.Height = commitsPaneHeight - 2  // account for borders
    
    return m, nil
```

### Step 6: Populate Viewport Content

In `RepoLoadedMsg` handler:
```go
// After loading commits
content := RenderCommitListForViewport(m.commits)
m.commitsPane.SetContent(content)
```

### Step 7: Handle Tab Navigation

```go
case tea.KeyMsg:
    switch msg.String() {
    case "tab", "shift+tab":
        if !m.showingDetail {
            if m.focusedPane == PaneVersion {
                m.focusedPane = PaneCommits
            } else {
                m.focusedPane = PaneVersion
            }
        }
        return m, nil
    }
```

### Step 8: Route Keys to Focused Pane

```go
// After tab handling
if m.focusedPane == PaneCommits {
    var cmd tea.Cmd
    m.commitsPane, cmd = m.commitsPane.Update(msg)
    return m, cmd
}
// Else: existing version selection key handling
```

### Step 9: Render Dual-Pane Layout

```go
func (m Model) renderVersionSelectView() string {
    // Commits pane
    commitsBorder := UnfocusedBorderStyle
    if m.focusedPane == PaneCommits {
        commitsBorder = FocusedBorderStyle
    }
    commitsBox := commitsBorder.Render(m.commitsPane.View())
    
    // Version pane
    versionBorder := UnfocusedBorderStyle
    if m.focusedPane == PaneVersion {
        versionBorder = FocusedBorderStyle
    }
    versionContent := RenderVersionSelector(m.versionOptions, m.selectedOption)
    versionBox := versionBorder.Render(versionContent)
    
    return lipgloss.JoinVertical(lipgloss.Top, commitsBox, versionBox)
}
```

---

## Testing Approach (TDD)

### Red Phase - Write Failing Tests First

1. **Test pane focus switching**:
   - Given version pane focused, when Tab pressed, then commits pane focused
   - Given commits pane focused, when Tab pressed, then version pane focused

2. **Test key routing**:
   - Given commits pane focused, when down arrow pressed, then viewport scrolls
   - Given version pane focused, when down arrow pressed, then version cursor moves

3. **Test layout rendering**:
   - Given model in version select state, when View() called, then both panes rendered

### Green Phase - Implement to Pass

Write minimal code to make each test pass.

### Refactor Phase

Clean up after tests pass:
- Extract pane rendering to separate functions
- Consolidate style definitions
- Remove duplication

---

## Keyboard Shortcuts

| Key | Commits Pane | Version Pane |
|-----|--------------|--------------|
| `↑` / `k` | Scroll up | Move selection up |
| `↓` / `j` | Scroll down | Move selection down |
| `Tab` | Switch to version pane | Switch to commits pane |
| `Enter` | Show commit detail | Select version |
| `Esc` | Close detail overlay | - |
| `q` | Quit | Quit |

---

## Common Gotchas

1. **Viewport must be sized before use** - Initialize with zero dimensions, resize on `tea.WindowSizeMsg`

2. **Viewport Update returns new model** - Always reassign: `m.commitsPane, cmd = m.commitsPane.Update(msg)`

3. **Focus indicator must be visual** - Users need to see which pane is active (border color change)

4. **Don't route keys when showing overlay** - Check `m.showingDetail` before pane-specific key handling
