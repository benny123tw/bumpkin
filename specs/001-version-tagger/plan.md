# Implementation Plan: Version Tagger CLI

**Branch**: `001-version-tagger` | **Date**: 2025-12-14 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/001-version-tagger/spec.md`

## Summary

Build a language-agnostic CLI tool for semantic version bumping and git tagging, inspired by antfu's bumpp. The tool provides both an interactive TUI mode (using Bubbletea) and a non-interactive CLI mode (using Cobra) for automation. It analyzes conventional commits to recommend version bumps and supports configurable hooks for custom actions like updating version files.

## Technical Context

**Language/Version**: Go 1.21+
**Primary Dependencies**:
- `github.com/charmbracelet/bubbletea` v0.25+ - TUI framework
- `github.com/charmbracelet/bubbles` v0.18+ - TUI components (list, input, spinner)
- `github.com/charmbracelet/lipgloss` v0.10+ - Terminal styling
- `github.com/spf13/cobra` v1.7+ - CLI framework
- `github.com/go-git/go-git/v5` v5.10+ - Git operations
- `github.com/Masterminds/semver/v3` v3.2+ - Semantic version parsing
- `github.com/leodido/go-conventionalcommits` - Conventional commit parsing

**Storage**: N/A (operates on git repository, optional config file)
**Testing**: Go standard testing + testify for assertions
**Target Platform**: Cross-platform (Linux, macOS, Windows)
**Project Type**: Single CLI application
**Performance Goals**: Interactive mode renders at 60fps, non-interactive completes in <5s
**Constraints**: Must work offline for local operations, network only for push
**Scale/Scope**: Single-user CLI tool

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principle | Status | Notes |
|-----------|--------|-------|
| I. Code Quality | ✅ PASS | golangci-lint configured, will run on all code |

**Pre-Phase 0 Gate**: ✅ PASSED - No violations

## Project Structure

### Documentation (this feature)

```text
specs/001-version-tagger/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output
├── quickstart.md        # Phase 1 output
├── contracts/           # Phase 1 output (CLI interface contracts)
└── tasks.md             # Phase 2 output (/speckit.tasks command)
```

### Source Code (repository root)

```text
bumpkin/
├── cmd/
│   └── bumpkin/
│       └── main.go          # Minimal entry point
├── internal/
│   ├── cli/
│   │   ├── root.go          # Root Cobra command
│   │   └── flags.go         # Flag definitions
│   ├── git/
│   │   ├── repository.go    # Repository wrapper
│   │   ├── tags.go          # Tag operations
│   │   ├── commits.go       # Commit listing
│   │   └── push.go          # Push operations
│   ├── version/
│   │   ├── semver.go        # Version struct and operations
│   │   ├── bump.go          # Bumping logic
│   │   └── prerelease.go    # Prerelease handling
│   ├── conventional/
│   │   ├── parser.go        # Commit message parsing
│   │   └── analyzer.go      # Bump recommendation
│   ├── tui/
│   │   ├── model.go         # Main Bubbletea model
│   │   ├── messages.go      # Message types
│   │   ├── styles.go        # Lipgloss styles
│   │   ├── commits.go       # Commit list view
│   │   ├── selector.go      # Version selector view
│   │   └── confirm.go       # Confirmation view
│   ├── hooks/
│   │   ├── runner.go        # Hook execution
│   │   └── types.go         # Hook definitions
│   ├── config/
│   │   └── config.go        # Configuration loading
│   └── executor/
│       └── bump.go          # Shared bump execution logic
├── .golangci.yml
├── go.mod
├── go.sum
└── README.md
```

**Structure Decision**: Single CLI application with `internal/` packages for all domain logic. The `cmd/bumpkin/main.go` is minimal (<20 lines), delegating to `internal/cli/root.go`. Both TUI and CLI modes share the `internal/executor` package for actual version bumping.

## Complexity Tracking

No constitution violations requiring justification.
