package platform

import (
	"runtime"
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

func TestIsSupportedFalse(t *testing.T) {
	supported := IsSupported()
	if runtime.GOOS == "darwin" && !supported {
		t.Error("IsSupported() should return true for darwin")
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

func TestPreflight(t *testing.T) {
	checks := Preflight()
	if len(checks) == 0 {
		t.Error("Preflight() returned empty checks")
	}
}

func TestPreflightCheckFields(t *testing.T) {
	checks := Preflight()
	for _, check := range checks {
		if check.Name == "" {
			t.Error("Check.Name is empty")
		}
		if check.Description == "" {
			t.Error("Check.Description is empty")
		}
		if check.Severity == "" {
			t.Error("Check.Severity is empty")
		}
		if check.Check == nil {
			t.Error("Check.Check is nil")
		}
	}
}

func TestRunPreflight(t *testing.T) {
	checks := RunPreflight()
	if len(checks) == 0 {
		t.Error("RunPreflight() returned empty checks")
	}
}

func TestCheckGit(t *testing.T) {
	err := checkGit()
	if err != nil {
		t.Errorf("checkGit() error = %v", err)
	}
}

func TestCheckShell(t *testing.T) {
	err := checkShell()
	if err != nil {
		t.Errorf("checkShell() error = %v", err)
	}
}

func TestCheckGo(t *testing.T) {
	err := checkGo()
	if err != nil {
		t.Errorf("checkGo() error = %v", err)
	}
}