package git

import (
	"fmt"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
)

// PushTag pushes a specific tag to the remote repository
func (r *Repository) PushTag(tagName, remoteName string) error {
	// Check if remote exists
	hasRemote, err := r.HasRemote(remoteName)
	if err != nil {
		return err
	}
	if !hasRemote {
		return fmt.Errorf("remote %q not found", remoteName)
	}

	// Check if tag exists locally
	tagRef, err := r.repo.Tag(tagName)
	if err != nil {
		return fmt.Errorf("tag %q not found: %w", tagName, err)
	}

	// Create refspec for the tag
	refSpec := config.RefSpec(fmt.Sprintf(
		"refs/tags/%s:refs/tags/%s",
		tagRef.Name().Short(),
		tagRef.Name().Short(),
	))

	// Push the tag
	err = r.repo.Push(&git.PushOptions{
		RemoteName: remoteName,
		RefSpecs:   []config.RefSpec{refSpec},
	})
	if err != nil {
		if err == git.NoErrAlreadyUpToDate {
			return nil // Tag already exists on remote
		}
		return fmt.Errorf("failed to push tag: %w", err)
	}

	return nil
}

// PushAllTags pushes all tags to the remote repository
func (r *Repository) PushAllTags(remoteName string) error {
	hasRemote, err := r.HasRemote(remoteName)
	if err != nil {
		return err
	}
	if !hasRemote {
		return fmt.Errorf("remote %q not found", remoteName)
	}

	err = r.repo.Push(&git.PushOptions{
		RemoteName: remoteName,
		RefSpecs:   []config.RefSpec{"refs/tags/*:refs/tags/*"},
	})
	if err != nil {
		if err == git.NoErrAlreadyUpToDate {
			return nil
		}
		return fmt.Errorf("failed to push tags: %w", err)
	}

	return nil
}

// HasRemote checks if a remote with the given name exists
func (r *Repository) HasRemote(name string) (bool, error) {
	remotes, err := r.repo.Remotes()
	if err != nil {
		return false, fmt.Errorf("failed to list remotes: %w", err)
	}

	for _, remote := range remotes {
		if remote.Config().Name == name {
			return true, nil
		}
	}

	return false, nil
}

// GetRemoteURL returns the URL of the specified remote
func (r *Repository) GetRemoteURL(name string) (string, error) {
	remote, err := r.repo.Remote(name)
	if err != nil {
		return "", fmt.Errorf("remote %q not found: %w", name, err)
	}

	urls := remote.Config().URLs
	if len(urls) == 0 {
		return "", fmt.Errorf("remote %q has no URLs configured", name)
	}

	return urls[0], nil
}

// GetCurrentBranch returns the name of the current branch
func (r *Repository) GetCurrentBranch() (string, error) {
	head, err := r.repo.Head()
	if err != nil {
		return "", fmt.Errorf("failed to get HEAD: %w", err)
	}

	if !head.Name().IsBranch() {
		return "", fmt.Errorf("HEAD is not on a branch (detached HEAD)")
	}

	return head.Name().Short(), nil
}

// GetHEAD returns the commit hash of HEAD
func (r *Repository) GetHEAD() (plumbing.Hash, error) {
	head, err := r.repo.Head()
	if err != nil {
		return plumbing.ZeroHash, fmt.Errorf("failed to get HEAD: %w", err)
	}
	return head.Hash(), nil
}
