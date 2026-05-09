package effect

import (
	"context"
	"fmt"
	"time"
)

type Effect interface {
	Run(ctx context.Context) error
}

type registry struct {
	effects map[string]func() Effect
}

var reg = &registry{effects: make(map[string]func() Effect)}

func Register(name string, factory func() Effect) {
	reg.effects[name] = factory
}

func Get(name string) (func() Effect, bool) {
	fn, ok := reg.effects[name]
	return fn, ok
}

func Run(ctx context.Context, name string) error {
	fn, ok := Get(name)
	if !ok {
		return nil
	}
	effect := fn()
	return effect.Run(ctx)
}

type Runner struct {
	defaultTimeout time.Duration
}

func NewRunner(timeout time.Duration) *Runner {
	if timeout <= 0 {
		timeout = 30 * time.Second
	}
	return &Runner{defaultTimeout: timeout}
}

type RunOptions struct {
	Timeout time.Duration
	Name    string
}

func (r *Runner) Run(ctx context.Context, name string, opts *RunOptions) error {
	fn, ok := Get(name)
	if !ok {
		return fmt.Errorf("effect %q not found", name)
	}

	effect := fn()

	timeout := r.defaultTimeout
	if opts != nil && opts.Timeout > 0 {
		timeout = opts.Timeout
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	errCh := make(chan error, 1)
	go func() {
		errCh <- effect.Run(ctx)
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		return fmt.Errorf("effect %q timed out after %v", name, timeout)
	}
}

func (r *Runner) RunSync(name string) error {
	return r.Run(context.Background(), name, nil)
}

func (r *Runner) RunPromise(name string) (err error) {
	done := make(chan struct{})
	go func() {
		err = r.Run(context.Background(), name, nil)
		close(done)
	}()
	<-done
	return
}
