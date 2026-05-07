package cli

import (
	"fmt"
	"os"

	"github.com/freecode/freecode/internal/config"
	"github.com/spf13/cobra"
)

var (
	accountListCmd   bool
	accountActiveCmd bool
)

var accountCmd = &cobra.Command{
	Use:   "account",
	Short: "Manage account settings and provider credentials",
	Long: `Manage your freecode account settings and configure provider credentials.

Examples:
  freecode account              # Show current account status
  freecode account list         # List all configured accounts/credentials
  freecode account active       # Show active account settings`,
	RunE: runAccount,
}

func init() {
	accountCmd.Flags().BoolVar(&accountListCmd, "list", false, "List all configured credentials")
	accountCmd.Flags().BoolVar(&accountActiveCmd, "active", false, "Show active account settings")
}

func runAccount(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load("")
	if err != nil {
		cfg = config.DefaultConfig()
	}

	if accountListCmd {
		return listAccounts(cfg)
	}

	if accountActiveCmd {
		return showActive(cfg)
	}

	return showAccountStatus(cfg)
}

func listAccounts(cfg *config.Config) error {
	fmt.Println("Configured Providers:")
	fmt.Println("")

	providers := []struct {
		name   string
		keyEnv string
		keySet bool
	}{
		{"Anthropic", "ANTHROPIC_API_KEY", os.Getenv("ANTHROPIC_API_KEY") != ""},
		{"OpenAI", "OPENAI_API_KEY", os.Getenv("OPENAI_API_KEY") != ""},
		{"Google", "GOOGLE_API_KEY", os.Getenv("GOOGLE_API_KEY") != ""},
		{"Mistral", "MISTRAL_API_KEY", os.Getenv("MISTRAL_API_KEY") != ""},
		{"Groq", "GROQ_API_KEY", os.Getenv("GROQ_API_KEY") != ""},
		{"OpenRouter", "OPENROUTER_API_KEY", os.Getenv("OPENROUTER_API_KEY") != ""},
		{"Ollama", "OLLAMA_BASE_URL", os.Getenv("OLLAMA_BASE_URL") != ""},
		{"LiteLLM", "LITELLM_BASE_URL", os.Getenv("LITELLM_BASE_URL") != ""},
	}

	for _, p := range providers {
		status := "not configured"
		if p.keySet {
			status = "configured"
		}
		fmt.Printf("  %-15s %s (%s)\n", p.name, status, p.keyEnv)
	}

	fmt.Println("")
	fmt.Println("Config locations to check:")
	fmt.Println("  ~/.config/freecode/config.yaml")
	fmt.Println("  ~/.config/freecode/config.json")

	return nil
}

func showActive(cfg *config.Config) error {
	fmt.Println("Active Configuration:")
	fmt.Println("")
	fmt.Println("  Config file: ~/.config/freecode/config.yaml")
	fmt.Println("")
	fmt.Println("  Provider configurations:")
	for name := range cfg.Providers {
		fmt.Printf("    %s: configured\n", name)
	}

	if cfg.LiteLLM.BaseURL != "" {
		fmt.Println("  LiteLLM:")
		fmt.Printf("    Base URL: %s\n", cfg.LiteLLM.BaseURL)
		if cfg.LiteLLM.APIKey != "" {
			fmt.Println("    API Key: configured")
		} else {
			fmt.Println("    API Key: not configured")
		}
	}

	if cfg.Minimax.APIKey != "" {
		fmt.Println("  MiniMax:")
		fmt.Printf("    Base URL: %s\n", cfg.Minimax.BaseURL)
		fmt.Println("    API Key: configured")
	}

	return nil
}

func showAccountStatus(cfg *config.Config) error {
	fmt.Println("Freecode Account Status")
	fmt.Println("======================")
	fmt.Println("")

	hasAnyKey := false
	providers := []string{
		"ANTHROPIC_API_KEY",
		"OPENAI_API_KEY",
		"GOOGLE_API_KEY",
		"MISTRAL_API_KEY",
		"GROQ_API_KEY",
		"OPENROUTER_API_KEY",
	}

	fmt.Println("Provider API Keys:")
	for _, env := range providers {
		if os.Getenv(env) != "" {
			hasAnyKey = true
			fmt.Printf("  ✓ %s is set\n", env)
		}
	}

	if !hasAnyKey && cfg.LiteLLM.APIKey == "" && cfg.Minimax.APIKey == "" {
		fmt.Println("  No provider API keys configured.")
		fmt.Println("")
		fmt.Println("  To use freecode, set at least one API key:")
		fmt.Println("    export ANTHROPIC_API_KEY=sk-...")
		fmt.Println("    export OPENAI_API_KEY=sk-...")
	}

	fmt.Println("")
	fmt.Println("Quick Commands:")
	fmt.Println("  freecode account --list    # See all configured providers")
	fmt.Println("  freecode account --active  # See active configuration")
	fmt.Println("  freecode models           # List available models")

	return nil
}
