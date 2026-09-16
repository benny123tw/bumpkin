package cli

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/creativeprojects/go-selfupdate"
	"github.com/spf13/cobra"

	"github.com/benny123tw/bumpkin/internal/version"
)

const (
	repoOwner = "benny123tw"
	repoName  = "bumpkin"

	// devVersion is the placeholder injected for non-release builds.
	devVersion = "unknown"

	// checksumsFile is the SHA-256 checksum asset published alongside binaries.
	checksumsFile = "checksums.txt"
)

// packageManagerMarkers are path fragments that indicate the binary is managed
// by a package manager and must not be self-replaced.
var packageManagerMarkers = []string{
	"/Cellar/",          // Homebrew (Intel default and casks)
	"/opt/homebrew/",    // Homebrew (Apple Silicon)
	"/home/linuxbrew/",  // Linuxbrew
	"/nix/store/",       // Nix
	"/var/lib/flatpak/", // Flatpak
	"/snap/",            // Snap
}

// release describes the latest available release. The real source stores the
// underlying go-selfupdate handle in `handle`; test fakes leave it nil.
type release struct {
	version   string
	assetName string
	handle    *selfupdate.Release
}

// updateSource abstracts GitHub release discovery and binary replacement so
// tests can exercise the command without network access or touching the binary.
type updateSource interface {
	// Latest returns the latest release, or found=false when none matches.
	Latest(ctx context.Context) (rel *release, found bool, err error)
	// Replace downloads rel's asset, verifies it, and replaces the executable at exePath.
	Replace(ctx context.Context, rel *release, exePath string) error
}

type updateCommand struct {
	cmd  *cobra.Command
	info BuildInfo

	// Seams for testing.
	source     updateSource
	executable func() (string, error)
	blockPath  func(path string) (reason string, blocked bool)
}

// newUpdateCommand creates the self-update command. It updates the running
// bumpkin binary to the latest GitHub release, with a --check dry-run and a
// --yes confirmation bypass. The `upgrade` alias behaves identically.
func newUpdateCommand(info BuildInfo) *updateCommand {
	c := &updateCommand{
		info:       info,
		executable: os.Executable,
		blockPath:  defaultPathBlocker,
	}

	updateCmd := &cobra.Command{
		Use:     "update",
		Aliases: []string{"upgrade"},
		Short:   "Update bumpkin to the latest release",
		Long: `Update bumpkin to the latest release published on GitHub.

Checks the latest GitHub release for benny123tw/bumpkin, and if a newer
version is available, downloads the binary for your platform, verifies it
against the published checksums, and replaces the running executable.

Use --check to see whether an update is available without installing it.`,
		Example: `  bumpkin update --check
  bumpkin update --yes
  bumpkin upgrade --yes`,
		Args: cobra.NoArgs,
		RunE: c.execute,
	}

	updateCmd.Flags().Bool("check", false, "Check for an update without installing it")
	updateCmd.Flags().BoolP("yes", "y", false, "Skip the confirmation prompt")

	c.cmd = updateCmd
	return c
}

func (c *updateCommand) execute(cmd *cobra.Command, _ []string) error {
	out := cmd.OutOrStdout()
	checkOnly, _ := cmd.Flags().GetBool("check")
	assumeYes, _ := cmd.Flags().GetBool("yes")

	// Guard: development builds have no release version to compare against.
	if isDevBuild(c.info.Version) {
		fmt.Fprintln(out, "Self-update is not available for development builds.")
		fmt.Fprintf(
			out,
			"Install a released binary from https://github.com/%s/%s/releases,\n",
			repoOwner,
			repoName,
		)
		fmt.Fprintf(
			out,
			"or run: go install github.com/%s/%s/cmd/bumpkin@latest\n",
			repoOwner,
			repoName,
		)
		return NewExitError(ExitGeneralError, "self-update unavailable for development builds", nil)
	}

	// Resolve the running executable path.
	exePath, err := c.executable()
	if err != nil {
		return fmt.Errorf("could not locate executable: %w", err)
	}

	// Guard: package-manager or non-writable installs cannot be self-replaced.
	// Skipped for --check, which never writes to disk.
	if !checkOnly {
		if reason, blocked := c.blockPath(exePath); blocked {
			fmt.Fprintf(out, "Cannot self-update: %s.\n", reason)
			fmt.Fprintln(
				out,
				"If you installed bumpkin with a package manager, update it there (e.g. `brew upgrade bumpkin`).",
			)
			return nil
		}
	}

	src := c.source
	if src == nil {
		src = githubSource{}
	}

	rel, found, err := src.Latest(cmd.Context())
	if err != nil {
		return fmt.Errorf("failed to check for updates: %w", err)
	}
	if !found || rel == nil {
		fmt.Fprintln(out, "No release found to update to.")
		return nil
	}

	available, err := updateAvailable(c.info.Version, rel.version)
	if err != nil {
		return err
	}

	if checkOnly {
		fmt.Fprintf(out, "Current version: %s\n", c.info.Version)
		fmt.Fprintf(out, "Latest version:  %s\n", rel.version)
		if available {
			fmt.Fprintln(out, "An update is available. Run `bumpkin update` to install it.")
		} else {
			fmt.Fprintln(out, "bumpkin is up to date.")
		}
		return nil
	}

	if !available {
		fmt.Fprintf(out, "bumpkin %s is already the latest version.\n", c.info.Version)
		return nil
	}

	if !assumeYes && !promptConfirm(cmd.InOrStdin(), out, c.info.Version, rel.version) {
		fmt.Fprintln(out, "Update cancelled.")
		return nil
	}

	if err := src.Replace(cmd.Context(), rel, exePath); err != nil {
		return fmt.Errorf("update failed: %w", err)
	}

	fmt.Fprintf(out, "Updated bumpkin %s -> %s\n", c.info.Version, rel.version)
	return nil
}

