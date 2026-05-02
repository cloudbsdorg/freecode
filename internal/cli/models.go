package cli

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/freecode/freecode/internal/config"
	"github.com/freecode/freecode/internal/provider"
	"github.com/spf13/cobra"
)

var (
	modelProvider string
	modelRefresh  bool
	modelList     bool
)

var modelsCmd = &cobra.Command{
	Use:   "models",
	Short: "List and manage available models",
	Long: `Discover and display available models from configured providers.

Examples:
  freecode models              # List all available models
  freecode models --provider openai  # List only OpenAI models
  freecode models --refresh   # Force refresh model cache`,
	RunE: runModels,
}

func init() {
	modelsCmd.Flags().StringVar(&modelProvider, "provider", "", "Filter by provider (openai, anthropic, minimax)")
	modelsCmd.Flags().BoolVar(&modelRefresh, "refresh", false, "Force refresh model cache")
	modelsCmd.Flags().BoolVar(&modelList, "list", false, "List models in simple format")
}

func runModels(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load("")
	if err != nil {
		cfg = config.DefaultConfig()
	}

	svc := provider.NewCatalogService()

	if err := svc.Load(); err != nil {
		fmt.Fprintf(os.Stderr, "Note: Could not load cached models: %v\n", err)
	}

	if modelRefresh {
		fmt.Println("Refreshing model catalog...")
		if err := refreshProviders(cfg, svc); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Some providers failed to refresh: %v\n", err)
		}
		if err := svc.Save(); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Could not save cache: %v\n", err)
		}
	}

	connected := discoverConnectedProviders(cfg)
	if len(connected) > 0 {
		fmt.Fprintf(os.Stderr, "Discovering models from %d connected providers...\n", len(connected))
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		for _, p := range connected {
			if err := svc.DiscoverFromProvider(ctx, p); err != nil {
				fmt.Fprintf(os.Stderr, "Warning: Failed to discover %s models: %v\n", p.Name(), err)
			}
		}
		_ = svc.Save()
	}

	catalog := svc.GetCatalog()
	providers := catalog.AllProviders()

	if modelProvider != "" {
		filtered := make(map[string]*provider.ProviderModels)
		for name, pm := range providers {
			if strings.EqualFold(name, modelProvider) || strings.EqualFold(name, modelProvider) {
				filtered[name] = pm
			}
		}
		providers = filtered
	}

	if len(providers) == 0 {
		fmt.Println("No providers configured. Set API keys in environment or config.")
		fmt.Println("\nSupported providers:")
		fmt.Println("  - ANTHROPIC_API_KEY  for Claude models")
		fmt.Println("  - OPENAI_API_KEY     for GPT models")
		fmt.Println("  - MINIMAX_API_KEY    for MiniMax models")
		fmt.Println("  - LITELLM_BASE_URL   for LiteLLM proxies")
		return nil
	}

	displayModels(providers, modelList)
	return nil
}

func discoverConnectedProviders(cfg *config.Config) []provider.Provider {
	var connected []provider.Provider

	if apiKey := os.Getenv("ANTHROPIC_API_KEY"); apiKey != "" {
		connected = append(connected, provider.NewAnthropicProvider(apiKey))
	}
	if apiKey := os.Getenv("OPENAI_API_KEY"); apiKey != "" {
		connected = append(connected, provider.NewOpenAIProvider(apiKey))
	}
	if apiKey := os.Getenv("MINIMAX_API_KEY"); apiKey != "" {
		connected = append(connected, provider.NewMinimaxProvider(apiKey))
	}

	// Ollama - discover if running locally (no API key required)
	if baseURL := os.Getenv("OLLAMA_BASE_URL"); baseURL != "" {
		connected = append(connected, provider.NewOllamaProvider(baseURL, ""))
	} else {
		// Try default Ollama URL
		connected = append(connected, provider.NewOllamaProvider("http://localhost:11434", ""))
	}

	if baseURL := os.Getenv("LITELLM_BASE_URL"); baseURL != "" {
		apiKey := os.Getenv("LITELLM_API_KEY")
		connected = append(connected, provider.NewLiteLLMProvider(baseURL, apiKey))
	}

	if cfg.LiteLLM.BaseURL != "" && cfg.LiteLLM.BaseURL != "http://localhost:4000" {
		connected = append(connected, provider.NewLiteLLMProvider(cfg.LiteLLM.BaseURL, cfg.LiteLLM.APIKey))
	}
	if cfg.Minimax.APIKey != "" {
		baseURL := cfg.Minimax.BaseURL
		if baseURL == "" {
			baseURL = "https://api.minimax.chat/v1"
		}
		p := provider.NewMinimaxProvider(cfg.Minimax.APIKey)
		p.BaseURL = baseURL
		connected = append(connected, p)
	}

	return connected
}

func refreshProviders(cfg *config.Config, svc *provider.CatalogService) error {
	connected := discoverConnectedProviders(cfg)
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	var lastErr error
	for _, p := range connected {
		if err := svc.DiscoverFromProvider(ctx, p); err != nil {
			lastErr = err
			fmt.Fprintf(os.Stderr, "Warning: %s refresh failed: %v\n", p.Name(), err)
		}
	}
	return lastErr
}

func displayModels(providers map[string]*provider.ProviderModels, simple bool) {
	type modelEntry struct {
		provider string
		model   provider.Model
	}

	var all []modelEntry
	for providerName, pm := range providers {
		for _, m := range pm.Models {
			all = append(all, modelEntry{provider: providerName, model: m})
		}
	}

	sort.Slice(all, func(i, j int) bool {
		if all[i].provider != all[j].provider {
			return all[i].provider < all[j].provider
		}
		return all[i].model.ID < all[j].model.ID
	})

	if simple {
		for _, entry := range all {
			fmt.Printf("%s/%s\n", entry.provider, entry.model.ID)
		}
		return
	}

	fmt.Printf("\n%-12s %-40s %-12s %s\n", "PROVIDER", "MODEL ID", "CONTEXT", "CAPABILITIES")
	fmt.Println(strings.Repeat("-", 80))

	currentProvider := ""
	for _, entry := range all {
		if entry.provider != currentProvider {
			currentProvider = entry.provider
			fmt.Println()
		}

		capabilities := []string{}
		if entry.model.Capabilities.Reasoning {
			capabilities = append(capabilities, "reasoning")
		}
		if entry.model.Capabilities.ToolCall {
			capabilities = append(capabilities, "tools")
		}
		if entry.model.Capabilities.Vision {
			capabilities = append(capabilities, "vision")
		}
		if entry.model.Capabilities.Audio {
			capabilities = append(capabilities, "audio")
		}
		if len(capabilities) == 0 {
			capabilities = append(capabilities, "text")
		}

		limit := entry.model.Limit.Context
		if limit == 0 {
			limit = 128000
		}

		fmt.Printf("%-12s %-40s %-12d %s\n",
			entry.provider,
			entry.model.ID,
			limit/1000,
			strings.Join(capabilities, ","),
		)
	}
	fmt.Println()
}