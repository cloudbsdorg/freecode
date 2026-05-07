package omo

import (
	"testing"

	"github.com/freecode/freecode/internal/config"
)

func TestOMOConfigStruct(t *testing.T) {
	cfg := &OMOConfig{
		Version:    "1.0",
		SkillsDir:  "/tmp/skills",
		SlopRemove: true,
	}

	if cfg.Version != "1.0" {
		t.Errorf("Version = %q, want %q", cfg.Version, "1.0")
	}
	if cfg.SkillsDir != "/tmp/skills" {
		t.Errorf("SkillsDir = %q, want %q", cfg.SkillsDir, "/tmp/skills")
	}
	if !cfg.SlopRemove {
		t.Error("SlopRemove should be true")
	}
}

func TestMerge(t *testing.T) {
	omoCfg := &OMOConfig{
		SkillsDir: "/tmp/skills",
	}
	cfg := config.DefaultConfig()

	Merge(cfg, omoCfg)

	if cfg.Platform.CacheDir != "/tmp/skills" {
		t.Errorf("CacheDir = %q, want %q", cfg.Platform.CacheDir, "/tmp/skills")
	}
}

func TestMergeWithSlopRemove(t *testing.T) {
	omoCfg := &OMOConfig{
		SlopRemove: true,
	}
	cfg := config.DefaultConfig()
	initialLen := len(cfg.Hooks.Session)

	Merge(cfg, omoCfg)

	if len(cfg.Hooks.Session) != initialLen+1 {
		t.Errorf("Hooks.Session length = %d, want %d", len(cfg.Hooks.Session), initialLen+1)
	}
}

func TestMergeIntoNonexistentPath(t *testing.T) {
	cfg := config.DefaultConfig()
	err := MergeInto(cfg, "/nonexistent/path/omo.jsonc")
	if err == nil {
		t.Error("MergeInto() expected error for nonexistent path")
	}
}

func TestReadEmptyPath(t *testing.T) {
	loader := &MockLoader{
		Config: &OMOConfig{},
	}
	cfg, err := loader.Read("")
	if err != nil {
		t.Errorf("MockLoader.Read() error = %v", err)
	}
	if cfg == nil {
		t.Error("MockLoader.Read() returned nil config")
	}
}

func TestReadWithMockLoader(t *testing.T) {
	loader := &MockLoader{
		Config: &OMOConfig{
			Version:    "2.0",
			SkillsDir:  "/custom/skills",
			SlopRemove: true,
		},
	}

	cfg, err := loader.Read("/test/path")
	if err != nil {
		t.Errorf("Read() error = %v", err)
	}
	if cfg.Version != "2.0" {
		t.Errorf("Version = %q, want %q", cfg.Version, "2.0")
	}
	if cfg.SkillsDir != "/custom/skills" {
		t.Errorf("SkillsDir = %q, want %q", cfg.SkillsDir, "/custom/skills")
	}
	if !cfg.SlopRemove {
		t.Error("SlopRemove should be true")
	}
}

func TestMockLoaderError(t *testing.T) {
	loader := &MockLoader{
		Err: ErrMock,
	}

	_, err := loader.Read("/test/path")
	if err != ErrMock {
		t.Errorf("Read() error = %v, want %v", err, ErrMock)
	}
}

var ErrMock = errMock{}

type errMock struct{}

func (e errMock) Error() string { return "mock error" }

func TestReadNonexistentPath(t *testing.T) {
	_, err := Read("/nonexistent/path/omo.jsonc")
	if err == nil {
		t.Error("Read() expected error for nonexistent path")
	}
}

func TestDefaultOMOPath(t *testing.T) {
	path := defaultOMOPath()
	if path == "" {
		t.Log("No OMO config found (expected on clean system)")
	}
}

func TestViperLoaderRead(t *testing.T) {
	loader := NewViperLoader()
	_, err := loader.Read("/nonexistent/path/omo.jsonc")
	if err == nil {
		t.Error("ViperLoader.Read() expected error for nonexistent path")
	}
}

func TestViperLoaderReadEmpty(t *testing.T) {
	loader := NewViperLoader()
	cfg, err := loader.Read("")
	if err != nil {
		t.Errorf("ViperLoader.Read() error = %v", err)
	}
	if cfg == nil {
		t.Error("ViperLoader.Read() returned nil config")
	}
}
