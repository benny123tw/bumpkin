# Feature Specification: Release Tool Integration

**Feature Branch**: `002-release-integration`  
**Created**: 2024-12-14  
**Status**: Draft  
**Input**: User description: "I want to integrate with release tool particularly goreleaser. We have discussed to enhance the new hook and add a release config. We hope this feature can be more flexible, customizable but also can provide an out-of-box experience."

## Overview

This feature enhances bumpkin with a new `post-push` hook phase to support CI/CD-triggered release workflows. Rather than executing release tools locally (which is atypical in the Go community where GitHub Actions handles releases via goreleaser-action), bumpkin focuses on providing hooks that integrate seamlessly with CI/CD pipelines triggered by tag pushes.

## Clarifications

### Session 2024-12-14

- Q: Should bumpkin focus on local development releases, or should it primarily support CI/CD-triggered workflows? → A: CI/CD-first - Remove local release execution; focus only on `post-push` hook for notifications/triggers; let CI/CD handle actual releases.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Post-Push Hook Phase (Priority: P1)

As a developer, I want to run commands after a successful push so that I can trigger notifications, update external systems, or perform cleanup tasks after my tag is pushed and CI/CD takes over.

**Why this priority**: This is the core value - enabling workflows that depend on the tag being pushed, complementing CI/CD-based release pipelines.

**Independent Test**: Can be fully tested by configuring post-push hooks and running a bump, then verifying the hooks executed after the push completed.

**Acceptance Scenarios**:

1. **Given** `.bumpkin.yml` contains `hooks.post-push` commands, **When** user runs a bump that pushes a tag, **Then** the post-push hooks execute after successful push.

2. **Given** the push fails (no remote, network error), **When** bumpkin cannot push the tag, **Then** post-push hooks do not execute.

3. **Given** a post-push hook fails, **When** the hook returns non-zero exit code, **Then** bumpkin reports the failure but the tag remains pushed (warning, not fatal).

4. **Given** `--no-push` flag is used, **When** user runs a bump, **Then** post-push hooks do not execute.

5. **Given** `--dry-run` flag is used, **When** user runs a bump, **Then** post-push hooks are shown but not executed.

---

### User Story 2 - Post-Push Notifications (Priority: P1)

As a developer, I want to send notifications after a tag is pushed so that my team knows a new release is being built by CI/CD.

**Why this priority**: Common use case that demonstrates the value of post-push hooks in CI/CD workflows.

**Independent Test**: Can be tested by configuring a notification command (e.g., curl to Slack webhook) and verifying it executes after push.

**Acceptance Scenarios**:

1. **Given** `.bumpkin.yml` contains a post-push hook like `curl -X POST $SLACK_WEBHOOK -d "New release $BUMPKIN_TAG"`, **When** tag is pushed, **Then** the notification is sent with correct version info.

2. **Given** notification hook uses BUMPKIN_* environment variables, **When** hook executes, **Then** all standard variables (BUMPKIN_VERSION, BUMPKIN_TAG, BUMPKIN_COMMIT, etc.) are available.

---

### User Story 3 - Multiple Post-Push Hooks (Priority: P2)

As a developer, I want to run multiple commands after push so that I can trigger several actions (notify, update docs, trigger downstream builds).

**Why this priority**: Extends the basic hook functionality with multi-command support.

**Independent Test**: Can be tested by configuring multiple post-push hooks and verifying all execute in order.

**Acceptance Scenarios**:

1. **Given** `.bumpkin.yml` contains multiple post-push hooks, **When** tag is pushed, **Then** hooks execute in the order defined.

2. **Given** one hook fails in a sequence, **When** failure occurs, **Then** subsequent hooks still execute (fail-open for notifications).

---

### User Story 4 - TUI Post-Push Hook Display (Priority: P3)

As a developer using the TUI, I want to see post-push hook execution results so that I know what actions were triggered after my tag was pushed.

**Why this priority**: Nice-to-have enhancement for TUI users; the hook execution itself is already covered.

