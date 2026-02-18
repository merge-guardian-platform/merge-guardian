package ai

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings" // Added import
	"testing"

	"github.com/google/generative-ai-go/genai"
	openai "github.com/sashabaranov/go-openai"
	"google.golang.org/api/option"
)

// setupMockOpenAIServer creates a mock HTTP server for OpenAI API.
func setupMockOpenAIServer(t *testing.T, response string, statusCode int) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/chat/completions" {
			t.Errorf("Unexpected OpenAI API path: %s", r.URL.Path)
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}
		if r.Method != "POST" {
			t.Errorf("Unexpected HTTP method: %s", r.Method)
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		w.WriteHeader(statusCode)
		fmt.Fprint(w, response)
	}))
}

// setupMockGeminiServer creates a mock HTTP server for Gemini API.
func setupMockGeminiServer(t *testing.T, response string, statusCode int) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1beta/models/gemini-pro:generateContent" {
			t.Errorf("Unexpected Gemini API path: %s", r.URL.Path)
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}
		if r.Method != "POST" {
			t.Errorf("Unexpected HTTP method: %s", r.Method)
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
			}
		w.WriteHeader(statusCode)
		fmt.Fprint(w, response)
	}))
}

func TestNewAIClient_OpenAI(t *testing.T) {
	// Temporarily override the newOpenAIClientFunc for this test
	originalNewClientFunc := newOpenAIClientFunc
	newOpenAIClientFunc = func(apiKey string) *openai.Client {
		return &openai.Client{} // Return a dummy client
	}
	defer func() { newOpenAIClientFunc = originalNewClientFunc }()

	client, err := NewAIClient(ProviderOpenAI, "test-key")
	if err != nil {
		t.Fatalf("NewAIClient failed: %v", err)
	}
	if _, ok := client.(*OpenAIClient); !ok {
		t.Errorf("NewAIClient returned wrong type for OpenAI: got %T, want *OpenAIClient", client)
	}
}

func TestNewAIClient_Gemini(t *testing.T) {
	// Temporarily override newGeminiClientFunc for this test
	originalNewGeminiClientFunc := newGeminiClientFunc
	newGeminiClientFunc = func(ctx context.Context, opts ...option.ClientOption) (*genai.Client, error) {
		return &genai.Client{}, nil // Return a dummy client
	}
	defer func() { newGeminiClientFunc = originalNewGeminiClientFunc }()

	// Temporarily override newGeminiGenerativeModelFunc for this test
	originalNewGeminiGenerativeModelFunc := newGeminiGenerativeModelFunc
	newGeminiGenerativeModelFunc = func(client *genai.Client, modelName string) *genai.GenerativeModel {
		return &genai.GenerativeModel{} // Return a dummy model
	}
	defer func() { newGeminiGenerativeModelFunc = originalNewGeminiGenerativeModelFunc }()

	client, err := NewAIClient(ProviderGemini, "test-key")
	if err != nil {
		t.Fatalf("NewAIClient failed: %v", err)
	}
	if _, ok := client.(*GeminiClient); !ok {
		t.Errorf("NewAIClient returned wrong type for Gemini: got %T, want *GeminiClient", client)
	}
}

func TestNewAIClient_Unsupported(t *testing.T) {
	_, err := NewAIClient("unsupported", "test-key")
	if err == nil {
		t.Fatal("NewAIClient for unsupported provider did not return an error")
	}
	expectedError := "unsupported AI provider: unsupported"
	if err.Error() != expectedError {
		t.Errorf("NewAIClient returned wrong error: got %q, want %q", err.Error(), expectedError)
	}
}

func TestOpenAIClient_GetConflictPrediction_Success(t *testing.T) {
	mockResponse := `{"choices": [{"message": {"content": "OpenAI Prediction"}}]}`
	server := setupMockOpenAIServer(t, mockResponse, http.StatusOK)
	defer server.Close()

	// Temporarily override the package-level newOpenAIClientFunc for testing
	originalNewClientFunc := newOpenAIClientFunc
	newOpenAIClientFunc = func(apiKey string) *openai.Client {
		cfg := openai.DefaultConfig(apiKey)
		cfg.BaseURL = server.URL + "/v1" // Match the mock server's base URL
		return openai.NewClientWithConfig(cfg)
	}
	defer func() { newOpenAIClientFunc = originalNewClientFunc }()

	client := NewOpenAIClient("test-key")
	prediction, err := client.GetConflictPrediction("test prompt")
	if err != nil {
		t.Fatalf("GetConflictPrediction failed: %v", err)
	}
	if prediction != "OpenAI Prediction" {
		t.Errorf("GetConflictPrediction returned wrong prediction: got %q, want %q", prediction, "OpenAI Prediction")
	}
}

func TestOpenAIClient_GetConflictPrediction_APIError(t *testing.T) {
	mockResponse := `{"error": {"message": "API Error"}}`
	server := setupMockOpenAIServer(t, mockResponse, http.StatusInternalServerError)
	defer server.Close()

	// Override the package-level newOpenAIClientFunc
	originalNewClientFunc := newOpenAIClientFunc
	newOpenAIClientFunc = func(apiKey string) *openai.Client {
		cfg := openai.DefaultConfig(apiKey)
		cfg.BaseURL = server.URL + "/v1"
		return openai.NewClientWithConfig(cfg)
	}
	defer func() { newOpenAIClientFunc = originalNewClientFunc }()

	client := NewOpenAIClient("test-key")
	_, err := client.GetConflictPrediction("test prompt")
	if err == nil {
		t.Fatal("GetConflictPrediction did not return an error for API error")
	}
	if !strings.Contains(err.Error(), "failed to get chat completion") {
		t.Errorf("GetConflictPrediction returned wrong error: got %q", err.Error())
	}
}

