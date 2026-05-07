package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

type Loader interface {
	Read(path string) (*Config, error)
}

type ViperLoader struct{}

func NewViperLoader() *ViperLoader {
	return &ViperLoader{}
}

func (vl *ViperLoader) Read(path string) (*Config, error) {
	v := viper.New()

	if path != "" {
		v.SetConfigFile(path)
	} else {
		configPaths := configSearchPaths()
		for _, p := range configPaths {
			if _, err := os.Stat(p); err == nil {
				v.SetConfigFile(p)
				break
			}
		}
		v.SetConfigType("yaml")
	}

	v.SetEnvPrefix("FREECODE")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config: %w", err)
		}
	}

	cfg := DefaultConfig()
	if err := v.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	cfg.applyDefaults()
	cfg.applyEnvOverrides()

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return cfg, nil
}

type MockLoader struct {
	Config *Config
	Err    error
}

func (ml *MockLoader) Read(path string) (*Config, error) {
	if ml.Err != nil {
		return nil, ml.Err
	}
	return ml.Config, nil
}

func Load(path string) (*Config, error) {
	loader := NewViperLoader()
	return loader.Read(path)
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
