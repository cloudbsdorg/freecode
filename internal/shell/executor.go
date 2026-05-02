package shell

import (
	"os"
	"os/exec"
	"strings"
)

type Executor struct {
	env    []string
	dir    string
	stdin  strings.Builder
	stdout strings.Builder
	stderr strings.Builder
}

type Result struct {
	ExitCode int
	Stdout   string
	Stderr   string
}

func NewExecutor() *Executor {
	return &Executor{
		env: os.Environ(),
	}
}

func (e *Executor) WithEnv(env []string) *Executor {
	e.env = env
	return e
}

func (e *Executor) WithDir(dir string) *Executor {
	e.dir = dir
	return e
}

func (e *Executor) Exec(command string, args ...string) (Result, error) {
	cmd := exec.Command(command, args...)
	cmd.Env = e.env
	cmd.Dir = e.dir

	out, err := cmd.CombinedOutput()

	result := Result{
		ExitCode: 0,
		Stdout:   string(out),
		Stderr:   "",
	}

	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			result.ExitCode = exitErr.ExitCode()
		} else {
			result.ExitCode = 1
		}
	}

	return result, nil
}

func (e *Executor) ExecShell(command string) (Result, error) {
	return e.Exec("sh", "-c", command)
}