func TestOpenAIClient_GetConflictPrediction_NoChoices(t *testing.T) {
	mockResponse := `{"choices": []}`
	server := setupMockOpenAIServer(t, mockResponse, http.StatusOK)
	defer server.Close()

	// Override the package-level newOpenAIClientFunc
	originalNewClientFunc := newOpenAIClientFunc
	newOpenAIClientFunc = func(apiKey string) *openai.Client {
		cfg := openai.DefaultConfig(apiKey)
		cfg.BaseURL = server.URL + "/v1"
		return openai.NewClientWithConfig(cfg)
	}
	defer func() { newOpenAIClientFunc = originalNewClientFunc }()

	client := NewOpenAIClient("test-key")
	_, err := client.GetConflictPrediction("test prompt")
	if err == nil {
		t.Fatal("GetConflictPrediction did not return an error for no choices")
	}
	expectedError := "no response choices from OpenAI API"
	if err.Error() != expectedError {
		t.Errorf("GetConflictPrediction returned wrong error: got %q, want %q", err.Error(), expectedError)
	}
}

func TestGeminiClient_GetConflictPrediction_Success(t *testing.T) {
	mockResponse := `{"candidates": [{"content": {"parts": [{"text": "Gemini Prediction"}]}}]}`
	server := setupMockGeminiServer(t, mockResponse, http.StatusOK)
	defer server.Close()

	// Override package-level newGeminiClientFunc
	originalNewGeminiClientFunc := newGeminiClientFunc
	newGeminiClientFunc = func(ctx context.Context, opts ...option.ClientOption) (*genai.Client, error) {
		return genai.NewClient(ctx, option.WithAPIKey("dummy-key"), option.WithEndpoint(server.URL))
	}
	defer func() { newGeminiClientFunc = originalNewGeminiClientFunc }()

	// Override package-level newGeminiGenerativeModelFunc
	originalNewGeminiGenerativeModelFunc := newGeminiGenerativeModelFunc
	newGeminiGenerativeModelFunc = func(client *genai.Client, modelName string) *genai.GenerativeModel {
		return client.GenerativeModel(modelName)
	}
	defer func() { newGeminiGenerativeModelFunc = originalNewGeminiGenerativeModelFunc }()

	client := NewGeminiClient("test-key")
	prediction, err := client.GetConflictPrediction("test prompt")
	if err != nil {
		t.Fatalf("GetConflictPrediction failed: %v", err)
	}
	if prediction != "Gemini Prediction" {
		t.Errorf("GetConflictPrediction returned wrong prediction: got %q, want %q", prediction, "Gemini Prediction")
	}
}

func TestGeminiClient_GetConflictPrediction_APIError(t *testing.T) {
	mockResponse := `{"error": {"message": "API Error"}}`
	server := setupMockGeminiServer(t, mockResponse, http.StatusInternalServerError)
	defer server.Close()

	originalNewGeminiClientFunc := newGeminiClientFunc
	newGeminiClientFunc = func(ctx context.Context, opts ...option.ClientOption) (*genai.Client, error) {
		return genai.NewClient(ctx, option.WithAPIKey("dummy-key"), option.WithEndpoint(server.URL))
	}
	defer func() { newGeminiClientFunc = originalNewGeminiClientFunc }()

	originalNewGeminiGenerativeModelFunc := newGeminiGenerativeModelFunc
	newGeminiGenerativeModelFunc = func(client *genai.Client, modelName string) *genai.GenerativeModel {
		return client.GenerativeModel(modelName)
	}
	defer func() { newGeminiGenerativeModelFunc = originalNewGeminiGenerativeModelFunc }()

	client := NewGeminiClient("test-key")
	_, err := client.GetConflictPrediction("test prompt")
	if err == nil {
		t.Fatal("GetConflictPrediction did not return an error for API error")
	}
	if !strings.Contains(err.Error(), "failed to generate content from Gemini") {
		t.Errorf("GetConflictPrediction returned wrong error: got %q", err.Error())
	}
}

func TestGeminiClient_GetConflictPrediction_NoContent(t *testing.T) {
	mockResponse := `{"candidates": []}`
	server := setupMockGeminiServer(t, mockResponse, http.StatusOK)
	defer server.Close()

	originalNewGeminiClientFunc := newGeminiClientFunc
	newGeminiClientFunc = func(ctx context.Context, opts ...option.ClientOption) (*genai.Client, error) {
		return genai.NewClient(ctx, option.WithAPIKey("dummy-key"), option.WithEndpoint(server.URL))
	}
	defer func() { newGeminiClientFunc = originalNewGeminiClientFunc }()

	originalNewGeminiGenerativeModelFunc := newGeminiGenerativeModelFunc
	newGeminiGenerativeModelFunc = func(client *genai.Client, modelName string) *genai.GenerativeModel {
		return client.GenerativeModel(modelName)
	}
	defer func() { newGeminiGenerativeModelFunc = originalNewGeminiGenerativeModelFunc }()

	client := NewGeminiClient("test-key")
	_, err := client.GetConflictPrediction("test prompt")
	if err == nil {
		t.Fatal("GetConflictPrediction did not return an error for no content")
	}
	expectedError := "no content generated from Gemini API"
	if err.Error() != expectedError {
		t.Errorf("GetConflictPrediction returned wrong error: got %q, want %q", err.Error(), expectedError)
	}
}

// Helper to check if a string contains a substring (for error messages)
func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[0:len(substr)] == substr
}
