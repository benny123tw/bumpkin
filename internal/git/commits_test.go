package git

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// T026: Test for listing commits since tag
func TestRepository_GetCommitsSinceTag(t *testing.T) {
	tmpDir := t.TempDir()
	initRealGitRepo(t, tmpDir)

	// Create a tag
	createTag(t, tmpDir, "v1.0.0", "Release 1.0.0")

	// Create some commits after the tag
	createCommit(t, tmpDir, "feat: add new feature")
	createCommit(t, tmpDir, "fix: resolve bug")
	createCommit(t, tmpDir, "docs: update readme")

	repo, err := Open(tmpDir)
	require.NoError(t, err)

	commits, err := repo.GetCommitsSinceTag("v1.0.0")
	require.NoError(t, err)

	assert.Len(t, commits, 3)

	// Commits should be in reverse chronological order
	assert.Contains(t, commits[0].Subject, "docs: update readme")
	assert.Contains(t, commits[1].Subject, "fix: resolve bug")
	assert.Contains(t, commits[2].Subject, "feat: add new feature")
}

// T027: Test for Commit struct with hash, message, author
func TestCommit_Fields(t *testing.T) {
	tmpDir := t.TempDir()
	initRealGitRepo(t, tmpDir)

	createTag(t, tmpDir, "v1.0.0", "Release 1.0.0")
	createCommit(t, tmpDir, "feat: add amazing feature")

	repo, err := Open(tmpDir)
	require.NoError(t, err)

	commits, err := repo.GetCommitsSinceTag("v1.0.0")
	require.NoError(t, err)
	require.Len(t, commits, 1)

	commit := commits[0]
	assert.NotEmpty(t, commit.Hash)
	assert.Len(t, commit.Hash, 40) // Full SHA-1 hash
	assert.NotEmpty(t, commit.ShortHash)
	assert.Len(t, commit.ShortHash, 7)
	assert.Equal(t, "feat: add amazing feature", commit.Subject)
	assert.Equal(t, "Test User", commit.Author)
	assert.Equal(t, "test@example.com", commit.AuthorEmail)
	assert.False(t, commit.Timestamp.IsZero())
}

// T029: Test for empty commit list when no commits since tag
func TestRepository_GetCommitsSinceTag_NoCommits(t *testing.T) {
	tmpDir := t.TempDir()
	initRealGitRepo(t, tmpDir)

	// Create a tag at HEAD - no commits after it
	createTag(t, tmpDir, "v1.0.0", "Release 1.0.0")

	repo, err := Open(tmpDir)
	require.NoError(t, err)

	commits, err := repo.GetCommitsSinceTag("v1.0.0")
	require.NoError(t, err)
	assert.Empty(t, commits)
}

func TestRepository_GetCommitsSinceTag_TagNotFound(t *testing.T) {
	tmpDir := t.TempDir()
	initRealGitRepo(t, tmpDir)

	repo, err := Open(tmpDir)
	require.NoError(t, err)

	commits, err := repo.GetCommitsSinceTag("v999.0.0")
	assert.Error(t, err)
	assert.Nil(t, commits)
	assert.Contains(t, err.Error(), "not found")
}

func TestRepository_GetCommitsSinceTag_MultilineMessage(t *testing.T) {
	tmpDir := t.TempDir()
	initRealGitRepo(t, tmpDir)

	createTag(t, tmpDir, "v1.0.0", "Release 1.0.0")

	// Create commit with multiline message using git directly
	runGit(
		t,
		tmpDir,
		"commit",
		"--allow-empty",
		"-m",
		"feat: add feature\n\nThis is a detailed description\nof the feature.",
	)

	repo, err := Open(tmpDir)
	require.NoError(t, err)

	commits, err := repo.GetCommitsSinceTag("v1.0.0")
	require.NoError(t, err)
	require.Len(t, commits, 1)

	// Subject should be just the first line
	assert.Equal(t, "feat: add feature", commits[0].Subject)
	// Message should contain everything
	assert.Contains(t, commits[0].Message, "detailed description")
}
