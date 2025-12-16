# Tasks: Expandable Commit History

**Input**: Design documents from `/specs/005-expandable-commit-history/`  
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, quickstart.md

**Tests**: Required per Constitution Principle II (TDD). Tests MUST be written BEFORE implementation.

**Organization**: Tasks grouped by user story for independent implementation and testing.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (US1, US2, US3, US4)
- Exact file paths included in descriptions

## Path Conventions

- **Source**: `internal/tui/` (existing TUI package)
- **Tests**: `internal/tui/*_test.go` (co-located with source)

---

## Phase 1: Setup

**Purpose**: Add viewport dependency and foundational types

- [ ] T001 Verify `charmbracelet/bubbles/viewport` is available (already in go.mod as bubbles v0.21.0)
- [ ] T002 [P] Define `PaneType` enum (`PaneVersion`, `PaneCommits`) in `internal/tui/pane.go`
- [ ] T003 [P] Add `FocusedBorderStyle` and `UnfocusedBorderStyle` in `internal/tui/styles.go`

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core model changes that ALL user stories depend on

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

### Tests for Foundational Phase

- [ ] T004 Write test for `Model` with new pane fields in `internal/tui/model_test.go` - verify `focusedPane` defaults to `PaneVersion`, `commitsPane` initialized with zero dimensions

### Implementation for Foundational Phase

- [ ] T005 Add `commitsPane viewport.Model` field to `Model` struct in `internal/tui/model.go`
- [ ] T006 Add `focusedPane PaneType` field to `Model` struct in `internal/tui/model.go`
- [ ] T007 Add `showingDetail bool` and `selectedCommitIndex int` fields to `Model` struct in `internal/tui/model.go`
- [ ] T008 Initialize new fields in `New()` function in `internal/tui/model.go` - `focusedPane = PaneVersion`, `commitsPane = viewport.New(0, 0)`
- [ ] T009 Update `tea.WindowSizeMsg` handler to calculate pane heights (~30% commits, ~70% version) and resize `commitsPane` in `internal/tui/model.go`

**Checkpoint**: Model has all new fields, viewport initializes and resizes correctly

---

## Phase 3: User Story 1 - Navigate Between Panes (Priority: P1) 🎯 MVP

**Goal**: Users can switch focus between commits pane and version pane using Tab/Shift-Tab

**Independent Test**: Run bumpkin, press Tab to switch between panes, verify focus indicator changes

### Tests for User Story 1

- [ ] T010 [US1] Write test: given `focusedPane == PaneVersion`, when Tab pressed, then `focusedPane == PaneCommits` in `internal/tui/model_test.go`
- [ ] T011 [US1] Write test: given `focusedPane == PaneCommits`, when Tab pressed, then `focusedPane == PaneVersion` in `internal/tui/model_test.go`
- [ ] T012 [US1] Write test: given `focusedPane == PaneCommits`, when Shift+Tab pressed, then `focusedPane == PaneVersion` in `internal/tui/model_test.go`
- [ ] T013 [P] [US1] Write test for dual-pane rendering: verify both panes rendered with correct border styles based on focus in `internal/tui/view_test.go`

### Implementation for User Story 1

- [ ] T014 [US1] Handle Tab and Shift+Tab key events to toggle `focusedPane` in `handleKeyPress()` in `internal/tui/model.go`
- [ ] T015 [US1] Add `RenderCommitListForViewport(commits []*git.Commit) string` function (no truncation) in `internal/tui/commits.go`
- [ ] T016 [US1] Populate `commitsPane.SetContent()` with rendered commits in `RepoLoadedMsg` handler in `internal/tui/model.go`
- [ ] T017 [US1] Refactor `renderVersionSelectView()` to render dual-pane layout using `lipgloss.JoinVertical()` in `internal/tui/model.go`
- [ ] T018 [US1] Apply `FocusedBorderStyle` or `UnfocusedBorderStyle` to each pane based on `focusedPane` value in `internal/tui/model.go`
- [ ] T019 [US1] Update help text at bottom to show "[Tab] switch pane" in `renderHelp()` in `internal/tui/model.go`

**Checkpoint**: User Story 1 complete - Tab switches focus, both panes visible with focus indicator

