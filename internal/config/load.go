package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

func Load(path string) (*Config, error) {
	cfg := DefaultConfig()

	if path != "" {
		viper.SetConfigFile(path)
	} else {
		configPaths := configSearchPaths()
		for _, p := range configPaths {
			if _, err := os.Stat(p); err == nil {
				viper.SetConfigFile(p)
				break
			}
		}
		viper.SetConfigType("yaml")
	}

	viper.AutomaticEnv()
	viper.SetEnvPrefix("FREECODE")
	viper.SetEnvKeyReplacer(cfgReplacer)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config: %w", err)
		}
	}

	if err := viper.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	cfg.applyDefaults()
	cfg.applyEnvOverrides()

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return cfg, nil
}

func configSearchPaths() []string {
	homeDir, _ := os.UserHomeDir()
	xdgConfigHome := os.Getenv("XDG_CONFIG_HOME")

	paths := []string{
		filepath.Join(homeDir, ".config", "freecode", "config.yaml"),
		filepath.Join(homeDir, ".config", "freecode", "config.yml"),
		filepath.Join(homeDir, ".config", "freecode", "config.json"),
		filepath.Join(xdgConfigHome, "freecode", "config.yaml"),
		filepath.Join(xdgConfigHome, "freecode", "config.yml"),
		filepath.Join(xdgConfigHome, "freecode", "config.json"),
		".freecode.yaml",
		".freecode.yml",
		".freecode.json",
	}

	return paths
}

var cfgReplacer = strings.NewReplacer(".", "_")

func (c *Config) applyDefaults() {
	if c.Shell == "" {
		c.Shell = "/bin/bash"
	}
	if c.LogLevel == "" {
		c.LogLevel = "info"
	}
	if c.Timeout == 0 {
		c.Timeout = 60
	}
	if c.Session.Dir == "" {
		homeDir, _ := os.UserHomeDir()
		c.Session.Dir = filepath.Join(homeDir, ".local", "share", "freecode")
	}
}

func (c *Config) applyEnvOverrides() {
	if dir := os.Getenv("FREECODE_DIR"); dir != "" {
		c.Session.Dir = dir
	}
	if proxy := os.Getenv("FREECODE_HTTP_PROXY"); proxy != "" {
		c.HTTPProxy = proxy
	}
	if logLevel := os.Getenv("FREECODE_LOG_LEVEL"); logLevel != "" {
		c.LogLevel = logLevel
	}
}

func Merge(other *Config) error {
	return nil
}
