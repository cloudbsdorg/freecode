package share

import (
	"context"
)

type Share struct {
	ID      string
	Content string
	URL     string
}

type Publisher interface {
	Publish(ctx context.Context, content string) (*Share, error)
	Get(ctx context.Context, id string) (*Share, error)
}

type localPublisher struct{}

func NewLocalPublisher() *localPublisher {
	return &localPublisher{}
}

func (p *localPublisher) Publish(ctx context.Context, content string) (*Share, error) {
	return &Share{ID: "local", Content: content}, nil
}

func (p *localPublisher) Get(ctx context.Context, id string) (*Share, error) {
	return &Share{ID: id}, nil
}
