package util

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestWhich(t *testing.T) {
	switch runtime.GOOS {
	case "windows":
		t.Run("Windows", testWhichWindows)
	default:
		t.Run("Unix", testWhichUnix)
	}
}

func testWhichUnix(t *testing.T) {
	t.Parallel()

	env := map[string]string{"PATH": "/usr/bin:/bin"}

	result := Which("ls", env)
	if result == "" {
		t.Fatal("expected to find ls")
	}
	if !filepath.IsAbs(result) {
		t.Errorf("expected absolute path, got %s", result)
	}
	if filepath.Base(result) != "ls" {
		t.Errorf("expected ls, got %s", filepath.Base(result))
	}
}

func testWhichWindows(t *testing.T) {
	t.Parallel()

	env := map[string]string{
		"PATH":    `C:\Windows\System32`,
		"PATHEXT": `.exe;.cmd;.bat;.com`,
	}

	result := Which("cmd", env)
	if result == "" {
		t.Fatal("expected to find cmd")
	}
	if !filepath.IsAbs(result) {
		t.Errorf("expected absolute path, got %s", result)
	}
}

func TestWhichNotFound(t *testing.T) {
	t.Parallel()

	env := map[string]string{"PATH": "/nonexistent"}
	result := Which("nonexistent_command_12345", env)
	if result != "" {
		t.Errorf("expected empty string, got %s", result)
	}
}

func TestWhichEmptyCmd(t *testing.T) {
	t.Parallel()

	result := Which("", map[string]string{"PATH": "/bin"})
	if result != "" {
		t.Errorf("expected empty string for empty cmd, got %s", result)
	}
}

func TestWhichWithPath(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	execPath := filepath.Join(tmpDir, "mytestcmd")
	if runtime.GOOS == "windows" {
		execPath += ".exe"
	}

	if err := os.WriteFile(execPath, []byte("#!/bin/sh\n"), 0755); err != nil {
		t.Fatalf("failed to write test executable: %v", err)
	}

	env := map[string]string{"PATH": tmpDir}
	result := Which(filepath.Base(execPath), env)
	if result == "" {
		t.Fatal("expected to find mytestcmd")
	}
	if result != execPath {
		t.Errorf("expected %s, got %s", execPath, result)
	}
}

func TestWhichFallsBackToOSEnv(t *testing.T) {
	t.Parallel()

	result := Which("sh", map[string]string{})
	if result == "" {
		t.Skip("sh not found in system PATH")
	}
	if !filepath.IsAbs(result) {
		t.Errorf("expected absolute path, got %s", result)
	}
}

func TestWhichPrefersProvidedEnv(t *testing.T) {
	t.Parallel()

	env := map[string]string{"PATH": "/nonexistent"}
	result := Which("sh", env)
	if result != "" {
		t.Errorf("expected not to find sh in /nonexistent, got %s", result)
	}
}

func TestWhichAbsolutePath(t *testing.T) {
	t.Parallel()

	if runtime.GOOS == "windows" {
		t.Skip("Windows absolute paths have different handling")
	}

	tmpDir := t.TempDir()
	execPath := filepath.Join(tmpDir, "mytestcmd")
	if err := os.WriteFile(execPath, []byte("#!/bin/sh\n"), 0755); err != nil {
		t.Fatalf("failed to write test executable: %v", err)
	}

	result := Which(execPath, map[string]string{})
	if result == "" {
		t.Fatal("expected to find executable at absolute path")
	}
}

func TestWhichDirectory(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	env := map[string]string{"PATH": tmpDir}

	result := Which(tmpDir, env)
	if result != "" {
		t.Errorf("expected empty string for directory, got %s", result)
	}
}

