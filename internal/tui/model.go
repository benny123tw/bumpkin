package tui

import (
	"context"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/benny123tw/bumpkin/internal/conventional"
	"github.com/benny123tw/bumpkin/internal/executor"
	"github.com/benny123tw/bumpkin/internal/git"
	"github.com/benny123tw/bumpkin/internal/version"
)

// Key constants
const (
	keyEnter = "enter"
	keySpace = " "
	keyUp    = "up"
	keyDown  = "down"
	keyK     = "k"
	keyJ     = "j"
)

// Config contains configuration for the TUI
type Config struct {
	Repository    *git.Repository
	Prefix        string
	Remote        string
	DryRun        bool
	NoPush        bool
	NoHooks       bool
	PreTagHooks   []string
	PostTagHooks  []string
	PostPushHooks []string
}

// Model is the main TUI model
type Model struct {
	config Config
	state  State
	err    error

	// Repository state
	currentVersion  *version.Version
	commits         []*git.Commit
	hasRemote       bool
	recommendedBump version.BumpType

	// Selection state
	versionOptions  []VersionOption
	selectedOption  int
	selectedConfirm int

	// Custom version input
	customInput textinput.Model

	// Selected bump
	selectedBumpType version.BumpType
	newVersion       string

	// Execution result
	result *executor.Result

	// UI components
	spinner spinner.Model

	// Dual-pane layout
	commitsPane         viewport.Model // Scrollable viewport for commits
	focusedPane         PaneType       // Which pane has focus (PaneVersion or PaneCommits)
	showingDetail       bool           // Whether commit detail overlay is shown
	selectedCommitIndex int            // Index of commit selected for detail view
	waitingForG         bool           // Whether we're waiting for second 'g' in 'gg' sequence

	// Window size
	width  int
	height int
}

// New creates a new TUI model
func New(cfg Config) Model {
	// Initialize spinner
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = SpinnerStyle

	// Initialize text input for custom version
	ti := textinput.New()
	ti.Placeholder = "e.g., 2.0.0"
	ti.CharLimit = 50
	ti.Width = 30

	return Model{
		config:          cfg,
		state:           StateLoading,
		spinner:         s,
		customInput:     ti,
		selectedOption:  0,
		selectedConfirm: 0,
		// Dual-pane layout initialization
		commitsPane:         viewport.New(0, 0), // Sized on WindowSizeMsg
		focusedPane:         PaneVersion,        // Default focus on version selection
		showingDetail:       false,
		selectedCommitIndex: 0,
	}
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		m.loadRepository,
	)
}

// loadRepository loads repository information
func (m Model) loadRepository() tea.Msg {
	// Get latest tag
	latestTag, err := m.config.Repository.LatestTag(m.config.Prefix)
	if err != nil {
		return ErrorMsg{Err: err}
	}

	var currentVersion *version.Version
	var commits []*git.Commit

	if latestTag != nil && latestTag.Version != nil {
		currentVersion = latestTag.Version
		commits, err = m.config.Repository.GetCommitsSinceTag(latestTag.Name)
		if err != nil {
			return ErrorMsg{Err: err}
		}
	} else {
		// No tags yet
		zero := version.Zero()
		currentVersion = &zero
		commits, err = m.config.Repository.GetAllCommits()
		if err != nil {
			return ErrorMsg{Err: err}
		}
	}

	// Check for remote
	hasRemote, _ := m.config.Repository.HasRemote(m.config.Remote)

	return RepoLoadedMsg{
		CurrentVersion: currentVersion,
		Commits:        commits,
		HasRemote:      hasRemote,
		RemoteName:     m.config.Remote,
	}
}

