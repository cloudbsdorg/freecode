package template

import (
	"os"
	"path/filepath"
)

type Loader struct {
	baseDir string
}

func NewLoader(baseDir string) *Loader {
	return &Loader{baseDir: baseDir}
}

func (l *Loader) Load(path string) (string, error) {
	fullPath := filepath.Join(l.baseDir, path)
	content, err := os.ReadFile(fullPath)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func (l *Loader) LoadViews() map[string]string {
	return l.loadTemplatesDir("views")
}

func (l *Loader) LoadComponents() map[string]string {
	return l.loadTemplatesDir("components")
}

func (l *Loader) loadTemplatesDir(subdir string) map[string]string {
	templates := make(map[string]string)
	dir := filepath.Join(l.baseDir, subdir)

	entries, err := os.ReadDir(dir)
	if err != nil {
		return templates
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		name := entry.Name()
		tmplPath := filepath.Join(dir, name, name+".tmpl")
		if content, err := os.ReadFile(tmplPath); err == nil {
			templates[name] = string(content)
		}
	}
	return templates
}

func (l *Loader) LoadRoot() (string, error) {
	return l.Load("ROOT.tmpl")
}

func (l *Loader) LoadView(name string) (string, error) {
	return l.Load(filepath.Join("views", name, name+".tmpl"))
}

func (l *Loader) LoadComponent(name string) (string, error) {
	return l.Load(filepath.Join("components", name, name+".tmpl"))
}
