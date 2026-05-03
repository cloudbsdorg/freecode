package config

import (
	"fmt"
	"os"

	"github.com/pelletier/go-toml/v2"
)

type TOMLConfig struct {
	Version string                 `toml:"version"`
	Config  map[string]interface{} `toml:"config"`
}

func LoadTOMLFile(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read TOML file: %w", err)
	}

	return LoadTOML(data)
}

func LoadTOML(data []byte) (*Config, error) {
	var raw map[string]interface{}
	if err := toml.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("failed to parse TOML: %w", err)
	}

	cfg := DefaultConfig()

	if shell, ok := raw["shell"].(string); ok {
		cfg.Shell = shell
	}
	if yolo, ok := raw["yolo"].(bool); ok {
		cfg.Yolo = yolo
	}
	if logLevel, ok := raw["log_level"].(string); ok {
		cfg.LogLevel = logLevel
	}
	if timeout, ok := raw["timeout"].(float64); ok {
		cfg.Timeout = int(timeout)
	}

	return cfg, nil
}

func (c *Config) ToTOML() ([]byte, error) {
	return toml.Marshal(c)
}

func (c *Config) ToTOMLPretty() ([]byte, error) {
	return toml.Marshal(c)
}

func (c *Config) MergeTOML(data []byte) error {
	var overlay map[string]interface{}
	if err := toml.Unmarshal(data, &overlay); err != nil {
		return fmt.Errorf("failed to parse TOML overlay: %w", err)
	}

	merged := mergeConfigMap(c.toMap(), overlay)

	cfg := DefaultConfig()
	tomlData, _ := toml.Marshal(merged)
	if err := toml.Unmarshal(tomlData, cfg); err != nil {
		return fmt.Errorf("failed to apply merged config: %w", err)
	}

	*c = *cfg
	return nil
}
