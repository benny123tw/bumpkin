# Tasks: TUI Commit Display Enhancement

**Input**: Design documents from `/specs/003-tui-commit-display/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, quickstart.md

**Organization**: Tasks are grouped by user story to enable independent implementation and testing.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (US1, US2, US3, US4)
- Include exact file paths in descriptions

## User Stories Summary

| Story | Title | Priority |
|-------|-------|----------|
| US1 | Direct Version Selection with Commit History | P1 |
| US2 | Colored Commit Type Badges | P1 |
| US3 | Breaking Change Highlighting | P1 |
| US4 | Commit Display Format | P2 |

---

## Phase 1: Setup

**Purpose**: Create feature branch and verify existing infrastructure

- [X] T001 Verify on feature branch `003-tui-commit-display`
- [X] T002 Verify all existing tests pass with `go test ./...`
- [X] T003 Verify linter passes with `golangci-lint run`

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Add CommitDisplay struct and parsing function that all user stories depend on

- [X] T004 [P] Add `CommitDisplay` struct to `internal/tui/commits.go`
- [X] T005 [P] Add `conventionalCommitRegex` pattern to `internal/tui/commits.go`
- [X] T006 Add `ParseCommitForDisplay` function to `internal/tui/commits.go`
- [X] T007 Verify parsing works with manual test

**Checkpoint**: ✅ Foundation ready - CommitDisplay parsing available

---

## Phase 3: User Story 1 - Direct Version Selection with Commit History (Priority: P1)

**Goal**: Skip commit preview screen, go directly to version selection with commits displayed

**Independent Test**: Run bumpkin, verify it goes directly to version selection with commits (no intermediate screen)

### Implementation

- [X] T008 [US1] Add `RenderCommitListWithBadges` function to `internal/tui/commits.go`
- [X] T009 [US1] Modify `renderVersionSelectView` in `internal/tui/model.go` to include commit history
- [X] T010 [US1] Add commit count header (e.g., "5 Commits since the last version:")
- [X] T011 [US1] Add truncation with "...and X more commit(s)" indicator when > 10 commits
- [X] T012 [US1] Handle edge case: show "No commits since last tag" when empty
- [X] T013a [US1] Skip `StateCommitList` - go directly from loading to `StateVersionSelect`
- [X] T013b [US1] Remove or deprecate `renderCommitListView` function (no longer used)
- [ ] T013c [US1] Manual test: verify bumpkin goes directly to version selection with commits

**Checkpoint**: ✅ User Story 1 complete - streamlined flow to version selection

---

## Phase 4: User Story 2 - Colored Commit Type Badges (Priority: P1)

**Goal**: Add colored badges for each commit type (feat, fix, docs, etc.)

**Independent Test**: Run bumpkin with conventional commits, verify each type has distinct color

### Implementation

- [X] T014 [P] [US2] Add `FeatStyle` (lime/green) to `internal/tui/styles.go`
- [X] T015 [P] [US2] Add `FixStyle` (yellow) to `internal/tui/styles.go`
- [X] T016 [P] [US2] Add `DocsStyle` (blue) to `internal/tui/styles.go`
- [X] T017 [P] [US2] Add `ChoreStyle` (gray) to `internal/tui/styles.go`
- [X] T018 [P] [US2] Add `RefactorStyle` (cyan) to `internal/tui/styles.go`
- [X] T019 [P] [US2] Add `TestStyle` (magenta) to `internal/tui/styles.go`
- [X] T020 [P] [US2] Add `PerfStyle` (orange) to `internal/tui/styles.go`
- [X] T021 [US2] Add `CommitTypeStyles` map linking types to styles in `internal/tui/styles.go`
- [X] T022 [US2] Add `GetCommitTypeStyle` function to `internal/tui/styles.go`
- [X] T023 [US2] Update `RenderCommitListWithBadges` to apply type styles in `internal/tui/commits.go`
- [ ] T024 [US2] Manual test: verify feat=green, fix=yellow, docs=blue, chore=gray

**Checkpoint**: ✅ User Story 2 complete - commit types have colored badges

---

## Phase 5: User Story 3 - Breaking Change Highlighting (Priority: P1)

**Goal**: Highlight breaking changes (commits with `!`) with red background

**Independent Test**: Create commits with `feat!:`, verify red background on type badge

### Implementation

- [X] T025 [P] [US3] Add `BreakingStyle` (red background) to `internal/tui/styles.go`
- [X] T026 [US3] Update `ParseCommitForDisplay` to detect `!` in commit type in `internal/tui/commits.go`
- [X] T027 [US3] Update `GetCommitTypeStyle` to return breaking style when `isBreaking=true` in `internal/tui/styles.go`
- [X] T028 [US3] Update `RenderCommitListWithBadges` to pass isBreaking flag in `internal/tui/commits.go`
- [ ] T029 [US3] Manual test: verify `feat!:` and `fix!:` show red background

**Checkpoint**: ✅ User Story 3 complete - breaking changes highlighted in red

---

## Phase 6: User Story 4 - Commit Display Format (Priority: P2)

**Goal**: Clean format matching bumpp: `<hash>  <type> : <description>`

**Independent Test**: Verify consistent formatting across all commit types

### Implementation

- [X] T030 [US4] Ensure hash is truncated to 7 characters in `ParseCommitForDisplay` in `internal/tui/commits.go`
- [X] T031 [US4] Format output as `<hash>  <type> : <description>` in `RenderCommitListWithBadges`
- [X] T032 [US4] Handle non-conventional commits: display as `<hash>  <message>` without badge
- [ ] T033 [US4] Manual test: verify format matches bumpp style

**Checkpoint**: ✅ User Story 4 complete - clean consistent formatting

---

## Phase 7: Polish & Cross-Cutting Concerns

**Purpose**: Final verification and cleanup

- [X] T034 Run all tests with `go test ./...`
- [X] T035 Run linter with `golangci-lint run` and fix any issues
- [X] T036 Apply formatting with `golangci-lint fmt`
- [ ] T037 Manual integration test: full TUI flow with various commit types
- [X] T038 Commit all changes with conventional commit message

---

## Dependencies & Execution Order

### Phase Dependencies

```
Phase 1 (Setup) → Phase 2 (Foundational) → Phase 3-6 (User Stories) → Phase 7 (Polish)
                                           ↓
                               Can run in parallel after Phase 2:
                               - US1 (P1) - Persistent display
                               - US2 (P1) - Colored badges
                               - US3 (P1) - Breaking changes
                               - US4 (P2) - Format (depends on US1-3)
