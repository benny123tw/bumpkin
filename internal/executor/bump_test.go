package executor

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/benny123tw/bumpkin/internal/git"
	"github.com/benny123tw/bumpkin/internal/version"
)

// T033: Test for executor with patch bump
func TestExecute_PatchBump(t *testing.T) {
	tmpDir := t.TempDir()
	initRealGitRepo(t, tmpDir)

	// Create initial tag
	createTag(t, tmpDir)
	createCommit(t, tmpDir, "fix: bug fix")

	repo, err := git.Open(tmpDir)
	require.NoError(t, err)

	result, err := Execute(context.Background(), Request{
		Repository: repo,
		BumpType:   version.BumpPatch,
		Prefix:     "v",
		NoPush:     true, // Don't push in tests
	})

	require.NoError(t, err)
	assert.Equal(t, "1.0.0", result.PreviousVersion)
	assert.Equal(t, "1.0.1", result.NewVersion)
	assert.Equal(t, "v1.0.1", result.TagName)
	assert.True(t, result.TagCreated)
	assert.False(t, result.Pushed)

	// Verify tag was actually created
	tags, err := repo.ListTags()
	require.NoError(t, err)

	tagNames := make([]string, len(tags))
	for i, tag := range tags {
		tagNames[i] = tag.Name
	}
	assert.Contains(t, tagNames, "v1.0.1")
}

// T034: Test for executor with dry-run mode
func TestExecute_DryRun(t *testing.T) {
	tmpDir := t.TempDir()
	initRealGitRepo(t, tmpDir)

	createTag(t, tmpDir)
	createCommit(t, tmpDir, "feat: new feature")

	repo, err := git.Open(tmpDir)
	require.NoError(t, err)

	result, err := Execute(context.Background(), Request{
		Repository: repo,
		BumpType:   version.BumpMinor,
		Prefix:     "v",
		DryRun:     true,
	})

	require.NoError(t, err)
	assert.Equal(t, "1.0.0", result.PreviousVersion)
	assert.Equal(t, "1.1.0", result.NewVersion)
	assert.Equal(t, "v1.1.0", result.TagName)
	assert.False(t, result.TagCreated) // Should NOT create tag in dry-run
	assert.False(t, result.Pushed)

	// Verify tag was NOT created
	tags, err := repo.ListTags()
	require.NoError(t, err)

	for _, tag := range tags {
		assert.NotEqual(t, "v1.1.0", tag.Name, "Tag should not be created in dry-run mode")
	}
}

// T036: Test for executor with no-push mode
func TestExecute_NoPush(t *testing.T) {
	// Create remote
	remoteDir := t.TempDir()
	runGit(t, remoteDir, "init", "--bare")

	// Create local repo with remote
	tmpDir := t.TempDir()
	initRealGitRepo(t, tmpDir)
	runGit(t, tmpDir, "remote", "add", "origin", remoteDir)

	branch := getCurrentBranch(t, tmpDir)
	runGit(t, tmpDir, "push", "-u", "origin", branch)

	createTag(t, tmpDir)
	runGit(t, tmpDir, "push", "origin", "v1.0.0")
	createCommit(t, tmpDir, "feat: major change")

	repo, err := git.Open(tmpDir)
	require.NoError(t, err)

	result, err := Execute(context.Background(), Request{
		Repository: repo,
		BumpType:   version.BumpMajor,
		Prefix:     "v",
		Remote:     "origin",
		NoPush:     true,
	})

	require.NoError(t, err)
	assert.Equal(t, "1.0.0", result.PreviousVersion)
	assert.Equal(t, "2.0.0", result.NewVersion)
	assert.True(t, result.TagCreated)
	assert.False(t, result.Pushed) // Should NOT push

	// Verify tag was created locally
	tags, err := repo.ListTags()
	require.NoError(t, err)

	tagNames := make([]string, len(tags))
	for i, tag := range tags {
		tagNames[i] = tag.Name
	}
	assert.Contains(t, tagNames, "v2.0.0")

	// Verify tag was NOT pushed to remote
	cloneDir := t.TempDir()
	runGit(t, cloneDir, "clone", remoteDir, ".")

	cloneRepo, err := git.Open(cloneDir)
	require.NoError(t, err)

	remoteTags, err := cloneRepo.ListTags()
	require.NoError(t, err)

	for _, tag := range remoteTags {
		assert.NotEqual(t, "v2.0.0", tag.Name, "Tag should not be pushed in no-push mode")
	}
}

func TestExecute_NoExistingTags(t *testing.T) {
	tmpDir := t.TempDir()
	initRealGitRepo(t, tmpDir)

	repo, err := git.Open(tmpDir)
	require.NoError(t, err)

	result, err := Execute(context.Background(), Request{
		Repository: repo,
		BumpType:   version.BumpMinor,
		Prefix:     "v",
		NoPush:     true,
	})

	require.NoError(t, err)
	assert.Equal(t, "0.0.0", result.PreviousVersion) // Start from 0.0.0
	assert.Equal(t, "0.1.0", result.NewVersion)
	assert.True(t, result.TagCreated)
}

func TestExecute_CustomVersion(t *testing.T) {
	tmpDir := t.TempDir()
	initRealGitRepo(t, tmpDir)

	createTag(t, tmpDir)
	createCommit(t, tmpDir, "feat: new feature")

	repo, err := git.Open(tmpDir)
	require.NoError(t, err)

	result, err := Execute(context.Background(), Request{
		Repository:    repo,
		BumpType:      version.BumpCustom,
		CustomVersion: "3.0.0",
		Prefix:        "v",
		NoPush:        true,
	})

	require.NoError(t, err)
	assert.Equal(t, "1.0.0", result.PreviousVersion)
	assert.Equal(t, "3.0.0", result.NewVersion)
	assert.Equal(t, "v3.0.0", result.TagName)
	assert.True(t, result.TagCreated)
}

// Helper functions

func initRealGitRepo(t *testing.T, dir string) {
	t.Helper()

	runGit(t, dir, "init")
	runGit(t, dir, "config", "user.email", "test@example.com")
	runGit(t, dir, "config", "user.name", "Test User")

	testFile := filepath.Join(dir, "README.md")
	//nolint:gosec // test file
	require.NoError(t, os.WriteFile(testFile, []byte("# Test\n"), 0o644))
	runGit(t, dir, "add", ".")
	runGit(t, dir, "commit", "-m", "Initial commit")
}

func createCommit(t *testing.T, dir, message string) {
	t.Helper()
	runGit(t, dir, "commit", "--allow-empty", "-m", message)
}

func createTag(t *testing.T, dir string) {
	t.Helper()
	runGit(t, dir, "tag", "-a", "v1.0.0", "-m", "Release 1.0.0")
}

func runGit(t *testing.T, dir string, args ...string) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "git", args...)
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "git %v failed: %s", args, string(output))
}

func getCurrentBranch(t *testing.T, dir string) string {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "git", "rev-parse", "--abbrev-ref", "HEAD")
	cmd.Dir = dir
	output, err := cmd.Output()
	require.NoError(t, err)
	return strings.TrimSpace(string(output))
}
