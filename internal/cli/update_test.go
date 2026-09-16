package cli

import (
	"bytes"
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// fakeSource is a test double for updateSource that records calls and never
// touches the network or the filesystem.
type fakeSource struct {
	rel        *release
	found      bool
	latestErr  error
	replaceErr error

	latestCalls  int
	replaceCalls int
	replacePath  string
}

func (f *fakeSource) Latest(_ context.Context) (*release, bool, error) {
	f.latestCalls++
	return f.rel, f.found, f.latestErr
}

func (f *fakeSource) Replace(_ context.Context, _ *release, exePath string) error {
	f.replaceCalls++
	f.replacePath = exePath
	return f.replaceErr
}

// newTestUpdateCommand builds an updateCommand wired with test seams: a fake
// source, a fixed executable path, and an unblocked path guard by default.
func newTestUpdateCommand(t *testing.T, version string, src updateSource) *updateCommand {
	t.Helper()
	c := newUpdateCommand(BuildInfo{Version: version})
	c.source = src
	c.executable = func() (string, error) { return "/tmp/bumpkin", nil }
	c.resolvePath = func(path string) (string, error) { return path, nil }
	c.blockPath = func(string) (string, bool) { return "", false }
	return c
}

func runUpdate(t *testing.T, c *updateCommand, stdin string, args ...string) (string, error) {
	t.Helper()
	buf := new(bytes.Buffer)
	c.cmd.SetOut(buf)
	c.cmd.SetErr(buf)
	c.cmd.SetIn(strings.NewReader(stdin))
	c.cmd.SetArgs(args)
	err := c.cmd.Execute()
	return buf.String(), err
}

// Task 2.1 / 2.2: command and upgrade alias resolve from the root command.
func TestUpdateCommand_Registration(t *testing.T) {
	for _, name := range []string{"update", "upgrade"} {
		cmd := NewRootCmd(testBuildInfo())
		found, _, err := cmd.Find([]string{name})
		require.NoError(t, err)
		assert.Equal(t, "update", found.Name(), "%q should resolve to the update command", name)
	}
}

func TestUpdateCommand_Help(t *testing.T) {
	buf := new(bytes.Buffer)
	cmd := NewRootCmd(testBuildInfo())
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"update", "--help"})

	require.NoError(t, cmd.Execute())

	out := buf.String()
	assert.Contains(t, out, "update")
	assert.Contains(t, out, "upgrade") // alias listed
	assert.Contains(t, out, "--check")
	assert.Contains(t, out, "--yes")
}

func TestUpdateCommand_ExamplesIncludeUpgradeAlias(t *testing.T) {
	c := newUpdateCommand(testBuildInfo())

	assert.Contains(t, c.cmd.Example, "bumpkin upgrade --yes")
}

// Task 4 (3.1): dev-build guard short-circuits before any network call.
func TestUpdateCommand_DevBuildGuard(t *testing.T) {
	for _, v := range []string{"unknown", "", "dev"} {
		src := &fakeSource{}
		c := newTestUpdateCommand(t, v, src)

		out, err := runUpdate(t, c, "")

		require.Error(t, err, "dev build %q should error", v)
		assert.Equal(t, ExitGeneralError, GetExitCode(err))
		assert.Equal(t, 0, src.latestCalls, "no network call for dev build %q", v)
		assert.Contains(t, out, "development builds")
	}
}

// Task 5 (3.2): package-manager / non-writable path guard refuses without replacing.
func TestUpdateCommand_PathGuard(t *testing.T) {
	src := &fakeSource{
		found: true,
		rel:   &release{version: "1.2.0"},
	}
	c := newTestUpdateCommand(t, "1.0.0", src)
	c.blockPath = func(string) (string, bool) {
		return "the binary at /usr/bin/bumpkin is not writable", true
	}

	out, err := runUpdate(t, c, "")

	require.NoError(t, err)
	assert.Equal(t, 0, src.replaceCalls, "must not replace a blocked binary")
	assert.Contains(t, out, "Cannot self-update")
	assert.Contains(t, out, "brew upgrade")
}

