package provider

import (
	"os"
	"strings"
	"sync"
)

type connectorType int

const (
	connectorOpenAI connectorType = iota
	connectorOpenAICompatible
	connectorAnthropic
	connectorAzure
	connectorBedrock
	connectorOllama
	connectorGoogle
	connectorGoogleVertex
	connectorUnknown
)

var connectorNames = map[connectorType]string{
	connectorOpenAI:           "@ai-sdk/openai",
	connectorOpenAICompatible: "@ai-sdk/openai-compatible",
	connectorAnthropic:        "@ai-sdk/anthropic",
	connectorAzure:            "@ai-sdk/azure",
	connectorBedrock:          "@ai-sdk/amazon-bedrock",
	connectorOllama:           "ollama",
	connectorGoogle:           "@ai-sdk/google",
	connectorGoogleVertex:     "@ai-sdk/google-vertex",
}

var (
	connectorHandlers     map[string]connectorType
	connectorHandlersOnce  sync.Once
)

func initConnectorHandlers() {
	connectorHandlersOnce.Do(func() {
		connectorHandlers = map[string]connectorType{
			"@ai-sdk/openai":              connectorOpenAI,
			"@ai-sdk/openai-compatible":    connectorOpenAICompatible,
			"@ai-sdk/anthropic":            connectorAnthropic,
			"@ai-sdk/azure":                connectorAzure,
			"@ai-sdk/amazon-bedrock":       connectorBedrock,
			"@ai-sdk/google":               connectorGoogle,
			"@ai-sdk/google-vertex":         connectorGoogleVertex,
			"@ai-sdk/cohere":               connectorOpenAICompatible,
			"@ai-sdk/groq":                 connectorOpenAICompatible,
			"@ai-sdk/mistral":              connectorOpenAICompatible,
			"@ai-sdk/perplexity":           connectorOpenAICompatible,
			"@ai-sdk/togetherai":           connectorOpenAICompatible,
			"@ai-sdk/xai":                  connectorOpenAICompatible,
			"@ai-sdk/deepinfra":            connectorOpenAICompatible,
			"@ai-sdk/cerebras":             connectorOpenAICompatible,
			"ollama":                       connectorOllama,
			"@openrouter/ai-sdk-provider":  connectorOpenAICompatible,
			"@jerome-benoit/sap-ai-provider-v2": connectorOpenAICompatible,
		}
	})
}

func getConnectorType(npm string) connectorType {
	initConnectorHandlers()
	if ct, ok := connectorHandlers[npm]; ok {
		return ct
	}
	if strings.Contains(npm, "openai") || strings.Contains(npm, "compatible") {
		return connectorOpenAICompatible
	}
	if npm == "" {
		return connectorOpenAICompatible
	}
	return connectorOpenAICompatible
}

type providerCreator struct {
	envVar string
	create func(providerID, baseURL, apiKey string) Provider
}

var (
	providerCreators     map[connectorType]*providerCreator
	providerCreatorsOnce sync.Once
)

func initProviderCreators() {
	providerCreatorsOnce.Do(func() {
		providerCreators = map[connectorType]*providerCreator{
			connectorOpenAI: {
				envVar: "OPENAI_API_KEY",
				create: func(id, baseURL, apiKey string) Provider {
					if baseURL == "" {
						baseURL = "https://api.openai.com/v1"
					}
					p := NewOpenAIProvider(apiKey)
					p.BaseURL = baseURL
					return p
				},
			},
			connectorOpenAICompatible: {
				envVar: "",
				create: func(id, baseURL, apiKey string) Provider {
					if baseURL == "" {
						baseURL = "http://localhost:4000"
					}
					return NewOpenAICompatibleProvider(id, baseURL, apiKey)
				},
			},
			connectorAnthropic: {
				envVar: "ANTHROPIC_API_KEY",
				create: func(id, baseURL, apiKey string) Provider {
					if baseURL == "" {
						baseURL = "https://api.anthropic.com"
					}
					p := NewAnthropicProvider(apiKey)
					p.BaseURL = baseURL
					return p
				},
			},
			connectorAzure: {
				envVar: "AZURE_API_KEY",
				create: func(id, baseURL, apiKey string) Provider {
					if baseURL == "" {
						baseURL = os.Getenv("AZURE_BASE_URL")
					}
					return NewAzureProvider(apiKey, baseURL)
				},
			},
			connectorBedrock: {
				envVar: "AWS_ACCESS_KEY_ID",
				create: func(id, baseURL, apiKey string) Provider {
					return NewBedrockProvider(
						os.Getenv("AWS_REGION"),
						os.Getenv("AWS_PROFILE"),
						baseURL,
						apiKey,
						os.Getenv("AWS_SECRET_ACCESS_KEY"),
					)
				},
			},
			connectorOllama: {
				envVar: "",
				create: func(id, baseURL, apiKey string) Provider {
					if baseURL == "" {
						baseURL = os.Getenv("OLLAMA_BASE_URL")
					}
					if baseURL == "" {
						baseURL = "http://localhost:11434"
					}
					return NewOllamaProvider(baseURL, apiKey)
				},
			},
			connectorGoogle: {
				envVar: "GOOGLE_API_KEY",
				create: func(id, baseURL, apiKey string) Provider {
					return NewGoogleProvider(apiKey)
				},
			},
			connectorGoogleVertex: {
				envVar: "VERTEX_ACCESS_TOKEN",
				create: func(id, baseURL, apiKey string) Provider {
					return NewVertexProvider(
						os.Getenv("VERTEX_PROJECT_ID"),
						os.Getenv("VERTEX_LOCATION"),
						apiKey,
					)
				},
			},
			connectorUnknown: {
				envVar: "",
				create: func(id, baseURL, apiKey string) Provider {
					if baseURL == "" {
						baseURL = "http://localhost:4000"
					}
					return NewLiteLLMProvider(baseURL, apiKey)
				},
			},
		}
	})
}

