package cli

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// T121: Test for exit codes
func TestExitCodes(t *testing.T) {
	tests := []struct {
		name     string
		code     int
		expected int
	}{
		{"success", ExitSuccess, 0},
		{"general error", ExitGeneralError, 1},
		{"invalid args", ExitInvalidArgs, 2},
		{"not git repo", ExitNotGitRepo, 3},
		{"no commits", ExitNoCommits, 4},
		{"user cancelled", ExitUserCancelled, 5},
		{"hook failed", ExitHookFailed, 6},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.code)
		})
	}
}

func TestExitError(t *testing.T) {
	t.Run("with message only", func(t *testing.T) {
		err := NewExitError(ExitInvalidArgs, "invalid flag", nil)
		assert.Equal(t, "invalid flag", err.Error())
		assert.Equal(t, ExitInvalidArgs, err.Code)
		assert.Nil(t, err.Unwrap())
	})

	t.Run("with wrapped error", func(t *testing.T) {
		cause := errors.New("underlying error")
		err := NewExitError(ExitGeneralError, "git failed", cause)
		assert.Equal(t, "git failed: underlying error", err.Error())
		assert.Equal(t, ExitGeneralError, err.Code)
		assert.Equal(t, cause, err.Unwrap())
	})

	t.Run("with wrapped error only", func(t *testing.T) {
		cause := errors.New("underlying error")
		err := NewExitError(ExitGeneralError, "", cause)
		assert.Equal(t, "underlying error", err.Error())
	})
}

func TestGetExitCode(t *testing.T) {
	t.Run("nil error returns success", func(t *testing.T) {
		assert.Equal(t, ExitSuccess, GetExitCode(nil))
	})

	t.Run("ExitError returns its code", func(t *testing.T) {
		err := NewExitError(ExitNotGitRepo, "not a git repo", nil)
		assert.Equal(t, ExitNotGitRepo, GetExitCode(err))
	})

	t.Run("regular error returns general error", func(t *testing.T) {
		err := errors.New("some error")
		assert.Equal(t, ExitGeneralError, GetExitCode(err))
	})
}
