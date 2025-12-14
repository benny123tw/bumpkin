# Quickstart: Post-Push Hooks Implementation

**Feature**: 002-release-integration
**Date**: 2024-12-14

## Prerequisites

- Go 1.24+
- golangci-lint v2
- Existing bumpkin codebase checked out

## Development Setup

```bash
# Clone and enter repo
cd /path/to/bumpkin
git checkout -b 002-release-integration

# Verify dependencies
go mod download
go mod verify

# Run existing tests
just test
# or: go test ./...

# Run linter
just lint
# or: golangci-lint run
```

## TDD Implementation Steps

### Step 1: Hook Types (5 min)

**Test First** (`internal/hooks/runner_test.go`):
```go
func TestPostPushHookType(t *testing.T) {
    assert.Equal(t, HookType("post-push"), PostPush)
}
```

**Implement** (`internal/hooks/types.go`):
```go
const (
    PreTag   HookType = "pre-tag"
    PostTag  HookType = "post-tag"
    PostPush HookType = "post-push"  // Add this
)
```

**Verify**: `go test ./internal/hooks/...`

---

### Step 2: Config Parsing (10 min)

**Test First** (`internal/config/config_test.go`):
```go
func TestConfigWithPostPushHooks(t *testing.T) {
    yaml := `
hooks:
  post-push:
    - "echo done"
    - "curl webhook"
`
    // Parse and verify cfg.Hooks.PostPush
}
```

**Implement** (`internal/config/config.go`):
```go
type Hooks struct {
    PreTag   []string `yaml:"pre-tag"`
    PostTag  []string `yaml:"post-tag"`
    PostPush []string `yaml:"post-push"`  // Add this
}
```

**Verify**: `go test ./internal/config/...`

---

### Step 3: Fail-Open Hook Runner (15 min)

**Test First** (`internal/hooks/runner_test.go`):
```go
func TestRunHooksFailOpen(t *testing.T) {
    hooks := []Hook{
        {Command: "exit 0", Type: PostPush},
        {Command: "exit 1", Type: PostPush},  // Fails
        {Command: "echo still runs", Type: PostPush},
    }
    results, warnings := RunHooksFailOpen(ctx, hooks, hookCtx)
    assert.Len(t, results, 3)  // All ran
    assert.Len(t, warnings, 1)  // One failed
}
```

**Implement** (`internal/hooks/runner.go`):
```go
func RunHooksFailOpen(ctx context.Context, hooks []Hook, hookCtx *HookContext) ([]*HookResult, []string) {
    var results []*HookResult
    var warnings []string
    
    for _, hook := range hooks {
        result := RunHook(ctx, hook, hookCtx)
        results = append(results, result)
        if !result.Success {
            warnings = append(warnings, fmt.Sprintf("hook '%s' failed: %v", hook.Command, result.Error))
        }
    }
    
    return results, warnings
}
```

**Verify**: `go test ./internal/hooks/...`

---

### Step 4: Executor Integration (20 min)

**Test First** (`internal/executor/bump_test.go`):
```go
func TestExecuteWithPostPushHooks(t *testing.T) {
    // Test: post-push hooks run after push
    // Test: post-push skipped when NoPush=true
    // Test: post-push warnings don't cause error
}
```

**Implement** (`internal/executor/bump.go`):
- Add `PostPushHooks []string` to Request
- Add `PostPushWarnings []string` to Result
- Add post-push execution after push

**Verify**: `go test ./internal/executor/...`

---

### Step 5: CLI Integration (10 min)

**Implement** (`internal/cli/root.go`):
- Pass `cfg.Hooks.PostPush` to executor
- Handle warnings in output

**Verify**: `go test ./internal/cli/...`

---

### Step 6: TUI Display (15 min)

**Implement** (`internal/tui/model.go`):
- Add `PostPushHooks` to Config
- Display warnings after completion

**Verify**: Manual testing with TUI

---

## Verification Checklist

```bash
# All tests pass
just test

# Linter passes
just lint

# Manual test: post-push hooks work
cd /tmp/test-repo
git init
bumpkin major  # Should run post-push hooks after push
```

## Common Issues

1. **Hook not running**: Check `--no-hooks` flag not set
2. **Hook output not visible**: Verify stdout/stderr streaming
3. **Push fails silently**: Check remote exists

## Files Changed Summary

| File | Lines Changed | Type |
|------|---------------|------|
| `internal/hooks/types.go` | +1 | Add constant |
| `internal/hooks/runner.go` | +15 | Add function |
| `internal/hooks/runner_test.go` | +30 | Add tests |
| `internal/config/config.go` | +5 | Add field, update Merge |
| `internal/config/config_test.go` | +20 | Add tests |
| `internal/executor/bump.go` | +25 | Add post-push phase |
| `internal/executor/bump_test.go` | +40 | Add tests |
| `internal/cli/root.go` | +10 | Pass hooks |
| `internal/tui/model.go` | +15 | Add display |

**Estimated Total**: ~160 lines of code + tests
