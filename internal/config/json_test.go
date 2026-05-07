package config

import (
	"testing"
)

func TestLoadJSON(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "valid minimal json",
			input:   `{"shell":"/bin/zsh","yolo":true,"log_level":"debug"}`,
			wantErr: false,
		},
		{
			name:    "empty input",
			input:   `{}`,
			wantErr: false,
		},
		{
			name:    "invalid json syntax",
			input:   `{"shell":}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := LoadJSON([]byte(tt.input))
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && cfg == nil {
				t.Error("LoadJSON() returned nil without error")
			}
		})
	}
}

func TestLoadJSONValues(t *testing.T) {
	input := `{"shell":"/bin/zsh","yolo":false,"log_level":"error","timeout":300}`

	cfg, err := LoadJSON([]byte(input))
	if err != nil {
		t.Fatalf("LoadJSON() error = %v", err)
	}

	if cfg.Shell != "/bin/zsh" {
		t.Errorf("Shell = %q, want %q", cfg.Shell, "/bin/zsh")
	}

	if cfg.Yolo != false {
		t.Errorf("Yolo = %v, want false", cfg.Yolo)
	}

	if cfg.LogLevel != "error" {
		t.Errorf("LogLevel = %q, want %q", cfg.LogLevel, "error")
	}

	if cfg.Timeout != 300 {
		t.Errorf("Timeout = %d, want %d", cfg.Timeout, 300)
	}
}

func TestConfigToJSON(t *testing.T) {
	cfg := DefaultConfig()

	data, err := cfg.ToJSON()
	if err != nil {
		t.Fatalf("ToJSON() error = %v", err)
	}

	if len(data) == 0 {
		t.Error("ToJSON() returned empty data")
	}
}

func TestConfigMergeJSON(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Shell = "/bin/bash"
	cfg.Timeout = 60

	overlay := `{"shell":"/bin/zsh","timeout":120}`

	err := cfg.MergeJSON([]byte(overlay))
	if err != nil {
		t.Fatalf("MergeJSON() error = %v", err)
	}

	if cfg.Shell != "/bin/zsh" {
		t.Errorf("Shell = %q, want %q", cfg.Shell, "/bin/zsh")
	}

	if cfg.Timeout != 120 {
		t.Errorf("Timeout = %d, want %d", cfg.Timeout, 120)
	}
}

func TestConfigToJSONCompact(t *testing.T) {
	cfg := DefaultConfig()

	data, err := cfg.ToJSONCompact()
	if err != nil {
		t.Fatalf("ToJSONCompact() error = %v", err)
	}

	if len(data) == 0 {
		t.Error("ToJSONCompact() returned empty data")
	}
}

func TestLoadJSONFile(t *testing.T) {
	cfg, err := LoadJSONFile("/nonexistent/file.json")
	if err == nil {
		t.Error("LoadJSONFile() should error for nonexistent file")
	}
	if cfg != nil {
		t.Error("LoadJSONFile() should return nil config for nonexistent file")
	}
}

func TestMergeJSONInvalid(t *testing.T) {
	cfg := DefaultConfig()
	err := cfg.MergeJSON([]byte(`{invalid`))
	if err == nil {
		t.Error("MergeJSON() should error for invalid JSON")
	}
}

func TestMergeConfigMap(t *testing.T) {
	base := map[string]interface{}{
		"a": 1,
		"b": map[string]interface{}{
			"c": 2,
		},
	}
	overlay := map[string]interface{}{
		"b": map[string]interface{}{
			"d": 3,
		},
		"e": 4,
	}

	result := mergeConfigMap(base, overlay)

	if result["a"] != 1 {
		t.Errorf("a = %v, want 1", result["a"])
	}
	if result["e"] != 4 {
		t.Errorf("e = %v, want 4", result["e"])
	}

	nested, ok := result["b"].(map[string]interface{})
	if !ok {
		t.Fatal("b should be map[string]interface{}")
	}
	if nested["c"] != 2 {
		t.Errorf("b.c = %v, want 2", nested["c"])
	}
	if nested["d"] != 3 {
		t.Errorf("b.d = %v, want 3", nested["d"])
	}
}
