# Quickstart: Hook Output Streaming

**Feature**: 006-hook-output-streaming
**Date**: 2025-12-20

## Prerequisites

- Go 1.24+
- Existing bumpkin development environment
- Familiarity with BubbleTea TUI framework

## Feature Overview

This feature adds real-time streaming of hook output to the bumpkin TUI. When hooks execute (pre-tag, post-tag, post-push), their stdout/stderr is displayed in a scrollable pane within the terminal interface.

## Implementation Order

Follow TDD (Red → Green → Refactor) for each component:

### 1. Output Buffer (internal/hooks/buffer.go)

Start with the core data structure:

```go
// Test first
func TestOutputBuffer_AddLine(t *testing.T) {
    buf := NewOutputBuffer(100)
    buf.AddLine(OutputLine{Text: "hello", Stream: Stdout})
    assert.Equal(t, 1, buf.LineCount())
}

// Then implement
type OutputBuffer struct {
    lines    []OutputLine
    maxLines int
    mu       sync.RWMutex
}

func NewOutputBuffer(maxLines int) *OutputBuffer {
    return &OutputBuffer{maxLines: maxLines}
}

func (b *OutputBuffer) AddLine(line OutputLine) {
    b.mu.Lock()
    defer b.mu.Unlock()
    b.lines = append(b.lines, line)
    // Trim if over limit
    if len(b.lines) > b.maxLines {
        b.lines = b.lines[len(b.lines)-b.maxLines:]
    }
}
```

### 2. Streaming Hook Runner (internal/hooks/stream.go)

Add streaming capability to hook execution:

```go
// Test the streaming function
func TestRunHookStreaming(t *testing.T) {
    lineChan := make(chan OutputLine, 10)
    doneChan := make(chan HookResult, 1)

    hook := Hook{Command: "echo hello", Type: PreTag}
    ctx := context.Background()

    go RunHookStreaming(ctx, hook, nil, lineChan, doneChan)

    // Should receive "hello" on stdout
    line := <-lineChan
    assert.Equal(t, "hello", line.Text)
    assert.Equal(t, Stdout, line.Stream)

    // Should complete successfully
    result := <-doneChan
    assert.True(t, result.Success)
}
```

### 3. TUI Messages (internal/tui/messages.go)

Add message types:

```go
// HookLineMsg is sent when hook produces output
type HookLineMsg struct {
    Line    hooks.OutputLine
    Command string
}

// HookCompleteMsg is sent when hook finishes
type HookCompleteMsg struct {
    Command string
    Success bool
    Error   error
}
```

### 4. Hook Output Pane (internal/tui/hookpane.go)

Create the TUI component:

```go
type HookPane struct {
    viewport viewport.Model
    buffer   *hooks.OutputBuffer
    width    int
    height   int
}

func NewHookPane(width, height int) HookPane {
    vp := viewport.New(width, height)
    return HookPane{
        viewport: vp,
        buffer:   hooks.NewOutputBuffer(10000),
    }
}

func (p *HookPane) AddLine(line hooks.OutputLine) {
    p.buffer.AddLine(line)
    p.viewport.SetContent(p.buffer.Render())
    p.viewport.GotoBottom()
}
```

### 5. Model Integration (internal/tui/model.go)

Extend the TUI model:

```go
// Add to Model struct
hookPane    HookPane
hookLineCh  chan hooks.OutputLine
hookDoneCh  chan hooks.HookResult

// Add new state
const StateExecutingHooks State = "executing_hooks"

// Handle in Update()
case HookLineMsg:
    m.hookPane.AddLine(msg.Line)
    return m, waitForHookLine(m.hookLineCh)

case HookCompleteMsg:
    if !msg.Success {
        m.err = msg.Error
        m.state = StateError
    }
    // Check if more hooks to run...
```

## Testing the Feature

### Manual Testing

1. Create a test hook in `.bumpkin.yaml`:
```yaml
hooks:
  pre-tag:
    - "for i in 1 2 3 4 5; do echo \"Step $i\"; sleep 1; done"
```

2. Run bumpkin in a git repository:
```bash
go run ./cmd/bumpkin
```

3. Select a version bump and confirm. You should see:
   - Hook output appearing line by line
   - Scrollable viewport for output
   - Hook header showing which hook is running

### Unit Tests

Run all tests:
```bash
go test ./internal/hooks/... ./internal/tui/... -v
```

### Performance Testing

Test high-throughput output:
```yaml
hooks:
  pre-tag:
    - "for i in $(seq 1 1000); do echo \"Line $i\"; done"
```

Verify:
- TUI remains responsive during output
- All 1000 lines are captured
- Scrollback works after completion

## Key Files

| File | Purpose |
|------|---------|
| `internal/hooks/buffer.go` | Thread-safe line buffer |
| `internal/hooks/stream.go` | Streaming hook runner |
| `internal/tui/hookpane.go` | Output viewport component |
| `internal/tui/messages.go` | Hook-related TUI messages |
| `internal/tui/model.go` | State and message handling |

## Common Issues

### Output Not Appearing

- Check that `lineChan` is being read in the Update loop
- Ensure `waitForHookLine` is re-queued after each message

### Viewport Not Scrolling

- Call `GotoBottom()` after `SetContent()`
- Ensure viewport height is set correctly on WindowSizeMsg

### Race Conditions

- Use `sync.RWMutex` for buffer access
- Never mutate model outside of Update()

## Next Steps

After implementing:
1. Run `golangci-lint run` to verify code quality
2. Run full test suite: `go test ./...`
3. Test with real hooks in various scenarios
4. Consider adding elapsed time indicator for long-running hooks
