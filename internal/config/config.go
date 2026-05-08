package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

type Config struct {
	Shell         string                    `mapstructure:"shell"`
	LogLevel      string                    `mapstructure:"log_level"`
	LogFormat     string                    `mapstructure:"log_format"`
	Yolo          bool                      `mapstructure:"yolo"`
	Quiet         bool                      `mapstructure:"quiet"`
	Verbose       bool                      `mapstructure:"verbose"`
	Debug         bool                      `mapstructure:"debug"`
	Color         bool                      `mapstructure:"color"`
	SoundEnabled  bool                      `mapstructure:"sound_enabled"`
	Theme         string                    `mapstructure:"theme"`
	Animation     string                    `mapstructure:"animation"`
	Width         int                       `mapstructure:"width"`
	Height        int                       `mapstructure:"height"`
	Editor        string                    `mapstructure:"editor"`
	Pager         string                    `mapstructure:"pager"`
	HTTPProxy     string                    `mapstructure:"http_proxy"`
	HTTPSProxy    string                    `mapstructure:"https_proxy"`
	NoProxy       string                    `mapstructure:"no_proxy"`
	Timeout       int                       `mapstructure:"timeout"`
	Retries       int                       `mapstructure:"retries"`
	Server        ServerConfig              `mapstructure:"server"`
	Agent         AgentConfig               `mapstructure:"agent"`
	Models        map[string]ModelConfig    `mapstructure:"models"`
	Providers     map[string]ProviderConfig `mapstructure:"providers"`
	Tools         ToolsConfig               `mapstructure:"tools"`
	Hooks         HooksConfig               `mapstructure:"hooks"`
	Permissions   PermissionsConfig         `mapstructure:"permissions"`
	Platform      PlatformConfig            `mapstructure:"platform"`
	Session       SessionConfig             `mapstructure:"session"`
	LiteLLM       LiteLLMConfig             `mapstructure:"litellm"`
	OpenAI        OpenAIConfig              `mapstructure:"openai"`
	Anthropic     AnthropicConfig           `mapstructure:"anthropic"`
	Ollama        OllamaConfig              `mapstructure:"ollama"`
	Azure         AzureConfig               `mapstructure:"azure"`
	Google        GoogleConfig              `mapstructure:"google"`
	Vertex        VertexConfig              `mapstructure:"vertex"`
	AWS           AWSConfig                 `mapstructure:"aws"`
	GitLab        GitLabConfig              `mapstructure:"gitlab"`
	GitHubCopilot GitHubCopilotConfig       `mapstructure:"github_copilot"`
	Vercel        VercelConfig              `mapstructure:"vercel"`
	Groq          GroqConfig                `mapstructure:"groq"`
	Perplexity    PerplexityConfig          `mapstructure:"perplexity"`
	Mistral       MistralConfig             `mapstructure:"mistral"`
	Cohere        CohereConfig              `mapstructure:"cohere"`
	TogetherAI    TogetherAIConfig          `mapstructure:"togetherai"`
	DeepInfra     DeepInfraConfig           `mapstructure:"deepinfra"`
	Cerebras      CerebrasConfig            `mapstructure:"cerebras"`
	XAI           XAIConfig                 `mapstructure:"xai"`
	Alibaba       AlibabaConfig             `mapstructure:"alibaba"`
	HuggingFace   HuggingFaceConfig         `mapstructure:"huggingface"`
	DeepSeek      DeepSeekConfig            `mapstructure:"deepseek"`
	Fireworks     FireworksConfig           `mapstructure:"fireworks"`
	Moonshot      MoonshotConfig            `mapstructure:"moonshot"`
	Nebius        NebiusConfig              `mapstructure:"nebius"`
	OpenRouter    OpenRouterConfig          `mapstructure:"openrouter"`
	Venice        VeniceConfig              `mapstructure:"venice"`
	ZAI           ZAIConfig                 `mapstructure:"zai"`
	ZenMux        ZenMuxConfig              `mapstructure:"zenmux"`
	Baseten       BasetenConfig             `mapstructure:"baseten"`
	Cortecs       CortecsConfig             `mapstructure:"cortecs"`
	Firmware      FirmwareConfig            `mapstructure:"firmware"`
	Ionet         IonetConfig               `mapstructure:"ionet"`
	NVIDIA        NVIDIAConfig              `mapstructure:"nvidia"`
	OllamaCloud   OllamaCloudConfig         `mapstructure:"ollamacloud"`
	Cloudflare    CloudflareConfig          `mapstructure:"cloudflare"`
	Helicone      HeliconeConfig            `mapstructure:"helicone"`
	LlamaCpp      LlamaCppConfig            `mapstructure:"llamacpp"`
	LMStudio      LMStudioConfig            `mapstructure:"lmstudio"`
	AtomicChat    AtomicChatConfig          `mapstructure:"atomic_chat"`
	Provider302AI Provider302AIConfig       `mapstructure:"302ai"`
	SAPAI         SAPAIConfig               `mapstructure:"sap_ai_core"`
	STACKIT       STACKITConfig             `mapstructure:"stackit"`
	OVHcloud      OVHcloudConfig            `mapstructure:"ovhcloud"`
	Scaleway      ScalewayConfig            `mapstructure:"scaleway"`
	Minimax       MinimaxConfig             `mapstructure:"minimax"`
}

