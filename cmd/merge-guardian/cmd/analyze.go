package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"merge-guardian/internal/ai"
	"merge-guardian/internal/github"
	"merge-guardian/internal/risk"

	ghlib "github.com/google/go-github/v58/github" // Alias for external github package
	"github.com/spf13/cobra"
)

// GitHubClient is an interface for GitHub API interactions, allowing for mocking in tests.
type GitHubClient interface {
	GetPullRequest(owner, repo string, prNumber int) (*ghlib.PullRequest, error)
	GetPullRequestFiles(owner, repo string, prNumber int) ([]*ghlib.CommitFile, error)
	GetRecentMergedPRs(owner, repo, targetBranch string) ([]*ghlib.PullRequest, error)
	GetOpenPRs(owner, repo, targetBranch string) ([]*ghlib.PullRequest, error)
}

// AIClient is an interface for AI service interactions, allowing for mocking in tests.
type AIClient interface {
	GetConflictPrediction(prompt string) (string, error)
}

// realGitHubClient implements the GitHubClient interface using the actual GitHub client.
type realGitHubClient struct {
	*github.Client
}

// realAIClient implements the AIClient interface using the actual AI client.
type realAIClient struct {
	ai.AIClient
}

var (
	githubClient GitHubClient
	aiClient     AIClient
)

// logFatalf is a package-level function variable that can be overridden for testing.
// Kept for backward compatibility if needed, but we will prefer returning errors.
var logFatalf = log.Fatalf

