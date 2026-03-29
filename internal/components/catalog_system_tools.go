package components

var catalogSystemTools = Section{
	Name: "System Tools",
	Items: []InstallItem{
		{
			Name:    "htop",
			Desc:    "Interactive process viewer and system monitor for Unix-based systems",
			Formula: "htop",
			Color:   rgb(100, 255, 100),
			Binary:  "htop",
		},
		{
			Name:    "btop",
			Desc:    "Resource monitor showing detailed usage and system statistics",
			Formula: "btop",
			Color:   rgb(50, 200, 255),
			Binary:  "btop",
		},
		{
			Name:    "NeoVim",
			Desc:    "Highly extensible terminal-based text editor built on Vim core",
			Formula: "neovim",
			Color:   rgb(100, 230, 50),
			Binary:  "nvim",
		},
		{
			Name:    "Vim",
			Desc:    "Classic terminal editor for efficient text and code manipulation",
			Formula: "vim",
			Color:   rgb(50, 180, 50),
			Binary:  "vim",
		},
		{
			Name:    "jq",
			Desc:    "Command-line processor for parsing and transforming JSON data",
			Formula: "jq",
			Color:   rgb(127, 219, 255),
			Binary:  "jq",
		},
		{
			Name:    "yq",
			Desc:    "Lightweight command-line processor for YAML and XML document formats",
			Formula: "yq",
			Color:   rgb(255, 150, 50),
			Binary:  "yq",
		},
		{
			Name:    "fd",
			Desc:    "Simple and fast alternative to the find command for file discovery",
			Formula: "fd",
			Color:   rgb(255, 100, 150),
			Binary:  "fd",
		},
		{
			Name:    "ncdu",
			Desc:    "Disk usage analyzer with an interactive ncurses-based interface",
			Formula: "ncdu",
			Color:   rgb(150, 150, 255),
			Binary:  "ncdu",
		},
		{
			Name:    "Tree",
			Desc:    "Recursive directory listing program that produces a depth-indented view",
			Formula: "tree",
			Color:   rgb(200, 150, 255),
			Binary:  "tree",
		},
		{
			Name:    "entr",
			Desc:    "Utility for running arbitrary commands when specific files change",
			Formula: "entr",
			Color:   rgb(255, 210, 50),
			Binary:  "entr",
		},
	},
}
