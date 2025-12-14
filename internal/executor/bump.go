package executor

import (
	"context"
	"fmt"

	"github.com/benny123tw/bumpkin/internal/git"
	"github.com/benny123tw/bumpkin/internal/version"
)

// Request contains the parameters for a version bump operation
type Request struct {
	Repository    *git.Repository
	BumpType      version.BumpType
	CustomVersion string // Only used when BumpType is BumpCustom
	Prefix        string // Tag prefix (default: "v")
	Remote        string // Remote name (default: "origin")
	DryRun        bool   // If true, don't actually create/push tags
	NoPush        bool   // If true, create tag but don't push
	NoHooks       bool   // If true, skip hook execution
}

// Result contains the outcome of a version bump operation
type Result struct {
	PreviousVersion string
	NewVersion      string
	TagName         string
	CommitHash      string
	TagCreated      bool
	Pushed          bool
}

// Execute performs a version bump operation
func Execute(ctx context.Context, req Request) (*Result, error) {
	// Set defaults
	if req.Prefix == "" {
		req.Prefix = "v"
	}
	if req.Remote == "" {
		req.Remote = "origin"
	}

	// Get the latest tag
	latestTag, err := req.Repository.LatestTag(req.Prefix)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest tag: %w", err)
	}

	// Determine previous version
	var prevVersion version.Version
	if latestTag == nil || latestTag.Version == nil {
		prevVersion = version.Zero()
	} else {
		prevVersion = *latestTag.Version
	}

	// Calculate new version
	var newVersion version.Version
	switch req.BumpType {
	case version.BumpCustom:
		if req.CustomVersion == "" {
			return nil, fmt.Errorf("custom version not specified")
		}
		parsed, err := version.Parse(req.CustomVersion)
		if err != nil {
			return nil, fmt.Errorf("invalid custom version: %w", err)
		}
		newVersion = parsed
	case version.BumpPatch, version.BumpMinor, version.BumpMajor, version.BumpRelease,
		version.BumpPrereleaseAlpha, version.BumpPrereleaseBeta, version.BumpPrereleaseRC:
		newVersion = version.Bump(prevVersion, req.BumpType)
	default:
		return nil, fmt.Errorf("unsupported bump type: %s", req.BumpType)
	}

	// Build tag name
	tagName := newVersion.StringWithPrefix(req.Prefix)

	// Get HEAD commit hash
	headHash, err := req.Repository.GetHEAD()
	if err != nil {
		return nil, fmt.Errorf("failed to get HEAD: %w", err)
	}

	result := &Result{
		PreviousVersion: prevVersion.String(),
		NewVersion:      newVersion.String(),
		TagName:         tagName,
		CommitHash:      headHash.String(),
		TagCreated:      false,
		Pushed:          false,
	}

	// Dry run - don't actually do anything
	if req.DryRun {
		return result, nil
	}

	// Create the tag
	tagMessage := fmt.Sprintf("Release %s", tagName)
	if err := req.Repository.CreateTag(tagName, tagMessage); err != nil {
		return nil, fmt.Errorf("failed to create tag: %w", err)
	}
	result.TagCreated = true

	// Push if requested
	if !req.NoPush {
		// Check if remote exists before pushing
		hasRemote, err := req.Repository.HasRemote(req.Remote)
		if err != nil {
			return result, fmt.Errorf("failed to check remote: %w", err)
		}

		if hasRemote {
			if err := req.Repository.PushTag(tagName, req.Remote); err != nil {
				return result, fmt.Errorf("failed to push tag: %w", err)
			}
			result.Pushed = true
		}
	}

	return result, nil
}