// Update handles messages
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeyPress(msg)

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		// Calculate pane heights for dual-pane layout
		// Layout: ~30% commits pane (top), ~70% version pane (bottom)
		headerHeight := 4 // title + dry run indicator + spacing
		footerHeight := 2 // help text
		availableHeight := m.height - headerHeight - footerHeight

		if availableHeight > 0 {
			// For small terminals (height < 16), we still show dual pane but minimal
			commitsPaneHeight := availableHeight * 30 / 100
			if commitsPaneHeight < 3 {
				commitsPaneHeight = 3 // Minimum height for usability
			}

			// Account for border (2 chars: top + bottom)
			m.commitsPane.Width = m.width - 2
			m.commitsPane.Height = commitsPaneHeight - 2
		}

		return m, nil

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case RepoLoadedMsg:
		m.currentVersion = msg.CurrentVersion
		m.commits = msg.Commits
		m.hasRemote = msg.HasRemote

		// Analyze commits for recommended bump
		var commitMessages []string
		for _, c := range m.commits {
			commitMessages = append(commitMessages, c.Message)
		}
		analysis := conventional.AnalyzeCommits(commitMessages)
		m.recommendedBump = analysis.RecommendedBump

		// Create version options with recommendation
		m.versionOptions = CreateVersionOptionsWithRecommendation(
			*m.currentVersion,
			m.config.Prefix,
			m.recommendedBump,
		)

		// Pre-select the recommended option
		for i, opt := range m.versionOptions {
			if opt.BumpType == m.recommendedBump {
				m.selectedOption = i
				break
			}
		}

		// Populate commits pane with rendered commit list
		m.commitsPane.SetContent(RenderCommitListForViewport(m.commits, m.selectedCommitIndex))

		// Go directly to version selection (skip commit preview screen)
		m.state = StateVersionSelect
		return m, nil

	case ExecuteResultMsg:
		m.result = msg.Result
		m.state = StateDone
		return m, nil

	case ErrorMsg:
		m.err = msg.Err
		m.state = StateError
		return m, nil
	}

	// Handle text input updates
	if m.state == StateCustomInput {
		var cmd tea.Cmd
		m.customInput, cmd = m.customInput.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		if m.state == StateDone || m.state == StateError {
			return m, tea.Quit
		}
		if m.state != StateExecuting {
			return m, tea.Quit
		}

	case "esc":
		// Dismiss overlay if showing
		if m.showingDetail {
			m.showingDetail = false
			return m, nil
		}

		switch m.state {
		case StateVersionSelect:
			// No action - this is the first screen now
		case StateCustomInput:
			m.state = StateVersionSelect
			m.customInput.Reset()
		case StateConfirm:
			m.state = StateVersionSelect
		case StateLoading, StateCommitList, StateExecuting, StateDone, StateError:
			// No action for these states
		}
		return m, nil

	case "tab", "shift+tab", "h", "l":
		// Toggle focus between panes in version select state
		if m.state == StateVersionSelect && !m.showingDetail {
			if m.focusedPane == PaneVersion {
				m.focusedPane = PaneCommits
			} else {
				m.focusedPane = PaneVersion
			}
			return m, nil
		}
	}

	switch m.state {
	case StateLoading, StateCommitList:
		// No key handling during loading (StateCommitList kept for backwards compatibility)
		return m, nil
	case StateVersionSelect:
		return m.handleVersionSelectKeys(msg)
	case StateCustomInput:
		return m.handleCustomInputKeys(msg)
	case StateConfirm:
		return m.handleConfirmKeys(msg)
	case StateExecuting:
		// No key handling during execution
		return m, nil
	case StateDone, StateError:
		return m, tea.Quit
	}

	return m, nil
}

func (m Model) handleVersionSelectKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Handle overlay dismiss first
	if m.showingDetail {
		if msg.String() == keyEnter || msg.String() == "esc" {
			m.showingDetail = false
			return m, nil
		}
		// Block other keys when overlay is showing
		return m, nil
	}

	// Route arrow keys based on focused pane
	if m.focusedPane == PaneCommits {
		// Handle gg sequence for jump to top
		if m.waitingForG {
			m.waitingForG = false
			if msg.String() == "g" && len(m.commits) > 0 {
				// gg: jump to top
				m.selectedCommitIndex = 0
				m.commitsPane.SetContent(
					RenderCommitListForViewport(m.commits, m.selectedCommitIndex),
				)
				m.commitsPane.SetYOffset(0)
				return m, nil
			}
			// Not 'g', fall through to normal handling
		}

		// When commits pane is focused, move selection and scroll viewport
		switch msg.String() {
		case keyUp, keyK:
			if m.selectedCommitIndex > 0 {
				m.selectedCommitIndex--
				// Update content to reflect new selection
				m.commitsPane.SetContent(
					RenderCommitListForViewport(m.commits, m.selectedCommitIndex),
				)
				// Scroll viewport to keep selection visible
				if m.selectedCommitIndex < m.commitsPane.YOffset {
					m.commitsPane.SetYOffset(m.selectedCommitIndex)
				}
			}
			return m, nil
		case keyDown, keyJ:
			if m.selectedCommitIndex < len(m.commits)-1 {
				m.selectedCommitIndex++
				// Update content to reflect new selection
				m.commitsPane.SetContent(
					RenderCommitListForViewport(m.commits, m.selectedCommitIndex),
				)
				// Scroll viewport to keep selection visible
				visibleEnd := m.commitsPane.YOffset + m.commitsPane.Height - 1
				if m.selectedCommitIndex > visibleEnd {
					m.commitsPane.SetYOffset(m.selectedCommitIndex - m.commitsPane.Height + 1)
				}
			}
			return m, nil
		case "g":
			// Start gg sequence
			m.waitingForG = true
			return m, nil
		case "G":
			// Jump to bottom
			if len(m.commits) > 0 {
				m.selectedCommitIndex = len(m.commits) - 1
				m.commitsPane.SetContent(
					RenderCommitListForViewport(m.commits, m.selectedCommitIndex),
				)
				// Scroll to show selection at bottom of viewport
				newOffset := m.selectedCommitIndex - m.commitsPane.Height + 1
				if newOffset < 0 {
					newOffset = 0
				}
				m.commitsPane.SetYOffset(newOffset)
			}
			return m, nil
		case keyEnter:
			// Enter on commits pane shows detail overlay
			if len(m.commits) > 0 && m.selectedCommitIndex < len(m.commits) {
				m.showingDetail = true
			}
			return m, nil
		}
		return m, nil
	}

	// Version pane is focused - handle version selection
	switch msg.String() {
	case keyUp, keyK:
		if m.selectedOption > 0 {
			m.selectedOption--
		}
	case keyDown, keyJ:
		if m.selectedOption < len(m.versionOptions)-1 {
			m.selectedOption++
		}
	case keyEnter, keySpace:
		selected := m.versionOptions[m.selectedOption]
		m.selectedBumpType = selected.BumpType

		if selected.BumpType == version.BumpCustom {
			m.state = StateCustomInput
			m.customInput.Focus()
			return m, textinput.Blink
		}

		m.newVersion = selected.NewVersion
		m.state = StateConfirm
		m.selectedConfirm = 0
	}
	return m, nil
}

