# Code Refactoring Plan

## Overview

This plan identifies code quality improvements, technical debt, and refactoring opportunities discovered during codebase analysis.

---

## 4. Viewport Dimension Validation

**Location**: `internal/tui/model.go:174-184`

**Issue**: Window resize handling can result in zero or negative pane heights if the window is very small.

```go
commitsPaneHeight := availableHeight * 30 / 100
if commitsPaneHeight < 3 {
    commitsPaneHeight = 3
}
m.commitsPane.Height = commitsPaneHeight - 2  // Could become 1
```

**Recommendation**: Add minimum dimension validation after all calculations:

```go
finalHeight := commitsPaneHeight - 2
if finalHeight < 1 {
    finalHeight = 1
}
m.commitsPane.Height = finalHeight
```

**Priority**: Low
**Effort**: Small

---

## 5. Prerelease Validation in Version Parsing

**Location**: `internal/version/prerelease.go`

**Issue**: Hardcoded dot separator in prerelease parsing. Malformed prerelease strings produce unclear errors late in the bump process.

**Recommendation**: Add early validation in `Version.Parse()` to reject malformed prerelease versions with clear error messages.

**Priority**: Low
**Effort**: Small

---

## 6. Remote URL Validation

**Location**: `internal/git/push.go`

**Issue**: `GetRemoteURL()` doesn't validate URLs. SSH, HTTPS, and file:// paths could fail silently downstream.

**Recommendation**: Add basic URL validation or log warnings for unusual formats:

```go
func isValidRemoteURL(url string) bool {
    return strings.HasPrefix(url, "git@") ||
           strings.HasPrefix(url, "https://") ||
           strings.HasPrefix(url, "ssh://") ||
           strings.HasPrefix(url, "git://")
}
```

**Priority**: Low
**Effort**: Small

---

## 7. Add Structured Logging

**Location**: Throughout codebase

**Issue**: No structured logging for debugging. Users cannot see what the tool is doing under the hood, especially during hook execution.

**Recommendation**:
1. Add `--verbose` or `--debug` flag
2. Use `log/slog` (Go 1.21+) for structured logging
3. Log key events: hook execution, git operations, version calculations

**Priority**: Medium
**Effort**: Medium

---

## 8. Improve Test Coverage

**Locations**:
- `internal/executor/bump.go` - no tests
- `internal/git/commits.go` - no direct tests
- `internal/config/config.go` - minimal tests

**Issue**: Critical paths lack test coverage. Current coverage is below recommended 80% threshold.

**Recommendation**:
1. Add executor tests with mocked git repository
2. Add git operation tests using in-memory repositories
3. Add config parsing edge case tests
4. Add integration tests for end-to-end workflows

**Priority**: High
**Effort**: Large

---

## 9. Hook Security Documentation

**Location**: `internal/hooks/runner.go`, `README.md`

**Issue**: Hooks are user-defined from config files with no sanitization. While intentional, security implications should be documented.

**Recommendation**: Add security section to README:
- Explain that hooks execute arbitrary commands
- Recommend reviewing hooks in shared configs
- Consider adding `--no-hooks` reminder in CI environments

**Priority**: Low
**Effort**: Small

---

## 10. Commit History Pagination

**Location**: `internal/git/commits.go`

**Issue**: Full commit history loaded into memory. Works for typical repos but could cause issues on very large repositories.

**Recommendation**: Add optional limit parameter:

```go
func (r *Repository) GetCommitsSinceTagWithLimit(tag string, limit int) ([]Commit, error)
```

Also consider adding `--max-commits` CLI flag.

**Priority**: Low
**Effort**: Medium

---

## Implementation Order

1. **Phase 1** (Quick Wins):
   - Deprecated state cleanup (#2)
   - Viewport dimension validation (#4)
   - Flag validation helper (#1)

2. **Phase 2** (Quality):
   - Post-tag error type (#3)
   - Prerelease validation (#5)
   - Remote URL validation (#6)

3. **Phase 3** (Observability):
   - Structured logging (#7)
   - Hook security documentation (#9)

4. **Phase 4** (Testing):
   - Improve test coverage (#8)

5. **Phase 5** (Performance):
   - Commit history pagination (#10)
