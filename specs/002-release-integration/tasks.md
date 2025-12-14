# Tasks: Release Tool Integration (Post-Push Hooks)

**Input**: Design documents from `/specs/002-release-integration/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md

**Development Approach**: Test-Driven Development (TDD) - Write tests first, implement to pass, refactor

**Organization**: Tasks are grouped by user story to enable independent implementation and testing.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (US1, US2, US3, US4)
- Include exact file paths in descriptions

## User Stories Summary

| Story | Title | Priority |
|-------|-------|----------|
| US1 | Post-Push Hook Phase | P1 |
| US2 | Post-Push Notifications | P1 |
| US3 | Multiple Post-Push Hooks | P2 |
| US4 | TUI Post-Push Hook Display | P3 |

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Create feature branch and verify existing test infrastructure

- [ ] T001 Create feature branch `002-release-integration` from main
- [ ] T002 Verify all existing tests pass with `go test ./...`
- [ ] T003 Verify linter passes with `golangci-lint run`

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before user stories

**Note**: This phase establishes the `PostPush` hook type and fail-open runner that all stories depend on.

### Tests First (TDD)

- [ ] T004 [P] Write test for `PostPush` HookType constant in `internal/hooks/runner_test.go`
- [ ] T005 [P] Write test for `RunHooksFailOpen` function (continues on failure) in `internal/hooks/runner_test.go`

### Implementation

- [ ] T006 Add `PostPush HookType = "post-push"` constant in `internal/hooks/types.go`
- [ ] T007 Implement `RunHooksFailOpen` function in `internal/hooks/runner.go`
- [ ] T008 Verify tests T004, T005 now pass

**Checkpoint**: Foundation ready - PostPush type and fail-open runner available

---

## Phase 3: User Story 1 - Post-Push Hook Phase (Priority: P1)

**Goal**: Enable running commands after successful tag push

**Independent Test**: Configure post-push hooks in `.bumpkin.yml`, run bump, verify hooks execute after push

### Tests First (TDD)

- [ ] T009 [P] [US1] Write test for config parsing `hooks.post-push` array in `internal/config/config_test.go`
- [ ] T010 [P] [US1] Write test for empty post-push returns empty slice in `internal/config/config_test.go`
- [ ] T011 [P] [US1] Write test for post-push hooks preserve order in `internal/config/config_test.go`
- [ ] T012 [P] [US1] Write test for executor post-push hooks execute after push in `internal/executor/bump_test.go`
- [ ] T013 [P] [US1] Write test for executor skips post-push when `--no-push` in `internal/executor/bump_test.go`
- [ ] T014 [P] [US1] Write test for executor skips post-push when push fails in `internal/executor/bump_test.go`
- [ ] T015 [P] [US1] Write test for post-push failure is warning (tag remains pushed) in `internal/executor/bump_test.go`

### Implementation

- [ ] T016 [US1] Add `PostPush []string` to `Hooks` struct in `internal/config/config.go`
- [ ] T017 [US1] Update `Merge` method to handle post-push hooks in `internal/config/config.go`
- [ ] T018 [US1] Add `PostPushHooks []string` to `Request` struct in `internal/executor/bump.go`
- [ ] T019 [US1] Add `PostPushWarnings []string` to `Result` struct in `internal/executor/bump.go`
- [ ] T020 [US1] Implement post-push hook execution after push in `internal/executor/bump.go`
- [ ] T021 [US1] Verify tests T009-T015 now pass

### CLI Integration

- [ ] T022 [P] [US1] Write test for config post-push hooks passed to executor in `internal/cli/root_test.go`
- [ ] T023 [US1] Pass `PostPushHooks` from config to executor in `internal/cli/root.go`
- [ ] T024 [US1] Verify test T022 passes

**Checkpoint**: User Story 1 complete - post-push hooks execute after successful push

---

## Phase 4: User Story 2 - Post-Push Notifications (Priority: P1)

**Goal**: Enable notifications after tag push with environment variables

**Independent Test**: Configure notification hook with `$BUMPKIN_TAG`, verify variable is substituted

### Tests First (TDD)

- [ ] T025 [P] [US2] Write test for BUMPKIN_* env vars available to post-push hooks in `internal/hooks/runner_test.go`

### Implementation

- [ ] T026 [US2] Verify post-push hooks receive same env vars as pre-tag/post-tag in `internal/executor/bump.go`
- [ ] T027 [US2] Verify test T025 passes

**Checkpoint**: User Story 2 complete - environment variables available in post-push hooks

---

## Phase 5: User Story 3 - Multiple Post-Push Hooks (Priority: P2)

**Goal**: Support multiple post-push hooks with fail-open behavior

**Independent Test**: Configure multiple hooks where one fails, verify all execute and warnings reported

### Tests First (TDD)

- [ ] T028 [P] [US3] Write test for multiple post-push hooks execute in order in `internal/executor/bump_test.go`
- [ ] T029 [P] [US3] Write test for fail-open: remaining hooks run even if one fails in `internal/executor/bump_test.go`

### Implementation

- [ ] T030 [US3] Ensure `RunHooksFailOpen` is used for post-push (not `RunHooks`) in `internal/executor/bump.go`
- [ ] T031 [US3] Collect all warnings from failed hooks in result in `internal/executor/bump.go`
- [ ] T032 [US3] Verify tests T028-T029 pass

**Checkpoint**: User Story 3 complete - multiple hooks with fail-open behavior

---

## Phase 6: User Story 4 - TUI Post-Push Hook Display (Priority: P3)

**Goal**: Display post-push hook results in TUI

**Independent Test**: Run TUI, complete bump with post-push hooks, verify results displayed

### Tests First (TDD)

- [ ] T033 [P] [US4] Write test for TUI model can receive post-push warnings in `internal/tui/model_test.go` (create if needed)

### Implementation

- [ ] T034 [US4] Add `PostPushHooks []string` to TUI Config in `internal/tui/model.go`
- [ ] T035 [US4] Add post-push hook display to result screen in `internal/tui/model.go`
- [ ] T036 [US4] Show warnings (yellow/orange) for failed post-push hooks in `internal/tui/model.go`
- [ ] T037 [US4] Verify test T033 passes

**Checkpoint**: User Story 4 complete - TUI shows post-push results

---

## Phase 7: CLI Enhancements

**Purpose**: Dry-run and --no-hooks support for post-push

### Tests First (TDD)

- [ ] T038 [P] Write test for `--dry-run` shows post-push hooks without executing in `internal/cli/root_test.go`
- [ ] T039 [P] Write test for `--no-hooks` skips post-push hooks in `internal/cli/root_test.go`

### Implementation

- [ ] T040 Update dry-run output to show post-push hooks in `internal/cli/root.go`
- [ ] T041 Verify `--no-hooks` already skips post-push (should work from US1) in `internal/cli/root.go`
- [ ] T042 Verify tests T038-T039 pass

---

## Phase 8: Polish & Cross-Cutting Concerns

**Purpose**: Final verification and cleanup

- [ ] T043 Run all tests with `go test ./...` - ensure all pass
- [ ] T044 Run linter with `golangci-lint run` - fix any issues
- [ ] T045 Apply formatting with `golangci-lint fmt`
- [ ] T046 Manual integration test: create test repo, configure post-push hooks, run bump
- [ ] T047 Update README.md with post-push hook documentation
- [ ] T048 Commit all changes with conventional commit message

---

## Dependencies & Execution Order

### Phase Dependencies

```
Phase 1 (Setup) → Phase 2 (Foundational) → Phase 3-6 (User Stories) → Phase 7-8 (Polish)
                                          ↓
                              Can run in parallel after Phase 2:
                              - US1 (P1) - Core post-push
                              - US2 (P1) - Env vars (depends on US1)
                              - US3 (P2) - Multi-hook (depends on US1)
                              - US4 (P3) - TUI (depends on US1)
