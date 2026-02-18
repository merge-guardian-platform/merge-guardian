package risk

import (
	"path/filepath"
	"strings"
)

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
	cleanFilename := filepath.Clean(filename)

	for _, h := range KnownHotspots {
		// 1. Check if the rule is intended for a directory (ends with /)
		isDirRule := strings.HasSuffix(h, "/")
		cleanRule := filepath.Clean(h)

		if isDirRule {
			// Directory Match:
			// "config/" should match "config/db.js" but NOT "my-config/db.js"
			// The clean filename must start with the clean rule + separator, OR be the directory itself.
			if strings.HasPrefix(cleanFilename, cleanRule+string(filepath.Separator)) || cleanFilename == cleanRule {
				return true
			}
		} else {
			// File Match:
			// If rule has no path separators (e.g. "package.json"), match filename base.
			if !strings.Contains(h, "/") && !strings.Contains(h, "\\") {
				if filepath.Base(cleanFilename) == cleanRule {
					return true
				}
			} else {
				// If rule has path separators (e.g. "src/config.js"), match exact suffix with boundary.
				// "app/src/config.js" matches "src/config.js".
				// "src/config.js" matches "src/config.js".
				if strings.HasSuffix(cleanFilename, cleanRule) {
					if cleanFilename == cleanRule || strings.HasSuffix(cleanFilename, string(filepath.Separator)+cleanRule) {
						return true
					}
				}
			}
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
