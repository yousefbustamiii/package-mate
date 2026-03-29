package components

var catalogDBGUI = Section{
	Name: "DB GUI & Dev Tools",
	Items: []InstallItem{
		{
			Name:   "TablePlus",
			Desc:   "Native database management client for macOS with a clean interface",
			Cask:   "tableplus",
			Color:  rgb(120, 220, 255),
			Binary: "TablePlus",
		},
		{
			Name:   "DBeaver",
			Desc:   "Universal database tool for developers supporting multiple SQL dialects",
			Cask:   "dbeaver-community",
			Color:  rgb(160, 190, 210),
			Binary: "dbeaver",
		},
		{
			Name:   "RedisInsight",
			Desc:   "Visual tool for managing and interacting with Redis data structures",
			Cask:   "redisinsight",
			Color:  rgb(255, 100, 100),
			Binary: "RedisInsight-v2",
		},
		{
			Name:   "Postman",
			Desc:   "API platform for designing, building, and testing RESTful services",
			Cask:   "postman",
			Color:  rgb(255, 150, 100),
			Binary: "Postman",
		},
		{
			Name:   "Insomnia",
			Desc:   "Open-source application for API design and collaborative testing",
			Cask:   "insomnia",
			Color:  rgb(150, 100, 255),
			Binary: "Insomnia",
		},
		{
			Name:    "HTTPie",
			Desc:    "User-friendly command-line interface for making HTTP requests",
			Formula: "httpie",
			Color:   rgb(30, 200, 255),
			Binary:  "http",
		},
		{
			Name:   "Lens",
			Desc:   "Integrated development environment for Kubernetes cluster management",
			Cask:   "lens",
			Color:  rgb(60, 180, 255),
			Binary: "Lens",
		},
		{
			Name:   "Sequel Ace",
			Desc:   "Lightweight MySQL and MariaDB database management tool for macOS",
			Cask:   "sequel-ace",
			Color:  rgb(100, 160, 255),
			Binary: "Sequel Ace",
		},
	},
}
