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
	PathIndex map[string][]string // binary name -> list of absolute paths
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

	// 1. Fetch Outdated (Managed)
	scan.Outdated, _ = AllOutdated()

	// 2. Fetch Installed (Managed)
	scan.Installed, _ = AllInstalled()

	// 3. Scan PATH exactly once
	pathStr := os.Getenv("PATH")
	dirs := strings.Split(pathStr, ":")
	for _, dir := range dirs {
		if dir == "" {
			continue
		}
		// Read directory once
		entries, err := os.ReadDir(dir)
		if err != nil {
			continue
		}

		for _, entry := range entries {
			if !entry.IsDir() {
				name := entry.Name()
				full := filepath.Join(dir, name)
				scan.PathIndex[name] = append(scan.PathIndex[name], full)
			}
		}
	}

	return scan
}

// ResolveStatus determines the DashboardStatus and whether multiple versions exist.
//
//goland:noinspection GoDfaInspectionRunner
func ResolveStatus(scan *SystemScan, item components.InstallItem) (components.DashboardStatus, bool) {
	status := components.StatusNotInstalled
	managedFormulas := make(map[string]bool)
	uniqueUnmanaged := make(map[string]bool)

	// 1. Check Managed Formulae (Exact + Older/Versioned)
	if item.Formula != "" {
		// Exact Match (e.g. "postgresql")
		if _, exists := scan.Installed.Formulae[item.Formula]; exists {
			managedFormulas[item.Formula] = true
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
			if scan.Outdated.Casks[item.Cask] {
				status = components.StatusOutdated
			} else if status == components.StatusNotInstalled {
				status = components.StatusInstalled
			}
		}
	}

	// 3. Scan PATH for Unmanaged Binaries
	if item.Binary != "" {
		for _, p := range scan.PathIndex[item.Binary] {
			resolved, err := filepath.EvalSymlinks(p)
			if err != nil {
				resolved = p
			}
			resolved, _ = filepath.Abs(resolved)

			// Is it in Homebrew?
			isBrew := strings.HasPrefix(resolved, scan.BrewRoot) || strings.HasPrefix(p, scan.BrewRoot)
			if isBrew {
				if status == components.StatusNotInstalled {
					status = components.StatusInstalled
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

	// 4. Check for unmanaged .app bundles (installed via DMG, not Homebrew)
	if status == components.StatusNotInstalled && item.Cask != "" {
		if _, ok := sys.AppExists(item.Name); ok {
			status = components.StatusUnmanaged
		} else {
			matches, _ := filepath.Glob(filepath.Join("/Applications", "*", item.Name+"*.app"))
			if len(matches) > 0 {
				status = components.StatusUnmanaged
			}
		}
	}

	totalCount := len(managedFormulas) + len(uniqueUnmanaged)
	isMultiple := totalCount >= 2
	return status, isMultiple
}
