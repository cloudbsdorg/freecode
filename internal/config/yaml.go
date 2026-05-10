package config

import (
	"bytes"
	"fmt"
	"os"

	"github.com/spf13/viper"
)

func LoadYAMLFile(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read YAML file: %w", err)
	}

	return LoadYAML(data)
}

func LoadYAML(data []byte) (*Config, error) {
	v := viper.New()
	v.SetConfigType("yaml")

	if err := v.ReadConfig(bytes.NewReader(data)); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	cfg := DefaultConfig()

	cfg.Shell = v.GetString("shell")
	cfg.Yolo = v.GetBool("yolo")
	cfg.LogLevel = v.GetString("log_level")
	cfg.Timeout = v.GetInt("timeout")
	cfg.Editor = v.GetString("editor")
	cfg.Pager = v.GetString("pager")
	cfg.Theme = v.GetString("theme")

	if v.IsSet("models") {
		cfg.Models = make(map[string]ModelConfig)
		if err := v.UnmarshalKey("models", &cfg.Models); err != nil {
			return nil, fmt.Errorf("failed to unmarshal models: %w", err)
		}
	}
	if v.IsSet("providers") {
		cfg.Providers = make(map[string]ProviderConfig)
		if err := v.UnmarshalKey("providers", &cfg.Providers); err != nil {
			return nil, fmt.Errorf("failed to unmarshal providers: %w", err)
		}
	}

	if v.IsSet("openai.api_key") {
		cfg.OpenAI.APIKey = v.GetString("openai.api_key")
		cfg.OpenAI.BaseURL = v.GetString("openai.base_url")
	}
	if v.IsSet("anthropic.api_key") {
		cfg.Anthropic.APIKey = v.GetString("anthropic.api_key")
		cfg.Anthropic.BaseURL = v.GetString("anthropic.base_url")
	}
	if v.IsSet("ollama.api_key") {
		cfg.Ollama.APIKey = v.GetString("ollama.api_key")
		cfg.Ollama.BaseURL = v.GetString("ollama.base_url")
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return cfg, nil
}

func (c *Config) SaveYAML(path string) error {
	v := viper.New()
	v.SetConfigType("yaml")

	v.Set("shell", c.Shell)
	v.Set("yolo", c.Yolo)
	v.Set("log_level", c.LogLevel)
	v.Set("timeout", c.Timeout)
	v.Set("editor", c.Editor)
	v.Set("pager", c.Pager)
	v.Set("theme", c.Theme)

	if len(c.Models) > 0 {
		v.Set("models", c.Models)
	}
	if len(c.Providers) > 0 {
		v.Set("providers", c.Providers)
	}

	if c.OpenAI.APIKey != "" || c.OpenAI.BaseURL != "" {
		v.Set("openai", c.OpenAI)
	}
	if c.Anthropic.APIKey != "" || c.Anthropic.BaseURL != "" {
		v.Set("anthropic", c.Anthropic)
	}
	if c.Ollama.APIKey != "" || c.Ollama.BaseURL != "" {
		v.Set("ollama", c.Ollama)
	}
	if c.Minimax.APIKey != "" || c.Minimax.BaseURL != "" {
		v.Set("minimax", c.Minimax)
	}

	v.Set("tools.tool_states", c.Tools.ToolStates)

	return v.WriteConfigAs(path)
}

func (c *Config) Save() error {
	return c.SaveYAML("")
}
