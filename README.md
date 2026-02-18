# Merge Guardian AI

## Project Name: "Merge Guardian AI"
_A GitHub Workflow + AI Integration for Intelligent Merge Management_

## Executive Summary
Merge Guardian AI is an enterprise solution that combines GitHub's native merge queue capabilities with custom AI-powered conflict prediction and resolution. It transforms the "4 PM Friday merge hell" into a smooth, automated, and predictable deployment pipeline. This CLI tool provides the core intelligence for this system, integrating with various AI providers to offer predictive insights and management recommendations.

---

## Why Merge Guardian? (vs GitHub Native)

GitHub is excellent at **reporting the present**. Merge Guardian is designed to **predict the future**.

| Feature | GitHub Native | Merge Guardian AI |
| :--- | :--- | :--- |
| **Conflict Detection** | Reactive (Only detects current git conflicts) | **Predictive** (Analyzes future risk based on open/recent PRs) |
| **Semantic Analysis** | No (Text-based diffs only) | **Yes** (Detects logic breaks, signature changes, API shifts) |
| **Risk Scoring** | None | **Deterministic + AI Score** (0-100% Risk Probability) |
| **Refactor Safety** | "Looks fine to me" | **Risky Refactor Detection** (Flags mass renames/moves) |
| **Hotspot Detection** | None | **Hotspot Tracking** (Identifies high-churn/fragile files) |
| **Merge Strategy** | Primitive (First-in-first-out) | **Intelligent** (Suggests optimal merge order to minimize breaks) |

---

## Intelligence Levels

Merge Guardian operates on three levels of intelligence to protect your main branch:

### Level 1: Reporting (The Basics)
- Aggregates context from Open PRs, Recently Merged PRs, and Current Changes.
- Provides a single view of all activity targeting the branch.

### Level 2: Deterministic Risk Engine (The Rules)
- **Hotspot Detection**: Flags changes to sensitive files (e.g., `package.json`, `migrations/`, `go.mod`).
- **Cross-Directory Impact**: Calculates risk when changes span multiple architectural domains.
- **Volume Analysis**: penalizes massive file changes that are hard to review.

### Level 3: AI Predictive Engine (The Brain)
- **Semantic Conflict Analysis**: "PR A changed the function signature, PR B is still calling the old one. GitHub says it merges, but it will break prod." -> **Merge Guardian catches this.**
- **Merge Strategy**: Recommends whether to merge immediately, wait for another PR, or reorder the queue.

---

## Architecture Overview

### Components:
1.  **GitHub Native Tools**
    *   Merge Queue (GitHub Actions)
    *   Required status checks
    *   Branch protection rules
2.  **Custom AI Layer**
    *   Predictive conflict analysis
    *   Intelligent merge strategy recommendation
    *   Automated conflict resolution suggestions
    *   Support for multiple AI providers (OpenAI, Google Gemini, Anthropic)
3.  **Integration Points**
    *   Slack/Teams notifications (future)
    *   Jira ticket linking (future)
    *   Metrics dashboard (future)

This repository contains the Go-based CLI tool that serves as the "Custom AI Layer" and integrates with GitHub Actions.

## Setup

### Prerequisites
*   Go (version 1.22 or higher)
*   GitHub Personal Access Token (with `repo` and `pull_requests` scopes)
*   API Key for your chosen AI Provider (OpenAI, Google Gemini, or Anthropic)

### Installation
Clone the repository:
```bash
git clone https://github.com/merge-guardian-platform/merge-guardian.git
cd merge-guardian
```

Build the CLI tool:
```bash
go build -o merge-guardian ./cmd/merge-guardian
```

Move the executable to your PATH (optional, but recommended):
```bash
mv merge-guardian /usr/local/bin/
```

### Configuration
Environment variables or direct flags can be used for API keys and tokens.
It is highly recommended to use GitHub Actions secrets for sensitive information.

## Usage

The `merge-guardian` CLI is designed to be integrated into GitHub Actions workflows.
The primary command currently implemented is `analyze pr`.

### `analyze pr` Command
This command analyzes a Pull Request and uses AI to predict potential merge conflicts.

```bash
merge-guardian analyze pr \
  --owner <github-owner> \
  --repo <github-repo> \
  --pr-number <pull-request-number> \
  --github-token <your-github-token> \
  --ai-provider <openai|gemini|anthropic> \
  --ai-api-key <your-ai-api-key>
```

**Flags:**
*   `--owner`, `-o`: GitHub repository owner (e.g., `octocat`)
*   `--repo`, `-r`: GitHub repository name (e.g., `hello-world`)
*   `--pr-number`, `-p`: Pull Request number (e.g., `123`)
*   `--github-token`, `-g`: GitHub Personal Access Token
*   `--ai-provider`, `-i`: AI Service Provider (`openai`, `gemini`, or `anthropic`, default `openai`)
*   `--ai-api-key`, `-k`: API Key for the chosen AI Service Provider

**Example in GitHub Actions (for OpenAI):**
```yaml
- name: Run AI Conflict Prediction
  run: |
    ./merge-guardian analyze pr 
      --owner ${{ github.repository_owner }} 
      --repo ${{ github.event.repository.name }} 
      --pr-number ${{ github.event.pull_request.number }} 
      --github-token ${{ secrets.GITHUB_TOKEN }} 
      --ai-provider openai 
      --ai-api-key ${{ secrets.OPENAI_API_KEY }}
```

## Future Roadmap
The current implementation focuses on the "Predictive Conflict Analyzer" (Prompt 1). Future phases will include:
*   Intelligent Merge Strategy Recommendation
*   Automated Conflict Resolution Assistant
*   Post-Merge Health Analysis
*   Integration with communication platforms (Slack, Teams) and issue trackers (Jira).
*   Fine-tuning AI models and further optimization.

## Contributing
Contributions are welcome! Please refer to the `CODEOWNERS` file for ownership details.

## License
This project is licensed under the Apache License 2.0. See the [LICENSE](LICENSE) file for details.

