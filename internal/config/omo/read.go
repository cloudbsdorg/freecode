package omo

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type OMOConfig struct {
	Version    string                 `mapstructure:"version"`
	SkillsDir  string                 `mapstructure:"skills_dir"`
	AgentsDir  string                 `mapstructure:"agents_dir"`
	Hooks      map[string]interface{} `mapstructure:"hooks"`
	Commands   map[string]interface{} `mapstructure:"commands"`
	Fleet      map[string]interface{} `mapstructure:"fleet"`
	SlopRemove bool                   `mapstructure:"slop_remove"`
}

type Loader interface {
	Read(path string) (*OMOConfig, error)
}

type ViperLoader struct{}

func NewViperLoader() *ViperLoader {
	return &ViperLoader{}
}

func (vl *ViperLoader) Read(path string) (*OMOConfig, error) {
	cfg := &OMOConfig{}

	if path == "" {
		path = defaultOMOPath()
		if path == "" {
			return cfg, nil
		}
	}

	v := viper.New()
	v.SetConfigFile(path)
	v.SetConfigType("jsonc")

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read OMO config: %w", err)
		}
		return cfg, nil
	}

	if err := v.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal OMO config: %w", err)
	}

	return cfg, nil
}

type MockLoader struct {
	Config *OMOConfig
	Err    error
}

func (ml *MockLoader) Read(path string) (*OMOConfig, error) {
	if ml.Err != nil {
		return nil, ml.Err
	}
	return ml.Config, nil
}

func Read(path string) (*OMOConfig, error) {
	loader := NewViperLoader()
	return loader.Read(path)
}

func defaultOMOPath() string {
	homeDir, _ := os.UserHomeDir()
	paths := []string{
		filepath.Join(homeDir, ".config", "oh-my-openagent", "oh-my-openagent.jsonc"),
		filepath.Join(homeDir, ".config", "oh-my-opencode", "oh-my-opencode.jsonc"),
	}

	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	return ""
}