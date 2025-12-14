# Feature Specification: Version Tagger CLI

**Feature Branch**: `001-version-tagger`  
**Created**: 2025-12-14  
**Status**: Draft  
**Input**: User description: "Build a CLI app that helps tag commits by analyzing conventional commit history, providing version options for users to select or customize, with TUI and non-interactive modes. Inspired by antfu's bumpp but language-agnostic and built in Go."

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Interactive Version Bump (Priority: P1)

A developer has completed work on their project and wants to create a new version tag. They run the CLI tool which displays the commit history since the last tag, shows the current version, and presents version bump options. The developer selects the appropriate version type and the tool creates and pushes the tag.

**Why this priority**: This is the core functionality - the primary reason for the tool's existence. Without interactive version bumping, the tool provides no value.

**Independent Test**: Can be fully tested by running the CLI in a git repository with commits since the last tag, selecting a version option, and verifying the tag is created correctly.

**Acceptance Scenarios**:

1. **Given** a git repository with tag v1.2.3 and 5 commits since that tag, **When** user runs the CLI, **Then** they see the 5 commits listed, current version "v1.2.3" displayed, and version options presented
2. **Given** the CLI is showing version options, **When** user selects "minor", **Then** the system creates tag v1.3.0 and pushes it to the remote
3. **Given** the CLI is showing version options, **When** user selects "custom" and enters "v2.0.0-beta.1", **Then** the system creates that exact tag and pushes it

---

### User Story 2 - Non-Interactive Mode (Priority: P2)

A developer wants to automate version bumping in their CI/CD pipeline or quickly bump a version without going through the interactive interface. They run a single command specifying the version bump type and the tool executes without prompts.

**Why this priority**: Enables automation and faster workflows for experienced users. Essential for CI/CD integration but requires the core tagging logic from P1.

**Independent Test**: Can be tested by running the CLI with command-line flags specifying version type and verifying the tag is created without any user prompts.

**Acceptance Scenarios**:

1. **Given** a repository at v1.0.0, **When** user runs `bumpkin --patch`, **Then** tag v1.0.1 is created and pushed without any prompts
2. **Given** a repository at v1.0.0, **When** user runs `bumpkin --version v2.0.0`, **Then** tag v2.0.0 is created and pushed
3. **Given** a repository at v1.0.0, **When** user runs `bumpkin --minor --dry-run`, **Then** the system shows what would happen but creates no tag

---

### User Story 3 - Conventional Commit Analysis (Priority: P3)

A developer wants the tool to automatically suggest the appropriate version bump based on their conventional commit messages. The tool analyzes commits for breaking changes, features, and fixes to recommend major, minor, or patch bumps.

**Why this priority**: Adds intelligence to the tool but is not required for basic functionality. Users can always manually select version type.

**Independent Test**: Can be tested by creating a repository with specific conventional commits and verifying the tool suggests the correct version bump type.

**Acceptance Scenarios**:

1. **Given** commits containing "feat:" since last tag, **When** user runs CLI in conventional mode, **Then** "minor" bump is recommended
2. **Given** commits containing "fix:" only since last tag, **When** user runs CLI in conventional mode, **Then** "patch" bump is recommended
3. **Given** commits containing "BREAKING CHANGE:" or "feat!:", **When** user runs CLI in conventional mode, **Then** "major" bump is recommended

---

### User Story 4 - Hook System for Custom Actions (Priority: P4)

A developer using the tool for a Node.js project wants to update their package.json version before tagging. They configure a hook that runs a custom script to update version files, then the tool creates the tag.

**Why this priority**: Extends the tool's usefulness to projects that need version file updates, but the core tagging works without it.

**Independent Test**: Can be tested by configuring a pre-tag hook script that creates a file, running a version bump, and verifying the script executed before the tag was created.

**Acceptance Scenarios**:

1. **Given** a configured pre-tag hook script, **When** user completes version selection, **Then** the hook script runs before the tag is created
2. **Given** a pre-tag hook that fails (non-zero exit), **When** user completes version selection, **Then** the tag is NOT created and error is displayed
3. **Given** a configured post-tag hook script, **When** tag is successfully created, **Then** the post-tag hook runs after tagging

