package opencode

import (
	"testing"
)

func TestOpenCodeConfigStruct(t *testing.T) {
	cfg := &OpenCodeConfig{
		Version:     "1.0",
		Model:       "gpt-4",
		Provider:    "openai",
		APIKey:      "test-key",
		Shell:       "/bin/bash",
		ContextSize: 4096,
	}

	if cfg.Version != "1.0" {
		t.Errorf("Version = %q, want %q", cfg.Version, "1.0")
	}
	if cfg.Model != "gpt-4" {
		t.Errorf("Model = %q, want %q", cfg.Model, "gpt-4")
	}
	if cfg.Provider != "openai" {
		t.Errorf("Provider = %q, want %q", cfg.Provider, "openai")
	}
	if cfg.Shell != "/bin/bash" {
		t.Errorf("Shell = %q, want %q", cfg.Shell, "/bin/bash")
	}
	if cfg.ContextSize != 4096 {
		t.Errorf("ContextSize = %d, want %d", cfg.ContextSize, 4096)
	}
}

func TestMigrate(t *testing.T) {
	oc := &OpenCodeConfig{
		Model:    "gpt-4",
		Provider: "openai",
		APIKey:   "test-key",
		Shell:    "/bin/zsh",
		Tools:    []string{"bash", "edit"},
	}

	cfg := Migrate(oc)

	if cfg.Agent.Default != "gpt-4" {
		t.Errorf("Agent.Default = %q, want %q", cfg.Agent.Default, "gpt-4")
	}
	if cfg.Tools.Bash.Shell != "/bin/zsh" {
		t.Errorf("Tools.Bash.Shell = %q, want %q", cfg.Tools.Bash.Shell, "/bin/zsh")
	}
	if len(cfg.Tools.Allowed) != 2 {
		t.Errorf("Tools.Allowed length = %d, want %d", len(cfg.Tools.Allowed), 2)
	}
}

func TestMigrateWithContextSize(t *testing.T) {
	oc := &OpenCodeConfig{
		Model:       "gpt-4",
		Provider:    "openai",
		ContextSize: 8192,
	}

	cfg := Migrate(oc)

	modelCfg, ok := cfg.Models["default"]
	if !ok {
		t.Fatal("default model not set")
	}
	if modelCfg.Provider != "openai" {
		t.Errorf("Model.Provider = %q, want %q", modelCfg.Provider, "openai")
	}
	if modelCfg.Name != "gpt-4" {
		t.Errorf("Model.Name = %q, want %q", modelCfg.Name, "gpt-4")
	}
}

func TestMigrateEmptyConfig(t *testing.T) {
	oc := &OpenCodeConfig{}

	cfg := Migrate(oc)
	if cfg.Agent.Default == "" {
		t.Error("Agent.Default should have default value from config.DefaultConfig()")
	}
}

func TestMigrateFileNonexistent(t *testing.T) {
	_, err := MigrateFile("/nonexistent/path/config.json")
	if err == nil {
		t.Error("MigrateFile() expected error for nonexistent path")
	}
}

func TestReadNonexistentPath(t *testing.T) {
	_, err := Read("/nonexistent/path/config.json")
	if err == nil {
		t.Error("Read() expected error for nonexistent path")
	}
}

func TestDefaultOpenCodePath(t *testing.T) {
	path := defaultOpenCodePath()
	if path == "" {
		t.Log("No opencode config found (expected on clean system)")
	}
}