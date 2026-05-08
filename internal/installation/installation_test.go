package installation

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestDetect(t *testing.T) {
	info, err := Detect()
	if err != nil {
		t.Fatalf("Detect() error = %v", err)
	}
	if info == nil {
		t.Fatal("Detect() returned nil")
	}
	if info.Platform == "" {
		t.Error("Platform is empty")
	}
	if info.DataDir == "" {
		t.Error("DataDir is empty")
	}
}

func TestDetectPlatform(t *testing.T) {
	info, _ := Detect()
	expectedPlatform := runtime.GOOS + "/" + runtime.GOARCH
	if info.Platform != expectedPlatform {
		t.Errorf("Platform = %q, want %q", info.Platform, expectedPlatform)
	}
}

func TestDetectDataDir(t *testing.T) {
	info, _ := Detect()

	switch runtime.GOOS {
	case "darwin":
		home, _ := os.UserHomeDir()
		expected := filepath.Join(home, "Library", "Application Support", "freecode")
		if info.DataDir != expected {
			t.Errorf("DataDir = %q, want %q", info.DataDir, expected)
		}
	case "linux", "freebsd":
		home := os.Getenv("HOME")
		expected := filepath.Join(home, ".local", "share", "freecode")
		if info.DataDir != expected {
			t.Errorf("DataDir = %q, want %q", info.DataDir, expected)
		}
	default:
		home := os.Getenv("HOME")
		expected := filepath.Join(home, ".freecode")
		if info.DataDir != expected {
			t.Errorf("DataDir = %q, want %q", info.DataDir, expected)
		}
	}
}

func TestIsInstalled(t *testing.T) {
	result := IsInstalled()
	if result {
		t.Log("IsInstalled() = true (freecode found in /usr/local/bin)")
	} else {
		t.Log("IsInstalled() = false (freecode not found in /usr/local/bin)")
	}
}

func TestInfo(t *testing.T) {
	info := &Info{
		Platform:   "linux/arm64",
		BinaryPath: "/usr/local/bin/freecode",
		DataDir:    "/home/user/.local/share/freecode",
	}

	if info.Platform != "linux/arm64" {
		t.Errorf("Platform = %q, want %q", info.Platform, "linux/arm64")
	}
	if info.BinaryPath != "/usr/local/bin/freecode" {
		t.Errorf("BinaryPath = %q, want %q", info.BinaryPath, "/usr/local/bin/freecode")
	}
	if info.DataDir != "/home/user/.local/share/freecode" {
		t.Errorf("DataDir = %q, want %q", info.DataDir, "/home/user/.local/share/freecode")
	}
}