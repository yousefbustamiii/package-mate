package components

var catalogPackageManagers = Section{
	Name: "Package Managers",
	Items: []InstallItem{
		{
			Name:    "npm",
			Desc:    "Standard package manager for the JavaScript and Node.js ecosystems",
			Special: "npm-check",
			Color:   rgb(255, 100, 90),
			Binary:  "npm",
		},
		{
			Name:    "Yarn",
			Desc:    "Fast and deterministic package manager for JavaScript dependencies",
			Formula: "yarn",
			Color:   rgb(100, 180, 220),
			Binary:  "yarn",
		},
		{
			Name:    "pnpm",
			Desc:    "Efficient alternative to npm using hard links for storage savings",
			Formula: "pnpm",
			Color:   rgb(255, 180, 60),
			Binary:  "pnpm",
		},
		{
			Name:    "uv",
			Desc:    "Extremely fast Python package installer and resolver written in Rust",
			Formula: "uv",
			Color:   rgb(100, 255, 200),
			Binary:  "uv",
		},
		{
			Name:    "Poetry",
			Desc:    "Dependency management and packaging tool for Python applications",
			Formula: "poetry",
			Color:   rgb(100, 140, 255),
			Binary:  "poetry",
		},
		{
			Name:    "pipx",
			Desc:    "Tool for installing and running Python applications in isolation",
			Formula: "pipx",
			Color:   rgb(255, 220, 120),
			Binary:  "pipx",
		},
		{
			Name:    "Composer",
			Desc:    "Dependency manager for PHP projects and library synchronization",
			Formula: "composer",
			Color:   rgb(200, 150, 100),
			Binary:  "composer",
		},
		{
			Name:    "Corepack",
			Desc:    "Zero-dependency bridge for managing Yarn and pnpm binary versions",
			Special: "npm-g",
			Formula: "corepack",
			Color:   rgb(150, 150, 150),
			Binary:  "corepack",
		},
	},
}
