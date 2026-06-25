package git

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	billy "github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/osfs"
	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/cache"
	"github.com/go-git/go-git/v5/storage"
	"github.com/go-git/go-git/v5/storage/filesystem"
	"github.com/go-git/go-git/v5/storage/filesystem/dotgit"
)

// Repository wraps a git repository
type Repository struct {
	Path string
	repo *gogit.Repository
}

// Open opens a git repository at the given path
func Open(path string) (*Repository, error) {
	repo, root, err := openPlain(path, false)
	if err != nil {
		if errors.Is(err, gogit.ErrRepositoryNotExists) {
			return nil, fmt.Errorf("not a git repository: %s", path)
		}
		return nil, fmt.Errorf("failed to open repository: %w", err)
	}

	return &Repository{
		Path: root,
		repo: repo,
	}, nil
}

// OpenFromCurrent opens a git repository from the current working directory
// It traverses up the directory tree to find the repository root
func OpenFromCurrent() (*Repository, error) {
	repo, root, err := openPlain(".", true)
	if err != nil {
		if errors.Is(err, gogit.ErrRepositoryNotExists) {
			return nil, fmt.Errorf("not a git repository (or any of the parent directories)")
		}
		return nil, fmt.Errorf("failed to open repository: %w", err)
	}

	return &Repository{
		Path: root,
		repo: repo,
	}, nil
}

// Raw returns the underlying go-git repository
func (r *Repository) Raw() *gogit.Repository {
	return r.repo
}

func openPlain(path string, detectDotGit bool) (*gogit.Repository, string, error) {
	repo, err := gogit.PlainOpenWithOptions(path, &gogit.PlainOpenOptions{
		DetectDotGit:          detectDotGit,
		EnableDotGitCommonDir: true,
	})
	if err == nil {
		root, rootErr := repositoryRoot(repo)
		if rootErr != nil {
			return nil, "", rootErr
		}

		return repo, root, nil
	}

	if !isWorktreeConfigExtensionError(err) {
		return nil, "", err
	}

	return openPlainWithWorktreeConfigExtension(path, detectDotGit)
}

func repositoryRoot(repo *gogit.Repository) (string, error) {
	wt, err := repo.Worktree()
	if err != nil {
		return "", fmt.Errorf("failed to get worktree: %w", err)
	}

	return wt.Filesystem.Root(), nil
}

func isWorktreeConfigExtensionError(err error) bool {
	if !errors.Is(err, gogit.ErrUnsupportedExtensionRepositoryFormatVersion) &&
		!errors.Is(err, gogit.ErrUnknownExtension) {
		return false
	}

	return strings.Contains(strings.ToLower(err.Error()), "worktreeconfig")
}

func openPlainWithWorktreeConfigExtension(
	path string,
	detectDotGit bool,
) (*gogit.Repository, string, error) {
	gitDir, worktree, err := resolveGitDir(path, detectDotGit)
	if err != nil {
		return nil, "", err
	}

	repositoryFS, err := repositoryFilesystem(gitDir)
	if err != nil {
		return nil, "", err
	}

	storer := worktreeConfigStorer{
		Storer: filesystem.NewStorage(repositoryFS, cache.NewObjectLRUDefault()),
	}

	repo, err := gogit.Open(storer, worktree)
	if err != nil {
		return nil, "", err
	}

	return repo, worktree.Root(), nil
}

type worktreeConfigStorer struct {
	storage.Storer
}

func (s worktreeConfigStorer) Config() (*config.Config, error) {
	cfg, err := s.Storer.Config()
	if err != nil {
		return nil, err
	}

	if cfg.Raw != nil && cfg.Raw.HasSection("extensions") {
		section := cfg.Raw.Section("extensions")
		section.RemoveOption("worktreeConfig")
		section.RemoveOption("worktreeconfig")
	}

	return cfg, nil
}

func resolveGitDir(path string, detectDotGit bool) (billy.Filesystem, billy.Filesystem, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		return nil, nil, err
	}

	// If detecting and the entry point is a file, start from its directory.
	// Every parent reached by the walk below is already a directory.
	if info, statErr := os.Stat(path); statErr == nil && !info.IsDir() && detectDotGit {
		path = filepath.Dir(path)
	}

	for {
		worktree := osfs.New(path)
		gitPath := filepath.Join(path, ".git")
		info, statErr := os.Stat(gitPath)
		if statErr == nil {
			if info.IsDir() {
				return osfs.New(gitPath), worktree, nil
			}

			gitDir, readErr := readGitDirFile(gitPath, path)
			if readErr != nil {
				return nil, nil, readErr
			}

			return osfs.New(gitDir), worktree, nil
		}
		if !os.IsNotExist(statErr) {
			return nil, nil, statErr
		}
		if !detectDotGit {
			return nil, nil, gogit.ErrRepositoryNotExists
		}

		parent := filepath.Dir(path)
		if parent == path {
			return nil, nil, gogit.ErrRepositoryNotExists
		}
		path = parent
	}
}

func readGitDirFile(gitPath, worktreePath string) (string, error) {
	content, err := os.ReadFile(gitPath)
	if err != nil {
		return "", err
	}

	line := strings.TrimSpace(strings.SplitN(string(content), "\n", 2)[0])
	gitDir, ok := strings.CutPrefix(line, "gitdir: ")
	if !ok {
		return "", fmt.Errorf(".git file has no gitdir prefix")
	}

	gitDir = strings.TrimSpace(gitDir)
	if filepath.IsAbs(gitDir) {
		return gitDir, nil
	}

	return filepath.Join(worktreePath, gitDir), nil
}

func repositoryFilesystem(gitDir billy.Filesystem) (billy.Filesystem, error) {
	commonDir, err := commonDirectory(gitDir)
	if err != nil {
		return nil, err
	}
	if commonDir == nil {
		return gitDir, nil
	}

	return dotgit.NewRepositoryFilesystem(gitDir, commonDir), nil
}

func commonDirectory(gitDir billy.Filesystem) (billy.Filesystem, error) {
	file, err := gitDir.Open("commondir")
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	path := strings.TrimSpace(string(content))
	if path == "" {
		return nil, nil
	}
	if filepath.IsAbs(path) {
		return osfs.New(path), nil
	}

	return osfs.New(filepath.Join(gitDir.Root(), path)), nil
}
