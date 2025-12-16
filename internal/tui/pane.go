package tui

// PaneType represents which pane has focus in the dual-pane layout
type PaneType int

const (
	// PaneVersion is the version selection pane (default focus)
	PaneVersion PaneType = iota
	// PaneCommits is the commit history pane
	PaneCommits
)

// String returns a string representation of the pane type
func (p PaneType) String() string {
	switch p {
	case PaneVersion:
		return "version"
	case PaneCommits:
		return "commits"
	default:
		return "unknown"
	}
}
