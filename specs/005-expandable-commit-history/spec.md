# Feature Specification: Expandable Commit History

**Feature Branch**: `005-expandable-commit-history`  
**Created**: 2025-12-16  
**Status**: Implemented  
**Input**: User description: "Since we removed the commit preview view and moved to version selection view. We truncate the commit history and the commit messages. This is not good for users to review the history, do you have any ideas on how to improve this situation?"

## Problem Statement

The current version selection view truncates commit history to 10 commits and shortens commit messages to 50-60 characters. This limits users' ability to:
1. Review all commits that will be included in the release
2. Read full commit messages to understand change context
3. Make informed decisions about version bumps based on complete information

## Proposed Solution

Implement a lazygit-style dual-pane interface where:
1. **Commit History Pane**: A dedicated scrollable section showing all commits with full navigation
2. **Version Selection Pane**: The existing version bump options

Users navigate between panes using Tab/Shift-Tab, similar to lazygit. Both sections are always visible, eliminating the need to toggle between views and allowing users to reference commits while selecting a version.

### Visual Layout

The layout uses a **version-focused ratio** (~30% commits, ~70% version). The primary task is selecting a version; the commits pane provides context when needed but doesn't dominate the screen. Version pane is the default focused pane on entry.

```
┌─ Commits [4/12] ─────────────────────────────────┐
│   abc1234  feat : add new feature                │
│   def5678  fix  : fix login bug                  │
│   ghi9012  feat!: breaking API change            │
│ ▸ jkl3456  docs : update readme                  │
└──────────────────────────────────────────────────┘
┌─ Version ════════════════════════════════════════┐
│ Current: v1.0.0                                  │
│                                                  │
│   > major  2.0.0  (recommended)                  │
│     minor  1.1.0                                 │
│     patch  1.0.1                                 │
│     custom                                       │
│                                                  │
│                                                  │
└──────────────════════════════════════════════════┘
  ↑/↓/j/k: navigate • h/l/tab: switch pane • gg/G: top/bottom • enter: select • q: quit
```

### Keyboard Navigation

| Key | Action |
|-----|--------|
| `↑` / `k` | Move selection up |
| `↓` / `j` | Move selection down |
| `Tab` / `h` / `l` | Switch between panes |
| `gg` | Jump to first commit |
| `G` | Jump to last commit |
| `Enter` | Select version / View commit details |
| `Escape` | Dismiss overlay / Go back |
| `q` | Quit |

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Navigate Between Panes (Priority: P1)

A developer using the version selection TUI wants to review commits before choosing a version bump. They use Tab to switch focus to the commits pane, scroll through commits, then Tab back to select their version.

**Why this priority**: This is the core interaction pattern. Without pane switching, users cannot leverage the dual-pane layout effectively.

**Independent Test**: Can be fully tested by running bumpkin, pressing Tab to switch between panes, and verifying focus indicator moves correctly.

**Acceptance Scenarios**:

1. **Given** user is in version selection pane, **When** user presses Tab (or h/l), **Then** focus moves to commits pane with visual indicator
2. **Given** user is in commits pane, **When** user presses Tab, Shift-Tab, h, or l, **Then** focus moves to version selection pane
3. **Given** user switches panes, **When** focus changes, **Then** the active pane has highlighted border and inactive pane is visually dimmed

---

### User Story 2 - Scroll Through All Commits (Priority: P1)

A developer with 50+ commits since the last release wants to review all of them before deciding on a version bump. They focus on the commits pane and use arrow keys to scroll through the complete list.

**Why this priority**: Viewing all commits is the primary problem this feature solves. Users must be able to see beyond the current 10-commit limit.

**Independent Test**: Can be fully tested by running bumpkin with 50+ commits, focusing on commits pane, and scrolling through entire list.

**Acceptance Scenarios**:

1. **Given** commits pane is focused with 50 commits, **When** user presses down arrow (or j), **Then** selection moves to next commit and viewport scrolls if needed
2. **Given** user is navigating commits, **When** selection changes, **Then** position indicator updates (e.g., "[25/50]") and selected commit is highlighted with `▸`
3. **Given** user is at bottom of list, **When** user presses down, **Then** selection stays at last commit (no wrap-around)
4. **Given** commits pane is focused, **When** user presses `gg`, **Then** selection jumps to first commit
5. **Given** commits pane is focused, **When** user presses `G`, **Then** selection jumps to last commit

---

### User Story 3 - View Full Commit Message (Priority: P2)

A developer sees a truncated commit message in the list and wants to read the full description to understand the change context.

**Why this priority**: While seeing all commits is critical, reading full messages provides additional context for informed decisions.

