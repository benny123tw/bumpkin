# Implementation Plan: Expandable Commit History

**Branch**: `005-expandable-commit-history` | **Date**: 2025-12-16 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/005-expandable-commit-history/spec.md`

## Summary

Implement a lazygit-style dual-pane TUI interface for version selection. The upper pane (~30% height) displays a scrollable commit history with all commits since the last tag. The lower pane (~70% height) contains the version selection options. Users navigate between panes using Tab/Shift-Tab, with the version pane being the default focus. This replaces the current single-view approach that truncates commits to 10 and shortens messages.

## Technical Context

**Language/Version**: Go 1.24+  
**Primary Dependencies**: charmbracelet/bubbletea v1.3.10, charmbracelet/bubbles v0.21.0, charmbracelet/lipgloss v1.1.0  
**Storage**: N/A (git repository read-only)  
**Testing**: go test with testify/assert  
**Target Platform**: Cross-platform CLI (macOS, Linux, Windows)  
**Project Type**: Single CLI application  
**Performance Goals**: Smooth scrolling through 500+ commits, no perceptible lag on navigation  
**Constraints**: Terminal minimum 80x24, responsive to resize events  
**Scale/Scope**: Single user CLI tool, no concurrent access concerns

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

### Principle I: Code Quality
- [ ] All new Go code MUST pass golangci-lint checks
- [ ] No `//nolint` directives without justification comments
- [ ] Code formatted via gofmt, gofumpt, goimports, golines

### Principle II: Test-Driven Development (TDD)
- [ ] Tests MUST be written BEFORE implementation code
- [ ] Tests MUST fail initially (Red phase)
- [ ] Implementation written only to make tests pass (Green phase)
- [ ] Code refactored after tests pass (Refactor phase)
- [ ] Commit history MUST show test commits before implementation commits

**Gate Status**: ✅ PASS - No violations. Implementation will follow TDD cycle.

## Project Structure

### Documentation (this feature)

```text
specs/005-expandable-commit-history/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output
├── quickstart.md        # Phase 1 output
├── contracts/           # Phase 1 output (N/A - no external APIs)
└── tasks.md             # Phase 2 output (/speckit.tasks command)
```

### Source Code (repository root)

```text
internal/
├── tui/
│   ├── model.go         # Main TUI model (modify)
│   ├── pane.go          # NEW: Pane abstraction and focus management
│   ├── commits_pane.go  # NEW: Commits pane component
│   ├── version_pane.go  # NEW: Version selection pane component
│   ├── overlay.go       # NEW: Commit detail overlay
│   ├── commits.go       # Existing commit rendering (refactor)
│   ├── selector.go      # Existing version selector (refactor)
│   ├── styles.go        # Styles (extend for pane borders, focus states)
│   ├── messages.go      # Messages (extend for pane events)
│   └── confirm.go       # Confirmation view (unchanged)
└── git/
    └── commits.go       # Commit retrieval (unchanged)

# Test files follow same structure with _test.go suffix
```

**Structure Decision**: Extend existing `internal/tui/` package with new pane components. No new packages needed.

## Complexity Tracking

No constitution violations to justify.
