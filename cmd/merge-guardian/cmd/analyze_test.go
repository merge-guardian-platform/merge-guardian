package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	ghlib "github.com/google/go-github/v58/github"
)

// MockGitHubClient is a mock implementation of GitHubClient for testing.
type MockGitHubClient struct {
	GetPullRequestFunc      func(owner, repo string, prNumber int) (*ghlib.PullRequest, error)
	GetPullRequestFilesFunc func(owner, repo string, prNumber int) ([]*ghlib.CommitFile, error)
	GetRecentMergedPRsFunc  func(owner, repo, targetBranch string) ([]*ghlib.PullRequest, error)
	GetOpenPRsFunc          func(owner, repo, targetBranch string) ([]*ghlib.PullRequest, error)
}

func (m *MockGitHubClient) GetPullRequest(owner, repo string, prNumber int) (*ghlib.PullRequest, error) {
	if m.GetPullRequestFunc != nil {
		return m.GetPullRequestFunc(owner, repo, prNumber)
	}
	return nil, nil
}
func (m *MockGitHubClient) GetPullRequestFiles(owner, repo string, prNumber int) ([]*ghlib.CommitFile, error) {
	if m.GetPullRequestFilesFunc != nil {
		return m.GetPullRequestFilesFunc(owner, repo, prNumber)
	}
	return nil, nil
}
func (m *MockGitHubClient) GetRecentMergedPRs(owner, repo, targetBranch string) ([]*ghlib.PullRequest, error) {
	if m.GetRecentMergedPRsFunc != nil {
		return m.GetRecentMergedPRsFunc(owner, repo, targetBranch)
	}
	return nil, nil
}
func (m *MockGitHubClient) GetOpenPRs(owner, repo, targetBranch string) ([]*ghlib.PullRequest, error) {
	if m.GetOpenPRsFunc != nil {
		return m.GetOpenPRsFunc(owner, repo, targetBranch)
	}
	return nil, nil
}

// MockAIClient is a mock implementation of AIClient for testing.
type MockAIClient struct {
	GetConflictPredictionFunc func(prompt string) (string, error)
}

func (m *MockAIClient) GetConflictPrediction(prompt string) (string, error) {
	if m.GetConflictPredictionFunc != nil {
		return m.GetConflictPredictionFunc(prompt)
	}
	return "", nil
}

// Helper function to capture stdout
func captureStdout(f func()) string {
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	buf.ReadFrom(r)
	return buf.String()
}

// Reset global variables for clean test runs
func resetGlobalClients() {
	githubClient = nil
	aiClient = nil
}

