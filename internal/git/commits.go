package git

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// Commit represents a git commit
type Commit struct {
	Hash        string
	ShortHash   string
	Message     string
	Subject     string
	Author      string
	AuthorEmail string
	Timestamp   time.Time
}

// GetCommitsSinceTag returns all commits between the given tag and HEAD
// Commits are returned in reverse chronological order (newest first)
func (r *Repository) GetCommitsSinceTag(tagName string) ([]*Commit, error) {
	// Find the tag reference
	tagRef, err := r.repo.Tag(tagName)
	if err != nil {
		return nil, fmt.Errorf("tag %q not found: %w", tagName, err)
	}

	// Resolve tag to commit hash
	tagCommitHash := r.resolveTagToCommit(tagRef)

	// Get HEAD
	head, err := r.repo.Head()
	if err != nil {
		return nil, fmt.Errorf("failed to get HEAD: %w", err)
	}

	// If HEAD is the same as tag, no commits since tag
	if head.Hash() == tagCommitHash {
		return []*Commit{}, nil
	}

	// Get all commits from HEAD
	commitIter, err := r.repo.Log(&git.LogOptions{
		From:  head.Hash(),
		Order: git.LogOrderCommitterTime,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get commit log: %w", err)
	}

	var commits []*Commit
	err = commitIter.ForEach(func(c *object.Commit) error {
		// Stop when we reach the tagged commit
		if c.Hash == tagCommitHash {
			return errStopIteration
		}

		commits = append(commits, commitFromObject(c))
		return nil
	})

	// errStopIteration is expected when we find the tag
	if err != nil && err != errStopIteration {
		return nil, fmt.Errorf("failed to iterate commits: %w", err)
	}

	return commits, nil
}

// GetAllCommits returns all commits from HEAD
func (r *Repository) GetAllCommits() ([]*Commit, error) {
	head, err := r.repo.Head()
	if err != nil {
		return nil, fmt.Errorf("failed to get HEAD: %w", err)
	}

	commitIter, err := r.repo.Log(&git.LogOptions{
		From:  head.Hash(),
		Order: git.LogOrderCommitterTime,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get commit log: %w", err)
	}

	var commits []*Commit
	err = commitIter.ForEach(func(c *object.Commit) error {
		commits = append(commits, commitFromObject(c))
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to iterate commits: %w", err)
	}

	return commits, nil
}

// resolveTagToCommit resolves a tag reference to its underlying commit hash
func (r *Repository) resolveTagToCommit(tagRef *plumbing.Reference) plumbing.Hash {
	// Try to get annotated tag object
	tagObj, err := r.repo.TagObject(tagRef.Hash())
	if err == nil {
		// Annotated tag - return the target commit
		return tagObj.Target
	}

	// Lightweight tag - ref points directly to commit
	return tagRef.Hash()
}

// commitFromObject converts a go-git Commit object to our Commit struct
func commitFromObject(c *object.Commit) *Commit {
	message := strings.TrimSpace(c.Message)
	subject, _, _ := strings.Cut(message, "\n")

	return &Commit{
		Hash:        c.Hash.String(),
		ShortHash:   c.Hash.String()[:7],
		Message:     message,
		Subject:     subject,
		Author:      c.Author.Name,
		AuthorEmail: c.Author.Email,
		Timestamp:   c.Author.When,
	}
}

// errStopIteration is used to stop commit iteration
var errStopIteration = fmt.Errorf("stop iteration")
