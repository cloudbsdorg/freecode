package platform

import (
	"testing"
)

func TestDetect(t *testing.T) {
	os := Detect()
	if os == "" {
		t.Error("Detect() returned empty string")
	}
}

func TestName(t *testing.T) {
	name := Name()
	if name == "" {
		t.Error("Name() returned empty string")
	}
}

func TestIsSupported(t *testing.T) {
	supported := IsSupported()
	if !supported {
		t.Error("IsSupported() should return true for current platform")
	}
}

func TestConfigureShell(t *testing.T) {
	err := ConfigureShell()
	if err != nil {
		t.Errorf("ConfigureShell() error = %v", err)
	}
}

func TestVersion(t *testing.T) {
	version, err := Version()
	if err != nil {
		t.Errorf("Version() error = %v", err)
	}
	if version == "" {
		t.Error("Version() returned empty string")
	}
}

func TestInstallDeps(t *testing.T) {
	err := InstallDeps()
	if err != nil {
		t.Errorf("InstallDeps() error = %v", err)
	}
}