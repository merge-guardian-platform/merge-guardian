package risk

// Rule-based risk scoring factors
const (
	BaseRiskScore           = 0
	RiskPerFileChanged      = 0.5
	RiskPerHotspotFile      = 10.0
	RiskCrossDirectory      = 5.0
	RiskRecentMergeConflict = 15.0
)

// CalculateRiskScore computes a deterministic risk score based on PR metadata.
func CalculateRiskScore(filesChangedCount int, hotspotFilesCount int, isCrossDirectory, hasRecentMergeConflict bool) int {
	score := float64(BaseRiskScore)

	score += float64(filesChangedCount) * RiskPerFileChanged
	score += float64(hotspotFilesCount) * RiskPerHotspotFile

	if isCrossDirectory {
		score += RiskCrossDirectory
	}

	if hasRecentMergeConflict {
		score += RiskRecentMergeConflict
	}

	// Cap at 100
	if score > 100 {
		score = 100
	}

	return int(score)
}
