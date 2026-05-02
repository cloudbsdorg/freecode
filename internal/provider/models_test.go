package provider

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestModelCatalog_NewModelCatalog(t *testing.T) {
	c := NewModelCatalog()
	if c == nil {
		t.Fatal("NewModelCatalog returned nil")
	}
	if c.providers == nil {
		t.Error("providers map is nil")
	}
}

func TestModelCatalog_GetProvider(t *testing.T) {
	c := NewModelCatalog()
	c.SetProvider("test", &ProviderModels{Provider: "test", Name: "Test Provider"})

	got := c.GetProvider("test")
	if got == nil {
		t.Fatal("GetProvider returned nil")
	}
	if got.Provider != "test" {
		t.Errorf("Provider = %q, want %q", got.Provider, "test")
	}
}

func TestModelCatalog_GetProviderNotFound(t *testing.T) {
	c := NewModelCatalog()
	got := c.GetProvider("nonexistent")
	if got != nil {
		t.Errorf("GetProvider() = %v, want nil", got)
	}
}

func TestModelCatalog_SetProvider(t *testing.T) {
	c := NewModelCatalog()
	pm := &ProviderModels{Provider: "openai", Name: "OpenAI"}
	c.SetProvider("openai", pm)

	got := c.GetProvider("openai")
	if got == nil {
		t.Fatal("GetProvider returned nil")
	}
	if got.Name != "OpenAI" {
		t.Errorf("Name = %q, want %q", got.Name, "OpenAI")
	}
}

func TestModelCatalog_AllProviders(t *testing.T) {
	c := NewModelCatalog()
	c.SetProvider("p1", &ProviderModels{Provider: "p1"})
	c.SetProvider("p2", &ProviderModels{Provider: "p2"})

	result := c.AllProviders()
	if len(result) != 2 {
		t.Errorf("len(AllProviders()) = %d, want 2", len(result))
	}
}

func TestModelCatalog_AllProvidersEmpty(t *testing.T) {
	c := NewModelCatalog()
	result := c.AllProviders()
	if len(result) != 0 {
		t.Errorf("len(AllProviders()) = %d, want 0", len(result))
	}
}

func TestModelCatalog_GetModel(t *testing.T) {
	c := NewModelCatalog()
	c.SetProvider("anthropic", &ProviderModels{
		Provider: "anthropic",
		Models: []Model{
			{ID: "claude-3", Name: "Claude 3"},
		},
	})

	got := c.GetModel("anthropic", "claude-3")
	if got == nil {
		t.Fatal("GetModel returned nil")
	}
	if got.ID != "claude-3" {
		t.Errorf("ID = %q, want %q", got.ID, "claude-3")
	}
}

func TestModelCatalog_GetModelProviderNotFound(t *testing.T) {
	c := NewModelCatalog()
	got := c.GetModel("nonexistent", "model")
	if got != nil {
		t.Errorf("GetModel() = %v, want nil", got)
	}
}

func TestModelCatalog_GetModelNotFound(t *testing.T) {
	c := NewModelCatalog()
	c.SetProvider("anthropic", &ProviderModels{
		Provider: "anthropic",
		Models:   []Model{},
	})

	got := c.GetModel("anthropic", "nonexistent")
	if got != nil {
		t.Errorf("GetModel() = %v, want nil", got)
	}
}

func TestProviderModels_IsStale(t *testing.T) {
	pm := &ProviderModels{
		LastUpdate: time.Now().Add(-10 * time.Minute),
	}
	if !pm.IsStale() {
		t.Error("IsStale() = false, want true for old update")
	}

	pm.LastUpdate = time.Now()
	if pm.IsStale() {
		t.Error("IsStale() = true, want false for recent update")
	}
}

