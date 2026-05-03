package pty

import (
	"context"
	"os"
)

type PTY struct {
	Stdout <-chan []byte
	Stderr <-chan []byte
	Input  chan<- []byte
}

type Options struct {
	Cols uint16
	Rows uint16
}

func Start(ctx context.Context, command string, args []string, opts Options) (*PTY, error) {
	return &PTY{}, nil
}

func (p *PTY) Close() error {
	return nil
}

func (p *PTY) Resize(cols, rows uint16) error {
	return nil
}

func NewTerminal(cols, rows uint16) *Terminal {
	return &Terminal{cols: cols, rows: rows}
}

type Terminal struct {
	cols uint16
	rows uint16
}

func (t *Terminal) Size() (cols, rows uint16) {
	return t.cols, t.rows
}

func (t *Terminal) SetSize(cols, rows uint16) {
	t.cols = cols
	t.rows = rows
}

func GetWindowSize() (cols, rows uint16, err error) {
	return 80, 24, nil
}

var (
	_ = os.Stdin
)