**Independent Test**: Can be tested by running TUI, completing a bump, and verifying post-push hook output is displayed.

**Acceptance Scenarios**:

1. **Given** post-push hooks are configured, **When** user completes a bump in TUI, **Then** TUI shows hook execution status after push completes.

2. **Given** a post-push hook fails, **When** displayed in TUI, **Then** the failure is shown as a warning (not blocking).

---

### Edge Cases

- What happens when post-push hook takes too long?
  - Hooks have a default timeout (30 seconds); timeout is treated as a warning, not fatal.

- What happens when multiple hooks are configured and one fails?
  - Remaining hooks continue to execute; all failures are reported at the end.

- What happens when hook command contains shell special characters?
  - Commands are executed via shell (sh -c on Unix, cmd /C on Windows) so standard shell features work.

- What happens when post-push hooks are configured but user uses --no-push?
  - Post-push hooks are skipped; bumpkin may optionally inform user that hooks were skipped.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST add a `post-push` hook phase that executes after successful tag push.
- **FR-002**: System MUST skip post-push hooks when `--no-push` flag is used.
- **FR-003**: System MUST skip post-push hooks when push fails.
- **FR-004**: System MUST pass all standard environment variables (BUMPKIN_VERSION, BUMPKIN_TAG, BUMPKIN_COMMIT, BUMPKIN_PREVIOUS_VERSION, BUMPKIN_PREFIX, BUMPKIN_REMOTE, BUMPKIN_DRY_RUN) to post-push hooks.
- **FR-005**: System MUST execute multiple post-push hooks in defined order.
- **FR-006**: System MUST continue executing remaining hooks if one fails (fail-open behavior).
- **FR-007**: System MUST report post-push hook failures as warnings, not fatal errors.
- **FR-008**: System MUST respect hook timeout (default 30 seconds) for post-push hooks.
- **FR-009**: System MUST show post-push hooks in dry-run output without executing them.
- **FR-010**: System MUST display post-push hook execution status in TUI mode.

### Configuration Schema

The `.bumpkin.yml` hooks section extended with post-push:

```yaml
hooks:
  pre-tag:
    - "./scripts/update-changelog.sh"
  post-tag:
    - "echo Tagged $BUMPKIN_TAG"
  post-push:
    - "curl -X POST $SLACK_WEBHOOK -d '{\"text\": \"Released $BUMPKIN_TAG\"}'"
    - "./scripts/notify-team.sh"
```

### Hook Execution Order

```
pre-tag → create tag → post-tag → push → post-push
```

### Key Entities

- **HookPhase**: Extended to include `post-push` in addition to existing `pre-tag` and `post-tag`.
- **HookResult**: Existing entity; post-push hooks use the same result structure.

## Out of Scope

The following features were considered but explicitly excluded to align with CI/CD-first workflows:

- **Local release execution** (`--release` flag) - CI/CD handles releases via goreleaser-action
- **GoReleaser auto-detection** - Not needed since CI/CD runs goreleaser
- **Release configuration section** - Hooks provide sufficient flexibility
- **Release provider abstraction** - Unnecessary complexity for hook-based approach

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Post-push hooks execute within 1 second of successful push completion.
- **SC-002**: Users can configure post-push hooks in under 1 minute by adding lines to `.bumpkin.yml`.
- **SC-003**: Hook failures do not prevent successful tag creation and push.
- **SC-004**: All BUMPKIN_* environment variables are available to post-push hooks.
- **SC-005**: Dry-run accurately shows which post-push hooks would execute.
- **SC-006**: TUI displays post-push hook results clearly.

## Assumptions

- CI/CD pipelines (e.g., GitHub Actions with goreleaser-action) handle the actual release build and artifact publishing.
- Post-push hooks are primarily used for notifications, triggering external systems, or lightweight cleanup tasks.
- Hook commands are trusted (user-defined in config file).
- Network-dependent hooks (e.g., Slack notifications) may fail due to connectivity issues; this is acceptable as a warning.
- Post-push hooks share the same environment variable set as existing pre-tag and post-tag hooks.
