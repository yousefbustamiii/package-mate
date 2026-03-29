package components

var catalogHomebrew = Section{
	Name: "Homebrew Setup",
	Items: []InstallItem{
		{
			Name:    "Homebrew",
			Desc:    "Standard package manager for macOS and Linux",
			Special: "brew-install",
			Color:   rgb(255, 170, 50),
			Binary:  "brew",
		},
		{
			Name:    "Homebrew Update",
			Desc:    "Synchronize local formula and cask metadata with upstream",
			Special: "brew-update",
			Color:   rgb(255, 230, 100),
			Binary:  "brew",
		},
	},
}
