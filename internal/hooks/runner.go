package hooks

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

// RunHook executes a single hook and returns the result
func RunHook(ctx context.Context, hook Hook, hookCtx *HookContext) *HookResult {
	start := time.Now()

	result := &HookResult{
		Hook:    hook,
		Success: true,
	}

	// Skip empty commands
	if strings.TrimSpace(hook.Command) == "" {
		result.Duration = time.Since(start)
		return result
	}

	// Create command
	// Note: G204 is expected here - hooks are user-defined commands from config
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		//nolint:gosec // User-defined hook command from config file
		cmd = exec.CommandContext(ctx, "cmd", "/C", hook.Command)
	} else {
		//nolint:gosec // User-defined hook command from config file
		cmd = exec.CommandContext(ctx, "sh", "-c", hook.Command)
	}

	// Set environment variables
	cmd.Env = append(os.Environ(), hookCtx.ToEnv()...)

	// Capture output
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Run command
	err := cmd.Run()
	result.Duration = time.Since(start)
	result.Output = stdout.String() + stderr.String()

	if err != nil {
		result.Success = false
		result.Error = fmt.Errorf("hook failed: %w", err)
	}

	return result
}

// RunHooks executes multiple hooks in sequence, stopping on first failure
func RunHooks(ctx context.Context, hooks []Hook, hookCtx *HookContext) ([]*HookResult, error) {
	var results []*HookResult

	for _, hook := range hooks {
		result := RunHook(ctx, hook, hookCtx)
		results = append(results, result)

		if !result.Success {
			return results, fmt.Errorf("hook '%s' failed: %w", hook.Command, result.Error)
		}
	}

	return results, nil
}

// CreateHooks creates Hook objects from command strings
func CreateHooks(commands []string, hookType HookType) []Hook {
	hooks := make([]Hook, len(commands))
	for i, cmd := range commands {
		hooks[i] = Hook{
			Command: cmd,
			Type:    hookType,
		}
	}
	return hooks
}
