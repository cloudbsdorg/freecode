package config

import (
	"os"
	"path/filepath"
	"strings"
)

type Paths struct {
	Home    string
	Config  string
	Data    string
	Cache   string
	State   string
	Temp    string
	Bin     string
	Log     string
}

var paths *Paths

func PathsGet() *Paths {
	if paths != nil {
		return paths
	}

	home := getEnvOr("FREECODE_HOME", "")
	if home == "" {
		home = getUserHome()
	}

	xdgConfigHome := getEnvOr("XDG_CONFIG_HOME", "")
	xdgDataHome := getEnvOr("XDG_DATA_HOME", "")
	xdgStateHome := getEnvOr("XDG_STATE_HOME", "")
	xdgCacheHome := getEnvOr("XDG_CACHE_HOME", "")

	if xdgConfigHome == "" {
		xdgConfigHome = filepath.Join(home, ".config")
	}
	if xdgDataHome == "" {
		xdgDataHome = filepath.Join(home, ".local", "share")
	}
	if xdgStateHome == "" {
		xdgStateHome = filepath.Join(home, ".local", "state")
	}
	if xdgCacheHome == "" {
		xdgCacheHome = filepath.Join(home, ".cache")
	}

	config := filepath.Join(xdgConfigHome, "freecode")
	data := filepath.Join(xdgDataHome, "freecode")
	cache := filepath.Join(xdgCacheHome, "freecode")
	state := filepath.Join(xdgStateHome, "freecode")

	paths = &Paths{
		Home:   home,
		Config: config,
		Data:   data,
		Cache:  cache,
		State:  state,
		Temp:   filepath.Join(os.TempDir(), "freecode"),
		Bin:    filepath.Join(cache, "bin"),
		Log:    filepath.Join(data, "log"),
	}

	return paths
}

func (p *Paths) Ensure() error {
	dirs := []string{
		p.Config,
		p.Data,
		p.Cache,
		p.State,
		p.Temp,
		p.Bin,
		p.Log,
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	return nil
}

func (p *Paths) ConfigFile(name string) string {
	return filepath.Join(p.Config, name)
}

func (p *Paths) DataFile(name string) string {
	return filepath.Join(p.Data, name)
}

func (p *Paths) StateFile(name string) string {
	return filepath.Join(p.State, name)
}

func (p *Paths) SessionDir() string {
	return filepath.Join(p.Data, "sessions")
}

func (p *Paths) SkillsDir() string {
	return filepath.Join(p.Config, "skills")
}

func getEnvOr(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func getUserHome() string {
	if home := os.Getenv("HOME"); home != "" {
		return home
	}
	if home := os.Getenv("USERPROFILE"); home != "" {
		return home
	}
	return "/tmp"
}

func IsFileNotFound(err error) bool {
	if err == nil {
		return false
	}
	return os.IsNotExist(err) || strings.Contains(err.Error(), "no such file")
}
