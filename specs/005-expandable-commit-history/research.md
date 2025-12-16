# Research: Expandable Commit History

**Feature**: 005-expandable-commit-history  
**Date**: 2025-12-16

## Overview

Research findings for implementing a lazygit-style dual-pane TUI layout using charmbracelet/bubbletea.

---

## Decision 1: Viewport Component for Scrollable Panes

**Decision**: Use `charmbracelet/bubbles/viewport` for the commits pane scrolling.

**Rationale**:
- Already a dependency in the project (bubbles v0.21.0)
- Provides built-in scrolling with `ScrollUp(n)`, `ScrollDown(n)`, `GotoTop()`, `GotoBottom()`
- Exposes `AtTop()`, `AtBottom()`, `ScrollPercent()` for scroll position indicators
- Handles `tea.KeyMsg` navigation automatically via configurable `KeyMap`
- Supports mouse wheel scrolling via `MouseWheelEnabled`

**Alternatives Considered**:
- **Custom scroll implementation**: Rejected - reinvents the wheel, more bugs
- **List component from bubbles**: Rejected - designed for item selection with cursor, not pure scrolling display

**Implementation Notes**:
- Initialize viewport on `tea.WindowSizeMsg` when dimensions are known
- Use `SetContent(string)` to populate commit list as rendered string
- Viewport height = ~30% of terminal height (version-focused layout)

---

## Decision 2: Focus Management Pattern

**Decision**: Manual focus state management with a `focusedPane` integer field.

**Rationale**:
- Simple and explicit - fits the two-pane requirement
- No external dependencies needed
- Matches existing codebase patterns (state machine in model.go)
- Easy to extend if more panes are added later

**Alternatives Considered**:
- **john-marinelli/panes library**: Rejected - adds external dependency for minimal benefit with only 2 panes
- **Focus() / Blur() on components**: Rejected - viewport doesn't have focus methods; this pattern is for textinput/textarea

**Implementation Pattern**:
```go
type Model struct {
    commitsPane   viewport.Model
    versionPane   // existing version selector state
    focusedPane   int  // 0 = version (default), 1 = commits
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        if msg.String() == "tab" || msg.String() == "shift+tab" {
            m.focusedPane = 1 - m.focusedPane  // toggle
            return m, nil
        }
        // Route to focused pane
        if m.focusedPane == 1 {
            m.commitsPane, cmd = m.commitsPane.Update(msg)
        } else {
            // existing version selection logic
        }
    }
}
```

---

## Decision 3: Layout Composition

**Decision**: Use `lipgloss.JoinVertical()` for stacking panes with calculated heights.

**Rationale**:
- Already using lipgloss v1.1.0 in the project
- `JoinVertical(lipgloss.Top, pane1, pane2)` cleanly stacks components
- Dynamic height calculation adapts to terminal resize
- Border styles can indicate focus state

**Height Calculation**:
```go
totalHeight := m.height - headerHeight - footerHeight
commitsPaneHeight := totalHeight * 30 / 100  // ~30%
versionPaneHeight := totalHeight - commitsPaneHeight  // ~70%
```

**Visual Focus Indication**:
- Focused pane: bright/highlighted border color
- Unfocused pane: dimmed/subtle border color
- Using lipgloss conditional styling based on `focusedPane` value

---

## Decision 4: Key Handling Delegation

**Decision**: Selective routing based on focused pane, with global handlers for quit/tab.

**Routing Strategy**:
1. **Global handlers** (always processed by root):
   - `ctrl+c`, `q` → `tea.Quit`
   - `tab`, `shift+tab` → toggle `focusedPane`
   - `enter` on version pane → proceed to confirmation

2. **Commits pane focused**:
   - `up/down`, `j/k` → viewport scrolling
   - `enter` → show commit detail overlay
   - `esc` → return from overlay to pane

3. **Version pane focused** (default):
   - `up/down`, `j/k` → version option navigation
   - `enter` → select version, proceed to confirm

4. **Broadcast to all**:
   - `tea.WindowSizeMsg` → recalculate all dimensions

---

## Decision 5: Commit Detail Overlay

**Decision**: Render overlay as a centered box on top of existing content.

**Rationale**:
- Modal overlays are common UX pattern for detail views
- Lipgloss `Place()` function can center content
- Pressing `esc` or `enter` dismisses and returns to pane navigation
- No need for separate state - just a boolean `showingDetail` flag

**Implementation Approach**:
- Add `showingDetail bool` and `selectedCommitIndex int` to model
- In `View()`, if `showingDetail`, render overlay on top of pane content
- Overlay shows: full commit hash, author, date, complete message, body

---

## Decision 6: Handling Small Terminals

**Decision**: Fall back to single-pane mode when terminal height < 16 lines.

**Rationale**:
- Dual-pane with only 3-4 lines per pane is unusable
- Single-pane with Tab to switch maintains functionality
- Threshold of 16 allows ~5 lines per pane minimum

**Implementation**:
```go
if m.height < 16 {
    // Single pane mode: show only focused pane full-screen
    if m.focusedPane == 1 {
        return m.renderFullScreenCommits()
    }
    return m.renderFullScreenVersion()
}
// Normal dual-pane layout
```

---

## Existing Code Analysis

### Current TUI Structure (`internal/tui/model.go`)

**State Machine**:
- `StateLoading` → `StateVersionSelect` → `StateConfirm` → `StateExecuting` → `StateDone`
- `StateCommitList` exists but is deprecated (skipped)

**Key Components**:
- `spinner.Model` for loading states
- `textinput.Model` for custom version input
- `versionOptions []VersionOption` for version selection
- `selectedOption int` for cursor position

**Commit Rendering** (`commits.go`):
- `RenderCommitListWithBadges(commits, maxDisplay)` renders commits with type badges
- Currently truncates to `maxCommitsToDisplay = 10`
- Uses conventional commit parsing for type badges

**Version Selector** (`selector.go`):
- `RenderVersionSelector(options, selected)` renders version options
- Handles recommended version highlighting

### Required Modifications

1. **model.go**:
   - Add `commitsPane viewport.Model`
   - Add `focusedPane int` (0 = version, 1 = commits)
   - Add `showingDetail bool`, `selectedCommitIndex int`
   - Modify `Update()` for pane focus switching
   - Modify `View()` for dual-pane layout

2. **commits.go**:
   - Add `RenderCommitListForViewport(commits)` - no truncation, full list
   - Keep existing `RenderCommitListWithBadges` for backwards compatibility

3. **styles.go**:
   - Add `FocusedBorderStyle`, `UnfocusedBorderStyle`
   - Add `OverlayStyle` for commit detail modal

4. **New files**:
   - Consider `overlay.go` for commit detail rendering (optional - could be in commits.go)

---

## References

- [Charmbracelet Bubbles - Viewport](https://github.com/charmbracelet/bubbles/tree/main/viewport)
- [Bubbletea Examples - Split Editors](https://github.com/charmbracelet/bubbletea/tree/main/examples/split-editors)
- [Lipgloss Layout Functions](https://github.com/charmbracelet/lipgloss#joining-paragraphs)
- [Building Bubble Tea Programs](https://leg100.github.io/en/posts/building-bubbletea-programs/)
