package installer

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/yousefbustamiii/package-mate/internal/components"
	"github.com/yousefbustamiii/package-mate/internal/sys"
)

// BrewJSON matches the JSON from `brew info --json=v2` (for single or all installed).
type BrewJSON struct {
	Formulae []struct {
		Name     string `json:"name"`
		Outdated bool   `json:"outdated"`
		Versions struct {
			Stable string `json:"stable"`
		} `json:"versions"`
		Installed []struct {
			Version            string `json:"version"`
			Time               int64  `json:"time"`
			InstalledOnRequest bool   `json:"installed_on_request"`
		} `json:"installed"`
		// Fields used by 'brew outdated'
		InstalledVersions []string `json:"installed_versions"`
		CurrentVersion    string   `json:"current_version"`
	} `json:"formulae"`
	Casks []struct {
		Token            string `json:"token"`
		Version          string `json:"version"`           // latest
		InstalledVersion string `json:"installed_version"` // installed
		InstalledTime    int64  `json:"installed_time"`
		Outdated         bool   `json:"outdated"`
		// Fields used by 'brew outdated'
		InstalledVersions string `json:"installed_versions"`
		CurrentVersion    string `json:"current_version"`
	} `json:"casks"`
}

// GetBrewInfo retrieves both version strings and installation date from Homebrew in one call.
func GetBrewInfo(item components.InstallItem) (versions []string, date time.Time, found bool) {
	if item.Formula == "" && item.Cask == "" {
		return nil, time.Time{}, false
	}

	arg := item.Formula
	if item.Cask != "" {
		arg = item.Cask
	}

	cmd := exec.Command(BrewExe(), "info", "--json=v2", arg)
	cmd.Env = append(os.Environ(), "HOMEBREW_NO_AUTO_UPDATE=1")
	out, err := cmd.Output()
	if err != nil {
		return nil, time.Time{}, false
	}

	var info BrewJSON
	if err := json.Unmarshal(out, &info); err != nil {
		return nil, time.Time{}, false
	}

	// Check formulae
	for _, f := range info.Formulae {
		if len(f.Installed) > 0 {
			vers := make([]string, len(f.Installed))
			for i, inst := range f.Installed {
				vers[i] = f.Name + " " + inst.Version
			}
			// Take the most recent installation time
			last := f.Installed[len(f.Installed)-1]
			var t time.Time
			if last.Time > 0 {
				t = time.Unix(last.Time, 0)
			}
			return vers, t, true
		}
	}

	// Check casks
	for _, c := range info.Casks {
		if c.InstalledVersion != "" || c.InstalledTime > 0 {
			ver := c.InstalledVersion
			if ver == "" {
				ver = c.Version
			}
			var t time.Time
			if c.InstalledTime > 0 {
				t = time.Unix(c.InstalledTime, 0)
			}
			return []string{c.Token + " " + ver}, t, true
		}
	}

	return nil, time.Time{}, false
}



// OutdatedStatus contains maps for quick lookup of outdated tools.
type OutdatedStatus struct {
	Formulae map[string]bool
	Casks    map[string]bool
}

// InstalledStatus contains maps of all installed tools and their versions.
type InstalledStatus struct {
	Formulae  map[string]string // formula -> version
	Casks     map[string]string // cask -> version
	Requested map[string]bool   // Name/Token -> true if installed on request (not as dependency)
}



// FetchFullBrewStatus retrieves installed tools and their outdated status in a single Homebrew call.
func FetchFullBrewStatus() (InstalledStatus, OutdatedStatus, error) {
	inst := InstalledStatus{
		Formulae:  make(map[string]string),
		Casks:     make(map[string]string),
		Requested: make(map[string]bool),
	}
	outdated := OutdatedStatus{
		Formulae: make(map[string]bool),
		Casks:    make(map[string]bool),
	}

	cmd := exec.Command(BrewExe(), "info", "--json=v2", "--installed")
	cmd.Env = append(os.Environ(), "HOMEBREW_NO_AUTO_UPDATE=1")
	out, err := cmd.Output()
	if err != nil {
		return inst, outdated, err
	}

	var info BrewJSON
	if err := json.Unmarshal(out, &info); err != nil {
		return inst, outdated, err
	}

	// 1. Map Formulae
	for _, f := range info.Formulae {
		if len(f.Installed) > 0 {
			// Get the version of the last (most recent) installation
			last := f.Installed[len(f.Installed)-1]
			inst.Formulae[f.Name] = last.Version
			if last.InstalledOnRequest {
				inst.Requested[f.Name] = true
			}
			if f.Outdated {
				outdated.Formulae[f.Name] = true
			}
		}
	}

	// 2. Map Casks
	for _, c := range info.Casks {
		ver := c.InstalledVersion
		if ver == "" {
			ver = c.Version
		}

		if ver != "" || c.InstalledVersion != "" {
			inst.Casks[c.Token] = ver
			// Check outdated flag or compare versions
			if c.Outdated || (c.Version != "" && c.InstalledVersion != "" && c.Version != c.InstalledVersion) {
				outdated.Casks[c.Token] = true
			}
		}
	}

	return inst, outdated, nil
}

