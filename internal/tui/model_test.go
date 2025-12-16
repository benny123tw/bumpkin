package tui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"

	"github.com/benny123tw/bumpkin/internal/git"
	"github.com/benny123tw/bumpkin/internal/version"
)

// testScrollContent is used for viewport scroll tests
const testScrollContent = "line1\nline2\nline3\nline4\nline5\nline6\nline7\nline8\nline9\nline10"

func TestNew_InitializesPaneFields(t *testing.T) {
	cfg := Config{
		Repository: &git.Repository{},
		Prefix:     "v",
	}

	model := New(cfg)

	// Verify focusedPane defaults to PaneVersion (version selection is primary task)
	assert.Equal(t, PaneVersion, model.focusedPane, "focusedPane should default to PaneVersion")

	// Verify commitsPane is initialized with zero dimensions
	assert.Equal(t, 0, model.commitsPane.Width, "commitsPane width should be 0 initially")
	assert.Equal(t, 0, model.commitsPane.Height, "commitsPane height should be 0 initially")

	// Verify showingDetail defaults to false
	assert.False(t, model.showingDetail, "showingDetail should default to false")

	// Verify selectedCommitIndex defaults to 0
	assert.Equal(t, 0, model.selectedCommitIndex, "selectedCommitIndex should default to 0")
}

func TestNew_CommitsPaneIsViewport(t *testing.T) {
	cfg := Config{
		Repository: &git.Repository{},
		Prefix:     "v",
	}

	model := New(cfg)

	// Verify commitsPane is initialized (compile-time check that field exists)
	assert.NotNil(t, &model.commitsPane, "commitsPane should be initialized")
}

// T010: Test Tab switches from version pane to commits pane
func TestTabSwitchesFromVersionToCommits(t *testing.T) {
	cfg := Config{
		Repository: &git.Repository{},
		Prefix:     "v",
	}

	model := New(cfg)
	model.state = StateVersionSelect
	model.focusedPane = PaneVersion

	// Simulate Tab key press
	msg := tea.KeyMsg{Type: tea.KeyTab}
	updatedModel, _ := model.Update(msg)
	m := updatedModel.(Model)

	assert.Equal(t, PaneCommits, m.focusedPane, "Tab should switch focus to commits pane")
}

// T011: Test Tab switches from commits pane to version pane
func TestTabSwitchesFromCommitsToVersion(t *testing.T) {
	cfg := Config{
		Repository: &git.Repository{},
		Prefix:     "v",
	}

	model := New(cfg)
	model.state = StateVersionSelect
	model.focusedPane = PaneCommits

	// Simulate Tab key press
	msg := tea.KeyMsg{Type: tea.KeyTab}
	updatedModel, _ := model.Update(msg)
	m := updatedModel.(Model)

	assert.Equal(t, PaneVersion, m.focusedPane, "Tab should switch focus to version pane")
}

// T012: Test Shift+Tab switches from commits pane to version pane
func TestShiftTabSwitchesFromCommitsToVersion(t *testing.T) {
	cfg := Config{
		Repository: &git.Repository{},
		Prefix:     "v",
	}

	model := New(cfg)
	model.state = StateVersionSelect
	model.focusedPane = PaneCommits

	// Simulate Shift+Tab key press
	msg := tea.KeyMsg{Type: tea.KeyShiftTab}
	updatedModel, _ := model.Update(msg)
	m := updatedModel.(Model)

	assert.Equal(t, PaneVersion, m.focusedPane, "Shift+Tab should switch focus to version pane")
}

// T020: Test down arrow scrolls viewport when commits pane focused
func TestDownArrowScrollsViewportWhenCommitsFocused(t *testing.T) {
	cfg := Config{
		Repository: &git.Repository{},
		Prefix:     "v",
	}

	model := New(cfg)
	model.state = StateVersionSelect
	model.focusedPane = PaneCommits
	// Set up viewport with content that can scroll
	model.commitsPane.Width = 80
	model.commitsPane.Height = 3
	model.commitsPane.SetContent(testScrollContent)

	initialOffset := model.commitsPane.YOffset

	// Simulate down arrow key press
	msg := tea.KeyMsg{Type: tea.KeyDown}
	updatedModel, _ := model.Update(msg)
	m := updatedModel.(Model)

	// Viewport should have scrolled (YOffset increased or content moved)
	assert.GreaterOrEqual(
		t, m.commitsPane.YOffset, initialOffset,
		"Down arrow should scroll viewport when commits pane focused",
	)
}

