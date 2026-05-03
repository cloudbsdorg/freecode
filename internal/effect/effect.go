package effect

import (
	"context"
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