---

## Phase 4: User Story 2 - Scroll Through All Commits (Priority: P1)

**Goal**: Users can scroll through all commits in the commits pane using arrow keys

**Independent Test**: Run bumpkin with 50+ commits, Tab to commits pane, scroll through entire list

### Tests for User Story 2

- [ ] T020 [US2] Write test: given commits pane focused, when down arrow pressed, then viewport scrolls down in `internal/tui/model_test.go`
- [ ] T021 [US2] Write test: given commits pane focused, when up arrow pressed, then viewport scrolls up in `internal/tui/model_test.go`
- [ ] T022 [US2] Write test: given version pane focused, when down arrow pressed, then version selector moves (not viewport) in `internal/tui/model_test.go`
- [ ] T023 [P] [US2] Write test for scroll position indicator rendering (e.g., "[5/25]") in `internal/tui/commits_test.go`

### Implementation for User Story 2

- [ ] T024 [US2] Route arrow key events to `commitsPane.Update()` when `focusedPane == PaneCommits` in `handleKeyPress()` in `internal/tui/model.go`
- [ ] T025 [US2] Keep existing arrow key handling for version selection when `focusedPane == PaneVersion` in `internal/tui/model.go`
- [ ] T026 [US2] Add scroll position indicator to commits pane header (e.g., "Commits (5/25)") using `commitsPane.YOffset` and total commit count in `internal/tui/model.go`
- [ ] T027 [US2] Ensure j/k keys also work for scrolling (match existing vim-style navigation) in `internal/tui/model.go`

**Checkpoint**: User Story 2 complete - commits pane scrolls, position indicator updates

---

## Phase 5: User Story 3 - View Full Commit Message (Priority: P2)

**Goal**: Users can view full commit message in a detail overlay

**Independent Test**: Tab to commits pane, press Enter on a commit with long message, verify overlay shows full message

### Tests for User Story 3

- [ ] T028 [US3] Write test: given commits pane focused, when Enter pressed, then `showingDetail == true` in `internal/tui/model_test.go`
- [ ] T029 [US3] Write test: given overlay showing, when Escape pressed, then `showingDetail == false` in `internal/tui/model_test.go`
- [ ] T030 [US3] Write test: given overlay showing, when Enter pressed, then `showingDetail == false` in `internal/tui/model_test.go`
- [ ] T031 [P] [US3] Write test for overlay rendering with full commit details in `internal/tui/overlay_test.go`

### Implementation for User Story 3

- [ ] T032 [P] [US3] Create `internal/tui/overlay.go` with `RenderCommitDetailOverlay(commit *git.Commit, width, height int) string` function
- [ ] T033 [US3] Handle Enter key in commits pane to set `showingDetail = true` and `selectedCommitIndex` in `internal/tui/model.go`
- [ ] T034 [US3] Handle Escape and Enter keys to dismiss overlay (`showingDetail = false`) in `internal/tui/model.go`
- [ ] T035 [US3] Block pane switching (Tab) when `showingDetail == true` in `internal/tui/model.go`
- [ ] T036 [US3] Render overlay on top of pane layout when `showingDetail == true` in `View()` in `internal/tui/model.go`
- [ ] T037 [US3] Add `OverlayStyle` with centered box styling in `internal/tui/styles.go`

**Checkpoint**: User Story 3 complete - Enter shows overlay, Escape/Enter dismisses it

---

## Phase 6: User Story 4 - Select Version While Viewing Commits (Priority: P2)

**Goal**: Version selection works correctly while commits pane is visible

**Independent Test**: Tab to version pane, select version, confirm - commits pane remains visible throughout

### Tests for User Story 4

- [ ] T038 [US4] Write test: given version pane focused, when Enter pressed on version, then proceeds to confirmation in `internal/tui/model_test.go`
- [ ] T039 [US4] Write test: given confirmation showing, commits pane scroll position preserved when returning to version select in `internal/tui/model_test.go`
- [ ] T040 [P] [US4] Write test: verify commits pane remains visible (dimmed) when version pane is focused in `internal/tui/view_test.go`

### Implementation for User Story 4

