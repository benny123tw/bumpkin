//nolint:goconst // "windows" string used for OS-specific test logic
package hooks

import (
	"context"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// T092: Test for running single hook command
func TestRunHook(t *testing.T) {
	ctx := context.Background()
	hookCtx := &HookContext{
		Version: "1.0.0",
		TagName: "v1.0.0",
	}

	hook := Hook{
		Command: "echo hello",
		Type:    PreTag,
	}

	result := RunHook(ctx, hook, hookCtx)

	require.True(t, result.Success)
	assert.NoError(t, result.Error)
	// Output is now streamed to stdout, not captured
}

// T094: Test for hook environment variables
func TestRunHook_EnvironmentVariables(t *testing.T) {
	ctx := context.Background()
	hookCtx := &HookContext{
		Version:         "1.2.3",
		PreviousVersion: "1.2.2",
		TagName:         "v1.2.3",
		Prefix:          "v",
		Remote:          "origin",
		CommitHash:      "abc123",
		DryRun:          false,
	}

	// Test that env vars are set by checking exit code of a test command
	var cmd string
	if runtime.GOOS == "windows" {
		cmd = "if \"%BUMPKIN_VERSION%\"==\"1.2.3\" exit 0"
	} else {
		cmd = "test \"$BUMPKIN_VERSION\" = \"1.2.3\""
	}

	hook := Hook{
		Command: cmd,
		Type:    PreTag,
	}

	result := RunHook(ctx, hook, hookCtx)

	require.True(t, result.Success)
	// Output is now streamed to stdout, env vars verified via exit code
}

// T096: Test for hook failure (non-zero exit)
func TestRunHook_Failure(t *testing.T) {
	ctx := context.Background()
	hookCtx := &HookContext{
		Version: "1.0.0",
	}

	hook := Hook{
		Command: "exit 1",
		Type:    PreTag,
	}

	result := RunHook(ctx, hook, hookCtx)

	assert.False(t, result.Success)
	assert.Error(t, result.Error)
}

// T098: Test for running multiple hooks in sequence
func TestRunHooks(t *testing.T) {
	ctx := context.Background()
	hookCtx := &HookContext{
		Version: "1.0.0",
	}

	hooks := []Hook{
		{Command: "echo first", Type: PreTag},
		{Command: "echo second", Type: PreTag},
		{Command: "echo third", Type: PreTag},
	}

	results, err := RunHooks(ctx, hooks, hookCtx)

	require.NoError(t, err)
	assert.Len(t, results, 3)
	for _, result := range results {
		assert.True(t, result.Success)
	}
}

func TestRunHooks_StopsOnFailure(t *testing.T) {
	ctx := context.Background()
	hookCtx := &HookContext{
		Version: "1.0.0",
	}

	hooks := []Hook{
		{Command: "echo first", Type: PreTag},
		{Command: "exit 1", Type: PreTag},
		{Command: "echo third", Type: PreTag},
	}

	results, err := RunHooks(ctx, hooks, hookCtx)

	require.Error(t, err)
	assert.Len(t, results, 2) // Only first two ran
	assert.True(t, results[0].Success)
	assert.False(t, results[1].Success)
}

func TestRunHook_Timeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	hookCtx := &HookContext{
		Version: "1.0.0",
	}

	hook := Hook{
		Command: "sleep 10",
		Type:    PreTag,
	}

	result := RunHook(ctx, hook, hookCtx)

	assert.False(t, result.Success)
	assert.Error(t, result.Error)
}

func TestRunHook_EmptyCommand(t *testing.T) {
	ctx := context.Background()
	hookCtx := &HookContext{
		Version: "1.0.0",
	}

	hook := Hook{
		Command: "",
		Type:    PreTag,
	}

	result := RunHook(ctx, hook, hookCtx)

	// Empty command should be skipped (success)
	assert.True(t, result.Success)
}

