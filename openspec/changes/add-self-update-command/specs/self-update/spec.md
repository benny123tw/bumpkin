## ADDED Requirements

### Requirement: Self-update command

The CLI SHALL provide an `update` subcommand, with `upgrade` as an alias, that updates the running bumpkin executable to the latest published GitHub Release for the `benny123tw/bumpkin` repository. The command SHALL take no positional arguments.

#### Scenario: A newer release is available

- **WHEN** the user runs `bumpkin update` from a released build and a newer release exists
- **THEN** the command downloads the asset matching the host OS and architecture, verifies it against the release `checksums.txt`, atomically replaces the running executable, prints the old and new versions, and exits 0

##### Example: Select only the host platform asset

- **GIVEN** a release contains binaries for Darwin, Linux, and Windows, plus `checksums.txt`
- **WHEN** `bumpkin update` runs on a supported host platform
- **THEN** only that host platform's binary matches the updater filter

| Host platform | Matching asset |
| ------------- | -------------- |
| `darwin/arm64` | `bumpkin-1.3.0-darwin-arm64` |
| `linux/amd64` | `bumpkin-1.3.0-linux-amd64` |
| `windows/amd64` | `bumpkin-1.3.0-windows-amd64.exe` |

#### Scenario: Already on the latest version

- **WHEN** the user runs `bumpkin update` and the running version is equal to or newer than the latest release
- **THEN** the command reports that bumpkin is already up to date, does not modify the binary, and exits 0

#### Scenario: Invoked via the upgrade alias

- **WHEN** the user runs `bumpkin upgrade`
- **THEN** the command behaves identically to `bumpkin update`

### Requirement: Version comparison

The command SHALL determine the running version from the build metadata (the same value reported by `bumpkin version`) and SHALL only replace the binary when the latest release version is strictly newer per semantic-version ordering.

#### Scenario: Newer, equal, and older releases

- **WHEN** the latest release version is compared against the running version
- **THEN** the command decides whether to update according to semantic-version ordering

##### Example: update decision by version pair

| Running version | Latest release | Action |
| --------------- | -------------- | ------ |
| v1.2.0 | v1.3.0 | update to v1.3.0 |
| v1.3.0 | v1.3.0 | already up to date, no change |
| v1.4.0 | v1.3.0 | already up to date, no change |

### Requirement: Check-only mode

The command SHALL accept a `--check` flag that reports update availability without downloading or replacing the binary.

#### Scenario: Check reports an available update

- **WHEN** the user runs `bumpkin update --check` and a newer release exists
- **THEN** the command prints the current version, the latest version, and that an update is available, makes no changes to disk, and exits 0

#### Scenario: Check reports up to date

- **WHEN** the user runs `bumpkin update --check` and no newer release exists
- **THEN** the command prints that bumpkin is up to date, makes no changes to disk, and exits 0

### Requirement: Confirmation prompt

Without the `--yes` flag, the command SHALL prompt for confirmation before replacing the binary. With the `--yes` (`-y`) flag, the command SHALL proceed without prompting.

#### Scenario: User confirms the update

- **WHEN** a newer release is available, the user runs `bumpkin update` without `--yes`, and answers the prompt affirmatively
- **THEN** the command proceeds to download and replace the binary

#### Scenario: User declines the update

- **WHEN** a newer release is available, the user runs `bumpkin update` without `--yes`, and answers the prompt negatively
- **THEN** the command aborts without modifying the binary and exits 0

#### Scenario: Non-interactive confirmation

- **WHEN** the user runs `bumpkin update --yes` and a newer release is available
- **THEN** the command proceeds without prompting

### Requirement: Refuse non-upgradable installs

The command SHALL detect installs that cannot be safely self-replaced and SHALL refuse with actionable guidance instead of attempting a replacement, before contacting the network.

#### Scenario: Development build

- **WHEN** the user runs `bumpkin update` from a build whose version metadata is missing or not valid semver (for example a local development build)
- **THEN** the command reports that self-update is unavailable for development builds, advises installing a released binary, makes no network request, and exits non-zero

#### Scenario: Package-manager or non-writable install

- **WHEN** the resolved executable path is not writable or resides under a known package-manager prefix
- **THEN** the command advises using the appropriate package manager (for example `brew upgrade`) instead of self-replacing, and does not attempt to replace the binary

### Requirement: Integrity verification and safe failure

The command SHALL verify the downloaded asset against the release `checksums.txt` before replacement, and SHALL leave the existing binary untouched when any step fails.

#### Scenario: Download or verification fails

- **WHEN** the network request fails, no asset matches the host OS and architecture, or the checksum does not match
- **THEN** the command returns a descriptive error, leaves the running binary unchanged, and exits non-zero

##### Example: Checksum mismatch preserves the installed binary

- **GIVEN** bumpkin v1.2.0 is installed and GitHub publishes v1.3.0 with a checksum that does not match the downloaded asset
- **WHEN** the user confirms `bumpkin update`
- **THEN** the command reports the checksum verification failure, exits non-zero, and the installed v1.2.0 binary remains unchanged
