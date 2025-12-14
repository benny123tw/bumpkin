# Tasks: Basic Subcommands

**Input**: Design documents from `/specs/004-basic-commands/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, quickstart.md
**Methodology**: Test-Driven Development (TDD)

**TDD Workflow**: For each user story:
1. Write failing tests first (Red)
2. Implement minimal code to pass tests (Green)
3. Refactor while keeping tests green (Refactor)

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

### TDD for User Story 1

- [x] T001 [US1] Write failing test for version subcommand in internal/cli/version_test.go
- [x] T002 [US1] Create version subcommand to pass tests in internal/cli/version.go
- [x] T003 [US1] Verify `bumpkin version` output matches `bumpkin --version` (add test if needed)

**Checkpoint**: User Story 1 complete - `bumpkin version` works with passing tests.

---

## Phase 4: User Story 2 - Help Subcommand (Priority: P1)

**Goal**: Allow users to run `bumpkin help` and `bumpkin help <subcommand>` to see usage information.

**Independent Test**: Run `bumpkin help` and verify it shows same output as `bumpkin --help`.

### TDD for User Story 2

- [x] T004 [US2] Write test to verify Cobra's built-in help subcommand in internal/cli/help_test.go
- [x] T005 [US2] Verify `bumpkin help version` shows version command help (test)
- [x] T006 [US2] Verify tests pass (Cobra provides help automatically - no code needed)

**Checkpoint**: User Story 2 complete - `bumpkin help` and `bumpkin help <cmd>` work with passing tests.

---

## Phase 5: User Story 3 - Init Subcommand (Priority: P2)

**Goal**: Allow users to run `bumpkin init` to create a starter `.bumpkin.yaml` configuration file.

**Independent Test**: Run `bumpkin init` in a directory without config and verify `.bumpkin.yaml` is created.

### TDD for User Story 3

- [x] T007 [US3] Write failing tests for init subcommand in internal/cli/init_test.go
- [x] T008 [US3] Create init subcommand with config template in internal/cli/init.go
- [x] T009 [US3] Write test for error when config exists, implement error handling
- [x] T010 [US3] Verify all init tests pass and refactor

**Checkpoint**: User Story 3 complete - `bumpkin init` creates config file with passing tests.

---

## Phase 6: User Story 4 - Current Subcommand (Priority: P2)

**Goal**: Allow users to run `bumpkin current` to see the latest version tag.

**Independent Test**: Run `bumpkin current` in a repo with tags and verify it shows the latest version.

### TDD for User Story 4

- [x] T011 [P] [US4] Write failing tests for current subcommand in internal/cli/current_test.go
- [x] T012 [US4] Create current subcommand in internal/cli/current.go
- [x] T013 [US4] Write test for --prefix flag, implement flag support
- [x] T014 [US4] Write tests for edge cases (no repo, no tags), implement handling

**Checkpoint**: User Story 4 complete - `bumpkin current` shows latest tag with passing tests.

---

## Phase 7: User Story 5 - Completion Subcommand (Priority: P3)

**Goal**: Allow users to run `bumpkin completion <shell>` to generate shell completion scripts.

**Independent Test**: Run `bumpkin completion bash` and verify it outputs a valid bash completion script.

### TDD for User Story 5

- [x] T015 [P] [US5] Write failing tests for completion subcommand in internal/cli/completion_test.go
- [x] T016 [US5] Create completion subcommand structure in internal/cli/completion.go
- [x] T017 [US5] Write test for bash completion, implement using Cobra
- [x] T018 [US5] Write test for zsh completion, implement using Cobra
- [x] T019 [US5] Write test for fish completion, implement using Cobra
- [x] T020 [US5] Write test for powershell completion, implement using Cobra
- [x] T021 [US5] Write test for missing shell arg, add help text with usage instructions

**Checkpoint**: User Story 5 complete - `bumpkin completion <shell>` generates scripts with passing tests.

---

## Phase 8: Polish & Cross-Cutting Concerns

**Purpose**: Final validation and cleanup

- [x] T022 Run golangci-lint and fix any issues
- [x] T023 Run full test suite: `go test ./...`
- [x] T024 Verify all commands appear in `bumpkin --help` output
- [x] T025 Run quickstart.md validation checklist

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

### TDD Workflow Within Each User Story

1. Write failing test (Red)
2. Implement minimal code to pass (Green)
3. Refactor while keeping tests green
4. Run golangci-lint before moving to next task
5. Story complete when all tests pass

### Parallel Opportunities

- US1 and US2 are P1 but US2 depends on US1 for full testing
- US3 and US4 are both P2 and can run in parallel after US1/US2
- T011 and T015 are marked [P] - different files, can be started together

---

## Parallel Example: User Stories 3 and 4

```bash
# After completing US1 and US2, start US3 and US4 in parallel:
Task: "Write failing tests for init subcommand" [US3]
Task: "Write failing tests for current subcommand" [US4]
```

---

## Implementation Strategy

### TDD MVP First (User Stories 1 + 2)

1. Write tests for User Story 1 → Implement → Verify green
2. Write tests for User Story 2 → Verify (Cobra built-in) → Green
3. **STOP and VALIDATE**: All tests pass, both commands work
4. Deploy/release if ready

### Incremental TDD Delivery

1. US1 + US2 → MVP ready (version + help) with tests
2. Add US3 (init) → Tests first → Config generation available
3. Add US4 (current) → Tests first → Version checking for scripts
4. Add US5 (completion) → Tests first → Power user feature
5. Each story adds value with full test coverage

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

- **TDD**: Tests are written BEFORE implementation for each feature
- Cobra provides `help` automatically - US2 is mostly test verification
- All subcommands are independent files in `internal/cli/`
- Each command file has a corresponding `_test.go` file
- Config template for `init` uses existing `Config` struct pattern
- Completion uses Cobra's built-in generation functions
- All code must pass golangci-lint before merge (Constitution requirement)
