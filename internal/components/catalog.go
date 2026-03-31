package components

import (
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
)

//go:embed data/*.json
var catalogFS embed.FS

// AllSections is the master ordered catalog of every installable tool.
var AllSections []Section

func init() {
	loadCatalog()
}

func loadCatalog() {
	var orderedSectionNames = []string{
		"Homebrew Setup",
		"Databases",
		"Caching & Messaging",
		"Containers & DevOps",
		"Backend & Runtime",
		"Languages & Runtimes",
		"DB GUI & Dev Tools",
		"Coding CLIs & AI",
		"Dev Essentials",
		"Package Managers",
		"Testing & Utilities",
		"System Tools",
		"Infrastructure & Cloud",
		"Editors & IDEs",
		"Security & Secrets",
		"Media & Graphics",
		"Low-Level & Embedded",
		"Reverse Engineering",
		"Docs & Static Sites",
		"Performance & Profiling",
		"Virtualization",
		"Terminal Glow-Up & Shell",
		"macOS Essentials",
		"Cloud & Kubernetes",
		"Data & Analytics",
		"Modern CLI Replacements",
		"Networking & API",
		"Dev Utilities",
		"Design & Frontend",
		"Web Browsers",
		"Communications",
		"Knowledge & Productivity",
		"Terminal Emulators",
		"JetBrains Mastery",
	}

	files := []string{
		"data/development.json",
		"data/infrastructure.json",
		"data/system.json",
		"data/security_networking.json",
		"data/misc_applications.json",
	}

	sectionsByName := make(map[string]Section)

	for _, fileName := range files {
		data, err := catalogFS.ReadFile(fileName)
		if err != nil {
			panic(fmt.Sprintf("Failed to load catalog file %s: %v", fileName, err))
		}

		var secs []Section
		decoder := json.NewDecoder(bytes.NewReader(data))
		// Important logic: ensure strict decoding if we want, but default is fine.
		if err := decoder.Decode(&secs); err != nil {
			panic(fmt.Sprintf("Failed to parse catalog JSON in %s: %v", fileName, err))
		}

		for _, sec := range secs {
			sectionsByName[sec.Name] = sec
		}
	}

	for _, name := range orderedSectionNames {
		if sec, exists := sectionsByName[name]; exists {
			AllSections = append(AllSections, sec)
		} else {
			// fallback warning in case a name mismatches later, but don't panic for resilience.
			fmt.Printf("Warning: Section '%s' missing from parsed JSON files.\n", name)
		}
	}
}
