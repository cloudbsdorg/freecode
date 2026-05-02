package provider

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestNewLiteLLMProvider(t *testing.T) {
	p := NewLiteLLMProvider("http://localhost:4000", "test-key")
	if p == nil {
		t.Fatal("NewLiteLLMProvider returned nil")
	}
	if p.BaseURL != "http://localhost:4000" {
		t.Errorf("BaseURL = %q, want %q", p.BaseURL, "http://localhost:4000")
	}
	if p.APIKey != "test-key" {
		t.Errorf("APIKey = %q, want %q", p.APIKey, "test-key")
	}
	if p.Client == nil {
		t.Error("Client should not be nil")
	}
	if p.Client.Timeout != 120*time.Second {
		t.Errorf("Client.Timeout = %v, want %v", p.Client.Timeout, 120*time.Second)
	}
}

func TestNewLiteLLMProviderTrailingSlash(t *testing.T) {
	p := NewLiteLLMProvider("http://localhost:4000/", "test-key")
	if p.BaseURL != "http://localhost:4000" {
		t.Errorf("BaseURL = %q, want %q (trailing slash should be removed)", p.BaseURL, "http://localhost:4000")
	}
}

func TestLiteLLMProviderName(t *testing.T) {
	p := NewLiteLLMProvider("http://localhost:4000", "test-key")
	if p.Name() != "litellm" {
		t.Errorf("Name() = %q, want %q", p.Name(), "litellm")
	}
}

func TestLiteLLMProviderGenerate(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Content-Type = %q, want %q", r.Header.Get("Content-Type"), "application/json")
		}
		if !strings.Contains(r.Header.Get("Authorization"), "test-key") {
			t.Errorf("Authorization header missing or wrong")
		}

		var body map[string]interface{}
		json.NewDecoder(r.Body).Decode(&body)

		if body["model"] != "ollama/llama3" {
			t.Errorf("model = %v, want %q", body["model"], "ollama/llama3")
		}
		if body["temperature"] != 0.7 {
			t.Errorf("temperature = %v, want %v", body["temperature"], 0.7)
		}
		if body["max_tokens"] != 100.0 {
			t.Errorf("max_tokens = %v, want %v", body["max_tokens"], 100.0)
		}

		json.NewEncoder(w).Encode(map[string]interface{}{
			"choices": []map[string]interface{}{{
				"message":       map[string]string{"content": "hello from llama3"},
				"finish_reason": "stop",
			}},
			"usage": map[string]int{
				"prompt_tokens":     10,
				"completion_tokens": 20,
			},
		})
	}))
	defer server.Close()

	p := NewLiteLLMProvider(server.URL, "test-key")
	resp, err := p.Generate(context.Background(), &Request{
		Model:       "ollama/llama3",
		Messages:    []Message{{Role: "user", Content: "hi"}},
		Temperature: 0.7,
		MaxTokens:   100,
	})

	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}
	if resp.Content != "hello from llama3" {
		t.Errorf("Content = %q, want %q", resp.Content, "hello from llama3")
	}
	if resp.StopReason != "stop" {
		t.Errorf("StopReason = %q, want %q", resp.StopReason, "stop")
	}
	if resp.Usage.InputTokens != 10 {
		t.Errorf("Usage.InputTokens = %d, want %d", resp.Usage.InputTokens, 10)
	}
	if resp.Usage.OutputTokens != 20 {
		t.Errorf("Usage.OutputTokens = %d, want %d", resp.Usage.OutputTokens, 20)
	}
}

func TestLiteLLMProviderGenerateNoAPIKey(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "" {
			t.Errorf("Authorization header should be empty when no API key, got %q", r.Header.Get("Authorization"))
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"choices": []map[string]interface{}{{
				"message":       map[string]string{"content": "response"},
				"finish_reason": "stop",
			}},
		})
	}))
	defer server.Close()

	p := NewLiteLLMProvider(server.URL, "")
	resp, err := p.Generate(context.Background(), &Request{
		Model:    "ollama/llama3",
		Messages: []Message{{Role: "user", Content: "hi"}},
	})

	if err != nil {
		t.Fatalf("Generate() with empty API key error = %v", err)
	}
	if resp.Content != "response" {
		t.Errorf("Content = %q, want %q", resp.Content, "response")
	}
}