- [ ] T041 [US4] Verify existing Enter handling for version selection still works in `handleVersionSelectKeys()` in `internal/tui/model.go`
- [ ] T042 [US4] Preserve `commitsPane.YOffset` (scroll position) when switching panes in `internal/tui/model.go`
- [ ] T043 [US4] Ensure confirmation view returns to version select state correctly (no regressions) in `internal/tui/model.go`

**Checkpoint**: User Story 4 complete - full workflow works: view commits → select version → confirm

---

## Phase 7: Polish & Cross-Cutting Concerns

**Purpose**: Edge cases, small terminal handling, and final cleanup

- [ ] T044 [P] Handle edge case: 0 commits - show "No new commits" message in commits pane in `internal/tui/model.go`
- [ ] T045 [P] Handle edge case: fewer than 5 commits - hide scroll indicator in `internal/tui/model.go`
- [ ] T046 [P] Handle edge case: empty commit message - show "(no message)" placeholder in `internal/tui/commits.go`
- [ ] T047 Implement small terminal fallback: if `height < 16`, show single pane with Tab switching in `internal/tui/model.go`
- [ ] T048 Run `golangci-lint run` and fix any issues
- [ ] T049 Run `go test ./internal/tui/...` and verify all tests pass
- [ ] T050 Manual testing: run `go run ./cmd/bumpkin` in a repo with various commit counts

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup - BLOCKS all user stories
- **User Stories (Phase 3-6)**: All depend on Foundational completion
  - US1 and US2 are both P1 and can run in parallel
  - US3 and US4 are both P2 and can run in parallel after US1/US2
- **Polish (Phase 7)**: Depends on all user stories complete

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Foundational - No dependencies on other stories
- **User Story 2 (P1)**: Can start after Foundational - No dependencies on other stories (can parallel with US1)
- **User Story 3 (P2)**: Can start after Foundational - Optionally after US1 for smoother integration
- **User Story 4 (P2)**: Can start after Foundational - Optionally after US1 for smoother integration

### Within Each User Story (TDD Order)

1. Write tests FIRST - they MUST FAIL
2. Implement minimum code to pass tests
3. Refactor if needed
4. Commit with test commit before implementation commit

### Parallel Opportunities

- T002 and T003 in Setup can run in parallel
- T010-T012 tests can run in parallel, T013 separately
- T020-T022 tests can run in parallel, T023 separately
- All [P] marked tasks can run in parallel with other [P] tasks in same phase
- US1 and US2 can run in parallel (both P1, independent)
- US3 and US4 can run in parallel (both P2, independent)

---

## Parallel Example: User Story 1

```bash
# Write all US1 tests first (in parallel):
Task: "T010 [US1] Write test: Tab from version pane to commits pane"
Task: "T011 [US1] Write test: Tab from commits pane to version pane"
Task: "T012 [US1] Write test: Shift+Tab from commits pane to version pane"

# Then implement (sequential due to dependencies):
Task: "T014 [US1] Handle Tab and Shift+Tab key events"
Task: "T015 [US1] Add RenderCommitListForViewport function"
...
```

---

## Implementation Strategy

### MVP First (User Story 1 + User Story 2)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational
3. Complete Phase 3: User Story 1 (pane navigation)
4. Complete Phase 4: User Story 2 (commit scrolling)
5. **STOP and VALIDATE**: Test dual-pane with scrolling works
6. This is a functional MVP - users can now view all commits

### Incremental Delivery

1. Setup + Foundational → Core model ready
2. Add US1 + US2 → MVP: dual-pane with scrolling
3. Add US3 → Enhanced: commit detail overlay
4. Add US4 → Complete: full workflow validation
5. Polish → Production-ready

### Single Developer Strategy

1. Complete Setup → Foundational → US1 → US2 (MVP)
2. Then US3 → US4 → Polish
3. TDD: Write test → Watch it fail → Implement → Pass → Commit

---

## Notes

- All tests follow TDD per Constitution Principle II
- Tests MUST fail before implementation (Red phase)
- Implementation MUST only make tests pass (Green phase)
- Commit history MUST show test commits before implementation commits
- [P] tasks = different files, no dependencies
- Stop at any checkpoint to validate independently
- Run `golangci-lint run` before final commit per Constitution Principle I
