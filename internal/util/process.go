// Package util provides process spawning and management utilities.
package util

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"syscall"
	"time"
)

type Stdio string

const (
	StdioInherit Stdio = "inherit"
	StdioPipe    Stdio = "pipe"
	StdioIgnore  Stdio = "ignore"
)

type Shell string

const ShellDefault Shell = ""

func ShellFromString(s string) Shell {
	return Shell(s)
}

type Options struct {
	Cwd     string
	Env     map[string]string
	Stdin   Stdio
	Stdout  Stdio
	Stderr  Stdio
	Shell   Shell
	Signal  *os.Signal
	Kill    syscall.Signal
	Timeout time.Duration
}

type Child struct {
	Cmd  *exec.Cmd
	Exit <-chan int
	done chan int
}

func (c *Child) Pid() int {
	if c.Cmd == nil || c.Cmd.Process == nil {
		return 0
	}
	return c.Cmd.Process.Pid
}

func (c *Child) Kill() error {
	if c.Cmd == nil || c.Cmd.Process == nil {
		return errors.New("process not running")
	}
	return c.Cmd.Process.Kill()
}

func (c *Child) Wait() (int, error) {
	err := c.Cmd.Wait()
	if c.Cmd.ProcessState == nil {
		return -1, err
	}
	return c.Cmd.ProcessState.ExitCode(), nil
}

type RunFailedError struct {
	Cmd    []string
	Code   int
	Stdout []byte
	Stderr []byte
}

func (e *RunFailedError) Error() string {
	if len(e.Stderr) > 0 {
		return fmt.Sprintf("Command failed with code %d: %v\n%s", e.Code, e.Cmd, string(e.Stderr))
	}
	return fmt.Sprintf("Command failed with code %d: %v", e.Code, e.Cmd)
}

type Result struct {
	Code   int
	Stdout []byte
	Stderr []byte
}

type TextResult struct {
	Result
	Text string
}

func buildCmd(cmd []string, opts Options) (*exec.Cmd, error) {
	if len(cmd) == 0 {
		return nil, errors.New("command is required")
	}

	command := exec.Command(cmd[0], cmd[1:]...)

	if opts.Cwd != "" {
		command.Dir = opts.Cwd
	}

	if opts.Env != nil {
		env := os.Environ()
		for k, v := range opts.Env {
			env = append(env, k+"="+v)
		}
		command.Env = env
	}

	stdin := opts.Stdin
	if stdin == "" {
		stdin = StdioIgnore
	}
	stdout := opts.Stdout
	if stdout == "" {
		stdout = StdioIgnore
	}
	stderr := opts.Stderr
	if stderr == "" {
		stderr = StdioIgnore
	}

	switch stdin {
	case StdioInherit:
		command.Stdin = os.Stdin
	}

	switch stdout {
	case StdioInherit:
		command.Stdout = os.Stdout
	}

	switch stderr {
	case StdioInherit:
		command.Stderr = os.Stderr
	}

	if opts.Shell != "" {
		shell := string(opts.Shell)
		if shell == "" {
			if runtime.GOOS == "windows" {
				command = exec.Command("cmd", "/c", joinCmd(cmd))
			} else {
				command = exec.Command("sh", "-c", joinCmd(cmd))
			}
		} else {
			command = exec.Command(shell, "-c", joinCmd(cmd))
		}
		if opts.Cwd != "" {
			command.Dir = opts.Cwd
		}
		if opts.Env != nil {
			env := os.Environ()
			for k, v := range opts.Env {
				env = append(env, k+"="+v)
			}
			command.Env = env
		}
	}

	applyPlatformAttrs(command)

	return command, nil
}

func joinCmd(cmd []string) string {
	result := ""
	for i, arg := range cmd {
		if i > 0 {
			result += " "
		}
		if shellEscapingNeeded(arg) {
			result += "'" + arg + "'"
		} else {
			result += arg
		}
	}
	return result
}

func shellEscapingNeeded(s string) bool {
	for _, c := range s {
		if c == ' ' || c == '\'' || c == '"' || c == '$' || c == '`' || c == '\\' {
			return true
		}
	}
	return false
}

func Spawn(cmd []string, opts Options) (*Child, error) {
	command, err := buildCmd(cmd, opts)
	if err != nil {
		return nil, err
	}

	var stdoutBuf, stderrBuf bytes.Buffer

	if opts.Stdout == StdioPipe || opts.Stdout == "" {
		command.Stdout = &stdoutBuf
	}
	if opts.Stderr == StdioPipe || opts.Stderr == "" {
		command.Stderr = &stderrBuf
	}
	if opts.Stdin == StdioPipe {
		command.Stdin = nil
	}

	child := &Child{
		Cmd:  command,
		done: make(chan int, 1),
	}
	child.Exit = child.done

	err = command.Start()
	if err != nil {
		return nil, err
	}

	go func() {
		command.Wait()
		code := 0
		if command.ProcessState != nil {
			code = command.ProcessState.ExitCode()
		}
		child.done <- code
	}()

	return child, nil
}