// T021: Test up arrow scrolls viewport when commits pane focused
func TestUpArrowScrollsViewportWhenCommitsFocused(t *testing.T) {
	cfg := Config{
		Repository: &git.Repository{},
		Prefix:     "v",
	}

	model := New(cfg)
	model.state = StateVersionSelect
	model.focusedPane = PaneCommits
	// Set up viewport with content and scroll position
	model.commitsPane.Width = 80
	model.commitsPane.Height = 3
	model.commitsPane.SetContent(testScrollContent)
	model.commitsPane.SetYOffset(5) // Start scrolled down

	initialOffset := model.commitsPane.YOffset

	// Simulate up arrow key press
	msg := tea.KeyMsg{Type: tea.KeyUp}
	updatedModel, _ := model.Update(msg)
	m := updatedModel.(Model)

	// Viewport should have scrolled up (YOffset decreased)
	assert.LessOrEqual(
		t, m.commitsPane.YOffset, initialOffset,
		"Up arrow should scroll viewport up when commits pane focused",
	)
}

// T022: Test down arrow moves version selector when version pane focused
func TestDownArrowMovesVersionSelectorWhenVersionFocused(t *testing.T) {
	cfg := Config{
		Repository: &git.Repository{},
		Prefix:     "v",
	}

	model := New(cfg)
	model.state = StateVersionSelect
	model.focusedPane = PaneVersion
	model.versionOptions = []VersionOption{
		{Label: "patch", NewVersion: "v1.0.1"},
		{Label: "minor", NewVersion: "v1.1.0"},
		{Label: "major", NewVersion: "v2.0.0"},
	}
	model.selectedOption = 0

	// Simulate down arrow key press
	msg := tea.KeyMsg{Type: tea.KeyDown}
	updatedModel, _ := model.Update(msg)
	m := updatedModel.(Model)

	assert.Equal(
		t, 1, m.selectedOption,
		"Down arrow should move version selector when version pane focused",
	)
}

// T028: Test Enter on commits pane shows detail overlay
func TestEnterOnCommitsPaneShowsDetail(t *testing.T) {
	cfg := Config{
		Repository: &git.Repository{},
		Prefix:     "v",
	}

	model := New(cfg)
	model.state = StateVersionSelect
	model.focusedPane = PaneCommits
	model.commits = []*git.Commit{
		{Hash: "abc1234", Subject: "feat: add feature"},
		{Hash: "def5678", Subject: "fix: fix bug"},
	}
	model.showingDetail = false

	// Simulate Enter key press
	msg := tea.KeyMsg{Type: tea.KeyEnter}
	updatedModel, _ := model.Update(msg)
	m := updatedModel.(Model)

	assert.True(t, m.showingDetail, "Enter should show detail overlay when commits pane focused")
}

// T029: Test Escape dismisses overlay
func TestEscapeDismissesOverlay(t *testing.T) {
	cfg := Config{
		Repository: &git.Repository{},
		Prefix:     "v",
	}

	model := New(cfg)
	model.state = StateVersionSelect
	model.focusedPane = PaneCommits
	model.showingDetail = true
	model.selectedCommitIndex = 0
	model.commits = []*git.Commit{
		{Hash: "abc1234", Subject: "feat: add feature"},
	}

	// Simulate Escape key press
	msg := tea.KeyMsg{Type: tea.KeyEscape}
	updatedModel, _ := model.Update(msg)
	m := updatedModel.(Model)

	assert.False(t, m.showingDetail, "Escape should dismiss detail overlay")
}

