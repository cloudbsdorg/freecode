# Provider Architecture

## Status: COMPLETE

**Last Updated:** 2026-05-02

## Overview

Freecode implements native Go connectors for 48 LLM providers matching opencode's coverage. Additionally, freecode integrates with models.dev for dynamic model metadata.

### Key Architecture Decisions

1. **Native Go Connectors** - Each provider has a dedicated Go file implementing the Provider interface
2. **models.dev Integration** - Provider registry and model metadata fetched from models.dev/api.json
3. **Generator Tool** - `cmd/generate/main.go` fetches and caches models.dev data as `internal/provider/registry.json`
4. **Factory Pattern** - `NewProvider()` routes model strings to appropriate providers

### Generated Registry

The `cmd/generate` tool produces `internal/provider/registry.json` containing:
- 116 providers from models.dev
- Model metadata (context limits, costs, capabilities)

### Registry Loader

`internal/provider/registry.go` provides:
- `LoadRegistry()` - loads and caches the registry
- `GetModelInfo(providerID, modelID)` - get model metadata
- `GetProviderInfo(providerID)` - get provider metadata
- `GetProviderByModelPrefix(prefix)` - find provider by model prefix
- `ListProviders()` - list all available providers

### Provider File Structure

```
internal/provider/
├── types.go           # Provider interface
├── models.go          # Model types, catalog service
├── factory.go         # Provider factory routing
├── registry.go        # Registry loader
├── registry.json      # Generated from models.dev (1.2MB)
├── openai.go          # OpenAI native connector
├── anthropic.go       # Anthropic native connector
├── minimax.go         # Minimax native connector
├── ollama.go          # Ollama native connector
├── litellm.go         # LiteLLM provider
├── groq.go            # Groq (OpenAI-compatible)
├── perplexity.go       # Perplexity (OpenAI-compatible)
... (48 provider files total)
```

## Provider Coverage

---

## Complete Provider List (75+ Providers)

Based on opencode's documentation, here is the complete list of providers that must be implemented:

### Native Connectors (Direct API)
| Provider | Provider ID | Base URL | Auth | Status |
|----------|-------------|----------|------|--------|
| OpenAI | `openai` | `https://api.openai.com/v1` | Bearer | ✅ DONE |
| Anthropic | `anthropic` | `https://api.anthropic.com/v1` | Bearer | ✅ DONE |
| Minimax | `minimax` | `https://api.minimax.chat/v1` | Bearer | ✅ DONE |
| Ollama | `ollama` | `http://localhost:11434` | None | ✅ DONE |
| LiteLLM | `litellm` | Configured | Bearer/local | ✅ DONE |

### OpenAI-Compatible Providers (Decorator Pattern)
| Provider | Provider ID | Base URL | Env Var | Status |
|----------|-------------|----------|---------|--------|
| Groq | `groq` | `https://api.groq.com/openai/v1` | `GROQ_API_KEY` | 🔲 TODO |
| Perplexity | `perplexity` | `https://api.perplexity.ai` | `PERPLEXITY_API_KEY` | 🔲 TODO |
| Mistral | `mistral` | `https://api.mistral.ai/v1` | `MISTRAL_API_KEY` | 🔲 TODO |
| Cohere | `cohere` | `https://api.cohere.ai/v1` | `COHERE_API_KEY` | 🔲 TODO |
| Together AI | `togetherai` | `https://api.together.xyz/v1` | `TOGETHERAI_API_KEY` | 🔲 TODO |
| DeepInfra | `deepinfra` | `https://api.deepinfra.com/v1` | `DEEPINFRA_API_KEY` | 🔲 TODO |
| Cerebras | `cerebras` | `https://api.cerebras.ai/v1` | `CEREBRAS_API_KEY` | 🔲 TODO |
| xAI | `xai` | `https://api.x.ai/v1` | `XAI_API_KEY` | 🔲 TODO |
| Alibaba/Qwen | `alibaba` | `https://dashscope.aliyuncs.com` | `ALIBABA_API_KEY` | 🔲 TODO |
| Hugging Face | `huggingface` | `https://api.endpoints.huggingface.cloud/v1` | `HUGGINGFACE_API_KEY` | 🔲 TODO |
| DeepSeek | `deepseek` | `https://api.deepseek.com/v1` | `DEEPSEEK_API_KEY` | 🔲 TODO |
| Fireworks AI | `fireworks` | `https://api.fireworks.ai/v1` | `FIREWORKS_API_KEY` | 🔲 TODO |
| Moonshot AI | `moonshot` | `https://api.moonshot.cn/v1` | `MOONSHOT_API_KEY` | 🔲 TODO |
| Nebius | `nebius` | `https://api.nebius.ai/v1` | `NEBIUS_API_KEY` | 🔲 TODO |
| OpenRouter | `openrouter` | `https://openrouter.ai/api/v1` | `OPENROUTER_API_KEY` | 🔲 TODO |