type LiteLLMConfig struct {
	BaseURL string `mapstructure:"base_url"`
	APIKey  string `mapstructure:"api_key"`
}

type OpenAIConfig struct {
	BaseURL string `mapstructure:"base_url"`
	APIKey  string `mapstructure:"api_key"`
}

type AnthropicConfig struct {
	BaseURL string `mapstructure:"base_url"`
	APIKey  string `mapstructure:"api_key"`
}

type OllamaConfig struct {
	BaseURL string `mapstructure:"base_url"`
	APIKey  string `mapstructure:"api_key"`
}

type AzureConfig struct {
	BaseURL string `mapstructure:"base_url"`
	APIKey  string `mapstructure:"api_key"`
}

type GoogleConfig struct {
	BaseURL string `mapstructure:"base_url"`
	APIKey  string `mapstructure:"api_key"`
}

type VertexConfig struct {
	ProjectID   string `mapstructure:"project_id"`
	Location    string `mapstructure:"location"`
	BaseURL     string `mapstructure:"base_url"`
	AccessToken string `mapstructure:"access_token"`
}

type AWSConfig struct {
	Region    string `mapstructure:"region"`
	Profile   string `mapstructure:"profile"`
	Endpoint  string `mapstructure:"endpoint"`
	AccessKey string `mapstructure:"access_key_id"`
	SecretKey string `mapstructure:"secret_access_key"`
}

type GitLabConfig struct {
	BaseURL string `mapstructure:"base_url"`
	Token   string `mapstructure:"token"`
}

type GitHubCopilotConfig struct {
	BaseURL string `mapstructure:"base_url"`
	Token   string `mapstructure:"token"`
}

type VercelConfig struct {
	BaseURL string `mapstructure:"base_url"`
	Token   string `mapstructure:"token"`
}

type GroqConfig struct {
	BaseURL string `mapstructure:"base_url"`
	APIKey  string `mapstructure:"api_key"`
}

type PerplexityConfig struct {
	BaseURL string `mapstructure:"base_url"`
	APIKey  string `mapstructure:"api_key"`
}

type MistralConfig struct {
	BaseURL string `mapstructure:"base_url"`
	APIKey  string `mapstructure:"api_key"`
}

type CohereConfig struct {
	BaseURL string `mapstructure:"base_url"`
	APIKey  string `mapstructure:"api_key"`
}

type TogetherAIConfig struct {
	BaseURL string `mapstructure:"base_url"`
	APIKey  string `mapstructure:"api_key"`
}

type DeepInfraConfig struct {
	BaseURL string `mapstructure:"base_url"`
	APIKey  string `mapstructure:"api_key"`
}

type CerebrasConfig struct {
	BaseURL string `mapstructure:"base_url"`
	APIKey  string `mapstructure:"api_key"`
}

type XAIConfig struct {
	BaseURL string `mapstructure:"base_url"`
	APIKey  string `mapstructure:"api_key"`
}

