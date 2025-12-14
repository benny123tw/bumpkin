# Tasks: Version Tagger CLI

**Input**: Design documents from `/specs/001-version-tagger/`
**Prerequisites**: plan.md (required), spec.md (required), data-model.md, contracts/cli.md

**TDD Approach**: Tests are written FIRST, must FAIL before implementation, then implementation makes tests pass.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions

- **Single project**: `cmd/bumpkin/`, `internal/` at repository root
- Tests alongside source: `internal/version/semver_test.go` for `internal/version/semver.go`

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and Go module setup

- [x] T001 Initialize Go module with `go mod init github.com/benny123tw/bumpkin`
- [x] T002 Create project directory structure per plan.md (cmd/, internal/ with subdirectories)
- [x] T003 [P] Add core dependencies to go.mod (bubbletea, bubbles, lipgloss, cobra, go-git, semver, conventionalcommits)
- [x] T004 [P] Create minimal main.go entry point in cmd/bumpkin/main.go
- [x] T005 [P] Verify golangci-lint runs successfully on empty project

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

### Version Entity (shared by all stories)

- [x] T006 [P] Write test for Version struct creation and string formatting in internal/version/semver_test.go
- [x] T007 [P] Write test for parsing version strings (with/without v prefix) in internal/version/semver_test.go
- [x] T008 Implement Version struct and Parse function in internal/version/semver.go (make T006, T007 pass)
- [x] T009 Write test for version comparison (LessThan, Equal) in internal/version/semver_test.go
- [x] T010 Implement version comparison methods in internal/version/semver.go (make T009 pass)

### Bump Types Entity

- [x] T011 [P] Write test for BumpType enum and string representation in internal/version/bump_test.go
- [x] T012 Implement BumpType constants and String method in internal/version/bump.go (make T011 pass)

### Git Repository Detection

- [x] T013 Write test for detecting git repository in internal/git/repository_test.go
- [x] T014 Implement Repository.Open function in internal/git/repository.go (make T013 pass)
- [x] T015 Write test for repository not found error in internal/git/repository_test.go
- [x] T016 Implement error handling for non-git directories in internal/git/repository.go (make T015 pass)

**Checkpoint**: Foundation ready - user story implementation can now begin

---

## Phase 3: User Story 1 - Interactive Version Bump (Priority: P1) 🎯 MVP

**Goal**: Core TUI for viewing commits, selecting version, and tagging

**Independent Test**: Run CLI in git repo with tags, select version, verify tag created

### Tests for US1 - Version Bumping Logic

- [x] T017 [P] [US1] Write test for BumpPatch (1.2.3 → 1.2.4) in internal/version/bump_test.go
- [x] T018 [P] [US1] Write test for BumpMinor (1.2.3 → 1.3.0) in internal/version/bump_test.go
- [x] T019 [P] [US1] Write test for BumpMajor (1.2.3 → 2.0.0) in internal/version/bump_test.go
- [x] T020 [US1] Implement BumpPatch, BumpMinor, BumpMajor in internal/version/bump.go (make T017-T019 pass)

### Tests for US1 - Git Tag Operations

- [ ] T021 [P] [US1] Write test for listing all tags in internal/git/tags_test.go
- [ ] T022 [P] [US1] Write test for finding latest semver tag in internal/git/tags_test.go
- [ ] T023 [US1] Implement ListTags and LatestTag in internal/git/tags.go (make T021-T022 pass)
- [ ] T024 [P] [US1] Write test for creating annotated tag in internal/git/tags_test.go
- [ ] T025 [US1] Implement CreateTag in internal/git/tags.go (make T024 pass)

### Tests for US1 - Git Commit Operations

- [ ] T026 [P] [US1] Write test for listing commits since tag in internal/git/commits_test.go
- [ ] T027 [P] [US1] Write test for Commit struct with hash, message, author in internal/git/commits_test.go
- [ ] T028 [US1] Implement GetCommitsSinceTag in internal/git/commits.go (make T026-T027 pass)
- [ ] T029 [P] [US1] Write test for empty commit list when no commits since tag in internal/git/commits_test.go
- [ ] T030 [US1] Handle edge case: no commits since tag in internal/git/commits.go (make T029 pass)

### Tests for US1 - Git Push Operations

- [ ] T031 [P] [US1] Write test for push to remote in internal/git/push_test.go
- [ ] T032 [US1] Implement PushTag in internal/git/push.go (make T031 pass)

### Tests for US1 - Executor (shared bump logic)

- [ ] T033 [P] [US1] Write test for executor with patch bump in internal/executor/bump_test.go
- [ ] T034 [P] [US1] Write test for executor with dry-run mode in internal/executor/bump_test.go
- [ ] T035 [US1] Implement Execute function in internal/executor/bump.go (make T033-T034 pass)
- [ ] T036 [P] [US1] Write test for executor with no-push mode in internal/executor/bump_test.go
- [ ] T037 [US1] Add NoPush support to executor in internal/executor/bump.go (make T036 pass)

