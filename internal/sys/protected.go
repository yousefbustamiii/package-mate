package sys

import "strings"

var protectedPrefixes = []string{
	"/System",        // Core macOS files (Locked by SIP)
	"/Library/Apple", // Apple-specific system resources
	"/usr/bin",       // Standard OS binaries (python3, ruby, etc.)
	"/usr/sbin",      // System admin binaries
	"/bin",           // Essential binaries (sh, zsh, ls)
	"/sbin",          // Essential system binaries
	"/var/root",      // Root user home
	"/opt/apple",     // Apple optional software
}

var systemPrefixes = []string{
	"/usr/local", // Shared system binaries (often root-owned)
}

// IsProtectedPath returns true if the given path is a vital macOS system file.
func IsProtectedPath(path string) bool {
	for _, p := range protectedPrefixes {
		if strings.HasPrefix(path, p) {
			return true
		}
	}
	return false
}

// IsSystemPath returns true if the path is a system-owned location where
// background operations might require elevations (e.g. /usr/local).
func IsSystemPath(path string) bool {
	if IsProtectedPath(path) {
		return true
	}
	for _, p := range systemPrefixes {
		if strings.HasPrefix(path, p) {
			return true
		}
	}
	return false
}
