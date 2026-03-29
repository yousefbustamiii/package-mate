package components

var catalogTesting = Section{
	Name: "Testing & Utilities",
	Items: []InstallItem{
		{
			Name:    "Playwright",
			Desc:    "Fast and reliable end-to-end testing tool for modern web apps",
			Special: "npm-g",
			Formula: "playwright",
			Color:   rgb(100, 220, 100),
			Binary:  "playwright",
		},
		{
			Name:    "Cypress",
			Desc:    "Testing framework for anything that runs in a web browser environment",
			Special: "npm-g",
			Formula: "cypress",
			Color:   rgb(150, 180, 200),
			Binary:  "cypress",
		},
		{
			Name:    "Vitest",
			Desc:    "Vite-native unit test framework for high-performance testing",
			Special: "npm-g",
			Formula: "vitest",
			Color:   rgb(230, 210, 50),
			Binary:  "vitest",
		},
		{
			Name:    "Pre-commit",
			Desc:    "Framework for managing and maintaining multi-language Git hooks",
			Formula: "pre-commit",
			Color:   rgb(255, 150, 100),
			Binary:  "pre-commit",
		},
		{
			Name:    "tox",
			Desc:    "Command-line tool for managing and testing Python environments",
			Formula: "tox",
			Color:   rgb(50, 180, 255),
			Binary:  "tox",
		},
		{
			Name:    "Ruff",
			Desc:    "Fast Python linter and formatter designed for high-performance projects",
			Formula: "ruff",
			Color:   rgb(100, 230, 255),
			Binary:  "ruff",
		},
		{
			Name:    "Black",
			Desc:    "Strict Python code formatter for maintaining consistent code style",
			Formula: "black",
			Color:   rgb(100, 100, 100),
			Binary:  "black",
		},
		{
			Name:    "ESLint",
			Desc:    "Static analysis tool for identifying and fixing problems in JavaScript",
			Special: "npm-g",
			Formula: "eslint",
			Color:   rgb(150, 100, 255),
			Binary:  "eslint",
		},
	},
}
