package tui

import (
	"time"

	"github.com/benny123tw/bumpkin/internal/executor"
	"github.com/benny123tw/bumpkin/internal/git"
	"github.com/benny123tw/bumpkin/internal/hooks"
	"github.com/benny123tw/bumpkin/internal/version"
)

// State represents the current state of the TUI
type State int

const (
	StateLoading State = iota
	StateVersionSelect
	StateCustomInput
	StateConfirm
	StateExecuting
	StateExecutingHooks // New state for hook execution with streaming output
	StateDone
	StateError
)

// Messages for state transitions

// RepoLoadedMsg is sent when repository info is loaded
type RepoLoadedMsg struct {
	CurrentVersion *version.Version
	Commits        []*git.Commit
	HasRemote      bool
	RemoteName     string
}

// ErrorMsg is sent when an error occurs
type ErrorMsg struct {
	Err error
}

// VersionSelectedMsg is sent when a version bump type is selected
type VersionSelectedMsg struct {
	BumpType version.BumpType
}

// CustomVersionMsg is sent when a custom version is entered
type CustomVersionMsg struct {
	Version string
}

// ConfirmMsg is sent when user confirms or cancels
type ConfirmMsg struct {
	Confirmed bool
}

// ExecuteStartMsg is sent to start execution
type ExecuteStartMsg struct{}

// ExecuteResultMsg is sent when execution completes
type ExecuteResultMsg struct {
	Result *executor.Result
}

// QuitMsg is sent when user wants to quit
type QuitMsg struct{}

// HookLineMsg is sent when a hook produces output
type HookLineMsg struct {
	Line        hooks.OutputLine // The output line
	HookCommand string           // Command string of the producing hook
	Phase       hooks.HookType   // Current execution phase
}

// HookStartMsg is sent when a hook begins execution
type HookStartMsg struct {
	Hook  hooks.Hook     // The hook starting execution
	Phase hooks.HookType // Execution phase
	Index int            // Hook index within phase (0-based)
	Total int            // Total hooks in phase
}

// HookCompleteMsg is sent when a single hook finishes
type HookCompleteMsg struct {
	Hook     hooks.Hook    // The completed hook
	Success  bool          // Whether hook succeeded
	Error    error         // Error if failed
	Duration time.Duration // Execution time
}

// HookPhaseCompleteMsg is sent when all hooks in a phase finish
type HookPhaseCompleteMsg struct {
	Phase        hooks.HookType      // Completed phase
	Results      []*hooks.HookResult // Results from all hooks
	AllSucceeded bool                // True if all hooks succeeded
}

// TagCreatedMsg is sent when the git tag has been created
type TagCreatedMsg struct {
	TagName    string
	CommitHash string
}

// PushCompleteMsg is sent when the tag has been pushed to remote
type PushCompleteMsg struct{}