func TestFetchModelsFromURL(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer test-key" {
			t.Errorf("Missing or wrong Authorization header")
		}
		resp := `{"object":"list","data":[{"id":"gpt-4","object":"model","created":1739836800,"owned_by":"openai"}]}`
		w.Write([]byte(resp))
	}))
	defer server.Close()

	models, err := fetchModelsFromURL(context.Background(), server.URL, "test-key")
	if err != nil {
		t.Fatalf("fetchModelsFromURL() error = %v", err)
	}
	if len(models) != 1 {
		t.Fatalf("len(models) = %d, want 1", len(models))
	}
	if models[0].ID != "gpt-4" {
		t.Errorf("models[0].ID = %q, want %q", models[0].ID, "gpt-4")
	}
	if models[0].OwnedBy != "openai" {
		t.Errorf("models[0].OwnedBy = %q, want %q", models[0].OwnedBy, "openai")
	}
}

func TestFetchModelsFromURLBadResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("not json"))
	}))
	defer server.Close()

	_, err := fetchModelsFromURL(context.Background(), server.URL, "test-key")
	if err == nil {
		t.Error("fetchModelsFromURL() expected error for bad JSON")
	}
}

func TestFetchModelsFromURLHTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "server error", 500)
	}))
	defer server.Close()

	_, err := fetchModelsFromURL(context.Background(), server.URL, "test-key")
	if err == nil {
		t.Error("fetchModelsFromURL() expected error for HTTP error")
	}
}

func TestFetchModelsFromURLContextCanceled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := fetchModelsFromURL(ctx, "http://localhost:99999", "test-key")
	if err == nil {
		t.Error("fetchModelsFromURL() expected error for canceled context")
	}
}

func TestOpenAIProvider_ListModels(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := `{"object":"list","data":[{"id":"gpt-4","object":"model","created":1739836800,"owned_by":"openai"},{"id":"gpt-3.5","object":"model","created":1739836800,"owned_by":"openai"}]}`
		w.Write([]byte(resp))
	}))
	defer server.Close()

	p := &OpenAIProvider{
		APIKey:  "test-key",
		BaseURL: server.URL,
		Client:  server.Client(),
	}

	models, err := p.ListModels(context.Background())
	if err != nil {
		t.Fatalf("ListModels() error = %v", err)
	}
	if len(models) != 2 {
		t.Errorf("len(models) = %d, want 2", len(models))
	}
}

func TestAnthropicProvider_ListModels(t *testing.T) {
	p := NewAnthropicProvider("test-key")
	models, err := p.ListModels(context.Background())
	if err != nil {
		t.Fatalf("ListModels() error = %v", err)
	}
	if len(models) == 0 {
		t.Error("len(models) = 0, want > 0 for anthropic")
	}
}

func TestMinimaxProvider_ListModels(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := `{"object":"list","data":[{"id":"MiniMax-M2.7","object":"model","created":1739836800,"owned_by":"minimax"}]}`
		w.Write([]byte(resp))
	}))
	defer server.Close()

	p := &MinimaxProvider{
		APIKey:  "test-key",
		BaseURL: server.URL,
		Client:  server.Client(),
	}

	models, err := p.ListModels(context.Background())
	if err != nil {
		t.Fatalf("ListModels() error = %v", err)
	}
	if len(models) != 1 {
		t.Errorf("len(models) = %d, want 1", len(models))
	}
}

func TestLiteLLMProvider_ListModels(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := `{"object":"list","data":[{"id":"gpt-4","object":"model","created":1739836800,"owned_by":"openai"}]}`
		w.Write([]byte(resp))
	}))
	defer server.Close()

	p := &LiteLLMProvider{
		APIKey:  "test-key",
		BaseURL: server.URL,
		Client:  server.Client(),
	}

	models, err := p.ListModels(context.Background())
	if err != nil {
		t.Fatalf("ListModels() error = %v", err)
	}
	if len(models) != 1 {
		t.Errorf("len(models) = %d, want 1", len(models))
	}
}

func TestGetAnthropicModels(t *testing.T) {
	models := getAnthropicModels()
	if len(models) == 0 {
		t.Fatal("getAnthropicModels() returned empty slice")
	}
	for _, m := range models {
		if m.Provider != "anthropic" {
			t.Errorf("Provider = %q, want %q", m.Provider, "anthropic")
		}
		if m.ID == "" {
			t.Error("Model ID is empty")
		}
	}
}

