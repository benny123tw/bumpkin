# Implementation Plan: Basic Subcommands

**Branch**: `004-basic-commands` | **Date**: 2024-12-14 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/004-basic-commands/spec.md`

## Summary

Add subcommand-style access to common operations (`version`, `help`, `init`, `current`, `completion`) following Go CLI conventions. Cobra natively supports these patterns, making implementation straightforward. Users will be able to use both `bumpkin version` and `bumpkin --version`.

## Technical Context

**Language/Version**: Go 1.24  
**Primary Dependencies**: Cobra v1.10.2 (already in use for CLI)  
**Storage**: YAML files (`.bumpkin.yml` configuration)  
**Testing**: Go standard testing with testify v1.11.1  
**Target Platform**: darwin, linux, windows (amd64, arm64)
**Project Type**: Single CLI application  
**Performance Goals**: Sub-second command execution (CLI tool)  
**Constraints**: None specific - simple CLI commands  
**Scale/Scope**: Single binary, 5 new subcommands

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principle | Status | Notes |
|-----------|--------|-------|
| I. Code Quality | ✅ PASS | All new code will pass golangci-lint checks |

**Pre-Phase 0 Gate**: PASSED - No violations. Simple feature with no architectural concerns.

## Project Structure

### Documentation (this feature)

```text
specs/004-basic-commands/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output
├── quickstart.md        # Phase 1 output
├── contracts/           # Phase 1 output (N/A - CLI commands, no API)
└── tasks.md             # Phase 2 output (/speckit.tasks command)
```

### Source Code (repository root)

```text
cmd/
└── bumpkin/
    └── main.go          # Entry point

internal/
├── cli/
│   ├── root.go          # Root command (existing)
│   ├── version.go       # NEW: version subcommand
│   ├── help.go          # NEW: help subcommand (if needed)
│   ├── init.go          # NEW: init subcommand
│   ├── current.go       # NEW: current subcommand
│   └── completion.go    # NEW: completion subcommand
├── config/
│   └── config.go        # Config loading (existing)
├── git/                 # Git operations (existing)
├── tui/                 # Interactive UI (existing)
└── version/             # Version logic (existing)
```

**Structure Decision**: Follows existing single-project structure. New subcommands will be added as separate files in `internal/cli/` and registered with the root command in their `init()` functions.

## Complexity Tracking

> No violations - section not applicable.
