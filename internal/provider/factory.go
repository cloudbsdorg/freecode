package provider

import (
	"os"
	"strings"
)

func NewProvider(model string) Provider {
	model = strings.ToLower(model)

	if strings.HasPrefix(model, "anthropic/") || strings.HasPrefix(model, "claude") {
		return NewAnthropicProvider(os.Getenv("ANTHROPIC_API_KEY"))
	}

	if strings.HasPrefix(model, "minimax/") || strings.HasPrefix(model, "minimax") {
		return NewMinimaxProvider(os.Getenv("MINIMAX_API_KEY"))
	}

	if strings.HasPrefix(model, "openai/") || strings.HasPrefix(model, "gpt") || strings.HasPrefix(model, "o1") || strings.HasPrefix(model, "o3") {
		return NewOpenAIProvider(os.Getenv("OPENAI_API_KEY"))
	}

	if strings.HasPrefix(model, "ollama/") {
		baseURL := os.Getenv("OLLAMA_BASE_URL")
		if baseURL == "" {
			baseURL = "http://localhost:11434"
		}
		return NewOllamaProvider(baseURL, os.Getenv("OLLAMA_API_KEY"))
	}

	if strings.HasPrefix(model, "groq/") {
		return NewGroqProvider(os.Getenv("GROQ_API_KEY"))
	}

	if strings.HasPrefix(model, "perplexity/") {
		return NewPerplexityProvider(os.Getenv("PERPLEXITY_API_KEY"))
	}

	if strings.HasPrefix(model, "google/") || strings.HasPrefix(model, "gemini/") {
		return NewGoogleProvider(os.Getenv("GOOGLE_API_KEY"))
	}

	if strings.HasPrefix(model, "cohere/") {
		return NewCohereProvider(os.Getenv("COHERE_API_KEY"))
	}

	if strings.HasPrefix(model, "mistral/") {
		return NewMistralProvider(os.Getenv("MISTRAL_API_KEY"))
	}

	if strings.HasPrefix(model, "togetherai/") || strings.HasPrefix(model, "together/") {
		return NewTogetherAIProvider(os.Getenv("TOGETHERAI_API_KEY"))
	}

	if strings.HasPrefix(model, "deepinfra/") {
		return NewDeepInfraProvider(os.Getenv("DEEPINFRA_API_KEY"))
	}

	if strings.HasPrefix(model, "cerebras/") {
		return NewCerebrasProvider(os.Getenv("CEREBRAS_API_KEY"))
	}

	if strings.HasPrefix(model, "xai/") {
		return NewxAIProvider(os.Getenv("XAI_API_KEY"))
	}

	if strings.HasPrefix(model, "alibaba/") || strings.HasPrefix(model, "qwen/") {
		return NewAlibabaProvider(os.Getenv("ALIBABA_API_KEY"))
	}

	if strings.HasPrefix(model, "huggingface/") {
		return NewHuggingFaceProvider(os.Getenv("HUGGINGFACE_API_KEY"))
	}

	if strings.HasPrefix(model, "deepseek/") {
		return NewDeepSeekProvider(os.Getenv("DEEPSEEK_API_KEY"))
	}

	if strings.HasPrefix(model, "fireworks/") {
		return NewFireworksProvider(os.Getenv("FIREWORKS_API_KEY"))
	}

	if strings.HasPrefix(model, "moonshot/") {
		return NewMoonshotProvider(os.Getenv("MOONSHOT_API_KEY"))
	}

	if strings.HasPrefix(model, "nebius/") {
		return NewNebiusProvider(os.Getenv("NEBIUS_API_KEY"))
	}

	if strings.HasPrefix(model, "openrouter/") {
		return NewOpenRouterProvider(os.Getenv("OPENROUTER_API_KEY"))
	}

	if strings.HasPrefix(model, "azure/") {
		return NewAzureProvider(os.Getenv("AZURE_API_KEY"), os.Getenv("AZURE_BASE_URL"))
	}

	if strings.HasPrefix(model, "aws/") || strings.HasPrefix(model, "bedrock/") {
		return NewBedrockProvider(
			os.Getenv("AWS_REGION"),
			os.Getenv("AWS_PROFILE"),
			os.Getenv("AWS_ENDPOINT"),
			os.Getenv("AWS_ACCESS_KEY_ID"),
			os.Getenv("AWS_SECRET_ACCESS_KEY"),
		)
	}

	if strings.HasPrefix(model, "gitlab/") {
		return NewGitLabProvider(os.Getenv("GITLAB_TOKEN"))
	}

	if strings.HasPrefix(model, "copilot/") || strings.HasPrefix(model, "github_copilot/") {
		return NewGitHubCopilotProvider(os.Getenv("GITHUB_COPILOT_TOKEN"))
	}

	if strings.HasPrefix(model, "vercel/") {
		return NewVercelProvider(os.Getenv("VERCEL_TOKEN"))
	}

	if strings.HasPrefix(model, "vertex/") {
		return NewVertexProvider(
			os.Getenv("VERTEX_PROJECT_ID"),
			os.Getenv("VERTEX_LOCATION"),
			os.Getenv("VERTEX_ACCESS_TOKEN"),
		)
	}

	if strings.HasPrefix(model, "venice/") {
		return NewVeniceProvider(os.Getenv("VENICE_API_KEY"))
	}

	if strings.HasPrefix(model, "zai/") {
		return NewZAIProvider(os.Getenv("ZAI_API_KEY"))
	}

	if strings.HasPrefix(model, "zenmux/") {
		return NewZenMuxProvider(os.Getenv("ZENMUX_API_KEY"))
	}

	if strings.HasPrefix(model, "baseten/") {
		return NewBasetenProvider(os.Getenv("BASETEN_API_KEY"))
	}

	if strings.HasPrefix(model, "cortecs/") {
		return NewCortecsProvider(os.Getenv("CORTECS_API_KEY"))
	}

	if strings.HasPrefix(model, "firmware/") {
		return NewFirmwareProvider(os.Getenv("FIRMWARE_API_KEY"))
	}

	if strings.HasPrefix(model, "ionet/") {
		return NewIonetProvider(os.Getenv("IONET_API_KEY"))
	}

	if strings.HasPrefix(model, "nvidia/") {
		return NewNVIDIAProvider(os.Getenv("NVIDIA_API_KEY"))
	}

	if strings.HasPrefix(model, "ollamacloud/") {
		return NewOllamaCloudProvider(os.Getenv("OLLAMA_CLOUD_API_KEY"))
	}

	if strings.HasPrefix(model, "cloudflare/") {
		return NewCloudflareGatewayProvider(
			os.Getenv("CLOUDFLARE_ACCOUNT_ID"),
			os.Getenv("CLOUDFLARE_GATEWAY_ID"),
			os.Getenv("CLOUDFLARE_API_TOKEN"),
		)
	}

	if strings.HasPrefix(model, "helicone/") {
		return NewHeliconeProvider(os.Getenv("HELICONE_API_KEY"))
	}

	if strings.HasPrefix(model, "llamacpp/") {
		return NewLlamaCppProvider(os.Getenv("LLAMACPP_BASE_URL"), os.Getenv("LLAMACPP_API_KEY"))
	}

	if strings.HasPrefix(model, "lmstudio/") {
		return NewLMStudioProvider(os.Getenv("LMSTUDIO_BASE_URL"), os.Getenv("LMSTUDIO_API_KEY"))
	}

	if strings.HasPrefix(model, "atomic/") {
		return NewAtomicChatProvider(os.Getenv("ATOMIC_CHAT_BASE_URL"), os.Getenv("ATOMIC_CHAT_API_KEY"))
	}

	if strings.HasPrefix(model, "302ai/") {
		return NewProvider302AI(os.Getenv("302AI_API_KEY"))
	}

	if strings.HasPrefix(model, "sap/") {
		return NewSAPAIProvider(os.Getenv("SAP_AI_CORE_SERVICE_KEY"))
	}

	if strings.HasPrefix(model, "stackit/") {
		return NewStackitProvider(os.Getenv("STACKIT_TOKEN"))
	}

	if strings.HasPrefix(model, "ovhcloud/") {
		return NewOVHcloudProvider(os.Getenv("OVHCLOUD_API_KEY"))
	}

	if strings.HasPrefix(model, "scaleway/") {
		return NewScalewayProvider(os.Getenv("SCALEWAY_API_KEY"))
	}

	defaultBaseURL := os.Getenv("LITELLM_BASE_URL")
	if defaultBaseURL == "" {
		defaultBaseURL = "http://localhost:4000"
	}
	return NewLiteLLMProvider(defaultBaseURL, os.Getenv("LITELLM_API_KEY"))
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
	case "minimax":
		if baseURL == "" {
			baseURL = "https://api.minimax.chat/v1"
		}
		p := NewMinimaxProvider(apiKey)
		p.BaseURL = baseURL
		return p
	case "ollama":
		if baseURL == "" {
			baseURL = "http://localhost:11434"
		}
		return NewOllamaProvider(baseURL, apiKey)
	case "groq":
		if baseURL == "" {
			baseURL = "https://api.groq.com/openai/v1"
		}
		p := NewGroqProvider(apiKey)
		p.BaseURL = baseURL
		return p
	case "perplexity":
		if baseURL == "" {
			baseURL = "https://api.perplexity.ai"
		}
		p := NewPerplexityProvider(apiKey)
		p.BaseURL = baseURL
		return p
	case "google", "gemini":
		p := NewGoogleProvider(apiKey)
		return p
	case "cohere":
		if baseURL == "" {
			baseURL = "https://api.cohere.ai/v1"
		}
		p := NewCohereProvider(apiKey)
		p.BaseURL = baseURL
		return p
	case "mistral":
		if baseURL == "" {
			baseURL = "https://api.mistral.ai/v1"
		}
		p := NewMistralProvider(apiKey)
		p.BaseURL = baseURL
		return p
	case "togetherai":
		if baseURL == "" {
			baseURL = "https://api.together.xyz/v1"
		}
		p := NewTogetherAIProvider(apiKey)
		p.BaseURL = baseURL
		return p
	case "deepinfra":
		if baseURL == "" {
			baseURL = "https://api.deepinfra.com/v1"
		}
		p := NewDeepInfraProvider(apiKey)
		p.BaseURL = baseURL
		return p
	case "cerebras":
		if baseURL == "" {
			baseURL = "https://api.cerebras.ai/v1"
		}
		p := NewCerebrasProvider(apiKey)
		p.BaseURL = baseURL
		return p
	case "xai":
		if baseURL == "" {
			baseURL = "https://api.x.ai/v1"
		}
		p := NewxAIProvider(apiKey)
		p.BaseURL = baseURL
		return p
	case "alibaba", "qwen":
		if baseURL == "" {
			baseURL = "https://dashscope.aliyuncs.com"
		}
		p := NewAlibabaProvider(apiKey)
		p.BaseURL = baseURL
		return p
	case "huggingface":
		if baseURL == "" {
			baseURL = "https://api.endpoints.huggingface.cloud/v1"
		}
		p := NewHuggingFaceProvider(apiKey)
		p.BaseURL = baseURL
		return p
	case "deepseek":
		if baseURL == "" {
			baseURL = "https://api.deepseek.com/v1"
		}
		p := NewDeepSeekProvider(apiKey)
		p.BaseURL = baseURL
		return p
	case "fireworks":
		if baseURL == "" {
			baseURL = "https://api.fireworks.ai/v1"
		}
		p := NewFireworksProvider(apiKey)
		p.BaseURL = baseURL
		return p
	case "moonshot":
		if baseURL == "" {
			baseURL = "https://api.moonshot.cn/v1"
		}
		p := NewMoonshotProvider(apiKey)
		p.BaseURL = baseURL
		return p
	case "nebius":
		if baseURL == "" {
			baseURL = "https://api.nebius.ai/v1"
		}
		p := NewNebiusProvider(apiKey)
		p.BaseURL = baseURL
		return p
	case "openrouter":
		if baseURL == "" {
			baseURL = "https://openrouter.ai/api/v1"
		}
		p := NewOpenRouterProvider(apiKey)
		p.BaseURL = baseURL
		return p
	case "azure":
		return NewAzureProvider(apiKey, baseURL)
	case "bedrock", "aws":
		return NewBedrockProvider(os.Getenv("AWS_REGION"), os.Getenv("AWS_PROFILE"), baseURL, apiKey, os.Getenv("AWS_SECRET_ACCESS_KEY"))
	case "vertex":
		return NewVertexProvider(os.Getenv("VERTEX_PROJECT_ID"), os.Getenv("VERTEX_LOCATION"), apiKey)
	case "gitlab":
		return NewGitLabProvider(apiKey)
	case "github_copilot":
		return NewGitHubCopilotProvider(apiKey)
	case "vercel":
		return NewVercelProvider(apiKey)
	case "venice":
		if baseURL == "" {
			baseURL = "https://api.venice.ai/api/v1"
		}
		p := NewVeniceProvider(apiKey)
		p.BaseURL = baseURL
		return p
	case "zai":
		if baseURL == "" {
			baseURL = "https://api.z-ai.ai/v1"
		}
		p := NewZAIProvider(apiKey)
		p.BaseURL = baseURL
		return p
	case "zenmux":
		if baseURL == "" {
			baseURL = "https://api.zenmux.ai/v1"
		}
		p := NewZenMuxProvider(apiKey)
		p.BaseURL = baseURL
		return p
	case "baseten":
		if baseURL == "" {
			baseURL = "https://app.baseten.co/v1"
		}
		p := NewBasetenProvider(apiKey)
		p.BaseURL = baseURL
		return p
	case "cortecs":
		if baseURL == "" {
			baseURL = "https://api.cortecs.ai/v1"
		}
		p := NewCortecsProvider(apiKey)
		p.BaseURL = baseURL
		return p
	case "firmware":
		if baseURL == "" {
			baseURL = "https://api.firmware.ai/v1"
		}
		p := NewFirmwareProvider(apiKey)
		p.BaseURL = baseURL
		return p
	case "ionet":
		if baseURL == "" {
			baseURL = "https://api.ionet.ai/v1"
		}
		p := NewIonetProvider(apiKey)
		p.BaseURL = baseURL
		return p
	case "nvidia":
		if baseURL == "" {
			baseURL = "https://ai.api.nvidia.com/v1"
		}
		p := NewNVIDIAProvider(apiKey)
		p.BaseURL = baseURL
		return p
	case "ollamacloud":
		if baseURL == "" {
			baseURL = "https://cloud.ollama.ai"
		}
		p := NewOllamaCloudProvider(apiKey)
		p.BaseURL = baseURL
		return p
	case "cloudflare_gateway":
		return NewCloudflareGatewayProvider(os.Getenv("CLOUDFLARE_ACCOUNT_ID"), os.Getenv("CLOUDFLARE_GATEWAY_ID"), apiKey)
	case "cloudflare_workers":
		return NewCloudflareWorkersProvider(os.Getenv("CLOUDFLARE_ACCOUNT_ID"), apiKey)
	case "helicone":
		if baseURL == "" {
			baseURL = "https://ai-gateway.helicone.ai"
		}
		p := NewHeliconeProvider(apiKey)
		p.BaseURL = baseURL
		return p
	case "llamacpp":
		if baseURL == "" {
			baseURL = "http://127.0.0.1:8080/v1"
		}
		return NewLlamaCppProvider(baseURL, apiKey)
	case "lmstudio":
		if baseURL == "" {
			baseURL = "http://127.0.0.1:1234/v1"
		}
		return NewLMStudioProvider(baseURL, apiKey)
	case "atomic_chat":
		if baseURL == "" {
			baseURL = "http://127.0.0.1:1337/v1"
		}
		return NewAtomicChatProvider(baseURL, apiKey)
	case "302ai":
		if baseURL == "" {
			baseURL = "https://api.302.ai/v1"
		}
		p := NewProvider302AI(apiKey)
		p.BaseURL = baseURL
		return p
	case "sap_ai_core":
		return NewSAPAIProvider(apiKey)
	case "stackit":
		return NewStackitProvider(apiKey)
	case "ovhcloud":
		if baseURL == "" {
			baseURL = "https://endpoints.ai.cloud.ovh.net/v1"
		}
		p := NewOVHcloudProvider(apiKey)
		p.BaseURL = baseURL
		return p
	case "scaleway":
		if baseURL == "" {
			baseURL = "https://api.scaleway.ai/v1"
		}
		p := NewScalewayProvider(apiKey)
		p.BaseURL = baseURL
		return p
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

	return "unknown"
}
