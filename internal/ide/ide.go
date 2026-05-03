package ide

import (
	"context"
)

type IDE interface {
	OpenFile(ctx context.Context, path string) error
	OpenTerminal(ctx context.Context, cwd string) error
	ShowNotification(ctx context.Context, message string) error
	GetOpenFiles(ctx context.Context) ([]string, error)
}

type stubIDE struct{}

func NewStubIDE() *stubIDE {
	return &stubIDE{}
}

func (i *stubIDE) OpenFile(ctx context.Context, path string) error {
	return nil
}

func (i *stubIDE) OpenTerminal(ctx context.Context, cwd string) error {
	return nil
}

func (i *stubIDE) ShowNotification(ctx context.Context, message string) error {
	return nil
}

func (i *stubIDE) GetOpenFiles(ctx context.Context) ([]string, error) {
	return []string{}, nil
}