func TestUpdateCommand_ResolvesExecutableBeforeGuardAndReplace(t *testing.T) {
	src := &fakeSource{found: true, rel: &release{version: "1.2.0"}}
	c := newTestUpdateCommand(t, "1.0.0", src)
	c.executable = func() (string, error) { return "/usr/local/bin/bumpkin", nil }
	c.resolvePath = func(path string) (string, error) {
		assert.Equal(t, "/usr/local/bin/bumpkin", path)
		return "/opt/homebrew/Cellar/bumpkin/1.2.0/bin/bumpkin", nil
	}

	var guardedPath string
	c.blockPath = func(path string) (string, bool) {
		guardedPath = path
		return "", false
	}

	_, err := runUpdate(t, c, "", "--yes")

	require.NoError(t, err)
	assert.Equal(t, "/opt/homebrew/Cellar/bumpkin/1.2.0/bin/bumpkin", guardedPath)
	assert.Equal(t, guardedPath, src.replacePath)
}

func TestDefaultPathBlocker_PackageManager(t *testing.T) {
	reason, blocked := defaultPathBlocker("/opt/homebrew/bin/bumpkin")
	assert.True(t, blocked)
	assert.Contains(t, reason, "package manager")
}

func TestDefaultPathBlocker_WritableDir(t *testing.T) {
	// A path in a writable temp dir is not blocked.
	exe := t.TempDir() + "/bumpkin"
	_, blocked := defaultPathBlocker(exe)
	assert.False(t, blocked)
}

// Task 6 (4.1): version comparison decides update strictly by semver.
func TestUpdateAvailable(t *testing.T) {
	cases := []struct {
		current, latest string
		want            bool
	}{
		{"v1.2.0", "v1.3.0", true},  // newer -> update
		{"v1.3.0", "v1.3.0", false}, // equal -> up to date
		{"v1.4.0", "v1.3.0", false}, // older latest -> up to date
		{"1.2.0", "1.3.0", true},    // works without the v prefix
	}
	for _, tc := range cases {
		got, err := updateAvailable(tc.current, tc.latest)
		require.NoError(t, err)
		assert.Equalf(t, tc.want, got, "updateAvailable(%q, %q)", tc.current, tc.latest)
	}
}

// Task 7 (4.2): asset matcher selects the host os/arch binary and excludes checksums.txt.
func TestSelectAsset(t *testing.T) {
	assets := []string{
		"bumpkin-1.3.0-darwin-amd64",
		"bumpkin-1.3.0-darwin-arm64",
		"bumpkin-1.3.0-linux-amd64",
		"bumpkin-1.3.0-linux-arm64",
		"bumpkin-1.3.0-windows-amd64.exe",
		"checksums.txt",
	}

	cases := []struct {
		goos, goarch, want string
	}{
		{"darwin", "arm64", "bumpkin-1.3.0-darwin-arm64"},
		{"linux", "amd64", "bumpkin-1.3.0-linux-amd64"},
		{"windows", "amd64", "bumpkin-1.3.0-windows-amd64.exe"},
	}
	for _, tc := range cases {
		got, ok := selectAsset(assets, tc.goos, tc.goarch)
		require.Truef(t, ok, "expected a match for %s/%s", tc.goos, tc.goarch)
		assert.Equal(t, tc.want, got)
		assert.NotEqual(t, checksumsFile, got, "checksums.txt must never be selected")
	}

	_, ok := selectAsset(assets, "plan9", "mips")
	assert.False(t, ok, "unknown os/arch should not match")
}

