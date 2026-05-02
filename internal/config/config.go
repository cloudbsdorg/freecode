package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

type Config struct {
	Shell       string            `mapstructure:"shell"`
	LogLevel    string            `mapstructure:"log_level"`
	LogFormat   string            `mapstructure:"log_format"`
	Yolo        bool              `mapstructure:"yolo"`
	Quiet       bool              `mapstructure:"quiet"`
	Verbose     bool              `mapstructure:"verbose"`
	Debug       bool              `mapstructure:"debug"`
	Color       bool              `mapstructure:"color"`
	Theme       string            `mapstructure:"theme"`
	Width       int               `mapstructure:"width"`
	Height      int               `mapstructure:"height"`
	Editor      string            `mapstructure:"editor"`
	Pager       string            `mapstructure:"pager"`
	HTTPProxy   string            `mapstructure:"http_proxy"`
	HTTPSProxy  string            `mapstructure:"https_proxy"`
	NoProxy     string            `mapstructure:"no_proxy"`
	Timeout     int               `mapstructure:"timeout"`
	Retries     int               `mapstructure:"retries"`
	Server      ServerConfig      `mapstructure:"server"`
	Agent       AgentConfig       `mapstructure:"agent"`
	Models      map[string]ModelConfig `mapstructure:"models"`
	Providers   map[string]ProviderConfig `mapstructure:"providers"`
	Tools       ToolsConfig       `mapstructure:"tools"`
	Hooks       HooksConfig      `mapstructure:"hooks"`
	Permissions PermissionsConfig `mapstructure:"permissions"`
	Platform    PlatformConfig   `mapstructure:"platform"`
	Session     SessionConfig    `mapstructure:"session"`
	LiteLLM     LiteLLMConfig    `mapstructure:"litellm"`
}

type LiteLLMConfig struct {
	BaseURL string `mapstructure:"base_url"`
	APIKey  string `mapstructure:"api_key"`
}

type ServerConfig struct {
	Host       string            `mapstructure:"host"`
	Port       int               `mapstructure:"port"`
	MCPHost    string            `mapstructure:"mcp_host"`
	MCPPort    int               `mapstructure:"mcp_port"`
	WebUIHost  string            `mapstructure:"webui_host"`
	WebUIPort  int               `mapstructure:"webui_port"`
	TLS        TLSConfig         `mapstructure:"tls"`
	Auth       AuthConfig        `mapstructure:"auth"`
	CORS       CORSConfig        `mapstructure:"cors"`
	RateLimit  RateLimitConfig   `mapstructure:"rate_limit"`
}

type TLSConfig struct {
	Enabled  bool   `mapstructure:"enabled"`
	CertFile string `mapstructure:"cert_file"`
	KeyFile  string `mapstructure:"key_file"`
}

type AuthConfig struct {
	Enabled  bool   `mapstructure:"enabled"`
	Password string `mapstructure:"password"`
	Token    string `mapstructure:"token"`
}

type CORSConfig struct {
	Enabled    bool     `mapstructure:"enabled"`
	Origins    []string `mapstructure:"origins"`
	Methods    []string `mapstructure:"methods"`
	Headers    []string `mapstructure:"headers"`
	MaxAge     int      `mapstructure:"max_age"`
}

type RateLimitConfig struct {
	Enabled       bool `mapstructure:"enabled"`
	RequestsPerMin int `mapstructure:"requests_per_min"`
	Burst         int  `mapstructure:"burst"`
}

type AgentConfig struct {
	Default      string            `mapstructure:"default"`
	Thinking     bool              `mapstructure:"thinking"`
	MaxTurns     int               `mapstructure:"max_turns"`
	Timeout      int               `mapstructure:"timeout"`
	Fallback     []string          `mapstructure:"fallback"`
	Capabilities CapabilitiesConfig `mapstructure:"capabilities"`
}

type CapabilitiesConfig struct {
	Vision       bool `mapstructure:"vision"`
	FunctionCall bool `mapstructure:"function_call"`
	Streaming     bool `mapstructure:"streaming"`
}

type ModelConfig struct {
	Provider   string                 `mapstructure:"provider"`
	Name       string                 `mapstructure:"name"`
	Variant    string                 `mapstructure:"variant"`
	APIKey     string                 `mapstructure:"api_key"`
	BaseURL    string                 `mapstructure:"base_url"`
	Timeout    int                    `mapstructure:"timeout"`
	MaxTokens  int                    `mapstructure:"max_tokens"`
	Temperature float64               `mapstructure:"temperature"`
	TopP       float64               `mapstructure:"top_p"`
	Tools      []string               `mapstructure:"tools"`
	Thinking   ThinkingConfig         `mapstructure:"thinking"`
}

type ThinkingConfig struct {
	Enabled bool    `mapstructure:"enabled"`
	Budget  int     `mapstructure:"budget"`
}

type ProviderConfig struct {
	Type     string            `mapstructure:"type"`
	APIKey   string            `mapstructure:"api_key"`
	BaseURL  string            `mapstructure:"base_url"`
	Models   []string         `mapstructure:"models"`
	Headers  map[string]string `mapstructure:"headers"`
}

type ToolsConfig struct {
	Bash       BashToolsConfig       `mapstructure:"bash"`
	Allowed    []string             `mapstructure:"allowed"`
	Denied     []string             `mapstructure:"denied"`
	MaxDepth   int                  `mapstructure:"max_depth"`
	Timeout    int                  `mapstructure:"timeout"`
	ASTGrep    ASTGrepConfig        `mapstructure:"ast_grep"`
	LSP        LSPConfig            `mapstructure:"lsp"`
}