func TestHookContext_ToEnv(t *testing.T) {
	ctx := &HookContext{
		Version:         "1.2.3",
		PreviousVersion: "1.2.2",
		TagName:         "v1.2.3",
		Prefix:          "v",
		Remote:          "origin",
		CommitHash:      "abc123def",
		DryRun:          true,
	}

	env := ctx.ToEnv()

	assert.Contains(t, env, "BUMPKIN_VERSION=1.2.3")
	assert.Contains(t, env, "BUMPKIN_PREVIOUS_VERSION=1.2.2")
	assert.Contains(t, env, "BUMPKIN_TAG=v1.2.3")
	assert.Contains(t, env, "BUMPKIN_PREFIX=v")
	assert.Contains(t, env, "BUMPKIN_REMOTE=origin")
	assert.Contains(t, env, "BUMPKIN_COMMIT=abc123def")
	assert.Contains(t, env, "BUMPKIN_DRY_RUN=true")
	assert.Contains(t, env, "VERSION=1.2.3")
	assert.Contains(t, env, "TAG=v1.2.3")
}

// T004: Test for PostPush HookType constant
func TestPostPushHookType(t *testing.T) {
	assert.Equal(t, HookType("post-push"), PostPush)
	assert.NotEqual(t, PostPush, PreTag)
	assert.NotEqual(t, PostPush, PostTag)
}

// T005: Test for RunHooksFailOpen function (continues on failure)
func TestRunHooksFailOpen(t *testing.T) {
	ctx := context.Background()
	hookCtx := &HookContext{
		Version: "1.0.0",
		TagName: "v1.0.0",
	}

	hooks := []Hook{
		{Command: "echo first", Type: PostPush},
		{Command: "exit 1", Type: PostPush},     // This fails
		{Command: "echo third", Type: PostPush}, // Should still run
	}

	results, warnings := RunHooksFailOpen(ctx, hooks, hookCtx)

	// All hooks should have run
	assert.Len(t, results, 3)
	assert.True(t, results[0].Success)
	assert.False(t, results[1].Success)
	assert.True(t, results[2].Success)

	// Should have one warning for the failed hook
	assert.Len(t, warnings, 1)
	assert.Contains(t, warnings[0], "exit 1")
}

func TestRunHooksFailOpen_AllSuccess(t *testing.T) {
	ctx := context.Background()
	hookCtx := &HookContext{
		Version: "1.0.0",
	}

	hooks := []Hook{
		{Command: "echo first", Type: PostPush},
		{Command: "echo second", Type: PostPush},
	}

	results, warnings := RunHooksFailOpen(ctx, hooks, hookCtx)

	assert.Len(t, results, 2)
	assert.True(t, results[0].Success)
	assert.True(t, results[1].Success)
	assert.Empty(t, warnings)
}

func TestRunHooksFailOpen_AllFail(t *testing.T) {
	ctx := context.Background()
	hookCtx := &HookContext{
		Version: "1.0.0",
	}

	hooks := []Hook{
		{Command: "exit 1", Type: PostPush},
		{Command: "exit 2", Type: PostPush},
	}

	results, warnings := RunHooksFailOpen(ctx, hooks, hookCtx)

	assert.Len(t, results, 2)
	assert.False(t, results[0].Success)
	assert.False(t, results[1].Success)
	assert.Len(t, warnings, 2)
}

// T016: Test RunHookStreaming basic output capture
func TestRunHookStreaming_BasicOutput(t *testing.T) {
	ctx := context.Background()
	hookCtx := &HookContext{
		Version: "1.0.0",
		TagName: "v1.0.0",
	}

	hook := Hook{
		Command: "echo hello",
		Type:    PreTag,
	}

	lineChan, doneChan := RunHookStreaming(ctx, hook, hookCtx)

	// Collect all output lines
	var lines []OutputLine
	done := false
	for !done {
		select {
		case line, ok := <-lineChan:
			if ok {
				lines = append(lines, line)
			}
		case result := <-doneChan:
			assert.True(t, result.Success)
			assert.NoError(t, result.Error)
			done = true
		case <-time.After(5 * time.Second):
			t.Fatal("timeout waiting for hook to complete")
		}
	}

	// Drain any remaining buffered lines after done signal
	for line := range lineChan {
		lines = append(lines, line)
	}

	// Should have received "hello" on stdout
	require.GreaterOrEqual(t, len(lines), 1)
	assert.Equal(t, "hello", lines[0].Text)
	assert.Equal(t, Stdout, lines[0].Stream)
}

