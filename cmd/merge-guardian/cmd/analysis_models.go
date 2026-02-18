package cmd

import "time"

type MergeAnalysisResponse struct {
	MergeAnalysis MergeAnalysis `json:"merge_analysis"`
}

type MergeAnalysis struct {
	PRNumber                int                 `json:"pr_number"`
	RiskScore               int                 `json:"risk_score"`
	RiskLevel               string              `json:"risk_level"`
	PotentialConflicts      []PotentialConflict `json:"potential_conflicts"`
	Hotspots                []string            `json:"hotspots"`
	MergeStrategyRec        string              `json:"merge_strategy_recommendation"`
	HumanReadableSummary    string              `json:"human_readable_summary"`
	AnalysisTimestamp       string              `json:"analysis_timestamp,omitempty"` // Injected field
	AnalysisTimestampParsed time.Time           `json:"-"`
}

type PotentialConflict struct {
	File        string `json:"file"`
	Severity    string `json:"severity"`
	Type        string `json:"type"`
	Description string `json:"description"`
}
