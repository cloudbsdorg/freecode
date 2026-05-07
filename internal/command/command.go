package command

import (
	"context"
)

type Command struct {
	Name string
	Args []string
	Dir  string
	Env  []string
}

type Runner interface {
	Run(ctx context.Context, cmd Command) error
	RunWithOutput(ctx context.Context, cmd Command) (string, error)
}

type registry struct {
	commands map[string]func(args []string) error
}

var reg = &registry{commands: make(map[string]func(args []string) error)}

func Register(name string, fn func(args []string) error) {
	reg.commands[name] = fn
}

func Get(name string) (func(args []string) error, bool) {
	fn, ok := reg.commands[name]
	return fn, ok
}

func List() []string {
	var names []string
	for name := range reg.commands {
		names = append(names, name)
	}
	return names
}

func Run(ctx context.Context, cmd Command, runner Runner) error {
	return runner.Run(ctx, cmd)
}

func RunOutput(ctx context.Context, cmd Command, runner Runner) (string, error) {
	return runner.RunWithOutput(ctx, cmd)
}
