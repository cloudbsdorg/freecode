package shell

import (
	"io"
	"os"
	"sync"
)

type PTY struct {
	ptyMaster *os.File
	ptySlave  *os.File
	writer    io.Writer
	mu        sync.Mutex
}

func StartPTY(command string, args ...string) (*PTY, error) {
	return &PTY{}, nil
}

func (p *PTY) Write(data []byte) (int, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	return 0, nil
}

func (p *PTY) Read(data []byte) (int, error) {
	return 0, nil
}

func (p *PTY) Close() error {
	return nil
}

func (p *PTY) Resize(rows, cols int) error {
	return nil
}