// NewAnalyzeCmd creates and returns a new Cobra command for analysis.
func NewAnalyzeCmd() *cobra.Command {
	analyzeCmd := &cobra.Command{
		Use:   "analyze",
		Short: "Analyze GitHub resources with AI",
		Long:  `Analyze GitHub pull requests and other resources using AI for predictive insights.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Please specify a subcommand: pr")
		},
	}
	analyzeCmd.AddCommand(NewPrCmd())
	return analyzeCmd
}

// NewPrCmd creates and returns a new Cobra subcommand for PR analysis.
func NewPrCmd() *cobra.Command {
	prCmd := &cobra.Command{
		Use:   "pr",
		Short: "Analyze a Pull Request for potential conflicts",
		Long:  `Analyzes a given Pull Request using AI to predict potential merge conflicts and suggest strategies.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Get flag values from the command
			githubOwner, _ := cmd.Flags().GetString("owner")
			githubRepo, _ := cmd.Flags().GetString("repo")
			prNumber, _ := cmd.Flags().GetInt("pr-number")
			githubToken, _ := cmd.Flags().GetString("github-token")
			aiAPIKey, _ := cmd.Flags().GetString("ai-api-key")
			aiProviderStr, _ := cmd.Flags().GetString("ai-provider")
			aiProvider := ai.Provider(aiProviderStr)

			log.Printf("Starting analysis for PR #%d in %s/%s using provider %s", prNumber, githubOwner, githubRepo, aiProvider)

			// Initialize real GitHub client if not already injected (e.g., for testing)
			if githubClient == nil {
				githubClient = &realGitHubClient{github.NewClient(githubToken)}
			}

			// Initialize real AI client if not already injected (e.g., for testing)
			if aiClient == nil {
				var err error
				aiClient, err = ai.NewAIClient(aiProvider, aiAPIKey)
				if err != nil {
					return fmt.Errorf("error creating AI client: %v", err)
				}
			}

			// 1. Get Current PR Details
			currentPR, err := githubClient.GetPullRequest(githubOwner, githubRepo, prNumber)
			if err != nil {
				return fmt.Errorf("error getting current PR details: %v", err)
			}
			if currentPR.Base == nil || currentPR.Base.Ref == nil {
				return fmt.Errorf("could not determine target branch for PR #%d: currentPR.Base is %v, currentPR.Base.Ref is %v", prNumber, currentPR.Base, currentPR.Base.Ref)
			}
			targetBranch := *currentPR.Base.Ref

			// 2. Get Files Changed in Current PR
			prFiles, err := githubClient.GetPullRequestFiles(githubOwner, githubRepo, prNumber)
			if err != nil {
				return fmt.Errorf("error getting PR files: %v", err)
			}
			var changedFiles []string
			for _, file := range prFiles {
				if file.Filename != nil {
					changedFiles = append(changedFiles, *file.Filename)
				}
			}

			// 3. Get Recently Merged PRs
			recentMerges, err := githubClient.GetRecentMergedPRs(githubOwner, githubRepo, targetBranch)
			if err != nil {
				return fmt.Errorf("error getting recent merged PRs: %v", err)
			}
			recentMergesJSON, _ := json.MarshalIndent(simplifyPRs(recentMerges), "", "  ")

			// 4. Get Currently Open PRs
			openPRs, err := githubClient.GetOpenPRs(githubOwner, githubRepo, targetBranch)
			if err != nil {
				return fmt.Errorf("error getting open PRs: %v", err)
			}
			openPRsJSON, _ := json.MarshalIndent(simplifyPRs(openPRs), "", "  ")

			// 5. Deterministic Risk Scoring (Level 2 Intelligence)
			filesChangedCount := len(changedFiles)
			hotspotFilesCount := risk.CountHotspots(changedFiles)
			isCrossDirectory := risk.IsCrossDirectory(changedFiles)
			deterministicRiskScore := risk.CalculateRiskScore(filesChangedCount, hotspotFilesCount, isCrossDirectory)

			// Construct AI Prompt
			prTitle := ""
			if currentPR.Title != nil {
				prTitle = *currentPR.Title
			}

			prompt := buildPredictiveConflictAnalyzerPrompt(
				prNumber,
				prTitle,
				targetBranch,
				changedFiles,
				string(recentMergesJSON),
				string(openPRsJSON),
				deterministicRiskScore, // Pass the calculated score
			)

			// Use the injected AI client
			// Get AI Conflict Prediction
			prediction, err := aiClient.GetConflictPrediction(prompt)
			if err != nil {
				return fmt.Errorf("error getting AI conflict prediction: %v", err)
			}

			// 6. Parse the AI response
			var analysisResp MergeAnalysisResponse
			jsonContent := extractJSON(prediction)

			if jsonContent != "" {
				if err := json.Unmarshal([]byte(jsonContent), &analysisResp); err != nil {
					// Fallback to raw output if parsing fails, but warn the user
					log.Printf("Warning: Failed to parse AI JSON response: %v", err)
					fmt.Println("\n--- AI Conflict Prediction Result (Raw) ---")
					fmt.Println(prediction)
					return nil
				}
			} else {
				log.Printf("Warning: No JSON found in AI response")
				fmt.Println("\n--- AI Conflict Prediction Result (Raw) ---")
				fmt.Println(prediction)
				return nil
			}

			// Injected fields (Real Intelligence)
			analysisResp.MergeAnalysis.AnalysisTimestampParsed = time.Now().UTC()
			analysisResp.MergeAnalysis.AnalysisTimestamp = analysisResp.MergeAnalysis.AnalysisTimestampParsed.Format(time.RFC3339)

			// Pretty print the structured result
			outputJSON, err := json.MarshalIndent(analysisResp, "", "  ")
			if err != nil {
				return fmt.Errorf("error marshalling output: %v", err)
			}

			fmt.Println(string(outputJSON))
			return nil
		},
	}

	prCmd.Flags().StringP("owner", "o", "", "GitHub repository owner")
	prCmd.Flags().StringP("repo", "r", "", "GitHub repository name")
	prCmd.Flags().IntP("pr-number", "p", 0, "Pull Request number")
	prCmd.Flags().StringP("github-token", "g", "", "GitHub Personal Access Token")
	prCmd.Flags().StringP("ai-api-key", "k", "", "AI Service API Key (OpenAI, Gemini, Anthropic)")
	prCmd.Flags().StringP("ai-provider", "i", string(ai.ProviderOpenAI), "AI Service Provider (openai, gemini, anthropic)")

	prCmd.MarkFlagRequired("owner")
	prCmd.MarkFlagRequired("repo")
	prCmd.MarkFlagRequired("pr-number")
	prCmd.MarkFlagRequired("github-token")
	prCmd.MarkFlagRequired("ai-api-key")
	return prCmd
}

