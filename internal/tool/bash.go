package tool

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

type BashTool struct{}

func NewBashTool() *BashTool {
	return &BashTool{}
}

func (t *BashTool) Name() string {
	return "bash"
}

func (t *BashTool) Description() string {
	return "Execute bash commands in a shell"
}

func (t *BashTool) Schema() ToolSchema {
	return ToolSchema{
		Name:        "bash",
		Description: "Execute bash commands in a shell",
		Parameters: map[string]Parameter{
			"command": {
				Type:        "string",
				Description: "The bash command to execute",
				Required:    true,
			},
			"timeout": {
				Type:        "integer",
				Description: "Timeout in seconds",
				Default:     60,
			},
			"workdir": {
				Type:        "string",
				Description: "Working directory",
			},
		},
	}
}

func (t *BashTool) Execute(ctx context.Context, req Request) (*Response, error) {
	cmdStr, ok := req.Arguments["command"].(string)
	if !ok {
		return nil, fmt.Errorf("command must be a string")
	}

	timeoutSecs := 60
	if to, ok := req.Arguments["timeout"].(int); ok {
		timeoutSecs = to
	}

	workdir := ""
	if wd, ok := req.Arguments["workdir"].(string); ok {
		workdir = wd
	}

	cmd := exec.Command("bash", "-c", cmdStr)
	if workdir != "" {
		cmd.Dir = workdir
	}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	done := make(chan error, 1)
	go func() {
		done <- cmd.Run()
	}()

	select {
	case <-ctx.Done():
		if cmd.Process != nil {
			cmd.Process.Kill()
		}
		return nil, ctx.Err()
	case err := <-done:
		_ = time.Duration(timeoutSecs) * time.Second
		if err != nil {
			return &Response{
				Result: strings.TrimSpace(stderr.String()),
				Error:  err,
			}, nil
		}
		return &Response{
			Result: strings.TrimSpace(stdout.String()),
		}, nil
	}
}
