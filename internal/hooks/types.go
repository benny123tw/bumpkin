package hooks

import "time"

// HookType represents the type of hook
type HookType string

const (
	// PreTag hooks run before the tag is created
	PreTag HookType = "pre-tag"
	// PostTag hooks run after the tag is created
	PostTag HookType = "post-tag"
)

// Hook represents a single hook command
type Hook struct {
	Command string
	Type    HookType
}

// HookResult contains the result of running a hook
type HookResult struct {
	Hook     Hook
	Success  bool
	Output   string
	Error    error
	Duration time.Duration
}

// HookContext contains information available to hooks via environment variables
type HookContext struct {
	Version         string
	PreviousVersion string
	TagName         string
	Prefix          string
	Remote          string
	CommitHash      string
	DryRun          bool
}

// ToEnv converts hook context to environment variables
func (c *HookContext) ToEnv() []string {
	return []string{
		"BUMPKIN_VERSION=" + c.Version,
		"BUMPKIN_PREVIOUS_VERSION=" + c.PreviousVersion,
		"BUMPKIN_TAG_NAME=" + c.TagName,
		"BUMPKIN_PREFIX=" + c.Prefix,
		"BUMPKIN_REMOTE=" + c.Remote,
		"BUMPKIN_COMMIT_HASH=" + c.CommitHash,
		"BUMPKIN_DRY_RUN=" + boolToString(c.DryRun),
		// Also provide VERSION for convenience
		"VERSION=" + c.Version,
		"TAG=" + c.TagName,
	}
}

func boolToString(b bool) string {
	if b {
		return "true"
	}
	return "false"
}
