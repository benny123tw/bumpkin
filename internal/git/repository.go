package git

import (
	"fmt"

	"github.com/go-git/go-git/v5"
)

// Repository wraps a git repository
type Repository struct {
	Path string
	repo *git.Repository
}

// Open opens a git repository at the given path
func Open(path string) (*Repository, error) {
	repo, err := git.PlainOpen(path)
	if err != nil {
		if err == git.ErrRepositoryNotExists {
			return nil, fmt.Errorf("not a git repository: %s", path)
		}
		return nil, fmt.Errorf("failed to open repository: %w", err)
	}

	return &Repository{
		Path: path,
		repo: repo,
	}, nil
}

// OpenFromCurrent opens a git repository from the current working directory
// It traverses up the directory tree to find the repository root
func OpenFromCurrent() (*Repository, error) {
	repo, err := git.PlainOpenWithOptions(".", &git.PlainOpenOptions{
		DetectDotGit: true,
	})
	if err != nil {
		if err == git.ErrRepositoryNotExists {
			return nil, fmt.Errorf("not a git repository (or any of the parent directories)")
		}
		return nil, fmt.Errorf("failed to open repository: %w", err)
	}

	// Get the worktree to find the root path
	wt, err := repo.Worktree()
	if err != nil {
		return nil, fmt.Errorf("failed to get worktree: %w", err)
	}

	return &Repository{
		Path: wt.Filesystem.Root(),
		repo: repo,
	}, nil
}

// Raw returns the underlying go-git repository
func (r *Repository) Raw() *git.Repository {
	return r.repo
}