func TestAnalyzePRCommand_Success(t *testing.T) {
	defer resetGlobalClients() // Ensure cleanup after test

	// Setup mock GitHub client
	mockGh := &MockGitHubClient{
		GetPullRequestFunc: func(owner, repo string, prNum int) (*ghlib.PullRequest, error) {
			if owner == "test-owner" && repo == "test-repo" && prNum == 123 {
				return &ghlib.PullRequest{
					Number: ghlib.Int(123),
					Title:  ghlib.String("Test PR"),
					Base:   &ghlib.PullRequestBranch{Ref: ghlib.String("main")},
				}, nil
			}
			return nil, fmt.Errorf("unexpected GetPullRequest call")
		},
		GetPullRequestFilesFunc: func(owner, repo string, prNum int) ([]*ghlib.CommitFile, error) {
			if owner == "test-owner" && repo == "test-repo" && prNum == 123 {
				return []*ghlib.CommitFile{
					{Filename: ghlib.String("file1.go")},
					{Filename: ghlib.String("file2.txt")},
				}, nil
			}
			return nil, fmt.Errorf("unexpected GetPullRequestFiles call")
		},
		GetRecentMergedPRsFunc: func(owner, repo, targetBranch string) ([]*ghlib.PullRequest, error) {
			if owner == "test-owner" && repo == "test-repo" && targetBranch == "main" {
				return []*ghlib.PullRequest{
					{Number: ghlib.Int(10), Title: ghlib.String("Recent Merge")},
				}, nil
			}
			return nil, fmt.Errorf("unexpected GetRecentMergedPRs call")
		},
		GetOpenPRsFunc: func(owner, repo, targetBranch string) ([]*ghlib.PullRequest, error) {
			if owner == "test-owner" && repo == "test-repo" && targetBranch == "main" {
				return []*ghlib.PullRequest{
					{Number: ghlib.Int(20), Title: ghlib.String("Open PR")},
				}, nil
			}
			return nil, fmt.Errorf("unexpected GetOpenPRs call")
		},
	}
	githubClient = mockGh

	// Setup mock AI client
	mockAI := &MockAIClient{
		GetConflictPredictionFunc: func(prompt string) (string, error) {
			if strings.Contains(prompt, "Test PR") && strings.Contains(prompt, "file1.go") {
				return "AI Prediction Result", nil
			}
			return "", fmt.Errorf("unexpected GetConflictPrediction call with prompt: %s", prompt)
		},
	}
	aiClient = mockAI

	// Create a new prCmd instance for this test
	cmd := NewPrCmd()

	// Set flags
	cmd.Flags().Set("owner", "test-owner")
	cmd.Flags().Set("repo", "test-repo")
	cmd.Flags().Set("pr-number", "123")
	cmd.Flags().Set("github-token", "test-token")
	cmd.Flags().Set("ai-api-key", "test-key")

	// Directly call RunE with a recover block to isolate panic
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Panic during cmd.RunE execution: %v", r)
			}
		}()
		// Suppress log output during test for cleaner stdout
		log.SetOutput(new(bytes.Buffer))
		defer log.SetOutput(os.Stderr) // Reset log output after test

		err := cmd.RunE(cmd, []string{}) // Use RunE to get error
		if err != nil {
			t.Fatalf("cmd.RunE returned error: %v", err)
		}
	}()

	// For now, no output assertion as stdout is not captured.
	// We are just checking for panic/error.
}

func TestAnalyzePRCommand_GitHubError(t *testing.T) {
	resetGlobalClients()
	defer resetGlobalClients()

	// Setup mock GitHub client to return an error
	mockGh := &MockGitHubClient{
		GetPullRequestFunc: func(owner, repo string, prNum int) (*ghlib.PullRequest, error) {
			return nil, errors.New("github API error")
		},
	}
	githubClient = mockGh

	// Setup mock AI client (should not be called)
	mockAI := &MockAIClient{
		GetConflictPredictionFunc: func(prompt string) (string, error) {
			t.Fatal("AI client should not be called on GitHub error")
			return "", nil
		},
	}
	aiClient = mockAI

	// Create a new prCmd instance for this test
	cmd := NewPrCmd()

	// Set flags
	cmd.Flags().Set("owner", "test-owner")
	cmd.Flags().Set("repo", "test-repo")
	cmd.Flags().Set("pr-number", "123")
	cmd.Flags().Set("github-token", "test-token")
	cmd.Flags().Set("ai-api-key", "test-key")

	// Expect an error return
	log.SetOutput(new(bytes.Buffer)) // Suppress log output
	defer log.SetOutput(os.Stderr)

	err := cmd.RunE(cmd, []string{})

	if err == nil || !strings.Contains(err.Error(), "github API error") {
		t.Errorf("Expected GitHub API error, got: %v", err)
	}
}

