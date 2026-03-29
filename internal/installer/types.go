package installer

// Status describes the result of an install operation.
type Status int

const (
	StatusInstalled   Status = iota // Freshly installed
	StatusAlreadyHave               // Was already present — skipped
	StatusFailed                    // Install attempt failed
)

// Result is returned by Install and Ensure* functions.
type Result struct {
	ItemName string
	Status   Status
	Version  string // populated when StatusAlreadyHave
	Err      error
}

type DetectionStatus int

const (
	DetectionNotFound DetectionStatus = iota
	DetectionExact                    // Exact formula/cask version found
	DetectionBinary                   // Generic binary found on path but no exact formula match
	DetectionDifferentBrew            // Found in Homebrew prefix but different formula/version
	DetectionOutdated                 // Exact formula/cask version found, but a newer version is available
	DetectionManualApp               // GUI App found in /Applications but not managed by Brew
	DetectionTrashedApp              // GUI App found in Trash
)

type Detection struct {
	Status      DetectionStatus
	Detail      string // The version or path found
	BinaryPath  string // The actual path of the unmanaged binary
	BrewFormula string // The extracted Homebrew formula name (if applicable)
	IsBrewCask  bool   // Whether the detected brew formula is a Cask
}

type VersionType int

const (
	VersionManaged VersionType = iota
	VersionManagedOlder
	VersionUnmanaged
)

type VersionEntry struct {
	Type    VersionType
	Version string
	Path    string
	Formula string // The specific Homebrew formula name (e.g. postgresql@16)
}
