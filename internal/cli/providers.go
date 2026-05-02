package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/freecode/freecode/internal/auth"
	"github.com/freecode/freecode/internal/provider"
	"github.com/spf13/cobra"
)

var providersLoginProvider string
var providersLoginMethod string

var providersCmd = &cobra.Command{
	Use:   "providers",
	Short: "Manage AI providers and credentials",
	Aliases: []string{"auth"},
}

var providersListCmd = &cobra.Command{
	Use:   "list",
	Short: "List configured providers and credentials",
	Aliases: []string{"ls"},
	RunE:  runProvidersList,
}

var providersLoginCmd = &cobra.Command{
	Use:   "login [provider]",
	Short: "Login to a provider with API key or OAuth",
	RunE:  runProvidersLogin,
}

var providersLogoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Remove credentials for a provider",
	RunE:  runProvidersLogout,
}

func init() {
	providersCmd.AddCommand(providersListCmd)
	providersCmd.AddCommand(providersLoginCmd)
	providersCmd.AddCommand(providersLogoutCmd)

	providersLoginCmd.Flags().StringVarP(&providersLoginProvider, "provider", "p", "", "Provider ID to login to")
	providersLoginCmd.Flags().StringVarP(&providersLoginMethod, "method", "m", "", "Login method (api, oauth)")
}

func runProvidersList(cmd *cobra.Command, args []string) error {
	store := auth.Default()
	creds := store.All()

	fmt.Println()
	if len(creds) == 0 {
		fmt.Println("No credentials configured.")
		fmt.Println("\nTo add credentials:")
		fmt.Println("  freecode providers login <provider>")
		fmt.Println("\nSupported providers:")
		for _, p := range supportedProviders() {
			fmt.Printf("  - %s\n", p)
		}
		return nil
	}

	displayPath := store.Path()
	fmt.Printf("Credentials: %s\n\n", displayPath)

	for providerID, info := range creds {
		providerName := getProviderDisplayName(providerID)
		fmt.Printf("%s %s\n", providerName, info.Type)
	}

	fmt.Printf("\n%d credential(s)\n", len(creds))

	showEnvProviders()

	return nil
}

func showEnvProviders() {
	envProviders := []struct {
		name    string
		envVars []string
	}{
		{"OpenAI", []string{"OPENAI_API_KEY"}},
		{"Anthropic", []string{"ANTHROPIC_API_KEY"}},
		{"MiniMax", []string{"MINIMAX_API_KEY"}},
		{"LiteLLM", []string{"LITELLM_BASE_URL", "LITELLM_API_KEY"}},
		{"Ollama", []string{"OLLAMA_BASE_URL"}},
	}

	var active []string
	for _, p := range envProviders {
		for _, envVar := range p.envVars {
			if os.Getenv(envVar) != "" {
				active = append(active, fmt.Sprintf("%s (%s)", p.name, envVar))
				break
			}
		}
	}

	if len(active) > 0 {
		fmt.Println("\nEnvironment:")
		for _, name := range active {
			fmt.Printf("  %s\n", name)
		}
		fmt.Printf("\n%d environment variable(s)\n", len(active))
	}
}

func runProvidersLogin(cmd *cobra.Command, args []string) error {
	store := auth.Default()

	providerID := providersLoginProvider
	if providerID == "" && len(args) > 0 {
		providerID = args[0]
	}

	if providerID == "" {
		providerID = selectProviderInteractive()
		if providerID == "" {
			return nil
		}
	}

	providerID = strings.ToLower(providerID)

	if !isValidProvider(providerID) {
		fmt.Fprintf(os.Stderr, "Unknown provider: %s\n", providerID)
		fmt.Fprintf(os.Stderr, "Supported providers: %s\n", strings.Join(supportedProviders(), ", "))
		return fmt.Errorf("unknown provider")
	}

	if providersLoginMethod == "" || providersLoginMethod == "api" {
		return loginWithAPIKey(store, providerID)
	}

	return fmt.Errorf("unsupported login method: %s", providersLoginMethod)
}