**Independent Test**: Can be fully tested by selecting a commit with a long message and verifying full message is displayed.

**Acceptance Scenarios**:

1. **Given** commits pane is focused, **When** user presses Enter on a commit, **Then** full commit message is displayed in a detail overlay
2. **Given** detail overlay is shown, **When** user presses Escape or Enter, **Then** overlay closes and returns to commits pane
3. **Given** commit has multi-line body, **When** detail view is shown, **Then** entire body is visible (scrollable if needed)

---

### User Story 4 - Select Version While Viewing Commits (Priority: P2)

A developer has reviewed the commits and wants to select a version bump. They switch to the version pane and make their selection without losing sight of the commit summary.

**Why this priority**: The dual-pane design's value is maintaining context. Users should see commit summary while selecting version.

**Independent Test**: Can be fully tested by switching to version pane and confirming commits pane remains visible.

**Acceptance Scenarios**:

1. **Given** user switches to version pane, **When** version pane is focused, **Then** commits pane remains visible (dimmed but readable)
2. **Given** version pane is focused, **When** user presses up/down arrows, **Then** version selection moves (not commits)
3. **Given** user selects a version, **When** Enter is pressed, **Then** confirmation flow proceeds as normal

---

### Edge Cases

- What happens when there are 0 commits since last tag? Display message "No new commits" in commits pane.
- What happens when there are fewer than 5 commits? Commits pane shows all commits; selection indicator still shows position (e.g., "[2/3]").
- What happens when terminal is too small for dual panes? Both panes shrink proportionally with minimum height of 3 lines for commits pane.
- What happens when commit message is empty? Display commit hash with placeholder text "(no message)".
- How does system handle very long commit messages? Full messages are displayed in the list (no truncation); detail overlay provides full view with metadata.
- What happens if terminal is resized during use? Both panes adapt to new dimensions dynamically.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST display two distinct panes: commits pane (top, ~30% height) and version selection pane (bottom, ~70% height)
- **FR-002**: System MUST allow users to switch focus between panes using Tab, Shift-Tab, h, or l keys
- **FR-003**: System MUST visually indicate which pane is currently focused (highlighted border vs dimmed)
- **FR-004**: System MUST show all commits in the commits pane with vertical scrolling
- **FR-005**: System MUST display selection position indicator showing current selection (e.g., "[5/25]")
- **FR-006**: System MUST allow users to view full commit message via Enter key on selected commit
- **FR-007**: System MUST allow users to dismiss commit detail overlay and return to pane view
- **FR-008**: System MUST preserve scroll position in commits pane when switching between panes
- **FR-009**: System MUST show keyboard shortcut hints at bottom of screen (e.g., "[Tab] switch pane [Enter] confirm")
- **FR-010**: System MUST maintain all existing version selection functionality in the version pane
- **FR-011**: System MUST clearly highlight the currently selected item in the focused pane with a visual indicator (▸)
- **FR-012**: System MUST default focus to version selection pane on initial display (version selection is the primary task)
- **FR-013**: System MUST support vim-style navigation: j/k for up/down, gg for top, G for bottom
- **FR-014**: System MUST display full commit messages without truncation in the commits list

### Key Entities

- **Commits Pane**: Upper section displaying scrollable list of all commits since last tag, showing hash, type badge, and message
- **Version Pane**: Lower section displaying current version and bump options (patch, minor, major, custom, etc.)
- **Commit Detail Overlay**: Modal view showing full commit information including complete message and body
- **Focus Indicator**: Visual distinction (border color, dimming) showing which pane is active

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Users can view 100% of commits since last tag without leaving the version selection workflow
- **SC-002**: Users can switch between panes in a single key press (Tab)
- **SC-003**: Users can read full commit messages (up to 500 characters) via detail overlay
- **SC-004**: Navigation through 100 commits takes no more than 10 seconds of scrolling
- **SC-005**: Both panes remain visible simultaneously during normal operation
- **SC-006**: Users can identify the recommended version bump type while viewing commit history
- **SC-007**: All existing version selection functionality works identically in the new layout

## Assumptions

- Terminal supports standard ANSI escape codes for cursor positioning and colors
- Minimum terminal size of 80x24 characters for dual-pane layout
- Users are familiar with Tab-based pane navigation (common in lazygit, vim splits)
- The existing charmbracelet/bubbletea framework supports the viewport/pane behavior required
- Performance will remain acceptable with up to 500 commits in the scrollable list

## Clarifications

### Session 2025-12-16

- Q: What should the pane height ratio be for the dual-pane layout? → A: Version-focused (~30% commits, ~70% version) - version selection is the primary task; commits pane provides context when needed but doesn't dominate the screen