```

### User Story Dependencies

| Story | Depends On | Can Parallel With |
|-------|------------|-------------------|
| US1 | Phase 2 (Foundational) | - |
| US2 | US1 (needs hooks working) | - |
| US3 | US1 (needs hooks working) | US2 |
| US4 | US1 (needs hooks working) | US2, US3 |

### Within Each Phase (TDD Flow)

1. Write ALL tests first (marked [P] can run in parallel)
2. Verify tests FAIL
3. Implement code
4. Verify tests PASS
5. Refactor if needed

---

## Parallel Opportunities

### Phase 2: Foundational Tests

```bash
# Can run in parallel:
T004: Test for PostPush constant
T005: Test for RunHooksFailOpen
```

### Phase 3: User Story 1 Tests

```bash
# Can run in parallel:
T009: Test config parsing post-push
T010: Test empty post-push
T011: Test post-push order
T012: Test executor post-push
T013: Test --no-push skips
T014: Test push fail skips
T015: Test warning not error
```

### After US1 Complete: Stories in Parallel

```bash
# Different team members can work on:
Developer A: US2 (T025-T027)
Developer B: US3 (T028-T032)
Developer C: US4 (T033-T037)
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. ✅ Phase 1: Setup (T001-T003)
2. ✅ Phase 2: Foundational (T004-T008)
3. ✅ Phase 3: User Story 1 (T009-T024)
4. **STOP and VALIDATE**: Test post-push hooks work end-to-end
5. Can ship MVP here!

### Full Implementation

Continue with:
- Phase 4: User Story 2 (env vars)
- Phase 5: User Story 3 (multi-hook)
- Phase 6: User Story 4 (TUI)
- Phase 7-8: Polish

---

## Notes

- All tests must be written FIRST and FAIL before implementation
- [P] tasks can run in parallel (different files)
- [US#] label maps task to specific user story
- Commit after each phase or logical group
- Stop at any checkpoint to validate independently
- Total: 48 tasks
