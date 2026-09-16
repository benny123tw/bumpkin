## Why

Bumpkin is distributed as a standalone binary via GitHub Releases (and `go install`), so users who installed a binary directly have no built-in way to upgrade — they must manually find the latest release, download the correct asset for their OS/arch, verify it, and replace the binary. A first-class `update` command removes that friction and keeps users on the latest version.

## What Changes

- Add a new `bumpkin update` subcommand (with `upgrade` as an alias) that:
  - Detects the latest GitHub Release for `benny123tw/bumpkin`.
  - Compares it against the running build's version (from the ldflags-injected `BuildInfo.Version`).
  - Downloads the correct asset for the host OS/arch, verifies it against the release `checksums.txt`, and atomically replaces the running executable in place.
- Add a `--check` (dry-run) flag that reports whether a newer release is available without modifying the binary.
- Add a `--yes`/`-y` flag to skip the confirmation prompt; default behavior prompts before replacing the binary.
- Detect non-upgradable installs and refuse gracefully with guidance:
  - Development builds with missing or unparseable version metadata (for example a local `go build`) — no release to compare against.
  - Print a clear message rather than attempting a replacement that would fail.
- Query the GitHub Releases API directly, validate the selected asset against `checksums.txt`, and use `github.com/minio/selfupdate` only for safe cross-platform binary replacement. This avoids importing the unmaintained `golang.org/x/crypto/openpgp` path pulled in by `github.com/creativeprojects/go-selfupdate`.

## Non-Goals

- Auto-update / background update checks on normal command runs — `update` is explicitly user-invoked only.
- Package-manager-managed installs (Homebrew, distro packages): the command detects when it cannot safely replace the binary and advises the appropriate package-manager command instead of self-replacing.
- Code signing / notarization of the downloaded binary beyond the existing `checksums.txt` verification.
- Downgrading to or pinning a specific arbitrary version (only "update to latest" is in scope).

## Capabilities

### New Capabilities

- `self-update`: CLI command that checks GitHub Releases for a newer version and replaces the running bumpkin binary in place, with dry-run, confirmation, and graceful handling of non-upgradable installs.

### Modified Capabilities

(none)

## Impact

- Affected specs: `self-update` (new)
- Affected code:
  - New: `internal/cli/update.go`, `internal/cli/update_test.go`
  - Modified: `internal/cli/root.go` (register the new subcommand), `go.mod`, `go.sum` (add `github.com/minio/selfupdate` and security updates)
  - Removed: (none)
- Dependencies: adds `github.com/minio/selfupdate`, updates vulnerable `go-git`, `go-billy`, and `x/text` versions, and removes `github.com/creativeprojects/go-selfupdate` plus its provider-client dependency tree.
