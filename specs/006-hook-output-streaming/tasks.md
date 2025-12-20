# Tasks: Hook Output Streaming

**Input**: Design documents from `/specs/006-hook-output-streaming/`
**Prerequisites**: plan.md (required), spec.md (required), research.md, data-model.md, contracts/

**Tests**: TDD workflow required per project constitution. Tests are written FIRST and MUST FAIL before implementation.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions

- **Single project**: `internal/` for source, tests adjacent to source files
- Paths based on plan.md structure for this Go CLI project

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and type definitions

- [ ] T001 Add StreamType enum (Stdout, Stderr) to internal/hooks/types.go
- [ ] T002 Add OutputLine struct (Text, Stream, Timestamp) to internal/hooks/types.go
- [ ] T003 [P] Add HookLineMsg, HookStartMsg, HookCompleteMsg, HookPhaseCompleteMsg message types to internal/tui/messages.go
- [ ] T004 [P] Add stdout/stderr display styles (StdoutStyle, StderrStyle, HookHeaderStyle) to internal/tui/styles.go

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

### Tests for Foundational Components

- [ ] T005 [P] Write test for OutputBuffer.AddLine() in internal/hooks/buffer_test.go - test adding lines and retrieving them
- [ ] T006 [P] Write test for OutputBuffer.LineCount() in internal/hooks/buffer_test.go - test accurate count tracking
- [ ] T007 [P] Write test for OutputBuffer.MaxLines eviction in internal/hooks/buffer_test.go - test oldest lines are removed when limit exceeded
- [ ] T008 [P] Write test for OutputBuffer.Render() in internal/hooks/buffer_test.go - test formatted output string generation
- [ ] T009 [P] Write test for OutputBuffer thread safety in internal/hooks/buffer_test.go - test concurrent AddLine calls

### Implementation for Foundational Components

- [ ] T010 Implement OutputBuffer struct with mutex in internal/hooks/buffer.go
- [ ] T011 Implement OutputBuffer.AddLine() method in internal/hooks/buffer.go
- [ ] T012 Implement OutputBuffer.LineCount() method in internal/hooks/buffer.go
- [ ] T013 Implement OutputBuffer.Render() method with styled stdout/stderr in internal/hooks/buffer.go
- [ ] T014 Implement OutputBuffer max lines eviction logic in internal/hooks/buffer.go
- [ ] T015 Run `go test ./internal/hooks/... -run Buffer` to verify all buffer tests pass

**Checkpoint**: OutputBuffer foundation ready - user story implementation can now begin

---

## Phase 3: User Story 1 - Real-time Hook Output in TUI (Priority: P1) 🎯 MVP

**Goal**: Display hook stdout/stderr output streaming in real-time within the TUI in a scrollable pane

**Independent Test**: Configure a hook that runs for 10+ seconds with periodic output and verify output appears progressively in the TUI during execution

### Tests for User Story 1

> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [ ] T016 [P] [US1] Write test for RunHookStreaming() basic output capture in internal/hooks/runner_test.go
- [ ] T017 [P] [US1] Write test for RunHookStreaming() stdout/stderr separation in internal/hooks/runner_test.go
- [ ] T018 [P] [US1] Write test for RunHookStreaming() channel closure on completion in internal/hooks/runner_test.go
- [ ] T019 [P] [US1] Write test for HookPane.AddLine() content update in internal/tui/hookpane_test.go
- [ ] T020 [P] [US1] Write test for HookPane.View() rendering with hook header in internal/tui/hookpane_test.go

### Implementation for User Story 1

