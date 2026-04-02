package components

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const (
	CatalogURL     = "https://package-mate.com/catalog.json"
	CatalogRefresh = 12 * time.Hour
	FetchTimeout   = 5 * time.Second
)

//go:embed data/catalog.json
var catalogFS embed.FS

// AllSections is the master ordered catalog of every installable tool.
var AllSections []Section

// CatalogMetadata tracks the sync status of the remote catalog.
type CatalogMetadata struct {
	LastFetch time.Time `json:"last_fetch"`
	ETag      string    `json:"etag"`
}

func init() {
	loadCatalog()
}

func configDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".package-mate")
}

func catalogPath() string {
	return filepath.Join(configDir(), "catalog.json")
}

func metadataPath() string {
	return filepath.Join(configDir(), "catalog.metadata.json")
}

// ClearCache removes the local catalog and metadata files.
// Returns true if anything was actually deleted, and an error if deletion failed.
func ClearCache() (bool, error) {
	catPath := catalogPath()
	metaPath := metadataPath()

	_, catErr := os.Stat(catPath)
	_, metaErr := os.Stat(metaPath)

	exists := catErr == nil || metaErr == nil
	if !exists {
		return false, nil
	}

	_ = os.Remove(catPath)
	_ = os.Remove(metaPath)
	return true, nil
}

// ForceUpdate performs a synchronous remote sync regardless of the 24h rule.
func ForceUpdate() error {
	metaData, metaErr := os.ReadFile(metadataPath())
	var meta CatalogMetadata
	if metaErr == nil {
		_ = json.Unmarshal(metaData, &meta)
	}

	newCatalog, newMeta, err := fetchRemoteCatalog(meta.ETag)
	if err != nil {
		return err
	}
	
	if newCatalog != nil {
		_ = os.WriteFile(catalogPath(), newCatalog, 0644)
		_ = saveMetadata(newMeta)
		// Reload current run
		var secs []Section
		if err := json.Unmarshal(newCatalog, &secs); err == nil {
			AllSections = sortSections(secs)
		}
	} else {
		// 304 - Just update timestamp
		meta.LastFetch = time.Now()
		_ = saveMetadata(meta)
	}
	return nil
}

func loadCatalog() {
	_ = os.MkdirAll(configDir(), 0755)

	// 1. Load from local cache (fast)
	localData, err := os.ReadFile(catalogPath())
	if err == nil {
		var secs []Section
		if err := json.Unmarshal(localData, &secs); err == nil {
			AllSections = sortSections(secs)
			return
		}
	}

	// 2. Clear or Missing Cache? Use Embed (fast)
	loadFromEmbed()
}

// CheckForUpdates handles the 24-hour sync rule.
// It is intended to be called in a background goroutine from main.go.
func CheckForUpdates() {
	metaData, metaErr := os.ReadFile(metadataPath())
	var meta CatalogMetadata
	if metaErr == nil {
		_ = json.Unmarshal(metaData, &meta)
	}

	if time.Since(meta.LastFetch) < CatalogRefresh {
		return
	}

	newCatalog, newMeta, syncErr := fetchRemoteCatalog(meta.ETag)
	if syncErr == nil && newCatalog != nil {
		_ = os.WriteFile(catalogPath(), newCatalog, 0644)
		_ = saveMetadata(newMeta)
	} else if syncErr == nil && newCatalog == nil {
		// 304 Not Modified
		meta.LastFetch = time.Now()
		_ = saveMetadata(meta)
	}
}

func saveMetadata(m CatalogMetadata) error {
	data, _ := json.MarshalIndent(m, "", "  ")
	return os.WriteFile(metadataPath(), data, 0644)
}

func fetchRemoteCatalog(currentETag string) ([]byte, CatalogMetadata, error) {
	ctx, cancel := context.WithTimeout(context.Background(), FetchTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", CatalogURL, nil)
	if err != nil {
		return nil, CatalogMetadata{}, err
	}

	if currentETag != "" {
		req.Header.Set("If-None-Match", currentETag)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, CatalogMetadata{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotModified {
		return nil, CatalogMetadata{}, nil
	}

	if resp.StatusCode != http.StatusOK {
		return nil, CatalogMetadata{}, fmt.Errorf("server returned %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, CatalogMetadata{}, err
	}

	newMeta := CatalogMetadata{
		LastFetch: time.Now(),
		ETag:      resp.Header.Get("ETag"),
	}

	return data, newMeta, nil
}

func loadFromEmbed() {
	data, err := catalogFS.ReadFile("data/catalog.json")
	if err != nil {
		return
	}
	var secs []Section
	_ = json.Unmarshal(data, &secs)
	AllSections = sortSections(secs)
}

func sortSections(secs []Section) []Section {
	var orderedNames = []string{
		"Homebrew Setup", "Databases", "Caching & Messaging", "Containers & DevOps",
		"Backend & Runtime", "Languages & Runtimes", "DB GUI & Dev Tools",
		"Coding CLIs & AI", "Dev Essentials", "Package Managers",
		"Testing & Utilities", "System Tools", "Infrastructure & Cloud",
		"Editors & IDEs", "Security & Secrets", "Media & Graphics",
		"Low-Level & Embedded", "Reverse Engineering", "Docs & Static Sites",
		"Performance & Profiling", "Virtualization", "Terminal Glow-Up & Shell",
		"macOS Essentials", "Cloud & Kubernetes", "Data & Analytics",
		"Modern CLI Replacements", "Networking & API", "Dev Utilities",
		"Design & Frontend", "Web Browsers", "Communications",
		"Knowledge & Productivity", "Terminal Emulators", "Developer Apps & Tooling", "JetBrains Mastery",
	}

	secsByName := make(map[string]Section)
	for _, s := range secs {
		secsByName[s.Name] = s
	}

	var sorted []Section
	for _, name := range orderedNames {
		if s, ok := secsByName[name]; ok {
			sorted = append(sorted, s)
		}
	}
	return sorted
}
