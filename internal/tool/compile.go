package tool

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func CompileExternalTools(toolsDir string, outputDir string) error {
	if toolsDir == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to get home directory: %w", err)
		}
		toolsDir = filepath.Join(homeDir, ".config", "freecode", "tools")
	}

	if _, err := os.Stat(toolsDir); os.IsNotExist(err) {
		return nil
	}

	if outputDir == "" {
		outputDir = filepath.Join(toolsDir, "compiled")
	}

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	entries, err := os.ReadDir(toolsDir)
	if err != nil {
		return fmt.Errorf("failed to read tools directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".go" {
			continue
		}

		srcPath := filepath.Join(toolsDir, entry.Name())
		outputName := strings.TrimSuffix(entry.Name(), ".go") + ".so"
		outputPath := filepath.Join(outputDir, outputName)

		if err := compileTool(srcPath, outputPath); err != nil {
			fmt.Printf("warning: failed to compile %s: %v\n", srcPath, err)
			continue
		}

		if err := loadPlugin(outputPath); err != nil {
			fmt.Printf("warning: failed to load compiled tool %s: %v\n", outputPath, err)
			continue
		}

		fmt.Printf("compiled and loaded: %s -> %s\n", entry.Name(), outputName)
	}

	return nil
}

func compileTool(srcPath, outputPath string) error {
	cmd := exec.Command("go", "build", "-buildmode=plugin", "-o", outputPath, srcPath)
	cmd.Env = append(os.Environ(), "GO111MODULE=on")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("go build failed: %w (output: %s)", err, string(output))
	}
	return nil
}

type ExternalToolSource struct {
	Path string
}

func (s *ExternalToolSource) Name() string {
	return strings.TrimSuffix(filepath.Base(s.Path), filepath.Ext(s.Path))
}

func (s *ExternalToolSource) Compile(outputPath string) error {
	return compileTool(s.Path, outputPath)
}