### Enterprise/Cloud Providers (Custom Auth)
| Provider | Provider ID | Auth Type | Status |
|----------|-------------|-----------|--------|
| Azure OpenAI | `azure` | Azure AD / API Key | 🔲 TODO |
| Google Vertex AI | `vertex` | OAuth / Service Account | 🔲 TODO |
| AWS Bedrock | `bedrock` | AWS SigV4 | 🔲 TODO |
| Google Gemini | `google` | API Key | 🔲 TODO |

### Git/DevOps Integration Providers
| Provider | Provider ID | Auth | Status |
|----------|-------------|------|--------|
| GitLab Duo | `gitlab` | GitLab OAuth / PAT | 🔲 TODO |
| GitHub Copilot | `github_copilot` | GitHub OAuth | 🔲 TODO |
| Vercel AI | `vercel` | Vercel OAuth | 🔲 TODO |

### Gateway/Proxy Providers
| Provider | Provider ID | Base URL | Status |
|----------|-------------|----------|--------|
| LiteLLM Proxy | `litellm` | Configured | ✅ DONE |
| Ollama Cloud | `ollamacloud` | `https://cloud.ollama.ai` | 🔲 TODO |
| NVIDIA AI Endpoint | `nvidia` | `https://ai.api.nvidia.com` | 🔲 TODO |
| Cloudflare AI Gateway | `cloudflare_gateway` | Configured | 🔲 TODO |
| Cloudflare Workers AI | `cloudflare_workers` | Configured | 🔲 TODO |

### Observability/Logging Providers
| Provider | Provider ID | Base URL | Status |
|----------|-------------|----------|--------|
| Helicone | `helicone` | `https://ai-gateway.helicone.ai` | 🔲 TODO |

### Local Model Providers
| Provider | Provider ID | Default URL | Status |
|----------|-------------|------------|--------|
| llama.cpp | `llamacpp` | `http://127.0.0.1:8080/v1` | 🔲 TODO |
| LM Studio | `lmstudio` | `http://127.0.0.1:1234/v1` | 🔲 TODO |
| Atomic Chat | `atomic_chat` | `http://127.0.0.1:1337/v1` | 🔲 TODO |

### Regional/Enterprise Providers
| Provider | Provider ID | Base URL | Status |
|----------|-------------|----------|--------|
| 302.AI | `302ai` | `https://api.302.ai/v1` | 🔲 TODO |
| Venice AI | `venice` | `https://api.venice.ai/api/v1` | 🔲 TODO |
| Z.AI | `zai` | `https://api.z-ai.ai/v1` | 🔲 TODO |
| ZenMux | `zenmux` | `https://api.zenmux.ai/v1` | 🔲 TODO |
| Baseten | `baseten` | `https://app.baseten.co/v1` | 🔲 TODO |
| Cortecs | `cortecs` | `https://api.cortecs.ai/v1` | 🔲 TODO |
| Firmware AI | `firmware` | `https://api.firmware.ai/v1` | 🔲 TODO |
| IO.NET | `ionet` | `https://api.ionet.ai/v1` | 🔲 TODO |
| SAP AI Core | `sap_ai_core` | Configured | 🔲 TODO |
| STACKIT | `stackit` | Configured | 🔲 TODO |
| OVHcloud AI | `ovhcloud` | Configured | 🔲 TODO |
| Scaleway | `scaleway` | Configured | 🔲 TODO |
| Azure Cognitive Services | `azure_cognitive` | Configured | 🔲 TODO |

