package tui

import (
	"testing"
	"time"

	"github.com/benny123tw/bumpkin/internal/hooks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// T019: Test HookPane.AddLine() content update
func TestHookPane_AddLine(t *testing.T) {
	pane := NewHookPane(80, 20)

	line := hooks.OutputLine{
		Text:      "Hello, World!",
		Stream:    hooks.Stdout,
		Timestamp: time.Now(),
	}

	pane.AddLine(line)

	// Verify line was added to buffer
	assert.Equal(t, 1, pane.buffer.LineCount())

	// View should contain the line
	view := pane.View()
	assert.Contains(t, view, "Hello, World!")
}

// T020: Test HookPane.View() rendering with hook header
func TestHookPane_View_WithHeader(t *testing.T) {
	pane := NewHookPane(80, 20)

	hook := hooks.Hook{
		Command: "./build.sh",
		Type:    hooks.PreTag,
	}
	pane.SetCurrentHook(hook, 0, 3)

	view := pane.View()

	// Should contain hook header information
	assert.Contains(t, view, "pre-tag")
	assert.Contains(t, view, "./build.sh")
	assert.Contains(t, view, "1/3") // Index 0 + 1 displayed as 1
}

// Test HookPane with multiple lines
func TestHookPane_MultipleLines(t *testing.T) {
	pane := NewHookPane(80, 20)

	lines := []hooks.OutputLine{
		{Text: "Line 1", Stream: hooks.Stdout, Timestamp: time.Now()},
		{Text: "Line 2", Stream: hooks.Stderr, Timestamp: time.Now()},
		{Text: "Line 3", Stream: hooks.Stdout, Timestamp: time.Now()},
	}

	for _, line := range lines {
		pane.AddLine(line)
	}

	assert.Equal(t, 3, pane.buffer.LineCount())

	view := pane.View()
	assert.Contains(t, view, "Line 1")
	assert.Contains(t, view, "Line 2")
	assert.Contains(t, view, "Line 3")
}

// Test HookPane empty state
func TestHookPane_Empty(t *testing.T) {
	pane := NewHookPane(80, 20)

	view := pane.View()
	// Empty pane should still render without error
	require.NotEmpty(t, view)
}

// Test HookPane resize
func TestHookPane_Resize(t *testing.T) {
	pane := NewHookPane(80, 20)

	pane.AddLine(hooks.OutputLine{Text: "Test line", Stream: hooks.Stdout, Timestamp: time.Now()})

	// Resize
	pane.SetSize(100, 30)

	// Should still work after resize
	view := pane.View()
	assert.Contains(t, view, "Test line")
}

// Test HookPane Clear
func TestHookPane_Clear(t *testing.T) {
	pane := NewHookPane(80, 20)

	pane.AddLine(hooks.OutputLine{Text: "Line 1", Stream: hooks.Stdout, Timestamp: time.Now()})
	pane.AddLine(hooks.OutputLine{Text: "Line 2", Stream: hooks.Stdout, Timestamp: time.Now()})

	assert.Equal(t, 2, pane.buffer.LineCount())

	pane.Clear()

	assert.Equal(t, 0, pane.buffer.LineCount())
}
