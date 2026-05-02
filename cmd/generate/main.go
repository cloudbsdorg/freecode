package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

const (
	modelsDevURL = "https://models.dev/api.json"
	outputFile   = "internal/provider/registry.json"
	timeout      = 30 * time.Second
)

type Provider struct {
	ID     string             `json:"id"`
	API    string             `json:"api,omitempty"`
	Name   string             `json:"name"`
	Env    []string           `json:"env"`
	NPM    string             `json:"npm,omitempty"`
	Models map[string]*Model `json:"models"`
	Doc    string             `json:"doc,omitempty"`
}

type Model struct {
	ID          string       `json:"id"`
	Name        string       `json:"name"`
	Family      string       `json:"family,omitempty"`
	Attachment  bool         `json:"attachment"`
	Reasoning   bool         `json:"reasoning"`
	ToolCall    bool         `json:"tool_call"`
	Temperature bool         `json:"temperature"`
	Modalities *Modalities `json:"modalities,omitempty"`
	Cost       *Cost       `json:"cost,omitempty"`
	Limit      *Limit      `json:"limit"`
	ReleaseDate string      `json:"release_date,omitempty"`
	LastUpdated string      `json:"last_updated,omitempty"`
}

type Modalities struct {
	Input  []string `json:"input,omitempty"`
	Output []string `json:"output,omitempty"`
}

type Cost struct {
	Input  float64 `json:"input"`
	Output float64 `json:"output"`
}

type Limit struct {
	Context int `json:"context"`
	Input   int `json:"input,omitempty"`
	Output  int `json:"output"`
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	data, err := fetchModelsDev()
	if err != nil {
		return fmt.Errorf("failed to fetch models.dev: %w", err)
	}

	var providers map[string]*Provider
	if err := json.Unmarshal(data, &providers); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	fmt.Printf("Fetched %d providers from models.dev\n", len(providers))

	generated, err := generateJSON(providers)
	if err != nil {
		return fmt.Errorf("failed to generate JSON: %w", err)
	}

	if err := os.WriteFile(outputFile, []byte(generated), 0644); err != nil {
		return fmt.Errorf("failed to write %s: %w", outputFile, err)
	}

	fmt.Printf("Generated %s\n", outputFile)
	return nil
}

func generateJSON(providers map[string]*Provider) (string, error) {
	type registryEntry struct {
		ID      string   `json:"id"`
		Name    string   `json:"name"`
		BaseURL string   `json:"baseURL"`
		EnvVars []string `json:"envVars"`
	}

	type modelEntry struct {
		ID          string  `json:"id"`
		Name        string  `json:"name"`
		Context     int     `json:"context"`
		OutputLimit int     `json:"outputLimit"`
		InputCost  float64 `json:"inputCost"`
		OutputCost float64 `json:"outputCost"`
		HasVision  bool    `json:"hasVision"`
		HasAudio   bool    `json:"hasAudio"`
	}

	result := struct {
		Providers map[string]registryEntry         `json:"providers"`
		Models   map[string]map[string]modelEntry `json:"models"`
	}{
		Providers: make(map[string]registryEntry),
		Models:   make(map[string]map[string]modelEntry),
	}

	for id, prov := range providers {
		result.Providers[id] = registryEntry{
			ID:      id,
			Name:    prov.Name,
			BaseURL: prov.API,
			EnvVars: prov.Env,
		}

		result.Models[id] = make(map[string]modelEntry)
		for mid, m := range prov.Models {
			if m == nil {
				continue
			}
			contextLimit := 128000
			if m.Limit != nil {
				contextLimit = m.Limit.Context
			}
			outputLimit := 4096
			if m.Limit != nil {
				outputLimit = m.Limit.Output
			}
			inputCost := 0.0
			if m.Cost != nil {
				inputCost = m.Cost.Input
			}
			outputCost := 0.0
			if m.Cost != nil {
				outputCost = m.Cost.Output
			}
			result.Models[id][mid] = modelEntry{
				ID:          m.ID,
				Name:        m.Name,
				Context:     contextLimit,
				OutputLimit: outputLimit,
				InputCost:  inputCost,
				OutputCost: outputCost,
				HasVision:  hasModality(m.Modalities, "image"),
				HasAudio:   hasModality(m.Modalities, "audio"),
			}
		}
	}

	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func fetchModelsDev() ([]byte, error) {
	client := &http.Client{Timeout: timeout}
	req, err := http.NewRequest("GET", modelsDevURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "freecode-generator/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

func hasModality(m *Modalities, modality string) bool {
	if m == nil {
		return false
	}
	for _, in := range m.Input {
		if in == modality {
			return true
		}
	}
	return false
}
