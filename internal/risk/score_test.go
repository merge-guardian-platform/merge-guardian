package risk

import "testing"

func TestCalculateRiskScore(t *testing.T) {
	tests := []struct {
		name             string
		filesChanged     int
		hotspotFiles     int
		isCrossDir       bool
		hasMergeConflict bool
		minScore         int
		maxScore         int
	}{
		{"Low Risk", 2, 0, false, false, 1, 5},
		{"Medium Risk", 10, 1, false, false, 15, 20},
		{"High Risk", 20, 2, true, false, 35, 45},
		{"Max Risk", 200, 10, true, true, 100, 100},
		{"Merge Conflict Risk", 2, 0, false, true, 15, 20},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score := CalculateRiskScore(tt.filesChanged, tt.hotspotFiles, tt.isCrossDir, tt.hasMergeConflict)
			if score < tt.minScore || score > tt.maxScore {
				t.Errorf("CalculateRiskScore() = %v, want between %v and %v", score, tt.minScore, tt.maxScore)
			}
		})
	}
}

func TestIsHotspot(t *testing.T) {
	tests := []struct {
		file     string
		expected bool
	}{
		{"package.json", true},
		{"src/utils/helper.ts", false},
		{"config/db.js", true},
		{"src/README.md", true},
	}
	for _, tt := range tests {
		t.Run(tt.file, func(t *testing.T) {
			if got := IsHotspot(tt.file); got != tt.expected {
				t.Errorf("IsHotspot(%q) = %v, want %v", tt.file, got, tt.expected)
			}
		})
	}
}
