package agent

import (
	"context"
	"fmt"
)

type SisyphusAgent struct {
	engine *Engine
}

func NewSisyphusAgent(engine *Engine) *SisyphusAgent {
	return &SisyphusAgent{engine: engine}
}

func (a *SisyphusAgent) Name() string {
	return "sisyphus"
}

func (a *SisyphusAgent) Run(ctx context.Context, req Request) (*Response, error) {
	cfg, ok := GetAgentConfig("sisyphus")
	if !ok {
		return nil, fmt.Errorf("sisyphus agent config not found")
	}
	return &Response{
		SessionID: req.SessionID,
		Message: Message{
			Role:    "assistant",
			Content: fmt.Sprintf("[Using Sisyphus agent with system prompt]\n\nUser request: %s", req.Message.Content),
		},
		AgentName:    cfg.Name,
		SystemPrompt: cfg.SystemPrompt,
	}, nil
}

type HephaestusAgent struct {
	engine *Engine
}

func NewHephaestusAgent(engine *Engine) *HephaestusAgent {
	return &HephaestusAgent{engine: engine}
}

func (a *HephaestusAgent) Name() string {
	return "hephaestus"
}

func (a *HephaestusAgent) Run(ctx context.Context, req Request) (*Response, error) {
	cfg, ok := GetAgentConfig("hephaestus")
	if !ok {
		return nil, fmt.Errorf("hephaestus agent config not found")
	}
	return &Response{
		SessionID: req.SessionID,
		Message: Message{
			Role:    "assistant",
			Content: fmt.Sprintf("[Using Hephaestus agent - code generation specialist]\n\nUser request: %s", req.Message.Content),
		},
		AgentName:    cfg.Name,
		SystemPrompt: cfg.SystemPrompt,
	}, nil
}

type OracleAgent struct {
	engine *Engine
}

func NewOracleAgent(engine *Engine) *OracleAgent {
	return &OracleAgent{engine: engine}
}

func (a *OracleAgent) Name() string {
	return "oracle"
}

func (a *OracleAgent) Run(ctx context.Context, req Request) (*Response, error) {
	cfg, ok := GetAgentConfig("oracle")
	if !ok {
		return nil, fmt.Errorf("oracle agent config not found")
	}
	return &Response{
		SessionID: req.SessionID,
		Message: Message{
			Role:    "assistant",
			Content: fmt.Sprintf("[Using Oracle agent - architecture consultation]\n\nUser request: %s", req.Message.Content),
		},
		AgentName:    cfg.Name,
		SystemPrompt: cfg.SystemPrompt,
	}, nil
}

type LibrarianAgent struct {
	engine *Engine
}

func NewLibrarianAgent(engine *Engine) *LibrarianAgent {
	return &LibrarianAgent{engine: engine}
}

func (a *LibrarianAgent) Name() string {
	return "librarian"
}

func (a *LibrarianAgent) Run(ctx context.Context, req Request) (*Response, error) {
	cfg, ok := GetAgentConfig("librarian")
	if !ok {
		return nil, fmt.Errorf("librarian agent config not found")
	}
	return &Response{
		SessionID: req.SessionID,
		Message: Message{
			Role:    "assistant",
			Content: fmt.Sprintf("[Using Librarian agent - research and documentation]\n\nUser request: %s", req.Message.Content),
		},
		AgentName:    cfg.Name,
		SystemPrompt: cfg.SystemPrompt,
	}, nil
}

type ExploreAgent struct {
	engine *Engine
}

func NewExploreAgent(engine *Engine) *ExploreAgent {
	return &ExploreAgent{engine: engine}
}

func (a *ExploreAgent) Name() string {
	return "explore"
}

func (a *ExploreAgent) Run(ctx context.Context, req Request) (*Response, error) {
	cfg, ok := GetAgentConfig("explore")
	if !ok {
		return nil, fmt.Errorf("explore agent config not found")
	}
	return &Response{
		SessionID: req.SessionID,
		Message: Message{
			Role:    "assistant",
			Content: fmt.Sprintf("[Using Explore agent - codebase exploration]\n\nUser request: %s", req.Message.Content),
		},
		AgentName:    cfg.Name,
		SystemPrompt: cfg.SystemPrompt,
	}, nil
}

type PrometheusAgent struct {
	engine *Engine
}

func NewPrometheusAgent(engine *Engine) *PrometheusAgent {
	return &PrometheusAgent{engine: engine}
}

func (a *PrometheusAgent) Name() string {
	return "prometheus"
}

func (a *PrometheusAgent) Run(ctx context.Context, req Request) (*Response, error) {
	cfg, ok := GetAgentConfig("prometheus")
	if !ok {
		return nil, fmt.Errorf("prometheus agent config not found")
	}
	return &Response{
		SessionID: req.SessionID,
		Message: Message{
			Role:    "assistant",
			Content: fmt.Sprintf("[Using Prometheus agent - task planning]\n\nUser request: %s", req.Message.Content),
		},
		AgentName:    cfg.Name,
		SystemPrompt: cfg.SystemPrompt,
	}, nil
}

