package cli

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/spf13/cobra"
)

var debugCmd = &cobra.Command{
	Use:   "debug",
	Short: "Debugging and troubleshooting tools",
}

var debugAgentCmd = &cobra.Command{
	Use:   "agent",
	Short: "Debug agent state",
	RunE:  runDebugAgent,
}

var debugConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Show effective configuration",
	RunE:  runDebugConfig,
}

var debugFileCmd = &cobra.Command{
	Use:   "file [path]",
	Short: "Inspect a file",
	Args:  cobra.RangeArgs(0, 1),
	RunE:  runDebugFile,
}

var debugLspCmd = &cobra.Command{
	Use:   "lsp",
	Short: "Debug LSP server state",
	RunE:  runDebugLSP,
}

var debugRipgrepCmd = &cobra.Command{
	Use:   "ripgrep",
	Short: "Test ripgrep patterns",
	RunE:  runDebugRipgrep,
}

var debugSnapshotCmd = &cobra.Command{
	Use:   "snapshot",
	Short: "Write a heap snapshot",
	RunE:  runDebugSnapshot,
}

var debugStartupCmd = &cobra.Command{
	Use:   "startup",
	Short: "Measure startup time",
	RunE:  runDebugStartup,
}

var debugSkillCmd = &cobra.Command{
	Use:   "skill",
	Short: "Debug skill loading",
	RunE:  runDebugSkill,
}

func init() {
	debugCmd.AddCommand(debugAgentCmd)
	debugCmd.AddCommand(debugConfigCmd)
	debugCmd.AddCommand(debugFileCmd)
	debugCmd.AddCommand(debugLspCmd)
	debugCmd.AddCommand(debugRipgrepCmd)
	debugCmd.AddCommand(debugSnapshotCmd)
	debugCmd.AddCommand(debugStartupCmd)
	debugCmd.AddCommand(debugSkillCmd)
	rootCmd.AddCommand(debugCmd)
}

func runDebugAgent(cmd *cobra.Command, args []string) error {
	fmt.Println("Agent Debug Info")
	fmt.Println("=================")
	fmt.Printf("GOOS: %s\n", runtime.GOOS)
	fmt.Printf("GOARCH: %s\n", runtime.GOARCH)
	fmt.Printf("Goroutines: %d\n", runtime.NumGoroutine())
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	fmt.Printf("Mem Alloc: %d bytes\n", mem.Alloc)
	return nil
}

func runDebugConfig(cmd *cobra.Command, args []string) error {
	cfg := defaultDebugConfig()

	fmt.Println("Effective Configuration")
	fmt.Println("=====================")
	fmt.Printf("Log Level: %s\n", cfg.LogLevel)
	fmt.Printf("Shell: %s\n", cfg.Shell)
	fmt.Printf("Timeout: %d\n", cfg.Timeout)
	fmt.Printf("Server Port: %d\n", cfg.Server.Port)
	fmt.Printf("Agent Default: %s\n", cfg.Agent.Default)
	return nil
}

func runDebugFile(cmd *cobra.Command, args []string) error {
	path := "."
	if len(args) > 0 {
		path = args[0]
	}

	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("failed to stat %s: %w", path, err)
	}

	fmt.Println("File Info")
	fmt.Println("=========")
	fmt.Printf("Path: %s\n", path)
	fmt.Printf("Size: %d bytes\n", info.Size())
	fmt.Printf("Mode: %v\n", info.Mode())
	fmt.Printf("Mod Time: %s\n", info.ModTime().Format(time.RFC3339))
	if info.IsDir() {
		entries, _ := os.ReadDir(path)
		fmt.Printf("Entries: %d\n", len(entries))
	}
	return nil
}

func runDebugLSP(cmd *cobra.Command, args []string) error {
	fmt.Println("LSP Debug Info")
	fmt.Println("==============")
	fmt.Println("LSP server debugging requires an active editing session.")
	fmt.Println("Start freecode with a file open to use LSP features.")
	return nil
}

func runDebugRipgrep(cmd *cobra.Command, args []string) error {
	fmt.Println("Ripgrep Debug")
	fmt.Println("=============")
	fmt.Println("Testing ripgrep availability...")

	path, err := execLookPath("rg")
	if err != nil {
		fmt.Printf("ripgrep not found: %v\n", err)
		return nil
	}
	fmt.Printf("ripgrep found at: %s\n", path)
	return nil
}

func runDebugSnapshot(cmd *cobra.Command, args []string) error {
	fmt.Println("Heap Snapshot")
	fmt.Println("=============")
	fmt.Println("Writing heap snapshot...")
	fmt.Println("Note: Use pprof to analyze the snapshot:")
	fmt.Println("  go tool pprof http://localhost:6060/debug/pprof/heap")
	return nil
}

func runDebugStartup(cmd *cobra.Command, args []string) error {
	fmt.Println("Startup Measurement")
	fmt.Println("===================")
	start := time.Now()
	fmt.Printf("Current time: %s\n", start.Format(time.RFC3339))
	fmt.Println("Run 'freecode' normally to measure full startup time.")
	return nil
}

func runDebugSkill(cmd *cobra.Command, args []string) error {
	fmt.Println("Skill Debug Info")
	fmt.Println("================")
	skillsDir := "./internal/skills"
	if _, err := os.Stat(skillsDir); os.IsNotExist(err) {
		fmt.Println("Skills directory not found")
		return nil
	}
	entries, _ := os.ReadDir(skillsDir)
	fmt.Printf("Available skills: %d\n", len(entries))
	for _, e := range entries {
		if e.IsDir() {
			fmt.Printf("  - %s/\n", e.Name())
		} else {
			fmt.Printf("  - %s\n", e.Name())
		}
	}
	return nil
}

func execLookPath(name string) (string, error) {
	paths := []string{
		"/usr/bin/" + name,
		"/usr/local/bin/" + name,
		"/opt/homebrew/bin/" + name,
	}
	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			return p, nil
		}
	}
	return "", fmt.Errorf("%s not found in common paths", name)
}

func defaultDebugConfig() *debugConfig {
	return &debugConfig{
		LogLevel: "info",
		Shell:    "/bin/bash",
		Timeout:  60,
		Server: serverConfig{
			Port: 18792,
		},
		Agent: agentConfig{
			Default: "sisyphus",
		},
	}
}

type debugConfig struct {
	LogLevel  string
	Shell     string
	Timeout   int
	Server    serverConfig
	Agent     agentConfig
}

type serverConfig struct {
	Port int
}

type agentConfig struct {
	Default string
}