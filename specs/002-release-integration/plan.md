# Implementation Plan: Release Tool Integration (Post-Push Hooks)

**Branch**: `002-release-integration` | **Date**: 2024-12-14 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/002-release-integration/spec.md`

## Summary

Add a `post-push` hook phase to bumpkin that executes after successful tag push. This enables CI/CD-first workflows where bumpkin handles versioning/tagging locally, and post-push hooks trigger notifications or external systems once the tag is pushed and CI/CD takes over for releases.

**Development Approach**: Test-Driven Development (TDD)
- Write tests first
- Implement code to pass tests
- Refactor for quality

## Technical Context

**Language/Version**: Go 1.24+
**Primary Dependencies**: 
- github.com/spf13/cobra (CLI)
- github.com/charmbracelet/bubbletea (TUI)
- github.com/go-git/go-git/v5 (git operations)
- gopkg.in/yaml.v3 (config parsing)

**Storage**: N/A (file-based config: `.bumpkin.yaml`)
**Testing**: go test with github.com/stretchr/testify
**Target Platform**: Cross-platform (Linux, macOS, Windows)
**Project Type**: Single CLI application
**Performance Goals**: Post-push hooks execute within 1 second of push completion
**Constraints**: 
- Hook timeout: 30 seconds default
- Fail-open behavior for post-push hooks (warnings, not fatal)

**Scale/Scope**: Single-user CLI tool

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principle | Requirement | Status | Notes |
|-----------|-------------|--------|-------|
| Code Quality | All Go code MUST pass golangci-lint | ✅ PASS | Will run `golangci-lint run` before commit |
| Code Quality | No suppressions without justification | ✅ PASS | Only existing justified `//nolint` directives for user-defined hooks |
| Code Quality | Formatting via gofmt/gofumpt/goimports/golines | ✅ PASS | Will apply `golangci-lint fmt` |
| Development Workflow | Run lint before commit | ✅ PASS | Using justfile tasks |
| CI/CD Integration | PRs must pass lint checks | ✅ PASS | CI workflow already configured |

**Gate Status**: ✅ PASSED - No violations

## Project Structure

### Documentation (this feature)

```text
specs/002-release-integration/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output
├── quickstart.md        # Phase 1 output
├── contracts/           # Phase 1 output (N/A - CLI, not API)
└── tasks.md             # Phase 2 output (/speckit.tasks command)
```

### Source Code (repository root)

```text
cmd/bumpkin/
└── main.go              # Entry point

internal/
├── cli/
│   ├── root.go          # CLI commands (modify for post-push)
│   ├── exitcodes.go     # Exit codes
│   └── *_test.go        # Tests
├── config/
│   ├── config.go        # Config loading (add post-push hooks)
│   └── config_test.go   # Tests
├── executor/
│   ├── bump.go          # Bump execution (add post-push phase)
│   └── bump_test.go     # Tests
├── hooks/
│   ├── types.go         # Hook types (add PostPush)
│   ├── runner.go        # Hook execution
│   └── runner_test.go   # Tests
├── tui/
│   ├── model.go         # TUI model (add post-push display)
│   └── *.go             # Other TUI components
├── git/                 # Git operations (no changes needed)
├── version/             # Version parsing (no changes needed)
└── conventional/        # Commit parsing (no changes needed)
```

**Structure Decision**: Single project structure. Changes are additive to existing packages: `hooks`, `config`, `executor`, `cli`, and `tui`.

## Complexity Tracking

> No violations - table not needed

## Implementation Phases (TDD Approach)

### Phase 1: Hook Types Extension

**Test First**:
- Test that `PostPush` constant exists and equals `"post-push"`
- Test that `HookType` can be compared

**Implement**:
- Add `PostPush HookType = "post-push"` to `internal/hooks/types.go`

**Refactor**: N/A (simple addition)

---

### Phase 2: Config Schema Extension

**Test First**:
- Test parsing `.bumpkin.yaml` with `hooks.post-push` array
- Test that empty post-push returns empty slice
- Test that post-push hooks preserve order

**Implement**:
- Add `PostPush []string` to `Hooks` struct in `internal/config/config.go`
- Update `Merge` method to handle post-push hooks

**Refactor**: Ensure consistent handling with pre-tag and post-tag

---

### Phase 3: Executor Post-Push Logic

**Test First**:
- Test post-push hooks execute after successful push
- Test post-push hooks skip when `--no-push` used
- Test post-push hooks skip when push fails
- Test post-push hook failure is warning (tag remains pushed)
- Test multiple post-push hooks execute in order
- Test fail-open: remaining hooks run even if one fails

**Implement**:
- Add `PostPushHooks []string` to `executor.Request`
- Add post-push execution after push in `executor.Execute`
- Implement fail-open behavior (continue on failure)
- Track warnings separately from errors

**Refactor**: Extract common hook execution logic if needed

---

### Phase 4: CLI Integration

**Test First**:
- Test `--no-hooks` skips post-push hooks
- Test `--dry-run` shows post-push hooks without executing
- Test config post-push hooks are passed to executor

**Implement**:
- Pass `PostPushHooks` from config to executor in `root.go`
- Update dry-run output to show post-push hooks

**Refactor**: N/A

---

### Phase 5: TUI Display

**Test First**:
- Test TUI model can receive post-push hook results
- Test TUI displays post-push warnings (not errors)

**Implement**:
- Add post-push hook display to TUI result screen
- Show warnings for failed post-push hooks

**Refactor**: Consistent styling with other hook displays

---

### Phase 6: Integration Testing

**Test First**:
- End-to-end test: bump with post-push hooks
- Test hook timeout (30 seconds)
- Test environment variables available to post-push hooks

**Implement**: Fix any issues found

**Refactor**: Clean up, run linter, ensure all tests pass

## Files to Modify

| File | Changes |
|------|---------|
| `internal/hooks/types.go` | Add `PostPush` HookType constant |
| `internal/hooks/runner.go` | Add fail-open `RunHooksFailOpen` function |
| `internal/hooks/runner_test.go` | Tests for post-push and fail-open behavior |
| `internal/config/config.go` | Add `PostPush` to `Hooks` struct |
| `internal/config/config_test.go` | Tests for post-push config parsing |
| `internal/executor/bump.go` | Add post-push execution phase |
| `internal/executor/bump_test.go` | Tests for post-push behavior |
| `internal/cli/root.go` | Pass post-push hooks, update dry-run output |
| `internal/cli/root_test.go` | Tests for CLI post-push handling |
| `internal/tui/model.go` | Add post-push results display |

## New Files

None required - all changes are extensions to existing files.

## Dependencies

No new dependencies required.
