# Data Model: Release Tool Integration (Post-Push Hooks)

**Feature**: 002-release-integration
**Date**: 2024-12-14

## Overview

This document describes the data model changes required for the post-push hooks feature. Changes are minimal and additive to existing structures.

## Entity Changes

### 1. HookType (Extended)

**File**: `internal/hooks/types.go`

```go
// HookType represents the type of hook
type HookType string

const (
    PreTag   HookType = "pre-tag"
    PostTag  HookType = "post-tag"
    PostPush HookType = "post-push"  // NEW
)
```

**Relationships**: Used by `Hook` struct to categorize hook phase

---

### 2. Hooks Config (Extended)

**File**: `internal/config/config.go`

```go
// Hooks contains hook commands for each phase
type Hooks struct {
    PreTag   []string `yaml:"pre-tag"`
    PostTag  []string `yaml:"post-tag"`
    PostPush []string `yaml:"post-push"`  // NEW
}
```

**Validation Rules**:
- Each entry is a shell command string
- Empty array is valid (no hooks for that phase)
- Order is preserved (execution order)

**YAML Schema**:
```yaml
hooks:
  pre-tag:
    - "command1"
  post-tag:
    - "command2"
  post-push:      # NEW
    - "command3"
    - "command4"
```

---

### 3. Executor Request (Extended)

**File**: `internal/executor/bump.go`

```go
type Request struct {
    Repository    *git.Repository
    BumpType      version.BumpType
    CustomVersion string
    Prefix        string
    Remote        string
    DryRun        bool
    NoPush        bool
    NoHooks       bool
    PreTagHooks   []string
    PostTagHooks  []string
    PostPushHooks []string  // NEW
}
```

**State Transitions**:
```
Request created
    ↓
Execute() called
    ↓
Pre-tag hooks run (if any)
    ↓
Tag created
    ↓
Post-tag hooks run (if any)
    ↓
Tag pushed (if !NoPush)
    ↓
Post-push hooks run (if any, if pushed)  // NEW PHASE
    ↓
Result returned
```

---

### 4. Executor Result (Extended)

**File**: `internal/executor/bump.go`

```go
type Result struct {
    PreviousVersion string
    NewVersion      string
    TagName         string
    CommitHash      string
    TagCreated      bool
    Pushed          bool
    HooksExecuted   int
    PostPushWarnings []string  // NEW: warnings from failed post-push hooks
}
```

**Validation Rules**:
- `PostPushWarnings` is empty if all post-push hooks succeeded
- `PostPushWarnings` contains error messages for failed hooks
- Failed post-push hooks do not affect `Pushed` status

---

### 5. TUI Config (Extended)

**File**: `internal/tui/model.go`

```go
type Config struct {
    Repository    *git.Repository
    Prefix        string
    Remote        string
    DryRun        bool
    NoPush        bool
    NoHooks       bool
    PreTagHooks   []string
    PostTagHooks  []string
    PostPushHooks []string  // NEW
}
```

---

## New Functions

### RunHooksFailOpen

**File**: `internal/hooks/runner.go`

```go
// RunHooksFailOpen executes hooks but continues on failure
// Returns all results and a slice of warnings for failed hooks
func RunHooksFailOpen(
    ctx context.Context, 
    hooks []Hook, 
    hookCtx *HookContext,
) ([]*HookResult, []string)
```

**Behavior**:
- Executes all hooks in order
- Does not stop on failure
- Collects warnings for failed hooks
- Returns both results and warnings

---

## Environment Variables

No changes to environment variables. Post-push hooks receive the same variables as other hooks:

| Variable | Description |
|----------|-------------|
| `BUMPKIN_VERSION` | New version (e.g., "1.2.3") |
| `BUMPKIN_PREVIOUS_VERSION` | Previous version |
| `BUMPKIN_TAG` | Full tag name (e.g., "v1.2.3") |
| `BUMPKIN_PREFIX` | Tag prefix (e.g., "v") |
| `BUMPKIN_REMOTE` | Remote name (e.g., "origin") |
| `BUMPKIN_COMMIT` | Commit hash |
| `BUMPKIN_DRY_RUN` | "true" or "false" |
| `VERSION` | Alias for BUMPKIN_VERSION |
| `TAG` | Alias for BUMPKIN_TAG |

---

## Execution Flow Diagram

```
┌─────────────────────────────────────────────────────────────────┐
│                     Bump Execution Flow                          │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  ┌──────────┐                                                   │
│  │  Start   │                                                   │
│  └────┬─────┘                                                   │
│       │                                                         │
│       ▼                                                         │
│  ┌──────────────────┐                                          │
│  │ Pre-tag hooks    │ ──fail──▶ ERROR (stop)                   │
│  └────────┬─────────┘                                          │
│           │ pass                                                │
│           ▼                                                     │
│  ┌──────────────────┐                                          │
│  │ Create tag       │ ──fail──▶ ERROR (stop)                   │
│  └────────┬─────────┘                                          │
│           │ pass                                                │
│           ▼                                                     │
│  ┌──────────────────┐                                          │
│  │ Post-tag hooks   │ ──fail──▶ ERROR (tag created)            │
│  └────────┬─────────┘                                          │
│           │ pass                                                │
│           ▼                                                     │
│  ┌──────────────────┐                                          │
│  │ Push tag         │ ──fail──▶ ERROR (skip post-push)         │
│  │ (if !NoPush)     │                                          │
│  └────────┬─────────┘                                          │
│           │ pass                                                │
│           ▼                                                     │
│  ┌──────────────────┐                                          │
│  │ Post-push hooks  │ ──fail──▶ WARNING (continue)  ◀── NEW    │
│  │ (fail-open)      │                                          │
│  └────────┬─────────┘                                          │
│           │                                                     │
│           ▼                                                     │
│  ┌──────────────────┐                                          │
│  │ Return Result    │                                          │
│  │ (with warnings)  │                                          │
│  └──────────────────┘                                          │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

---

## Summary of Changes

| Entity | Change Type | Description |
|--------|-------------|-------------|
| `HookType` | Add constant | `PostPush = "post-push"` |
| `Hooks` | Add field | `PostPush []string` |
| `Request` | Add field | `PostPushHooks []string` |
| `Result` | Add field | `PostPushWarnings []string` |
| `TUI Config` | Add field | `PostPushHooks []string` |
| `hooks` package | Add function | `RunHooksFailOpen()` |
