# Implementation Plan: TUI Commit Display Enhancement

**Branch**: `003-tui-commit-display` | **Date**: 2024-12-14 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/003-tui-commit-display/spec.md`

## Summary

Enhance the TUI to persistently display commit history during version selection, and add visual highlighting for commit types with colored badges. Breaking changes (commits with `!`) get special red highlighting.

## Technical Context

**Language/Version**: Go 1.24+
**Primary Dependencies**: 
- github.com/charmbracelet/bubbletea (TUI framework)
- github.com/charmbracelet/lipgloss (styling)

**Storage**: N/A
**Testing**: go test with github.com/stretchr/testify
**Target Platform**: Cross-platform (Linux, macOS, Windows)
**Project Type**: Single CLI application
**Performance Goals**: TUI remains responsive with up to 100 commits
**Constraints**: Terminal width varies; must handle narrow terminals gracefully
**Scale/Scope**: Single-user CLI tool

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principle | Requirement | Status | Notes |
|-----------|-------------|--------|-------|
| Code Quality | All Go code MUST pass golangci-lint | ✅ PASS | Will run `golangci-lint run` before commit |
| Code Quality | No suppressions without justification | ✅ PASS | No new suppressions needed |
| Code Quality | Formatting via gofmt/gofumpt/goimports/golines | ✅ PASS | Will apply `golangci-lint fmt` |
| Development Workflow | Run lint before commit | ✅ PASS | Using justfile tasks |
| CI/CD Integration | PRs must pass lint checks | ✅ PASS | CI workflow already configured |

**Gate Status**: ✅ PASSED - No violations

## Project Structure

### Documentation (this feature)

```text
specs/003-tui-commit-display/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output
├── quickstart.md        # Phase 1 output
├── contracts/           # Phase 1 output (N/A - TUI only)
└── tasks.md             # Phase 2 output (/speckit.tasks command)
```

### Source Code (repository root)

```text
internal/
├── tui/
│   ├── model.go         # Main TUI model (modify view rendering)
│   ├── styles.go        # Lipgloss styles (add commit type styles)
│   ├── commits.go       # Commit list rendering (enhance with badges)
│   ├── selector.go      # Version selector (integrate commit display)
│   └── messages.go      # Tea messages
├── conventional/
│   ├── parser.go        # Commit parsing (already exists)
│   └── parser_test.go   # Parser tests
└── git/
    └── commits.go       # Commit data (already exists)
```

**Structure Decision**: Single project structure. Changes are primarily in `internal/tui/` package with possible minor enhancements to `internal/conventional/` for commit type parsing.

## Complexity Tracking

> No violations - table not needed

## Implementation Approach

### Phase 1: Commit Type Parsing Enhancement

Enhance the conventional commit parser to extract:
- Commit type (feat, fix, docs, etc.)
- Breaking change indicator (`!`)
- Commit description

### Phase 2: Commit Display Styles

Add lipgloss styles for each commit type:
- feat: Green/Lime
- fix: Yellow
- docs: Blue
- chore/style/ci/build: Gray
- refactor: Cyan
- test: Magenta
- perf: Orange
- Breaking (`!`): Red background

### Phase 3: Persistent Commit Display

Modify the TUI model to:
- Show commits on version selection screen
- Keep commit history visible during all selection states
- Truncate to max 10 commits with "and X more..." indicator

### Phase 4: Commit Formatting

Format each commit as:
```
<hash>  <type-badge> : <description>
```

Example:
```
383d79a  feat : add power function
d677b41  docs : add package documentation
6bf3686  feat : add divide function with zero check
```

With breaking change:
```
383d79a  feat (red background) : add power function
```

## Files to Modify

| File | Changes |
|------|---------|
| `internal/tui/styles.go` | Add commit type badge styles |
| `internal/tui/commits.go` | Enhance commit rendering with badges |
| `internal/tui/model.go` | Show commits on version selection screen |
| `internal/conventional/parser.go` | Extract commit type and breaking flag |

## Dependencies

No new dependencies required - uses existing lipgloss for styling.
