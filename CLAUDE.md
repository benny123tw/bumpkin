<!-- SPECTRA:START v1.0.2 -->

# Spectra Instructions

This project uses Spectra for Spec-Driven Development(SDD). Specs live in `openspec/specs/`, change proposals in `openspec/changes/`.

## Use `/spectra-*` skills when:

- A discussion needs structure before coding → `/spectra-discuss`
- User wants to plan, propose, or design a change → `/spectra-propose`
- Tasks are ready to implement → `/spectra-apply`
- There's an in-progress change to continue → `/spectra-ingest`
- User asks about specs or how something works → `/spectra-ask`
- Implementation is done → `/spectra-archive`
- Commit only files related to a specific change → `/spectra-commit`

## Workflow

discuss? → propose → apply ⇄ ingest → archive

- `discuss` is optional — skip if requirements are clear
- Requirements change mid-work? Plan mode → `ingest` → resume `apply`

## Parked Changes

Changes can be parked（暫存）— temporarily moved out of `openspec/changes/`. Parked changes won't appear in `spectra list` but can be found with `spectra list --parked`. To restore: `spectra unpark <name>`. The `/spectra-apply` and `/spectra-ingest` skills handle parked changes automatically.

<!-- SPECTRA:END -->

# bumpkin Development Guidelines

## Overview

Bumpkin is a semantic version tagger CLI for git repositories, built with Go 1.25+ and BubbleTea TUI framework.

## Project Structure

```text
cmd/bumpkin/       # CLI entry point
internal/
  cli/             # Cobra command definitions and flags
  config/          # YAML config loading (.bumpkin.yaml)
  conventional/    # Conventional commit parsing
  executor/        # Tag creation and push execution
  git/             # Git repository operations (go-git)
  hooks/           # Hook execution and output streaming
  tui/             # BubbleTea TUI components
  version/         # Semver parsing and bumping
specs/             # Feature specifications
```

## Commands

```bash
# Development
just build          # Build binary to bin/
just test           # Run all tests
just test-v         # Run tests with verbose output
just test-cov       # Run tests with coverage report
just lint           # Run golangci-lint
just fmt            # Format code
just check          # Run tests + lint
just modernize      # Check for Go 1.21+ modernization opportunities

# Running
just run            # Run interactive mode
just run --help     # Show CLI help
just dry-patch      # Dry run patch bump
just dry-minor      # Dry run minor bump
just dry-conventional  # Dry run with conventional commit analysis
```

## Code Style

- Go 1.25+ with modern idioms (use `min`/`max`, `range over int`, `strings.Cut`)
- Follow golangci-lint rules (see `.golangci.yml`)
- Use BubbleTea patterns for TUI (tea.Model, tea.Cmd, tea.Msg)
- Thread-safe buffers with sync.RWMutex for concurrent access
- Context-based cancellation for hook execution

## Key Libraries

- `github.com/spf13/cobra` - CLI framework
- `github.com/charmbracelet/bubbletea` - TUI framework
- `github.com/charmbracelet/lipgloss` - TUI styling
- `github.com/go-git/go-git/v5` - Git operations
- `gopkg.in/yaml.v3` - Config parsing

## Testing

- Use `testify/assert` and `testify/require` for assertions
- Tests are in `*_test.go` files alongside implementation
- Run specific tests: `go test ./internal/hooks/... -run TestRunHook`

<!-- MANUAL ADDITIONS START -->
<!-- MANUAL ADDITIONS END -->