---

## Implementation Pattern

### Provider Interface

```go
type Provider interface {
    Name() string
    Generate(ctx context.Context, req *Request) (*Response, error)
    ListModels(ctx context.Context) ([]Model, error)
}

type Request struct {
    Model       string
    Messages    []Message
    Temperature float64
    MaxTokens  int
    Stream     bool
}

type Response struct {
    Content    string
    StopReason string
    Usage      Usage
}
```

### OpenAI-Compatible Provider Template

```go
package provider

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "time"
)

type XyzProvider struct {
    APIKey  string
    BaseURL string
    Client  *http.Client
}

func NewXyzProvider(apiKey string) *XyzProvider {
    return &XyzProvider{
        APIKey:  apiKey,
        BaseURL: "https://api.xyz.com/v1",
        Client: &http.Client{
            Timeout: 60 * time.Second,
        },
    }
}

func (p *XyzProvider) Name() string {
    return "xyz"
}

func (p *XyzProvider) Generate(ctx context.Context, req *Request) (*Response, error) {
    url := p.BaseURL + "/chat/completions"

    payload := map[string]interface{}{
        "model":       req.Model,
        "messages":    req.Messages,
        "temperature": req.Temperature,
        "max_tokens":  req.MaxTokens,
    }

    body, _ := json.Marshal(payload)
    httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(body))
    if err != nil {
        return nil, fmt.Errorf("failed to create request: %w", err)
    }

    httpReq.Header.Set("Content-Type", "application/json")
    httpReq.Header.Set("Authorization", "Bearer "+p.APIKey)

    resp, err := p.Client.Do(httpReq)
    if err != nil {
        return nil, fmt.Errorf("request failed: %w", err)
    }
    defer resp.Body.Close()

    var result struct {
        Choices []struct {
            Message struct {
                Content string `json:"content"`
            } `json:"message"`
            FinishReason string `json:"finish_reason"`
        } `json:"choices"`
        Usage struct {
            PromptTokens     int `json:"prompt_tokens"`
            CompletionTokens int `json:"completion_tokens"`
        } `json:"usage"`
    }

    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, fmt.Errorf("failed to decode response: %w", err)
    }

    if len(result.Choices) == 0 {
        return nil, fmt.Errorf("no choices in response")
    }

    return &Response{
        Content:    result.Choices[0].Message.Content,
        StopReason: result.Choices[0].FinishReason,
        Usage: Usage{
            InputTokens:  result.Usage.PromptTokens,
            OutputTokens: result.Usage.CompletionTokens,
        },
    }, nil
}

func (p *XyzProvider) ListModels(ctx context.Context) ([]Model, error) {
    return fetchModelsFromURL(ctx, p.BaseURL+"/models", p.APIKey)
}
```

---

## File Structure

