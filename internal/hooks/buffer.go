package hooks

import (
	"fmt"
	"strings"
	"sync"
)

// OutputBuffer is a thread-safe buffer holding output lines from hook execution
type OutputBuffer struct {
	lines    []OutputLine
	maxLines int
	mu       sync.RWMutex
}

// NewOutputBuffer creates a new OutputBuffer with the specified maximum line capacity
func NewOutputBuffer(maxLines int) *OutputBuffer {
	return &OutputBuffer{
		lines:    make([]OutputLine, 0),
		maxLines: maxLines,
	}
}

// AddLine appends a line to the buffer, evicting oldest lines if capacity exceeded
func (b *OutputBuffer) AddLine(line OutputLine) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.lines = append(b.lines, line)

	// Evict oldest lines if over capacity
	if b.maxLines > 0 && len(b.lines) > b.maxLines {
		excess := len(b.lines) - b.maxLines
		b.lines = b.lines[excess:]
	} else if b.maxLines == 0 {
		// With 0 max lines, keep nothing
		b.lines = b.lines[:0]
	}
}

// LineCount returns the current number of lines in the buffer
func (b *OutputBuffer) LineCount() int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return len(b.lines)
}

// Lines returns a copy of all lines in the buffer
func (b *OutputBuffer) Lines() []OutputLine {
	b.mu.RLock()
	defer b.mu.RUnlock()

	result := make([]OutputLine, len(b.lines))
	copy(result, b.lines)
	return result
}

// Render returns the formatted output string for display
func (b *OutputBuffer) Render() string {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if len(b.lines) == 0 {
		return ""
	}

	var sb strings.Builder
	for i, line := range b.lines {
		if i > 0 {
			sb.WriteString("\n")
		}
		// Format: [stream] text
		prefix := fmt.Sprintf("[%s] ", line.Stream.String())
		sb.WriteString(prefix)
		sb.WriteString(line.Text)
	}
	return sb.String()
}

// Clear removes all lines from the buffer
func (b *OutputBuffer) Clear() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.lines = b.lines[:0]
}