### Implementation for US1 - TUI Components

- [ ] T038 [P] [US1] Write test for TUI model state transitions in internal/tui/model_test.go
- [ ] T039 [US1] Implement base Model struct and Init in internal/tui/model.go (make T038 pass)
- [ ] T040 [P] [US1] Define TUI message types in internal/tui/messages.go
- [ ] T041 [P] [US1] Create Lipgloss styles in internal/tui/styles.go
- [ ] T042 [P] [US1] Write test for commit list view rendering in internal/tui/commits_test.go
- [ ] T043 [US1] Implement commit list View in internal/tui/commits.go (make T042 pass)
- [ ] T044 [P] [US1] Write test for version selector view in internal/tui/selector_test.go
- [ ] T045 [US1] Implement version selector View in internal/tui/selector.go (make T044 pass)
- [ ] T046 [P] [US1] Write test for confirmation view in internal/tui/confirm_test.go
- [ ] T047 [US1] Implement confirmation View in internal/tui/confirm.go (make T046 pass)
- [ ] T048 [US1] Implement Update function handling all state transitions in internal/tui/model.go
- [ ] T049 [US1] Implement main View function composing all views in internal/tui/model.go

### Integration for US1

- [ ] T050 [US1] Write integration test for full interactive bump flow in internal/tui/integration_test.go
- [ ] T051 [US1] Wire TUI to executor and verify end-to-end in internal/tui/model.go

**Checkpoint**: User Story 1 complete - interactive version bump works independently

---

## Phase 4: User Story 2 - Non-Interactive Mode (Priority: P2)

**Goal**: CLI flags for automation without TUI

**Independent Test**: Run `bumpkin --patch` and verify tag created without prompts

### Tests for US2 - CLI Root Command

- [ ] T052 [P] [US2] Write test for root command creation in internal/cli/root_test.go
- [ ] T053 [US2] Implement root command with Cobra in internal/cli/root.go (make T052 pass)

### Tests for US2 - CLI Flags

- [ ] T054 [P] [US2] Write test for --patch flag in internal/cli/flags_test.go
- [ ] T055 [P] [US2] Write test for --minor flag in internal/cli/flags_test.go
- [ ] T056 [P] [US2] Write test for --major flag in internal/cli/flags_test.go
- [ ] T057 [P] [US2] Write test for --version custom flag in internal/cli/flags_test.go
- [ ] T058 [US2] Implement version bump flags in internal/cli/flags.go (make T054-T057 pass)
- [ ] T059 [P] [US2] Write test for --dry-run flag in internal/cli/flags_test.go
- [ ] T060 [P] [US2] Write test for --no-push flag in internal/cli/flags_test.go
- [ ] T061 [P] [US2] Write test for --yes flag in internal/cli/flags_test.go
- [ ] T062 [US2] Implement behavior flags in internal/cli/flags.go (make T059-T061 pass)
- [ ] T063 [P] [US2] Write test for --json output flag in internal/cli/flags_test.go
- [ ] T064 [US2] Implement JSON output in internal/cli/root.go (make T063 pass)

### Tests for US2 - Mode Selection Logic

- [ ] T065 [P] [US2] Write test: flags present → non-interactive mode in internal/cli/root_test.go
- [ ] T066 [P] [US2] Write test: no flags → interactive mode in internal/cli/root_test.go
- [ ] T067 [US2] Implement mode selection in root command RunE in internal/cli/root.go (make T065-T066 pass)

### Integration for US2

- [ ] T068 [US2] Write integration test for `bumpkin --patch --yes` in internal/cli/integration_test.go
- [ ] T069 [US2] Wire main.go to execute root command in cmd/bumpkin/main.go

**Checkpoint**: User Story 2 complete - non-interactive mode works independently

---

## Phase 5: User Story 3 - Conventional Commit Analysis (Priority: P3)

**Goal**: Automatically suggest version bump based on commit types

**Independent Test**: Create repo with "feat:" commits, verify minor bump suggested

### Tests for US3 - Commit Parsing

- [ ] T070 [P] [US3] Write test for parsing "feat:" commit in internal/conventional/parser_test.go
- [ ] T071 [P] [US3] Write test for parsing "fix:" commit in internal/conventional/parser_test.go
- [ ] T072 [P] [US3] Write test for parsing "feat!:" breaking change in internal/conventional/parser_test.go
- [ ] T073 [P] [US3] Write test for parsing "BREAKING CHANGE:" footer in internal/conventional/parser_test.go
- [ ] T074 [US3] Implement ParseCommit function in internal/conventional/parser.go (make T070-T073 pass)
- [ ] T075 [P] [US3] Write test for parsing commit with scope in internal/conventional/parser_test.go
- [ ] T076 [US3] Add scope parsing to ParseCommit in internal/conventional/parser.go (make T075 pass)