- [ ] T021 [US1] Add RunHookStreaming() function signature in internal/hooks/runner.go - accepts context, hook, hookCtx, returns (chan OutputLine, chan HookResult)
- [ ] T022 [US1] Implement io.Pipe setup for stdout/stderr in RunHookStreaming() in internal/hooks/runner.go
- [ ] T023 [US1] Implement goroutine readers for stdout/stderr pipes with bufio.Scanner in internal/hooks/runner.go
- [ ] T024 [US1] Implement channel send logic for OutputLine messages in internal/hooks/runner.go
- [ ] T025 [US1] Implement hook completion signaling via done channel in internal/hooks/runner.go
- [ ] T026 [US1] Create HookPane struct wrapping viewport.Model in internal/tui/hookpane.go
- [ ] T027 [US1] Implement NewHookPane(width, height int) constructor in internal/tui/hookpane.go
- [ ] T028 [US1] Implement HookPane.AddLine(line OutputLine) method with auto-scroll in internal/tui/hookpane.go
- [ ] T029 [US1] Implement HookPane.SetCurrentHook(hook Hook, index, total int) for header in internal/tui/hookpane.go
- [ ] T030 [US1] Implement HookPane.View() rendering with header and styled content in internal/tui/hookpane.go
- [ ] T031 [US1] Implement HookPane.Update(msg tea.Msg) for scroll handling in internal/tui/hookpane.go
- [ ] T032 [US1] Add StateExecutingHooks constant to internal/tui/model.go
- [ ] T033 [US1] Add hookPane, hookOutputChan, hookDoneChan fields to Model struct in internal/tui/model.go
- [ ] T034 [US1] Implement waitForHookLine(chan OutputLine) tea.Cmd function in internal/tui/model.go
- [ ] T035 [US1] Implement waitForHookDone(chan HookResult) tea.Cmd function in internal/tui/model.go
- [ ] T036 [US1] Handle HookLineMsg in Update() - append to hookPane, re-queue listener in internal/tui/model.go
- [ ] T037 [US1] Handle HookStartMsg in Update() - update hookPane header in internal/tui/model.go
- [ ] T038 [US1] Handle HookCompleteMsg in Update() - transition state or start next hook in internal/tui/model.go
- [ ] T039 [US1] Add StateExecutingHooks case to View() - render hookPane in internal/tui/model.go
- [ ] T040 [US1] Add StateExecutingHooks case to handleKeyPress() - delegate scroll to hookPane in internal/tui/model.go
- [ ] T041 [US1] Modify executeVersion() to return tea.Cmd that starts streaming hooks in internal/tui/model.go
- [ ] T042 [US1] Handle WindowSizeMsg for hookPane resizing in internal/tui/model.go
- [ ] T043 [US1] Run `go test ./internal/hooks/... ./internal/tui/... -v` to verify all US1 tests pass
- [ ] T044 [US1] Run `golangci-lint run` to verify code quality

**Checkpoint**: At this point, User Story 1 should be fully functional - hook output streams in TUI with scrolling

---

## Phase 4: User Story 2 - Non-Interactive Output Streaming (Priority: P2)

**Goal**: Stream hook output to terminal with proper line buffering for CI/CD pipelines

**Independent Test**: Run bumpkin in non-interactive mode with a hook that outputs lines with delays, verify each line appears immediately after newline

### Tests for User Story 2

- [ ] T045 [P] [US2] Write test for line-buffered output in non-interactive mode in internal/hooks/runner_test.go
- [ ] T046 [P] [US2] Write test for stdout/stderr preservation in non-interactive mode in internal/hooks/runner_test.go

### Implementation for User Story 2

- [ ] T047 [US2] Add streaming bool parameter or mode detection to runner functions in internal/hooks/runner.go
- [ ] T048 [US2] Implement line-buffered writes to os.Stdout/os.Stderr for non-TUI mode in internal/hooks/runner.go
- [ ] T049 [US2] Add prefix markers ("[stdout]", "[stderr]") for non-TUI mode output in internal/hooks/runner.go
- [ ] T050 [US2] Ensure existing RunHook() behavior preserved when TUI not active in internal/hooks/runner.go
- [ ] T051 [US2] Run `go test ./internal/hooks/... -run NonInteractive` to verify US2 tests pass

**Checkpoint**: Non-interactive mode streams output correctly with line buffering

---

## Phase 5: User Story 3 - Hook Output State Persistence (Priority: P3)

**Goal**: Retain complete output history from all hooks for post-execution review

**Independent Test**: Run hooks that succeed and fail, verify output pane retains complete output with scrollback after execution

### Tests for User Story 3

- [ ] T052 [P] [US3] Write test for OutputBuffer retaining output after hook failure in internal/hooks/buffer_test.go
- [ ] T053 [P] [US3] Write test for OutputBuffer multi-hook separation in internal/hooks/buffer_test.go
- [ ] T054 [P] [US3] Write test for HookPane scrollback after completion in internal/tui/hookpane_test.go

### Implementation for User Story 3