type AlibabaConfig struct {
	BaseURL string `mapstructure:"base_url"`
	APIKey  string `mapstructure:"api_key"`
}

type HuggingFaceConfig struct {
	BaseURL string `mapstructure:"base_url"`
	APIKey  string `mapstructure:"api_key"`
}

type DeepSeekConfig struct {
	BaseURL string `mapstructure:"base_url"`
	APIKey  string `mapstructure:"api_key"`
}

type FireworksConfig struct {
	BaseURL string `mapstructure:"base_url"`
	APIKey  string `mapstructure:"api_key"`
}

type MoonshotConfig struct {
	BaseURL string `mapstructure:"base_url"`
	APIKey  string `mapstructure:"api_key"`
}

type NebiusConfig struct {
	BaseURL string `mapstructure:"base_url"`
	APIKey  string `mapstructure:"api_key"`
}

type OpenRouterConfig struct {
	BaseURL string `mapstructure:"base_url"`
	APIKey  string `mapstructure:"api_key"`
}

type VeniceConfig struct {
	BaseURL string `mapstructure:"base_url"`
	APIKey  string `mapstructure:"api_key"`
}

type ZAIConfig struct {
	BaseURL string `mapstructure:"base_url"`
	APIKey  string `mapstructure:"api_key"`
}

type ZenMuxConfig struct {
	BaseURL string `mapstructure:"base_url"`
	APIKey  string `mapstructure:"api_key"`
}

type BasetenConfig struct {
	BaseURL string `mapstructure:"base_url"`
	APIKey  string `mapstructure:"api_key"`
}

type CortecsConfig struct {
	BaseURL string `mapstructure:"base_url"`
	APIKey  string `mapstructure:"api_key"`
}

type FirmwareConfig struct {
	BaseURL string `mapstructure:"base_url"`
	APIKey  string `mapstructure:"api_key"`
}

type IonetConfig struct {
	BaseURL string `mapstructure:"base_url"`
	APIKey  string `mapstructure:"api_key"`
}

type NVIDIAConfig struct {
	BaseURL string `mapstructure:"base_url"`
	APIKey  string `mapstructure:"api_key"`
}

type OllamaCloudConfig struct {
	BaseURL string `mapstructure:"base_url"`
	APIKey  string `mapstructure:"api_key"`
}

type CloudflareConfig struct {
	AccountID string `mapstructure:"account_id"`
	GatewayID string `mapstructure:"gateway_id"`
	APIKey    string `mapstructure:"api_key"`
	BaseURL   string `mapstructure:"base_url"`
}

type HeliconeConfig struct {
	BaseURL string `mapstructure:"base_url"`
	APIKey  string `mapstructure:"api_key"`
}

type LlamaCppConfig struct {
	BaseURL string `mapstructure:"base_url"`
	APIKey  string `mapstructure:"api_key"`
}

type LMStudioConfig struct {
	BaseURL string `mapstructure:"base_url"`
	APIKey  string `mapstructure:"api_key"`
}

type AtomicChatConfig struct {
	BaseURL string `mapstructure:"base_url"`
	APIKey  string `mapstructure:"api_key"`
}

type Provider302AIConfig struct {
	BaseURL string `mapstructure:"base_url"`
	APIKey  string `mapstructure:"api_key"`
}

type SAPAIConfig struct {
	ServiceKey string `mapstructure:"service_key"`
	BaseURL    string `mapstructure:"base_url"`
}

type STACKITConfig struct {
	Token   string `mapstructure:"token"`
	BaseURL string `mapstructure:"base_url"`
}

type OVHcloudConfig struct {
	BaseURL string `mapstructure:"base_url"`
	APIKey  string `mapstructure:"api_key"`
}

type ScalewayConfig struct {
	BaseURL string `mapstructure:"base_url"`
	APIKey  string `mapstructure:"api_key"`
}

type MinimaxConfig struct {
	BaseURL string `mapstructure:"base_url"`
	APIKey  string `mapstructure:"api_key"`
}

