package ai

import (
	"context"
	"fmt"

	openai "github.com/sashabaranov/go-openai"
)

// newOpenAIClientFunc is a package-level variable to allow mocking openai.NewClient in tests.
var newOpenAIClientFunc = func(apiKey string) *openai.Client {
	return openai.NewClient(apiKey)
}

// OpenAIClient provides a wrapper around the OpenAI API client.
type OpenAIClient struct {
	aiClient *openai.Client
	ctx      context.Context
}

// NewOpenAIClient creates a new OpenAI client.
func NewOpenAIClient(apiKey string) *OpenAIClient {
	ctx := context.Background()
	return &OpenAIClient{
		aiClient: newOpenAIClientFunc(apiKey),
		ctx:      ctx,
	}
}

// GetConflictPrediction sends a prompt to the OpenAI API for conflict prediction.
func (c *OpenAIClient) GetConflictPrediction(prompt string) (string, error) {
	resp, err := c.aiClient.CreateChatCompletion(
		c.ctx,
		openai.ChatCompletionRequest{
			Model: openai.GPT4o, // Using GPT-4o as specified, can be configured later
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
		},
	)

	if err != nil {
		return "", fmt.Errorf("failed to get chat completion: %w", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no response choices from OpenAI API")
	}

	return resp.Choices[0].Message.Content, nil
}
