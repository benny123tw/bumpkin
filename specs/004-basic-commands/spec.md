# Feature Specification: Basic Subcommands

**Feature Branch**: `004-basic-commands`  
**Created**: 2024-12-14  
**Status**: Ready for Planning  
**Input**: User description: "Add basic subcommands like help and version for Go CLI conventions"

## Overview

Add subcommand-style access to common operations (`help`, `version`, `init`, `current`, `completion`) following Go CLI conventions. Users expect both `bumpkin version` and `bumpkin --version` to work, plus utility commands for initialization and shell completion.

## Clarifications

### Session 2024-12-14

- Q: Command name for showing current version? → A: `current`
- Q: Init behavior when config exists? → A: Fail with error (like goreleaser)
- Q: Which shells for completion? → A: Bash, Zsh, Fish, and PowerShell

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Version Subcommand (Priority: P1)

As a user, I want to run `bumpkin version` to see the version information so that I can quickly check what version is installed using familiar CLI patterns.

**Why this priority**: Version checking is the most common utility command. Users frequently run `<tool> version` before using a CLI tool.

**Independent Test**: Run `bumpkin version` and verify it displays version, commit, and build date.

**Acceptance Scenarios**:

1. **Given** bumpkin is installed, **When** user runs `bumpkin version`, **Then** version information is displayed (version, commit hash, build date).

2. **Given** bumpkin is installed, **When** user runs `bumpkin version`, **Then** output matches the existing `bumpkin --version` output.

3. **Given** user runs `bumpkin version --help`, **When** help is displayed, **Then** it shows a brief description of the version command.

---

### User Story 2 - Help Subcommand (Priority: P1)

As a user, I want to run `bumpkin help` to see usage information so that I can learn how to use the tool using familiar CLI patterns.

**Why this priority**: Help is essential for discoverability. New users often try `<tool> help` first.

**Independent Test**: Run `bumpkin help` and verify it displays the same help as `bumpkin --help`.

**Acceptance Scenarios**:

1. **Given** bumpkin is installed, **When** user runs `bumpkin help`, **Then** full help text is displayed (same as `--help`).

2. **Given** bumpkin is installed, **When** user runs `bumpkin help version`, **Then** help for the version subcommand is displayed.

3. **Given** bumpkin is installed, **When** user runs `bumpkin help <subcommand>`, **Then** help for that specific subcommand is displayed.

---

### User Story 3 - Init Subcommand (Priority: P2)

As a user, I want to run `bumpkin init` to create a starter configuration file so that I can quickly set up bumpkin in a new project.

**Why this priority**: Useful for new users but not essential - bumpkin works without config.

**Independent Test**: Run `bumpkin init` in a directory without `.bumpkin.yaml` and verify it creates the file.

**Acceptance Scenarios**:

1. **Given** no `.bumpkin.yaml` exists, **When** user runs `bumpkin init`, **Then** a `.bumpkin.yaml` file is created with default config and helpful comments.

2. **Given** `.bumpkin.yaml` already exists, **When** user runs `bumpkin init`, **Then** command fails with error message explaining the file exists.

3. **Given** user runs `bumpkin init --help`, **When** help is displayed, **Then** it shows usage and explains the command.

---

### User Story 4 - Current Subcommand (Priority: P2)

As a user, I want to run `bumpkin current` to quickly see the current version (latest tag) without launching the TUI.

**Why this priority**: Convenient for scripting and quick checks, but not core functionality.

**Independent Test**: Run `bumpkin current` in a repo with tags and verify it shows the latest version.

**Acceptance Scenarios**:

1. **Given** repository has tags, **When** user runs `bumpkin current`, **Then** the latest version is displayed (e.g., `v0.1.0`).

2. **Given** repository has no tags, **When** user runs `bumpkin current`, **Then** display message indicating no version found (e.g., `No version tags found`).

3. **Given** user runs `bumpkin current --prefix "ver"`, **When** command executes, **Then** it uses the specified prefix to find tags.

---

### User Story 5 - Completion Subcommand (Priority: P3)