func loginWithAPIKey(store *auth.Store, providerID string) error {
	fmt.Printf("\nEnter API key for %s: ", getProviderDisplayName(providerID))

	reader := bufio.NewReader(os.Stdin)
	apiKey, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read API key: %w", err)
	}

	apiKey = strings.TrimSpace(apiKey)
	if apiKey == "" {
		return fmt.Errorf("API key cannot be empty")
	}

	info := auth.Info{
		Type: auth.AuthTypeAPI,
		Key:  apiKey,
	}

	if err := store.Set(providerID, info); err != nil {
		return fmt.Errorf("failed to save credentials: %w", err)
	}

	fmt.Printf("\nSuccessfully logged in to %s\n", getProviderDisplayName(providerID))
	return nil
}

func runProvidersLogout(cmd *cobra.Command, args []string) error {
	store := auth.Default()
	creds := store.All()

	if len(creds) == 0 {
		fmt.Println("No credentials to remove.")
		return nil
	}

	providerID := ""
	if len(args) > 0 {
		providerID = args[0]
	} else {
		providerID = selectProviderToRemove(creds)
		if providerID == "" {
			return nil
		}
	}

	providerID = strings.ToLower(providerID)

	info, ok := creds[providerID]
	if !ok {
		fmt.Fprintf(os.Stderr, "No credentials found for: %s\n", providerID)
		return fmt.Errorf("provider not found")
	}

	if err := store.Remove(providerID); err != nil {
		return fmt.Errorf("failed to remove credentials: %w", err)
	}

	fmt.Printf("Successfully logged out from %s (%s)\n", getProviderDisplayName(providerID), info.Type)
	return nil
}

func selectProviderInteractive() string {
	providers := supportedProviders()

	fmt.Println("\nSelect a provider:")
	for i, p := range providers {
		fmt.Printf("  %d. %s\n", i+1, getProviderDisplayName(p))
	}
	fmt.Println("  0. Cancel")

	fmt.Print("\nEnter number: ")

	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil {
		return ""
	}

	line = strings.TrimSpace(line)
	var idx int
	if _, err := fmt.Sscanf(line, "%d", &idx); err != nil {
		return ""
	}

	if idx < 0 || idx > len(providers) {
		return ""
	}

	if idx == 0 {
		return ""
	}

	return providers[idx-1]
}

func selectProviderToRemove(creds map[string]auth.Info) string {
	var providers []string
	for id := range creds {
		providers = append(providers, id)
	}

	if len(providers) == 0 {
		return ""
	}

	fmt.Println("\nSelect a provider to remove:")
	for i, p := range providers {
		fmt.Printf("  %d. %s\n", i+1, getProviderDisplayName(p))
	}
	fmt.Println("  0. Cancel")

	fmt.Print("\nEnter number: ")

	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil {
		return ""
	}

	line = strings.TrimSpace(line)
	var idx int
	if _, err := fmt.Sscanf(line, "%d", &idx); err != nil {
		return ""
	}

	if idx < 0 || idx > len(providers) {
		return ""
	}

	if idx == 0 {
		return ""
	}

	return providers[idx-1]
}

func supportedProviders() []string {
	return []string{
		"openai",
		"anthropic",
		"google",
		"azure",
		"bedrock",
		"vertex",
		"openrouter",
		"ollama",
		"litellm",
		"minimax",
		"groq",
		"deepseek",
		"mistral",
		"cohere",
		"ai21",
		"opencode",
		"github-copilot",
		"vercel",
		"cloudflare",
		"cloudflare-ai-gateway",
		"openai-compatible",
		"anthropic-compatible",
		"azure-openai",
		"aws-bedrock",
		"aws-sagemaker",
		"together",
		"perplexity",
		"replicate",
		"anyscale",
		"mistral-api",
		"nvidia",
		"fireworks",
		"voyage",
		"voyageai",
		"codestral",
		"cerebras",
		"hyperbolic",
		"genhuan",
		"zhipu",
		"baichuan",
		"tencent",
		"tencent-hunyuan",
		"baidu",
		" volcengine",
		"moonshot",
		"qwen",
		"yi",
		"stepfun",
		"ali-bailian",
		"01-ai",
		"meta",
		"meta-llama",
		"x-ai",
		"grok",
		"localai",
		"定西",
		"inference",
	}
}