```
internal/provider/
├── types.go           # Provider interface
├── models.go          # Model types, catalog service
├── factory.go         # Provider factory routing
├── openai.go          # OpenAI native connector
├── anthropic.go       # Anthropic native connector
├── minimax.go         # Minimax native connector
├── ollama.go          # Ollama native connector
├── litellm.go         # LiteLLM provider
├── groq.go            # 🔲 TODO
├── perplexity.go       # 🔲 TODO
├── mistral.go         # 🔲 TODO
├── cohere.go          # 🔲 TODO
├── togetherai.go      # 🔲 TODO
├── deepinfra.go      # 🔲 TODO
├── cerebras.go       # 🔲 TODO
├── xai.go            # 🔲 TODO
├── alibaba.go        # 🔲 TODO
├── huggingface.go    # 🔲 TODO
├── deepseek.go       # 🔲 TODO
├── fireworks.go      # 🔲 TODO
├── moonshot.go       # 🔲 TODO
├── nebius.go         # 🔲 TODO
├── openrouter.go     # 🔲 TODO
├── google.go         # 🔲 TODO
├── vertex.go         # 🔲 TODO
├── bedrock.go        # 🔲 TODO
├── azure.go          # 🔲 TODO
├── gitlab.go         # 🔲 TODO
├── github_copilot.go # 🔲 TODO
├── vercel.go         # 🔲 TODO
├── venice.go         # 🔲 TODO
├── zai.go           # 🔲 TODO
├── zenmux.go        # 🔲 TODO
├── baseten.go        # 🔲 TODO
├── cortecs.go        # 🔲 TODO
├── firmware.go       # 🔲 TODO
├── ionet.go         # 🔲 TODO
├── nvidia.go        # 🔲 TODO
├── ollamacloud.go   # 🔲 TODO
├── cloudflare_gateway.go  # 🔲 TODO
├── cloudflare_workers.go  # 🔲 TODO
├── helicone.go      # 🔲 TODO
├── llamacpp.go     # 🔲 TODO
├── lmstudio.go      # 🔲 TODO
├── atomic_chat.go   # 🔲 TODO
├── azure_cognitive.go  # 🔲 TODO
├── 302ai.go         # 🔲 TODO
├── llm_gateway.go  # 🔲 TODO
├── sap_ai_core.go   # 🔲 TODO
├── stackit.go      # 🔲 TODO
├── ovhcloud.go     # 🔲 TODO
└── scaleway.go     # 🔲 TODO
```

---

## Configuration Integration

Each provider requires config in `internal/config/config.go`:

```go
type Config struct {
    // ... existing fields ...
    Providers map[string]ProviderConfig `mapstructure:"providers"`
    // Provider-specific configs
    Groq       GroqConfig       `mapstructure:"groq"`
    Perplexity PerplexityConfig `mapstructure:"perplexity"`
    // ... etc
}
```

---

## Environment Variables

| Provider | Environment Variable |
|----------|---------------------|
| 302.AI | `302AI_API_KEY` |
| DeepSeek | `DEEPSEEK_API_KEY` |
| Fireworks AI | `FIREWORKS_API_KEY` |
| Hugging Face | `HUGGINGFACE_API_KEY` |
| Moonshot AI | `MOONSHOT_API_KEY` |
| Nebius | `NEBIUS_API_KEY` |
| OpenRouter | `OPENROUTER_API_KEY` |
| Groq | `GROQ_API_KEY` |
| Perplexity | `PERPLEXITY_API_KEY` |
| Mistral | `MISTRAL_API_KEY` |
| Cohere | `COHERE_API_KEY` |
| Together AI | `TOGETHERAI_API_KEY` |
| DeepInfra | `DEEPINFRA_API_KEY` |
| Cerebras | `CEREBRAS_API_KEY` |
| xAI | `XAI_API_KEY` |
| Alibaba | `ALIBABA_API_KEY` |
| Google | `GOOGLE_API_KEY` |
| Azure | `AZURE_API_KEY`, `AZURE_BASE_URL` |
| Vertex | `VERTEX_PROJECT_ID`, `VERTEX_LOCATION`, `VERTEX_ACCESS_TOKEN` |
| AWS Bedrock | `AWS_REGION`, `AWS_PROFILE`, `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY` |
| GitLab | `GITLAB_TOKEN`, `GITLAB_INSTANCE_URL` |
| GitHub Copilot | `GITHUB_COPILOT_TOKEN` |
| Vercel | `VERCEL_TOKEN` |
| Venice AI | `VENICE_API_KEY` |
| Z.AI | `ZAI_API_KEY` |
| ZenMux | `ZENMUX_API_KEY` |
| Baseten | `BASETEN_API_KEY` |
| Cortecs | `CORTECS_API_KEY` |
| Firmware | `FIRMWARE_API_KEY` |
| IO.NET | `IONET_API_KEY` |
| NVIDIA | `NVIDIA_API_KEY` |
| Ollama Cloud | `OLLAMA_CLOUD_API_KEY` |
| Cloudflare | `CLOUDFLARE_ACCOUNT_ID`, `CLOUDFLARE_API_TOKEN` |
| Helicone | `HELICONE_API_KEY` |
| SAP AI Core | `AICORE_SERVICE_KEY` |
| STACKIT | `STACKIT_TOKEN` |
| OVHcloud | `OVHCLOUD_API_KEY` |
| Scaleway | `SCALEWAY_API_KEY` |

