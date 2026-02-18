package ai

import "fmt"

// AIClient is an interface for different AI providers.
type AIClient interface {
	GetConflictPrediction(prompt string) (string, error)
}

// Provider represents the AI service provider.
type Provider string

const (
	ProviderOpenAI    Provider = "openai"
	ProviderGemini    Provider = "gemini"
	ProviderAnthropic Provider = "anthropic"
)

// NewAIClient creates a new AI client based on the specified provider and API key.
func NewAIClient(provider Provider, apiKey string) (AIClient, error) {
	switch provider {
	case ProviderOpenAI:
		return NewOpenAIClient(apiKey), nil
	case ProviderGemini:
		return NewGeminiClient(apiKey), nil
	case ProviderAnthropic:
		return NewAnthropicClient(apiKey), nil
	default:
		return nil, fmt.Errorf("unsupported AI provider: %s", provider)
	}
}
