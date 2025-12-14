package git

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"

	"github.com/benny123tw/bumpkin/internal/version"
)

// Tag represents a git tag
type Tag struct {
	Name        string
	CommitHash  string
	Tagger      string
	TaggerEmail string
	Message     string
	Timestamp   time.Time
	IsAnnotated bool
	Version     *version.Version
}

// ListTags returns all tags in the repository
func (r *Repository) ListTags() ([]*Tag, error) {
	var tags []*Tag

	tagRefs, err := r.repo.Tags()
	if err != nil {
		return nil, fmt.Errorf("failed to list tags: %w", err)
	}

	err = tagRefs.ForEach(func(ref *plumbing.Reference) error {
		tag := &Tag{
			Name: ref.Name().Short(),
		}

		// Try to get annotated tag object
		tagObj, err := r.repo.TagObject(ref.Hash())
		if err == nil {
			// Annotated tag
			tag.IsAnnotated = true
			tag.Tagger = tagObj.Tagger.Name
			tag.TaggerEmail = tagObj.Tagger.Email
			tag.Message = strings.TrimSpace(tagObj.Message)
			tag.Timestamp = tagObj.Tagger.When
			tag.CommitHash = tagObj.Target.String()
		} else {
			// Lightweight tag - ref points directly to commit
			tag.IsAnnotated = false
			tag.CommitHash = ref.Hash().String()

			// Get commit timestamp
			commit, err := r.repo.CommitObject(ref.Hash())
			if err == nil {
				tag.Timestamp = commit.Committer.When
			}
		}

		// Try to parse as semver
		if v, err := version.Parse(tag.Name); err == nil {
			tag.Version = &v
		}

		tags = append(tags, tag)
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to iterate tags: %w", err)
	}

	return tags, nil
}

// LatestTag returns the most recent semver tag with the given prefix
// Returns nil if no matching tags found
func (r *Repository) LatestTag(prefix string) (*Tag, error) {
	tags, err := r.ListTags()
	if err != nil {
		return nil, err
	}

	var latest *Tag
	for _, tag := range tags {
		// Skip tags that don't match prefix
		if !strings.HasPrefix(tag.Name, prefix) {
			continue
		}

		// Skip non-semver tags
		if tag.Version == nil {
			continue
		}

		if latest == nil || latest.Version.LessThan(*tag.Version) {
			latest = tag
		}
	}

	return latest, nil
}

// CreateTag creates an annotated tag at HEAD
func (r *Repository) CreateTag(name, message string) error {
	// Check if tag already exists
	tags, err := r.ListTags()
	if err != nil {
		return err
	}

	for _, tag := range tags {
		if tag.Name == name {
			return fmt.Errorf("tag %q already exists", name)
		}
	}

	// Get HEAD reference
	head, err := r.repo.Head()
	if err != nil {
		return fmt.Errorf("failed to get HEAD: %w", err)
	}

	// Get the commit object
	commit, err := r.repo.CommitObject(head.Hash())
	if err != nil {
		return fmt.Errorf("failed to get commit: %w", err)
	}

	// Create annotated tag
	_, err = r.repo.CreateTag(name, commit.Hash, &git.CreateTagOptions{
		Tagger: &object.Signature{
			Name:  commit.Author.Name,
			Email: commit.Author.Email,
			When:  time.Now(),
		},
		Message: message,
	})
	if err != nil {
		return fmt.Errorf("failed to create tag: %w", err)
	}

	return nil
}
