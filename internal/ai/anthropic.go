package ai

import (
	"context"

	anthropic "github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

// newAnthropicClientFunc is a package-level variable to allow mocking anthropic.NewClient in tests.
var newAnthropicClientFunc = func(apiKey string, opts ...option.RequestOption) *anthropic.Client {
	c := anthropic.NewClient(append(opts, option.WithAPIKey(apiKey))...)
	return &c
}

// AnthropicClient provides a wrapper around the Anthropic API client.
type AnthropicClient struct {
	client *anthropic.Client
	ctx    context.Context
}

// NewAnthropicClient creates a new Anthropic client.
func NewAnthropicClient(apiKey string) *AnthropicClient {
	ctx := context.Background()
	return &AnthropicClient{
		client: newAnthropicClientFunc(apiKey),
		ctx:    ctx,
	}
}

// GetConflictPrediction sends a prompt to the Anthropic API (Claude) for conflict prediction.
func (c *AnthropicClient) GetConflictPrediction(prompt string) (string, error) {
	resp, err := c.client.Messages.New(c.ctx, anthropic.MessageNewParams{
		Model:     anthropic.ModelClaude3_7Sonnet20250219,
		MaxTokens: 4000,
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(prompt)),
		},
	})

	if err != nil {
		return "", err
	}

	if len(resp.Content) == 0 {
		return "", nil
	}

	return resp.Content[0].Text, nil
}