type MetisAgent struct {
	engine *Engine
}

func NewMetisAgent(engine *Engine) *MetisAgent {
	return &MetisAgent{engine: engine}
}

func (a *MetisAgent) Name() string {
	return "metis"
}

func (a *MetisAgent) Run(ctx context.Context, req Request) (*Response, error) {
	cfg, ok := GetAgentConfig("metis")
	if !ok {
		return nil, fmt.Errorf("metis agent config not found")
	}
	return &Response{
		SessionID: req.SessionID,
		Message: Message{
			Role:    "assistant",
			Content: fmt.Sprintf("[Using Metis agent - pre-planning consultation]\n\nUser request: %s", req.Message.Content),
		},
		AgentName:    cfg.Name,
		SystemPrompt: cfg.SystemPrompt,
	}, nil
}

type MomusAgent struct {
	engine *Engine
}

func NewMomusAgent(engine *Engine) *MomusAgent {
	return &MomusAgent{engine: engine}
}

func (a *MomusAgent) Name() string {
	return "momus"
}

func (a *MomusAgent) Run(ctx context.Context, req Request) (*Response, error) {
	cfg, ok := GetAgentConfig("momus")
	if !ok {
		return nil, fmt.Errorf("momus agent config not found")
	}
	return &Response{
		SessionID: req.SessionID,
		Message: Message{
			Role:    "assistant",
			Content: fmt.Sprintf("[Using Momus agent - code review]\n\nUser request: %s", req.Message.Content),
		},
		AgentName:    cfg.Name,
		SystemPrompt: cfg.SystemPrompt,
	}, nil
}

type AtlasAgent struct {
	engine *Engine
}

func NewAtlasAgent(engine *Engine) *AtlasAgent {
	return &AtlasAgent{engine: engine}
}

func (a *AtlasAgent) Name() string {
	return "atlas"
}

func (a *AtlasAgent) Run(ctx context.Context, req Request) (*Response, error) {
	cfg, ok := GetAgentConfig("atlas")
	if !ok {
		return nil, fmt.Errorf("atlas agent config not found")
	}
	return &Response{
		SessionID: req.SessionID,
		Message: Message{
			Role:    "assistant",
			Content: fmt.Sprintf("[Using Atlas agent - task tracking]\n\nUser request: %s", req.Message.Content),
		},
		AgentName:    cfg.Name,
		SystemPrompt: cfg.SystemPrompt,
	}, nil
}

type MultimodalLookerAgent struct {
	engine *Engine
}

func NewMultimodalLookerAgent(engine *Engine) *MultimodalLookerAgent {
	return &MultimodalLookerAgent{engine: engine}
}

func (a *MultimodalLookerAgent) Name() string {
	return "multimodal-looker"
}

func (a *MultimodalLookerAgent) Run(ctx context.Context, req Request) (*Response, error) {
	cfg, ok := GetAgentConfig("multimodal-looker")
	if !ok {
		return nil, fmt.Errorf("multimodal-looker agent config not found")
	}
	return &Response{
		SessionID: req.SessionID,
		Message: Message{
			Role:    "assistant",
			Content: fmt.Sprintf("[Using Multimodal-Looker agent - image/document analysis]\n\nUser request: %s", req.Message.Content),
		},
		AgentName:    cfg.Name,
		SystemPrompt: cfg.SystemPrompt,
	}, nil
}

type SisyphusJuniorAgent struct {
	engine *Engine
}

func NewSisyphusJuniorAgent(engine *Engine) *SisyphusJuniorAgent {
	return &SisyphusJuniorAgent{engine: engine}
}

func (a *SisyphusJuniorAgent) Name() string {
	return "sisyphus-junior"
}

func (a *SisyphusJuniorAgent) Run(ctx context.Context, req Request) (*Response, error) {
	cfg, ok := GetAgentConfig("sisyphus-junior")
	if !ok {
		return nil, fmt.Errorf("sisyphus-junior agent config not found")
	}
	return &Response{
		SessionID: req.SessionID,
		Message: Message{
			Role:    "assistant",
			Content: fmt.Sprintf("[Using Sisyphus-Junior agent - simple tasks]\n\nUser request: %s", req.Message.Content),
		},
		AgentName:    cfg.Name,
		SystemPrompt: cfg.SystemPrompt,
	}, nil
}

func RegisterBuiltinAgents(e *Engine) {
	e.RegisterAgent(NewSisyphusAgent(e))
	e.RegisterAgent(NewHephaestusAgent(e))
	e.RegisterAgent(NewOracleAgent(e))
	e.RegisterAgent(NewLibrarianAgent(e))
	e.RegisterAgent(NewExploreAgent(e))
	e.RegisterAgent(NewPrometheusAgent(e))
	e.RegisterAgent(NewMetisAgent(e))
	e.RegisterAgent(NewMomusAgent(e))
	e.RegisterAgent(NewAtlasAgent(e))
	e.RegisterAgent(NewMultimodalLookerAgent(e))
	e.RegisterAgent(NewSisyphusJuniorAgent(e))
}
