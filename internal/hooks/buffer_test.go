package hooks

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// T005: Test OutputBuffer.AddLine() - adding lines and retrieving them
func TestOutputBuffer_AddLine(t *testing.T) {
	buf := NewOutputBuffer(100)

	line1 := OutputLine{
		Text:      "hello world",
		Stream:    Stdout,
		Timestamp: time.Now(),
	}
	line2 := OutputLine{
		Text:      "error occurred",
		Stream:    Stderr,
		Timestamp: time.Now(),
	}

	buf.AddLine(line1)
	buf.AddLine(line2)

	lines := buf.Lines()
	require.Len(t, lines, 2)
	assert.Equal(t, "hello world", lines[0].Text)
	assert.Equal(t, Stdout, lines[0].Stream)
	assert.Equal(t, "error occurred", lines[1].Text)
	assert.Equal(t, Stderr, lines[1].Stream)
}

// T006: Test OutputBuffer.LineCount() - accurate count tracking
func TestOutputBuffer_LineCount(t *testing.T) {
	buf := NewOutputBuffer(100)

	assert.Equal(t, 0, buf.LineCount())

	buf.AddLine(OutputLine{Text: "line 1", Stream: Stdout, Timestamp: time.Now()})
	assert.Equal(t, 1, buf.LineCount())

	buf.AddLine(OutputLine{Text: "line 2", Stream: Stdout, Timestamp: time.Now()})
	assert.Equal(t, 2, buf.LineCount())

	buf.AddLine(OutputLine{Text: "line 3", Stream: Stderr, Timestamp: time.Now()})
	assert.Equal(t, 3, buf.LineCount())
}

// T007: Test OutputBuffer.MaxLines eviction - oldest lines removed when limit exceeded
func TestOutputBuffer_MaxLinesEviction(t *testing.T) {
	buf := NewOutputBuffer(3) // Only keep 3 lines

	// Add 5 lines
	for i := 1; i <= 5; i++ {
		buf.AddLine(OutputLine{
			Text:      "line " + string(rune('0'+i)),
			Stream:    Stdout,
			Timestamp: time.Now(),
		})
	}

	// Should only have 3 lines (the last 3)
	assert.Equal(t, 3, buf.LineCount())

	lines := buf.Lines()
	require.Len(t, lines, 3)
	// Lines 3, 4, 5 should remain (lines 1 and 2 evicted)
	assert.Equal(t, "line 3", lines[0].Text)
	assert.Equal(t, "line 4", lines[1].Text)
	assert.Equal(t, "line 5", lines[2].Text)
}

// T008: Test OutputBuffer.Render() - formatted output string generation
func TestOutputBuffer_Render(t *testing.T) {
	buf := NewOutputBuffer(100)

	buf.AddLine(OutputLine{Text: "Building project...", Stream: Stdout, Timestamp: time.Now()})
	buf.AddLine(OutputLine{Text: "Warning: deprecated API", Stream: Stderr, Timestamp: time.Now()})
	buf.AddLine(OutputLine{Text: "Build complete", Stream: Stdout, Timestamp: time.Now()})

	rendered := buf.Render()

	// Should contain all lines with appropriate formatting
	assert.Contains(t, rendered, "Building project...")
	assert.Contains(t, rendered, "Warning: deprecated API")
	assert.Contains(t, rendered, "Build complete")
}

// T009: Test OutputBuffer thread safety - concurrent AddLine calls
func TestOutputBuffer_ThreadSafety(t *testing.T) {
	buf := NewOutputBuffer(1000)
	var wg sync.WaitGroup

	// Spawn 10 goroutines, each adding 100 lines
	for range 10 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range 100 {
				buf.AddLine(OutputLine{
					Text:      "line from goroutine",
					Stream:    Stdout,
					Timestamp: time.Now(),
				})
			}
		}()
	}

	wg.Wait()

	// Should have 1000 lines total (10 goroutines * 100 lines each)
	assert.Equal(t, 1000, buf.LineCount())
}

// Test empty buffer
func TestOutputBuffer_Empty(t *testing.T) {
	buf := NewOutputBuffer(100)

	assert.Equal(t, 0, buf.LineCount())
	assert.Empty(t, buf.Lines())
	assert.Equal(t, "", buf.Render())
}

// Test buffer with max lines of 0 (should still work)
func TestOutputBuffer_ZeroMaxLines(t *testing.T) {
	buf := NewOutputBuffer(0)

	buf.AddLine(OutputLine{Text: "test", Stream: Stdout, Timestamp: time.Now()})

	// With 0 max lines, all lines should be evicted
	assert.Equal(t, 0, buf.LineCount())
}
