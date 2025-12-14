package cli

// Exit codes per CLI contract
const (
	ExitSuccess       = 0 // Success
	ExitGeneralError  = 1 // General error (git operation failed, hook failed, etc.)
	ExitInvalidArgs   = 2 // Invalid arguments
	ExitNotGitRepo    = 3 // Not a git repository
	ExitNoCommits     = 4 // No commits since last tag
	ExitUserCancelled = 5 // User cancelled operation
	ExitHookFailed    = 6 // Hook execution failed
)

// ExitError is an error that carries an exit code
type ExitError struct {
	Code    int
	Message string
	Err     error
}

func (e *ExitError) Error() string {
	if e.Err != nil {
		if e.Message != "" {
			return e.Message + ": " + e.Err.Error()
		}
		return e.Err.Error()
	}
	return e.Message
}

func (e *ExitError) Unwrap() error {
	return e.Err
}

// NewExitError creates a new ExitError
func NewExitError(code int, message string, err error) *ExitError {
	return &ExitError{Code: code, Message: message, Err: err}
}

// GetExitCode returns the exit code for an error
// Returns ExitSuccess (0) for nil errors
func GetExitCode(err error) int {
	if err == nil {
		return ExitSuccess
	}

	if exitErr, ok := err.(*ExitError); ok {
		return exitErr.Code
	}

	return ExitGeneralError
}
