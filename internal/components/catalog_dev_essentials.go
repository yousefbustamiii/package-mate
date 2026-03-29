package components

var catalogDevEssentials = Section{
	Name: "Dev Essentials",
	Items: []InstallItem{
		{
			Name:    "Git",
			Desc:    "Distributed version control system for tracking source code changes",
			Formula: "git",
			Color:   rgb(255, 100, 80),
			Binary:  "git",
		},
		{
			Name:    "Git LFS",
			Desc:    "Extension for managing large files within Git version control",
			Formula: "git-lfs",
			Color:   rgb(255, 140, 100),
			Binary:  "git-lfs",
		},
		{
			Name:    "LazyGit",
			Desc:    "Simple terminal-based user interface for managing Git operations",
			Formula: "lazygit",
			Color:   rgb(100, 230, 180),
			Binary:  "lazygit",
		},
		{
			Name:    "fzf",
			Desc:    "General-purpose command-line fuzzy finder for files and processes",
			Formula: "fzf",
			Color:   rgb(255, 140, 230),
			Binary:  "fzf",
		},
		{
			Name:    "ripgrep",
			Desc:    "Fast line-oriented search tool that respects gitignore rules",
			Formula: "ripgrep",
			Color:   rgb(180, 230, 50),
			Binary:  "rg",
		},
		{
			Name:    "bat",
			Desc:    "Modern cat replacement with syntax highlighting and Git integration",
			Formula: "bat",
			Color:   rgb(100, 230, 255),
			Binary:  "bat",
		},
		{
			Name:    "eza",
			Desc:    "Modern replacement for the ls command with enhanced file metadata",
			Formula: "eza",
			Color:   rgb(180, 220, 255),
			Binary:  "eza",
		},
		{
			Name:    "zoxide",
			Desc:    "Smarter navigation tool inspired by the cd command and autojump",
			Formula: "zoxide",
			Color:   rgb(255, 220, 100),
			Binary:  "zoxide",
		},
		{
			Name:    "Starship",
			Desc:    "Minimal and customizable shell prompt for any terminal environment",
			Formula: "starship",
			Color:   rgb(255, 100, 255),
			Binary:  "starship",
		},
		{
			Name:    "direnv",
			Desc:    "Shell extension to load or unload environment variables per directory",
			Formula: "direnv",
			Color:   rgb(255, 255, 100),
			Binary:  "direnv",
		},
		{
			Name:    "NVM",
			Desc:    "Node Version Manager for switching between multiple Node.js environments",
			Special: "nvm",
			Color:   rgb(140, 200, 75),
			Binary:  "nvm",
		},
	},
}