### Tests for US3 - Bump Recommendation

- [ ] T077 [P] [US3] Write test: feat commits → recommend minor in internal/conventional/analyzer_test.go
- [ ] T078 [P] [US3] Write test: fix only commits → recommend patch in internal/conventional/analyzer_test.go
- [ ] T079 [P] [US3] Write test: breaking change → recommend major in internal/conventional/analyzer_test.go
- [ ] T080 [US3] Implement AnalyzeCommits function in internal/conventional/analyzer.go (make T077-T079 pass)
- [ ] T081 [P] [US3] Write test: mixed commits use highest priority in internal/conventional/analyzer_test.go
- [ ] T082 [US3] Implement priority logic (major > minor > patch) in internal/conventional/analyzer.go (make T081 pass)

### Integration for US3

- [ ] T083 [US3] Integrate analyzer with TUI to show recommendation in internal/tui/selector.go
- [ ] T084 [US3] Add --conventional flag for CLI mode in internal/cli/flags.go

**Checkpoint**: User Story 3 complete - conventional commit analysis works

---

## Phase 6: User Story 4 - Hook System (Priority: P4)

**Goal**: Run custom scripts before/after tagging

**Independent Test**: Configure pre-tag hook, verify it runs before tag creation

### Tests for US4 - Configuration Loading

- [ ] T085 [P] [US4] Write test for loading .bumpkin.yml in internal/config/config_test.go
- [ ] T086 [P] [US4] Write test for default config when file missing in internal/config/config_test.go
- [ ] T087 [US4] Implement Load function in internal/config/config.go (make T085-T086 pass)
- [ ] T088 [P] [US4] Write test for config with hooks defined in internal/config/config_test.go
- [ ] T089 [US4] Parse hook arrays in config in internal/config/config.go (make T088 pass)

### Tests for US4 - Hook Types

- [ ] T090 [P] [US4] Write test for Hook struct in internal/hooks/types_test.go
- [ ] T091 [US4] Implement Hook and HookResult structs in internal/hooks/types.go (make T090 pass)

### Tests for US4 - Hook Execution

- [ ] T092 [P] [US4] Write test for running single hook command in internal/hooks/runner_test.go
- [ ] T093 [US4] Implement RunHook function in internal/hooks/runner.go (make T092 pass)
- [ ] T094 [P] [US4] Write test for hook environment variables in internal/hooks/runner_test.go
- [ ] T095 [US4] Pass environment variables to hooks in internal/hooks/runner.go (make T094 pass)
- [ ] T096 [P] [US4] Write test for hook failure (non-zero exit) in internal/hooks/runner_test.go
- [ ] T097 [US4] Return error on hook failure in internal/hooks/runner.go (make T096 pass)
- [ ] T098 [P] [US4] Write test for running multiple hooks in sequence in internal/hooks/runner_test.go
- [ ] T099 [US4] Implement RunHooks for multiple commands in internal/hooks/runner.go (make T098 pass)

### Integration for US4

- [ ] T100 [US4] Integrate pre-tag hooks into executor (abort on failure) in internal/executor/bump.go
- [ ] T101 [US4] Integrate post-tag hooks into executor in internal/executor/bump.go
- [ ] T102 [US4] Add --no-hooks flag support in internal/cli/flags.go
- [ ] T103 [US4] Write integration test with actual hook script in internal/hooks/integration_test.go

**Checkpoint**: User Story 4 complete - hook system works

---

## Phase 7: User Story 5 - Prerelease Versions (Priority: P5)

**Goal**: Support alpha, beta, rc prerelease versions

**Independent Test**: Bump to alpha, verify correct prerelease format

### Tests for US5 - Prerelease Parsing

- [ ] T104 [P] [US5] Write test for parsing prerelease version in internal/version/prerelease_test.go
- [ ] T105 [P] [US5] Write test for extracting prerelease type and number in internal/version/prerelease_test.go
- [ ] T106 [US5] Implement prerelease parsing in internal/version/prerelease.go (make T104-T105 pass)

### Tests for US5 - Prerelease Bumping

- [ ] T107 [P] [US5] Write test: v1.0.0 → v1.0.1-alpha.0 in internal/version/prerelease_test.go
- [ ] T108 [P] [US5] Write test: v1.0.1-alpha.0 → v1.0.1-alpha.1 in internal/version/prerelease_test.go
- [ ] T109 [P] [US5] Write test: v1.0.1-alpha.1 → v1.0.1-beta.0 in internal/version/prerelease_test.go
- [ ] T110 [US5] Implement BumpPrerelease in internal/version/prerelease.go (make T107-T109 pass)
- [ ] T111 [P] [US5] Write test: v1.0.1-rc.0 → v1.0.1 (release) in internal/version/prerelease_test.go
- [ ] T112 [US5] Implement Release (strip prerelease) in internal/version/prerelease.go (make T111 pass)

