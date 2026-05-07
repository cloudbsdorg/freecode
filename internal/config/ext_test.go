package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestConfig_applyDefaults(t *testing.T) {
	cfg := &Config{}
	cfg.applyDefaults()
	if cfg.Shell != "/bin/bash" {
		t.Fatalf("Default Shell = %q, want %q", cfg.Shell, "/bin/bash")
	}
	if cfg.LogLevel != "info" {
		t.Fatalf("Default LogLevel = %q, want %q", cfg.LogLevel, "info")
	}
	if cfg.Timeout != 60 {
		t.Fatalf("Default Timeout = %d, want %d", cfg.Timeout, 60)
	}
	if cfg.Session.Dir == "" {
		t.Fatalf("Default Session.Dir should be set, got empty string")
	}
}

func TestConfig_applyEnvOverrides(t *testing.T) {
	// Only a subset of env vars are used by ApplyEnvOverrides
	os.Setenv("FREECODE_SHELL", "/bin/zsh")
	os.Setenv("FREECODE_LOG_LEVEL", "debug")
	os.Setenv("FREECODE_TIMEOUT", "45")
	defer func() {
		os.Unsetenv("FREECODE_SHELL")
		os.Unsetenv("FREECODE_LOG_LEVEL")
		os.Unsetenv("FREECODE_TIMEOUT")
	}()

	cfg := DefaultConfig()
	cfg.ApplyEnvOverrides()

	if cfg.Shell != "/bin/zsh" {
		t.Fatalf("Shell = %q, want %q", cfg.Shell, "/bin/zsh")
	}
	if cfg.LogLevel != "debug" {
		t.Fatalf("LogLevel = %q, want %q", cfg.LogLevel, "debug")
	}
	if cfg.Timeout != 45 {
		t.Fatalf("Timeout = %d, want %d", cfg.Timeout, 45)
	}
}

func TestConfigMergeJSON_Nested(t *testing.T) {
	cfg := DefaultConfig()
	overlay := `{"server":{"port":12345},"tools":{"bash":{"timeout":99}}}`
	if err := cfg.MergeJSON([]byte(overlay)); err != nil {
		t.Fatalf("MergeJSON() error: %v", err)
	}
	if cfg.Server.Port != 12345 {
		t.Fatalf("Server.Port = %d, want %d", cfg.Server.Port, 12345)
	}
	if cfg.Tools.Bash.Timeout != 99 {
		t.Fatalf("Tools.Bash.Timeout = %d, want %d", cfg.Tools.Bash.Timeout, 99)
	}
}

func TestConfigLoadYamlPath(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	content := []byte("shell: /bin/zsh\nlog_level: debug\nserver:\n  port: 3000\n")
	if err := os.WriteFile(path, content, 0644); err != nil {
		t.Fatalf("failed to write temp config: %v", err)
	}
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if cfg.Shell != "/bin/zsh" {
		t.Fatalf("Shell = %q, want %q", cfg.Shell, "/bin/zsh")
	}
	if cfg.LogLevel != "debug" {
		t.Fatalf("LogLevel = %q, want %q", cfg.LogLevel, "debug")
	}
	if cfg.Server.Port != 3000 {
		t.Fatalf("Server.Port = %d, want %d", cfg.Server.Port, 3000)
	}
}
