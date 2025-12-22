# Contracts: Hook Output Streaming

**Feature**: 006-hook-output-streaming

## Overview

This feature does not expose any external APIs. It is an internal TUI enhancement that streams hook output within the bumpkin application.

## Internal Interfaces

The following Go interfaces define the contracts between components:

### StreamingHookRunner

```go
// StreamingHookRunner executes hooks with output streaming capability
type StreamingHookRunner interface {
    // RunHookStreaming executes a hook and streams output to channels
    // Returns immediately after starting the hook
    // Output lines are sent to lineChan
    // Completion is signaled via doneChan
    RunHookStreaming(
        ctx context.Context,
        hook Hook,
        hookCtx *HookContext,
        lineChan chan<- OutputLine,
        doneChan chan<- HookResult,
    )
}
```

### OutputBufferWriter

```go
// OutputBufferWriter accumulates hook output
type OutputBufferWriter interface {
    // AddLine appends a line to the buffer
    AddLine(line OutputLine)

    // StartHook marks the beginning of a new hook's output
    StartHook(hook Hook)

    // EndHook marks the completion of a hook
    EndHook(hook Hook, success bool, err error)

    // Render returns the formatted output for display
    Render() string

    // LineCount returns total lines in buffer
    LineCount() int
}
```

### HookOutputPane

```go
// HookOutputPane is a TUI component for displaying hook output
type HookOutputPane interface {
    // SetContent updates the viewport content
    SetContent(content string)

    // GotoBottom scrolls to the latest output
    GotoBottom()

    // View renders the pane
    View() string

    // Update handles input messages
    Update(msg tea.Msg) (HookOutputPane, tea.Cmd)
}
```

## Message Protocol

Messages flow from the hook runner to the TUI model:

```
┌─────────────┐     HookStartMsg      ┌───────────┐
│  Hook       │ ───────────────────▶  │   TUI     │
│  Runner     │     HookLineMsg       │   Model   │
│             │ ───────────────────▶  │           │
│             │     HookCompleteMsg   │           │
│             │ ───────────────────▶  │           │
└─────────────┘                       └───────────┘
```

All messages are immutable value types.
