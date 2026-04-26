package git

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/go-git/go-git/v5/plumbing"
)

// PushTag pushes a specific tag to the remote repository.
//
// This shells out to the system `git` binary instead of using go-git's
// in-process Push so that the user's git config is honored — insteadOf
// URL rewrites, http.sslVerify, http.proxy, credential helpers, and SSH
// keys. go-git does not read any of those on its own, which breaks pushes
// to enterprise hosts that rely on them.
func (r *Repository) PushTag(ctx context.Context, tagName, remoteName string) error {
	hasRemote, err := r.HasRemote(remoteName)
	if err != nil {
		return err
	}
	if !hasRemote {
		return fmt.Errorf("remote %q not found", remoteName)
	}

	if _, err := r.repo.Tag(tagName); err != nil {
		return fmt.Errorf("tag %q not found: %w", tagName, err)
	}

	refSpec := "refs/tags/" + tagName + ":refs/tags/" + tagName
	cmd := exec.CommandContext(ctx, "git", "push", remoteName, refSpec)
	cmd.Dir = r.Path
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf(
			"failed to push tag: %w: %s",
			err, strings.TrimSpace(string(out)),
		)
	}
	return nil
}

// PushAllTags pushes all tags to the remote repository.
// Shells out to `git push --tags` for the same reasons as PushTag.
func (r *Repository) PushAllTags(ctx context.Context, remoteName string) error {
	hasRemote, err := r.HasRemote(remoteName)
	if err != nil {
		return err
	}
	if !hasRemote {
		return fmt.Errorf("remote %q not found", remoteName)
	}

	cmd := exec.CommandContext(ctx, "git", "push", remoteName, "--tags")
	cmd.Dir = r.Path
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf(
			"failed to push tags: %w: %s",
			err, strings.TrimSpace(string(out)),
		)
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
