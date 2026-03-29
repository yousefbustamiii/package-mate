package specials

import (
	"fmt"

	"github.com/yousefbustamiii/package-mate/internal/components"
)

// IsInstalled reports whether a special-tagged tool is already on the system.
func IsInstalled(item components.InstallItem) (bool, string) {
	switch item.Special {
	case "nvm":
		return nvmCheck()
	case "rustup":
		return rustCheck()
	case "npm-check":
		return npmCheck()
	case "claude":
		return npmGCheck("@anthropic-ai/claude-code")
	case "gemini":
		return npmGCheck("@google/gemini-cli")
	case "jest":
		return npmGCheck("jest")
	case "npm-g":
		return npmGCheck(item.Formula)
	case "pytest":
		return pytestCheck()
	case "pipx-g":
		return pipxCheck(item.Formula)
	}
	return false, ""
}

// Install runs the installer for the given special-tagged item.
// Returns nil on success, error on failure.
func Install(item components.InstallItem) error {
	switch item.Special {
	case "nvm":
		return nvmInstall()
	case "rustup":
		return rustInstall()
	case "npm-check":
		return npmInstall()
	case "claude":
		return npmGInstall("Claude Code", "@anthropic-ai/claude-code")
	case "gemini":
		return npmGInstall("Gemini CLI", "@google/gemini-cli")
	case "jest":
		return npmGInstall("Jest", "jest")
	case "npm-g":
		return npmGInstall(item.Name, item.Formula)
	case "pytest":
		return pytestInstall()
	case "pipx-g":
		return pipxInstall(item.Name, item.Formula)
	}
	return fmt.Errorf("unknown special: %s", item.Special)
}
