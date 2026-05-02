package provider

import (
	_ "embed"
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
)

type RegistryLoader struct {
	once     sync.Once
	registry *Registry
	err      error
}

type Registry struct {
	Providers map[string]ProviderInfo `json:"providers"`
	Models    map[string]map[string]ModelInfo `json:"models"`
}

type ProviderInfo struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	BaseURL string   `json:"baseURL"`
	EnvVars []string `json:"envVars"`
}

type ModelInfo struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Context     int     `json:"context"`
	OutputLimit int     `json:"outputLimit"`
	InputCost  float64 `json:"inputCost"`
	OutputCost float64 `json:"outputCost"`
	HasVision  bool    `json:"hasVision"`
	HasAudio   bool    `json:"hasAudio"`
}

var (
	registryLoader RegistryLoader
)

func LoadRegistry() (*Registry, error) {
	registryLoader.once.Do(func() {
		registryLoader.registry, registryLoader.err = loadRegistryFromFile()
		if registryLoader.err != nil {
			registryLoader.registry, registryLoader.err = loadRegistryFromEmbed()
		}
	})
	return registryLoader.registry, registryLoader.err
}

func loadRegistryFromFile() (*Registry, error) {
	exePath, err := os.Executable()
	if err != nil {
		return nil, err
	}
	registryPath := filepath.Join(filepath.Dir(exePath), "internal", "provider", "registry.json")
	data, err := os.ReadFile(registryPath)
	if err != nil {
		return nil, err
	}
	return parseRegistry(data)
}

func loadRegistryFromEmbed() (*Registry, error) {
	data, err := os.ReadFile("internal/provider/registry.json")
	if err != nil {
		return nil, err
	}
	return parseRegistry(data)
}

func parseRegistry(data []byte) (*Registry, error) {
	var reg Registry
	if err := json.Unmarshal(data, &reg); err != nil {
		return nil, err
	}
	return &reg, nil
}

func GetModelInfo(providerID, modelID string) (ModelInfo, bool) {
	reg, err := LoadRegistry()
	if err != nil {
		return ModelInfo{}, false
	}
	models, ok := reg.Models[providerID]
	if !ok {
		return ModelInfo{}, false
	}
	model, ok := models[modelID]
	return model, ok
}

func GetProviderInfo(providerID string) (ProviderInfo, bool) {
	reg, err := LoadRegistry()
	if err != nil {
		return ProviderInfo{}, false
	}
	prov, ok := reg.Providers[providerID]
	return prov, ok
}

func GetProviderByModelPrefix(prefix string) string {
	reg, err := LoadRegistry()
	if err != nil {
		return ""
	}
	prefixLower := prefix
	for providerID := range reg.Providers {
		for modelID := range reg.Models[providerID] {
			if len(modelID) >= len(prefixLower) && modelID[:len(prefixLower)] == prefixLower {
				return providerID
			}
		}
	}
	return ""
}

func ListProviders() []string {
	reg, err := LoadRegistry()
	if err != nil {
		return nil
	}
	ids := make([]string, 0, len(reg.Providers))
	for id := range reg.Providers {
		ids = append(ids, id)
	}
	return ids
}
