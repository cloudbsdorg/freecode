package shell

import (
	"testing"
)

func TestNewExecutor(t *testing.T) {
	e := NewExecutor()

	if e == nil {
		t.Fatal("NewExecutor() returned nil")
	}

	if e.env == nil {
		t.Error("Executor.env is nil")
	}
}

func TestExecutorWithEnv(t *testing.T) {
	e := NewExecutor()

	custom := []string{"TEST=value"}
	e = e.WithEnv(custom)

	if len(e.env) != 1 {
		t.Errorf("len(Executor.env) = %d, want 1", len(e.env))
	}
}

func TestExecutorWithDir(t *testing.T) {
	e := NewExecutor()

	e = e.WithDir("/tmp")

	if e.dir != "/tmp" {
		t.Errorf("Executor.dir = %q, want %q", e.dir, "/tmp")
	}
}

func TestExecutorExec(t *testing.T) {
	e := NewExecutor()

	result, err := e.Exec("echo", "hello")
	if err != nil {
		t.Fatalf("Exec() error = %v", err)
	}

	if result.ExitCode != 0 {
		t.Errorf("ExitCode = %d, want 0", result.ExitCode)
	}

	if result.Stdout == "" {
		t.Error("Stdout is empty")
	}
}

func TestExecutorExecShell(t *testing.T) {
	e := NewExecutor()

	result, err := e.ExecShell("echo 'shell test'")
	if err != nil {
		t.Fatalf("ExecShell() error = %v", err)
	}

	if result.ExitCode != 0 {
		t.Errorf("ExitCode = %d, want 0", result.ExitCode)
	}
}

func TestExecutorExecWithArgs(t *testing.T) {
	e := NewExecutor()

	result, err := e.Exec("printf", "test-%s", "value")
	if err != nil {
		t.Fatalf("Exec() error = %v", err)
	}

	if result.ExitCode != 0 {
		t.Errorf("ExitCode = %d, want 0", result.ExitCode)
	}
}
