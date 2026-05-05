package pty

import (
	"context"
	"io"
	"os"
	"os/exec"
	"sync"

	"github.com/creack/pty"
)

type PTY struct {
	ptty *os.File
	cmd  *exec.Cmd
	mu   sync.Mutex

	Stdout <-chan []byte
	Stderr <-chan []byte
	Input  chan<- []byte

	stdoutCh chan []byte
	stderrCh chan []byte
	inputCh  chan []byte
}

type Options struct {
	Cols uint16
	Rows uint16
	Dir  string
	Env  []string
}

func Start(ctx context.Context, command string, args []string, opts Options) (*PTY, error) {
	cmd := exec.CommandContext(ctx, command, args...)
	if opts.Dir != "" {
		cmd.Dir = opts.Dir
	}
	if opts.Env != nil {
		cmd.Env = opts.Env
	} else {
		cmd.Env = os.Environ()
	}

	ptty, err := pty.StartWithSize(cmd, &pty.Winsize{
		Cols: opts.Cols,
		Rows: opts.Rows,
	})
	if err != nil {
		return nil, err
	}

	p := &PTY{
		ptty:    ptty,
		cmd:     cmd,
		stdoutCh: make(chan []byte, 1024),
		stderrCh: make(chan []byte, 1024),
		inputCh:  make(chan []byte, 1024),
	}
	p.Stdout = p.stdoutCh
	p.Stderr = p.stderrCh
	p.Input = p.inputCh

	go p.readLoop()
	go p.writeLoop()

	return p, nil
}

func (p *PTY) readLoop() {
	buf := make([]byte, 4096)
	for {
		n, err := p.ptty.Read(buf)
		if n > 0 {
			p.stdoutCh <- append([]byte{}, buf[:n]...)
		}
		if err != nil {
			if err != io.EOF {
				p.stderrCh <- []byte(err.Error())
			}
			close(p.stdoutCh)
			close(p.stderrCh)
			return
		}
	}
}

func (p *PTY) writeLoop() {
	for data := range p.inputCh {
		p.mu.Lock()
		_, err := p.ptty.Write(data)
		p.mu.Unlock()
		if err != nil {
			return
		}
	}
}

func (p *PTY) Close() error {
	close(p.inputCh)
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.ptty != nil {
		p.ptty.Close()
	}
	if p.cmd != nil && p.cmd.Process != nil {
		p.cmd.Process.Kill()
		p.cmd.Wait()
	}
	return nil
}

func (p *PTY) Resize(cols, rows uint16) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.ptty == nil {
		return nil
	}
	return pty.Setsize(p.ptty, &pty.Winsize{Cols: cols, Rows: rows})
}

func (p *PTY) OutputReader() io.Reader {
	return p.ptty
}

func (p *PTY) InputWriter() io.Writer {
	return p.ptty
}

func (p *PTY) TTY() *os.File {
	return p.ptty
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
	ws, err := pty.GetsizeFull(os.Stdout)
	if err != nil {
		return 80, 24, nil
	}
	return ws.Cols, ws.Rows, nil
}
