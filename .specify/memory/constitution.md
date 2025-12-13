<!--
SYNC IMPACT REPORT
==================
Version change: N/A → 1.0.0 (initial ratification)
Modified principles: N/A (initial)
Added sections:
  - Principle I: Code Quality (golangci-lint)
  - Section: Code Quality Configuration
  - Section: Development Workflow
  - Section: Governance
Removed sections: N/A (initial)
Templates requiring updates:
  - .specify/templates/plan-template.md: ✅ No updates required (Constitution Check generic)
  - .specify/templates/spec-template.md: ✅ No updates required
  - .specify/templates/tasks-template.md: ✅ No updates required
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
1. Run `golangci-lint run` to check for linting errors
2. Run `golangci-lint fmt` to apply formatting (or use editor integration)
3. Address all reported issues or add justified `//nolint` directives

**CI/CD Integration**:
- Linting MUST be part of the CI pipeline
- PRs MUST pass lint checks before merge

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

**Version**: 1.0.0 | **Ratified**: 2025-12-13 | **Last Amended**: 2025-12-13
