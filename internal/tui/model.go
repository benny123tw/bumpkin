package tui

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/benny123tw/bumpkin/internal/conventional"
	"github.com/benny123tw/bumpkin/internal/executor"
	"github.com/benny123tw/bumpkin/internal/git"
	"github.com/benny123tw/bumpkin/internal/hooks"
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

	// Hook output streaming
	hookPane       *HookPane             // Scrollable output pane for hooks
	hookLineChan   chan hooks.OutputLine // Channel for receiving output lines
	hookDoneChan   chan hooks.HookResult // Channel for hook completion
	currentHooks   []hooks.Hook          // Current hooks being executed
	currentHookIdx int                   // Index of current hook in sequence
	hookPhase      hooks.HookType        // Current hook phase (pre-tag, post-tag, etc.)

	// Hook cancellation
	hookCancelFunc    context.CancelFunc // Function to cancel current hook
	cancelPending     bool               // Whether first ctrl+c was pressed
	cancelPendingTime time.Time          // When first ctrl+c was pressed

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

		// Resize hook pane if it exists
		if m.hookPane != nil {
			width, height := m.calculateHookPaneDimensions()
			m.hookPane.SetSize(width, height)
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

	case HookLineMsg:
		// Add line to hook pane
		if m.hookPane != nil {
			m.hookPane.AddLine(msg.Line)
		}
		// Continue listening for more lines
		return m, waitForHookLine(m.hookLineChan)

	case HookStartMsg:
		// Update hook pane with current hook info
		if m.hookPane != nil {
			m.hookPane.SetCurrentHook(msg.Hook, msg.Index, msg.Total)
		}
		return m, nil

	case HookCompleteMsg:
		if !msg.Success {
			// For post-push hooks, we use fail-open (warnings, not errors)
			if m.hookPhase == hooks.PostPush {
				// Add warning and continue
				if m.result != nil {
					m.result.PostPushWarnings = append(m.result.PostPushWarnings,
						fmt.Sprintf("hook '%s' failed: %v", msg.Hook.Command, msg.Error))
				}
			} else {
				m.err = msg.Error
				m.state = StateError
				return m, nil
			}
		}
		// Move to next hook or complete
		m.currentHookIdx++
		if m.currentHookIdx < len(m.currentHooks) {
			// Start next hook
			return m, m.startNextHook()
		}
		// All hooks complete, proceed with execution
		return m, m.continueExecution()

	case TagCreatedMsg:
		// Store result from message
		m.result = &executor.Result{
			TagName:    msg.TagName,
			CommitHash: msg.CommitHash,
			Pushed:     false,
		}
		// Tag created, run post-tag hooks
		return m, m.startPostTagHooks()

	case PushCompleteMsg:
		// Mark as pushed
		if m.result != nil {
			m.result.Pushed = true
		}
		// Push complete, run post-push hooks
		return m, m.startPostPushHooks()

	case ExecuteStartMsg:
		// Start the execution flow with pre-tag hooks
		return m, m.startPreTagHooks()
	}

	// Handle text input updates
	if m.state == StateCustomInput {
		var cmd tea.Cmd
		m.customInput, cmd = m.customInput.Update(msg)
		return m, cmd
	}

	return m, nil
}

// handleCtrlC handles ctrl+c key press with double-press cancellation for hooks
func (m Model) handleCtrlC() (tea.Model, tea.Cmd) {
	if m.state == StateDone || m.state == StateError {
		return m, tea.Quit
	}
	// Handle cancellation during hook execution
	if m.state == StateExecutingHooks {
		// Check if this is second ctrl+c within 3 seconds
		if m.cancelPending && time.Since(m.cancelPendingTime) < 3*time.Second {
			// Cancel the hook
			if m.hookCancelFunc != nil {
				m.hookCancelFunc()
			}
			m.cancelPending = false
			m.err = fmt.Errorf("hook cancelled by user")
			m.state = StateError
			return m, nil
		}
		// First ctrl+c - show warning
		m.cancelPending = true
		m.cancelPendingTime = time.Now()
		return m, nil
	}
	// Don't allow quit during non-hook execution
	if m.state == StateExecuting {
		return m, nil
	}
	return m, tea.Quit
}

