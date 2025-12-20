# Implementation Plan: Hook Output Streaming

**Branch**: `006-hook-output-streaming` | **Date**: 2025-12-20 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/006-hook-output-streaming/spec.md`

## Summary

Enable real-time streaming of hook stdout/stderr output during execution. Currently hooks output directly to os.Stdout/os.Stderr which works for non-interactive mode but provides no visibility in the TUI. This feature adds a scrollable output pane to the TUI that displays hook output as it occurs, while maintaining proper line buffering for non-interactive mode.

## Technical Context

**Language/Version**: Go 1.24+
**Primary Dependencies**: charmbracelet/bubbletea v1.3.10, charmbracelet/bubbles v0.21.0, charmbracelet/lipgloss v1.1.0
**Storage**: N/A (in-memory buffer only)
**Testing**: go test with testify/assert, TDD workflow per constitution
**Target Platform**: Darwin/Linux/Windows (cross-platform CLI)
**Project Type**: Single CLI application
**Performance Goals**: 500ms output latency, handle 1,000 lines/sec, 100ms input responsiveness
**Constraints**: <100MB memory for output buffer, no external dependencies beyond existing
**Scale/Scope**: Single-user CLI, buffer up to 10,000 lines

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principle | Status | Notes |
|-----------|--------|-------|
| I. Code Quality | ✅ PASS | All code will pass golangci-lint; no suppressions expected |
| II. Test-Driven Development | ✅ PASS | TDD workflow: Red → Green → Refactor for all components |

**Pre-Phase 0 Gate**: PASSED - No violations, proceed with research.

**Post-Phase 1 Re-check**: PASSED - Design adheres to all principles:
- Code Quality: Design uses existing patterns (viewport, channels, tea.Cmd)
- TDD: Quickstart includes test-first examples for each component

## Project Structure

### Documentation (this feature)

```text
specs/006-hook-output-streaming/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output
├── quickstart.md        # Phase 1 output
├── contracts/           # Phase 1 output (N/A - no external API)
└── tasks.md             # Phase 2 output (/speckit.tasks command)
```

### Source Code (repository root)

```text
internal/
├── hooks/
│   ├── runner.go        # MODIFY: Add streaming support
│   ├── runner_test.go   # MODIFY: Add streaming tests
│   ├── types.go         # MODIFY: Add StreamingResult type
│   ├── buffer.go        # NEW: Line buffer implementation
│   └── buffer_test.go   # NEW: Buffer tests
├── tui/
│   ├── model.go         # MODIFY: Add hook output state
│   ├── hookpane.go      # NEW: Hook output pane component
│   ├── hookpane_test.go # NEW: Hook pane tests
│   ├── messages.go      # MODIFY: Add hook output messages
│   └── styles.go        # MODIFY: Add hook output styles
└── executor/
    └── bump.go          # MODIFY: Wire streaming to TUI
```

**Structure Decision**: Single project layout. New files added to existing `internal/hooks/` and `internal/tui/` packages. No new packages required.

## Complexity Tracking

> No violations to justify - design follows existing patterns.
