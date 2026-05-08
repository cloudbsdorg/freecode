package provider

import (
	"context"
)

type Provider interface {
	Name() string
	Generate(ctx context.Context, req *Request) (*Response, error)
	ListModels(ctx context.Context) ([]Model, error)
}

type Request struct {
	Model       string
	Messages    []Message
	Temperature float64
	MaxTokens   int
	Stream      bool
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Part struct {
	Type    string `json:"type"`
	Content string `json:"content,omitempty"`
	Tool    string `json:"tool,omitempty"`
}

type Response struct {
	Content    string
	StopReason string
	Usage      Usage
	Parts      []Part
}

type Usage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
}
