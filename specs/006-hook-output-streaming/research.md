# Research: Hook Output Streaming

**Feature**: 006-hook-output-streaming
**Date**: 2025-12-20
**Status**: Complete

## Research Questions

### 1. How to stream subprocess output to BubbleTea?

**Decision**: Use the channel + tea.Cmd pattern with io.Pipe

**Rationale**: BubbleTea's architecture requires all external data to flow through `tea.Cmd` functions that return `tea.Msg`. The recommended pattern uses:
1. A goroutine that reads subprocess output and sends lines to a channel
2. A `tea.Cmd` that waits on the channel and returns a message
3. Re-queuing the wait command in `Update()` to create a continuous listener

**Alternatives considered**:
- `Program.Send()` from external goroutines: Rejected because it requires passing the Program reference, adds coupling, and the channel pattern is more idiomatic
- Direct goroutine writes: Rejected because BubbleTea's `Update()` must be the sole point of model mutation

**Implementation pattern**:
```go
type hookLineMsg struct {
    line   string
    stream StreamType // stdout or stderr
}

func waitForLine(lineChan <-chan hookLineMsg) tea.Cmd {
    return func() tea.Msg {
        line, ok := <-lineChan
        if !ok {
            return nil // Channel closed
        }
        return line
    }
}
```

### 2. How to capture stdout/stderr separately?

**Decision**: Use `io.Pipe()` to create separate pipes for stdout and stderr

**Rationale**: The current implementation uses `*os.File` for output (`RunHookWithOutput`). To capture and stream separately:
1. Create `io.Pipe()` for stdout and stderr
2. Set `cmd.Stdout` and `cmd.Stderr` to the write ends
3. Read from read ends in separate goroutines
4. Merge into a single channel with stream type tags

**Alternatives considered**:
- Combine stdout/stderr into one stream (`cmd.Stderr = cmd.Stdout`): Rejected because spec requires visual distinction between streams
- Use `CombinedOutput()`: Rejected because it blocks until completion, no streaming

**Implementation pattern**:
```go
stdoutR, stdoutW := io.Pipe()
stderrR, stderrW := io.Pipe()

cmd.Stdout = stdoutW
cmd.Stderr = stderrW

// Read stdout in goroutine
go func() {
    scanner := bufio.NewScanner(stdoutR)
    for scanner.Scan() {
        lineChan <- hookLineMsg{line: scanner.Text(), stream: Stdout}
    }
}()
```

### 3. How to handle line buffering for non-interactive mode?

**Decision**: Use `bufio.Scanner` for line-by-line reading; write complete lines atomically

**Rationale**: `bufio.Scanner` with `ScanLines` split function ensures complete lines are processed. For non-interactive output, write to `os.Stdout` with a single `fmt.Println` call (atomic for typical line lengths).

**Alternatives considered**:
- Custom line buffer with explicit flush: Rejected as unnecessary complexity; Go's bufio handles this
- Unbuffered writes: Rejected because partial lines could interleave

### 4. How to display scrollable output in TUI?

**Decision**: Use `charmbracelet/bubbles/viewport` component

**Rationale**: The project already uses viewport for the commits pane. Viewport provides:
- Automatic keyboard navigation (arrow keys, j/k, Page Up/Down)
- Mouse wheel scrolling
- `SetContent()` for dynamic updates
- `GotoBottom()` for auto-scroll during streaming

**Implementation pattern**:
```go
m.hookOutputPane = viewport.New(width, height)
m.hookOutputPane.SetContent(strings.Join(m.hookLines, "\n"))
m.hookOutputPane.GotoBottom() // Auto-scroll to latest
```

### 5. How to distinguish stdout vs stderr visually?

**Decision**: Prefix stderr lines with styled marker and use different colors

**Rationale**: Terminal-standard approach. Use lipgloss for styling:
- stdout: normal text color
- stderr: red/orange color with "stderr:" prefix

**Implementation pattern**:
```go
var StdoutStyle = lipgloss.NewStyle()
var StderrStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("9")) // Red
```

### 6. How to handle high-throughput output?

**Decision**: Buffer lines in memory with configurable limit; use rate-limited viewport updates

**Rationale**: Spec requires handling 1,000 lines/sec while maintaining 100ms input responsiveness. Strategy:
1. Append all lines to buffer immediately (for scrollback)
2. Limit viewport re-renders to max 10/sec (every 100ms)
3. Use circular buffer or line limit to prevent memory exhaustion

**Alternatives considered**:
- Drop lines when buffer full: Rejected; spec requires preserving all content
- Throttle at read side: Rejected; could block subprocess

**Implementation notes**:
- Default buffer limit: 10,000 lines (per spec SC-002)
- Render batching: Accumulate lines between renders

### 7. How to handle hook state transitions in TUI?

**Decision**: Add new TUI state `StateExecutingHooks` with hook output pane visible

**Rationale**: Current `StateExecuting` just shows a spinner. New state:
1. Displays hook output pane (scrollable viewport)
2. Shows current hook name as header
3. Allows scrolling while hook runs
4. Transitions to normal `StateDone` or `StateError` on completion

**State flow**:
```
StateConfirm -> StateExecutingHooks -> StateDone/StateError
                     ↑
              (hook output streaming)
```

## Architecture Summary

### Components to Create

1. **`internal/hooks/buffer.go`**: Line buffer with stream type tracking
   - Thread-safe append
   - Configurable max lines
   - Render method returning formatted string

2. **`internal/hooks/stream.go`**: Streaming hook runner
   - Returns channel of line messages
   - Handles stdout/stderr pipes
   - Signals completion via done channel

3. **`internal/tui/hookpane.go`**: Hook output viewport wrapper
   - Viewport with auto-scroll
   - Header showing current hook
   - Styled stdout/stderr lines

### Message Types to Add

```go
// Hook output line received
type HookLineMsg struct {
    Line   string
    Stream StreamType
    Hook   string // Which hook produced this
}

// Hook execution completed
type HookCompleteMsg struct {
    Hook    string
    Success bool
    Error   error
}

// All hooks in phase completed
type HookPhaseCompleteMsg struct {
    Phase   HookPhase // pre-tag, post-tag, post-push
    Results []*HookResult
}
```

### Non-Interactive Mode

For non-interactive mode (when TUI is not running), the current behavior of streaming to `os.Stdout`/`os.Stderr` continues to work. The new streaming infrastructure is only activated when the TUI is running.

Detection: Check if a `tea.Program` is active or use a configuration flag.

## References

- [BubbleTea Realtime Example](https://github.com/charmbracelet/bubbletea/blob/main/examples/realtime/main.go)
- [Bubbles Viewport Documentation](https://pkg.go.dev/github.com/charmbracelet/bubbles/viewport)
- [Go os/exec Patterns](https://www.dolthub.com/blog/2022-11-28-go-os-exec-patterns/)
- [Commands in Bubble Tea](https://charm.land/blog/commands-in-bubbletea/)