// T017: Test RunHookStreaming stdout/stderr separation
func TestRunHookStreaming_StdoutStderrSeparation(t *testing.T) {
	ctx := context.Background()
	hookCtx := &HookContext{
		Version: "1.0.0",
	}

	// Command that writes to both stdout and stderr (platform-specific)
	var cmd string
	if runtime.GOOS == "windows" {
		cmd = "echo stdout_line & echo stderr_line 1>&2"
	} else {
		cmd = "echo stdout_line && echo stderr_line >&2"
	}

	hook := Hook{
		Command: cmd,
		Type:    PreTag,
	}

	lineChan, doneChan := RunHookStreaming(ctx, hook, hookCtx)

	var stdoutLines, stderrLines []OutputLine
	done := false
	for !done {
		select {
		case line, ok := <-lineChan:
			if ok {
				if line.Stream == Stdout {
					stdoutLines = append(stdoutLines, line)
				} else {
					stderrLines = append(stderrLines, line)
				}
			}
		case result := <-doneChan:
			assert.True(t, result.Success)
			done = true
		case <-time.After(5 * time.Second):
			t.Fatal("timeout waiting for hook to complete")
		}
	}

	// Drain any remaining buffered lines after done signal
	for line := range lineChan {
		if line.Stream == Stdout {
			stdoutLines = append(stdoutLines, line)
		} else {
			stderrLines = append(stderrLines, line)
		}
	}

	// Should have at least one line on each stream
	require.GreaterOrEqual(t, len(stdoutLines), 1)
	require.GreaterOrEqual(t, len(stderrLines), 1)
	assert.Equal(t, "stdout_line", stdoutLines[0].Text)
	assert.Equal(t, "stderr_line", stderrLines[0].Text)
}

// T018: Test RunHookStreaming channel closure on completion
func TestRunHookStreaming_ChannelClosure(t *testing.T) {
	ctx := context.Background()
	hookCtx := &HookContext{
		Version: "1.0.0",
	}

	hook := Hook{
		Command: "echo done",
		Type:    PreTag,
	}

	lineChan, doneChan := RunHookStreaming(ctx, hook, hookCtx)

	// Wait for completion
	select {
	case <-doneChan:
		// After done, line channel should be closed
		// Drain remaining lines
		for range lineChan {
			// consume remaining lines
		}
		// Channel should now be closed
		_, ok := <-lineChan
		assert.False(t, ok, "line channel should be closed after hook completion")
	case <-time.After(5 * time.Second):
		t.Fatal("timeout waiting for hook to complete")
	}
}

// Test RunHookStreaming with failing hook
func TestRunHookStreaming_Failure(t *testing.T) {
	ctx := context.Background()
	hookCtx := &HookContext{
		Version: "1.0.0",
	}

	hook := Hook{
		Command: "echo before_fail && exit 1",
		Type:    PreTag,
	}

	lineChan, doneChan := RunHookStreaming(ctx, hook, hookCtx)

	var lines []OutputLine
	done := false
	for !done {
		select {
		case line, ok := <-lineChan:
			if ok {
				lines = append(lines, line)
			}
		case result := <-doneChan:
			assert.False(t, result.Success)
			assert.Error(t, result.Error)
			done = true
		case <-time.After(5 * time.Second):
			t.Fatal("timeout waiting for hook to complete")
		}
	}

	// Drain any remaining buffered lines after done signal
	for line := range lineChan {
		lines = append(lines, line)
	}

	// Should have received output before failure
	require.GreaterOrEqual(t, len(lines), 1)
	assert.Equal(t, "before_fail", lines[0].Text)
}
