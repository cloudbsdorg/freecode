package provider

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNewAnthropicProvider(t *testing.T) {
	p := NewAnthropicProvider("test-key")
	if p == nil {
		t.Fatal("NewAnthropicProvider returned nil")
	}
	if p.APIKey != "test-key" {
		t.Errorf("APIKey = %q, want %q", p.APIKey, "test-key")
	}
	if p.BaseURL != "https://api.anthropic.com/v1" {
		t.Errorf("BaseURL = %q, want %q", p.BaseURL, "https://api.anthropic.com/v1")
	}
}

func TestAnthropicProviderName(t *testing.T) {
	p := NewAnthropicProvider("test-key")
	if p.Name() != "anthropic" {
		t.Errorf("Name() = %q, want %q", p.Name(), "anthropic")
	}
}

func TestAnthropicProviderGenerate(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("x-api-key") != "test-key" {
			t.Errorf("Missing or wrong x-api-key header")
		}
		if r.Header.Get("anthropic-version") != "2023-06-01" {
			t.Errorf("Missing or wrong anthropic-version header")
		}

		resp := `{"content":[{"type":"text","text":"hello"}],"stop_reason":"end_turn","usage":{"input_tokens":10,"output_tokens":20}}`
		w.Write([]byte(resp))
	}))
	defer server.Close()

	p := &AnthropicProvider{
		APIKey:  "test-key",
		BaseURL: server.URL,
		Client:  server.Client(),
	}

	req := &Request{
		Model: "claude-3",
		Messages: []Message{
			{Role: "user", Content: "hi"},
		},
		Temperature: 0.7,
		MaxTokens:   100,
	}

	resp, err := p.Generate(context.Background(), req)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}
	if resp.Content != "hello" {
		t.Errorf("Content = %q, want %q", resp.Content, "hello")
	}
	if resp.StopReason != "end_turn" {
		t.Errorf("StopReason = %q, want %q", resp.StopReason, "end_turn")
	}
	if resp.Usage.InputTokens != 10 {
		t.Errorf("InputTokens = %d, want %d", resp.Usage.InputTokens, 10)
	}
	if resp.Usage.OutputTokens != 20 {
		t.Errorf("OutputTokens = %d, want %d", resp.Usage.OutputTokens, 20)
	}
}

func TestAnthropicProviderGenerateWithSystemMessage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]interface{}
		json.NewDecoder(r.Body).Decode(&body)

		if body["system"] != "you are helpful" {
			t.Errorf("system = %v, want %q", body["system"], "you are helpful")
		}

		resp := `{"content":[{"type":"text","text":"result"}],"stop_reason":"end_turn","usage":{"input_tokens":5,"output_tokens":5}}`
		w.Write([]byte(resp))
	}))
	defer server.Close()

	p := &AnthropicProvider{
		APIKey:  "test-key",
		BaseURL: server.URL,
		Client:  server.Client(),
	}

	req := &Request{
		Model: "claude-3",
		Messages: []Message{
			{Role: "system", Content: "you are helpful"},
			{Role: "user", Content: "hi"},
		},
	}

	_, err := p.Generate(context.Background(), req)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}
}

func TestAnthropicProviderGenerateRequestError(t *testing.T) {
	p := &AnthropicProvider{
		APIKey:  "test-key",
		BaseURL: "http://localhost:99999",
		Client:  &http.Client{Timeout: 1},
	}

	req := &Request{
		Model: "claude-3",
		Messages: []Message{
			{Role: "user", Content: "hi"},
		},
	}

	_, err := p.Generate(context.Background(), req)
	if err == nil {
		t.Error("Generate() expected error for invalid URL")
	}
}

func TestNewOpenAIProvider(t *testing.T) {
	p := NewOpenAIProvider("test-key")
	if p == nil {
		t.Fatal("NewOpenAIProvider returned nil")
	}
	if p.APIKey != "test-key" {
		t.Errorf("APIKey = %q, want %q", p.APIKey, "test-key")
	}
	if p.BaseURL != "https://api.openai.com/v1" {
		t.Errorf("BaseURL = %q, want %q", p.BaseURL, "https://api.openai.com/v1")
	}
}

func TestOpenAIProviderName(t *testing.T) {
	p := NewOpenAIProvider("test-key")
	if p.Name() != "openai" {
		t.Errorf("Name() = %q, want %q", p.Name(), "openai")
	}
}

func TestOpenAIProviderGenerate(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Authorization"), "test-key") {
			t.Errorf("Missing or wrong Authorization header")
		}

		resp := `{"choices":[{"message":{"content":"hello"},"finish_reason":"stop"}],"usage":{"prompt_tokens":10,"completion_tokens":20}}`
		w.Write([]byte(resp))
	}))
	defer server.Close()

	p := &OpenAIProvider{
		APIKey:  "test-key",
		BaseURL: server.URL,
		Client:  server.Client(),
	}

	req := &Request{
		Model: "gpt-4",
		Messages: []Message{
			{Role: "user", Content: "hi"},
		},
		Temperature: 0.7,
		MaxTokens:   100,
	}

	resp, err := p.Generate(context.Background(), req)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}
	if resp.Content != "hello" {
		t.Errorf("Content = %q, want %q", resp.Content, "hello")
	}
	if resp.StopReason != "stop" {
		t.Errorf("StopReason = %q, want %q", resp.StopReason, "stop")
	}
	if resp.Usage.InputTokens != 10 {
		t.Errorf("InputTokens = %d, want %d", resp.Usage.InputTokens, 10)
	}
	if resp.Usage.OutputTokens != 20 {
		t.Errorf("OutputTokens = %d, want %d", resp.Usage.OutputTokens, 20)
	}
}

func TestOpenAIProviderGenerateRequestError(t *testing.T) {
	p := &OpenAIProvider{
		APIKey:  "test-key",
		BaseURL: "http://localhost:99999",
		Client:  &http.Client{Timeout: 1},
	}

	req := &Request{
		Model: "gpt-4",
		Messages: []Message{
			{Role: "user", Content: "hi"},
		},
	}

	_, err := p.Generate(context.Background(), req)
	if err == nil {
		t.Error("Generate() expected error for invalid URL")
	}
}

func TestOpenAIProviderGenerateNoChoices(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := `{"choices":[],"usage":{"prompt_tokens":10,"completion_tokens":20}}`
		w.Write([]byte(resp))
	}))
	defer server.Close()

	p := &OpenAIProvider{
		APIKey:  "test-key",
		BaseURL: server.URL,
		Client:  server.Client(),
	}

	req := &Request{
		Model: "gpt-4",
		Messages: []Message{
			{Role: "user", Content: "hi"},
		},
	}

	_, err := p.Generate(context.Background(), req)
	if err == nil {
		t.Error("Generate() expected error for no choices")
	}
}