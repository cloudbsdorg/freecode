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

func Read(path string) (*OMOConfig, error) {
	cfg := &OMOConfig{}

	if path == "" {
		path = defaultOMOPath()
	}

	viper.SetConfigFile(path)
	viper.SetConfigType("jsonc")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read OMO config: %w", err)
		}
	}

	if err := viper.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal OMO config: %w", err)
	}

	return cfg, nil
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
