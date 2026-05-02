package opencode

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type OpenCodeConfig struct {
	Version     string                 `mapstructure:"version"`
	Model       string                 `mapstructure:"model"`
	Provider    string                 `mapstructure:"provider"`
	APIKey      string                 `mapstructure:"api_key"`
	BaseURL     string                 `mapstructure:"base_url"`
	Shell       string                 `mapstructure:"shell"`
	Prompt      string                 `mapstructure:"prompt"`
	ContextSize int                    `mapstructure:"context_size"`
	Tools       []string               `mapstructure:"tools"`
	Agents      map[string]interface{} `mapstructure:"agents"`
	Hooks       map[string]interface{} `mapstructure:"hooks"`
}

type Loader interface {
	Read(path string) (*OpenCodeConfig, error)
}

type ViperLoader struct{}

func NewViperLoader() *ViperLoader {
	return &ViperLoader{}
}

func (vl *ViperLoader) Read(path string) (*OpenCodeConfig, error) {
	cfg := &OpenCodeConfig{}

	if path == "" {
		path = defaultOpenCodePath()
		if path == "" {
			return cfg, nil
		}
	}

	ext := filepath.Ext(path)
	v := viper.New()
	v.SetConfigFile(path)

	switch ext {
	case ".json":
		v.SetConfigType("json")
	case ".jsonc":
		v.SetConfigType("json")
	case ".toml":
		v.SetConfigType("toml")
	case ".yaml", ".yml":
		v.SetConfigType("yaml")
	default:
		v.SetConfigType("yaml")
	}

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read opencode config: %w", err)
	}

	if err := v.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal opencode config: %w", err)
	}

	return cfg, nil
}

type MockLoader struct {
	Config *OpenCodeConfig
	Err    error
}

func (ml *MockLoader) Read(path string) (*OpenCodeConfig, error) {
	if ml.Err != nil {
		return nil, ml.Err
	}
	return ml.Config, nil
}

func Read(path string) (*OpenCodeConfig, error) {
	loader := NewViperLoader()
	return loader.Read(path)
}

func defaultOpenCodePath() string {
	homeDir, _ := os.UserHomeDir()
	paths := []string{
		filepath.Join(homeDir, ".config", "opencode", "config.json"),
		filepath.Join(homeDir, ".config", "opencode", "config.jsonc"),
		filepath.Join(homeDir, ".config", "opencode", "config.toml"),
		filepath.Join(homeDir, ".config", "opencode", "opencode.json"),
		filepath.Join(homeDir, ".config", "opencode", "opencode.jsonc"),
		filepath.Join(homeDir, ".config", "oh-my-opencode", "oh-my-opencode.jsonc"),
		filepath.Join(homeDir, ".config", "oh-my-openagent", "oh-my-openagent.jsonc"),
	}

	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	return ""
}