### Integration for US5

- [ ] T113 [US5] Add prerelease options to TUI selector in internal/tui/selector.go
- [ ] T114 [US5] Add --prerelease flag to CLI in internal/cli/flags.go
- [ ] T115 [US5] Add --release flag to CLI in internal/cli/flags.go
- [ ] T116 [US5] Wire prerelease to executor in internal/executor/bump.go

**Checkpoint**: User Story 5 complete - prerelease versions work

---

## Phase 8: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect multiple user stories

### Error Handling & Edge Cases

- [ ] T117 [P] Write test for no tags in repo (start from v0.0.0) in internal/git/tags_test.go
- [ ] T118 Handle no existing tags gracefully in internal/git/tags.go (make T117 pass)
- [ ] T119 [P] Write test for custom version validation in internal/version/semver_test.go
- [ ] T120 Implement version validation for custom input in internal/version/semver.go (make T119 pass)
- [ ] T121 [P] Write test for exit codes in internal/cli/root_test.go
- [ ] T122 Implement proper exit codes per CLI contract in internal/cli/root.go (make T121 pass)

### Configuration Flags

- [ ] T123 [P] Write test for --remote flag in internal/cli/flags_test.go
- [ ] T124 [P] Write test for --prefix flag in internal/cli/flags_test.go
- [ ] T125 [P] Write test for --config flag in internal/cli/flags_test.go
- [ ] T126 Implement configuration flags in internal/cli/flags.go (make T123-T125 pass)

### Final Integration

- [ ] T127 Run full end-to-end test: interactive bump in test repository
- [ ] T128 Run full end-to-end test: non-interactive bump with all flags
- [ ] T129 Run golangci-lint on entire codebase and fix any issues
- [ ] T130 Verify quickstart.md examples work correctly

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
- **User Stories (Phase 3+)**: All depend on Foundational phase completion
  - US1 (P1): No dependencies on other stories
  - US2 (P2): Can proceed in parallel with US1 (shares executor)
  - US3 (P3): Can proceed in parallel with US1/US2
  - US4 (P4): Can proceed in parallel with US1/US2/US3
  - US5 (P5): Can proceed in parallel with US1/US2/US3/US4
- **Polish (Phase 8)**: Depends on all desired user stories being complete

### Within Each User Story (TDD Flow)

For every feature:
1. Write test → verify it FAILS (Red)
2. Implement minimal code → verify test PASSES (Green)
3. Refactor if needed → verify tests still pass (Refactor)
4. Move to next test

### Parallel Opportunities

- All Setup tasks marked [P] can run in parallel
- All Foundational tests marked [P] can run in parallel
- Once Foundational phase completes, user stories can start in parallel
- Within each story, tests marked [P] can run in parallel

---

## Parallel Example: User Story 1 Tests

```bash
# Launch all version bump tests together:
Task: "Write test for BumpPatch in internal/version/bump_test.go"
Task: "Write test for BumpMinor in internal/version/bump_test.go"
Task: "Write test for BumpMajor in internal/version/bump_test.go"

# Then implement to make all pass:
Task: "Implement BumpPatch, BumpMinor, BumpMajor in internal/version/bump.go"
```

---

## Parallel Example: Multiple User Stories

```bash
# After Foundational complete, launch US2 tests while implementing US1:
Developer A: Working on US1 TUI implementation
Developer B: Writing US2 CLI tests in parallel
Developer C: Writing US3 conventional commit tests in parallel
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational (CRITICAL - blocks all stories)
3. Complete Phase 3: User Story 1 (TDD: test → implement → refactor)
4. **STOP and VALIDATE**: Test interactive bump end-to-end
5. Deploy/demo if ready - this is your MVP!

### Incremental Delivery (TDD)

For each phase:
1. Write ALL tests first (they will fail) - Red phase
2. Implement code to make tests pass - Green phase
3. Refactor for code quality - Refactor phase
4. Run golangci-lint before moving on
5. Commit after each logical group

### Task Breakdown Reminder

Each task in this file is intentionally small to support TDD:
- Test tasks are separate from implementation tasks
- One test file per feature area
- Implementation follows immediately after tests
- Refactoring is implicit after each green phase

---

## Notes

- [P] tasks = different files, no dependencies
- [Story] label maps task to specific user story for traceability
- Each user story should be independently completable and testable
- TDD cycle: Write test → Fail → Implement → Pass → Refactor
- Commit after each task or logical TDD cycle
- Stop at any checkpoint to validate story independently
- Run `golangci-lint run` frequently to catch issues early
