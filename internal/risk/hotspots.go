package risk

import "strings"

var KnownHotspots = []string{
	"package.json",
	"go.mod",
	"go.sum",
	"README.md",
	"config/",
	"migrations/",
}

// IsHotspot checks if a file is a known hotspot.
func IsHotspot(filename string) bool {
	for _, h := range KnownHotspots {
		if strings.Contains(filename, h) {
			return true
		}
	}
	return false
}

// CountHotspots returns the number of hotspot files in a list.
func CountHotspots(files []string) int {
	count := 0
	for _, f := range files {
		if IsHotspot(f) {
			count++
		}
	}
	return count
}

// IsCrossDirectory checks if changes span multiple top-level directories.
func IsCrossDirectory(files []string) bool {
	dirs := make(map[string]bool)
	for _, f := range files {
		parts := strings.Split(f, "/")
		if len(parts) > 1 {
			dirs[parts[0]] = true
		}
	}
	return len(dirs) > 1
}