// CheckOutdated returns (isOutdated, currentVersion, latestVersion)
func CheckOutdated(item components.InstallItem) (bool, string, string) {
	if item.Formula == "" && item.Cask == "" {
		return false, "", ""
	}

	arg := item.Formula
	if item.Cask != "" {
		arg = item.Cask
	}

	cmd := exec.Command(BrewExe(), "outdated", "--json=v2", arg)
	cmd.Env = append(os.Environ(), "HOMEBREW_NO_AUTO_UPDATE=1")
	out, _ := cmd.Output()
	// brew outdated exits with code 1 when the formula IS outdated — that's valid output.
	// Only bail if we received nothing at all (e.g. brew not found).
	if len(out) == 0 {
		return false, "", ""
	}

	var info BrewJSON
	if err := json.Unmarshal(out, &info); err != nil {
		return false, "", ""
	}

	// Check if it's an outdated formula
	if len(info.Formulae) > 0 && len(info.Formulae[0].InstalledVersions) > 0 {
		return true, info.Formulae[0].InstalledVersions[0], info.Formulae[0].CurrentVersion
	}

	// Check if it's an outdated cask
	if len(info.Casks) > 0 && info.Casks[0].InstalledVersions != "" {
		return true, info.Casks[0].InstalledVersions, info.Casks[0].CurrentVersion
	}

	return false, "", ""
}

// GetAllVersions returns a list of all detected versions (Managed, ManagedOlder, Unmanaged).
func GetAllVersions(item components.InstallItem) ([]VersionEntry, error) {
	var entries []VersionEntry
	uniqueResolved := make(map[string]bool)
	prefix := sys.BrewPrefix()

	// 1. Scan PATH for all instances of the primary binary
	if item.Binary != "" {
		pathStr := os.Getenv("PATH")
		dirs := strings.Split(pathStr, ":")
		for _, dir := range dirs {
			if dir == "" {
				continue
			}
			p := filepath.Join(dir, item.Binary)
			if _, err := os.Stat(p); err == nil {
				resolved, err := filepath.EvalSymlinks(p)
				if err != nil {
					resolved = p
				}
				resolved, _ = filepath.Abs(resolved)

				if uniqueResolved[resolved] {
					continue
				}
				uniqueResolved[resolved] = true

				t := VersionUnmanaged
				ver := "System/Manual"
				path := p
				formulaName := ""

				// Check if this resolved path is actually managed by Homebrew
				if strings.HasPrefix(resolved, prefix) {
					rel := strings.TrimPrefix(resolved, prefix)
					parts := strings.Split(rel, "/")

					for i, part := range parts {
						if (part == "Cellar" || part == "opt") && i+1 < len(parts) {
							formulaName = parts[i+1]

							// ❯ Helper: Strip tap prefix (e.g. timescale/tap/timescaledb -> timescaledb)
							baseFormula := item.Formula
							if strings.Contains(baseFormula, "/") {
								parts := strings.Split(baseFormula, "/")
								baseFormula = parts[len(parts)-1]
							}

							// Only mark as Managed if we successfully extracted a formula name
							isMain := formulaName == item.Formula || formulaName == baseFormula
							isVersioned := strings.HasPrefix(formulaName, item.Formula+"@") || strings.HasPrefix(formulaName, baseFormula+"@")

							if isMain {
								t = VersionManaged
							} else if isVersioned {
								t = VersionManagedOlder
							}

							if part == "Cellar" && i+2 < len(parts) {
								ver = formulaName + " " + parts[i+2]
							} else {
								ver = formulaName
							}
							break
						}
					}
				} else if strings.HasPrefix(p, prefix) && item.Cask != "" {
					t = VersionManaged
					formulaName = item.Cask
					ver = item.Cask
				}

				entries = append(entries, VersionEntry{
					Type:    t,
					Version: ver,
					Path:    path,
					Formula: formulaName,
				})
			}
		}
	}

	// 2. Also check 'brew info' as backup for unlinked brew versions
	if item.Formula != "" {
		cmd := exec.Command(BrewExe(), "info", "--json=v2", item.Formula)
		cmd.Env = append(os.Environ(), "HOMEBREW_NO_AUTO_UPDATE=1")
		out, err := cmd.Output()
		if err == nil {
			var info BrewJSON
			if json.Unmarshal(out, &info) == nil {
				for _, f := range info.Formulae {
					stableVer := f.Versions.Stable
					for _, inst := range f.Installed {
						// Check if we already have this version from the PATH scan
						prefixMatch := false
						verStr := f.Name + " " + inst.Version
						for _, e := range entries {
							if e.Version == verStr {
								prefixMatch = true
								break
							}
						}

						if !prefixMatch {
							vType := VersionManaged
							if stableVer != "" && inst.Version != stableVer {
								vType = VersionManagedOlder
							}
							entries = append(entries, VersionEntry{
								Type:    vType,
								Version: verStr,
								Path:    "(unlinked)",
								Formula: f.Name,
							})
						}
					}
				}
			}
		}
	}

	// 3. Check /Applications for unmanaged .app bundles.
	if item.Cask != "" {
		hasManagedEntry := false
		for _, e := range entries {
			if e.Type == VersionManaged || e.Type == VersionManagedOlder {
				hasManagedEntry = true
				break
			}
		}
		if !hasManagedEntry {
			var appPaths []string
			if path, ok := sys.AppExists(item.Name); ok {
				appPaths = append(appPaths, path)
			}
			// One level deep: e.g. pgAdmin inside /Applications/PostgreSQL 18/
			if matches, _ := filepath.Glob(filepath.Join("/Applications", "*", item.Name+"*.app")); len(matches) > 0 {
				appPaths = append(appPaths, matches...)
			}
			for _, appPath := range appPaths {
				entries = append(entries, VersionEntry{
					Type:    VersionUnmanaged,
					Version: filepath.Base(appPath),
					Path:    appPath,
				})
			}
		}
	}

	hasManaged := false
	for _, e := range entries {
		if e.Type == VersionManaged {
			hasManaged = true
			break
		}
	}
	if !hasManaged {
		for i, e := range entries {
			if e.Type == VersionManagedOlder && e.Path != "(unlinked)" {
				entries[i].Type = VersionManaged
				break
			}
		}
	}

	return entries, nil
}
