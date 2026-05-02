package platform

import (
	"os/exec"
	"runtime"
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

func InstallDeps() error {
	switch runtime.GOOS {
	case "freebsd":
		cmd := exec.Command("pkg", "install", "-y", "go124")
		return cmd.Run()
	case "linux":
		cmd := exec.Command("apt-get", "update")
		if err := cmd.Run(); err != nil {
			return err
		}
		cmd = exec.Command("apt-get", "install", "-y", "build-essential")
		return cmd.Run()
	case "darwin":
		return nil
	default:
		return nil
	}
}

func ConfigureShell() error {
	return nil
}
