package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type JSONConfig struct {
	Version string                 `json:"version"`
	Config  map[string]interface{} `json:"config"`
}

func LoadJSONFile(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read JSON file: %w", err)
	}

	return LoadJSON(data)
}

func LoadJSON(data []byte) (*Config, error) {
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
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

func (c *Config) ToJSON() ([]byte, error) {
	return json.MarshalIndent(c, "", "  ")
}

func (c *Config) ToJSONCompact() ([]byte, error) {
	return json.Marshal(c)
}

func (c *Config) MergeJSON(data []byte) error {
	var overlay map[string]interface{}
	if err := json.Unmarshal(data, &overlay); err != nil {
		return fmt.Errorf("failed to parse JSON overlay: %w", err)
	}

	merged := mergeConfigMap(c.toMap(), overlay)

	cfg := DefaultConfig()
	jsonData, _ := json.Marshal(merged)
	if err := json.Unmarshal(jsonData, cfg); err != nil {
		return fmt.Errorf("failed to apply merged config: %w", err)
	}

	*c = *cfg
	return nil
}

func (c *Config) toMap() map[string]interface{} {
	data, _ := json.Marshal(c)
	var result map[string]interface{}
	json.Unmarshal(data, &result)
	return result
}

func mergeConfigMap(base, overlay map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	for k, v := range base {
		result[k] = v
	}

	for k, v := range overlay {
		if ov, ok := v.(map[string]interface{}); ok {
			if bv, ok := result[k].(map[string]interface{}); ok {
				result[k] = mergeConfigMap(bv, ov)
			} else {
				result[k] = ov
			}
		} else {
			result[k] = v
		}
	}

	return result
}
