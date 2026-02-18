package github

import (
	"context"
	"fmt"
	"time"

	"github.com/google/go-github/v58/github"
	"golang.org/x/oauth2"
)

// Client provides a wrapper around the GitHub API client.
type Client struct {
	ghClient *github.Client
	ctx      context.Context
}

// NewClient creates a new GitHub client.
func NewClient(token string) *Client {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	return &Client{
		ghClient: github.NewClient(tc),
		ctx:      ctx,
	}
}

// GetPullRequest fetches details of a specific pull request.
func (c *Client) GetPullRequest(owner, repo string, prNumber int) (*github.PullRequest, error) {
	pr, _, err := c.ghClient.PullRequests.Get(c.ctx, owner, repo, prNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to get pull request %d: %w", prNumber, err)
	}
	return pr, nil
}

// GetPullRequestFiles fetches files changed in a specific pull request.
func (c *Client) GetPullRequestFiles(owner, repo string, prNumber int) ([]*github.CommitFile, error) {
	files, _, err := c.ghClient.PullRequests.ListFiles(c.ctx, owner, repo, prNumber, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list files for pull request %d: %w", prNumber, err)
	}
	return files, nil
}

// GetRecentMergedPRs fetches recently merged pull requests (last 24 hours) for a given branch.
func (c *Client) GetRecentMergedPRs(owner, repo, targetBranch string) ([]*github.PullRequest, error) {
	// GitHub API doesn't allow direct filtering by merged date, so we fetch recent PRs and filter.
	// This might be inefficient for very active repos, but serves the initial purpose.
	opt := &github.PullRequestListOptions{
		State:     "closed",
		Sort:      "updated",
		Direction: "desc",
		Base:      targetBranch,
		ListOptions: github.ListOptions{
			PerPage: 100, // Fetch a good number to ensure we get recent ones
		},
	}

	prs, _, err := c.ghClient.PullRequests.List(c.ctx, owner, repo, opt)
	if err != nil {
		return nil, fmt.Errorf("failed to list recent closed pull requests: %w", err)
	}

	var recentMerged []*github.PullRequest
	threshold := time.Now().Add(-24 * time.Hour) // Last 24 hours

	for _, pr := range prs {
		if pr.MergedAt != nil && pr.MergedAt.After(threshold) {
			recentMerged = append(recentMerged, pr)
		}
	}
	return recentMerged, nil
}

// GetOpenPRs fetches currently open pull requests targeting a specific branch.
func (c *Client) GetOpenPRs(owner, repo, targetBranch string) ([]*github.PullRequest, error) {
	opt := &github.PullRequestListOptions{
		State: "open",
		Sort:  "created",
		Direction: "desc",
		Base:  targetBranch,
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}

	prs, _, err := c.ghClient.PullRequests.List(c.ctx, owner, repo, opt)
	if err != nil {
		return nil, fmt.Errorf("failed to list open pull requests: %w", err)
	}
	return prs, nil
}
