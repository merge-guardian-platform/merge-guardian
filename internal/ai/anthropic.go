package ai

import (
	"context"

	"github.com/liushuangls/go-anthropic/v2"
)

// newAnthropicClientFunc is a package-level variable to allow mocking anthropic.NewClient in tests.
var newAnthropicClientFunc = func(apiKey string, opts ...anthropic.ClientOption) *anthropic.Client {
	return anthropic.NewClient(apiKey, opts...)
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
	resp, err := c.client.CreateMessages(
		c.ctx,
		anthropic.MessagesRequest{
			Model: anthropic.ModelClaude3Dot5Sonnet20241022,
			Messages: []anthropic.Message{
				{
					Role: anthropic.RoleUser,
					Content: []anthropic.MessageContent{
						anthropic.NewTextMessageContent(prompt),
					},
				},
			},
			MaxTokens: 4000,
		})

	if err != nil {
		return "", err
	}

	if len(resp.Content) == 0 {
		return "", nil
	}

	return *resp.Content[0].Text, nil
}
