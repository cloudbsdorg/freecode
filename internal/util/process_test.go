package util

import (
	"context"
	"errors"
	"os"
	"runtime"
	"testing"
	"time"
)

func TestRunSimpleCommand(t *testing.T) {
	result, err := Run([]string{"echo", "hello"}, Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Code != 0 {
		t.Errorf("expected exit code 0, got %d", result.Code)
	}
	if string(result.Stdout) == "" && runtime.GOOS != "windows" {
		t.Error("expected stdout output")
	}
}

func TestRunWithArgs(t *testing.T) {
	result, err := Run([]string{"printf", "hello %s", "world"}, Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Code != 0 {
		t.Errorf("expected exit code 0, got %d", result.Code)
	}
	if string(result.Stdout) != "hello world" {
		t.Errorf("expected 'hello world', got '%s'", string(result.Stdout))
	}
}

func TestRunFailedCommand(t *testing.T) {
	var cmd []string
	if runtime.GOOS == "windows" {
		cmd = []string{"cmd", "/c", "exit 1"}
	} else {
		cmd = []string{"sh", "-c", "exit 1"}
	}

	_, err := Run(cmd, Options{})
	if err == nil {
		t.Fatal("expected error for non-zero exit")
	}
	rfe, ok := err.(*RunFailedError)
	if !ok {
		t.Fatalf("expected RunFailedError, got %T", err)
	}
	if rfe.Code != 1 {
		t.Errorf("expected exit code 1, got %d", rfe.Code)
	}
	if len(rfe.Stdout) != 0 {
		t.Error("expected no stdout on failure")
	}
}

func TestRunMissingCommand(t *testing.T) {
	_, err := Run([]string{"nonexistent_command_12345"}, Options{})
	if err == nil {
		t.Fatal("expected error for missing command")
	}
}

func TestRunWithEnv(t *testing.T) {
	result, err := Run([]string{"printenv", "TEST_VAR"}, Options{
		Env: map[string]string{"TEST_VAR": "hello_env"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	output := string(result.Stdout)
	if output == "" {
		t.Error("expected stdout output")
	}
	if output != "hello_env" && output != "hello_env\n" && output != "hello_env\r\n" {
		t.Errorf("expected 'hello_env' (possibly with newline), got '%s'", output)
	}
}

func TestRunWithCwd(t *testing.T) {
	result, err := Run([]string{"pwd"}, Options{
		Cwd: "/tmp",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	output := string(result.Stdout)
	if output == "" && runtime.GOOS != "windows" {
		t.Error("expected pwd output")
	}
}

func TestRunWithTimeout(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("skipping timeout test on windows")
	}
	result, err := Run([]string{"sleep", "10"}, Options{
		Timeout: 100 * time.Millisecond,
	})
	if err == nil {
		t.Fatal("expected error for timeout")
	}
	if result.Code == 0 {
		t.Error("expected non-zero exit code on timeout")
	}
}

func TestTextSimpleCommand(t *testing.T) {
	result, err := Text([]string{"printf", "hello"}, Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Text != "hello" {
		t.Errorf("expected 'hello', got '%s'", result.Text)
	}
	if result.Code != 0 {
		t.Errorf("expected exit code 0, got %d", result.Code)
	}
}

func TestTextFailedCommand(t *testing.T) {
	var cmd []string
	if runtime.GOOS == "windows" {
		cmd = []string{"cmd", "/c", "echo stderr >&2 && exit 1"}
	} else {
		cmd = []string{"sh", "-c", "echo stderr >&2 && exit 1"}
	}

	_, err := Text(cmd, Options{})
	if err == nil {
		t.Fatal("expected error for non-zero exit")
	}
	rfe, ok := err.(*RunFailedError)
	if !ok {
		t.Fatalf("expected RunFailedError, got %T", err)
	}
	if rfe.Code != 1 {
		t.Errorf("expected exit code 1, got %d", rfe.Code)
	}
}

func TestSpawnAndWait(t *testing.T) {
	child, err := Spawn([]string{"echo", "hello"}, Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	pid := child.Pid()
	if pid <= 0 {
		t.Error("expected valid pid")
	}

	code := <-child.Exit
	if code != 0 {
		t.Errorf("expected exit code 0, got %d", code)
	}
}

func TestSpawnWithStdioPipe(t *testing.T) {
	child, err := Spawn([]string{"printf", "hello world"}, Options{
		Stdout: StdioPipe,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	code := <-child.Exit
	if code != 0 {
		t.Errorf("expected exit code 0, got %d", code)
	}
}

func TestSpawnFailedCommand(t *testing.T) {
	_, err := Spawn([]string{"nonexistent_command_12345"}, Options{})
	if err == nil {
		t.Fatal("expected error for missing command")
	}
}

func TestChildKill(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("skipping kill test on windows")
	}

	child, err := Spawn([]string{"sleep", "10"}, Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err = child.Kill()
	if err != nil {
		t.Fatalf("unexpected error killing process: %v", err)
	}

	code := <-child.Exit
	if code == 0 {
		t.Error("expected non-zero exit code after kill")
	}
}

func TestStop(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("skipping stop test on windows")
	}

	child, err := Spawn([]string{"sleep", "10"}, Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	pid := child.Pid()
	if pid <= 0 {
		t.Fatal("expected valid pid")
	}

	err = Stop(pid)
	if err != nil {
		t.Fatalf("unexpected error stopping process: %v", err)
	}

	code := <-child.Exit
	if code == 0 {
		t.Error("expected non-zero exit code after stop")
	}
}

func TestStopInvalidPid(t *testing.T) {
	err := Stop(0)
	if err == nil {
		t.Error("expected error for invalid pid")
	}
}

func TestStopNegativePid(t *testing.T) {
	err := Stop(-1)
	if err == nil {
		t.Error("expected error for negative pid")
	}
}

func TestRunContext(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	result, err := RunContext(ctx, []string{"sleep", "10"}, Options{})
	if err == nil {
		t.Fatal("expected error for context timeout")
	}
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("expected DeadlineExceeded, got %v", err)
	}
	if result.Code == 0 {
		t.Error("expected non-zero exit code on context cancellation")
	}
}

func TestRunContextCompleted(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := RunContext(ctx, []string{"echo", "hello"}, Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Code != 0 {
		t.Errorf("expected exit code 0, got %d", result.Code)
	}
}

func TestRunContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		time.Sleep(50 * time.Millisecond)
		cancel()
	}()

	result, err := RunContext(ctx, []string{"sleep", "10"}, Options{})
	if err == nil {
		t.Fatal("expected error for context cancellation")
	}
	if !errors.Is(err, context.Canceled) {
		t.Errorf("expected Canceled, got %v", err)
	}
	if result.Code == 0 {
		t.Error("expected non-zero exit code on context cancellation")
	}
}

func TestLines(t *testing.T) {
	tests := []struct {
		input    string
		expected []string
	}{
		{"hello\nworld", []string{"hello", "world"}},
		{"hello\nworld\n", []string{"hello", "world"}},
		{"hello\r\nworld", []string{"hello", "world"}},
		{"hello\r\n\r\nworld", []string{"hello", "world"}},
		{"", []string{}},
		{"hello", []string{"hello"}},
		{"hello\n", []string{"hello"}},
		{"\n", []string{}},
		{"\n\n", []string{}},
	}

	for _, tc := range tests {
		result := Lines(tc.input)
		if len(result) != len(tc.expected) {
			t.Errorf("Lines(%q): expected %v, got %v", tc.input, tc.expected, result)
			continue
		}
		for i, line := range result {
			if line != tc.expected[i] {
				t.Errorf("Lines(%q): expected %v, got %v", tc.input, tc.expected, result)
				break
			}
		}
	}
}

func TestSplitLines(t *testing.T) {
	result := splitLines("hello\nworld\n")
	if len(result) != 2 {
		t.Errorf("expected 2 lines, got %d", len(result))
	}
	if result[0] != "hello" || result[1] != "world" {
		t.Errorf("unexpected lines: %v", result)
	}
}

func TestLinesWithEmptyStrings(t *testing.T) {
	result := Lines("hello\n\nworld")
	if len(result) != 2 {
		t.Errorf("expected 2 non-empty lines, got %d", len(result))
	}
}

func TestLinesFiltersEmpty(t *testing.T) {
	result := Lines("hello\n\nworld")
	if len(result) != 2 {
		t.Errorf("expected 2 non-empty lines, got %d", len(result))
	}
	if result[0] != "hello" || result[1] != "world" {
		t.Errorf("unexpected result: %v", result)
	}
}

func TestOutputCollector(t *testing.T) {
	collector := &OutputCollector{}

	result, err := RunWithCollector([]string{"printf", "hello world"}, Options{}, collector)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Code != 0 {
		t.Errorf("expected exit code 0, got %d", result.Code)
	}

	stdout := collector.Stdout.String()
	if stdout != "hello world" {
		t.Errorf("expected 'hello world', got '%s'", stdout)
	}
}

func TestOutputCollectorWithNil(t *testing.T) {
	result, err := RunWithCollector([]string{"echo", "hello"}, Options{}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Code != 0 {
		t.Errorf("expected exit code 0, got %d", result.Code)
	}
}

func TestShellEscapingNeeded(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"hello", false},
		{"hello world", true},
		{"hello$world", true},
		{"hello`world", true},
		{"hello\\world", true},
		{"hello'world", true},
		{"hello\"world", true},
		{"simple", false},
		{"no_special_chars", false},
	}

	for _, tc := range tests {
		result := shellEscapingNeeded(tc.input)
		if result != tc.expected {
			t.Errorf("shellEscapingNeeded(%q): expected %v, got %v", tc.input, tc.expected, result)
		}
	}
}

func TestJoinCmd(t *testing.T) {
	result := joinCmd([]string{"echo", "hello"})
	if result != "echo hello" {
		t.Errorf("expected 'echo hello', got '%s'", result)
	}

	result = joinCmd([]string{"echo", "hello world"})
	if result != "echo 'hello world'" {
		t.Errorf("expected \"echo 'hello world'\", got '%s'", result)
	}
}

func TestRunFailedError(t *testing.T) {
	err := &RunFailedError{
		Cmd:    []string{"test", "cmd"},
		Code:   1,
		Stdout: []byte("stdout"),
		Stderr: []byte("stderr"),
	}

	msg := err.Error()
	if msg == "" {
		t.Error("expected non-empty error message")
	}
}

func TestRunFailedErrorWithEmptyStderr(t *testing.T) {
	err := &RunFailedError{
		Cmd:    []string{"test", "cmd"},
		Code:   1,
		Stdout: []byte("stdout"),
		Stderr: []byte{},
	}

	msg := err.Error()
	if msg == "" {
		t.Error("expected non-empty error message")
	}
}

func TestShellFromString(t *testing.T) {
	if ShellFromString("") != ShellDefault {
		t.Error("expected ShellDefault for empty string")
	}
	if ShellFromString("bash") != "bash" {
		t.Error("expected 'bash'")
	}
}

func TestRunWithShell(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("skipping shell test on windows")
	}
	if os.Getenv("CI") != "" {
		t.Skip("skipping shell test in CI container environment")
	}

	result, err := Run([]string{"echo $SHELL"}, Options{
		Shell: "/bin/sh",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Code != 0 {
		t.Errorf("expected exit code 0, got %d", result.Code)
	}
}

func TestRunWithStdioInherit(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("skipping inherit test on windows")
	}

	result, err := Run([]string{"echo", "hello"}, Options{
		Stdout: StdioInherit,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Code != 0 {
		t.Errorf("expected exit code 0, got %d", result.Code)
	}
}

func TestRunWithStdioIgnore(t *testing.T) {
	result, err := Run([]string{"echo", "hello"}, Options{
		Stdout: StdioIgnore,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Code != 0 {
		t.Errorf("expected exit code 0, got %d", result.Code)
	}
	if len(result.Stdout) != 0 {
		t.Error("expected no stdout when StdioIgnore is set")
	}
}

func TestChildWait(t *testing.T) {
	child, err := Spawn([]string{"echo", "hello"}, Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	code, err := child.Wait()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if code != 0 {
		t.Errorf("expected exit code 0, got %d", code)
	}
}

func TestChildPid(t *testing.T) {
	child, err := Spawn([]string{"echo", "hello"}, Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	pid := child.Pid()
	if pid <= 0 {
		t.Error("expected valid pid")
	}

	<-child.Exit
}

func TestChildPidWithNilProcess(t *testing.T) {
	child := &Child{}
	if child.Pid() != 0 {
		t.Error("expected 0 for nil process")
	}
}

func TestChildKillWithNilProcess(t *testing.T) {
	child := &Child{}
	err := child.Kill()
	if err == nil {
		t.Error("expected error for nil process")
	}
}