As a user, I want to run `bumpkin completion <shell>` to generate shell completion scripts so that I can enable autocompletion in my terminal.

**Why this priority**: Nice-to-have for power users, not essential for basic usage.

**Independent Test**: Run `bumpkin completion bash` and verify it outputs a valid bash completion script.

**Acceptance Scenarios**:

1. **Given** user runs `bumpkin completion bash`, **When** command executes, **Then** valid bash completion script is output to stdout.

2. **Given** user runs `bumpkin completion zsh`, **When** command executes, **Then** valid zsh completion script is output to stdout.

3. **Given** user runs `bumpkin completion fish`, **When** command executes, **Then** valid fish completion script is output to stdout.

4. **Given** user runs `bumpkin completion powershell`, **When** command executes, **Then** valid PowerShell completion script is output to stdout.

5. **Given** user runs `bumpkin completion`, **When** no shell specified, **Then** help is displayed listing supported shells.

---

### Edge Cases

- What happens when user runs `bumpkin help nonexistent`?
  - Display error message: "Unknown command: nonexistent" and show available commands.

- What happens when user runs `bumpkin version` with other flags like `--json`?
  - The `--json` flag should not affect version output (version is always plain text).

- What happens when user runs both `bumpkin version --help` and `bumpkin help version`?
  - Both should show the same help text for the version command.

- What happens when user runs `bumpkin init` in a directory without write permissions?
  - Display error message explaining the permission issue.

- What happens when user runs `bumpkin init` and `.bumpkin.yaml` already exists?
  - Command fails with error: "Error: .bumpkin.yaml already exists" (no overwrite, like goreleaser).

- What happens when user runs `bumpkin current` outside a git repository?
  - Display error message: "Error: not a git repository".

- What happens when user runs `bumpkin current` with no tags in the repository?
  - Display message: "No version tags found".

- What happens when user runs `bumpkin completion` without specifying a shell?
  - Display help listing supported shells (bash, zsh, fish, powershell).

- What happens when user runs `bumpkin completion unsupported`?
  - Display error message: "Unsupported shell: unsupported" and list valid options.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST provide a `version` subcommand that displays version information.
- **FR-002**: System MUST provide a `help` subcommand that displays usage information.
- **FR-003**: The `version` subcommand MUST display the same information as `--version` flag.
- **FR-004**: The `help` subcommand MUST display the same information as `--help` flag.
- **FR-005**: Running `bumpkin help <command>` MUST display help for that specific command.
- **FR-006**: Both flags (`--version`, `--help`) and subcommands (`version`, `help`) MUST continue to work.
- **FR-007**: Unknown subcommands MUST display a helpful error with available commands listed.
- **FR-008**: System MUST provide an `init` subcommand that creates a `.bumpkin.yaml` configuration file.
- **FR-009**: The `init` subcommand MUST fail with an error if `.bumpkin.yaml` already exists.
- **FR-010**: System MUST provide a `current` subcommand that displays the latest version tag.
- **FR-011**: The `current` subcommand MUST support `--prefix` flag to filter tags.
- **FR-012**: System MUST provide a `completion` subcommand that generates shell completion scripts.
- **FR-013**: The `completion` subcommand MUST support bash, zsh, fish, and powershell.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: `bumpkin version` produces identical output to `bumpkin --version`.
- **SC-002**: `bumpkin help` produces identical output to `bumpkin --help`.
- **SC-003**: Users familiar with Go CLI tools can use bumpkin without reading documentation.
- **SC-004**: All existing flag-based functionality continues to work unchanged.
- **SC-005**: `bumpkin init` creates a valid, commented `.bumpkin.yaml` configuration file.
- **SC-006**: `bumpkin current` correctly identifies the latest semver tag.
- **SC-007**: Generated shell completion scripts work correctly for all supported shells.

## Assumptions

- Cobra (the CLI framework already in use) natively supports subcommands, making this a straightforward addition.
- The root command (no subcommand) will continue to launch the interactive TUI or execute based on flags.
- Subcommands take precedence over positional arguments if there's ambiguity.
