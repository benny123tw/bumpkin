package tui

import (
	"github.com/benny123tw/bumpkin/internal/executor"
	"github.com/benny123tw/bumpkin/internal/git"
	"github.com/benny123tw/bumpkin/internal/version"
)

// State represents the current state of the TUI
type State int

const (
	StateLoading State = iota
	StateCommitList
	StateVersionSelect
	StateCustomInput
	StateConfirm
	StateExecuting
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
