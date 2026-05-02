package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type Model struct {
	ID            string
	Name          string
	Provider      string
	Capabilities  ModelCapabilities
	Cost          ModelCost
	Limit        ModelLimit
	Created      int64
	OwnedBy      string
}

type ModelCapabilities struct {
	Temperature bool
	Reasoning  bool
	ToolCall   bool
	Vision     bool
	Audio      bool
}

type ModelCost struct {
	Input  float64
	Output float64
}

type ModelLimit struct {
	Context int
	Input   int
	Output  int
}

type ProviderModels struct {
	Provider   string
	Name       string
	APIKey     string
	BaseURL    string
	APIType    string
	Models     []Model
	LastUpdate time.Time
}

var (
	modelCache     = make(map[string]*ProviderModels)
	modelCacheMu   sync.RWMutex
	cacheTTL       = 5 * time.Minute
	cacheDir       string
)

func init() {
	cacheDir = filepath.Join(os.TempDir(), "freecode-models")
	_ = os.MkdirAll(cacheDir, 0755)
}

type ModelCatalog struct {
	providers map[string]*ProviderModels
	mu        sync.RWMutex
}

func NewModelCatalog() *ModelCatalog {
	return &ModelCatalog{
		providers: make(map[string]*ProviderModels),
	}
}

func (c *ModelCatalog) GetProvider(providerName string) *ProviderModels {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.providers[providerName]
}

func (c *ModelCatalog) SetProvider(providerName string, pm *ProviderModels) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.providers[providerName] = pm
}

func (c *ModelCatalog) AllProviders() map[string]*ProviderModels {
	c.mu.RLock()
	defer c.mu.RUnlock()
	result := make(map[string]*ProviderModels, len(c.providers))
	for k, v := range c.providers {
		result[k] = v
	}
	return result
}

func (c *ModelCatalog) GetModel(providerName, modelID string) *Model {
	pm := c.GetProvider(providerName)
	if pm == nil {
		return nil
	}
	for i := range pm.Models {
		if pm.Models[i].ID == modelID {
			return &pm.Models[i]
		}
	}
	return nil
}

func (pm *ProviderModels) IsStale() bool {
	return time.Since(pm.LastUpdate) > cacheTTL
}

type modelListResponse struct {
	Object string `json:"object"`
	Data   []struct {
		ID       string `json:"id"`
		Object   string `json:"object"`
		Created  int64  `json:"created"`
		OwnedBy  string `json:"owned_by"`
	} `json:"data"`
}

func fetchModelsFromURL(ctx context.Context, url, apiKey string) ([]Model, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	var result modelListResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	models := make([]Model, 0, len(result.Data))
	for _, m := range result.Data {
		models = append(models, Model{
			ID:       m.ID,
			Name:     m.ID,
			Provider: "",
			Created:  m.Created,
			OwnedBy:  m.OwnedBy,
			Capabilities: ModelCapabilities{
				Temperature: true,
				Reasoning:  false,
				ToolCall:   true,
				Vision:     false,
				Audio:     false,
			},
			Cost: ModelCost{Input: 0, Output: 0},
			Limit: ModelLimit{
				Context: 128000,
				Input:   0,
				Output:  0,
			},
		})
	}
	return models, nil
}

func (p *OpenAIProvider) ListModels(ctx context.Context) ([]Model, error) {
	url := p.BaseURL + "/models"
	return fetchModelsFromURL(ctx, url, p.APIKey)
}

func (p *AnthropicProvider) ListModels(ctx context.Context) ([]Model, error) {
	return getAnthropicModels(), nil
}

func (p *MinimaxProvider) ListModels(ctx context.Context) ([]Model, error) {
	url := p.BaseURL + "/models"
	return fetchModelsFromURL(ctx, url, p.APIKey)
}

func (p *LiteLLMProvider) ListModels(ctx context.Context) ([]Model, error) {
	url := p.BaseURL + "/models"
	return fetchModelsFromURL(ctx, url, p.APIKey)
}

