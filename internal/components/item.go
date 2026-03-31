package components

import (
	"encoding/json"
	"fmt"
	"strings"
)

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
	Name    string `json:"Name"`    // Display name shown in UI
	Desc    string `json:"Desc"`    // One-line description
	Formula string `json:"Formula"` // `brew install <formula>`
	Cask    string `json:"Cask"`    // `brew install --cask <cask>`
	Special string `json:"Special"` // special install method
	Color   string `json:"Color"`   // raw ANSI prefix from ui.RGB(), applied by caller
	Binary  string `json:"Binary"`  // binary name for LookPath check (e.g. "psql")
}

// Section groups related InstallItems under a named header.
type Section struct {
	Name  string        `json:"Name"`
	Items []InstallItem `json:"Items"`
}

// UnmarshalJSON parses the custom "rgb(r, g, b)" JSON string into an ANSI code string.
func (i *InstallItem) UnmarshalJSON(data []byte) error {
	type Alias InstallItem
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(i),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	
	// Convert "rgb(r, g, b)" or "rgb(r,g,b)" to ANSI if it matches
	if strings.HasPrefix(i.Color, "rgb(") && strings.HasSuffix(i.Color, ")") {
		inner := strings.TrimSuffix(strings.TrimPrefix(i.Color, "rgb("), ")")
		parts := strings.Split(inner, ",")
		if len(parts) == 3 {
			r := strings.TrimSpace(parts[0])
			g := strings.TrimSpace(parts[1])
			b := strings.TrimSpace(parts[2])
			
			var rInt, gInt, bInt int
			fmt.Sscanf(r, "%d", &rInt)
			fmt.Sscanf(g, "%d", &gInt)
			fmt.Sscanf(b, "%d", &bInt)
			i.Color = rgb(rInt, gInt, bInt)
		}
	}
	
	return nil
}
