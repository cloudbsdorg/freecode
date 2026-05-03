package installation

import (
	"os"
	"path/filepath"
	"runtime"
)

type Info struct {
	Platform   string
	BinaryPath string
	DataDir    string
}

func Detect() (*Info, error) {
	exe, err := os.Executable()
	if err != nil {
		exe = "/usr/local/bin/freecode"
	}

	var dataDir string
	switch runtime.GOOS {
	case "darwin":
		home, _ := os.UserHomeDir()
		dataDir = filepath.Join(home, "Library", "Application Support", "freecode")
	case "linux":
		dataDir = filepath.Join(os.Getenv("HOME"), ".local", "share", "freecode")
	case "freebsd":
		dataDir = filepath.Join(os.Getenv("HOME"), ".local", "share", "freecode")
	default:
		dataDir = filepath.Join(os.Getenv("HOME"), ".freecode")
	}

	return &Info{
		Platform:   runtime.GOOS + "/" + runtime.GOARCH,
		BinaryPath: exe,
		DataDir:    dataDir,
	}, nil
}

func IsInstalled() bool {
	_, err := os.Stat("/usr/local/bin/freecode")
	if err == nil {
		return true
	}
	return false
}
