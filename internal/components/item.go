package components

// DashboardStatus represents the installation/update state for the dashboard grid.
type DashboardStatus int

const (
	StatusNotInstalled DashboardStatus = iota
	StatusInstalled
	StatusOutdated
	StatusUnmanaged
)

// InstallItem represents a single installable tool.
type InstallItem struct {
	Name    string // Display name shown in UI
	Desc    string // One-line description
	Formula string // `brew install <formula>`
	Cask    string // `brew install --cask <cask>`
	Special string // special install method: "nvm", "claude", "gemini", "brew-install", "brew-update", etc.
	Color   string // raw ANSI prefix from ui.RGB(), applied by caller
	Binary  string // binary name for LookPath check (e.g. "psql")
}

// Section groups related InstallItems under a named header.
type Section struct {
	Name  string
	Items []InstallItem
}