// simplifyPRs extracts essential information from GitHub PR objects for the AI prompt.
func simplifyPRs(prs []*ghlib.PullRequest) []map[string]interface{} {
	var simplified []map[string]interface{}
	for _, pr := range prs {
		if pr.Number != nil && pr.Title != nil {
			simplified = append(simplified, map[string]interface{}{
				"number": *pr.Number,
				"title":  *pr.Title,
				// Optionally add more fields if the AI needs them
			})
		}
		if pr.Base != nil && pr.Base.Ref != nil { // Add this check for safety
			simplified[len(simplified)-1]["base_ref"] = *pr.Base.Ref
		}
	}
	return simplified
}

// buildPredictiveConflictAnalyzerPrompt constructs the AI prompt based on the user's template.
func buildPredictiveConflictAnalyzerPrompt(
	prNum int,
	prTitle string,
	targetBranch string,
	changedFiles []string,
	recentMergesJSON string,
	openPRsJSON string,
	deterministicRiskScore int,
) string {
	filesList := "-\n"
	if len(changedFiles) > 0 {
		filesList = "- " + strings.Join(changedFiles, "\n- ")
	}

	promptTemplate := `You are Merge Guardian AI, an expert in analyzing code conflicts and predicting merge issues.

Given the following context:
- Current PR: %[1]d "%[2]s" targeting branch: %[3]s
- Calculated Pre-Analysis Risk Score: %[7]d/100 (Use this as a baseline, but adjust if semantic analysis reveals deeper issues)
- Files changed:
%[4]s
- Recently merged PRs (last 24 hours):
%[5]s
- Currently open PRs targeting same branch:
%[6]s

Task:
1. **Analyze Future Conflict Risk**: Compare files changed in this PR against recently merged and currently open PRs.
2. **Identify Severity**:
   - 🔴 HIGH: Same files, overlapping line ranges, or complex refactors.
   - 🟡 MEDIUM: Same files, different sections but related logic.
   - 🟢 LOW: Different files, minimal risk.
3. **Semantic Conflict Analysis**: Detect logical conflicts (e.g., function signature changes, dependency updates) that might check out fine in git but break runtime.
4. **Detect Risky Refactors**: Flag large-scale renames or structural changes across many files.
5. **Score Merge Risk**: Start with the baseline score (%[7]d). Explain why you increased or decreased it based on your semantic analysis.
6. **Identify Hotspots**: Highlight files that are being touched by multiple PRs or have a history of conflict (inferred from context).
7. **Suggest Strategy**: Recommend optimal merge order.

Output Format (JSON):
{
  "merge_analysis": {
    "pr_number": %[1]d,
    "risk_score": 0-100,
    "risk_level": "HIGH|MEDIUM|LOW",
    "potential_conflicts": [
      {
        "file": "path/to/file",
        "severity": "HIGH|MEDIUM|LOW",
        "type": "DIRECT|SEMANTIC|REFACTOR",
        "description": "Explanation..."
      }
    ],
    "hotspots": ["file1", "file2"],
    "merge_strategy_recommendation": "...",
    "human_readable_summary": "..."
  }
}
`
	return fmt.Sprintf(
		promptTemplate,
		prNum,
		prTitle, // 2
		targetBranch, // 3
		filesList, // 4
		recentMergesJSON, // 5
		openPRsJSON, // 6
		deterministicRiskScore, // 7
	)
}

// extractJSON extracts the first valid JSON object from a string by matching braces.
func extractJSON(s string) string {
	start := strings.Index(s, "{")
	if start == -1 {
		return ""
	}

	depth := 0
	for i := start; i < len(s); i++ {
		if s[i] == '{' {
			depth++
		} else if s[i] == '}' {
			depth--
			if depth == 0 {
				return s[start : i+1]
			}
		}
	}

	return ""
}