func TestLiteLLMProviderGenerateAnthropicModel(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]interface{}
		json.NewDecoder(r.Body).Decode(&body)

		if body["model"] != "anthropic/claude-3-opus" {
			t.Errorf("model = %v, want %q", body["model"], "anthropic/claude-3-opus")
		}

		json.NewEncoder(w).Encode(map[string]interface{}{
			"choices": []map[string]interface{}{{
				"message":       map[string]string{"content": "claude response"},
				"finish_reason": "end_turn",
			}},
		})
	}))
	defer server.Close()

	p := NewLiteLLMProvider(server.URL, "sk-ant")
	resp, err := p.Generate(context.Background(), &Request{
		Model:    "anthropic/claude-3-opus",
		Messages: []Message{{Role: "user", Content: "hello"}},
	})

	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}
	if resp.Content != "claude response" {
		t.Errorf("Content = %q, want %q", resp.Content, "claude response")
	}
}

func TestLiteLLMProviderGenerateOpenAIModel(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]interface{}
		json.NewDecoder(r.Body).Decode(&body)

		if body["model"] != "openai/gpt-4-turbo" {
			t.Errorf("model = %v, want %q", body["model"], "openai/gpt-4-turbo")
		}

		json.NewEncoder(w).Encode(map[string]interface{}{
			"choices": []map[string]interface{}{{
				"message":       map[string]string{"content": "gpt response"},
				"finish_reason": "stop",
			}},
		})
	}))
	defer server.Close()

	p := NewLiteLLMProvider(server.URL, "sk-openai")
	resp, err := p.Generate(context.Background(), &Request{
		Model:    "openai/gpt-4-turbo",
		Messages: []Message{{Role: "user", Content: "hello"}},
	})

	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}
	if resp.Content != "gpt response" {
		t.Errorf("Content = %q, want %q", resp.Content, "gpt response")
	}
}

func TestLiteLLMProviderGenerateWithSystemMessage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]interface{}
		json.NewDecoder(r.Body).Decode(&body)

		messages, ok := body["messages"].([]interface{})
		if !ok {
			t.Fatal("messages is not a slice")
		}
		if len(messages) != 2 {
			t.Errorf("len(messages) = %d, want 2", len(messages))
		}

		json.NewEncoder(w).Encode(map[string]interface{}{
			"choices": []map[string]interface{}{{
				"message":       map[string]string{"content": "result"},
				"finish_reason": "stop",
			}},
		})
	}))
	defer server.Close()

	p := NewLiteLLMProvider(server.URL, "test-key")
	_, err := p.Generate(context.Background(), &Request{
		Model: "ollama/llama3",
		Messages: []Message{
			{Role: "system", Content: "you are helpful"},
			{Role: "user", Content: "hi"},
		},
	})

	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}
}

func TestLiteLLMProviderGenerateNoChoices(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"choices": []map[string]interface{}{},
		})
	}))
	defer server.Close()

	p := NewLiteLLMProvider(server.URL, "test-key")
	_, err := p.Generate(context.Background(), &Request{
		Model:    "ollama/llama3",
		Messages: []Message{{Role: "user", Content: "hi"}},
	})

	if err == nil {
		t.Error("Generate() expected error for empty choices")
	}
	if err != nil && err.Error() != "no choices in response" {
		t.Errorf("error = %q, want %q", err.Error(), "no choices in response")
	}
}

func TestLiteLLMProviderGenerateHTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": map[string]string{
				"message": "internal server error",
				"type":    "internal_error",
			},
		})
	}))
	defer server.Close()

	p := NewLiteLLMProvider(server.URL, "test-key")
	_, err := p.Generate(context.Background(), &Request{
		Model:    "ollama/llama3",
		Messages: []Message{{Role: "user", Content: "hi"}},
	})

	if err == nil {
		t.Error("Generate() expected error for HTTP 500")
	}
}

func TestLiteLLMProviderGenerateLiteLLMError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": map[string]string{
				"message": "model not found",
				"type":    "invalid_request_error",
			},
		})
	}))
	defer server.Close()

	p := NewLiteLLMProvider(server.URL, "test-key")
	_, err := p.Generate(context.Background(), &Request{
		Model:    "ollama/nonexistent",
		Messages: []Message{{Role: "user", Content: "hi"}},
	})

	if err == nil {
		t.Error("Generate() expected error for LiteLLM error response")
	}
}

func TestLiteLLMProviderGenerateTimeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"choices": []map[string]interface{}{{
				"message": map[string]string{"content": "late"},
			}},
		})
	}))
	defer server.Close()

	p := &LiteLLMProvider{
		BaseURL: server.URL,
		APIKey:  "test-key",
		Client:  &http.Client{Timeout: 50 * time.Millisecond},
	}

	_, err := p.Generate(context.Background(), &Request{
		Model:    "ollama/llama3",
		Messages: []Message{{Role: "user", Content: "hi"}},
	})

	if err == nil {
		t.Error("Generate() expected error for timeout")
	}
}

func TestLiteLLMProviderGenerateConnectionError(t *testing.T) {
	p := NewLiteLLMProvider("http://localhost:99999", "test-key")
	_, err := p.Generate(context.Background(), &Request{
		Model:    "ollama/llama3",
		Messages: []Message{{Role: "user", Content: "hi"}},
	})

	if err == nil {
		t.Error("Generate() expected error for connection failure")
	}
}

func TestLiteLLMProviderGenerateInvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("not json"))
	}))
	defer server.Close()

	p := NewLiteLLMProvider(server.URL, "test-key")
	_, err := p.Generate(context.Background(), &Request{
		Model:    "ollama/llama3",
		Messages: []Message{{Role: "user", Content: "hi"}},
	})

	if err == nil {
		t.Error("Generate() expected error for invalid JSON response")
	}
}

func TestLiteLLMProviderGenerateMultipleMessages(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]interface{}
		json.NewDecoder(r.Body).Decode(&body)

		messages, ok := body["messages"].([]interface{})
		if !ok {
			t.Fatal("messages is not a slice")
		}
		if len(messages) != 4 {
			t.Errorf("len(messages) = %d, want 4", len(messages))
		}

		json.NewEncoder(w).Encode(map[string]interface{}{
			"choices": []map[string]interface{}{{
				"message":       map[string]string{"content": "response"},
				"finish_reason": "stop",
			}},
		})
	}))
	defer server.Close()

	p := NewLiteLLMProvider(server.URL, "test-key")
	_, err := p.Generate(context.Background(), &Request{
		Model: "ollama/llama3",
		Messages: []Message{
			{Role: "system", Content: "you are a helpful assistant"},
			{Role: "user", Content: "hello"},
			{Role: "assistant", Content: "hi there"},
			{Role: "user", Content: "how are you?"},
		},
	})

	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}
}

func TestLiteLLMProviderGenerateZeroTemperature(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]interface{}
		json.NewDecoder(r.Body).Decode(&body)

		if body["temperature"] != 0.0 {
			t.Errorf("temperature = %v, want 0.0", body["temperature"])
		}

		json.NewEncoder(w).Encode(map[string]interface{}{
			"choices": []map[string]interface{}{{
				"message":       map[string]string{"content": "deterministic"},
				"finish_reason": "stop",
			}},
		})
	}))
	defer server.Close()

	p := NewLiteLLMProvider(server.URL, "test-key")
	resp, err := p.Generate(context.Background(), &Request{
		Model:       "ollama/llama3",
		Messages:    []Message{{Role: "user", Content: "hi"}},
		Temperature: 0,
		MaxTokens:   50,
	})

	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}
	if resp.Content != "deterministic" {
		t.Errorf("Content = %q, want %q", resp.Content, "deterministic")
	}
}

func TestLiteLLMProviderGenerateMaxTokens(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]interface{}
		json.NewDecoder(r.Body).Decode(&body)

		if body["max_tokens"] != 500.0 {
			t.Errorf("max_tokens = %v, want 500", body["max_tokens"])
		}

		json.NewEncoder(w).Encode(map[string]interface{}{
			"choices": []map[string]interface{}{{
				"message":       map[string]string{"content": "limited response"},
				"finish_reason": "length",
			}},
			"usage": map[string]int{
				"prompt_tokens":     10,
				"completion_tokens": 500,
			},
		})
	}))
	defer server.Close()

	p := NewLiteLLMProvider(server.URL, "test-key")
	resp, err := p.Generate(context.Background(), &Request{
		Model:       "ollama/llama3",
		Messages:    []Message{{Role: "user", Content: "hi"}},
		Temperature: 0.7,
		MaxTokens:   500,
	})

	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}
	if resp.Usage.OutputTokens != 500 {
		t.Errorf("Usage.OutputTokens = %d, want 500", resp.Usage.OutputTokens)
	}
}