```

### User Story Dependencies

| Story | Depends On | Can Parallel With |
|-------|------------|-------------------|
| US1 | Phase 2 (Foundational) | US2, US3 (styles) |
| US2 | Phase 2 (Foundational) | US1, US3 |
| US3 | Phase 2 (Foundational) | US1, US2 |
| US4 | US1, US2, US3 | - |

---

## Parallel Opportunities

### Phase 2: Foundational

```bash
# Can run in parallel:
T004: Add CommitDisplay struct
T005: Add regex pattern
```

### Phase 4: User Story 2 Styles

```bash
# All style definitions can run in parallel:
T014: FeatStyle
T015: FixStyle
T016: DocsStyle
T017: ChoreStyle
T018: RefactorStyle
T019: TestStyle
T020: PerfStyle
T025: BreakingStyle
```

---

## Implementation Strategy

### MVP First (User Story 1 + 2 + 3)

1. ✅ Phase 1: Setup (T001-T003)
2. ✅ Phase 2: Foundational (T004-T007)
3. ✅ Phase 3: User Story 1 - Persistent Display (T008-T013)
4. ✅ Phase 4: User Story 2 - Colored Badges (T014-T024)
5. ✅ Phase 5: User Story 3 - Breaking Changes (T025-T029)
6. **STOP and VALIDATE**: Core functionality complete
7. Can ship MVP here!

### Full Implementation

Continue with:
- Phase 6: User Story 4 - Format polish (T030-T033)
- Phase 7: Polish (T034-T038)

---

## Notes

- [P] tasks can run in parallel (different files)
- [US#] label maps task to specific user story
- All P1 stories should be completed together for MVP
- Commit after each phase or logical group
- Total: 38 tasks