func (m Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		return m.handleCtrlC()

	case "q":
		if m.state == StateDone || m.state == StateError {
			return m, tea.Quit
		}
		// Don't allow quit during execution states
		if m.state == StateExecuting || m.state == StateExecutingHooks {
			return m, nil
		}
		return m, tea.Quit

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
		case StateLoading, StateExecuting, StateExecutingHooks, StateDone, StateError:
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
	case StateLoading:
		// No key handling during loading
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
	case StateExecutingHooks:
		// Allow scrolling in hook output pane
		return m.handleHookExecutionKeys(msg)
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

func (m Model) handleHookExecutionKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Delegate scroll keys to hookPane
	if m.hookPane != nil {
		switch msg.String() {
		case keyUp, keyK, keyDown, keyJ, "g", "G", "ctrl+u", "ctrl+d":
			var cmd tea.Cmd
			m.hookPane, cmd = m.hookPane.Update(msg)
			return m, cmd
		}
	}
	return m, nil
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
	// Start pre-tag hooks (which will chain to tagging, post-tag, push, post-push)
	return ExecuteStartMsg{}
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

	case StateVersionSelect:
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

	case StateExecutingHooks:
		sb.WriteString(m.renderHookExecutionView())

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

func (m Model) renderHookExecutionView() string {
	var sb strings.Builder

	// Show current phase
	phaseLabel := "Running hooks"
	switch m.hookPhase {
	case hooks.PreTag:
		phaseLabel = "Running pre-tag hooks"
	case hooks.PostTag:
		phaseLabel = "Running post-tag hooks"
	case hooks.PostPush:
		phaseLabel = "Running post-push hooks"
	}

	sb.WriteString(m.spinner.View())
	sb.WriteString(" ")
	sb.WriteString(SubtitleStyle.Render(phaseLabel))
	sb.WriteString("\n")

	// Show cancellation warning if pending
	if m.cancelPending && time.Since(m.cancelPendingTime) < 3*time.Second {
		sb.WriteString(WarningStyle.Render("  Press ctrl+c again to cancel"))
		sb.WriteString("\n")
	}
	sb.WriteString("\n")

	// Render hook output pane
	if m.hookPane != nil {
		sb.WriteString(m.hookPane.View())
	}

	return sb.String()
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
	case StateVersionSelect:
		help = "↑/↓/j/k: navigate • h/l/tab: switch pane • gg/G: top/bottom • enter: select • q: quit"
	case StateCustomInput:
		help = "enter: confirm • esc: back"
	case StateConfirm:
		help = "↑/↓: select • enter: confirm • y/n: yes/no • esc: back"
	case StateExecuting:
		help = "please wait..."
	case StateExecutingHooks:
		help = "↑/↓/j/k: scroll output • running hooks..."
	case StateDone, StateError:
		// Help text already included in RenderSuccess/RenderError
		return ""
	}

	return HelpStyle.Render(help)
}

// waitForHookLine returns a tea.Cmd that waits for output lines from a hook
func waitForHookLine(lineChan chan hooks.OutputLine) tea.Cmd {
	return func() tea.Msg {
		line, ok := <-lineChan
		if !ok {
			return nil // Channel closed
		}
		return HookLineMsg{Line: line}
	}
}

// waitForHookDone returns a tea.Cmd that waits for hook completion
func waitForHookDone(doneChan chan hooks.HookResult) tea.Cmd {
	return func() tea.Msg {
		result, ok := <-doneChan
		if !ok {
			return nil // Channel closed
		}
		return HookCompleteMsg{
			Hook:     result.Hook,
			Success:  result.Success,
			Error:    result.Error,
			Duration: result.Duration,
		}
	}
}

// startNextHook starts the next hook in the sequence and returns the tea.Cmd
func (m *Model) startNextHook() tea.Cmd {
	if m.currentHookIdx >= len(m.currentHooks) {
		return nil
	}

	hook := m.currentHooks[m.currentHookIdx]

	// Reset cancel state for new hook
	m.cancelPending = false

	// Create hook context
	hookCtx := &hooks.HookContext{
		TagName:         m.newVersion,
		PreviousVersion: m.currentVersion.String(),
		Version:         strings.TrimPrefix(m.newVersion, m.config.Prefix),
		Prefix:          m.config.Prefix,
		Remote:          m.config.Remote,
		DryRun:          m.config.DryRun,
	}

	// Create cancellable context for hook execution
	ctx, cancel := context.WithCancel(context.Background())
	m.hookCancelFunc = cancel

	// Start streaming hook execution
	lineChan, doneChan := hooks.RunHookStreaming(ctx, hook, hookCtx)
	m.hookLineChan = lineChan
	m.hookDoneChan = doneChan

	// Update hookPane with current hook info
	if m.hookPane != nil {
		m.hookPane.SetCurrentHook(hook, m.currentHookIdx, len(m.currentHooks))
	}

	// Return commands to listen for output and completion, plus spinner
	return tea.Batch(
		m.spinner.Tick,
		waitForHookLine(lineChan),
		waitForHookDone(doneChan),
	)
}

// continueExecution proceeds with the execution after hooks complete
func (m *Model) continueExecution() tea.Cmd {
	// Determine what to do next based on hook phase
	switch m.hookPhase {
	case hooks.PreTag:
		// Pre-tag hooks done, now execute the actual tagging
		return m.doTagging
	case hooks.PostTag:
		// Post-tag hooks done, check if we need to push
		if m.config.NoPush || !m.hasRemote {
			m.state = StateDone
			return nil
		}
		// Push and then run post-push hooks
		return m.doPushAndPostPush
	case hooks.PostPush:
		// All done
		m.state = StateDone
		return nil
	default:
		m.state = StateDone
		return nil
	}
}