func TestNewCatalogService(t *testing.T) {
	s := NewCatalogService()
	if s == nil {
		t.Fatal("NewCatalogService returned nil")
	}
	if s.catalog == nil {
		t.Error("catalog is nil")
	}
	if s.httpClient == nil {
		t.Error("httpClient is nil")
	}
}

func TestCatalogService_DiscoverFromProvider(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := `{"object":"list","data":[{"id":"gpt-4","object":"model","created":1739836800,"owned_by":"openai"}]}`
		w.Write([]byte(resp))
	}))
	defer server.Close()

	p := &OpenAIProvider{
		APIKey:  "test-key",
		BaseURL: server.URL,
		Client:  server.Client(),
	}

	s := NewCatalogService()
	err := s.DiscoverFromProvider(context.Background(), p)
	if err != nil {
		t.Fatalf("DiscoverFromProvider() error = %v", err)
	}

	catalog := s.GetCatalog()
	pm := catalog.GetProvider("openai")
	if pm == nil {
		t.Fatal("Provider not found in catalog")
	}
	if len(pm.Models) != 1 {
		t.Errorf("len(pm.Models) = %d, want 1", len(pm.Models))
	}
}

func TestCatalogService_DiscoverFromProviderAnthropic(t *testing.T) {
	p := NewAnthropicProvider("test-key")
	s := NewCatalogService()

	err := s.DiscoverFromProvider(context.Background(), p)
	if err != nil {
		t.Fatalf("DiscoverFromProvider() error = %v", err)
	}

	catalog := s.GetCatalog()
	pm := catalog.GetProvider("anthropic")
	if pm == nil {
		t.Fatal("Provider not found in catalog")
	}
	if len(pm.Models) == 0 {
		t.Error("len(pm.Models) = 0, want > 0 for anthropic")
	}
}

func TestCatalogService_DiscoverFromProviderMinimax(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := `{"object":"list","data":[{"id":"MiniMax-M2.7","object":"model","created":1739836800,"owned_by":"minimax"}]}`
		w.Write([]byte(resp))
	}))
	defer server.Close()

	p := &MinimaxProvider{
		APIKey:  "test-key",
		BaseURL: server.URL,
		Client:  server.Client(),
	}

	s := NewCatalogService()
	err := s.DiscoverFromProvider(context.Background(), p)
	if err != nil {
		t.Fatalf("DiscoverFromProvider() error = %v", err)
	}

	catalog := s.GetCatalog()
	pm := catalog.GetProvider("minimax")
	if pm == nil {
		t.Fatal("Provider not found in catalog")
	}
}

func TestCatalogService_DiscoverFromProviderLiteLLM(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := `{"object":"list","data":[{"id":"gpt-4","object":"model","created":1739836800,"owned_by":"openai"}]}`
		w.Write([]byte(resp))
	}))
	defer server.Close()

	p := &LiteLLMProvider{
		APIKey:  "test-key",
		BaseURL: server.URL,
		Client:  server.Client(),
	}

	s := NewCatalogService()
	err := s.DiscoverFromProvider(context.Background(), p)
	if err != nil {
		t.Fatalf("DiscoverFromProvider() error = %v", err)
	}

	catalog := s.GetCatalog()
	pm := catalog.GetProvider("litellm")
	if pm == nil {
		t.Fatal("Provider not found in catalog")
	}
}

func TestCatalogService_GetCatalog(t *testing.T) {
	s := NewCatalogService()
	catalog := s.GetCatalog()
	if catalog == nil {
		t.Fatal("GetCatalog() returned nil")
	}
}

