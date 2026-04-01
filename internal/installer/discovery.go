package installer

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/yousefbustamiii/package-mate/internal/components"
	"github.com/yousefbustamiii/package-mate/internal/sys"
)

// SystemScan represents a point-in-time snapshot of the system's PATH and Homebrew inventory.
type SystemScan struct {
	PathIndex map[string][]string // binary name -> list of absolute paths (pre-resolved)
	AppNames  []string            // List of detected .app names across common locations
	Installed InstalledStatus     // Homebrew formulae and casks
	Outdated  OutdatedStatus      // Tools needing updates
	BrewRoot  string              // Homebrew prefix
}

// PerformSystemScan executes a single-pass scan of the user's PATH and Homebrew.
func PerformSystemScan() *SystemScan {
	scan := &SystemScan{
		PathIndex: make(map[string][]string),
		BrewRoot:  sys.BrewPrefix(),
	}

	// 1. Fetch all Homebrew data in one pass (Installed + Outdated)
	scan.Installed, scan.Outdated, _ = FetchFullBrewStatus()

	// 2. Scan PATH exactly once and pre-resolve symlinks
	pathStr := os.Getenv("PATH")
	dirs := strings.Split(pathStr, ":")
	for _, dir := range dirs {
		if dir == "" {
			continue
		}
		entries, err := os.ReadDir(dir)
		if err != nil {
			continue
		}

		for _, entry := range entries {
			if !entry.IsDir() {
				name := entry.Name()
				full := filepath.Join(dir, name)

				// Pre-resolve symlink for the PATH entry to avoid expensive calls later
				resolved, err := filepath.EvalSymlinks(full)
				if err != nil {
					resolved = full
				}
				abs, _ := filepath.Abs(resolved)
				scan.PathIndex[name] = append(scan.PathIndex[name], abs)
			}
		}
	}

	// 3. Index Applications (for unmanaged GUI apps)
	appDirs := []string{"/Applications", filepath.Join(os.Getenv("HOME"), "Applications")}
	for _, appDir := range appDirs {
		entries, err := os.ReadDir(appDir)
		if err != nil {
			continue
		}
		for _, entry := range entries {
			name := entry.Name()
			if strings.HasSuffix(name, ".app") {
				scan.AppNames = append(scan.AppNames, name)
			} else if entry.IsDir() {
				subDir := filepath.Join(appDir, name)
				subEntries, _ := os.ReadDir(subDir)
				for _, sub := range subEntries {
					if strings.HasSuffix(sub.Name(), ".app") {
						scan.AppNames = append(scan.AppNames, sub.Name())
					}
				}
			}
		}
	}

	return scan
}

// ResolveStatus determines the DashboardStatus and whether multiple versions exist.
//
//goland:noinspection GoDfaInspectionRunner
func ResolveStatus(scan *SystemScan, item components.InstallItem) (components.DashboardStatus, bool, bool) {
	status := components.StatusNotInstalled
	isRequested := false
	managedFormulas := make(map[string]bool)
	uniqueUnmanaged := make(map[string]bool)

	// 1. Check Managed Formulae (Exact + Older/Versioned)
	if item.Formula != "" {
		// Exact Match (e.g. "postgresql")
		if _, exists := scan.Installed.Formulae[item.Formula]; exists {
			managedFormulas[item.Formula] = true
			if scan.Installed.Requested[item.Formula] {
				isRequested = true
			}
			if scan.Outdated.Formulae[item.Formula] {
				status = components.StatusOutdated
			} else if status == components.StatusNotInstalled {
				status = components.StatusInstalled
			}
		}
		// Versioned Matches (e.g. "postgresql@18")
		for name := range scan.Installed.Formulae {
			if strings.HasPrefix(name, item.Formula+"@") {
				managedFormulas[name] = true
				if scan.Installed.Requested[name] {
					isRequested = true
				}
				if !scan.Outdated.Formulae[name] {
					if status == components.StatusNotInstalled || status == components.StatusOutdated {
						status = components.StatusInstalled
					}
				} else {
					if status == components.StatusNotInstalled {
						status = components.StatusOutdated
					}
				}
			}
		}
	}

	// 2. Count Managed Casks
	if item.Cask != "" {
		if _, exists := scan.Installed.Casks[item.Cask]; exists {
			managedFormulas[item.Cask] = true
			isRequested = true // Assume casks are always on request
			if scan.Outdated.Casks[item.Cask] {
				status = components.StatusOutdated
			} else if status == components.StatusNotInstalled {
				status = components.StatusInstalled
			}
		}
	}

	// 3. Scan PATH for Unmanaged Binaries (using pre-resolved indices)
	if item.Binary != "" {
		for _, resolved := range scan.PathIndex[item.Binary] {
			// Is it in Homebrew? (Check prefix)
			isBrew := strings.HasPrefix(resolved, scan.BrewRoot)
			if isBrew {
				if status == components.StatusNotInstalled {
					status = components.StatusInstalled
					// We don't mark as requested here specifically, but 
					// we could check common brew prefixes like /opt/homebrew/Cellar/<item.Formula>
				}
			}

			if !isBrew {
				if !uniqueUnmanaged[resolved] {
					uniqueUnmanaged[resolved] = true
					if status == components.StatusNotInstalled {
						status = components.StatusUnmanaged
					}
				}
			}
		}
	}

	// 4. Check for unmanaged .app bundles (using pre-indexed AppNames)
	if status == components.StatusNotInstalled && item.Cask != "" {
		searchName := strings.ToLower(item.Name)
		found := false

		for _, appName := range scan.AppNames {
			lowerApp := strings.ToLower(appName)
			// Match if item name is part of app name (e.g. "PostgreSQL" matches "PostgreSQL 16.app")
			if strings.Contains(lowerApp, searchName) {
				found = true
				break
			}
		}

		if found {
			status = components.StatusUnmanaged
		}
	}

	totalCount := len(managedFormulas) + len(uniqueUnmanaged)
	isMultiple := totalCount >= 2
	return status, isMultiple, isRequested
}