func TestWhichWindowsPATHEXT(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("Windows-specific test")
	}

	t.Parallel()

	tmpDir := t.TempDir()
	execPath := filepath.Join(tmpDir, "testcmd.cmd")

	if err := os.WriteFile(execPath, []byte("@echo off\n"), 0644); err != nil {
		t.Fatalf("failed to write test executable: %v", err)
	}

	env := map[string]string{
		"PATH":    tmpDir,
		"PATHEXT": `.exe;.cmd;.bat`,
	}

	result := Which("testcmd", env)
	if result == "" {
		t.Fatal("expected to find testcmd.cmd")
	}
	if filepath.Ext(result) != ".cmd" {
		t.Errorf("expected .cmd extension, got %s", filepath.Ext(result))
	}
}

func TestWhichWindowsWithExtension(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("Windows-specific test")
	}

	t.Parallel()

	tmpDir := t.TempDir()
	execPath := filepath.Join(tmpDir, "testcmd.exe")

	if err := os.WriteFile(execPath, []byte("@echo off\n"), 0644); err != nil {
		t.Fatalf("failed to write test executable: %v", err)
	}

	env := map[string]string{"PATH": tmpDir}

	result := Which("testcmd.exe", env)
	if result == "" {
		t.Fatal("expected to find testcmd.exe")
	}
}

func TestGetEnv(t *testing.T) {
	t.Parallel()

	env := map[string]string{"PATH": "/usr/bin", "Path": "/local/bin"}

	result := getEnv(env, "PATH", "Path")
	if result != "/usr/bin" {
		t.Errorf("expected /usr/bin, got %s", result)
	}

	result = getEnv(env, "Path", "PATH")
	if result != "/local/bin" {
		t.Errorf("expected /local/bin, got %s", result)
	}

	result = getEnv(env, "NONEXISTENT")
	if result != "" {
		t.Errorf("expected empty string, got %s", result)
	}

	result = getEnv(env, "NONEXISTENT", "ALSO_NOT_THERE")
	if result != "" {
		t.Errorf("expected empty string, got %s", result)
	}
}

func TestIsExecutable(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()

	regularFile := filepath.Join(tmpDir, "regular")
	if err := os.WriteFile(regularFile, []byte("test"), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	if isExecutable(regularFile) {
		t.Error("expected regular file without execute bit to not be executable")
	}

	execFile := filepath.Join(tmpDir, "exec")
	if err := os.WriteFile(execFile, []byte("test"), 0755); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	if !isExecutable(execFile) {
		t.Error("expected file with execute bit to be executable")
	}

	dir := filepath.Join(tmpDir, "dir")
	if err := os.Mkdir(dir, 0755); err != nil {
		t.Fatalf("failed to create test directory: %v", err)
	}

	if isExecutable(dir) {
		t.Error("expected directory to not be executable")
	}

	if isExecutable(filepath.Join(tmpDir, "nonexistent")) {
		t.Error("expected nonexistent file to not be executable")
	}
}

func TestWhichPathWithEmptyEntries(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	execPath := filepath.Join(tmpDir, "emptypathtest")
	if runtime.GOOS == "windows" {
		execPath += ".exe"
	}
	if err := os.WriteFile(execPath, []byte("#!/bin/sh\n"), 0755); err != nil {
		t.Fatalf("failed to write test executable: %v", err)
	}

	env := map[string]string{"PATH": ":" + tmpDir}
	result := Which(filepath.Base(execPath), env)
	if result == "" {
		t.Error("expected to find executable with empty PATH entries")
	}
}

func BenchmarkWhich(b *testing.B) {
	env := map[string]string{"PATH": "/usr/bin:/bin:/usr/local/bin"}
	for i := 0; i < b.N; i++ {
		Which("ls", env)
	}
}

func BenchmarkWhichNotFound(b *testing.B) {
	env := map[string]string{"PATH": "/nonexistent"}
	for i := 0; i < b.N; i++ {
		Which("nonexistent_command", env)
	}
}

func BenchmarkBuiltinLookPath(b *testing.B) {
	for i := 0; i < b.N; i++ {
		os.Stat("ls")
	}
}