func getAnthropicModels() []Model {
	return []Model{
		{ID: "claude-opus-4-5", Name: "Claude Opus 4", Provider: "anthropic", OwnedBy: "anthropic", Created: 1739836800, Capabilities: ModelCapabilities{Temperature: true, Reasoning: true, ToolCall: true, Vision: true}, Limit: ModelLimit{Context: 200000, Output: 8192}},
		{ID: "claude-sonnet-4-5", Name: "Claude Sonnet 4", Provider: "anthropic", OwnedBy: "anthropic", Created: 1739836800, Capabilities: ModelCapabilities{Temperature: true, Reasoning: true, ToolCall: true, Vision: true}, Limit: ModelLimit{Context: 200000, Output: 8192}},
		{ID: "claude-haiku-4-5", Name: "Claude Haiku 4", Provider: "anthropic", OwnedBy: "anthropic", Created: 1739836800, Capabilities: ModelCapabilities{Temperature: true, Reasoning: false, ToolCall: true, Vision: true}, Limit: ModelLimit{Context: 200000, Output: 8192}},
		{ID: "claude-3-5-opus", Name: "Claude 3.5 Opus", Provider: "anthropic", OwnedBy: "anthropic", Created: 1714080000, Capabilities: ModelCapabilities{Temperature: true, Reasoning: true, ToolCall: true, Vision: true}, Limit: ModelLimit{Context: 200000, Output: 8192}},
		{ID: "claude-3-5-sonnet", Name: "Claude 3.5 Sonnet", Provider: "anthropic", OwnedBy: "anthropic", Created: 1714080000, Capabilities: ModelCapabilities{Temperature: true, Reasoning: true, ToolCall: true, Vision: true}, Limit: ModelLimit{Context: 200000, Output: 8192}},
		{ID: "claude-3-5-haiku", Name: "Claude 3.5 Haiku", Provider: "anthropic", OwnedBy: "anthropic", Created: 1714080000, Capabilities: ModelCapabilities{Temperature: true, Reasoning: false, ToolCall: true, Vision: true}, Limit: ModelLimit{Context: 200000, Output: 8192}},
		{ID: "claude-3-opus", Name: "Claude 3 Opus", Provider: "anthropic", OwnedBy: "anthropic", Created: 1709251200, Capabilities: ModelCapabilities{Temperature: true, Reasoning: true, ToolCall: true, Vision: true}, Limit: ModelLimit{Context: 200000, Output: 4096}},
		{ID: "claude-3-sonnet", Name: "Claude 3 Sonnet", Provider: "anthropic", OwnedBy: "anthropic", Created: 1709251200, Capabilities: ModelCapabilities{Temperature: true, Reasoning: true, ToolCall: true, Vision: true}, Limit: ModelLimit{Context: 200000, Output: 4096}},
		{ID: "claude-3-haiku", Name: "Claude 3 Haiku", Provider: "anthropic", OwnedBy: "anthropic", Created: 1709251200, Capabilities: ModelCapabilities{Temperature: true, Reasoning: false, ToolCall: true, Vision: true}, Limit: ModelLimit{Context: 200000, Output: 4096}},
	}
}

type CatalogService struct {
	catalog    *ModelCatalog
	httpClient *http.Client
	cachePath  string
	mu         sync.Mutex
}

func NewCatalogService() *CatalogService {
	return &CatalogService{
		catalog: NewModelCatalog(),
		httpClient: &http.Client{Timeout: 30 * time.Second},
		cachePath: filepath.Join(cacheDir, "catalog.json"),
	}
}

func (s *CatalogService) DiscoverFromProvider(ctx context.Context, p Provider) error {
	models, err := p.ListModels(ctx)
	if err != nil {
		return err
	}

	pm := &ProviderModels{
		Provider:   p.Name(),
		Name:       p.Name(),
		Models:     models,
		LastUpdate: time.Now(),
	}

	if sp, ok := p.(*OpenAIProvider); ok {
		pm.BaseURL = sp.BaseURL
		pm.APIKey = sp.APIKey
		pm.APIType = "openai"
	} else if sp, ok := p.(*AnthropicProvider); ok {
		pm.BaseURL = sp.BaseURL
		pm.APIKey = sp.APIKey
		pm.APIType = "anthropic"
	} else if sp, ok := p.(*MinimaxProvider); ok {
		pm.BaseURL = sp.BaseURL
		pm.APIKey = sp.APIKey
		pm.APIType = "minimax"
	} else if sp, ok := p.(*LiteLLMProvider); ok {
		pm.BaseURL = sp.BaseURL
		pm.APIKey = sp.APIKey
		pm.APIType = "litellm"
	}

	s.catalog.SetProvider(p.Name(), pm)
	return nil
}

func (s *CatalogService) GetCatalog() *ModelCatalog {
	return s.catalog
}

type ConfiguredModel struct {
	Name string
}

func (s *CatalogService) MergeWithConfig(configured map[string][]ConfiguredModel) {
	for providerName, modelCfgs := range configured {
		pm := s.catalog.GetProvider(providerName)
		if pm == nil {
			pm = &ProviderModels{Provider: providerName, Name: providerName}
			s.catalog.SetProvider(providerName, pm)
		}

		for _, cfg := range modelCfgs {
			model := Model{
				ID:   cfg.Name,
				Name: cfg.Name,
			}
			pm.Models = append(pm.Models, model)
		}
	}
}

func (s *CatalogService) Save() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data := s.catalog.AllProviders()
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.cachePath, jsonData, 0644)
}

func (s *CatalogService) Load() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := os.ReadFile(s.cachePath)
	if err != nil {
		return err
	}

	var providers map[string]*ProviderModels
	if err := json.Unmarshal(data, &providers); err != nil {
		return err
	}

	for name, pm := range providers {
		s.catalog.SetProvider(name, pm)
	}
	return nil
}