func TestPlatformAssetFilter(t *testing.T) {
	assets := []string{
		"bumpkin-1.3.0-darwin-amd64",
		"bumpkin-1.3.0-darwin-arm64",
		"bumpkin-1.3.0-linux-amd64",
		"bumpkin-1.3.0-linux-arm64",
		"bumpkin-1.3.0-windows-amd64.exe",
		"checksums.txt",
	}

	cases := []struct {
		goos, goarch, want string
	}{
		{"darwin", "arm64", "bumpkin-1.3.0-darwin-arm64"},
		{"linux", "amd64", "bumpkin-1.3.0-linux-amd64"},
		{"windows", "amd64", "bumpkin-1.3.0-windows-amd64.exe"},
	}

	for _, tc := range cases {
		t.Run(tc.goos+"/"+tc.goarch, func(t *testing.T) {
			filter := regexp.MustCompile(platformAssetFilter(tc.goos, tc.goarch))
			for _, asset := range assets {
				assert.Equalf(t, asset == tc.want, filter.MatchString(asset), "asset %q", asset)
			}
		})
	}
}

func TestGitHubSourceLatest_NoRelease(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	t.Cleanup(server.Close)

	src := githubSource{
		client:    server.Client(),
		latestURL: server.URL,
		goos:      "linux",
		goarch:    "arm64",
	}

	rel, found, err := src.Latest(context.Background())

	require.NoError(t, err)
	assert.False(t, found)
	assert.Nil(t, rel)
}

func TestGitHubSourceLatest_NoMatchingPlatformAsset(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{
  "tag_name": "v1.3.0",
  "assets": [
    {"name": "bumpkin-1.3.0-darwin-arm64", "browser_download_url": "https://example.invalid/darwin"},
    {"name": "checksums.txt", "browser_download_url": "https://example.invalid/checksums"}
  ]
}`)
	}))
	t.Cleanup(server.Close)

	src := githubSource{
		client:    server.Client(),
		latestURL: server.URL,
		goos:      "linux",
		goarch:    "arm64",
	}

	rel, found, err := src.Latest(context.Background())

	require.Error(t, err)
	assert.Contains(t, err.Error(), "no release asset for linux/arm64")
	assert.False(t, found)
	assert.Nil(t, rel)
}

func TestGitHubSourceLatest_SelectsPlatformAsset(t *testing.T) {
	var server *httptest.Server
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{
  "tag_name": "v1.3.0",
  "assets": [
    {"name": "bumpkin-1.3.0-darwin-arm64", "browser_download_url": %q},
    {"name": "bumpkin-1.3.0-linux-arm64", "browser_download_url": %q},
    {"name": "checksums.txt", "browser_download_url": %q}
  ]
}`, server.URL+"/darwin", server.URL+"/binary", server.URL+"/checksums")
	}))
	t.Cleanup(server.Close)

	src := githubSource{
		client:    server.Client(),
		latestURL: server.URL,
		goos:      "linux",
		goarch:    "arm64",
	}

	rel, found, err := src.Latest(context.Background())

	require.NoError(t, err)
	require.True(t, found)
	require.NotNil(t, rel)
	assert.Equal(t, "v1.3.0", rel.version)
	assert.Equal(t, "bumpkin-1.3.0-linux-arm64", rel.assetName)
	assert.Equal(t, server.URL+"/binary", rel.assetURL)
	assert.Equal(t, server.URL+"/checksums", rel.checksumURL)
}

func TestGitHubSourceReplace_ChecksumMismatchLeavesBinaryUntouched(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/binary":
			fmt.Fprint(w, "new binary")
		case "/checksums":
			fmt.Fprintln(w, strings.Repeat("0", sha256.Size*2), " bumpkin-1.3.0-linux-amd64")
		default:
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(server.Close)

	exePath := filepath.Join(t.TempDir(), "bumpkin")
	//nolint:gosec // The updater requires an executable test fixture.
	require.NoError(t, os.WriteFile(exePath, []byte("old binary"), 0o755))
	src := githubSource{client: server.Client()}
	rel := &release{
		assetName:   "bumpkin-1.3.0-linux-amd64",
		assetURL:    server.URL + "/binary",
		checksumURL: server.URL + "/checksums",
	}

	err := src.Replace(context.Background(), rel, exePath)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "checksum")
	contents, readErr := os.ReadFile(exePath)
	require.NoError(t, readErr)
	assert.Equal(t, "old binary", string(contents))
}

