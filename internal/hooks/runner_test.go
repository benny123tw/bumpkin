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
	assert.Contains(t, result.Output, "hello")
	assert.NoError(t, result.Error)
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

	var cmd string
	if runtime.GOOS == "windows" {
		cmd = "echo %BUMPKIN_VERSION%"
	} else {
		cmd = "echo $BUMPKIN_VERSION"
	}

	hook := Hook{
		Command: cmd,
		Type:    PreTag,
	}

	result := RunHook(ctx, hook, hookCtx)

	require.True(t, result.Success)
	assert.Contains(t, result.Output, "1.2.3")
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
	assert.Contains(t, env, "BUMPKIN_TAG_NAME=v1.2.3")
	assert.Contains(t, env, "BUMPKIN_PREFIX=v")
	assert.Contains(t, env, "BUMPKIN_REMOTE=origin")
	assert.Contains(t, env, "BUMPKIN_COMMIT_HASH=abc123def")
	assert.Contains(t, env, "BUMPKIN_DRY_RUN=true")
	assert.Contains(t, env, "VERSION=1.2.3")
	assert.Contains(t, env, "TAG=v1.2.3")
}
