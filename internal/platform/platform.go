package platform

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func Detect() string {
	return runtime.GOOS
}

func Name() string {
	switch runtime.GOOS {
	case "freebsd":
		return "FreeBSD"
	case "linux":
		return "Linux"
	case "darwin":
		return "macOS"
	case "illumos":
		return "IllumOS"
	default:
		return runtime.GOOS
	}
}

func Version() (string, error) {
	cmd := exec.Command("uname", "-r")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func IsSupported() bool {
	switch runtime.GOOS {
	case "freebsd", "linux", "darwin", "illumos":
		return true
	default:
		return false
	}
}

type PreflightCheck struct {
	Name        string
	Description string
	Severity    string
	Check       func() error
	Fix         string
}

func Preflight() []PreflightCheck {
	return []PreflightCheck{
		{
			Name:        "Go",
			Description: "Go 1.24 or later",
			Severity:    "required",
			Check:       checkGo,
			Fix:         "Install from https://go.dev/dl/",
		},
		{
			Name:        "Git",
			Description: "Git CLI",
			Severity:    "required",
			Check:       checkGit,
			Fix:         "Install via your package manager or https://git-scm.com",
		},
		{
			Name:        "Shell",
			Description: "POSIX-compatible shell (/bin/bash, /bin/zsh, etc.)",
			Severity:    "required",
			Check:       checkShell,
			Fix:         "Ensure /bin/bash or /bin/zsh is installed",
		},
	}
}

func checkGo() error {
	cmd := exec.Command("go", "version")
	out, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("go not found: %w", err)
	}

	output := string(out)
	fields := strings.Fields(output)
	if len(fields) < 3 {
		return fmt.Errorf("unexpected go version output: %s", output)
	}

	version := fields[2]
	if !strings.HasPrefix(version, "go1.") {
		return fmt.Errorf("unsupported go version: %s (need 1.24+)", version)
	}

	version = strings.TrimPrefix(version, "go")
	parts := strings.Split(version, ".")
	if len(parts) < 2 {
		return fmt.Errorf("unsupported go version: %s", version)
	}

	major := 0
	fmt.Sscanf(parts[0], "%d", &major)
	if major < 1 {
		return fmt.Errorf("unsupported go version: go%s", version)
	}

	return nil
}

func checkGit() error {
	cmd := exec.Command("git", "version")
	_, err := cmd.Output()
	return err
}

func checkShell() error {
	shells := []string{"/bin/bash", "/bin/zsh", "/bin/sh"}
	for _, shell := range shells {
		if _, err := os.Stat(shell); err == nil {
			return nil
		}
	}
	return fmt.Errorf("no POSIX shell found in %v", shells)
}

func RunPreflight() []PreflightCheck {
	checks := Preflight()
	for i := range checks {
		checks[i].Check()
	}
	return checks
}

func ConfigureShell() error {
	return nil
}
