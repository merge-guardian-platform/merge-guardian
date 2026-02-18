# Merge Guardian AI

## Project Name: "Merge Guardian AI"
_A GitHub Workflow + AI Integration for Intelligent Merge Management_

## Executive Summary
Merge Guardian AI is an enterprise solution that combines GitHub's native merge queue capabilities with custom AI-powered conflict prediction and resolution. It transforms the "4 PM Friday merge hell" into a smooth, automated, and predictable deployment pipeline. This CLI tool provides the core intelligence for this system, integrating with various AI providers to offer predictive insights and management recommendations.

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
    *   Support for multiple AI providers (OpenAI, Google Gemini)
3.  **Integration Points**
    *   Slack/Teams notifications (future)
    *   Jira ticket linking (future)
    *   Metrics dashboard (future)

This repository contains the Go-based CLI tool that serves as the "Custom AI Layer" and integrates with GitHub Actions.

## Setup

### Prerequisites
*   Go (version 1.22 or higher)
*   GitHub Personal Access Token (with `repo` and `pull_requests` scopes)
*   API Key for your chosen AI Provider (OpenAI or Google Gemini)

### Installation
Clone the repository:
```bash
git clone https://github.com/iamvirul/merge-guardian.git
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
merge-guardian analyze pr --owner <github-owner> --repo <github-repo> --pr-number <pull-request-number> --github-token <your-github-token> --ai-provider <openai|gemini> --ai-api-key <your-ai-api-key>
```

**Flags:**
*   `--owner`, `-o`: GitHub repository owner (e.g., `octocat`)
*   `--repo`, `-r`: GitHub repository name (e.g., `hello-world`)
*   `--pr-number`, `-p`: Pull Request number (e.g., `123`)
*   `--github-token`, `-g`: GitHub Personal Access Token
*   `--ai-provider`, `-i`: AI Service Provider (`openai` or `gemini`, default `openai`)
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

## Code Owners
All code in this repository is currently owned by @iamvirul.
```
* @iamvirul
```
