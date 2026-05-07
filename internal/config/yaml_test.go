package config

import (
	"testing"
)

func TestLoadYAML(t *testing.T) {
	data := []byte(`
shell: /bin/zsh
yolo: true
log_level: debug
timeout: 120
editor: vim
pager: less
theme: dark
`)

	cfg, err := LoadYAML(data)
	if err != nil {
		t.Fatalf("LoadYAML() error = %v", err)
	}

	if cfg.Shell != "/bin/zsh" {
		t.Errorf("Shell = %q, want %q", cfg.Shell, "/bin/zsh")
	}
	if cfg.Yolo != true {
		t.Errorf("Yolo = %v, want true", cfg.Yolo)
	}
	if cfg.LogLevel != "debug" {
		t.Errorf("LogLevel = %q, want %q", cfg.LogLevel, "debug")
	}
	if cfg.Timeout != 120 {
		t.Errorf("Timeout = %d, want %d", cfg.Timeout, 120)
	}
}

func TestLoadYAMLInvalid(t *testing.T) {
	data := []byte(`invalid: yaml: content: {`)

	_, err := LoadYAML(data)
	if err == nil {
		t.Error("LoadYAML() should error for invalid YAML")
	}
}

func TestLoadYAMLFile(t *testing.T) {
	cfg, err := LoadYAMLFile("/nonexistent/file.yaml")
	if err == nil {
		t.Error("LoadYAMLFile() should error for nonexistent file")
	}
	if cfg != nil {
		t.Error("LoadYAMLFile() should return nil config for nonexistent file")
	}
}

func TestConfigSaveYAML(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Shell = "/bin/zsh"
	cfg.Yolo = true
	cfg.LogLevel = "debug"

	tmpFile := t.TempDir() + "/config.yaml"
	err := cfg.SaveYAML(tmpFile)
	if err != nil {
		t.Fatalf("SaveYAML() error = %v", err)
	}

	cfg2, err := LoadYAMLFile(tmpFile)
	if err != nil {
		t.Fatalf("LoadYAMLFile() error = %v", err)
	}

	if cfg2.Shell != "/bin/zsh" {
		t.Errorf("Shell = %q, want %q", cfg2.Shell, "/bin/zsh")
	}
}
