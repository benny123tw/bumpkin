## Context

Bumpkin ships as a standalone binary. Releases are built by GoReleaser (`.goreleaser.yml`) and published to GitHub Releases under `benny123tw/bumpkin`. Key facts that constrain the design:

- Release assets are **raw binaries** (`archives.formats: ["binary"]`), named `bumpkin-{version}-{os}-{arch}` on Unix-like systems and `bumpkin-{version}-windows-{arch}.exe` on Windows (e.g. `bumpkin-1.2.3-darwin-arm64`, `bumpkin-1.2.3-windows-amd64.exe`). They are NOT tar.gz/zip archives.
- A `checksums.txt` (SHA-256) asset is published alongside the binaries.
- The running build's version metadata is injected via ldflags into `main.version` / `main.commit` / `main.date` and surfaced through `cli.BuildInfo`. Local development builds generally expose `unknown` or another non-semver value; versioned `go install module@version` builds can derive their version from Go build info.
- Commands follow a struct + factory + `execute(cmd, args)` pattern (see the `current` and `version` commands) and are registered in `internal/cli/root.go` via `rootCmd.AddCommand(...)`.
- No HTTP client or release-fetching code exists in the project today.

## Goals / Non-Goals

**Goals:**

- Provide a user-invoked `bumpkin update` command (alias `upgrade`) that upgrades the running binary to the latest GitHub Release.
- Verify integrity via the published `checksums.txt` before replacing the binary.
- Replace the executable atomically and cross-platform (including Windows' running-exe constraint).
- Fail safe and informatively for non-upgradable installs (dev builds, package-manager installs, unwritable binary path).
- Offer a `--check` dry-run and a `--yes` non-interactive confirmation bypass.

**Non-Goals:**

- Background/automatic update checks during normal command runs.
- Installing to a managed location on the user's behalf (Homebrew/distro packages own their binaries).
- Arbitrary version pinning or downgrade (only "latest" is targeted).
- Additional signature/notarization beyond checksum verification.

## Decisions

### Use go-selfupdate for release discovery and binary replacement

Use `github.com/creativeprojects/go-selfupdate` (maintained successor to `rhysd/go-github-selfupdate`). It queries the GitHub Releases API, selects the asset matching the host OS/arch, validates against `checksums.txt`, and uses its built-in cross-platform replacement logic.

Alternatives considered:
- `minio/selfupdate` directly — only does the byte-level replacement; we would have to hand-write GitHub API calls, OS/arch asset matching, version compare, and checksum parsing. Rejected as re-implementing go-selfupdate.
- Hand-rolled `net/http` against the GitHub API — same drawback, plus more surface area to maintain and test. Rejected.

### Configure the asset matcher for GoReleaser's binary naming

go-selfupdate's default asset filtering targets archive-style names. Because bumpkin publishes raw binaries, configure the updater with an explicit `Filters`/asset-matching regex anchored to the current runtime platform: `^bumpkin-.*-{runtime.GOOS}-{runtime.GOARCH}$` on Unix-like systems and the equivalent pattern ending in `\.exe` on Windows. In go-selfupdate, a non-empty `Filters` list replaces the library's built-in OS/architecture suffix matching, so a broad product-only filter such as `^bumpkin-.*$` is unsafe: it can select the first asset for another platform. The platform-specific filter must uniquely match the host binary and exclude `checksums.txt`. `amd64` maps directly; verify `arm64` (Apple Silicon, linux/arm64) and the Windows executable suffix resolve. The updater is pointed at the `benny123tw/bumpkin` repository slug.

### Version source and comparison

Read the current version from the `BuildInfo` passed into the command (same value `bumpkin version` prints). Strip a leading `v` as needed and let go-selfupdate's semver comparison decide whether the latest release is newer. If the latest release is not newer than the current version, report "already up to date" and exit 0 without downloading.

### Guard non-upgradable installs before attempting replacement

Before contacting GitHub, detect and refuse with guidance when:
- The version is missing, `unknown`, `dev`, or otherwise not valid semver — instruct the user to install a released binary or use `go install ...@latest`.
- The resolved executable path is not writable, or sits under a known package-manager prefix (e.g. a Homebrew Cellar/`bin` path) — instruct the user to use their package manager (`brew upgrade bumpkin`) instead of self-replacing.

These checks produce a clear message and a non-zero exit only on genuine error; an intentional "use your package manager" refusal is communicated explicitly.

### Flags and confirmation UX

- Default (interactive): print current → latest version, prompt for confirmation, then replace.
- `--yes`/`-y`: skip the prompt (for scripts/CI).
- `--check`: report availability (current, latest, whether newer) and exit WITHOUT downloading or replacing; exit 0.
- `--prefix` is NOT needed (this command is about the tool's own release tags, not repo tags).

## Implementation Contract

**Behavior:**

- `bumpkin update` (and `bumpkin upgrade`) checks GitHub Releases for `benny123tw/bumpkin`, and when a newer release exists, downloads the host-matching asset, verifies it against `checksums.txt`, and replaces the currently running executable in place. On success it prints the old and new version (e.g. `Updated bumpkin v1.2.0 -> v1.3.0`).
- When already current, it prints a message such as `bumpkin v1.3.0 is already the latest version` and exits 0.
- `bumpkin update --check` performs only the version check: prints current version, latest version, and whether an update is available; never writes to disk; exits 0.
- Confirmation: without `--yes`, the command prompts before replacing and aborts (exit 0) if the user declines. With `--yes`, it proceeds without prompting.

**Interface / data shape:**

- New file `internal/cli/update.go` defines `updateCommand` (struct), `newUpdateCommand(info BuildInfo) *updateCommand`, and `(*updateCommand).execute(cmd *cobra.Command, args []string) error`, mirroring `newVersionCommand` / `newCurrentCommand`.
- Cobra command: `Use: "update"`, `Aliases: []string{"upgrade"}`, `Args: cobra.NoArgs`, flags `--check` (bool), `--yes`/`-y` (bool).
- Registered in `internal/cli/root.go` alongside the existing `rootCmd.AddCommand(...)` calls, passing the root command's `BuildInfo`.
- Output is written via `cmd.OutOrStdout()`; prompts read from `cmd.InOrStdin()` so the command is testable.

**Failure modes:**

- Development build (`version == "unknown"`): print guidance, return a non-nil error (or exit non-zero) indicating self-update is unavailable for dev builds.
- Package-manager / non-writable binary path: print the package-manager guidance and exit without attempting replacement.
- Network/API failure, no matching asset, or checksum mismatch: return a wrapped error (`fmt.Errorf("...: %w", err)`) surfaced through the existing root error handling; the binary is left untouched.

**Acceptance criteria:**

- `internal/cli/update_test.go` covers: command/alias registration and flag wiring; the dev-build guard short-circuits before any network call; `--check` reports status without replacing; version-comparison decision (newer vs. up-to-date). GitHub interaction and the actual binary swap are abstracted behind a small seam so tests do not hit the network or modify the test binary.
- `just build` succeeds with the new dependency; `bumpkin update --help` shows the command, alias, and flags.

**Scope boundaries:**

- In scope: the `update`/`upgrade` command, its flags, guards, dependency addition, registration, and tests.
- Out of scope: changes to release/build configuration, auto-update behavior, Homebrew tap automation, and any modification of existing commands beyond registration in `root.go`.

## Risks / Trade-offs

- [GoReleaser raw-binary naming not matched by go-selfupdate defaults] → Configure an explicit asset-matching filter and add a test/asserted regex; verify against actual release asset names for each supported OS/arch.
- [A broad custom filter bypasses go-selfupdate's platform matching] → Construct the production filter from the target OS and architecture and test it against a mixed-platform release asset list.
- [arm64 / Apple Silicon asset not selected] → Confirm the OS/arch token mapping the library uses matches the `{os}-{arch}` tokens GoReleaser emits; cover in the asset-matcher decision.
- [Self-replacement fails on locked/managed binary] → Detect non-writable / package-manager paths up front and refuse with guidance rather than producing a partial write.
- [New transitive dependency surface] → go-selfupdate pulls in provider API clients and related dependencies; run `go mod tidy` and `just check` to confirm the tree builds and lints cleanly.
- [Network/API rate limits without auth] → Unauthenticated GitHub API is sufficient for a single release check; document that repeated checks may hit anonymous rate limits.
