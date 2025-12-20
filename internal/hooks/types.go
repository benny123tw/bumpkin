package hooks

import "time"

// StreamType identifies the source stream of hook output
type StreamType int

const (
	// Stdout represents standard output stream
	Stdout StreamType = iota
	// Stderr represents standard error stream
	Stderr
)

// String returns the string representation of StreamType
func (s StreamType) String() string {
	switch s {
	case Stdout:
		return "stdout"
	case Stderr:
		return "stderr"
	default:
		return "unknown"
	}
}

// OutputLine represents a single line of output from a hook execution
type OutputLine struct {
	Text      string     // The line content (without trailing newline)
	Stream    StreamType // Source stream (stdout or stderr)
	Timestamp time.Time  // When the line was received
}

// HookType represents the type of hook
type HookType string

const (
	// PreTag hooks run before the tag is created
	PreTag HookType = "pre-tag"
	// PostTag hooks run after the tag is created
	PostTag HookType = "post-tag"
	// PostPush hooks run after the tag is pushed (fail-open behavior)
	PostPush HookType = "post-push"
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
		"BUMPKIN_TAG=" + c.TagName,
		"BUMPKIN_PREFIX=" + c.Prefix,
		"BUMPKIN_REMOTE=" + c.Remote,
		"BUMPKIN_COMMIT=" + c.CommitHash,
		"BUMPKIN_DRY_RUN=" + boolToString(c.DryRun),
		// Short aliases for convenience
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
