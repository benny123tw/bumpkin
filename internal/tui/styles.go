package tui

import (
	"github.com/charmbracelet/lipgloss"
)

// Colors
var (
	primaryColor   = lipgloss.Color("212") // Bright blue
	secondaryColor = lipgloss.Color("240") // Gray
	successColor   = lipgloss.Color("42")  // Green
	errorColor     = lipgloss.Color("196") // Red
	warningColor   = lipgloss.Color("214") // Orange
	mutedColor     = lipgloss.Color("245") // Light gray
)

// Styles
var (
	// Title style
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(primaryColor).
			MarginBottom(1)

	// Subtitle style
	SubtitleStyle = lipgloss.NewStyle().
			Foreground(secondaryColor).
			MarginBottom(1)

	// Version styles
	CurrentVersionStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(successColor)

	NewVersionStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(primaryColor)

	// Commit styles
	CommitHashStyle = lipgloss.NewStyle().
			Foreground(warningColor)

	CommitMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("255"))

	CommitAuthorStyle = lipgloss.NewStyle().
				Foreground(mutedColor)

	// Selection styles
	SelectedStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(primaryColor).
			Background(lipgloss.Color("236"))

	UnselectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("255"))

	// Status styles
	SuccessStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(successColor)

	ErrorStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(errorColor)

	WarningStyle = lipgloss.NewStyle().
			Foreground(warningColor)

	// Box style
	BoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(secondaryColor).
			Padding(1, 2)

	// Help style
	HelpStyle = lipgloss.NewStyle().
			Foreground(mutedColor).
			MarginTop(1)

	// Spinner style
	SpinnerStyle = lipgloss.NewStyle().
			Foreground(primaryColor)

	// Dry run indicator
	DryRunStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(warningColor).
			Background(lipgloss.Color("235")).
			Padding(0, 1)

	// Recommended indicator
	RecommendedStyle = lipgloss.NewStyle().
				Foreground(successColor).
				Italic(true)
)

// Icons
const (
	IconCheck    = "✓"
	IconCross    = "✗"
	IconArrow    = "→"
	IconBullet   = "•"
	IconSelected = "❯"
)
