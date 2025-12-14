package git

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// T021: Test for listing all tags
func TestRepository_ListTags(t *testing.T) {
	tmpDir := t.TempDir()
	initRealGitRepo(t, tmpDir)

	// Create some tags
	createTag(t, tmpDir, "v1.0.0", "Release 1.0.0")
	createTag(t, tmpDir, "v1.1.0", "Release 1.1.0")
	createTag(t, tmpDir, "v2.0.0", "Release 2.0.0")

	repo, err := Open(tmpDir)
	require.NoError(t, err)

	tags, err := repo.ListTags()
	require.NoError(t, err)

	assert.Len(t, tags, 3)

	// Check tag names exist
	tagNames := make([]string, len(tags))
	for i, tag := range tags {
		tagNames[i] = tag.Name
	}
	assert.Contains(t, tagNames, "v1.0.0")
	assert.Contains(t, tagNames, "v1.1.0")
	assert.Contains(t, tagNames, "v2.0.0")
}

func TestRepository_ListTags_Empty(t *testing.T) {
	tmpDir := t.TempDir()
	initRealGitRepo(t, tmpDir)

	repo, err := Open(tmpDir)
	require.NoError(t, err)

	tags, err := repo.ListTags()
	require.NoError(t, err)
	assert.Empty(t, tags)
}

// T022: Test for finding latest semver tag
func TestRepository_LatestTag(t *testing.T) {
	tmpDir := t.TempDir()
	initRealGitRepo(t, tmpDir)

	// Create tags in non-chronological order
	createTag(t, tmpDir, "v1.0.0", "Release 1.0.0")
	createCommit(t, tmpDir, "second commit")
	createTag(t, tmpDir, "v2.0.0", "Release 2.0.0")
	createCommit(t, tmpDir, "third commit")
	createTag(t, tmpDir, "v1.5.0", "Release 1.5.0")

	repo, err := Open(tmpDir)
	require.NoError(t, err)

	tag, err := repo.LatestTag("v")
	require.NoError(t, err)
	require.NotNil(t, tag)

	// Should return v2.0.0 as it's the highest semver
	assert.Equal(t, "v2.0.0", tag.Name)
}

func TestRepository_LatestTag_NoTags(t *testing.T) {
	tmpDir := t.TempDir()
	initRealGitRepo(t, tmpDir)

	repo, err := Open(tmpDir)
	require.NoError(t, err)

	tag, err := repo.LatestTag("v")
	require.NoError(t, err)
	assert.Nil(t, tag)
}

func TestRepository_LatestTag_WithPrerelease(t *testing.T) {
	tmpDir := t.TempDir()
	initRealGitRepo(t, tmpDir)

	createTag(t, tmpDir, "v1.0.0", "Release 1.0.0")
	createCommit(t, tmpDir, "second commit")
	createTag(t, tmpDir, "v1.1.0-alpha.0", "Alpha release")
	createCommit(t, tmpDir, "third commit")
	createTag(t, tmpDir, "v1.0.1", "Patch release")

	repo, err := Open(tmpDir)
	require.NoError(t, err)

	tag, err := repo.LatestTag("v")
	require.NoError(t, err)
	require.NotNil(t, tag)

	// v1.1.0-alpha.0 < v1.1.0, but > v1.0.1, so it should be the latest
	assert.Equal(t, "v1.1.0-alpha.0", tag.Name)
}

// T024: Test for creating annotated tag
func TestRepository_CreateTag(t *testing.T) {
	tmpDir := t.TempDir()
	initRealGitRepo(t, tmpDir)

	repo, err := Open(tmpDir)
	require.NoError(t, err)

	err = repo.CreateTag("v1.0.0", "Release v1.0.0")
	require.NoError(t, err)

	// Verify tag was created
	tags, err := repo.ListTags()
	require.NoError(t, err)
	require.Len(t, tags, 1)
	assert.Equal(t, "v1.0.0", tags[0].Name)
	assert.True(t, tags[0].IsAnnotated)
}

func TestRepository_CreateTag_AlreadyExists(t *testing.T) {
	tmpDir := t.TempDir()
	initRealGitRepo(t, tmpDir)

	repo, err := Open(tmpDir)
	require.NoError(t, err)

	err = repo.CreateTag("v1.0.0", "Release v1.0.0")
	require.NoError(t, err)

	// Try to create same tag again
	err = repo.CreateTag("v1.0.0", "Release v1.0.0")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already exists")
}

// Helper to initialize a real git repo for testing
func initRealGitRepo(t *testing.T, dir string) {
	t.Helper()

	runGit(t, dir, "init")
	runGit(t, dir, "config", "user.email", "test@example.com")
	runGit(t, dir, "config", "user.name", "Test User")

	// Create initial commit
	testFile := filepath.Join(dir, "README.md")
	//nolint:gosec // test file
	require.NoError(t, os.WriteFile(testFile, []byte("# Test\n"), 0o644))
	runGit(t, dir, "add", ".")
	runGit(t, dir, "commit", "-m", "Initial commit")
}

// Helper to create a commit
func createCommit(t *testing.T, dir, message string) {
	t.Helper()

	testFile := filepath.Join(dir, "file-"+time.Now().Format("20060102150405.000"))
	//nolint:gosec // test file
	require.NoError(t, os.WriteFile(testFile, []byte(message+"\n"), 0o644))
	runGit(t, dir, "add", ".")
	runGit(t, dir, "commit", "-m", message)
}

// Helper to create an annotated tag
func createTag(t *testing.T, dir, name, message string) {
	t.Helper()
	runGit(t, dir, "tag", "-a", name, "-m", message)
}

// Helper to run git command
func runGit(t *testing.T, dir string, args ...string) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "git", args...)
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "git %v failed: %s", args, string(output))
}