func (m Model) handleCustomInputKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Handle Enter key for submission
	if msg.String() == keyEnter {
		customVer := m.customInput.Value()
		if customVer == "" {
			return m, nil
		}

		// Validate version
		_, err := version.Parse(customVer)
		if err != nil {
			m.err = fmt.Errorf("invalid version: %s", customVer)
			return m, nil
		}

		// Add prefix if not present
		if !strings.HasPrefix(customVer, m.config.Prefix) {
			m.newVersion = m.config.Prefix + customVer
		} else {
			m.newVersion = customVer
		}

		m.state = StateConfirm
		m.selectedConfirm = 0
		return m, nil
	}

	// Pass all other keys to the text input
	var cmd tea.Cmd
	m.customInput, cmd = m.customInput.Update(msg)
	return m, cmd
}

func (m Model) handleConfirmKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case keyUp, keyK, keyDown, keyJ, "tab":
		m.selectedConfirm = 1 - m.selectedConfirm // Toggle between 0 and 1
	case keyEnter, keySpace:
		if m.selectedConfirm == 0 {
			// Confirmed - execute
			m.state = StateExecuting
			return m, m.executeVersion
		}
		// Cancelled
		m.state = StateVersionSelect
	case "y", "Y":
		m.state = StateExecuting
		return m, m.executeVersion
	case "n", "N":
		m.state = StateVersionSelect
	}
	return m, nil
}

func (m Model) executeVersion() tea.Msg {
	// Parse the new version (strip prefix)
	newVerStr := strings.TrimPrefix(m.newVersion, m.config.Prefix)

	req := executor.Request{
		Repository:    m.config.Repository,
		BumpType:      m.selectedBumpType,
		CustomVersion: newVerStr,
		Prefix:        m.config.Prefix,
		Remote:        m.config.Remote,
		DryRun:        m.config.DryRun,
		NoPush:        m.config.NoPush || !m.hasRemote,
		NoHooks:       m.config.NoHooks,
		PreTagHooks:   m.config.PreTagHooks,
		PostTagHooks:  m.config.PostTagHooks,
		PostPushHooks: m.config.PostPushHooks,
	}

	result, err := executor.Execute(context.Background(), req)
	if err != nil {
		return ErrorMsg{Err: err}
	}

	return ExecuteResultMsg{Result: result}
}

// View renders the UI
func (m Model) View() string {
	var sb strings.Builder

	// Header
	sb.WriteString(TitleStyle.Render("🎃 bumpkin"))
	sb.WriteString("\n")

	if m.config.DryRun {
		sb.WriteString(DryRunStyle.Render(" DRY RUN "))
		sb.WriteString("\n")
	}
	sb.WriteString("\n")

	switch m.state {
	case StateLoading:
		sb.WriteString(m.spinner.View())
		sb.WriteString(" Loading repository...")

	case StateCommitList, StateVersionSelect:
		// StateCommitList is deprecated - both render the same view now
		sb.WriteString(m.renderVersionSelectView())

		// Render overlay on top if showing detail
		if m.showingDetail && len(m.commits) > 0 && m.selectedCommitIndex < len(m.commits) {
			commit := m.commits[m.selectedCommitIndex]
			overlay := RenderCommitDetailOverlay(commit, m.width, m.height)
			return overlay
		}

	case StateCustomInput:
		sb.WriteString(m.renderCustomInputView())

	case StateConfirm:
		sb.WriteString(m.renderConfirmView())

	case StateExecuting:
		sb.WriteString(m.spinner.View())
		sb.WriteString(" Creating tag...")

	case StateDone:
		sb.WriteString(m.renderDoneView())

	case StateError:
		sb.WriteString(RenderError(m.err))
	}

	// Footer with help
	sb.WriteString("\n")
	sb.WriteString(m.renderHelp())

	return sb.String()
}

