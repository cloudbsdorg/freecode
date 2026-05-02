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

	return v.WriteConfigAs(path)
}
