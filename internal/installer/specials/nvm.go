package specials

import (
	"fmt"
	"os"
	"strings"

	"github.com/yousefbustamiii/package-mate/internal/sys"
	"github.com/yousefbustamiii/package-mate/internal/ui"
)

func nvmCheck() (bool, string) {
	nvmDir := os.ExpandEnv("$HOME/.nvm")
	if info, err := os.Stat(nvmDir); err == nil && info.IsDir() {
		return true, "~/.nvm exists"
	}
	return false, ""
}

func nvmInstall() error {
	ui.Doing("Installing %s", "NVM")
	cmd := sys.ShellCommand(`curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.7/install.sh | bash`)
	var stderr strings.Builder
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%s", strings.TrimSpace(stderr.String()))
	}
	ui.Done("Successfully installed %s", "NVM")
	return nil
}
