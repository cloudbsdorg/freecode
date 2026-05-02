package provider

import (
	"os"
	"strings"
)

func NewProvider(model string) Provider {
	baseURL := os.Getenv("LITELLM_BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:4000"
	}

	apiKey := os.Getenv("LITELLM_API_KEY")
	if apiKey == "" {
		apiKey = "local"
	}

	return NewLiteLLMProvider(baseURL, apiKey)
}

func NewProviderWithConfig(baseURL, apiKey string) Provider {
	if baseURL == "" {
		baseURL = "http://localhost:4000"
	}
	if apiKey == "" {
		apiKey = "local"
	}
	return NewLiteLLMProvider(baseURL, apiKey)
}

func GetModelProvider(model string) string {
	model = strings.ToLower(model)

	if strings.HasPrefix(model, "ollama/") {
		return "ollama"
	}
	if strings.HasPrefix(model, "anthropic/") {
		return "anthropic"
	}
	if strings.HasPrefix(model, "openai/") {
		return "openai"
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

	return "litellm"
}
