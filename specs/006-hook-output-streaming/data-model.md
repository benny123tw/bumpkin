# Data Model: Hook Output Streaming

**Feature**: 006-hook-output-streaming
**Date**: 2025-12-20

## Entities

### StreamType

Enumeration identifying the source stream of hook output.

| Value | Description |
|-------|-------------|
| `Stdout` | Standard output stream |
| `Stderr` | Standard error stream |

### OutputLine

A single line of output from a hook execution.

| Field | Type | Description |
|-------|------|-------------|
| `Text` | string | The line content (without trailing newline) |
| `Stream` | StreamType | Source stream (stdout or stderr) |
| `Timestamp` | time.Time | When the line was received |

### HookOutput

Accumulated output from a single hook execution.

| Field | Type | Description |
|-------|------|-------------|
| `Hook` | Hook | Reference to the hook being executed |
| `Lines` | []OutputLine | All output lines in order received |
| `StartTime` | time.Time | When hook execution started |
| `EndTime` | *time.Time | When hook execution completed (nil if running) |
| `Success` | *bool | Execution result (nil if running) |
| `Error` | error | Error if execution failed |

**Validation rules**:
- `Lines` may be empty (hook produced no output)
- `EndTime` is nil while hook is running
- `Success` is nil while hook is running
- `Error` is nil if `Success` is true

### OutputBuffer

Thread-safe buffer holding output from multiple hooks.

| Field | Type | Description |
|-------|------|-------------|
| `HookOutputs` | []HookOutput | Output from all hooks in execution order |
| `MaxLines` | int | Maximum total lines to retain (default: 10,000) |
| `TotalLines` | int | Current total line count across all hooks |

**Validation rules**:
- When `TotalLines` exceeds `MaxLines`, oldest lines are discarded
- `HookOutputs` maintains insertion order (first hook first)

**State transitions**:
1. **Empty** → **Streaming**: First hook starts, first line received
2. **Streaming** → **Streaming**: Additional lines received or new hook starts
3. **Streaming** → **Complete**: All hooks finished

### HookPhase

Enumeration for hook execution phases (existing, extended for this feature).

| Value | Description |
|-------|-------------|
| `PreTag` | Pre-tag hooks (fail-closed) |
| `PostTag` | Post-tag hooks (fail-closed) |
| `PostPush` | Post-push hooks (fail-open) |

## Relationships

```
OutputBuffer
    └── contains 0..* HookOutput (ordered by execution)
            └── contains 0..* OutputLine (ordered by receipt)
                    └── has 1 StreamType
            └── references 1 Hook
            └── has 1 HookPhase (derived from Hook.Type)
```

## Message Types (TUI Integration)

### HookLineMsg

Sent when a new output line is available.

| Field | Type | Description |
|-------|------|-------------|
| `Line` | OutputLine | The output line |
| `HookCommand` | string | Command string of the producing hook |
| `Phase` | HookPhase | Current execution phase |

### HookStartMsg

Sent when a hook begins execution.

| Field | Type | Description |
|-------|------|-------------|
| `Hook` | Hook | The hook starting execution |
| `Phase` | HookPhase | Execution phase |
| `Index` | int | Hook index within phase (0-based) |
| `Total` | int | Total hooks in phase |

### HookCompleteMsg

Sent when a single hook finishes.

| Field | Type | Description |
|-------|------|-------------|
| `Hook` | Hook | The completed hook |
| `Success` | bool | Whether hook succeeded |
| `Error` | error | Error if failed |
| `Duration` | time.Duration | Execution time |

### HookPhaseCompleteMsg

Sent when all hooks in a phase finish.

| Field | Type | Description |
|-------|------|-------------|
| `Phase` | HookPhase | Completed phase |
| `Results` | []*HookResult | Results from all hooks |
| `AllSucceeded` | bool | True if all hooks succeeded |

## TUI Model Extensions

### Model (extended)

| New Field | Type | Description |
|-----------|------|-------------|
| `hookOutputBuffer` | *OutputBuffer | Buffer for hook output |
| `hookOutputPane` | viewport.Model | Scrollable output viewport |
| `currentHook` | *Hook | Currently executing hook |
| `hookOutputChan` | chan HookLineMsg | Channel for receiving output |
| `hookDoneChan` | chan HookCompleteMsg | Channel for hook completion |

### State (extended)

| New Value | Description |
|-----------|-------------|
| `StateExecutingHooks` | Hooks are running with output streaming |

## Display Format

### Hook Header Format
```
─── pre-tag [1/3]: ./build.sh ───
```

### Output Line Format
```
[stdout] Building project...
[stderr] Warning: deprecated API usage
```

### Separator Between Hooks
```
─── pre-tag [2/3]: ./lint.sh ───
```
