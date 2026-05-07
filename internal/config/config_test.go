package config

import (
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg == nil {
		t.Fatal("DefaultConfig() returned nil")
	}

	if cfg.Shell != "/bin/bash" {
		t.Errorf("DefaultConfig().Shell = %q, want %q", cfg.Shell, "/bin/bash")
	}

	if cfg.LogLevel != "info" {
		t.Errorf("DefaultConfig().LogLevel = %q, want %q", cfg.LogLevel, "info")
	}

	if cfg.Yolo != true {
		t.Errorf("DefaultConfig().Yolo = %v, want true", cfg.Yolo)
	}

	if cfg.Timeout != 60 {
		t.Errorf("DefaultConfig().Timeout = %d, want %d", cfg.Timeout, 60)
	}
}

func TestConfigValidate(t *testing.T) {
	tests := []struct {
		name    string
		cfg     Config
		wantErr bool
	}{
		{
			name: "valid config",
			cfg: Config{
				Server:  ServerConfig{Port: 18792},
				Timeout: 60,
				Session: SessionConfig{HistorySize: 100},
			},
			wantErr: false,
		},
		{
			name: "invalid server port - zero",
			cfg: Config{
				Server: ServerConfig{Port: 0},
			},
			wantErr: true,
		},
		{
			name: "invalid server port - negative",
			cfg: Config{
				Server: ServerConfig{Port: -1},
			},
			wantErr: true,
		},
		{
			name: "invalid server port - too high",
			cfg: Config{
				Server: ServerConfig{Port: 70000},
			},
			wantErr: true,
		},
		{
			name: "invalid timeout - zero",
			cfg: Config{
				Server:  ServerConfig{Port: 18792},
				Timeout: 0,
			},
			wantErr: true,
		},
		{
			name: "invalid timeout - negative",
			cfg: Config{
				Server:  ServerConfig{Port: 18792},
				Timeout: -1,
			},
			wantErr: true,
		},
		{
			name: "invalid history size - negative",
			cfg: Config{
				Server:  ServerConfig{Port: 18792},
				Timeout: 60,
				Session: SessionConfig{HistorySize: -1},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Config.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConfigServerDefaults(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.Server.Host != "127.0.0.1" {
		t.Errorf("Server.Host = %q, want %q", cfg.Server.Host, "127.0.0.1")
	}

	if cfg.Server.Port != 18792 {
		t.Errorf("Server.Port = %d, want %d", cfg.Server.Port, 18792)
	}

	if cfg.Server.MCPHost != "127.0.0.1" {
		t.Errorf("Server.MCPHost = %q, want %q", cfg.Server.MCPHost, "127.0.0.1")
	}

	if cfg.Server.MCPPort != 18793 {
		t.Errorf("Server.MCPPort = %d, want %d", cfg.Server.MCPPort, 18793)
	}

	if cfg.Server.WebUIHost != "127.0.0.1" {
		t.Errorf("Server.WebUIHost = %q, want %q", cfg.Server.WebUIHost, "127.0.0.1")
	}

	if cfg.Server.WebUIPort != 18791 {
		t.Errorf("Server.WebUIPort = %d, want %d", cfg.Server.WebUIPort, 18791)
	}
}

func TestConfigAgentDefaults(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.Agent.Default != "sisyphus" {
		t.Errorf("Agent.Default = %q, want %q", cfg.Agent.Default, "sisyphus")
	}

	if cfg.Agent.Thinking != false {
		t.Errorf("Agent.Thinking = %v, want false", cfg.Agent.Thinking)
	}

	if cfg.Agent.MaxTurns != 100 {
		t.Errorf("Agent.MaxTurns = %d, want %d", cfg.Agent.MaxTurns, 100)
	}

	if cfg.Agent.Timeout != 120 {
		t.Errorf("Agent.Timeout = %d, want %d", cfg.Agent.Timeout, 120)
	}
}

func TestConfigToolsDefaults(t *testing.T) {
	cfg := DefaultConfig()

	if len(cfg.Tools.Allowed) == 0 {
		t.Error("Tools.Allowed is empty")
	}

	expectedTools := []string{"bash", "read", "write", "edit", "glob", "grep", "webfetch", "websearch"}
	for _, tool := range expectedTools {
		found := false
		for _, allowed := range cfg.Tools.Allowed {
			if allowed == tool {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Tools.Allowed missing %q", tool)
		}
	}

	if cfg.Tools.MaxDepth != 10 {
		t.Errorf("Tools.MaxDepth = %d, want %d", cfg.Tools.MaxDepth, 10)
	}

	if cfg.Tools.Timeout != 30 {
		t.Errorf("Tools.Timeout = %d, want %d", cfg.Tools.Timeout, 30)
	}
}

func TestConfigSessionDefaults(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.Session.HistorySize != 1000 {
		t.Errorf("Session.HistorySize = %d, want %d", cfg.Session.HistorySize, 1000)
	}

	if !cfg.Session.Compaction.Enabled {
		t.Error("Session.Compaction.Enabled should be true")
	}

	if cfg.Session.Compaction.ThresholdTokens != 100000 {
		t.Errorf("Session.Compaction.ThresholdTokens = %d, want %d",
			cfg.Session.Compaction.ThresholdTokens, 100000)
	}

	if cfg.Session.Compaction.KeepLastN != 10 {
		t.Errorf("Session.Compaction.KeepLastN = %d, want %d",
			cfg.Session.Compaction.KeepLastN, 10)
	}
}

func TestConfigPlatform(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.Platform.OS == "" {
		t.Error("Platform.OS is empty")
	}

	if cfg.Platform.Arch == "" {
		t.Error("Platform.Arch is empty")
	}

	if cfg.Platform.HomeDir == "" {
		t.Error("Platform.HomeDir is empty")
	}

	if cfg.Platform.TempDir == "" {
		t.Error("Platform.TempDir is empty")
	}
}