func TestCatalogService_MergeWithConfig(t *testing.T) {
	c := NewModelCatalog()
	c.SetProvider("anthropic", &ProviderModels{
		Provider: "anthropic",
		Name:      "anthropic",
		Models: []Model{
			{ID: "claude-3", Name: "Claude 3"},
		},
	})

	s := &CatalogService{catalog: c}
	s.MergeWithConfig(map[string][]ConfiguredModel{
		"anthropic": {
			{Name: "claude-3.5"},
		},
		"newprovider": {
			{Name: "new-model"},
		},
	})

	pm := c.GetProvider("anthropic")
	if len(pm.Models) != 2 {
		t.Errorf("len(pm.Models) = %d, want 2", len(pm.Models))
	}

	pmNew := c.GetProvider("newprovider")
	if pmNew == nil {
		t.Fatal("newprovider not found in catalog")
	}
	if len(pmNew.Models) != 1 {
		t.Errorf("len(pmNew.Models) = %d, want 1", len(pmNew.Models))
	}
}

func TestCatalogService_SaveAndLoad(t *testing.T) {
	tmpDir := t.TempDir()
	cachePath := filepath.Join(tmpDir, "catalog.json")

	c := NewModelCatalog()
	c.SetProvider("anthropic", &ProviderModels{
		Provider:   "anthropic",
		Name:       "anthropic",
		LastUpdate: time.Now(),
		Models: []Model{
			{ID: "claude-3", Name: "Claude 3"},
		},
	})

	s := &CatalogService{
		catalog:    c,
		cachePath:  cachePath,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}

	err := s.Save()
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	c2 := NewModelCatalog()
	s2 := &CatalogService{
		catalog:    c2,
		cachePath:  cachePath,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}

	err = s2.Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	pm := c2.GetProvider("anthropic")
	if pm == nil {
		t.Fatal("Provider not found after load")
	}
	if len(pm.Models) != 1 {
		t.Errorf("len(pm.Models) = %d, want 1", len(pm.Models))
	}
}

func TestCatalogService_LoadFileNotFound(t *testing.T) {
	c := NewModelCatalog()
	s := &CatalogService{
		catalog:    c,
		cachePath:  "/nonexistent/path/catalog.json",
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}

	err := s.Load()
	if err == nil {
		t.Error("Load() expected error for nonexistent file")
	}
}

func TestCatalogService_SaveInvalidData(t *testing.T) {
	c := NewModelCatalog()
	cacheDir := t.TempDir()
	cachePath := filepath.Join(cacheDir, "catalog.json")

	s := &CatalogService{
		catalog:    c,
		cachePath:  cachePath,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}

	_ = os.MkdirAll(cachePath+".dir", 0755)
	_ = os.WriteFile(cachePath+".dir", []byte("not a file"), 0644)

	s.cachePath = cachePath + ".dir"
	err := s.Save()
	if err == nil {
		t.Error("Save() expected error when path is directory")
	}
}

func TestCatalogService_SaveUnwritablePath(t *testing.T) {
	c := NewModelCatalog()
	s := &CatalogService{
		catalog:    c,
		cachePath:  "/proc/fake/catalog.json",
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}

	err := s.Save()
	if err == nil {
		t.Error("Save() expected error for unwritable path")
	}
}

func TestCatalogService_LoadInvalidJSON(t *testing.T) {
	tmpDir := t.TempDir()
	cachePath := filepath.Join(tmpDir, "catalog.json")
	_ = os.WriteFile(cachePath, []byte("not valid json"), 0644)

	c := NewModelCatalog()
	s := &CatalogService{
		catalog:    c,
		cachePath:  cachePath,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}

	err := s.Load()
	if err == nil {
		t.Error("Load() expected error for invalid JSON")
	}
}

func TestOllamaProvider_ListModels(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := `{"models":[{"name":"llama3","model":"llama3","size":123}]}`
		w.Write([]byte(resp))
	}))
	defer server.Close()

	p := NewOllamaProvider(server.URL, "test-key")
	models, err := p.ListModels(context.Background())
	if err != nil {
		t.Fatalf("ListModels() error = %v", err)
	}
	if len(models) != 1 {
		t.Errorf("len(models) = %d, want 1", len(models))
	}
}

