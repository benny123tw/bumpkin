package hooks

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"time"
)

// RunHook executes a single hook and returns the result
func RunHook(ctx context.Context, hook Hook, hookCtx *HookContext) *HookResult {
	return RunHookWithOutput(ctx, hook, hookCtx, os.Stdout, os.Stderr)
}

// RunHookWithOutput executes a single hook with custom output writers
func RunHookWithOutput(
	ctx context.Context,
	hook Hook,
	hookCtx *HookContext,
	stdout, stderr *os.File,
) *HookResult {
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

	// Stream output directly to stdout/stderr
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	// Run command
	err := cmd.Run()
	result.Duration = time.Since(start)

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

// RunHooksFailOpen executes multiple hooks in sequence, continuing on failure.
// Returns all results and a slice of warning messages for failed hooks.
// This is used for post-push hooks where failures should not block the workflow.
func RunHooksFailOpen(
	ctx context.Context,
	hooks []Hook,
	hookCtx *HookContext,
) ([]*HookResult, []string) {
	var results []*HookResult
	var warnings []string

	for _, hook := range hooks {
		result := RunHook(ctx, hook, hookCtx)
		results = append(results, result)

		if !result.Success {
			warnings = append(
				warnings,
				fmt.Sprintf("hook '%s' failed: %v", hook.Command, result.Error),
			)
		}
	}

	return results, warnings
}

// RunHookStreaming executes a hook and streams output lines to channels.
// Returns a channel for output lines and a channel for the final result.
// The line channel will be closed after the result is sent.
func RunHookStreaming(
	ctx context.Context,
	hook Hook,
	hookCtx *HookContext,
) (chan OutputLine, chan HookResult) {
	lineChan := make(chan OutputLine, 100) // Buffered to prevent blocking
	doneChan := make(chan HookResult, 1)

	go func() {
		defer close(lineChan)
		defer close(doneChan)

		start := time.Now()

		result := HookResult{
			Hook:    hook,
			Success: true,
		}

		// Skip empty commands
		if strings.TrimSpace(hook.Command) == "" {
			result.Duration = time.Since(start)
			doneChan <- result
			return
		}

		// Create command
		var cmd *exec.Cmd
		if runtime.GOOS == "windows" {
			//nolint:gosec // User-defined hook command from config file
			cmd = exec.CommandContext(ctx, "cmd", "/C", hook.Command)
		} else {
			//nolint:gosec // User-defined hook command from config file
			cmd = exec.CommandContext(ctx, "sh", "-c", hook.Command)
		}

		// Set environment variables
		if hookCtx != nil {
			cmd.Env = append(os.Environ(), hookCtx.ToEnv()...)
		} else {
			cmd.Env = os.Environ()
		}

		// Create pipes for stdout and stderr
		stdoutPipe, err := cmd.StdoutPipe()
		if err != nil {
			result.Success = false
			result.Error = fmt.Errorf("failed to create stdout pipe: %w", err)
			result.Duration = time.Since(start)
			doneChan <- result
			return
		}

		stderrPipe, err := cmd.StderrPipe()
		if err != nil {
			result.Success = false
			result.Error = fmt.Errorf("failed to create stderr pipe: %w", err)
			result.Duration = time.Since(start)
			doneChan <- result
			return
		}

		// Start the command
		if err := cmd.Start(); err != nil {
			result.Success = false
			result.Error = fmt.Errorf("failed to start hook: %w", err)
			result.Duration = time.Since(start)
			doneChan <- result
			return
		}

		// Read stdout and stderr in separate goroutines
		var wg sync.WaitGroup
		wg.Add(2)

		// Read stdout
		go func() {
			defer wg.Done()
			readPipeToChannel(stdoutPipe, Stdout, lineChan)
		}()

		// Read stderr
		go func() {
			defer wg.Done()
			readPipeToChannel(stderrPipe, Stderr, lineChan)
		}()

		// Wait for all readers to complete
		wg.Wait()

		// Wait for command to finish
		err = cmd.Wait()
		result.Duration = time.Since(start)

		if err != nil {
			result.Success = false
			result.Error = fmt.Errorf("hook failed: %w", err)
		}

		doneChan <- result
	}()

	return lineChan, doneChan
}

// readPipeToChannel reads lines from a pipe and sends them to the channel
func readPipeToChannel(pipe io.Reader, streamType StreamType, lineChan chan<- OutputLine) {
	scanner := bufio.NewScanner(pipe)
	for scanner.Scan() {
		lineChan <- OutputLine{
			Text:      scanner.Text(),
			Stream:    streamType,
			Timestamp: time.Now(),
		}
	}
	// Check for scanner errors (e.g., line too long, I/O error)
	// Send as stderr line so user sees the issue in output
	if err := scanner.Err(); err != nil {
		lineChan <- OutputLine{
			Text:      fmt.Sprintf("[scanner error: %v]", err),
			Stream:    Stderr,
			Timestamp: time.Now(),
		}
	}
}