func Run(cmd []string, opts Options) (Result, error) {
	if opts.Stdout == "" {
		opts.Stdout = StdioPipe
	}
	if opts.Stderr == "" {
		opts.Stderr = StdioPipe
	}

	command, err := buildCmd(cmd, opts)
	if err != nil {
		return Result{Code: -1}, err
	}

	var stdoutBuf, stderrBuf bytes.Buffer
	if opts.Stdout != StdioIgnore {
		command.Stdout = &stdoutBuf
	}
	if opts.Stderr != StdioIgnore {
		command.Stderr = &stderrBuf
	}

	runFunc := func() error {
		return command.Run()
	}

	if opts.Timeout > 0 {
		_, err = WithTimeout(func() (int, error) {
			runErr := runFunc()
			return 0, runErr
		}, int(opts.Timeout.Milliseconds()))
	} else {
		err = runFunc()
	}

	code := 0
	if command.ProcessState != nil {
		code = command.ProcessState.ExitCode()
	} else if err != nil {
		code = -1
	}

	result := Result{
		Code:   code,
		Stdout: stdoutBuf.Bytes(),
		Stderr: stderrBuf.Bytes(),
	}

	if code != 0 && err != nil {
		return result, &RunFailedError{
			Cmd:    cmd,
			Code:   code,
			Stdout: result.Stdout,
			Stderr: result.Stderr,
		}
	}

	return result, nil
}

func Text(cmd []string, opts Options) (TextResult, error) {
	result, err := Run(cmd, opts)
	if err != nil {
		if rfe, ok := err.(*RunFailedError); ok {
			return TextResult{
				Result: Result{
					Code:   rfe.Code,
					Stdout: rfe.Stdout,
					Stderr: rfe.Stderr,
				},
				Text: string(rfe.Stdout),
			}, err
		}
		return TextResult{Result: result}, err
	}
	return TextResult{
		Result: result,
		Text:   string(result.Stdout),
	}, nil
}

func Stop(pid int) error {
	if pid <= 0 {
		return errors.New("invalid pid")
	}

	proc, err := os.FindProcess(pid)
	if err != nil {
		return err
	}

	if runtime.GOOS == "windows" {
		taskkill := exec.Command("taskkill", "/pid", strconv.Itoa(pid), "/T", "/F")
		if err := taskkill.Run(); err == nil {
			return nil
		}
	}

	if err := proc.Signal(syscall.SIGTERM); err != nil {
		if runtime.GOOS != "windows" {
			return proc.Kill()
		}
		return err
	}

	go func() {
		time.Sleep(5 * time.Second)
		proc.Kill()
	}()

	return nil
}

type OutputCollector struct {
	Stdout bytes.Buffer
	Stderr bytes.Buffer
}

type stderrWriter struct {
	oc *OutputCollector
}

func (w stderrWriter) Write(p []byte) (n int, err error) {
	return w.oc.Stderr.Write(p)
}

func (oc *OutputCollector) Write(p []byte) (n int, err error) {
	return oc.Stdout.Write(p)
}

func RunWithCollector(cmd []string, opts Options, collector *OutputCollector) (Result, error) {
	command, err := buildCmd(cmd, opts)
	if err != nil {
		return Result{Code: -1}, err
	}

	if collector != nil {
		command.Stdout = collector
		command.Stderr = stderrWriter{oc: collector}
	}

	err = command.Run()

	code := 0
	if command.ProcessState != nil {
		code = command.ProcessState.ExitCode()
	} else if err != nil {
		code = -1
	}

	var stdout, stderr []byte
	if collector != nil {
		stdout = collector.Stdout.Bytes()
		stderr = collector.Stderr.Bytes()
	}

	result := Result{
		Code:   code,
		Stdout: stdout,
		Stderr: stderr,
	}

	if code != 0 && err != nil {
		return result, &RunFailedError{
			Cmd:    cmd,
			Code:   code,
			Stdout: stdout,
			Stderr: stderr,
		}
	}

	return result, nil
}

func Lines(text string) []string {
	var lines []string
	for _, line := range splitLines(text) {
		if line != "" {
			lines = append(lines, line)
		}
	}
	return lines
}

func splitLines(s string) []string {
	var lines []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			if i > start {
				line := s[start:i]
				if len(line) > 0 && line[len(line)-1] == '\r' {
					line = line[:len(line)-1]
				}
				lines = append(lines, line)
			} else {
				lines = append(lines, "")
			}
			start = i + 1
		}
	}
	if start < len(s) {
		line := s[start:]
		if len(line) > 0 && line[len(line)-1] == '\r' {
			line = line[:len(line)-1]
		}
		lines = append(lines, line)
	}
	return lines
}

func RunContext(ctx context.Context, cmd []string, opts Options) (Result, error) {
	command, err := buildCmd(cmd, opts)
	if err != nil {
		return Result{Code: -1}, err
	}

	var stdoutBuf, stderrBuf bytes.Buffer
	command.Stdout = &stdoutBuf
	command.Stderr = &stderrBuf

	err = command.Start()
	if err != nil {
		return Result{Code: -1}, err
	}

	done := make(chan error, 1)

	go func() {
		done <- command.Wait()
	}()

	select {
	case <-ctx.Done():
		command.Process.Kill()
		<-done
		return Result{
			Code:   -1,
			Stdout: stdoutBuf.Bytes(),
			Stderr: stderrBuf.Bytes(),
		}, ctx.Err()
	case <-done:
		code := 0
		if command.ProcessState != nil {
			code = command.ProcessState.ExitCode()
		}

		result := Result{
			Code:   code,
			Stdout: stdoutBuf.Bytes(),
			Stderr: stderrBuf.Bytes(),
		}

		if code != 0 && err != nil {
			return result, &RunFailedError{
				Cmd:    cmd,
				Code:   code,
				Stdout: result.Stdout,
				Stderr: result.Stderr,
			}
		}

		return result, err
	}
}
