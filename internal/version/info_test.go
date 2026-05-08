package version

import (
	"runtime"
	"testing"
)

func TestGet(t *testing.T) {
	info := Get()

	if info.Version == "" {
		t.Error("Version is empty")
	}
	if info.Platform == "" {
		t.Error("Platform is empty")
	}
	if info.GoVersion == "" {
		t.Error("GoVersion is empty")
	}
}

func TestPlatform(t *testing.T) {
	p := Platform()
	expected := runtime.GOOS + "/" + runtime.GOARCH
	if p != expected {
		t.Errorf("Platform() = %q, want %q", p, expected)
	}
}

func TestGoVersion(t *testing.T) {
	v := GoVersion()
	expected := runtime.Version()
	if v != expected {
		t.Errorf("GoVersion() = %q, want %q", v, expected)
	}
}

func TestInfo(t *testing.T) {
	info := Info{
		Version:   "1.0.0",
		Platform:  "linux/arm64",
		GoVersion: "go1.24.0",
		Commit:    "abc123",
		BuildDate: "2024-01-01",
	}

	if info.Version != "1.0.0" {
		t.Errorf("Version = %q, want %q", info.Version, "1.0.0")
	}
	if info.Platform != "linux/arm64" {
		t.Errorf("Platform = %q, want %q", info.Platform, "linux/arm64")
	}
	if info.GoVersion != "go1.24.0" {
		t.Errorf("GoVersion = %q, want %q", info.GoVersion, "go1.24.0")
	}
	if info.Commit != "abc123" {
		t.Errorf("Commit = %q, want %q", info.Commit, "abc123")
	}
	if info.BuildDate != "2024-01-01" {
		t.Errorf("BuildDate = %q, want %q", info.BuildDate, "2024-01-01")
	}
}

func TestVersionVariable(t *testing.T) {
	if Version == "" {
		t.Error("Version variable is empty")
	}
}

func TestGetReturnsConsistentInfo(t *testing.T) {
	info1 := Get()
	info2 := Get()

	if info1.Version != info2.Version {
		t.Errorf("Get() returned different versions: %q vs %q", info1.Version, info2.Version)
	}
	if info1.Platform != info2.Platform {
		t.Errorf("Get() returned different platforms: %q vs %q", info1.Platform, info2.Platform)
	}
	if info1.GoVersion != info2.GoVersion {
		t.Errorf("Get() returned different GoVersions: %q vs %q", info1.GoVersion, info2.GoVersion)
	}
}