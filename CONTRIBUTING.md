# Contributing to Merge Guardian AI

Thank you for your interest in contributing to Merge Guardian AI! We welcome contributions from the community to help make merge management smarter and more predictable.

## Project Philosophy
Merge Guardian AI aims to predict and resolve merge conflicts using a combination of deterministic rules and AI-powered semantic analysis. Our goal is to ensure that the main branch remains stable and that conflicts are identified before they cause issues in the deployment pipeline.

## Getting Started

### Prerequisites
- **Go**: Version 1.25 or higher.
- **GitHub Account**: To create issues and pull requests.
- **Personal Access Token (PAT)**: With `repo` and `pull_requests` scopes for testing.
- **AI API Keys**: Optional, but needed for testing AI-specific logic (OpenAI, Gemini, or Anthropic).

### Development Setup
1. Fork the repository on GitHub.
2. Clone your fork locally:
   ```bash
   git clone https://github.com/merge-guardian-platform/merge-guardian.git
   cd merge-guardian
   ```
3. Initialize dependencies:
   ```bash
   go mod download
   ```
4. Run tests to ensure everything is working correctly:
   ```bash
   go test ./...
   ```

## Contribution Workflow

### 1. Opening an Issue
Before making a major change, please open an issue to discuss your proposal. This helps avoid redundant work and ensures alignment with the project's roadmap.

### 2. Branching Strategy
- Create a new branch for your feature or bug fix:
  ```bash
  git checkout -b feature/your-feature-name
  ```
- Use descriptive branch names (e.g., `fix/conflict-parsing`, `feat/new-ai-provider`).

### 3. Making Changes
- Follow Go coding conventions and best practices.
- Ensure your code is properly formatted:
  ```bash
  go fmt ./...
  ```
- Add unit tests for new functionality in the corresponding `_test.go` files.

### 4. Submitting a Pull Request
- Push your branch to your fork:
  ```bash
  git push origin feature/your-feature-name
  ```
- Open a Pull Request (PR) against the `main` branch of the official repository.
- Provide a clear and concise description of your changes in the PR body.
- Link any related issues using keywords like `Closes #123`.

## Coding Guidelines

### Testing
- We aim for high test coverage, especially for core risk scoring and AI integration logic.
- Mock external dependencies (like GitHub and AI APIs) where possible to keep tests fast and deterministic.
- Use `httptest` for mocking API responses.

### Commit Messages
- Use clear, professional commit messages.
- Follow the [Conventional Commits](https://www.conventionalcommits.org/) format where appropriate (e.g., `feat:`, `fix:`, `docs:`, `test:`).

## Community and Conduct
We are committed to fostering a welcoming and inclusive community. Please be respectful and professional in all interactions.

## License
By contributing to Merge Guardian AI, you agree that your contributions will be licensed under the [Apache License 2.0](LICENSE).