func (m Model) renderVersionSelectView() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Current version: %s\n\n",
		CurrentVersionStyle.Render(m.currentVersion.StringWithPrefix(m.config.Prefix)),
	))

	// Dual-pane layout: commits pane (top) + version pane (bottom)

	// Commits pane with border based on focus
	commitsBorderStyle := UnfocusedBorderStyle
	if m.focusedPane == PaneCommits {
		commitsBorderStyle = FocusedBorderStyle
	}

	// Build commits pane header with selection indicator
	var commitsHeader string
	if len(m.commits) > 0 {
		// Show selected commit position
		commitsHeader = fmt.Sprintf(
			" Commits [%d/%d] ",
			m.selectedCommitIndex+1, len(m.commits),
		)
	} else {
		commitsHeader = " Commits "
	}

	// Render commits pane content
	var commitsContent string
	if len(m.commits) > 0 {
		commitsContent = m.commitsPane.View()
	} else {
		commitsContent = WarningStyle.Render("No new commits")
	}

	// Apply border and width to commits pane
	commitsPaneWidth := m.width - 2 // Account for left/right margins
	if commitsPaneWidth < 20 {
		commitsPaneWidth = 20
	}
	commitsBox := commitsBorderStyle.
		Width(commitsPaneWidth).
		Render(commitsHeader + "\n" + commitsContent)

	sb.WriteString(commitsBox)
	sb.WriteString("\n")

	// Version pane with border based on focus
	versionBorderStyle := UnfocusedBorderStyle
	if m.focusedPane == PaneVersion {
		versionBorderStyle = FocusedBorderStyle
	}

	// Build version pane content
	versionHeader := " Version "
	versionContent := SubtitleStyle.Render(" Select version bump:") + "\n" +
		RenderVersionSelector(m.versionOptions, m.selectedOption)

	// Apply border and width to version pane
	versionBox := versionBorderStyle.
		Width(commitsPaneWidth).
		Render(versionHeader + "\n" + versionContent)

	sb.WriteString(versionBox)

	return sb.String()
}

func (m Model) renderCustomInputView() string {
	var sb strings.Builder

	sb.WriteString(SubtitleStyle.Render("Enter custom version:"))
	sb.WriteString("\n\n")
	sb.WriteString(fmt.Sprintf("  %s%s\n", m.config.Prefix, m.customInput.View()))

	if m.err != nil {
		sb.WriteString("\n")
		sb.WriteString(ErrorStyle.Render(fmt.Sprintf("  %s", m.err.Error())))
		// Note: error is displayed but not cleared here since View() is immutable
		// The error will be cleared on the next key press that sets a new error or succeeds
	}

	return sb.String()
}

func (m Model) renderConfirmView() string {
	return RenderConfirmation(
		m.currentVersion.String(),
		strings.TrimPrefix(m.newVersion, m.config.Prefix),
		m.newVersion,
		len(m.commits),
		m.config.Remote,
		m.config.NoPush,
		m.config.DryRun,
		m.selectedConfirm,
	)
}

func (m Model) renderDoneView() string {
	if m.result == nil {
		return ""
	}

	summary := &ExecutionSummary{
		TagName:          m.result.TagName,
		CommitHash:       m.result.CommitHash,
		Pushed:           m.result.Pushed,
		Remote:           m.config.Remote,
		PostPushWarnings: m.result.PostPushWarnings,
	}

	return RenderSuccess(summary)
}

func (m Model) renderHelp() string {
	var help string

	switch m.state {
	case StateLoading:
		help = "loading..."
	case StateCommitList:
		// Deprecated - kept for backwards compatibility
		help = "enter: select version • q: quit"
	case StateVersionSelect:
		help = "↑/↓/j/k: navigate • h/l/tab: switch pane • gg/G: top/bottom • enter: select • q: quit"
	case StateCustomInput:
		help = "enter: confirm • esc: back"
	case StateConfirm:
		help = "↑/↓: select • enter: confirm • y/n: yes/no • esc: back"
	case StateExecuting:
		help = "please wait..."
	case StateDone, StateError:
		// Help text already included in RenderSuccess/RenderError
		return ""
	}

	return HelpStyle.Render(help)
}

// Run starts the TUI
func Run(cfg Config) error {
	p := tea.NewProgram(New(cfg))
	_, err := p.Run()
	return err
}
