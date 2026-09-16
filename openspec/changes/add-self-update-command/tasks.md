## 1. Dependency setup

- [x] 1.1 Use go-selfupdate for release discovery and binary replacement: add `github.com/creativeprojects/go-selfupdate` by running `go get github.com/creativeprojects/go-selfupdate` and `go mod tidy`. Verified by `just build` succeeding and the dependency appearing in `go.mod`/`go.sum`.

## 2. Self-update command skeleton

- [x] 2.1 Implement the self-update command in `internal/cli/update.go` following the existing struct + factory + `execute` pattern: define `updateCommand`, `newUpdateCommand(info BuildInfo) *updateCommand`, and `execute`, with cobra `Use: "update"`, `Aliases: []string{"upgrade"}`, `Args: cobra.NoArgs`. Implements the upgrade alias behavior from the flags and confirmation UX decision. Verified by a test asserting the command and `upgrade` alias resolve and by `bumpkin update --help`.
- [x] 2.2 Register the command in `internal/cli/root.go` via `rootCmd.AddCommand(newUpdateCommand(info).cmd)`, passing the root `BuildInfo`. Verified by `bumpkin update` being reachable from the root command in a test and in `bumpkin --help` output.

## 3. Guard non-upgradable installs

- [x] 3.1 Implement the dev-build guard satisfying "Refuse non-upgradable installs" and the guard non-upgradable installs before attempting replacement decision: when `BuildInfo.Version` is `unknown`, print guidance to install a released binary, make no network request, and return a non-zero error. Verified by a unit test asserting no network seam is invoked and a non-nil error is returned for a dev build.
- [x] 3.2 Implement the package-manager / non-writable path guard: resolve the executable path and refuse with `brew upgrade`-style guidance when the path is not writable or under a known package-manager prefix, without attempting replacement. Verified by a unit test driving the path-check seam with a non-writable path and asserting no replacement attempt.

## 4. Release check and version comparison

- [x] 4.1 Implement release discovery and version comparison satisfying "Version comparison" and the version source and comparison decision: read the running version from `BuildInfo`, query the latest `benny123tw/bumpkin` release, and decide to update only when the latest is strictly newer by semver. Verified by table-driven tests over the newer/equal/older version pairs from the spec example.
- [x] 4.2 Configure the asset matcher for GoReleaser's binary naming so the host OS/arch raw-binary asset (`bumpkin-{version}-{os}-{arch}`) is selected and `checksums.txt` is excluded. Verified by a unit test asserting the matcher selects the correct asset name for representative os/arch pairs (including `darwin-arm64` and `linux-amd64`).
- [x] 4.3 Correct the production go-selfupdate filter so its non-empty `Filters` configuration uniquely matches the runtime OS/architecture instead of bypassing platform selection with a product-only regex, including GoReleaser's `.exe` suffix on Windows. Verify with a table-driven test that the configured regex accepts only the requested OS/architecture from a mixed-platform asset list and rejects `checksums.txt` and every incompatible binary.

## 5. Check, confirm, and apply

- [x] 5.1 Implement "Check-only mode": add the `--check` flag that prints current version, latest version, and availability, makes no disk changes, and exits 0. Verified by a test invoking `update --check` against a stubbed release and asserting output plus that the replace seam is never called.
- [x] 5.2 Implement "Confirmation prompt" per the flags and confirmation UX decision: add `--yes`/`-y`, prompt on stdin (`cmd.InOrStdin()`) when absent, abort (exit 0) on decline, and proceed on accept or with `--yes`. Verified by tests feeding "y"/"n" to the prompt and asserting replace-called vs. aborted, plus that `--yes` skips the prompt.
- [x] 5.3 Implement "Self-update command" apply path and "Integrity verification and safe failure": on confirmed newer release, download via go-selfupdate, verify against `checksums.txt`, atomically replace the executable, and print old → new version; on any failure leave the binary untouched and return a wrapped error. Verified by tests over the stubbed update seam asserting success output and that an injected download/verify error yields a non-zero error with no binary change.

## 6. Validation

- [x] 6.1 Run `just check` (tests + lint) and confirm `bumpkin update --help` shows the command, `upgrade` alias, and `--check`/`--yes` flags. Verified by a clean `just check` run and the help output.

## 7. Review follow-up and CI security

- [x] 7.1 Resolve the executable symlink to its canonical target before the package-manager guard and replacement so managed installations cannot be overwritten through a launcher symlink. Verify with a test whose executable seam returns a symlink path and asserts both the guard and replacement receive only the canonical target.
- [x] 7.2 Distinguish an existing GitHub release with no matching host asset from an absent release: the former returns a descriptive non-zero error while the latter retains the existing no-release result. Verify both states with HTTP-backed source tests and command-level assertions.
- [x] 7.3 Implement "Use the GitHub API with a focused replacement primitive": replace the reachable `go-selfupdate` OpenPGP dependency with direct GitHub release discovery, SHA-256 checksum verification, and the narrow `minio/selfupdate` replacement primitive; upgrade `go-git`, `go-billy`, and `x/text` to fixed versions. Verify with `just check`, `govulncheck ./...` on a patched supported Go toolchain, and a live `bumpkin update --check` release lookup.