---

## Model Naming Convention

| Provider | Model Prefix | Example |
|----------|--------------|---------|
| OpenAI | `openai/` | `openai/gpt-4o` |
| Anthropic | `anthropic/` | `anthropic/claude-3-5-sonnet` |
| Groq | `groq/` | `groq/llama-3.3-70b` |
| DeepSeek | `deepseek/` | `deepseek/deepseek-chat` |
| OpenRouter | `openrouter/` | `openrouter/anthropic/claude-3.5-sonnet` |
| Azure | `azure/` | `azure/gpt-4o` |
| Vertex | `vertex/` | `vertex/gemini-2.5-pro` |

---

## Progress Tracker

### Core Providers (DONE)
- [x] OpenAI
- [x] Anthropic
- [x] Minimax
- [x] Ollama
- [x] LiteLLM

### OpenAI-Compatible Providers (DONE)
- [x] Groq
- [x] Perplexity
- [x] Mistral
- [x] Cohere
- [x] Together AI
- [x] DeepInfra
- [x] Cerebras
- [x] xAI
- [x] Alibaba
- [x] Hugging Face
- [x] DeepSeek
- [x] Fireworks AI
- [x] Moonshot AI
- [x] Nebius
- [x] OpenRouter

### Cloud/Enterprise (DONE)
- [x] Google Gemini
- [x] Azure OpenAI
- [x] Google Vertex AI
- [x] AWS Bedrock

### DevOps/Git Integration (DONE)
- [x] GitLab Duo
- [x] GitHub Copilot
- [x] Vercel AI

### Regional/Specialty (DONE)
- [x] Venice AI
- [x] Z.AI
- [x] ZenMux
- [x] Baseten
- [x] Cortecs
- [x] Firmware AI
- [x] IO.NET
- [x] 302.AI
- [x] NVIDIA
- [x] Ollama Cloud

### Gateway/Proxy (DONE)
- [x] Cloudflare AI Gateway
- [x] Cloudflare Workers AI
- [x] Helicone

### Local Models (DONE)
- [x] llama.cpp
- [x] LM Studio
- [x] Atomic Chat

### Enterprise Cloud (DONE)
- [x] SAP AI Core
- [x] STACKIT
- [x] OVHcloud AI
- [x] Scaleway
- [x] Azure Cognitive Services

### models.dev Integration (DONE)
- [x] Registry generator (`cmd/generate/main.go`)
- [x] Generated registry (`internal/provider/registry.json`)
- [x] 116 providers from models.dev

## Optimization Notes

### Known Duplication
25+ OpenAI-compatible providers have identical 90-line implementations differing only in:
- Provider struct name
- Provider name string
- Base URL

**Opportunity**: Create a single `OpenAICompatibleProvider` that takes name, baseURL, apiKey as parameters.

### Not Done (Acceptable)
- models.dev runtime fetching (build-time generation only)
- Full model metadata in ProviderModels (only basic info)

---

## Author

Mark LaPointe <mark@cloudbsd.org>
