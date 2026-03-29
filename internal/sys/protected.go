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

// IsProtectedPath returns true if the given path is a vital macOS system file.
func IsProtectedPath(path string) bool {
	for _, p := range protectedPrefixes {
		if strings.HasPrefix(path, p) {
			return true
		}
	}
	return false
}
