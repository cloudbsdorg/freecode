package opencode

import (
	"github.com/freecode/freecode/internal/config"
)

func Migrate(oc *OpenCodeConfig) *config.Config {
	cfg := config.DefaultConfig()

	if oc.Model != "" {
		cfg.Agent.Default = oc.Model
	}

	if oc.Provider != "" {
		cfg.Providers["default"] = config.ProviderConfig{
			Type:   oc.Provider,
			APIKey: oc.APIKey,
			BaseURL: oc.BaseURL,
		}
	}

	if oc.Shell != "" {
		cfg.Tools.Bash.Shell = oc.Shell
	}

	if oc.ContextSize > 0 {
		cfg.Models["default"] = config.ModelConfig{
			Provider: oc.Provider,
			Name:     oc.Model,
		}
	}

	cfg.Tools.Allowed = oc.Tools

	return cfg
}

func MigrateFile(opencodePath string) (*config.Config, error) {
	oc, err := Read(opencodePath)
	if err != nil {
		return nil, err
	}
	return Migrate(oc), nil
}