func TestNewOllamaProviderEmptyBaseURL(t *testing.T) {
	p := NewOllamaProvider("", "test-key")
	if p.BaseURL != "http://localhost:11434" {
		t.Errorf("BaseURL = %q, want %q", p.BaseURL, "http://localhost:11434")
	}
	if p.APIKey != "test-key" {
		t.Errorf("APIKey = %q, want %q", p.APIKey, "test-key")
	}
}

func TestOllamaProvider_ListModelsEmptyAPIKey(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := `{"models":[{"name":"llama3","model":"llama3","size":123}]}`
		w.Write([]byte(resp))
	}))
	defer server.Close()

	p := NewOllamaProvider(server.URL, "")
	_, err := p.ListModels(context.Background())
	if err != nil {
		t.Fatalf("ListModels() error = %v", err)
	}
}

func TestCatalogService_DiscoverFromProviderOllama(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := `{"models":[{"name":"llama3","model":"llama3","size":123}]}`
		w.Write([]byte(resp))
	}))
	defer server.Close()

	p := NewOllamaProvider(server.URL, "test-key")
	s := NewCatalogService()

	err := s.DiscoverFromProvider(context.Background(), p)
	if err != nil {
		t.Fatalf("DiscoverFromProvider() error = %v", err)
	}

	catalog := s.GetCatalog()
	pm := catalog.GetProvider("ollama")
	if pm == nil {
		t.Fatal("Provider not found in catalog")
	}
}

func TestOllamaProvider_ListModelsContextCanceled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	p := NewOllamaProvider("http://localhost:99999", "test-key")
	_, err := p.ListModels(ctx)
	if err == nil {
		t.Error("ListModels() expected error for canceled context")
	}
}

func TestOllamaProvider_GenerateContextCanceled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	p := NewOllamaProvider("http://localhost:99999", "test-key")
	_, err := p.Generate(ctx, &Request{
		Model: "llama3",
		Messages: []Message{
			{Role: "user", Content: "hi"},
		},
	})
	if err == nil {
		t.Error("Generate() expected error for canceled context")
	}
}

func TestOllamaProvider_GenerateServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "server error", 500)
	}))
	defer server.Close()

	p := NewOllamaProvider(server.URL, "test-key")
	_, err := p.Generate(context.Background(), &Request{
		Model: "llama3",
		Messages: []Message{
			{Role: "user", Content: "hi"},
		},
	})
	if err == nil {
		t.Error("Generate() expected error for server error")
	}
}

func TestOllamaProvider_GenerateBadJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("not json"))
	}))
	defer server.Close()

	p := NewOllamaProvider(server.URL, "test-key")
	_, err := p.Generate(context.Background(), &Request{
		Model: "llama3",
		Messages: []Message{
			{Role: "user", Content: "hi"},
		},
	})
	if err == nil {
		t.Error("Generate() expected error for bad JSON")
	}
}

func TestOllamaProvider_ListModelsServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "server error", 500)
	}))
	defer server.Close()

	p := NewOllamaProvider(server.URL, "test-key")
	_, err := p.ListModels(context.Background())
	if err == nil {
		t.Error("ListModels() expected error for server error")
	}
}

func TestOllamaProvider_ListModelsBadJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("not json"))
	}))
	defer server.Close()

	p := NewOllamaProvider(server.URL, "test-key")
	_, err := p.ListModels(context.Background())
	if err == nil {
		t.Error("ListModels() expected error for bad JSON")
	}
}

func TestOllamaProvider_GenerateSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]interface{}
		json.NewDecoder(r.Body).Decode(&body)
		if body["model"] != "llama3" {
			t.Errorf("model = %v, want llama3", body["model"])
		}
		resp := `{"message":{"content":"hello","role":"assistant"},"done":true,"eval_count":10,"prompt_eval_count":5}`
		w.Write([]byte(resp))
	}))
	defer server.Close()

	p := NewOllamaProvider(server.URL, "test-key")
	resp, err := p.Generate(context.Background(), &Request{
		Model: "llama3",
		Messages: []Message{
			{Role: "user", Content: "hi"},
		},
	})
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}
	if resp.Content != "hello" {
		t.Errorf("Content = %q, want %q", resp.Content, "hello")
	}
}