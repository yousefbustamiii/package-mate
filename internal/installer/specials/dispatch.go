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
	case "aider":
		return pipxCheck("aider-chat")
	case "continue-cli":
		return npmGCheck("@continuedev/cli")
	case "corepack":
		return npmGCheck("corepack")
	case "playwright":
		return npmGCheck("playwright")
	case "cypress":
		return npmGCheck("cypress")
	case "vitest":
		return npmGCheck("vitest")
	case "eslint":
		return npmGCheck("eslint")
	case "locust":
		return pipxCheck("locust")
	case "firebase":
		return npmGCheck("firebase-tools")
	case "prism":
		return npmGCheck("@stoplight/prism-cli")
	case "eleventy":
		return npmGCheck("@11ty/eleventy")
	case "spectral":
		return npmGCheck("@stoplight/spectral-cli")
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
	case "aider":
		return pipxInstall("Aider", "aider-chat")
	case "continue-cli":
		return npmGInstall("Continue CLI", "@continuedev/cli")
	case "corepack":
		return npmGInstall("Corepack", "corepack")
	case "playwright":
		return npmGInstall("Playwright", "playwright")
	case "cypress":
		return npmGInstall("Cypress", "cypress")
	case "vitest":
		return npmGInstall("Vitest", "vitest")
	case "eslint":
		return npmGInstall("ESLint", "eslint")
	case "locust":
		return pipxInstall("Locust", "locust")
	case "firebase":
		return npmGInstall("Firebase CLI", "firebase-tools")
	case "prism":
		return npmGInstall("Prism", "@stoplight/prism-cli")
	case "eleventy":
		return npmGInstall("Eleventy", "@11ty/eleventy")
	case "spectral":
		return npmGInstall("Spectral", "@stoplight/spectral-cli")
	case "pytest":
		return pytestInstall()
	case "pipx-g":
		return pipxInstall(item.Name, item.Formula)
	}
	return fmt.Errorf("unknown special: %s", item.Special)
}
