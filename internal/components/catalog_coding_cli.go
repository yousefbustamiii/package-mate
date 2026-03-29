package components

var catalogCodingCLI = Section{
	Name: "Coding CLIs & AI",
	Items: []InstallItem{
		{
			Name:    "Claude Code",
			Desc:    "Agentic coding CLI from Anthropic for terminal-based development",
			Special: "claude",
			Color:   rgb(230, 150, 120),
			Binary:  "claude",
		},
		{
			Name:    "Gemini CLI",
			Desc:    "Developer interface for interacting with Google's Gemini AI models",
			Special: "gemini",
			Color:   rgb(100, 170, 255),
			Binary:  "gemini",
		},
		{
			Name:    "Aider",
			Desc:    "AI-powered pair programming tool for terminal-oriented coding tasks",
			Special: "pipx-g",
			Formula: "aider-chat",
			Color:   rgb(255, 150, 50),
			Binary:  "aider",
		},
		{
			Name:    "Continue CLI",
			Desc:    "Command-line companion for the Continue AI code assistant platform",
			Special: "npm-g",
			Formula: "@continuedev/cli",
			Color:   rgb(100, 220, 255),
			Binary:  "typegen",
		},
		{
			Name:    "GitHub CLI",
			Desc:    "Official interface for interacting with GitHub repositories and PRs",
			Formula: "gh",
			Color:   rgb(255, 255, 255),
			Binary:  "gh",
		},
	},
}
