package executor

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPartialSuccessError_Error(t *testing.T) {
	tests := []struct {
		name     string
		phase    HookPhase
		err      error
		expected string
	}{
		{
			name:     "post-tag phase",
			phase:    PhasePostTag,
			err:      errors.New("hook command failed"),
			expected: "post-tag hook failed (tag already created): hook command failed",
		},
		{
			name:     "with wrapped error",
			phase:    PhasePostTag,
			err:      fmt.Errorf("script error: %w", errors.New("exit code 1")),
			expected: "post-tag hook failed (tag already created): script error: exit code 1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := &PartialSuccessError{
				Phase: tt.phase,
				Err:   tt.err,
				Result: &Result{
					TagName:    "v1.0.0",
					TagCreated: true,
				},
			}
			assert.Equal(t, tt.expected, err.Error())
		})
	}
}

func TestPartialSuccessError_Unwrap(t *testing.T) {
	originalErr := errors.New("original error")
	partialErr := &PartialSuccessError{
		Phase:  PhasePostTag,
		Err:    originalErr,
		Result: &Result{TagCreated: true},
	}

	// Test Unwrap
	assert.Equal(t, originalErr, partialErr.Unwrap())

	// Test errors.Is works with wrapped error
	wrappedErr := fmt.Errorf("wrapped: %w", originalErr)
	partialErr2 := &PartialSuccessError{
		Phase:  PhasePostTag,
		Err:    wrappedErr,
		Result: &Result{TagCreated: true},
	}
	assert.True(t, errors.Is(partialErr2, originalErr))
}

func TestPartialSuccessError_ErrorsAs(t *testing.T) {
	originalErr := errors.New("hook failed")
	partialErr := &PartialSuccessError{
		Phase: PhasePostTag,
		Err:   originalErr,
		Result: &Result{
			TagName:    "v1.0.0",
			TagCreated: true,
			Pushed:     false,
		},
	}

	// Test errors.As works
	var target *PartialSuccessError
	assert.True(t, errors.As(partialErr, &target))
	assert.Equal(t, PhasePostTag, target.Phase)
	assert.Equal(t, "v1.0.0", target.Result.TagName)
	assert.True(t, target.Result.TagCreated)

	// Test errors.As with wrapped PartialSuccessError
	wrappedPartial := fmt.Errorf("operation failed: %w", partialErr)
	var target2 *PartialSuccessError
	assert.True(t, errors.As(wrappedPartial, &target2))
	assert.Equal(t, PhasePostTag, target2.Phase)
}

func TestIsPartialSuccess(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "nil error",
			err:      nil,
			expected: false,
		},
		{
			name:     "regular error",
			err:      errors.New("some error"),
			expected: false,
		},
		{
			name: "partial success error",
			err: &PartialSuccessError{
				Phase:  PhasePostTag,
				Err:    errors.New("hook failed"),
				Result: &Result{TagCreated: true},
			},
			expected: true,
		},
		{
			name: "wrapped partial success error",
			err: fmt.Errorf("wrapped: %w", &PartialSuccessError{
				Phase:  PhasePostTag,
				Err:    errors.New("hook failed"),
				Result: &Result{TagCreated: true},
			}),
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, IsPartialSuccess(tt.err))
		})
	}
}

func TestGetPartialResult(t *testing.T) {
	result := &Result{
		TagName:         "v1.0.0",
		NewVersion:      "1.0.0",
		PreviousVersion: "0.9.0",
		TagCreated:      true,
		Pushed:          false,
	}

	tests := []struct {
		name           string
		err            error
		expectedResult *Result
	}{
		{
			name:           "nil error",
			err:            nil,
			expectedResult: nil,
		},
		{
			name:           "regular error",
			err:            errors.New("some error"),
			expectedResult: nil,
		},
		{
			name: "partial success error",
			err: &PartialSuccessError{
				Phase:  PhasePostTag,
				Err:    errors.New("hook failed"),
				Result: result,
			},
			expectedResult: result,
		},
		{
			name: "wrapped partial success error",
			err: fmt.Errorf("wrapped: %w", &PartialSuccessError{
				Phase:  PhasePostTag,
				Err:    errors.New("hook failed"),
				Result: result,
			}),
			expectedResult: result,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetPartialResult(tt.err)
			assert.Equal(t, tt.expectedResult, got)
		})
	}
}

func TestHookPhase_Constants(t *testing.T) {
	// Verify the phase constant value
	assert.Equal(t, HookPhase("post-tag"), PhasePostTag)
}
