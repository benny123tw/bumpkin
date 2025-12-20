package executor

import (
	"errors"
	"fmt"
)

// HookPhase represents the phase of hook execution
type HookPhase string

const (
	// PhasePostTag indicates the post-tag hook phase
	PhasePostTag HookPhase = "post-tag"
)

// PartialSuccessError represents an error that occurred after partial completion.
// This error type is used when an operation has partially succeeded (e.g., tag created)
// but a subsequent step failed (e.g., post-tag hook).
//
// Callers can use errors.As to detect this error type and access the partial result:
//
//	var partialErr *executor.PartialSuccessError
//	if errors.As(err, &partialErr) {
//	    // Tag was created, but hook failed
//	    fmt.Printf("Tag %s created, but %s hook failed: %v\n",
//	        partialErr.Result.TagName, partialErr.Phase, partialErr.Unwrap())
//	}
type PartialSuccessError struct {
	// Phase indicates which phase failed (e.g., "post-tag")
	Phase HookPhase

	// Err is the underlying error that caused the failure
	Err error

	// Result contains the partial result up to the point of failure.
	// Result is always non-nil when this error is returned.
	// For post-tag failures, Result.TagCreated will be true.
	Result *Result
}

// Error implements the error interface
func (e *PartialSuccessError) Error() string {
	return fmt.Sprintf("%s hook failed (tag already created): %v", e.Phase, e.Err)
}

// Unwrap returns the underlying error for use with errors.Is and errors.As
func (e *PartialSuccessError) Unwrap() error {
	return e.Err
}

// IsPartialSuccess checks if the given error represents a partial success.
// Returns true if the error is a PartialSuccessError, indicating that some
// operations completed successfully before the failure occurred.
func IsPartialSuccess(err error) bool {
	var partialErr *PartialSuccessError
	return errors.As(err, &partialErr)
}

// GetPartialResult extracts the partial result from a PartialSuccessError.
// Returns nil if the error is not a PartialSuccessError.
func GetPartialResult(err error) *Result {
	var partialErr *PartialSuccessError
	if errors.As(err, &partialErr) {
		return partialErr.Result
	}
	return nil
}