func NewProvider(model string) Provider {
	modelLower := strings.ToLower(model)

	if strings.HasPrefix(modelLower, "claude") {
		return NewAnthropicProvider(os.Getenv("ANTHROPIC_API_KEY"))
	}
	if strings.HasPrefix(modelLower, "gpt") || strings.HasPrefix(modelLower, "o1") || strings.HasPrefix(modelLower, "o3") {
		return NewOpenAIProvider(os.Getenv("OPENAI_API_KEY"))
	}
	if strings.HasPrefix(modelLower, "minimax") {
		return NewMinimaxProvider(os.Getenv("MINIMAX_API_KEY"))
	}

	parts := strings.SplitN(model, "/", 2)
	if len(parts) < 2 {
		return NewLiteLLMProvider(os.Getenv("LITELLM_BASE_URL"), os.Getenv("LITELLM_API_KEY"))
	}
	providerID := strings.ToLower(parts[0])

	providerAliases := map[string]string{
		"aws": "bedrock",
	}

	if actualID, ok := providerAliases[providerID]; ok {
		providerID = actualID
	}

	ct := connectorOpenAICompatible
	envVar := ""
	baseURL := ""

	info, ok := GetProviderInfo(providerID)
	if ok {
		ct = getConnectorType(info.ConnectorType)
		if len(info.EnvVars) > 0 {
			envVar = info.EnvVars[0]
		}
		baseURL = info.BaseURL
	}

	if ct == connectorOpenAICompatible && envVar == "" {
		envVar = strings.ToUpper(providerID) + "_API_KEY"
	}

	if baseURL == "" {
		baseURL = os.Getenv(strings.ToUpper(providerID) + "_BASE_URL")
	}

	initProviderCreators()
	creator, ok := providerCreators[ct]
	if !ok {
		creator = providerCreators[connectorOpenAICompatible]
	}

	apiKey := os.Getenv(envVar)
	return creator.create(providerID, baseURL, apiKey)
}

func NewProviderWithConfig(providerType, apiKey, baseURL string) Provider {
	providerType = strings.ToLower(providerType)

	info, ok := GetProviderInfo(providerType)
	ct := connectorOpenAICompatible
	if ok {
		ct = getConnectorType(info.ConnectorType)
	}

	switch providerType {
	case "openai":
		ct = connectorOpenAI
	case "anthropic":
		ct = connectorAnthropic
	case "azure":
		ct = connectorAzure
	case "bedrock", "aws":
		ct = connectorBedrock
	case "ollama":
		ct = connectorOllama
	case "google", "gemini":
		ct = connectorGoogle
	case "vertex":
		ct = connectorGoogleVertex
	case "litellm", "unknown":
		ct = connectorUnknown
	}

	initProviderCreators()
	creator, ok := providerCreators[ct]
	if !ok {
		creator = providerCreators[connectorUnknown]
	}

	if baseURL == "" && ok && info.BaseURL != "" {
		baseURL = info.BaseURL
	}

	return creator.create(providerType, baseURL, apiKey)
}