// isDevBuild reports whether v is a development build with no comparable
// release version (empty, the "unknown"/"dev" placeholder, or unparseable).
func isDevBuild(v string) bool {
	switch strings.TrimSpace(v) {
	case "", devVersion, "dev":
		return true
	}
	if _, err := version.Parse(v); err != nil {
		return true
	}
	return false
}

// updateAvailable reports whether latest is strictly newer than current per
// semantic-version ordering.
func updateAvailable(current, latest string) (bool, error) {
	cur, err := version.Parse(current)
	if err != nil {
		return false, fmt.Errorf("invalid current version %q: %w", current, err)
	}
	lat, err := version.Parse(latest)
	if err != nil {
		return false, fmt.Errorf("invalid latest version %q: %w", latest, err)
	}
	return cur.LessThan(lat), nil
}

// selectAsset returns the asset built for goos/goarch from the given release
// asset names, applying the same os/arch tokens GoReleaser emits. It documents
// and verifies the naming contract the updater's Filters rely on.
func selectAsset(assets []string, goos, goarch string) (string, bool) {
	re := regexp.MustCompile(platformAssetFilter(goos, goarch))
	for _, a := range assets {
		if re.MatchString(a) {
			return a, true
		}
	}
	return "", false
}

// promptConfirm asks the user to confirm the update, returning true only for an
// affirmative answer. Input is read from in and prompts are written to out.
func promptConfirm(in io.Reader, out io.Writer, current, latest string) bool {
	fmt.Fprintf(out, "Update available: %s -> %s\n", current, latest)
	fmt.Fprint(out, "Do you want to update now? [y/N]: ")

	line, _ := bufio.NewReader(in).ReadString('\n')
	switch strings.ToLower(strings.TrimSpace(line)) {
	case "y", "yes":
		return true
	default:
		return false
	}
}

// defaultPathBlocker refuses installs that cannot be safely self-replaced: those
// under a known package-manager prefix, or whose directory is not writable
// (self-update writes a temp file alongside the binary before renaming).
func defaultPathBlocker(exePath string) (string, bool) {
	for _, marker := range packageManagerMarkers {
		if strings.Contains(exePath, marker) {
			return fmt.Sprintf("%s is managed by a package manager", exePath), true
		}
	}

	dir := filepath.Dir(exePath)
	tmp, err := os.CreateTemp(dir, ".bumpkin-update-*")
	if err != nil {
		return fmt.Sprintf("the binary at %s is not writable", exePath), true
	}
	name := tmp.Name()
	_ = tmp.Close()
	_ = os.Remove(name)
	return "", false
}

// githubSource is the production updateSource backed by go-selfupdate.
type githubSource struct{}

// platformAssetFilter returns a filter that uniquely matches the raw
// GoReleaser binary for a target OS and architecture. A custom go-selfupdate
// filter replaces the library's built-in platform matching, so the platform
// must be part of the expression.
func platformAssetFilter(goos, goarch string) string {
	extension := ""
	if goos == "windows" {
		extension = `\.exe`
	}

	return fmt.Sprintf(
		`^bumpkin-.*-%s-%s%s$`,
		regexp.QuoteMeta(goos),
		regexp.QuoteMeta(goarch),
		extension,
	)
}

func updaterConfig(goos, goarch string) selfupdate.Config {
	return selfupdate.Config{
		Validator: &selfupdate.ChecksumValidator{UniqueFilename: checksumsFile},
		Filters:   []string{platformAssetFilter(goos, goarch)},
	}
}

func newUpdater() (*selfupdate.Updater, error) {
	return selfupdate.NewUpdater(updaterConfig(runtime.GOOS, runtime.GOARCH))
}

func (githubSource) Latest(ctx context.Context) (*release, bool, error) {
	updater, err := newUpdater()
	if err != nil {
		return nil, false, fmt.Errorf("failed to create updater: %w", err)
	}

	rel, found, err := updater.DetectLatest(ctx, selfupdate.NewRepositorySlug(repoOwner, repoName))
	if err != nil {
		return nil, false, fmt.Errorf("failed to detect latest release: %w", err)
	}
	if !found || rel == nil {
		return nil, false, nil
	}

	return &release{version: rel.Version(), assetName: rel.AssetName, handle: rel}, true, nil
}

func (githubSource) Replace(ctx context.Context, rel *release, exePath string) error {
	updater, err := newUpdater()
	if err != nil {
		return fmt.Errorf("failed to create updater: %w", err)
	}
	if err := updater.UpdateTo(ctx, rel.handle, exePath); err != nil {
		return fmt.Errorf("failed to replace binary: %w", err)
	}
	return nil
}