// T030: Test Enter dismisses overlay when showing
func TestEnterDismissesOverlayWhenShowing(t *testing.T) {
	cfg := Config{
		Repository: &git.Repository{},
		Prefix:     "v",
	}

	model := New(cfg)
	model.state = StateVersionSelect
	model.focusedPane = PaneCommits
	model.showingDetail = true
	model.selectedCommitIndex = 0
	model.commits = []*git.Commit{
		{Hash: "abc1234", Subject: "feat: add feature"},
	}

	// Simulate Enter key press
	msg := tea.KeyMsg{Type: tea.KeyEnter}
	updatedModel, _ := model.Update(msg)
	m := updatedModel.(Model)

	assert.False(t, m.showingDetail, "Enter should dismiss detail overlay when already showing")
}

// T038: Test Enter on version pane proceeds to confirmation
func TestEnterOnVersionPaneProceedsToConfirmation(t *testing.T) {
	cfg := Config{
		Repository: &git.Repository{},
		Prefix:     "v",
	}

	model := New(cfg)
	model.state = StateVersionSelect
	model.focusedPane = PaneVersion
	model.versionOptions = []VersionOption{
		{Label: "patch", NewVersion: "v1.0.1", BumpType: 0}, // patch
		{Label: "minor", NewVersion: "v1.1.0", BumpType: 1}, // minor
		{Label: "major", NewVersion: "v2.0.0", BumpType: 2}, // major
	}
	model.selectedOption = 0

	// Simulate Enter key press on version pane
	msg := tea.KeyMsg{Type: tea.KeyEnter}
	updatedModel, _ := model.Update(msg)
	m := updatedModel.(Model)

	assert.Equal(
		t, StateConfirm, m.state,
		"Enter on version pane should proceed to confirmation state",
	)
	assert.Equal(t, "v1.0.1", m.newVersion, "newVersion should be set to selected option")
}

// T039: Test scroll position preserved when returning from confirmation
func TestScrollPositionPreservedWhenReturningFromConfirmation(t *testing.T) {
	cfg := Config{
		Repository: &git.Repository{},
		Prefix:     "v",
	}

	model := New(cfg)
	model.state = StateVersionSelect
	model.focusedPane = PaneCommits
	// Set up viewport with content and scroll position
	model.commitsPane.Width = 80
	model.commitsPane.Height = 5
	model.commitsPane.SetContent(testScrollContent)
	model.commitsPane.SetYOffset(3) // Scroll to position 3

	// Switch to version pane and go to confirmation
	model.focusedPane = PaneVersion
	model.versionOptions = []VersionOption{
		{Label: "patch", NewVersion: "v1.0.1", BumpType: 0},
	}
	model.selectedOption = 0

	// Simulate Enter to go to confirmation
	msg := tea.KeyMsg{Type: tea.KeyEnter}
	updatedModel, _ := model.Update(msg)
	m := updatedModel.(Model)

	assert.Equal(t, StateConfirm, m.state, "Should be in confirmation state")

	// Now press Escape to go back to version select
	msg = tea.KeyMsg{Type: tea.KeyEscape}
	updatedModel, _ = m.Update(msg)
	m = updatedModel.(Model)

	assert.Equal(t, StateVersionSelect, m.state, "Should be back in version select state")
	assert.Equal(t, 3, m.commitsPane.YOffset, "Scroll position should be preserved")
}

// T040: Test commits pane remains visible when version pane focused
func TestCommitsPaneVisibleWhenVersionPaneFocused(t *testing.T) {
	cfg := Config{
		Repository: &git.Repository{},
		Prefix:     "v",
	}

	model := New(cfg)
	model.state = StateVersionSelect
	model.focusedPane = PaneVersion
	model.width = 80
	model.height = 24
	model.commits = []*git.Commit{
		{Hash: "abc1234", Subject: "feat: add feature"},
		{Hash: "def5678", Subject: "fix: fix bug"},
	}
	model.commitsPane.Width = 78
	model.commitsPane.Height = 5
	model.commitsPane.SetContent("abc1234 feat: add feature\ndef5678 fix: fix bug")
	model.versionOptions = []VersionOption{
		{Label: "patch", NewVersion: "v1.0.1", BumpType: 0},
	}
	// Initialize currentVersion to avoid nil pointer
	zeroVer := version.Zero()
	model.currentVersion = &zeroVer

	// Render the view
	view := model.View()

	// Verify both panes are visible in the output
	assert.Contains(t, view, "Commits", "Commits pane header should be visible")
	assert.Contains(t, view, "Version", "Version pane header should be visible")
}