func GetModelProvider(model string) string {
	model = strings.ToLower(model)

	if strings.HasPrefix(model, "anthropic/") || strings.HasPrefix(model, "claude") {
		return "anthropic"
	}
	if strings.HasPrefix(model, "minimax/") || strings.HasPrefix(model, "minimax") {
		return "minimax"
	}
	if strings.HasPrefix(model, "openai/") || strings.HasPrefix(model, "gpt") || strings.HasPrefix(model, "o1") || strings.HasPrefix(model, "o3") {
		return "openai"
	}
	if strings.HasPrefix(model, "ollama/") {
		return "ollama"
	}
	if strings.HasPrefix(model, "groq/") {
		return "groq"
	}
	if strings.HasPrefix(model, "perplexity/") {
		return "perplexity"
	}
	if strings.HasPrefix(model, "google/") || strings.HasPrefix(model, "gemini/") {
		return "google"
	}
	if strings.HasPrefix(model, "cohere/") {
		return "cohere"
	}
	if strings.HasPrefix(model, "mistral/") {
		return "mistral"
	}
	if strings.HasPrefix(model, "togetherai/") || strings.HasPrefix(model, "together/") {
		return "togetherai"
	}
	if strings.HasPrefix(model, "deepinfra/") {
		return "deepinfra"
	}
	if strings.HasPrefix(model, "cerebras/") {
		return "cerebras"
	}
	if strings.HasPrefix(model, "xai/") {
		return "xai"
	}
	if strings.HasPrefix(model, "alibaba/") || strings.HasPrefix(model, "qwen/") {
		return "alibaba"
	}
	if strings.HasPrefix(model, "huggingface/") {
		return "huggingface"
	}
	if strings.HasPrefix(model, "deepseek/") {
		return "deepseek"
	}
	if strings.HasPrefix(model, "fireworks/") {
		return "fireworks"
	}
	if strings.HasPrefix(model, "moonshot/") {
		return "moonshot"
	}
	if strings.HasPrefix(model, "nebius/") {
		return "nebius"
	}
	if strings.HasPrefix(model, "openrouter/") {
		return "openrouter"
	}
	if strings.HasPrefix(model, "azure/") {
		return "azure"
	}
	if strings.HasPrefix(model, "aws/") || strings.HasPrefix(model, "bedrock/") {
		return "bedrock"
	}
	if strings.HasPrefix(model, "gitlab/") {
		return "gitlab"
	}
	if strings.HasPrefix(model, "copilot/") || strings.HasPrefix(model, "github_copilot/") {
		return "github_copilot"
	}
	if strings.HasPrefix(model, "vercel/") {
		return "vercel"
	}
	if strings.HasPrefix(model, "vertex/") {
		return "vertex"
	}
	if strings.HasPrefix(model, "venice/") {
		return "venice"
	}
	if strings.HasPrefix(model, "zai/") {
		return "zai"
	}
	if strings.HasPrefix(model, "zenmux/") {
		return "zenmux"
	}
	if strings.HasPrefix(model, "baseten/") {
		return "baseten"
	}
	if strings.HasPrefix(model, "cortecs/") {
		return "cortecs"
	}
	if strings.HasPrefix(model, "firmware/") {
		return "firmware"
	}
	if strings.HasPrefix(model, "ionet/") {
		return "ionet"
	}
	if strings.HasPrefix(model, "nvidia/") {
		return "nvidia"
	}
	if strings.HasPrefix(model, "ollamacloud/") {
		return "ollamacloud"
	}
	if strings.HasPrefix(model, "cloudflare/") {
		return "cloudflare_gateway"
	}
	if strings.HasPrefix(model, "helicone/") {
		return "helicone"
	}
	if strings.HasPrefix(model, "llamacpp/") {
		return "llamacpp"
	}
	if strings.HasPrefix(model, "lmstudio/") {
		return "lmstudio"
	}
	if strings.HasPrefix(model, "atomic/") {
		return "atomic_chat"
	}
	if strings.HasPrefix(model, "302ai/") {
		return "302ai"
	}
	if strings.HasPrefix(model, "sap/") {
		return "sap_ai_core"
	}
	if strings.HasPrefix(model, "stackit/") {
		return "stackit"
	}
	if strings.HasPrefix(model, "ovhcloud/") {
		return "ovhcloud"
	}
	if strings.HasPrefix(model, "scaleway/") {
		return "scaleway"
	}

	parts := strings.SplitN(model, "/", 2)
	if len(parts) >= 2 {
		return parts[0]
	}

	return "unknown"
}
