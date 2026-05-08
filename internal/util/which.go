package util

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func Which(cmd string, env map[string]string) string {
	if cmd == "" {
		return ""
	}

	if runtime.GOOS == "windows" {
		return whichWindows(cmd, env)
	}

	return whichUnix(cmd, env)
}

func whichUnix(cmd string, env map[string]string) string {
	if strings.ContainsAny(cmd, "/") {
		if isExecutable(cmd) {
			return filepath.Clean(cmd)
		}
		return ""
	}

	pathEnv := getEnv(env, "PATH", "Path")
	if pathEnv == "" {
		pathEnv = os.Getenv("PATH")
	}
	if pathEnv == "" {
		return ""
	}

	for _, dir := range filepath.SplitList(pathEnv) {
		if dir == "" {
			continue
		}
		path := filepath.Join(filepath.Clean(dir), cmd)
		if isExecutable(path) {
			return path
		}
	}

	return ""
}

func whichWindows(cmd string, env map[string]string) string {
	pathEnv := getEnv(env, "PATH", "Path")
	if pathEnv == "" {
		pathEnv = os.Getenv("PATH")
	}

	pathExtEnv := getEnv(env, "PATHEXT", "PathExt")
	if pathExtEnv == "" {
		pathExtEnv = os.Getenv("PATHEXT")
	}

	pathExts := filepath.SplitList(pathExtEnv)
	if len(pathExts) == 0 {
		pathExts = []string{".exe", ".cmd", ".bat", ".com"}
	}

	if strings.ContainsAny(cmd, `\/`) {
		if filepath.Ext(cmd) != "" {
			if isExecutable(cmd) {
				return filepath.Clean(cmd)
			}
		} else {
			for _, ext := range pathExts {
				if isExecutable(cmd + ext) {
					return filepath.Clean(cmd + ext)
				}
			}
		}
		return ""
	}

	for _, dir := range filepath.SplitList(pathEnv) {
		if dir == "" {
			continue
		}
		dir = filepath.Clean(dir)

		if filepath.Ext(cmd) != "" {
			if isExecutable(filepath.Join(dir, cmd)) {
				return filepath.Join(dir, cmd)
			}
		} else {
			for _, ext := range pathExts {
				path := filepath.Join(dir, cmd+ext)
				if isExecutable(path) {
					return filepath.Clean(path)
				}
			}
		}
	}

	return ""
}

func getEnv(env map[string]string, keys ...string) string {
	for _, key := range keys {
		if val, ok := env[key]; ok && val != "" {
			return val
		}
	}
	return ""
}

func isExecutable(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	if info.IsDir() {
		return false
	}

	if runtime.GOOS == "windows" {
		return true
	}

	return info.Mode().IsRegular() && (info.Mode()&0111) != 0
}