---

### User Story 5 - Prerelease Versions (Priority: P5)

A developer wants to create prerelease versions like alpha, beta, or release candidates. They select prerelease options and the tool creates appropriate prerelease tags following semver conventions.

**Why this priority**: Important for projects with formal release processes but not required for basic version bumping.

**Independent Test**: Can be tested by selecting prerelease options and verifying tags follow correct semver prerelease format.

**Acceptance Scenarios**:

1. **Given** current version v1.0.0, **When** user selects "prerelease alpha", **Then** tag v1.0.1-alpha.0 is created
2. **Given** current version v1.0.1-alpha.0, **When** user selects "prerelease alpha", **Then** tag v1.0.1-alpha.1 is created
3. **Given** current version v1.0.1-beta.2, **When** user selects "release", **Then** tag v1.0.1 is created (removing prerelease suffix)

---

### Edge Cases

- What happens when the repository has no existing tags? System assumes v0.0.0 as starting point.
- What happens when there are no commits since last tag? System warns user but allows manual version bump.
- What happens when user is not on a git repository? System displays clear error and exits.
- What happens when git push fails (no remote, auth failure)? System displays error but the local tag remains; user can retry push.
- What happens with merge commits in the history? System includes them in commit list; conventional commit parsing only uses first line.
- What happens with lightweight vs annotated tags? System creates annotated tags by default.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST detect the current git repository and its remote configuration
- **FR-002**: System MUST find the most recent semver tag in the repository
- **FR-003**: System MUST list all commits between the last tag and HEAD
- **FR-004**: System MUST display commits in reverse chronological order with hash, message, and author
- **FR-005**: System MUST support version bump types: major, minor, patch, and custom
- **FR-006**: System MUST support prerelease versions: alpha, beta, rc (release candidate)
- **FR-007**: System MUST create annotated git tags with the selected version
- **FR-008**: System MUST push tags to the configured remote (default: origin)
- **FR-009**: System MUST provide a TUI interface for interactive mode
- **FR-010**: System MUST provide a non-interactive CLI mode with flags for automation
- **FR-011**: System MUST parse conventional commits to recommend version bump type
- **FR-012**: System MUST support a "dry-run" mode that shows actions without executing them
- **FR-013**: System MUST support configurable hooks that run before and after tagging
- **FR-014**: System MUST abort tagging if pre-tag hooks fail (non-zero exit code)
- **FR-015**: System MUST support custom tag prefixes (default: "v", e.g., v1.0.0)
- **FR-016**: System MUST validate that custom versions follow semver format
- **FR-017**: System MUST handle repositories with no existing tags (starting from v0.0.0)
- **FR-018**: System MUST provide clear error messages for git operation failures

### Key Entities

- **Version**: Represents a semantic version with major, minor, patch, and optional prerelease/build metadata
- **Commit**: Represents a git commit with hash, message, author, and timestamp
- **Tag**: Represents a git tag with name, associated commit, and optional annotation
- **Hook**: Represents a user-defined script that runs at specific lifecycle points (pre-tag, post-tag)
- **Configuration**: User settings including tag prefix, remote name, hook definitions

## Assumptions

- Users have git installed and configured on their system
- Repositories use semantic versioning (semver) for tags
- The default remote is named "origin" unless configured otherwise
- Conventional commits follow the standard format: `type(scope): description`
- Hook scripts are executable files on the user's system
- Users have push access to the remote repository

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Users can complete a version bump and push in under 30 seconds using interactive mode
- **SC-002**: Non-interactive mode completes version bump in under 5 seconds (excluding network latency)
- **SC-003**: Conventional commit analysis correctly identifies bump type for 95% of standard commit messages
- **SC-004**: Tool works on any git repository regardless of programming language or project type
- **SC-005**: Users can configure and run custom hooks without modifying tool source code
- **SC-006**: Error messages enable users to resolve issues without external documentation in 90% of cases
- **SC-007**: Tool correctly handles prerelease version sequences (incrementing prerelease numbers)
