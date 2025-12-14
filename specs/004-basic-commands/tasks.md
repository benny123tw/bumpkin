# Tasks: Basic Subcommands

**Input**: Design documents from `/specs/004-basic-commands/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, quickstart.md

**Tests**: Not explicitly requested - tests will be added for key functionality.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3, US4, US5)
- Include exact file paths in descriptions

## Path Conventions

- **Single project**: `internal/cli/` for CLI commands
- All new subcommands go in `internal/cli/`

---

## Phase 1: Setup

**Purpose**: No setup needed - using existing project structure

This feature builds on the existing codebase. No new dependencies or project structure changes required.

**Checkpoint**: Ready to proceed to user story implementation.

---

## Phase 2: Foundational

**Purpose**: No foundational blocking work - Cobra already supports subcommands natively.

All user stories are independent and can proceed directly.

**Checkpoint**: Foundation ready - user story implementation can begin.

---

## Phase 3: User Story 1 - Version Subcommand (Priority: P1) 🎯 MVP

**Goal**: Allow users to run `bumpkin version` to see version information.

**Independent Test**: Run `bumpkin version` and verify it displays version, commit, and build date (same as `bumpkin --version`).

### Implementation for User Story 1

- [ ] T001 [US1] Create version subcommand in internal/cli/version.go
- [ ] T002 [US1] Add unit test for version command in internal/cli/version_test.go
- [ ] T003 [US1] Verify `bumpkin version` output matches `bumpkin --version`

**Checkpoint**: User Story 1 complete - `bumpkin version` works.

---

## Phase 4: User Story 2 - Help Subcommand (Priority: P1)

**Goal**: Allow users to run `bumpkin help` and `bumpkin help <subcommand>` to see usage information.

**Independent Test**: Run `bumpkin help` and verify it shows same output as `bumpkin --help`.

### Implementation for User Story 2

- [ ] T004 [US2] Verify Cobra's built-in help subcommand works (no code needed)
- [ ] T005 [US2] Verify `bumpkin help version` shows version command help
- [ ] T006 [US2] Add test to confirm help subcommand behavior in internal/cli/help_test.go

**Checkpoint**: User Story 2 complete - `bumpkin help` and `bumpkin help <cmd>` work.

---

## Phase 5: User Story 3 - Init Subcommand (Priority: P2)

**Goal**: Allow users to run `bumpkin init` to create a starter `.bumpkin.yaml` configuration file.

**Independent Test**: Run `bumpkin init` in a directory without config and verify `.bumpkin.yaml` is created.

### Implementation for User Story 3

- [ ] T007 [US3] Create init subcommand in internal/cli/init.go
- [ ] T008 [US3] Define config template with comments in internal/cli/init.go
- [ ] T009 [US3] Add error handling for existing config file
- [ ] T010 [US3] Add unit tests for init command in internal/cli/init_test.go

**Checkpoint**: User Story 3 complete - `bumpkin init` creates config file.

---

## Phase 6: User Story 4 - Current Subcommand (Priority: P2)

**Goal**: Allow users to run `bumpkin current` to see the latest version tag.

**Independent Test**: Run `bumpkin current` in a repo with tags and verify it shows the latest version.

### Implementation for User Story 4

- [ ] T011 [P] [US4] Create current subcommand in internal/cli/current.go
- [ ] T012 [US4] Add --prefix flag support for tag filtering
- [ ] T013 [US4] Handle edge cases (no repo, no tags)
- [ ] T014 [US4] Add unit tests for current command in internal/cli/current_test.go

**Checkpoint**: User Story 4 complete - `bumpkin current` shows latest tag.

---

## Phase 7: User Story 5 - Completion Subcommand (Priority: P3)

**Goal**: Allow users to run `bumpkin completion <shell>` to generate shell completion scripts.

**Independent Test**: Run `bumpkin completion bash` and verify it outputs a valid bash completion script.

### Implementation for User Story 5

- [ ] T015 [P] [US5] Create completion subcommand in internal/cli/completion.go
- [ ] T016 [US5] Implement bash completion generation using Cobra
- [ ] T017 [US5] Implement zsh completion generation using Cobra
- [ ] T018 [US5] Implement fish completion generation using Cobra
- [ ] T019 [US5] Implement powershell completion generation using Cobra
- [ ] T020 [US5] Add help text with usage instructions for each shell
- [ ] T021 [US5] Add unit tests for completion command in internal/cli/completion_test.go

**Checkpoint**: User Story 5 complete - `bumpkin completion <shell>` generates scripts for all shells.

---

## Phase 8: Polish & Cross-Cutting Concerns

**Purpose**: Final validation and cleanup

- [ ] T022 Run golangci-lint and fix any issues
- [ ] T023 Run full test suite: `go test ./...`
- [ ] T024 Verify all commands appear in `bumpkin --help` output
- [ ] T025 Run quickstart.md validation checklist

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: N/A - existing project
- **Foundational (Phase 2)**: N/A - Cobra supports subcommands natively
- **User Stories (Phase 3-7)**: Can proceed in priority order or in parallel
- **Polish (Phase 8)**: Depends on all user stories being complete

### User Story Dependencies

- **User Story 1 (P1)**: No dependencies - can start immediately
- **User Story 2 (P1)**: Depends on US1 (version command must exist to test `help version`)
- **User Story 3 (P2)**: No dependencies - can start after US1/US2 or in parallel
- **User Story 4 (P2)**: No dependencies - can start after US1/US2 or in parallel
- **User Story 5 (P3)**: Depends on all other commands being registered (for complete completion)

### Within Each User Story

- Implementation before tests (tests verify behavior)
- Core functionality before edge cases
- Story complete before moving to next priority

### Parallel Opportunities

- US1 and US2 are P1 but US2 depends on US1 for full testing
- US3 and US4 are both P2 and can run in parallel after US1/US2
- T011 and T015 are marked [P] - different files, can be started together
- All test tasks can run after their respective implementations

---

## Parallel Example: User Stories 3 and 4

```bash
# After completing US1 and US2, start US3 and US4 in parallel:
Task: "Create init subcommand in internal/cli/init.go" [US3]
Task: "Create current subcommand in internal/cli/current.go" [US4]
```

---

## Implementation Strategy

### MVP First (User Stories 1 + 2)

1. Complete User Story 1: version subcommand
2. Complete User Story 2: help subcommand (mostly verification)
3. **STOP and VALIDATE**: Test both commands
4. Deploy/release if ready

### Incremental Delivery

1. US1 + US2 → MVP ready (version + help)
2. Add US3 (init) → Config generation available
3. Add US4 (current) → Version checking for scripts
4. Add US5 (completion) → Power user feature
5. Each story adds value without breaking previous stories

---

## Summary

| Metric | Count |
|--------|-------|
| Total tasks | 25 |
| User Story 1 (version) | 3 tasks |
| User Story 2 (help) | 3 tasks |
| User Story 3 (init) | 4 tasks |
| User Story 4 (current) | 4 tasks |
| User Story 5 (completion) | 7 tasks |
| Polish | 4 tasks |
| Parallel opportunities | US3+US4 can run together, T011+T015 can start together |
| MVP scope | US1 + US2 (6 tasks) |

---

## Notes

- Cobra provides `help` automatically - US2 is mostly verification
- All subcommands are independent files in `internal/cli/`
- Config template for `init` uses existing `Config` struct pattern
- Completion uses Cobra's built-in generation functions
- All code must pass golangci-lint before merge (Constitution requirement)
