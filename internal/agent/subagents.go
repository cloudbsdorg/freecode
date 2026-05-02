package agent

import (
	"context"
	"fmt"
)

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
	return &Response{
		SessionID: req.SessionID,
		Message: Message{
			Role:    "assistant",
			Content: "Prometheus agent - planning placeholder",
		},
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
	return &Response{
		SessionID: req.SessionID,
		Message: Message{
			Role:    "assistant",
			Content: "Metis agent - plan consultation placeholder",
		},
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
	return &Response{
		SessionID: req.SessionID,
		Message: Message{
			Role:    "assistant",
			Content: "Momus agent - code review placeholder",
		},
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
	return &Response{
		SessionID: req.SessionID,
		Message: Message{
			Role:    "assistant",
			Content: "Atlas agent - session tracking placeholder",
		},
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
	return &Response{
		SessionID: req.SessionID,
		Message: Message{
			Role:    "assistant",
			Content: "Multimodal-Looker agent - multimodal analysis placeholder",
		},
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
	return &Response{
		SessionID: req.SessionID,
		Message: Message{
			Role:    "assistant",
			Content: fmt.Sprintf("Sisyphus-Junior agent - simpler tasks placeholder"),
		},
	}, nil
}
