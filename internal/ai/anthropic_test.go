package ai

import (

	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/liushuangls/go-anthropic/v2"
)

// setupMockAnthropicServer creates a mock HTTP server for Anthropic API.
func setupMockAnthropicServer(t *testing.T, response string, statusCode int) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// The client appends /messages to the base URL
		if r.URL.Path != "/messages" {
			t.Errorf("Unexpected Anthropic API path: %s", r.URL.Path)
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


func TestNewAIClient_Anthropic(t *testing.T) {
	client, err := NewAIClient(ProviderAnthropic, "test-key")
	if err != nil {
		t.Fatalf("NewAIClient failed: %v", err)
	}
	if _, ok := client.(*AnthropicClient); !ok {
		t.Errorf("NewAIClient returned wrong type for Anthropic: got %T, want *AnthropicClient", client)
	}
}

func TestAnthropicClient_GetConflictPrediction_Success(t *testing.T) {
	mockResponse := `{
		"id": "msg_123",
		"type": "message",
		"role": "assistant",
		"model": "claude-3-5-sonnet-20241022",
		"content": [
			{
				"type": "text",
				"text": "Anthropic Prediction"
			}
		],
		"stop_reason": "end_turn",
		"stop_sequence": null,
		"usage": {
			"input_tokens": 10,
			"output_tokens": 10
		}
	}`
	server := setupMockAnthropicServer(t, mockResponse, http.StatusOK)
	defer server.Close()

	// Override package-level newAnthropicClientFunc
	originalNewAnthropicClientFunc := newAnthropicClientFunc
	newAnthropicClientFunc = func(apiKey string, opts ...anthropic.ClientOption) *anthropic.Client {
		// Create a client that points to the mock server
		return anthropic.NewClient(apiKey, anthropic.WithBaseURL(server.URL))
	}
	defer func() { newAnthropicClientFunc = originalNewAnthropicClientFunc }()

	client := NewAnthropicClient("test-key")
	prediction, err := client.GetConflictPrediction("test prompt")
	if err != nil {
		t.Fatalf("GetConflictPrediction failed: %v", err)
	}
	if prediction != "Anthropic Prediction" {
		t.Errorf("GetConflictPrediction returned wrong prediction: got %q, want %q", prediction, "Anthropic Prediction")
	}
}

func TestAnthropicClient_GetConflictPrediction_APIError(t *testing.T) {
	mockResponse := `{"error": {"type": "api_error", "message": "API Error"}}`
	server := setupMockAnthropicServer(t, mockResponse, http.StatusInternalServerError)
	defer server.Close()

	originalNewAnthropicClientFunc := newAnthropicClientFunc
	newAnthropicClientFunc = func(apiKey string, opts ...anthropic.ClientOption) *anthropic.Client {
		return anthropic.NewClient(apiKey, anthropic.WithBaseURL(server.URL))
	}
	defer func() { newAnthropicClientFunc = originalNewAnthropicClientFunc }()

	client := NewAnthropicClient("test-key")
	_, err := client.GetConflictPrediction("test prompt")
	if err == nil {
		t.Fatal("GetConflictPrediction did not return an error for API error")
	}
}

func TestAnthropicClient_GetConflictPrediction_NoContent(t *testing.T) {
	// Simulate response with empty content list
	mockResponse := `{
		"id": "msg_123",
		"type": "message",
		"role": "assistant",
		"content": [],
		"usage": {}
	}`
	server := setupMockAnthropicServer(t, mockResponse, http.StatusOK)
	defer server.Close()

	originalNewAnthropicClientFunc := newAnthropicClientFunc
	newAnthropicClientFunc = func(apiKey string, opts ...anthropic.ClientOption) *anthropic.Client {
		return anthropic.NewClient(apiKey, anthropic.WithBaseURL(server.URL))
	}
	defer func() { newAnthropicClientFunc = originalNewAnthropicClientFunc }()

	client := NewAnthropicClient("test-key")
	prediction, err := client.GetConflictPrediction("test prompt")
	
	if err != nil {
		t.Fatalf("GetConflictPrediction failed with unexpected error: %v", err)
	}
	if prediction != "" {
		t.Errorf("Expected empty prediction for no content, got %q", prediction)
	}
}