- [ ] T055 [US3] Add StartHook(hook Hook) method to OutputBuffer for hook separation headers in internal/hooks/buffer.go
- [ ] T056 [US3] Add EndHook(hook Hook, success bool, err error) method to OutputBuffer in internal/hooks/buffer.go
- [ ] T057 [US3] Update Render() to include hook separation headers in output in internal/hooks/buffer.go
- [ ] T058 [US3] Ensure hookPane retains OutputBuffer reference after StateExecutingHooks exits in internal/tui/model.go
- [ ] T059 [US3] Allow hookPane scrolling in StateDone state in internal/tui/model.go
- [ ] T060 [US3] Add help text for scrolling in StateDone when hooks executed in internal/tui/model.go
- [ ] T061 [US3] Run `go test ./internal/hooks/... ./internal/tui/... -run Persist` to verify US3 tests pass

**Checkpoint**: All user stories should now be independently functional

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect multiple user stories

- [ ] T062 [P] Handle edge case: hooks with no output - display "No output" message in internal/tui/hookpane.go
- [ ] T063 [P] Handle edge case: very long lines (>10,000 chars) - truncate with indicator in internal/hooks/buffer.go
- [ ] T064 [P] Handle edge case: non-printable characters - escape or replace in internal/hooks/buffer.go
- [ ] T065 [P] Add elapsed time indicator during hook execution in internal/tui/hookpane.go
- [ ] T066 Run full test suite: `go test ./... -v`
- [ ] T067 Run linter: `golangci-lint run`
- [ ] T068 Manual test: Configure long-running hook and verify real-time streaming
- [ ] T069 Manual test: Verify scroll works during and after hook execution
- [ ] T070 Manual test: Verify non-interactive mode with piped output

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
- **User Stories (Phase 3+)**: All depend on Foundational phase completion
  - User stories can then proceed in priority order (P1 → P2 → P3)
- **Polish (Final Phase)**: Depends on all desired user stories being complete

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Foundational (Phase 2) - No dependencies on other stories
- **User Story 2 (P2)**: Can start after Foundational (Phase 2) - Independent of US1, uses same buffer
- **User Story 3 (P3)**: Can start after Foundational (Phase 2) - Enhances US1 buffer, but independently testable

### Within Each User Story

- Tests MUST be written and FAIL before implementation (TDD per constitution)
- Types/structs before methods
- Core implementation before integration
- Story complete before moving to next priority

### Parallel Opportunities

- All Setup tasks marked [P] can run in parallel (T003, T004)
- All Foundational tests marked [P] can run in parallel (T005-T009)
- All tests for a user story marked [P] can run in parallel
- Different developers could work on US2 and US3 in parallel after US1 core is done

---

## Parallel Example: Foundational Tests

```bash
# Launch all buffer tests together:
Task: "Write test for OutputBuffer.AddLine()"
Task: "Write test for OutputBuffer.LineCount()"
Task: "Write test for OutputBuffer.MaxLines eviction"
Task: "Write test for OutputBuffer.Render()"
Task: "Write test for OutputBuffer thread safety"
```

## Parallel Example: User Story 1 Tests

```bash
# Launch all US1 tests together:
Task: "Write test for RunHookStreaming() basic output capture"
Task: "Write test for RunHookStreaming() stdout/stderr separation"
Task: "Write test for RunHookStreaming() channel closure"
Task: "Write test for HookPane.AddLine() content update"
Task: "Write test for HookPane.View() rendering"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup (T001-T004)
2. Complete Phase 2: Foundational (T005-T015)
3. Complete Phase 3: User Story 1 (T016-T044)
4. **STOP and VALIDATE**: Test hook output streaming in TUI
5. Deploy/demo if ready - core value delivered

### Incremental Delivery

1. Complete Setup + Foundational → OutputBuffer ready
2. Add User Story 1 → Real-time TUI streaming works (MVP!)
3. Add User Story 2 → Non-interactive mode enhanced
4. Add User Story 3 → Post-execution review enabled
5. Each story adds value without breaking previous stories

### TDD Cycle Reminder

For each task group:
1. **Red**: Write test, run it, confirm it FAILS
2. **Green**: Implement minimum code to pass test
3. **Refactor**: Clean up while keeping tests green
4. **Commit**: After each logical group of tests pass

---

## Notes

- [P] tasks = different files, no dependencies
- [Story] label maps task to specific user story for traceability
- Each user story should be independently completable and testable
- Verify tests fail before implementing (TDD)
- Commit after each task or logical group
- Stop at any checkpoint to validate story independently
- All paths are relative to repository root `/Users/benny/Documents/Projects/golang/bumpkin/`
