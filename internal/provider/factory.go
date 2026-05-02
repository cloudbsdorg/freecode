package provider

import (
	"os"
	"strings"
)

func NewProvider(model string) Provider {
	model = strings.ToLower(model)

	if strings.HasPrefix(model, "anthropic/") || strings.HasPrefix(model, "claude") {
		apiKey := os.Getenv("ANTHROPIC_API_KEY")
		return NewAnthropicProvider(apiKey)
	}

	if strings.HasPrefix(model, "openai/") || strings.HasPrefix(model, "gpt") || strings.HasPrefix(model, "o1") || strings.HasPrefix(model, "o3") {
		apiKey := os.Getenv("OPENAI_API_KEY")
		return NewOpenAIProvider(apiKey)
	}

	if strings.HasPrefix(model, "ollama/") {
		baseURL := os.Getenv("OLLAMA_BASE_URL")
		if baseURL == "" {
			baseURL = "http://localhost:11434"
		}
		return NewLiteLLMProvider(baseURL, "local")
	}

	if strings.HasPrefix(model, "google/") || strings.HasPrefix(model, "gemini/") {
		apiKey := os.Getenv("GOOGLE_API_KEY")
		return NewLiteLLMProvider("http://localhost:4000", apiKey)
	}

	if strings.HasPrefix(model, "cohere/") {
		apiKey := os.Getenv("COHERE_API_KEY")
		return NewLiteLLMProvider("http://localhost:4000", apiKey)
	}

	if strings.HasPrefix(model, "azure/") {
		apiKey := os.Getenv("AZURE_API_KEY")
		baseURL := os.Getenv("AZURE_BASE_URL")
		return NewLiteLLMProvider(baseURL, apiKey)
	}

	if strings.HasPrefix(model, "aws/") || strings.HasPrefix(model, "bedrock/") {
		return NewLiteLLMProvider("http://localhost:4000", "aws")
	}

	defaultBaseURL := os.Getenv("LITELLM_BASE_URL")
	if defaultBaseURL == "" {
		defaultBaseURL = "http://localhost:4000"
	}
	defaultAPIKey := os.Getenv("LITELLM_API_KEY")
	return NewLiteLLMProvider(defaultBaseURL, defaultAPIKey)
}

func NewProviderWithConfig(providerType, apiKey, baseURL string) Provider {
	switch strings.ToLower(providerType) {
	case "openai":
		if baseURL == "" {
			baseURL = "https://api.openai.com/v1"
		}
		p := NewOpenAIProvider(apiKey)
		p.BaseURL = baseURL
		return p
	case "anthropic":
		if baseURL == "" {
			baseURL = "https://api.anthropic.com/v1"
		}
		p := NewAnthropicProvider(apiKey)
		p.BaseURL = baseURL
		return p
	case "ollama":
		if baseURL == "" {
			baseURL = "http://localhost:11434"
		}
		return NewLiteLLMProvider(baseURL, "local")
	default:
		if baseURL == "" {
			baseURL = "http://localhost:4000"
		}
		return NewLiteLLMProvider(baseURL, apiKey)
	}
}

func GetModelProvider(model string) string {
	model = strings.ToLower(model)

	if strings.HasPrefix(model, "anthropic/") || strings.HasPrefix(model, "claude") {
		return "anthropic"
	}
	if strings.HasPrefix(model, "openai/") || strings.HasPrefix(model, "gpt") || strings.HasPrefix(model, "o1") || strings.HasPrefix(model, "o3") {
		return "openai"
	}
	if strings.HasPrefix(model, "ollama/") {
		return "ollama"
	}
	if strings.HasPrefix(model, "google/") || strings.HasPrefix(model, "gemini/") {
		return "google"
	}
	if strings.HasPrefix(model, "cohere/") {
		return "cohere"
	}
	if strings.HasPrefix(model, "azure/") {
		return "azure"
	}
	if strings.HasPrefix(model, "aws/") || strings.HasPrefix(model, "bedrock/") {
		return "aws"
	}

	return "unknown"
}
