package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	generateType  string
	generateForce bool
	generateDir   string
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate code from templates",
	Long: `Generate code from templates for common patterns and structures.

Examples:
  freecode generate          # Show generation options
  freecode generate tool     # Generate a new tool
  freecode generate agent    # Generate a new agent
  freecode generate hook     # Generate a new hook
  freecode generate cli       # Generate a new CLI command`,
	RunE: runGenerate,
}

func init() {
	generateCmd.Flags().StringVar(&generateType, "type", "", "Generation type (tool, agent, hook, cli)")
	generateCmd.Flags().BoolVar(&generateForce, "force", false, "Overwrite existing files")
	generateCmd.Flags().StringVar(&generateDir, "dir", "", "Output directory")
}

func runGenerate(cmd *cobra.Command, args []string) error {
	if len(args) == 0 && generateType == "" {
		return showGenerationOptions()
	}

	genType := generateType
	if genType == "" && len(args) > 0 {
		genType = args[0]
	}

	switch genType {
	case "tool":
		name := ""
		if len(args) > 1 {
			name = args[1]
		}
		return generateTool(name)
	case "agent":
		name := ""
		if len(args) > 1 {
			name = args[1]
		}
		return generateAgent(name)
	case "hook":
		name := ""
		if len(args) > 1 {
			name = args[1]
		}
		return generateHook(name)
	case "cli":
		name := ""
		if len(args) > 1 {
			name = args[1]
		}
		return generateCLI(name)
	default:
		return fmt.Errorf("unknown generation type: %s", genType)
	}
}

func showGenerationOptions() error {
	fmt.Println("Freecode Code Generation")
	fmt.Println("========================")
	fmt.Println("")
	fmt.Println("  Available generation types:")
	fmt.Println("")
	fmt.Println("    freecode generate tool   # Generate a new tool")
	fmt.Println("    freecode generate agent  # Generate a new agent")
	fmt.Println("    freecode generate hook   # Generate a new hook")
	fmt.Println("    freecode generate cli    # Generate a new CLI command")
	fmt.Println("")
	fmt.Println("  Options:")
	fmt.Println("    --type <type>    # Specify generation type")
	fmt.Println("    --force          # Overwrite existing files")
	fmt.Println("    --dir <path>     # Output directory")
	fmt.Println("")

	return nil
}

func generateTool(name string) error {
	if name == "" {
		name = "my-tool"
	}

	dir := generateDir
	if dir == "" {
		dir = "internal/tool"
	}

	file := filepath.Join(dir, name+".go")

	if !generateForce && fileExists(file) {
		return fmt.Errorf("file already exists: %s (use --force to overwrite)", file)
	}

	content := fmt.Sprintf(`package tool

import (
	"context"
	"fmt"
)

type %s struct{}

func New%s() *%s {
	return &%s{}
}

func (t *%s) Name() string {
	return "%s"
}

func (t *%s) Description() string {
	return "Description of %s"
}

func (t *%s) Schema() ToolSchema {
	return ToolSchema{
		Name:        "%s",
		Description: "Description of %s",
		Parameters: map[string]Parameter{},
	}
}

func (t *%s) Execute(ctx context.Context, req Request) (*Response, error) {
	return &Response{
		Result: "Executed tool",
	}, nil
}
`, titleCase(name), titleCase(name), titleCase(name), titleCase(name), titleCase(name), name, titleCase(name), name, titleCase(name), name, name, name)

	return os.WriteFile(file, []byte(content), 0644)
}

func generateAgent(name string) error {
	if name == "" {
		name = "my-agent"
	}

	dir := generateDir
	if dir == "" {
		dir = "internal/agent"
	}

	file := filepath.Join(dir, name+".go")

	if !generateForce && fileExists(file) {
		return fmt.Errorf("file already exists: %s (use --force to overwrite)", file)
	}

	content := fmt.Sprintf(`package agent

import (
	"context"
	"fmt"
)

type %sAgent struct{}

func New%sAgent() *%sAgent {
	return &%sAgent{}
}

func (a *%sAgent) Name() string {
	return "%s"
}

func (a *%sAgent) Run(ctx context.Context, req Request) (*Response, error) {
	return &Response{
		Message: Message{
			Role:    "assistant",
			Content: fmt.Sprintf("Hello from %s agent"),
		},
	}, nil
}
`, titleCase(name), titleCase(name), titleCase(name), titleCase(name), titleCase(name), name, titleCase(name), name)

	return os.WriteFile(file, []byte(content), 0644)
}

func generateHook(name string) error {
	if name == "" {
		name = "my-hook"
	}

	dir := generateDir
	if dir == "" {
		dir = "internal/hook"
	}

	file := filepath.Join(dir, name+".go")

	if !generateForce && fileExists(file) {
		return fmt.Errorf("file already exists: %s (use --force to overwrite)", file)
	}

	content := fmt.Sprintf(`package hook

func init() {
	// Register your hook here
}

// %s is a custom hook
func %s() error {
	return nil
}
`, titleCase(name), titleCase(name))

	return os.WriteFile(file, []byte(content), 0644)
}

func generateCLI(name string) error {
	if name == "" {
		name = "my-command"
	}

	dir := generateDir
	if dir == "" {
		dir = "internal/cli"
	}

	file := filepath.Join(dir, name+".go")

	if !generateForce && fileExists(file) {
		return fmt.Errorf("file already exists: %s (use --force to overwrite)", file)
	}

	longDesc := "Description of " + name + " command.\n\nExamples:\n  freecode " + name + "          # Run " + name + " command"

	content := fmt.Sprintf(`package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var %sCmd = &cobra.Command{
	Use:   "%s",
	Short: "Description of %s",
	Long:  "%s",
	RunE:  run%s,
}

func init() {
	// Add flags here
}

func run%s(cmd *cobra.Command, args []string) error {
	fmt.Println("Running %s command")
	return nil
}
`, name, name, name, longDesc, titleCase(name), titleCase(name), name)

	return os.WriteFile(file, []byte(content), 0644)
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func titleCase(s string) string {
	if len(s) == 0 {
		return s
	}
	return string(s[0]-32) + s[1:]
}