type BashToolsConfig struct {
	Shell    string   `mapstructure:"shell"`
	Timeout  int      `mapstructure:"timeout"`
	Env      []string `mapstructure:"env"`
	WorkDir  string   `mapstructure:"work_dir"`
}

type ASTGrepConfig struct {
	Enabled bool              `mapstructure:"enabled"`
	Languages []string        `mapstructure:"languages"`
}

type LSPConfig struct {
	Enabled  bool                   `mapstructure:"enabled"`
	Servers  map[string]LSPServerConfig `mapstructure:"servers"`
	AutoStart bool                 `mapstructure:"auto_start"`
}

type LSPServerConfig struct {
	Command    []string `mapstructure:"command"`
	RootURI    string   `mapstructure:"root_uri"`
	LanguageID string   `mapstructure:"language_id"`
}

type HooksConfig struct {
	Session     []string `mapstructure:"session"`
	Tool        []string `mapstructure:"tool"`
	Transform   []string `mapstructure:"transform"`
	Continuation []string `mapstructure:"continuation"`
	Ralph       []string `mapstructure:"ralph"`
	Skill       []string `mapstructure:"skill"`
}

type PermissionsConfig struct {
	AllowedDirs    []string                  `mapstructure:"allowed_dirs"`
	DeniedDirs     []string                  `mapstructure:"denied_dirs"`
	AllowedHosts   []string                  `mapstructure:"allowed_hosts"`
	DeniedHosts    []string                  `mapstructure:"denied_hosts"`
	ToolPatterns   map[string][]string       `mapstructure:"tool_patterns"`
	AgentRules     map[string]AgentRuleConfig `mapstructure:"agent_rules"`
}

type AgentRuleConfig struct {
	Allow []string `mapstructure:"allow"`
	Deny  []string `mapstructure:"deny"`
}

type PlatformConfig struct {
	OS           string                 `mapstructure:"os"`
	Arch         string                 `mapstructure:"arch"`
	HomeDir      string                 `mapstructure:"home_dir"`
	TempDir      string                 `mapstructure:"temp_dir"`
	CacheDir     string                 `mapstructure:"cache_dir"`
	FreeBSD      FreeBSDConfig          `mapstructure:"freebsd"`
	Darwin       DarwinConfig           `mapstructure:"darwin"`
	Linux        LinuxConfig            `mapstructure:"linux"`
}

type FreeBSDConfig struct {
	UseBaseGit bool `mapstructure:"use_base_git"`
}

type DarwinConfig struct {
	UseGitCredential bool `mapstructure:"use_git_credential"`
}

type LinuxConfig struct {
	XDGConfigHome string `mapstructure:"xdg_config_home"`
}

type SessionConfig struct {
	Dir         string            `mapstructure:"dir"`
	Compaction  CompactionConfig  `mapstructure:"compaction"`
	HistorySize int               `mapstructure:"history_size"`
}

type CompactionConfig struct {
	Enabled        bool    `mapstructure:"enabled"`
	ThresholdTokens int    `mapstructure:"threshold_tokens"`
	KeepLastN      int     `mapstructure:"keep_last_n"`
}

func DefaultConfig() *Config {
	homeDir, _ := os.UserHomeDir()
	tempDir := os.TempDir()
	cacheDir := filepath.Join(tempDir, "freecode-cache")

	return &Config{
		Shell:     "/bin/bash",
		LogLevel:  "info",
		LogFormat: "text",
		Yolo:      true,
		Color:     true,
		Theme:     "default",
		Width:     120,
		Height:    80,
		Editor:    os.Getenv("EDITOR"),
		Pager:     os.Getenv("PAGER"),
		Timeout:   60,
		Retries:   3,
		Server: ServerConfig{
			Host:      "127.0.0.1",
			Port:      18792,
			MCPHost:   "127.0.0.1",
			MCPPort:   18793,
			WebUIHost: "127.0.0.1",
			WebUIPort: 18791,
		},
		Agent: AgentConfig{
			Default:  "sisyphus",
			Thinking: false,
			MaxTurns: 100,
			Timeout:  120,
			Capabilities: CapabilitiesConfig{
				Vision:       true,
				FunctionCall: true,
				Streaming:    true,
			},
		},
		Models:    make(map[string]ModelConfig),
		Providers: make(map[string]ProviderConfig),
		Tools: ToolsConfig{
			Allowed:  []string{"bash", "read", "write", "edit", "glob", "grep", "webfetch", "websearch"},
			MaxDepth: 10,
			Timeout:  30,
			Bash: BashToolsConfig{
				Shell:   "/bin/bash",
				Timeout: 60,
			},
		},
		Session: SessionConfig{
			Dir:         filepath.Join(homeDir, ".local", "share", "freecode"),
			HistorySize: 1000,
			Compaction: CompactionConfig{
				Enabled:         true,
				ThresholdTokens: 100000,
				KeepLastN:       10,
			},
		},
		Platform: PlatformConfig{
			OS:       runtime.GOOS,
			Arch:     runtime.GOARCH,
			HomeDir:  homeDir,
			TempDir:  tempDir,
			CacheDir: cacheDir,
		},
		LiteLLM: LiteLLMConfig{
			BaseURL: "http://localhost:4000",
			APIKey:  "local",
		},
	}
}

func (c *Config) Validate() error {
	if c.Server.Port < 1 || c.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", c.Server.Port)
	}
	if c.Timeout < 1 {
		return fmt.Errorf("timeout must be positive")
	}
	if c.Session.HistorySize < 0 {
		return fmt.Errorf("history_size must be non-negative")
	}
	return nil
}