type ServerConfig struct {
	Host      string          `mapstructure:"host"`
	Port      int             `mapstructure:"port"`
	MCPHost   string          `mapstructure:"mcp_host"`
	MCPPort   int             `mapstructure:"mcp_port"`
	WebUIHost string          `mapstructure:"webui_host"`
	WebUIPort int             `mapstructure:"webui_port"`
	TLS       TLSConfig       `mapstructure:"tls"`
	Auth      AuthConfig      `mapstructure:"auth"`
	CORS      CORSConfig      `mapstructure:"cors"`
	RateLimit RateLimitConfig `mapstructure:"rate_limit"`
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
	Enabled bool     `mapstructure:"enabled"`
	Origins []string `mapstructure:"origins"`
	Methods []string `mapstructure:"methods"`
	Headers []string `mapstructure:"headers"`
	MaxAge  int      `mapstructure:"max_age"`
}

type RateLimitConfig struct {
	Enabled        bool `mapstructure:"enabled"`
	RequestsPerMin int  `mapstructure:"requests_per_min"`
	Burst          int  `mapstructure:"burst"`
}

type AgentConfig struct {
	Default      string             `mapstructure:"default"`
	Thinking     bool               `mapstructure:"thinking"`
	MaxTurns     int                `mapstructure:"max_turns"`
	Timeout      int                `mapstructure:"timeout"`
	Fallback     []string           `mapstructure:"fallback"`
	Capabilities CapabilitiesConfig `mapstructure:"capabilities"`
}

type CapabilitiesConfig struct {
	Vision       bool `mapstructure:"vision"`
	FunctionCall bool `mapstructure:"function_call"`
	Streaming    bool `mapstructure:"streaming"`
}

type ModelConfig struct {
	Provider    string         `mapstructure:"provider"`
	Name        string         `mapstructure:"name"`
	Variant     string         `mapstructure:"variant"`
	APIKey      string         `mapstructure:"api_key"`
	BaseURL     string         `mapstructure:"base_url"`
	Timeout     int            `mapstructure:"timeout"`
	MaxTokens   int            `mapstructure:"max_tokens"`
	Temperature float64        `mapstructure:"temperature"`
	TopP        float64        `mapstructure:"top_p"`
	Tools       []string       `mapstructure:"tools"`
	Thinking    ThinkingConfig `mapstructure:"thinking"`
}

type ThinkingConfig struct {
	Enabled bool `mapstructure:"enabled"`
	Budget  int  `mapstructure:"budget"`
}

type ProviderConfig struct {
	Type    string            `mapstructure:"type"`
	APIKey  string            `mapstructure:"api_key"`
	BaseURL string            `mapstructure:"base_url"`
	Models  []string          `mapstructure:"models"`
	Headers map[string]string `mapstructure:"headers"`
}

type ToolsConfig struct {
	Bash     BashToolsConfig `mapstructure:"bash"`
	Allowed  []string        `mapstructure:"allowed"`
	Denied   []string        `mapstructure:"denied"`
	MaxDepth int             `mapstructure:"max_depth"`
	Timeout  int             `mapstructure:"timeout"`
	ASTGrep  ASTGrepConfig   `mapstructure:"ast_grep"`
	LSP      LSPConfig       `mapstructure:"lsp"`
}

type BashToolsConfig struct {
	Shell   string   `mapstructure:"shell"`
	Timeout int      `mapstructure:"timeout"`
	Env     []string `mapstructure:"env"`
	WorkDir string   `mapstructure:"work_dir"`
}

type ASTGrepConfig struct {
	Enabled   bool     `mapstructure:"enabled"`
	Languages []string `mapstructure:"languages"`
}

type LSPConfig struct {
	Enabled   bool                       `mapstructure:"enabled"`
	Servers   map[string]LSPServerConfig `mapstructure:"servers"`
	AutoStart bool                       `mapstructure:"auto_start"`
}

type LSPServerConfig struct {
	Command    []string `mapstructure:"command"`
	RootURI    string   `mapstructure:"root_uri"`
	LanguageID string   `mapstructure:"language_id"`
}

