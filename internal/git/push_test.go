package git

import (
	"context"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// T031: Test for push to remote
func TestRepository_PushTag(t *testing.T) {
	// Create a bare repository to act as remote
	remoteDir := t.TempDir()
	runGit(t, remoteDir, "init", "--bare")

	// Create a local repository
	localDir := t.TempDir()
	initRealGitRepo(t, localDir)

	// Add the bare repo as remote
	runGit(t, localDir, "remote", "add", "origin", remoteDir)

	// Get current branch name and push
	branch := getCurrentBranch(t, localDir)
	runGit(t, localDir, "push", "-u", "origin", branch)

	// Create a tag
	repo, err := Open(localDir)
	require.NoError(t, err)

	err = repo.CreateTag("v1.0.0", "Release v1.0.0")
	require.NoError(t, err)

	// Push the tag
	err = repo.PushTag("v1.0.0", "origin")
	require.NoError(t, err)

	// Verify tag exists in remote by cloning to a new location
	cloneDir := t.TempDir()
	runGit(t, cloneDir, "clone", remoteDir, ".")

	// Check that tag exists in clone
	cloneRepo, err := Open(cloneDir)
	require.NoError(t, err)

	tags, err := cloneRepo.ListTags()
	require.NoError(t, err)

	tagNames := make([]string, len(tags))
	for i, tag := range tags {
		tagNames[i] = tag.Name
	}
	assert.Contains(t, tagNames, "v1.0.0")
}

func TestRepository_PushTag_NoRemote(t *testing.T) {
	tmpDir := t.TempDir()
	initRealGitRepo(t, tmpDir)

	repo, err := Open(tmpDir)
	require.NoError(t, err)

	err = repo.CreateTag("v1.0.0", "Release v1.0.0")
	require.NoError(t, err)

	// Try to push without a remote configured
	err = repo.PushTag("v1.0.0", "origin")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "remote")
}

func TestRepository_PushTag_TagNotFound(t *testing.T) {
	// Create a bare repository to act as remote
	remoteDir := t.TempDir()
	runGit(t, remoteDir, "init", "--bare")

	localDir := t.TempDir()
	initRealGitRepo(t, localDir)
	runGit(t, localDir, "remote", "add", "origin", remoteDir)

	repo, err := Open(localDir)
	require.NoError(t, err)

	// Try to push a tag that doesn't exist
	err = repo.PushTag("v999.0.0", "origin")
	assert.Error(t, err)
}

func TestRepository_HasRemote(t *testing.T) {
	t.Run("with remote", func(t *testing.T) {
		remoteDir := t.TempDir()
		runGit(t, remoteDir, "init", "--bare")

		localDir := t.TempDir()
		initRealGitRepo(t, localDir)
		runGit(t, localDir, "remote", "add", "origin", remoteDir)

		repo, err := Open(localDir)
		require.NoError(t, err)

		hasRemote, err := repo.HasRemote("origin")
		require.NoError(t, err)
		assert.True(t, hasRemote)
	})

	t.Run("without remote", func(t *testing.T) {
		localDir := t.TempDir()
		initRealGitRepo(t, localDir)

		repo, err := Open(localDir)
		require.NoError(t, err)

		hasRemote, err := repo.HasRemote("origin")
		require.NoError(t, err)
		assert.False(t, hasRemote)
	})
}

func TestRepository_GetRemoteURL(t *testing.T) {
	remoteDir := t.TempDir()
	runGit(t, remoteDir, "init", "--bare")

	localDir := t.TempDir()
	initRealGitRepo(t, localDir)
	runGit(t, localDir, "remote", "add", "origin", remoteDir)

	repo, err := Open(localDir)
	require.NoError(t, err)

	url, err := repo.GetRemoteURL("origin")
	require.NoError(t, err)

	// The URL should contain the remote directory path
	assert.Contains(t, url, filepath.Base(remoteDir))
}

// Helper to get current branch name
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