func isValidProvider(id string) bool {
	providers := supportedProviders()
	for _, p := range providers {
		if p == id {
			return true
		}
	}
	return false
}

func getProviderDisplayName(id string) string {
	names := map[string]string{
		"openai":    "OpenAI",
		"anthropic": "Anthropic",
		"google":    "Google",
		"azure":     "Azure",
		"bedrock":   "Amazon Bedrock",
		"vertex":    "Google Vertex AI",
		"openrouter": "OpenRouter",
		"ollama":    "Ollama",
		"litellm":   "LiteLLM",
		"minimax":   "MiniMax",
		"groq":      "Groq",
		"deepseek":  "DeepSeek",
		"mistral":   "Mistral",
		"cohere":    "Cohere",
		"ai21":      "AI21",
		"opencode":  "OpenCode",
	}

	if name, ok := names[id]; ok {
		return name
	}
	return id
}

func getProviderModel(id string) string {
	defaultModels := map[string]string{
		"openai":    "gpt-4o",
		"anthropic": "claude-sonnet-4-5",
		"google":    "gemini-2.0-flash",
		"azure":     "gpt-4o",
		"bedrock":   "anthropic.claude-3-5-sonnet",
		"vertex":    "gemini-2.0-flash",
		"openrouter": "openai/gpt-4o",
		"ollama":    "llama3.2",
		"litellm":   "gpt-4o",
		"minimax":   "MiniMax-Text-01",
		"groq":      "llama-3.3-70b",
		"deepseek":  "deepseek-chat",
		"mistral":   "mistral-large-latest",
		"cohere":    "command-r-plus",
		"ai21":      "jamba-1-5-large",
		"opencode":  "opencode",
	}

	if model, ok := defaultModels[id]; ok {
		return model
	}
	return ""
}

type providerConfig struct {
	envVar        string
	defaultBaseURL string
	defaultModel  string
}

var providerConfigs = map[string]providerConfig{
	"openai": {
		envVar:        "OPENAI_API_KEY",
		defaultBaseURL: "https://api.openai.com/v1",
		defaultModel:  "gpt-4o",
	},
	"anthropic": {
		envVar:        "ANTHROPIC_API_KEY",
		defaultBaseURL: "https://api.anthropic.com",
		defaultModel:  "claude-sonnet-4-5",
	},
	"google": {
		envVar:        "GOOGLE_API_KEY",
		defaultBaseURL: "https://generativelanguage.googleapis.com/v1",
		defaultModel:  "gemini-2.0-flash",
	},
	"ollama": {
		envVar:        "OLLAMA_BASE_URL",
		defaultBaseURL: "http://localhost:11434",
		defaultModel:  "llama3.2",
	},
	"litellm": {
		envVar:        "LITELLM_BASE_URL",
		defaultBaseURL: "http://localhost:4000",
		defaultModel:  "gpt-4o",
	},
}

func (p *providerConfig) getAPIKey() string {
	return os.Getenv(p.envVar)
}

func getConfiguredProviders() []*providerInfo {
	var result []*providerInfo

	for id, cfg := range providerConfigs {
		if apiKey := cfg.getAPIKey(); apiKey != "" {
			result = append(result, &providerInfo{
				id:       id,
				apiKey:   apiKey,
				baseURL:  os.Getenv(cfg.envVar + "_BASE_URL"),
				model:    getProviderModel(id),
			})
		}
	}

	return result
}

type providerInfo struct {
	id      string
	apiKey  string
	baseURL string
	model   string
}

func (p *providerInfo) createProvider() provider.Provider {
	switch p.id {
	case "openai":
		return provider.NewOpenAIProvider(p.apiKey)
	case "anthropic":
		return provider.NewAnthropicProvider(p.apiKey)
	case "ollama":
		return provider.NewOllamaProvider(p.baseURL, p.apiKey)
	case "litellm":
		return provider.NewLiteLLMProvider(p.baseURL, p.apiKey)
	default:
		return nil
	}
}