type HooksConfig struct {
	Session      []string `mapstructure:"session"`
	Tool         []string `mapstructure:"tool"`
	Transform    []string `mapstructure:"transform"`
	Continuation []string `mapstructure:"continuation"`
	Ralph        []string `mapstructure:"ralph"`
	Skill        []string `mapstructure:"skill"`
}

type PermissionsConfig struct {
	AllowedDirs  []string                   `mapstructure:"allowed_dirs"`
	DeniedDirs   []string                   `mapstructure:"denied_dirs"`
	AllowedHosts []string                   `mapstructure:"allowed_hosts"`
	DeniedHosts  []string                   `mapstructure:"denied_hosts"`
	ToolPatterns map[string][]string        `mapstructure:"tool_patterns"`
	AgentRules   map[string]AgentRuleConfig `mapstructure:"agent_rules"`
}

type AgentRuleConfig struct {
	Allow []string `mapstructure:"allow"`
	Deny  []string `mapstructure:"deny"`
}

type PlatformConfig struct {
	OS       string        `mapstructure:"os"`
	Arch     string        `mapstructure:"arch"`
	HomeDir  string        `mapstructure:"home_dir"`
	TempDir  string        `mapstructure:"temp_dir"`
	CacheDir string        `mapstructure:"cache_dir"`
	FreeBSD  FreeBSDConfig `mapstructure:"freebsd"`
	Darwin   DarwinConfig  `mapstructure:"darwin"`
	Linux    LinuxConfig   `mapstructure:"linux"`
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
	Dir         string           `mapstructure:"dir"`
	Compaction  CompactionConfig `mapstructure:"compaction"`
	HistorySize int              `mapstructure:"history_size"`
}

type CompactionConfig struct {
	Enabled         bool `mapstructure:"enabled"`
	ThresholdTokens int  `mapstructure:"threshold_tokens"`
	KeepLastN       int  `mapstructure:"keep_last_n"`
}