func TestAnalyzePRCommand_AIError(t *testing.T) {
	resetGlobalClients()
	defer resetGlobalClients()

	// Setup mock GitHub client (success)
	mockGh := &MockGitHubClient{
		GetPullRequestFunc: func(owner, repo string, prNum int) (*ghlib.PullRequest, error) {
			return &ghlib.PullRequest{
				Number: ghlib.Int(123),
				Title:  ghlib.String("Test PR"),
				Base:   &ghlib.PullRequestBranch{Ref: ghlib.String("main")},
			}, nil
		},
		GetPullRequestFilesFunc: func(owner, repo string, prNum int) ([]*ghlib.CommitFile, error) {
			return []*ghlib.CommitFile{}, nil
		},
		GetRecentMergedPRsFunc: func(owner, repo, targetBranch string) ([]*ghlib.PullRequest, error) {
			return []*ghlib.PullRequest{}, nil
		},
		GetOpenPRsFunc: func(owner, repo, targetBranch string) ([]*ghlib.PullRequest, error) {
			return []*ghlib.PullRequest{}, nil
		},
	}
	githubClient = mockGh

	// Setup mock AI client to return an error
	mockAI := &MockAIClient{
		GetConflictPredictionFunc: func(prompt string) (string, error) {
			return "", errors.New("AI service error")
		},
	}
	aiClient = mockAI

	// Create a new prCmd instance for this test
	cmd := NewPrCmd()

	// Set flags
	cmd.Flags().Set("owner", "test-owner")
	cmd.Flags().Set("repo", "test-repo")
	cmd.Flags().Set("pr-number", "123")
	cmd.Flags().Set("github-token", "test-token")
	cmd.Flags().Set("ai-api-key", "test-key")

	// Expect an error return
	log.SetOutput(new(bytes.Buffer)) // Suppress log output
	defer log.SetOutput(os.Stderr)

	err := cmd.RunE(cmd, []string{})

	if err == nil || !strings.Contains(err.Error(), "AI service error") {
		t.Errorf("Expected AI service error, got: %v", err)
	}
}

func TestAnalyzePRCommand_InvalidAIProvider(t *testing.T) {
	resetGlobalClients()
	defer resetGlobalClients()

	// Create a new prCmd instance for this test
	cmd := NewPrCmd()

	// Set flags - specifically invalid AI provider
	cmd.Flags().Set("ai-provider", "unknown")
	cmd.Flags().Set("owner", "test-owner")
	cmd.Flags().Set("repo", "test-repo")
	cmd.Flags().Set("pr-number", "123")
	cmd.Flags().Set("github-token", "test-token")
	cmd.Flags().Set("ai-api-key", "test-key")


	log.SetOutput(new(bytes.Buffer))
	defer log.SetOutput(os.Stderr)

	err := cmd.RunE(cmd, []string{})

	if err == nil || !strings.Contains(err.Error(), "unsupported AI provider: unknown") {
		t.Errorf("Expected unsupported AI provider error, got: %v", err)
	}
}

func TestAnalyzePRCommand_Minimal(t *testing.T) {
	defer resetGlobalClients() // Ensure cleanup after test

	// Setup mock GitHub client
	mockGh := &MockGitHubClient{
		GetPullRequestFunc: func(owner, repo string, prNum int) (*ghlib.PullRequest, error) {
			return &ghlib.PullRequest{
				Number: ghlib.Int(123),
				Title:  ghlib.String("Minimal PR"),
				Base:   &ghlib.PullRequestBranch{Ref: ghlib.String("main")},
			}, nil
		},
		GetPullRequestFilesFunc: func(owner, repo string, prNum int) ([]*ghlib.CommitFile, error) {
			return []*ghlib.CommitFile{}, nil
		},
		GetRecentMergedPRsFunc: func(owner, repo, targetBranch string) ([]*ghlib.PullRequest, error) {
			return []*ghlib.PullRequest{}, nil
		},
		GetOpenPRsFunc: func(owner, repo, targetBranch string) ([]*ghlib.PullRequest, error) {
			return []*ghlib.PullRequest{}, nil
		},
	}
	githubClient = mockGh

	// Setup mock AI client
	mockAI := &MockAIClient{
		GetConflictPredictionFunc: func(prompt string) (string, error) {
			return "Minimal AI Prediction", nil
		},
	}
	aiClient = mockAI

	// Create a new prCmd instance for this test
	cmd := NewPrCmd()

	// Set flags
	cmd.Flags().Set("owner", "test-owner")
	cmd.Flags().Set("repo", "test-repo")
	cmd.Flags().Set("pr-number", "123")
	cmd.Flags().Set("github-token", "test-token")
	cmd.Flags().Set("ai-api-key", "test-key")

	// Directly call RunE to isolate panic
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Panic during cmd.RunE execution: %v", r)
			}
		}()
		log.SetOutput(new(bytes.Buffer))
		defer log.SetOutput(os.Stderr)

		err := cmd.RunE(cmd, []string{})
		if err != nil {
			t.Fatalf("cmd.RunE returned error: %v", err)
		}
	}()
}