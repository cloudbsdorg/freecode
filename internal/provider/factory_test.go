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
	if p.Name() != "ollama" {
		t.Errorf("Name() = %q, want %q", p.Name(), "ollama")
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

func TestNewProvider_AnthropicByPrefix(t *testing.T) {
	os.Setenv("ANTHROPIC_API_KEY", "test-key")
	defer os.Unsetenv("ANTHROPIC_API_KEY")

	p := NewProvider("anthropic/claude-3")
	if p.Name() != "anthropic" {
		t.Errorf("Name() = %q, want %q", p.Name(), "anthropic")
	}
}

func TestNewProvider_AnthropicByName(t *testing.T) {
	os.Setenv("ANTHROPIC_API_KEY", "test-key")
	defer os.Unsetenv("ANTHROPIC_API_KEY")

	p := NewProvider("claude-3-opus")
	if p.Name() != "anthropic" {
		t.Errorf("Name() = %q, want %q", p.Name(), "anthropic")
	}
}

func TestNewProvider_MinimaxByPrefix(t *testing.T) {
	os.Setenv("MINIMAX_API_KEY", "test-key")
	defer os.Unsetenv("MINIMAX_API_KEY")

	p := NewProvider("minimax/agent")
	if p.Name() != "minimax" {
		t.Errorf("Name() = %q, want %q", p.Name(), "minimax")
	}
}

func TestNewProvider_OpenAIByPrefix(t *testing.T) {
	os.Setenv("OPENAI_API_KEY", "test-key")
	defer os.Unsetenv("OPENAI_API_KEY")

	p := NewProvider("openai/gpt-4")
	if p.Name() != "openai" {
		t.Errorf("Name() = %q, want %q", p.Name(), "openai")
	}
}

func TestNewProvider_OpenAIByGPT(t *testing.T) {
	os.Setenv("OPENAI_API_KEY", "test-key")
	defer os.Unsetenv("OPENAI_API_KEY")

	p := NewProvider("gpt-4-turbo")
	if p.Name() != "openai" {
		t.Errorf("Name() = %q, want %q", p.Name(), "openai")
	}
}

func TestNewProvider_OpenAIByO1(t *testing.T) {
	os.Setenv("OPENAI_API_KEY", "test-key")
	defer os.Unsetenv("OPENAI_API_KEY")

	p := NewProvider("o1-preview")
	if p.Name() != "openai" {
		t.Errorf("Name() = %q, want %q", p.Name(), "openai")
	}
}

func TestNewProvider_Google(t *testing.T) {
	os.Setenv("GOOGLE_API_KEY", "test-key")
	defer os.Unsetenv("GOOGLE_API_KEY")

	p := NewProvider("google/gemini-pro")
	if p.Name() != "litellm" {
		t.Errorf("Name() = %q, want %q", p.Name(), "litellm")
	}
}

func TestNewProvider_Cohere(t *testing.T) {
	os.Setenv("COHERE_API_KEY", "test-key")
	defer os.Unsetenv("COHERE_API_KEY")

	p := NewProvider("cohere/command")
	if p.Name() != "litellm" {
		t.Errorf("Name() = %q, want %q", p.Name(), "litellm")
	}
}

func TestNewProvider_Azure(t *testing.T) {
	os.Setenv("AZURE_API_KEY", "test-key")
	os.Setenv("AZURE_BASE_URL", "https://example.openai.azure.com")
	defer os.Unsetenv("AZURE_API_KEY")
	defer os.Unsetenv("AZURE_BASE_URL")

	p := NewProvider("azure/gpt-4")
	if p.Name() != "litellm" {
		t.Errorf("Name() = %q, want %q", p.Name(), "litellm")
	}
}

func TestNewProvider_AWS(t *testing.T) {
	p := NewProvider("aws/bedrock")
	if p.Name() != "litellm" {
		t.Errorf("Name() = %q, want %q", p.Name(), "litellm")
	}
}

func TestNewProvider_Bedrock(t *testing.T) {
	p := NewProvider("bedrock/anthropic")
	if p.Name() != "litellm" {
		t.Errorf("Name() = %q, want %q", p.Name(), "litellm")
	}
}

func TestNewProvider_DefaultLitellm(t *testing.T) {
	os.Setenv("LITELLM_BASE_URL", "http://custom:4000")
	os.Setenv("LITELLM_API_KEY", "test-litellm-key")
	defer func() {
		os.Unsetenv("LITELLM_BASE_URL")
		os.Unsetenv("LITELLM_API_KEY")
	}()

	p := NewProvider("some-random-model")
	if p.Name() != "litellm" {
		t.Errorf("Name() = %q, want %q", p.Name(), "litellm")
	}
}

func TestNewProviderWithConfig(t *testing.T) {
	p := NewProviderWithConfig("ollama", "custom-key", "http://custom:4000")
	if p == nil {
		t.Fatal("NewProviderWithConfig returned nil")
	}
	if p.Name() != "ollama" {
		t.Errorf("Name() = %q, want %q", p.Name(), "ollama")
	}
}

func TestNewProviderWithConfigEmpty(t *testing.T) {
	p := NewProviderWithConfig("ollama", "", "")
	if p == nil {
		t.Fatal("NewProviderWithConfig returned nil")
	}
}

func TestNewProviderWithConfigOpenAI(t *testing.T) {
	p := NewProviderWithConfig("openai", "test-key", "")
	if p == nil {
		t.Fatal("NewProviderWithConfig returned nil")
	}
	if p.Name() != "openai" {
		t.Errorf("Name() = %q, want %q", p.Name(), "openai")
	}
}

func TestNewProviderWithConfigOpenAIWithBaseURL(t *testing.T) {
	p := NewProviderWithConfig("openai", "test-key", "https://custom.openai.com/v1")
	if p == nil {
		t.Fatal("NewProviderWithConfig returned nil")
	}
	if p.Name() != "openai" {
		t.Errorf("Name() = %q, want %q", p.Name(), "openai")
	}
}

func TestNewProviderWithConfigAnthropic(t *testing.T) {
	p := NewProviderWithConfig("anthropic", "test-key", "")
	if p == nil {
		t.Fatal("NewProviderWithConfig returned nil")
	}
	if p.Name() != "anthropic" {
		t.Errorf("Name() = %q, want %q", p.Name(), "anthropic")
	}
}

func TestNewProviderWithConfigAnthropicWithBaseURL(t *testing.T) {
	p := NewProviderWithConfig("anthropic", "test-key", "https://custom.anthropic.com/v1")
	if p == nil {
		t.Fatal("NewProviderWithConfig returned nil")
	}
	if p.Name() != "anthropic" {
		t.Errorf("Name() = %q, want %q", p.Name(), "anthropic")
	}
}

func TestNewProviderWithConfigMinimax(t *testing.T) {
	p := NewProviderWithConfig("minimax", "test-key", "")
	if p == nil {
		t.Fatal("NewProviderWithConfig returned nil")
	}
	if p.Name() != "minimax" {
		t.Errorf("Name() = %q, want %q", p.Name(), "minimax")
	}
}

func TestNewProviderWithConfigMinimaxWithBaseURL(t *testing.T) {
	p := NewProviderWithConfig("minimax", "test-key", "https://custom.minimax.com/v1")
	if p == nil {
		t.Fatal("NewProviderWithConfig returned nil")
	}
	if p.Name() != "minimax" {
		t.Errorf("Name() = %q, want %q", p.Name(), "minimax")
	}
}

func TestNewProviderWithConfigOllama(t *testing.T) {
	p := NewProviderWithConfig("ollama", "", "")
	if p == nil {
		t.Fatal("NewProviderWithConfig returned nil")
	}
	if p.Name() != "ollama" {
		t.Errorf("Name() = %q, want %q", p.Name(), "ollama")
	}
}

func TestNewProviderWithConfigUnknown(t *testing.T) {
	p := NewProviderWithConfig("unknown", "test-key", "http://localhost:4000")
	if p == nil {
		t.Fatal("NewProviderWithConfig returned nil")
	}
	if p.Name() != "litellm" {
		t.Errorf("Name() = %q, want %q", p.Name(), "litellm")
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
		{"unknown/model", "unknown"},
		{"some-random-model", "unknown"},
		{"", "unknown"},
		{"claude-3-opus", "anthropic"},
		{"gpt-4", "openai"},
		{"minimax/MiniMax-M2.7", "minimax"},
		{"MiniMax-M2.7", "minimax"},
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