func DefaultConfig() *Config {
	homeDir, _ := os.UserHomeDir()
	tempDir := os.TempDir()
	cacheDir := filepath.Join(tempDir, "freecode-cache")

	return &Config{
		Shell:        "/bin/bash",
		LogLevel:     "info",
		LogFormat:    "text",
		Yolo:         true,
		Color:        true,
		SoundEnabled: true,
		Theme:        "default",
		Animation:    "full",
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
		OpenAI: OpenAIConfig{
			BaseURL: "https://api.openai.com/v1",
			APIKey:  "",
		},
		Anthropic: AnthropicConfig{
			BaseURL: "https://api.anthropic.com/v1",
			APIKey:  "",
		},
		Ollama: OllamaConfig{
			BaseURL: "http://localhost:11434",
			APIKey:  "",
		},
		Azure: AzureConfig{
			BaseURL: "",
			APIKey:  "",
		},
		Google: GoogleConfig{
			BaseURL: "https://generativelanguage.googleapis.com",
			APIKey:  "",
		},
		Vertex: VertexConfig{
			ProjectID:   "",
			Location:    "us-central1",
			BaseURL:     "",
			AccessToken: "",
		},
		AWS: AWSConfig{
			Region:    "",
			Profile:   "",
			Endpoint:  "",
			AccessKey: "",
			SecretKey: "",
		},
		GitLab: GitLabConfig{
			BaseURL: "https://gitlab.com",
			Token:   "",
		},
		GitHubCopilot: GitHubCopilotConfig{
			BaseURL: "https://api.github.com",
			Token:   "",
		},
		Vercel: VercelConfig{
			BaseURL: "https://api.vercel.com",
			Token:   "",
		},
		Groq: GroqConfig{
			BaseURL: "https://api.groq.com/openai/v1",
			APIKey:  "",
		},
		Perplexity: PerplexityConfig{
			BaseURL: "https://api.perplexity.ai",
			APIKey:  "",
		},
		Mistral: MistralConfig{
			BaseURL: "https://api.mistral.ai/v1",
			APIKey:  "",
		},
		Cohere: CohereConfig{
			BaseURL: "https://api.cohere.ai/v1",
			APIKey:  "",
		},
		TogetherAI: TogetherAIConfig{
			BaseURL: "https://api.together.xyz/v1",
			APIKey:  "",
		},
		DeepInfra: DeepInfraConfig{
			BaseURL: "https://api.deepinfra.com/v1",
			APIKey:  "",
		},
		Cerebras: CerebrasConfig{
			BaseURL: "https://api.cerebras.ai/v1",
			APIKey:  "",
		},
		XAI: XAIConfig{
			BaseURL: "https://api.x.ai/v1",
			APIKey:  "",
		},
		Alibaba: AlibabaConfig{
			BaseURL: "https://dashscope.aliyuncs.com",
			APIKey:  "",
		},
		HuggingFace: HuggingFaceConfig{
			BaseURL: "https://api.endpoints.huggingface.cloud/v1",
			APIKey:  "",
		},
		DeepSeek: DeepSeekConfig{
			BaseURL: "https://api.deepseek.com/v1",
			APIKey:  "",
		},
		Fireworks: FireworksConfig{
			BaseURL: "https://api.fireworks.ai/v1",
			APIKey:  "",
		},
		Moonshot: MoonshotConfig{
			BaseURL: "https://api.moonshot.cn/v1",
			APIKey:  "",
		},
		Nebius: NebiusConfig{
			BaseURL: "https://api.nebius.ai/v1",
			APIKey:  "",
		},
		OpenRouter: OpenRouterConfig{
			BaseURL: "https://openrouter.ai/api/v1",
			APIKey:  "",
		},
		Venice: VeniceConfig{
			BaseURL: "https://api.venice.ai/api/v1",
			APIKey:  "",
		},
		ZAI: ZAIConfig{
			BaseURL: "https://api.z-ai.ai/v1",
			APIKey:  "",
		},
		ZenMux: ZenMuxConfig{
			BaseURL: "https://api.zenmux.ai/v1",
			APIKey:  "",
		},
		Baseten: BasetenConfig{
			BaseURL: "https://app.baseten.co/v1",
			APIKey:  "",
		},
		Cortecs: CortecsConfig{
			BaseURL: "https://api.cortecs.ai/v1",
			APIKey:  "",
		},
		Firmware: FirmwareConfig{
			BaseURL: "https://api.firmware.ai/v1",
			APIKey:  "",
		},
		Ionet: IonetConfig{
			BaseURL: "https://api.ionet.ai/v1",
			APIKey:  "",
		},
		NVIDIA: NVIDIAConfig{
			BaseURL: "https://ai.api.nvidia.com/v1",
			APIKey:  "",
		},
		OllamaCloud: OllamaCloudConfig{
			BaseURL: "https://cloud.ollama.ai",
			APIKey:  "",
		},
		Cloudflare: CloudflareConfig{
			AccountID: "",
			GatewayID: "",
			APIKey:    "",
			BaseURL:   "",
		},
		Helicone: HeliconeConfig{
			BaseURL: "https://ai-gateway.helicone.ai",
			APIKey:  "",
		},
		LlamaCpp: LlamaCppConfig{
			BaseURL: "http://127.0.0.1:8080/v1",
			APIKey:  "",
		},
		LMStudio: LMStudioConfig{
			BaseURL: "http://127.0.0.1:1234/v1",
			APIKey:  "",
		},
		AtomicChat: AtomicChatConfig{
			BaseURL: "http://127.0.0.1:1337/v1",
			APIKey:  "",
		},
		Provider302AI: Provider302AIConfig{
			BaseURL: "https://api.302.ai/v1",
			APIKey:  "",
		},
		SAPAI: SAPAIConfig{
			ServiceKey: "",
			BaseURL:    "",
		},
		STACKIT: STACKITConfig{
			Token:   "",
			BaseURL: "",
		},
		OVHcloud: OVHcloudConfig{
			BaseURL: "https://endpoints.ai.cloud.ovh.net/v1",
			APIKey:  "",
		},
		Scaleway: ScalewayConfig{
			BaseURL: "https://api.scaleway.ai/v1",
			APIKey:  "",
		},
		Minimax: MinimaxConfig{
			BaseURL: "https://api.minimax.chat/v1",
			APIKey:  "",
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
