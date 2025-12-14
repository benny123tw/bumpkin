package git

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// T013: Test for detecting git repository
func TestRepository_Open(t *testing.T) {
	// Create a temporary directory with a git repo
	tmpDir := t.TempDir()

	// Initialize a git repo
	initGitRepo(t, tmpDir)

	repo, err := Open(tmpDir)
	require.NoError(t, err)
	assert.NotNil(t, repo)
	assert.Equal(t, tmpDir, repo.Path)
}

// T015: Test for repository not found error
func TestRepository_Open_NotGitRepo(t *testing.T) {
	// Create a temporary directory without git
	tmpDir := t.TempDir()

	repo, err := Open(tmpDir)
	assert.Error(t, err)
	assert.Nil(t, repo)
	assert.Contains(t, err.Error(), "not a git repository")
}

func TestRepository_Open_NonExistentPath(t *testing.T) {
	repo, err := Open("/nonexistent/path/that/does/not/exist")
	assert.Error(t, err)
	assert.Nil(t, repo)
}

func TestRepository_Open_CurrentDir(t *testing.T) {
	// Get the current working directory
	cwd, err := os.Getwd()
	require.NoError(t, err)

	// Find git root by traversing up
	gitRoot := findGitRoot(cwd)
	if gitRoot == "" {
		t.Skip("Not running inside a git repository")
	}

	repo, err := Open(gitRoot)
	require.NoError(t, err)
	assert.NotNil(t, repo)
}

// Helper to initialize a git repo for testing
func initGitRepo(t *testing.T, dir string) {
	t.Helper()

	// Create .git directory structure manually for testing
	gitDir := filepath.Join(dir, ".git")
	require.NoError(t, os.MkdirAll(filepath.Join(gitDir, "objects"), 0o755))
	require.NoError(t, os.MkdirAll(filepath.Join(gitDir, "refs", "heads"), 0o755))

	// Create HEAD file
	headContent := []byte("ref: refs/heads/main\n")
	//nolint:gosec // test file permissions are fine
	require.NoError(t, os.WriteFile(filepath.Join(gitDir, "HEAD"), headContent, 0o644))

	// Create config file
	configContent := "[core]\n\trepositoryformatversion = 0\n\tfilemode = true\n\tbare = false\n"
	//nolint:gosec // test file permissions are fine
	require.NoError(t, os.WriteFile(filepath.Join(gitDir, "config"), []byte(configContent), 0o644))
}

// Helper to find git root directory
func findGitRoot(dir string) string {
	for {
		if _, err := os.Stat(filepath.Join(dir, ".git")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return ""
		}
		dir = parent
	}
}
