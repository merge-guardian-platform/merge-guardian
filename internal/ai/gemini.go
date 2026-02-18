package ai

import (
	"context"
	"fmt"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

// newGeminiClientFunc is a package-level variable to allow mocking genai.NewClient in tests.
var newGeminiClientFunc = func(ctx context.Context, opts ...option.ClientOption) (*genai.Client, error) {
	return genai.NewClient(ctx, opts...)
}

// newGeminiGenerativeModelFunc is a package-level variable to allow mocking model creation in tests.
var newGeminiGenerativeModelFunc = func(client *genai.Client, modelName string) *genai.GenerativeModel {
	return client.GenerativeModel(modelName)
}

// GeminiClient provides a wrapper around the Google Gemini API client.
type GeminiClient struct {
	model *genai.GenerativeModel
	ctx   context.Context
}

// NewGeminiClient creates a new Gemini client.
func NewGeminiClient(apiKey string) *GeminiClient {
	ctx := context.Background()
	client, err := newGeminiClientFunc(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		// Log error or handle it as appropriate for a constructor
		// For now, we'll panic as NewAIClient expects a valid client.
		panic(fmt.Sprintf("Failed to create Gemini client: %v", err))
	}

	// Use a generative model, e.g., "gemini-2.5-flash"
	model := newGeminiGenerativeModelFunc(client, "gemini-2.5-flash")

	return &GeminiClient{
		model: model,
		ctx:   ctx,
	}
}

// GetConflictPrediction sends a prompt to the Gemini API for conflict prediction.
func (c *GeminiClient) GetConflictPrediction(prompt string) (string, error) {
	resp, err := c.model.GenerateContent(c.ctx, genai.Text(prompt))
	if err != nil {
		return "", fmt.Errorf("failed to generate content from Gemini: %w", err)
	}

	if len(resp.Candidates) == 0 || resp.Candidates[0].Content == nil || len(resp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no content generated from Gemini API")
	}

	var prediction string
	for _, part := range resp.Candidates[0].Content.Parts {
		if text, ok := part.(genai.Text); ok {
			prediction += string(text)
		}
	}
	return prediction, nil
}
