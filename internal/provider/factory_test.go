package provider

import (
	"os"
	"testing"
)

func TestNewProvider(t *testing.T) {
	os.Setenv("LITELLM_BASE_URL", "")
	os.Setenv("LITELLM_API_KEY", "")
	defer func() {
		os.Unsetenv("LITELLM_BASE_URL")
		os.Unsetenv("LITELLM_API_KEY")
	}()

	p := NewProvider("ollama/llama3")
	if p == nil {
		t.Fatal("NewProvider returned nil")
	}
	if p.Name() != "litellm" {
		t.Errorf("Name() = %q, want %q", p.Name(), "litellm")
	}
}

func TestNewProviderWithEnv(t *testing.T) {
	os.Setenv("LITELLM_BASE_URL", "http://custom:4000")
	os.Setenv("LITELLM_API_KEY", "test-key")
	defer func() {
		os.Unsetenv("LITELLM_BASE_URL")
		os.Unsetenv("LITELLM_API_KEY")
	}()

	p := NewProvider("ollama/llama3")
	if p == nil {
		t.Fatal("NewProvider returned nil")
	}
}

func TestNewProviderWithConfig(t *testing.T) {
	p := NewProviderWithConfig("http://custom:4000", "custom-key")
	if p == nil {
		t.Fatal("NewProviderWithConfig returned nil")
	}
	if p.Name() != "litellm" {
		t.Errorf("Name() = %q, want %q", p.Name(), "litellm")
	}
}

func TestNewProviderWithConfigEmpty(t *testing.T) {
	p := NewProviderWithConfig("", "")
	if p == nil {
		t.Fatal("NewProviderWithConfig returned nil")
	}
}

func TestGetModelProvider(t *testing.T) {
	tests := []struct {
		model      string
		wantPrefix string
	}{
		{"ollama/llama3", "ollama"},
		{"ollama/mistral", "ollama"},
		{"anthropic/claude-3-opus", "anthropic"},
		{"anthropic/claude-3-sonnet", "anthropic"},
		{"openai/gpt-4", "openai"},
		{"openai/gpt-3.5-turbo", "openai"},
		{"google/gemini-pro", "google"},
		{"gemini/gemini-1.5-pro", "google"},
		{"cohere/command-r", "cohere"},
		{"azure/gpt-4", "azure"},
		{"aws/bedrock/claude", "aws"},
		{"bedrock/anthropic/claude", "aws"},
		{"unknown/model", "litellm"},
		{"some-random-model", "litellm"},
		{"", "litellm"},
	}

	for _, tt := range tests {
		t.Run(tt.model, func(t *testing.T) {
			got := GetModelProvider(tt.model)
			if got != tt.wantPrefix {
				t.Errorf("GetModelProvider(%q) = %q, want %q", tt.model, got, tt.wantPrefix)
			}
		})
	}
}
