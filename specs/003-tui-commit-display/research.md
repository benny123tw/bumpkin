# Research: TUI Commit Display Enhancement

**Feature**: 003-tui-commit-display
**Date**: 2024-12-14

## Research Summary

This document captures research decisions for the TUI commit display enhancement. All technical context is based on existing codebase patterns.

## Decisions

### 1. Commit Type Color Scheme

**Decision**: Use lipgloss ANSI colors matching bumpp's visual style

**Color Mapping**:
| Type | Lipgloss Color | Hex Equivalent |
|------|----------------|----------------|
| feat | Color("154") | Lime/Green |
| fix | Color("214") | Yellow/Orange |
| docs | Color("75") | Light Blue |
| chore | Color("245") | Gray |
| refactor | Color("81") | Cyan |
| test | Color("213") | Magenta |
| style | Color("245") | Gray |
| perf | Color("214") | Orange |
| ci | Color("245") | Gray |
| build | Color("245") | Gray |
| Breaking | Color("196") bg | Red background |

**Rationale**: 
- Colors are distinct and accessible
- Matches common conventions (green=feature, yellow=fix, red=breaking)
- Uses existing lipgloss color system

**Alternatives Considered**:
- Custom RGB colors: Rejected - less portable across terminals
- Only highlight breaking: Rejected - loses information density

---

### 2. Breaking Change Detection

**Decision**: Detect breaking changes via `!` after commit type (e.g., `feat!:`)

**Rationale**:
- Follows conventional commit specification
- Already partially supported in existing parser
- Simple regex pattern: `^(\w+)(!)?:`

**Alternatives Considered**:
- Parse commit body for `BREAKING CHANGE:`: Deferred - body not currently parsed
- Use separate marker: Rejected - not conventional commit standard

---

### 3. Commit Display Format

**Decision**: Format as `<hash>  <type> : <description>`

Example:
```
383d79a  feat : add power function
d677b41  docs : add package documentation
```

**Rationale**:
- Matches bumpp's format for familiarity
- Hash provides quick reference
- Type badge is visually prominent
- Colon separates type from description

**Alternatives Considered**:
- `[type] description`: Rejected - less visual distinction
- No hash: Rejected - loses useful reference

---

### 4. Persistent Commit Display

**Decision**: Show commits above version selector on all selection screens

**Layout**:
```
🎃 bumpkin

Current version: v1.0.0

5 Commits since the last version:

383d79a  feat : add power function
d677b41  docs : add package documentation

? Current version 1.0.0 >
        major 2.0.0
        minor 1.1.0
      > patch 1.0.1
        ...
```

**Rationale**:
- Users need commit context when selecting version
- Avoids screen switching/memory burden
- Matches bumpp's behavior

**Alternatives Considered**:
- Collapsible commits: Rejected - adds complexity
- Side-by-side: Rejected - terminal width constraints

---

### 5. Commit Truncation

**Decision**: Show max 10 commits, then "and X more commits..."

**Rationale**:
- Prevents screen overflow
- 10 commits is sufficient context for most decisions
- "X more" indicator informs user of hidden commits

**Alternatives Considered**:
- Scrollable list: Deferred - adds complexity
- Dynamic based on terminal height: Deferred - adds complexity

---

### 6. Non-Conventional Commit Handling

**Decision**: Display without badge, just `<hash>  <message>`

Example:
```
a1b2c3d  Updated readme
```

**Rationale**:
- Graceful degradation
- Doesn't force conventional commit format
- Still shows relevant information

---

## Existing Patterns to Follow

### From `internal/tui/styles.go`:
- Use existing color variables where applicable
- Follow naming convention: `TypeStyle` suffix
- Use lipgloss.NewStyle() pattern

### From `internal/tui/commits.go`:
- Existing `RenderCommitList` function to extend
- Already has hash and message extraction

### From `internal/conventional/parser.go`:
- Existing `ParseCommit` function
- Has `Type` field in `ConventionalCommit` struct

## No External Research Needed

All decisions are based on:
- Existing codebase patterns
- Spec requirements
- Bumpp reference implementation
- Lipgloss documentation

## Next Steps

Proceed to Phase 1: Design & Contracts
- Generate data-model.md
- Generate quickstart.md
- No API contracts needed (TUI-only feature)
