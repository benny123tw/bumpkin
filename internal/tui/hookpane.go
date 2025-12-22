package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/benny123tw/bumpkin/internal/hooks"
)

// HookPane is a TUI component for displaying streaming hook output
type HookPane struct {
	viewport    viewport.Model
	buffer      *hooks.OutputBuffer
	width       int
	height      int
	currentHook *hooks.Hook
	hookIndex   int
	hookTotal   int
}

// NewHookPane creates a new HookPane with the specified dimensions
func NewHookPane(width, height int) *HookPane {
	vp := viewport.New(width, height-2) // Reserve space for header
	return &HookPane{
		viewport:  vp,
		buffer:    hooks.NewOutputBuffer(10000), // 10,000 lines per spec
		width:     width,
		height:    height,
		hookIndex: 0,
		hookTotal: 0,
	}
}

// AddLine adds an output line to the pane and updates the viewport
func (p *HookPane) AddLine(line hooks.OutputLine) {
	p.buffer.AddLine(line)
	p.updateViewport()
	p.viewport.GotoBottom() // Auto-scroll to latest
}

// SetCurrentHook sets the currently executing hook for the header display
func (p *HookPane) SetCurrentHook(hook hooks.Hook, index, total int) {
	p.currentHook = &hook
	p.hookIndex = index
	p.hookTotal = total
}

// SetSize updates the pane dimensions
func (p *HookPane) SetSize(width, height int) {
	p.width = width
	p.height = height
	p.viewport.Width = width
	p.viewport.Height = height - 2 // Reserve space for header
	p.updateViewport()
}

// Clear removes all content from the pane
func (p *HookPane) Clear() {
	p.buffer.Clear()
	p.currentHook = nil
	p.hookIndex = 0
	p.hookTotal = 0
	p.updateViewport()
}

// Update handles input messages for scrolling
func (p *HookPane) Update(msg tea.Msg) (*HookPane, tea.Cmd) {
	var cmd tea.Cmd
	p.viewport, cmd = p.viewport.Update(msg)
	return p, cmd
}

// View renders the hook output pane
func (p *HookPane) View() string {
	var sb strings.Builder

	// Render header
	header := p.renderHeader()
	sb.WriteString(header)
	sb.WriteString("\n")

	// Render viewport content
	sb.WriteString(p.viewport.View())

	// Apply border style
	return HookOutputPaneStyle.Width(p.width).Render(sb.String())
}

// renderHeader creates the hook header line
func (p *HookPane) renderHeader() string {
	if p.currentHook == nil {
		return HookHeaderStyle.Render("─── Hook Output ───")
	}

	// Format: ─── pre-tag [1/3]: ./build.sh ───
	indexStr := fmt.Sprintf("%d/%d", p.hookIndex+1, p.hookTotal)
	header := fmt.Sprintf("─── %s [%s]: %s ───",
		p.currentHook.Type,
		indexStr,
		p.currentHook.Command,
	)
	return HookHeaderStyle.Render(header)
}

// updateViewport refreshes the viewport content from the buffer
func (p *HookPane) updateViewport() {
	content := p.renderContent()
	p.viewport.SetContent(content)
}

// renderContent formats the buffer content for display
func (p *HookPane) renderContent() string {
	lines := p.buffer.Lines()
	if len(lines) == 0 {
		return MutedStyle.Render("  Waiting for output...")
	}

	var sb strings.Builder
	for i, line := range lines {
		if i > 0 {
			sb.WriteString("\n")
		}

		// Format with stream prefix and styling
		prefix := StreamPrefixStyle.Render(fmt.Sprintf("[%s] ", line.Stream.String()))

		var styledText string
		if line.Stream == hooks.Stderr {
			styledText = StderrStyle.Render(line.Text)
		} else {
			styledText = StdoutStyle.Render(line.Text)
		}

		sb.WriteString(prefix)
		sb.WriteString(styledText)
	}

	return sb.String()
}

// Buffer returns the underlying output buffer (for testing/inspection)
func (p *HookPane) Buffer() *hooks.OutputBuffer {
	return p.buffer
}
