<!--
SYNC IMPACT REPORT
==================
Version change: 1.0.0 → 1.1.0
Modified principles: None
Added sections:
  - Principle II: Test-Driven Development (TDD)
Removed sections: None
Templates requiring updates:
  - .specify/templates/plan-template.md: ✅ No updates required (Constitution Check generic)
  - .specify/templates/spec-template.md: ✅ No updates required (already supports test scenarios)
  - .specify/templates/tasks-template.md: ✅ Already aligned (tests-first pattern documented)
Follow-up TODOs: None
-->

# Bumpkin Constitution

## Core Principles

### I. Code Quality

All Go code MUST pass golangci-lint checks before merge.

**Requirements**:
- Code MUST pass all enabled linters without suppressions unless explicitly justified
- All `//nolint` directives MUST include a comment explaining the justification
- Formatting MUST be applied via the configured formatters (gofmt, gofumpt, goimports, golines)
- Generated code, mocks, third-party code, builtins, and examples are excluded from linting

**Rationale**: Consistent code quality and formatting reduces cognitive load during review,
catches bugs early, and ensures the codebase remains maintainable over time.

### II. Test-Driven Development (TDD)

All feature implementation MUST follow the TDD cycle: Red → Green → Refactor.

**Requirements**:
- Tests MUST be written BEFORE implementation code
- Tests MUST fail initially (Red phase) to confirm they detect the missing functionality
- Implementation code MUST be written only to make failing tests pass (Green phase)
- Code MUST be refactored after tests pass while maintaining all tests green (Refactor phase)
- No production code changes are permitted without a corresponding failing test first

**Workflow**:
1. **Red**: Write a test that defines expected behavior. Run it. It MUST fail.
2. **Green**: Write the minimum code necessary to make the test pass. No more.
3. **Refactor**: Improve code structure, remove duplication, enhance readability—all tests MUST remain green.

**Rationale**: TDD ensures comprehensive test coverage by design, produces modular and decoupled code,
provides immediate feedback during development, and creates living documentation of expected behavior.
The discipline of writing tests first prevents over-engineering and keeps implementation focused on
actual requirements.

## Code Quality Configuration

The project uses golangci-lint version 2 configuration stored at `.golangci.yml`.

**Enabled Linters**:
- bodyclose, dogsled, dupl, errcheck, exhaustive, goconst, gocritic
- gocyclo, goprintffuncname, gosec, govet, ineffassign, misspell
- nakedret, noctx, nolintlint, rowserrcheck, staticcheck, unconvert
- unparam, unused, whitespace, copyloopvar, predeclared

**Enabled Formatters**:
- gofmt (with simplify enabled)
- gofumpt (with extra-rules enabled)
- goimports
- golines

**Excluded Paths**:
- `internal/mocks` - mock implementations
- `third_party$` - vendored third-party code
- `builtin$` - builtin overrides
- `examples$` - example code

## Development Workflow

**Before committing**:
1. Ensure all new code follows TDD cycle (test written first, failed, then implementation)
2. Run `go test ./...` to verify all tests pass
3. Run `golangci-lint run` to check for linting errors
4. Run `golangci-lint fmt` to apply formatting (or use editor integration)
5. Address all reported issues or add justified `//nolint` directives

**CI/CD Integration**:
- All tests MUST pass before merge
- Linting MUST be part of the CI pipeline
- PRs MUST pass both test and lint checks before merge

## Governance

This constitution establishes the foundational principles for the Bumpkin project.

**Amendment Process**:
1. Propose changes via PR with clear rationale
2. Document impact on existing code and workflows
3. Update version following semantic versioning:
   - MAJOR: Backward-incompatible principle changes or removals
   - MINOR: New principles or material expansions
   - PATCH: Clarifications, wording, typo fixes
4. All team members must acknowledge changes

**Compliance**:
- All PRs and code reviews MUST verify compliance with these principles
- Violations MUST be addressed before merge unless explicitly waived with documented justification
- TDD compliance MUST be evidenced by commit history showing test commits before implementation commits

**Version**: 1.1.0 | **Ratified**: 2025-12-13 | **Last Amended**: 2025-12-15
