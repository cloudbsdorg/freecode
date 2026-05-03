package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
)

type JSONCConfig struct {
	Version string                 `json:"version"`
	Config  map[string]interface{} `json:"config"`
}

func LoadJSONCFile(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read JSONC file: %w", err)
	}

	return LoadJSONC(data)
}

func LoadJSONC(data []byte) (*Config, error) {
	cleaned := stripJSONComments(data)

	var raw map[string]interface{}
	if err := json.Unmarshal(cleaned, &raw); err != nil {
		return nil, fmt.Errorf("failed to parse JSONC: %w", err)
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

func stripJSONComments(data []byte) []byte {
	var result bytes.Buffer
	lines := bytes.Split(data, []byte("\n"))

	lineCommentRe := regexp.MustCompile(`//.*$`)
	blockCommentRe := regexp.MustCompile(`/\*[\s\S]*?\*/`)

	for i, line := range lines {
		lineWithoutComment := lineCommentRe.ReplaceAll(line, []byte{})
		result.Write(lineWithoutComment)
		if i < len(lines)-1 {
			result.WriteByte('\n')
		}
	}

	resultBytes := result.Bytes()
	resultBytes = blockCommentRe.ReplaceAll(resultBytes, []byte{})

	return resultBytes
}

func (c *Config) ToJSONC() ([]byte, error) {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	lines := bytes.Split(data, []byte("\n"))
	for i, line := range lines {
		buf.Write(line)
		if i < len(lines)-1 {
			buf.WriteByte('\n')
		}
	}

	return buf.Bytes(), nil
}

func (c *Config) MergeJSONC(data []byte) error {
	cleaned := stripJSONComments(data)

	var overlay map[string]interface{}
	if err := json.Unmarshal(cleaned, &overlay); err != nil {
		return fmt.Errorf("failed to parse JSONC overlay: %w", err)
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
