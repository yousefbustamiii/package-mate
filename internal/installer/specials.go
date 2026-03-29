package installer

import (
	"github.com/yousefbustamiii/package-mate/internal/components"
	"github.com/yousefbustamiii/package-mate/internal/installer/specials"
)

// isSpecialInstalled checks if a special-tagged tool is already on the system.
// brew-install and brew-update are handled here since they depend on this package.
func isSpecialInstalled(item components.InstallItem) (bool, string) {
	switch item.Special {
	case "brew-install":
		if IsBrewInstalled() {
			return true, brewVersion()
		}
		return false, ""
	case "brew-update":
		return false, ""
	}
	return specials.IsInstalled(item)
}

// dispatchSpecial routes a special-tagged item to its installer.
// brew-install and brew-update are handled here since they depend on this package.
func dispatchSpecial(item components.InstallItem) Result {
	switch item.Special {
	case "brew-install":
		return EnsureBrew()
	case "brew-update":
		return UpdateBrew()
	}
	if err := specials.Install(item); err != nil {
		return Result{ItemName: item.Name, Status: StatusFailed, Err: err}
	}
	return Result{ItemName: item.Name, Status: StatusInstalled}
}
