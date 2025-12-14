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

// Commit type styles for colored badges
var (
	FeatStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("154")) // Lime/Green

	FixStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("214")) // Yellow

	DocsStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("75")) // Blue

	ChoreStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("245")) // Gray

	RefactorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("81")) // Cyan

	TestStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("213")) // Magenta

	PerfStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("208")) // Orange

	BreakingStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("255")).
			Background(lipgloss.Color("196")) // Red background
)

// CommitTypeStyles maps commit types to their lipgloss styles
var CommitTypeStyles = map[string]lipgloss.Style{
	"feat":     FeatStyle,
	"fix":      FixStyle,
	"docs":     DocsStyle,
	"chore":    ChoreStyle,
	"refactor": RefactorStyle,
	"test":     TestStyle,
	"style":    ChoreStyle,
	"perf":     PerfStyle,
	"ci":       ChoreStyle,
	"build":    ChoreStyle,
}

// GetCommitTypeStyle returns the appropriate style for a commit type
func GetCommitTypeStyle(commitType string, isBreaking bool) lipgloss.Style {
	if isBreaking {
		// Apply breaking style (red background) to the type
		if baseStyle, ok := CommitTypeStyles[commitType]; ok {
			return baseStyle.Background(lipgloss.Color("196"))
		}
		return BreakingStyle
	}
	if style, ok := CommitTypeStyles[commitType]; ok {
		return style
	}
	return lipgloss.NewStyle() // Default unstyled
}

// Icons
const (
	IconCheck    = "✓"
	IconCross    = "✗"
	IconArrow    = "→"
	IconBullet   = "•"
	IconSelected = "❯"
)