// Task 8 (5.1): --check reports status without replacing.
func TestUpdateCommand_CheckAvailable(t *testing.T) {
	src := &fakeSource{found: true, rel: &release{version: "1.2.0"}}
	c := newTestUpdateCommand(t, "1.0.0", src)

	out, err := runUpdate(t, c, "", "--check")

	require.NoError(t, err)
	assert.Equal(t, 0, src.replaceCalls, "--check must not replace")
	assert.Contains(t, out, "Current version: 1.0.0")
	assert.Contains(t, out, "Latest version:  1.2.0")
	assert.Contains(t, out, "update is available")
}

func TestUpdateCommand_CheckUpToDate(t *testing.T) {
	src := &fakeSource{found: true, rel: &release{version: "1.0.0"}}
	c := newTestUpdateCommand(t, "1.0.0", src)

	out, err := runUpdate(t, c, "", "--check")

	require.NoError(t, err)
	assert.Equal(t, 0, src.replaceCalls)
	assert.Contains(t, out, "up to date")
}

// Task 9 (5.2): confirmation prompt gates replacement; --yes skips it.
func TestUpdateCommand_ConfirmAccept(t *testing.T) {
	src := &fakeSource{found: true, rel: &release{version: "1.2.0"}}
	c := newTestUpdateCommand(t, "1.0.0", src)

	_, err := runUpdate(t, c, "y\n")

	require.NoError(t, err)
	assert.Equal(t, 1, src.replaceCalls, "affirmative answer should replace")
}

func TestUpdateCommand_ConfirmDecline(t *testing.T) {
	src := &fakeSource{found: true, rel: &release{version: "1.2.0"}}
	c := newTestUpdateCommand(t, "1.0.0", src)

	out, err := runUpdate(t, c, "n\n")

	require.NoError(t, err)
	assert.Equal(t, 0, src.replaceCalls, "declining should not replace")
	assert.Contains(t, out, "cancelled")
}

func TestUpdateCommand_YesSkipsPrompt(t *testing.T) {
	src := &fakeSource{found: true, rel: &release{version: "1.2.0"}}
	c := newTestUpdateCommand(t, "1.0.0", src)

	// Empty stdin: replacement must still happen because --yes bypasses the prompt.
	_, err := runUpdate(t, c, "", "--yes")

	require.NoError(t, err)
	assert.Equal(t, 1, src.replaceCalls)
}

// Task 10 (5.3): apply path prints old -> new on success; a failure surfaces an
// error and leaves the binary untouched.
func TestUpdateCommand_ApplySuccess(t *testing.T) {
	src := &fakeSource{found: true, rel: &release{version: "1.2.0"}}
	c := newTestUpdateCommand(t, "1.0.0", src)

	out, err := runUpdate(t, c, "", "--yes")

	require.NoError(t, err)
	assert.Equal(t, 1, src.replaceCalls)
	assert.Contains(t, out, "Updated bumpkin 1.0.0 -> 1.2.0")
}

func TestUpdateCommand_ApplyFailure(t *testing.T) {
	src := &fakeSource{
		found:      true,
		rel:        &release{version: "1.2.0"},
		replaceErr: errors.New("checksum mismatch"),
	}
	c := newTestUpdateCommand(t, "1.0.0", src)

	_, err := runUpdate(t, c, "", "--yes")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "update failed")
	assert.Equal(t, 1, src.replaceCalls, "replace was attempted")
}

func TestUpdateCommand_NoReleaseFound(t *testing.T) {
	src := &fakeSource{found: false}
	c := newTestUpdateCommand(t, "1.0.0", src)

	out, err := runUpdate(t, c, "", "--yes")

	require.NoError(t, err)
	assert.Equal(t, 0, src.replaceCalls)
	assert.Contains(t, out, "No release found")
}