// calculateHookPaneDimensions returns the width and height for the hook pane
// based on the current window dimensions
func (m *Model) calculateHookPaneDimensions() (width, height int) {
	// Calculate pane dimensions - cap height to maintain consistent layout
	headerHeight := 4 // title + phase label + spacing
	footerHeight := 2 // help text
	availableHeight := m.height - headerHeight - footerHeight

	// Cap pane height to ~60% of available space (similar to version pane)
	height = availableHeight * 60 / 100
	if height < 10 {
		height = 10
	}
	if height > 20 {
		height = 20 // Max 20 lines to keep it manageable
	}

	width = m.width - 4
	if width < 40 {
		width = 40
	}

	return width, height
}

// initHookPane initializes the hook output pane with current dimensions
func (m *Model) initHookPane() {
	width, height := m.calculateHookPaneDimensions()
	m.hookPane = NewHookPane(width, height)
}

// startPreTagHooks initializes and starts pre-tag hook execution
func (m *Model) startPreTagHooks() tea.Cmd {
	if m.config.NoHooks || len(m.config.PreTagHooks) == 0 {
		// No pre-tag hooks, proceed directly to tagging
		return m.doTagging
	}

	// Initialize hook pane
	m.initHookPane()

	// Set up pre-tag hooks
	m.currentHooks = hooks.CreateHooks(m.config.PreTagHooks, hooks.PreTag)
	m.currentHookIdx = 0
	m.hookPhase = hooks.PreTag
	m.state = StateExecutingHooks

	return m.startNextHook()
}

// startPostTagHooks initializes and starts post-tag hook execution
func (m *Model) startPostTagHooks() tea.Cmd {
	if m.config.NoHooks || len(m.config.PostTagHooks) == 0 {
		// No post-tag hooks, check if we need to push
		if m.config.NoPush || !m.hasRemote {
			m.state = StateDone
			return nil
		}
		return m.doPushAndPostPush
	}

	// Initialize hook pane if not already done, or clear it for new phase
	if m.hookPane == nil {
		m.initHookPane()
	} else {
		m.hookPane.Clear() // Clear buffer for new phase
	}

	// Set up post-tag hooks
	m.currentHooks = hooks.CreateHooks(m.config.PostTagHooks, hooks.PostTag)
	m.currentHookIdx = 0
	m.hookPhase = hooks.PostTag
	m.state = StateExecutingHooks

	return m.startNextHook()
}

// startPostPushHooks initializes and starts post-push hook execution
func (m *Model) startPostPushHooks() tea.Cmd {
	if m.config.NoHooks || len(m.config.PostPushHooks) == 0 {
		m.state = StateDone
		return nil
	}

	// Initialize hook pane if not already done, or clear it for new phase
	if m.hookPane == nil {
		m.initHookPane()
	} else {
		m.hookPane.Clear() // Clear buffer for new phase
	}

	// Set up post-push hooks
	m.currentHooks = hooks.CreateHooks(m.config.PostPushHooks, hooks.PostPush)
	m.currentHookIdx = 0
	m.hookPhase = hooks.PostPush
	m.state = StateExecutingHooks

	return m.startNextHook()
}

// doTagging creates the git tag
func (m Model) doTagging() tea.Msg {
	newVerStr := strings.TrimPrefix(m.newVersion, m.config.Prefix)

	// Get HEAD commit hash first
	headHash, err := m.config.Repository.GetHEAD()
	if err != nil {
		return ErrorMsg{Err: fmt.Errorf("failed to get HEAD: %w", err)}
	}

	if m.config.DryRun {
		// Dry run - just pretend we created the tag
		return TagCreatedMsg{
			TagName:    m.newVersion,
			CommitHash: headHash.String(),
		}
	}

	// Create the tag
	err = m.config.Repository.CreateTag(m.newVersion, fmt.Sprintf("Release %s", newVerStr))
	if err != nil {
		return ErrorMsg{Err: fmt.Errorf("failed to create tag: %w", err)}
	}

	return TagCreatedMsg{
		TagName:    m.newVersion,
		CommitHash: headHash.String(),
	}
}

// doPushAndPostPush pushes the tag and starts post-push hooks
func (m Model) doPushAndPostPush() tea.Msg {
	if m.config.DryRun {
		return PushCompleteMsg{}
	}

	// Push the tag
	err := m.config.Repository.PushTag(m.newVersion, m.config.Remote)
	if err != nil {
		return ErrorMsg{Err: fmt.Errorf("failed to push tag: %w", err)}
	}

	return PushCompleteMsg{}
}

// Run starts the TUI
func Run(cfg Config) error {
	p := tea.NewProgram(New(cfg))
	_, err := p.Run()
	return err
}
