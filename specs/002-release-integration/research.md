# Research: Release Tool Integration (Post-Push Hooks)

**Feature**: 002-release-integration
**Date**: 2024-12-14

## Research Summary

This document captures research decisions for the post-push hooks feature. All technical context items were already known from the existing codebase - no external research was needed.

## Decisions

### 1. Hook Phase Naming

**Decision**: Use `post-push` as the hook phase name

**Rationale**: 
- Follows existing naming convention (`pre-tag`, `post-tag`)
- Clear semantic meaning: runs after push operation
- Consistent with other tools (git hooks use similar naming)

**Alternatives Considered**:
- `after-push`: Rejected - inconsistent with existing `post-*` convention
- `on-push`: Rejected - implies during, not after

---

### 2. Fail-Open Behavior for Post-Push Hooks

**Decision**: Post-push hooks use fail-open behavior (continue on failure, report warnings)

**Rationale**:
- The tag is already pushed when post-push hooks run
- Failing would not undo the push
- Notifications/integrations should not block the workflow
- User can still see failures in output

**Alternatives Considered**:
- Fail-closed (stop on first failure): Rejected - no benefit since tag already pushed
- Configurable behavior: Rejected - adds complexity, fail-open is the right default

---

### 3. Hook Timeout

**Decision**: 30 second default timeout for all hooks (including post-push)

**Rationale**:
- Matches existing hook timeout behavior
- Long enough for most HTTP calls (Slack, webhooks)
- Short enough to not block the terminal indefinitely

**Alternatives Considered**:
- Per-hook configurable timeout: Deferred - can add later if needed
- Longer default (60s): Rejected - 30s is sufficient for notifications

---

### 4. Hook Execution Context

**Decision**: Post-push hooks receive the same environment variables as pre-tag/post-tag

**Rationale**:
- Consistent behavior across all hook phases
- All relevant version info is available
- No new environment variables needed

**Alternatives Considered**:
- Add `BUMPKIN_PUSHED=true` variable: Deferred - can add later if useful
- Add push result details: Rejected - not needed for notifications

---

### 5. TDD Implementation Order

**Decision**: Follow this order for TDD:
1. Types (foundation)
2. Config (data input)
3. Executor (core logic)
4. CLI (user interface)
5. TUI (visual interface)

**Rationale**:
- Bottom-up approach builds on stable foundations
- Each layer depends on the previous
- Tests can use real implementations, not mocks

**Alternatives Considered**:
- Top-down (CLI first): Rejected - would require mocking
- Parallel development: Rejected - dependencies exist between layers

---

### 6. Warning vs Error Handling

**Decision**: Create a `Warning` type or return warnings separately from errors

**Rationale**:
- Post-push failures are warnings, not errors
- Exit code should still be 0 if only warnings
- User should see warning output but not be blocked

**Alternatives Considered**:
- Use stderr for warnings: Partial - also need to track for exit code
- Log-only warnings: Rejected - user might miss them

---

## Existing Patterns to Follow

### From `internal/hooks/runner.go`:
- Use `RunHooks` pattern but create `RunHooksFailOpen` variant
- Keep `HookResult` structure unchanged
- Stream output to stdout/stderr

### From `internal/config/config.go`:
- Add `PostPush []string` following `PreTag`/`PostTag` pattern
- Update `Merge` method consistently

### From `internal/executor/bump.go`:
- Add post-push phase after push, before return
- Track warnings in Result struct

## No External Research Needed

All decisions are based on:
- Existing codebase patterns
- Spec requirements (CI/CD-first approach)
- Go best practices already in use

## Next Steps

Proceed to Phase 1: Design & Contracts
- Generate data-model.md
- Generate quickstart.md
- No API contracts needed (CLI tool)
