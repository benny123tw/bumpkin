# Feature Specification: TUI Commit Display Enhancement

**Feature Branch**: `003-tui-commit-display`  
**Created**: 2024-12-14  
**Status**: Draft  
**Input**: User description: "Improve TUI commit history display: show commits during version selection, highlight breaking changes with colored type badges like bumpp"

## Overview

Enhance the TUI to persistently display commit history during the version selection process, and add visual highlighting for commit types (feat, fix, docs, etc.) with special emphasis on breaking changes (commits with `!`).

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Persistent Commit History Display (Priority: P1)

As a developer, I want to see the commit history while selecting a version so that I can make an informed decision about which version bump to choose without having to remember what commits were made.

**Why this priority**: This is the core value - users need context about their commits when deciding on a version bump. Currently, commits disappear when entering version selection.

**Independent Test**: Run bumpkin, navigate to version selection, verify commits remain visible above the version options.

**Acceptance Scenarios**:

1. **Given** user runs bumpkin in interactive mode, **When** they navigate from commit list to version selection, **Then** the commit history remains visible above the version options.

2. **Given** user is on version selection screen, **When** they view the display, **Then** they see both commit history and version options in a single view.

3. **Given** there are more commits than can fit on screen, **When** user views commit history, **Then** commits are truncated with a count indicator (e.g., "and 5 more commits...").

---

### User Story 2 - Colored Commit Type Badges (Priority: P1)

As a developer, I want commit types (feat, fix, docs, chore, etc.) to be visually highlighted with colors so that I can quickly scan and understand the nature of changes.

**Why this priority**: Visual differentiation makes it much faster to understand the commit history at a glance, similar to how bumpp displays commits.

**Independent Test**: Run bumpkin with conventional commits, verify each commit type has a distinct colored badge.

**Acceptance Scenarios**:

1. **Given** a commit with type `feat:`, **When** displayed in commit history, **Then** the type is shown with a green/lime badge.

2. **Given** a commit with type `fix:`, **When** displayed in commit history, **Then** the type is shown with a yellow badge.

3. **Given** a commit with type `docs:`, **When** displayed in commit history, **Then** the type is shown with a blue badge.

4. **Given** a commit with type `chore:`, **When** displayed in commit history, **Then** the type is shown with a gray badge.

5. **Given** a commit without conventional format, **When** displayed in commit history, **Then** no type badge is shown, just the commit message.

---

### User Story 3 - Breaking Change Highlighting (Priority: P1)

As a developer, I want breaking change commits (those with `!`) to be prominently highlighted so that I know a major version bump may be required.

**Why this priority**: Breaking changes are critical to identify as they require a major version bump. Missing them could lead to incorrect versioning.

**Independent Test**: Create commits with `feat!:` or `fix!:`, run bumpkin, verify these are highlighted distinctly (red background on type badge).

**Acceptance Scenarios**:

1. **Given** a commit with type `feat!:`, **When** displayed in commit history, **Then** the type badge has a red/warning background color indicating breaking change.

2. **Given** a commit with type `fix!:`, **When** displayed in commit history, **Then** the type badge has a red/warning background color.

3. **Given** a commit with `BREAKING CHANGE:` in the body, **When** displayed in commit history, **Then** it is treated as a breaking change with red highlight.

4. **Given** multiple commits with one breaking change, **When** displayed, **Then** the breaking change commit stands out visually from other commits.

---

### User Story 4 - Commit Display Format (Priority: P2)

As a developer, I want commits displayed in a clean, consistent format showing hash, type badge, and description so that I can quickly read and understand each commit.

**Why this priority**: A clean format improves readability and user experience, but the core functionality works without it.

**Independent Test**: Run bumpkin with various commits, verify consistent formatting across all commit types.

**Acceptance Scenarios**:

1. **Given** a conventional commit, **When** displayed, **Then** format is: `<short-hash>  <type-badge> : <description>` (like bumpp).

2. **Given** commit hash is 7+ characters, **When** displayed, **Then** only the first 7 characters are shown.

3. **Given** commit description is very long, **When** displayed, **Then** it is truncated to fit the terminal width.

---

### Edge Cases

- What happens when there are no commits since last tag?
  - Show message "No commits since last tag" instead of empty commit list.

- What happens when commit message doesn't follow conventional format?
  - Display without type badge, just: `<hash>  <full-message>`.

- What happens when terminal is very narrow?
  - Truncate commit messages appropriately to prevent wrapping.

- What happens when there are 50+ commits?
  - Show first N commits (configurable, default 10) with "and X more commits..." indicator.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST display commit history on the version selection screen.
- **FR-002**: System MUST parse conventional commit types from commit messages.
- **FR-003**: System MUST display commit type as a colored badge before the commit description.
- **FR-004**: System MUST highlight breaking changes (commits with `!` after type) with a distinct red/warning color.
- **FR-005**: System MUST display commit hash (first 7 characters) for each commit.
- **FR-006**: System MUST limit displayed commits to prevent screen overflow (default: 10).
- **FR-007**: System MUST show count of remaining commits if truncated.
- **FR-008**: System MUST handle non-conventional commits gracefully (display without badge).

### Color Scheme

| Commit Type | Badge Color |
|-------------|-------------|
| `feat` | Green/Lime |
| `fix` | Yellow |
| `docs` | Blue |
| `chore` | Gray |
| `refactor` | Cyan |
| `test` | Magenta |
| `style` | Gray |
| `perf` | Orange |
| `ci` | Gray |
| `build` | Gray |
| Breaking (`!`) | Red background |

### Key Entities

- **CommitType**: Parsed conventional commit type (feat, fix, docs, etc.)
- **CommitDisplay**: Formatted commit for TUI display (hash, type, description, isBreaking)

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Commit history is visible on version selection screen without requiring user to scroll back.
- **SC-002**: Breaking change commits are immediately identifiable by red highlighting.
- **SC-003**: Users can identify commit types in under 1 second per commit via color coding.
- **SC-004**: TUI remains responsive with up to 100 commits in history.
- **SC-005**: Display format matches bumpp style: `<hash>  <type> : <description>`.
