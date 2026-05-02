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

func Read(path string) (*OpenCodeConfig, error) {
	cfg := &OpenCodeConfig{}

	if path == "" {
		path = defaultOpenCodePath()
	}

	ext := filepath.Ext(path)
	viper.SetConfigFile(path)

	switch ext {
	case ".json":
		viper.SetConfigType("json")
	case ".jsonc":
		viper.SetConfigType("json")
	case ".toml":
		viper.SetConfigType("toml")
	case ".yaml", ".yml":
		viper.SetConfigType("yaml")
	default:
		viper.SetConfigType("yaml")
	}

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read opencode config: %w", err)
	}

	if err := viper.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal opencode config: %w", err)
	}

	return cfg, nil
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
