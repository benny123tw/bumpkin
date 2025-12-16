# Data Model: Expandable Commit History

**Feature**: 005-expandable-commit-history  
**Date**: 2025-12-16

## Overview

This feature extends the existing TUI model with pane management and overlay state. No new persistent data storage is required - all state is in-memory during TUI session.

---

## Entities

### 1. Model (extended)

The main TUI model gains additional fields for dual-pane management.

**New Fields**:

| Field | Type | Description |
|-------|------|-------------|
| `commitsPane` | `viewport.Model` | Scrollable viewport containing rendered commit list |
| `focusedPane` | `PaneType` | Which pane currently has focus (0=Version, 1=Commits) |
| `showingDetail` | `bool` | Whether commit detail overlay is displayed |
| `selectedCommitIndex` | `int` | Index of commit selected for detail view |

**Existing Fields** (unchanged):
- `commits []*git.Commit` - Already loaded in current implementation
- `versionOptions []VersionOption` - Already exists
- `selectedOption int` - Cursor for version selection

### 2. PaneType (new enum)

Enumeration for pane focus management.

| Value | Name | Description |
|-------|------|-------------|
| 0 | `PaneVersion` | Version selection pane (default focus) |
| 1 | `PaneCommits` | Commit history pane |

### 3. CommitDetail (view model)

Structure for rendering commit detail overlay. Derived from existing `git.Commit`.

| Field | Type | Description |
|-------|------|-------------|
| `Hash` | `string` | Full commit SHA (40 chars) |
| `ShortHash` | `string` | Abbreviated hash (7 chars) |
| `Author` | `string` | Commit author name |
| `AuthorEmail` | `string` | Commit author email |
| `Date` | `time.Time` | Commit timestamp |
| `Subject` | `string` | First line of commit message |
| `Body` | `string` | Remaining lines of commit message |
| `Type` | `string` | Conventional commit type (feat, fix, etc.) |
| `Scope` | `string` | Conventional commit scope |
| `IsBreaking` | `bool` | Whether commit is a breaking change |

**Source**: Derived from `*git.Commit` and `conventional.ParsedCommit`

---

## State Transitions

### Pane Focus State Machine

```
┌─────────────────────────────────────────────────┐
│              StateVersionSelect                 │
│  ┌─────────────────────────────────────────┐   │
│  │    focusedPane: PaneVersion (default)   │   │
│  │         ↑                               │   │
│  │   Tab/Shift+Tab                         │   │
│  │         ↓                               │   │
│  │    focusedPane: PaneCommits             │   │
│  └─────────────────────────────────────────┘   │
└─────────────────────────────────────────────────┘
```

### Commit Detail Overlay State

```
focusedPane == PaneCommits
        │
        │ [Enter] on commit
        ↓
showingDetail = true
selectedCommitIndex = <cursor position>
        │
        │ [Esc] or [Enter]
        ↓
showingDetail = false
```

---

## Relationships

```
Model
├── commitsPane: viewport.Model
│   └── Content: string (rendered commit list)
├── commits: []*git.Commit (source data)
│   └── CommitDetail (derived for overlay)
├── versionOptions: []VersionOption
└── focusedPane: PaneType
    ├── PaneVersion → routes keys to version selection
    └── PaneCommits → routes keys to viewport scrolling
```

---

## Validation Rules

1. **focusedPane**: Must be 0 (PaneVersion) or 1 (PaneCommits)
2. **selectedCommitIndex**: Must be in range `[0, len(commits)-1]` when `showingDetail` is true
3. **commitsPane dimensions**: Must be recalculated on `tea.WindowSizeMsg`
4. **Overlay dismissal**: `showingDetail` must be false before pane switching

---

## Initialization

1. `focusedPane` defaults to `PaneVersion` (version selection is primary task)
2. `commitsPane` initialized with zero dimensions; resized on first `tea.WindowSizeMsg`
3. `showingDetail` defaults to `false`
4. Viewport content set after commits are loaded in `RepoLoadedMsg` handler
