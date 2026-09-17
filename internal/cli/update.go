package cli

import (
	"bufio"
	"context"
	"crypto"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	minioselfupdate "github.com/minio/selfupdate"
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

	githubAPIBaseURL = "https://api.github.com"
	maxMetadataBytes = 4 << 20
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

// release describes the latest available release and its downloadable assets.
type release struct {
	version     string
	assetName   string
	assetURL    string
	checksumURL string
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
	source      updateSource
	executable  func() (string, error)
	resolvePath func(string) (string, error)
	blockPath   func(path string) (reason string, blocked bool)
}

// newUpdateCommand creates the self-update command. It updates the running
// bumpkin binary to the latest GitHub release, with a --check dry-run and a
// --yes confirmation bypass. The `upgrade` alias behaves identically.
func newUpdateCommand(info BuildInfo) *updateCommand {
	c := &updateCommand{
		info:        info,
		executable:  os.Executable,
		resolvePath: filepath.EvalSymlinks,
		blockPath:   defaultPathBlocker,
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
		exePath, err = c.resolvePath(exePath)
		if err != nil {
			return fmt.Errorf("could not resolve executable path: %w", err)
		}

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

// githubSource is the production updateSource backed by GitHub's Releases API.
type githubSource struct {
	client    *http.Client
	latestURL string
	goos      string
	goarch    string
}

// platformAssetFilter returns a filter that uniquely matches the raw
// GoReleaser binary for a target OS and architecture.
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

func (s githubSource) effectiveClient() *http.Client {
	if s.client != nil {
		return s.client
	}
	return http.DefaultClient
}

func (s githubSource) effectivePlatform() (string, string) {
	goos := s.goos
	if goos == "" {
		goos = runtime.GOOS
	}
	goarch := s.goarch
	if goarch == "" {
		goarch = runtime.GOARCH
	}
	return goos, goarch
}

func (s githubSource) effectiveLatestURL() string {
	if s.latestURL != "" {
		return s.latestURL
	}
	return fmt.Sprintf("%s/repos/%s/%s/releases/latest", githubAPIBaseURL, repoOwner, repoName)
}

func (s githubSource) Latest(ctx context.Context) (*release, bool, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, s.effectiveLatestURL(), nil)
	if err != nil {
		return nil, false, fmt.Errorf("failed to create latest-release request: %w", err)
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", repoName)

	resp, err := s.effectiveClient().Do(req)
	if err != nil {
		return nil, false, fmt.Errorf("failed to fetch latest release: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, false, nil
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, false, fmt.Errorf(
			"GitHub latest-release request returned HTTP %d",
			resp.StatusCode,
		)
	}

	var payload struct {
		TagName string `json:"tag_name"`
		Assets  []struct {
			Name string `json:"name"`
			URL  string `json:"browser_download_url"`
		} `json:"assets"`
	}
	metadata := io.LimitReader(resp.Body, maxMetadataBytes)
	if err := json.NewDecoder(metadata).Decode(&payload); err != nil {
		return nil, false, fmt.Errorf("failed to decode latest release: %w", err)
	}

	goos, goarch := s.effectivePlatform()
	assetPattern := regexp.MustCompile(platformAssetFilter(goos, goarch))
	var assetName, assetURL, checksumURL string
	for _, asset := range payload.Assets {
		switch {
		case assetPattern.MatchString(asset.Name):
			assetName = asset.Name
			assetURL = asset.URL
		case asset.Name == checksumsFile:
			checksumURL = asset.URL
		}
	}

	if assetURL == "" {
		return nil, false, fmt.Errorf("no release asset for %s/%s", goos, goarch)
	}
	if checksumURL == "" {
		return nil, false, fmt.Errorf("release is missing %s", checksumsFile)
	}
	if strings.TrimSpace(payload.TagName) == "" {
		return nil, false, fmt.Errorf("latest release has no tag name")
	}

	return &release{
		version:     payload.TagName,
		assetName:   assetName,
		assetURL:    assetURL,
		checksumURL: checksumURL,
	}, true, nil
}

func (s githubSource) Replace(ctx context.Context, rel *release, exePath string) error {
	if rel == nil {
		return fmt.Errorf("invalid release")
	}

	checksumBody, err := s.download(ctx, rel.checksumURL)
	if err != nil {
		return fmt.Errorf("failed to download %s: %w", checksumsFile, err)
	}
	checksumData, err := io.ReadAll(io.LimitReader(checksumBody, maxMetadataBytes))
	closeErr := checksumBody.Close()
	if err != nil {
		return fmt.Errorf("failed to read %s: %w", checksumsFile, err)
	}
	if closeErr != nil {
		return fmt.Errorf("failed to close %s response: %w", checksumsFile, closeErr)
	}

	checksum, err := checksumForAsset(checksumData, rel.assetName)
	if err != nil {
		return err
	}

	binaryBody, err := s.download(ctx, rel.assetURL)
	if err != nil {
		return fmt.Errorf("failed to download %s: %w", rel.assetName, err)
	}
	defer binaryBody.Close()

	if err := minioselfupdate.Apply(binaryBody, minioselfupdate.Options{
		TargetPath: exePath,
		Checksum:   checksum,
		Hash:       crypto.SHA256,
	}); err != nil {
		return fmt.Errorf("failed to replace binary: %w", err)
	}
	return nil
}

func (s githubSource) download(ctx context.Context, url string) (io.ReadCloser, error) {
	if strings.TrimSpace(url) == "" {
		return nil, fmt.Errorf("download URL is empty")
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/octet-stream")
	req.Header.Set("User-Agent", repoName)

	resp, err := s.effectiveClient().Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		resp.Body.Close()
		return nil, fmt.Errorf("download returned HTTP %d", resp.StatusCode)
	}
	return resp.Body, nil
}

func checksumForAsset(content []byte, assetName string) ([]byte, error) {
	for _, line := range strings.Split(string(content), "\n") {
		fields := strings.Fields(line)
		if len(fields) != 2 || strings.TrimPrefix(fields[1], "*") != assetName {
			continue
		}

		checksum, err := hex.DecodeString(fields[0])
		if err != nil || len(checksum) != crypto.SHA256.Size() {
			return nil, fmt.Errorf("invalid SHA-256 checksum for %s", assetName)
		}
		return checksum, nil
	}
	return nil, fmt.Errorf("checksum for %s not found", assetName)
}
