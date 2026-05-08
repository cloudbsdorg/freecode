package agent

import (
	"context"
	"fmt"

	"github.com/freecode/freecode/internal/provider"
)

func runAgent(ctx context.Context, agentName string, req Request) (*Response, error) {
	cfg, ok := GetAgentConfig(agentName)
	if !ok {
		return nil, fmt.Errorf("agent config not found: %s", agentName)
	}

	model := req.Model
	if model == "" {
		model = cfg.DefaultModel
	}

	p := provider.NewProvider(model)

	messages := []provider.Message{
		{Role: "system", Content: cfg.SystemPrompt},
		{Role: "user", Content: req.Message.Content},
	}

	providerReq := &provider.Request{
		Model:       model,
		Messages:    messages,
		Temperature: 0.7,
		MaxTokens:   4096,
		Stream:      req.Stream,
	}

	resp, err := p.Generate(ctx, providerReq)
	if err != nil || resp == nil || resp.Content == "" {
		return stubResponse(agentName, req, cfg)
	}

	var parts []MessagePart
	for _, p := range resp.Parts {
		parts = append(parts, MessagePart{
			Type:    p.Type,
			Content: p.Content,
			Tool:    p.Tool,
		})
	}

	return &Response{
		SessionID: req.SessionID,
		Message:   Message{Role: "assistant", Content: resp.Content, Parts: parts},
		AgentName: cfg.Name,
	}, nil
}

func stubResponse(agentName string, req Request, cfg AgentConfig) (*Response, error) {
	agentDescriptions := map[string]string{
		"sisyphus":          "main orchestrator - coordinates all other agents",
		"hephaestus":         "code generation and refactoring specialist",
		"oracle":            "architecture and design consultation",
		"librarian":         "research and documentation lookup",
		"explore":           "codebase exploration and search",
		"prometheus":        "task planning and decomposition",
		"metis":             "pre-planning consultation and risk assessment",
		"momus":             "code review and quality assessment",
		"atlas":             "task tracking and progress monitoring",
		"multimodal-looker": "image and document analysis",
		"sisyphus-junior":   "simple task assistant",
	}

	desc, ok := agentDescriptions[agentName]
	if !ok {
		desc = "general purpose agent"
	}

	stubContent := fmt.Sprintf("[%s] %s\n\nUser request: %s\n\nNote: Configure an AI provider (ANTHROPIC_API_KEY, OPENAI_API_KEY, etc.) for full responses.", cfg.Name, desc, req.Message.Content)

	return &Response{
		SessionID:    req.SessionID,
		Message:      Message{Role: "assistant", Content: stubContent},
		AgentName:    cfg.Name,
		SystemPrompt: cfg.SystemPrompt,
	}, nil
}

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
	return runAgent(ctx, "sisyphus", req)
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
	return runAgent(ctx, "hephaestus", req)
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
	return runAgent(ctx, "oracle", req)
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
	return runAgent(ctx, "librarian", req)
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
	return runAgent(ctx, "explore", req)
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
	return runAgent(ctx, "prometheus", req)
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
	return runAgent(ctx, "metis", req)
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
	return runAgent(ctx, "momus", req)
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
	return runAgent(ctx, "atlas", req)
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
	return runAgent(ctx, "multimodal-looker", req)
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
	return runAgent(ctx, "sisyphus-junior", req)
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