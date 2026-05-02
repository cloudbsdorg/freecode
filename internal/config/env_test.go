package config

import (
	"os"
	"testing"
)

func TestGetEnv(t *testing.T) {
	os.Setenv("FREECODE_SHELL", "/bin/tcsh")
	defer os.Unsetenv("FREECODE_SHELL")

	got := getEnv("SHELL", "/bin/bash")
	want := "/bin/tcsh"

	if got != want {
		t.Errorf("getEnv() = %q, want %q", got, want)
	}
}

func TestGetEnvDefault(t *testing.T) {
	os.Unsetenv("FREECODE_SHELL")

	got := getEnv("SHELL", "/bin/bash")
	want := "/bin/bash"

	if got != want {
		t.Errorf("getEnv() = %q, want %q", got, want)
	}
}

func TestGetEnvBool(t *testing.T) {
	tests := []struct {
		name    string
		envVal  string
		def     bool
		want    bool
	}{
		{"true string", "true", false, true},
		{"TRUE uppercase", "TRUE", false, true},
		{"1 digit", "1", false, true},
		{"false string", "false", true, false},
		{"0 digit", "0", true, false},
		{"empty uses default", "", true, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envVal != "" {
				os.Setenv("FREECODE_TEST_BOOL", tt.envVal)
				defer os.Unsetenv("FREECODE_TEST_BOOL")
			} else {
				os.Unsetenv("FREECODE_TEST_BOOL")
			}

			got := getEnvBool("TEST_BOOL", tt.def)
			if got != tt.want {
				t.Errorf("getEnvBool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetEnvInt(t *testing.T) {
	tests := []struct {
		name   string
		envVal string
		def    int
		want   int
	}{
		{"valid number", "300", 60, 300},
		{"zero", "0", 60, 0},
		{"invalid uses default", "notanumber", 60, 60},
		{"empty uses default", "", 60, 60},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envVal != "" {
				os.Setenv("FREECODE_TEST_INT", tt.envVal)
				defer os.Unsetenv("FREECODE_TEST_INT")
			} else {
				os.Unsetenv("FREECODE_TEST_INT")
			}

			got := getEnvInt("TEST_INT", tt.def)
			if got != tt.want {
				t.Errorf("getEnvInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestApplyEnvOverrides(t *testing.T) {
	os.Setenv("FREECODE_SHELL", "/bin/zsh")
	os.Setenv("FREECODE_LOG_LEVEL", "debug")
	os.Setenv("FREECODE_YOLO", "false")
	os.Setenv("FREECODE_TIMEOUT", "180")
	os.Setenv("FREECODE_SERVER_HOST", "0.0.0.0")
	os.Setenv("FREECODE_SERVER_PORT", "19999")
	os.Setenv("FREECODE_AGENT_DEFAULT", "oracle")
	defer func() {
		os.Unsetenv("FREECODE_SHELL")
		os.Unsetenv("FREECODE_LOG_LEVEL")
		os.Unsetenv("FREECODE_YOLO")
		os.Unsetenv("FREECODE_TIMEOUT")
		os.Unsetenv("FREECODE_SERVER_HOST")
		os.Unsetenv("FREECODE_SERVER_PORT")
		os.Unsetenv("FREECODE_AGENT_DEFAULT")
	}()

	cfg := DefaultConfig()
	cfg.ApplyEnvOverrides()

	if cfg.Shell != "/bin/zsh" {
		t.Errorf("Shell = %q, want %q", cfg.Shell, "/bin/zsh")
	}

	if cfg.LogLevel != "debug" {
		t.Errorf("LogLevel = %q, want %q", cfg.LogLevel, "debug")
	}

	if cfg.Yolo != false {
		t.Errorf("Yolo = %v, want false", cfg.Yolo)
	}

	if cfg.Timeout != 180 {
		t.Errorf("Timeout = %d, want %d", cfg.Timeout, 180)
	}

	if cfg.Server.Host != "0.0.0.0" {
		t.Errorf("Server.Host = %q, want %q", cfg.Server.Host, "0.0.0.0")
	}

	if cfg.Server.Port != 19999 {
		t.Errorf("Server.Port = %d, want %d", cfg.Server.Port, 19999)
	}

	if cfg.Agent.Default != "oracle" {
		t.Errorf("Agent.Default = %q, want %q", cfg.Agent.Default, "oracle")
	}
}
