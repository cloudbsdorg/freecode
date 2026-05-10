package config

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

const modelsDevURL = "https://models.dev/api.json"

type BootstrapOptions struct {
	Force       bool
	Interactive bool
	Provider    string
	APIKey      string
	Model       string
}

type BootstrapResult struct {
	ConfigPath string
	Provider   string
	Model      string
	SessionDir string
	SkillsDir  string
	WizardRan  bool
}

type modelsDevProvider struct {
	ID     string            `json:"id"`
	Name   string            `json:"name"`
	API    string            `json:"api"`
	Env    []string          `json:"env"`
	Models map[string]struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"models"`
}

type modelsDevRegistry map[string]modelsDevProvider

func Bootstrap(opts BootstrapOptions) (*BootstrapResult, error) {
	paths := PathsGet()

	if err := paths.Ensure(); err != nil {
		return nil, fmt.Errorf("failed to create directories: %w", err)
	}

	result := &BootstrapResult{
		ConfigPath: paths.ConfigFile("config.yaml"),
		SessionDir: paths.SessionDir(),
		SkillsDir:  paths.SkillsDir(),
	}

	configExists := fileExists(result.ConfigPath)

	if opts.Force || !configExists {
		wizardResult, err := runWizard(paths, opts)
		if err != nil {
			return nil, fmt.Errorf("setup failed: %w", err)
		}
		result.Provider = wizardResult.Provider
		result.Model = wizardResult.Model
		result.WizardRan = true
	} else {
		cfg, err := Load(result.ConfigPath)
		if err == nil && cfg != nil {
			result.Provider = getConfiguredProvider(cfg)
			result.Model = getConfiguredModel(cfg)
		}
	}

	return result, nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

type wizardResult struct {
	Provider string
	Model   string
}

func FetchProvidersFromFeed() (modelsDevRegistry, error) {
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Get(modelsDevURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch provider list: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("received status %d from %s", resp.StatusCode, modelsDevURL)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var registry modelsDevRegistry
	if err := json.Unmarshal(body, &registry); err != nil {
		return nil, fmt.Errorf("failed to parse provider list: %w", err)
	}

	return registry, nil
}

func runWizard(paths *Paths, opts BootstrapOptions) (*wizardResult, error) {
	fmt.Println("=== Freecode Initial Setup ===")
	fmt.Println()

	providerID := opts.Provider
	model := opts.Model
	apiKey := opts.APIKey

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Fetching provider list from models.dev...")

	registry, err := FetchProvidersFromFeed()
	if err != nil {
		fmt.Printf("Warning: failed to fetch providers: %v\n", err)
		fmt.Println("Cannot continue without provider list.")
		return nil, err
	}
	fmt.Printf("Found %d providers\n\n", len(registry))

	if providerID == "" {
		providerID = promptProviderSelection(reader, registry)
	}

	prov := registry[providerID]

	if model == "" {
		model = promptModelSelection(reader, providerID, prov)
	}

	if apiKey == "" && !isLocalProvider(prov) {
		apiKey = promptAPIKey(reader, providerID, prov)
	}

	if err := createConfigForProvider(paths, providerID, model, apiKey, prov); err != nil {
		return nil, fmt.Errorf("failed to create config: %w", err)
	}

	fmt.Println()
	fmt.Printf("Configuration saved to: %s\n", paths.ConfigFile("config.yaml"))
	fmt.Printf("Sessions stored in: %s\n", paths.SessionDir())
	fmt.Println()

	return &wizardResult{Provider: providerID, Model: model}, nil
}

func promptProviderSelection(reader *bufio.Reader, registry modelsDevRegistry) string {
	providerIDs := sortedProviderIDs(registry)

	fmt.Println("Available providers:")
	for i, id := range providerIDs {
		prov := registry[id]
		modelCount := len(prov.Models)
		if modelCount > 0 {
			fmt.Printf("  %2d. %s (%s - %d models)\n", i+1, prov.Name, id, modelCount)
		} else {
			fmt.Printf("  %2d. %s (%s)\n", i+1, prov.Name, id)
		}
	}
	fmt.Println()

	fmt.Print("Select provider number: ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	idx := 0
	fmt.Sscanf(input, "%d", &idx)
	if idx < 1 || idx > len(providerIDs) {
		idx = 1
	}
	return providerIDs[idx-1]
}

func promptModelSelection(reader *bufio.Reader, providerID string, prov modelsDevProvider) string {
	if len(prov.Models) == 0 {
		fmt.Print("Enter model ID: ")
		input, _ := reader.ReadString('\n')
		return strings.TrimSpace(input)
	}

	modelIDs := sortedModelIDs(prov.Models)

	fmt.Println("\nAvailable models:")
	for i, id := range modelIDs {
		if i >= 20 {
			fmt.Printf("  ... and %d more\n", len(modelIDs)-20)
			break
		}
		m := prov.Models[id]
		if m.Name != "" && m.Name != id {
			fmt.Printf("  %2d. %s (%s)\n", i+1, m.Name, id)
		} else {
			fmt.Printf("  %2d. %s\n", i+1, id)
		}
	}
	fmt.Println()

	fmt.Print("Select model number [1]: ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if input == "" {
		input = "1"
	}

	idx := 0
	fmt.Sscanf(input, "%d", &idx)
	if idx < 1 || idx > len(modelIDs) {
		idx = 1
	}
	return modelIDs[idx-1]
}

func promptAPIKey(reader *bufio.Reader, providerID string, prov modelsDevProvider) string {
	envVar := ""
	if len(prov.Env) > 0 {
		envVar = prov.Env[0]
	}

	if envVar != "" {
		fmt.Printf("Enter API key for %s (env: %s): ", providerID, envVar)
	} else {
		fmt.Printf("Enter API key for %s: ", providerID)
	}

	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func isLocalProvider(prov modelsDevProvider) bool {
	if prov.API == "" {
		return false
	}
	return strings.Contains(prov.API, "localhost") ||
		strings.Contains(prov.API, "127.0.0.1") ||
		strings.Contains(prov.API, "ollama") ||
		strings.Contains(prov.API, "lmstudio")
}

func sortedProviderIDs(registry modelsDevRegistry) []string {
	ids := make([]string, 0, len(registry))
	for id := range registry {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	return ids
}

func sortedModelIDs(models map[string]struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}) []string {
	ids := make([]string, 0, len(models))
	for id := range models {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	return ids
}

func createConfigForProvider(paths *Paths, providerID, model, apiKey string, prov modelsDevProvider) error {
	cfg := DefaultConfig()
	cfg.Session.Dir = paths.SessionDir()

	cfg.Models = map[string]ModelConfig{
		"default": {
			Provider: providerID,
			Name:     model,
		},
	}

	baseURL := prov.API
	if baseURL == "" {
		baseURL = inferBaseURL(providerID)
	}

	cfg.Providers = map[string]ProviderConfig{
		providerID: {
			APIKey:  apiKey,
			BaseURL: baseURL,
		},
	}

	return SaveConfig(paths.ConfigFile("config.yaml"), cfg)
}

func inferBaseURL(providerID string) string {
	guesses := map[string]string{
		"ollama":       "http://localhost:11434",
		"lmstudio":     "http://localhost:1234",
		"ollama-cloud": "https://ollama.cloud",
	}
	if url, ok := guesses[providerID]; ok {
		return url
	}
	return ""
}

func SaveConfig(path string, cfg *Config) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	if err := cfg.SaveYAML(path); err != nil {
		return fmt.Errorf("failed to save YAML: %w", err)
	}

	return nil
}

func getConfiguredProvider(cfg *Config) string {
	if cfg == nil {
		return ""
	}
	if len(cfg.Models) > 0 {
		for _, m := range cfg.Models {
			if m.Provider != "" {
				return m.Provider
			}
		}
	}
	if cfg.OpenAI.APIKey != "" {
		return "openai"
	}
	if cfg.Anthropic.APIKey != "" {
		return "anthropic"
	}
	return ""
}

func getConfiguredModel(cfg *Config) string {
	if cfg == nil {
		return ""
	}
	if m, ok := cfg.Models["default"]; ok && m.Name != "" {
		return m.Name
	}
	return ""
}
