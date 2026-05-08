package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
	"sync"

	"github.com/spf13/cobra"
)

type MCPServer struct {
	Name       string `json:"name"`
	Package    string `json:"package"`
	Version    string `json:"version"`
	Installed  bool   `json:"installed"`
	Running    bool   `json:"running"`
	PID        int    `json:"pid,omitempty"`
}

type MCPServerManager struct {
	mu      sync.RWMutex
	servers map[string]*MCPServer
	running map[string]int
}

var mcpManager = &MCPServerManager{
	servers: make(map[string]*MCPServer),
	running: make(map[string]int),
}

func mcpConfigDir() string {
	usr, _ := user.Current()
	return filepath.Join(usr.HomeDir, ".config", "freecode", "mcp")
}

func mcpServersFile() string {
	return filepath.Join(mcpConfigDir(), "servers.json")
}

func mcpPIDFile(name string) string {
	return filepath.Join(mcpConfigDir(), "pids", name+".pid")
}

func (m *MCPServerManager) loadServers() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	data, err := os.ReadFile(mcpServersFile())
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	var servers []MCPServer
	if err := json.Unmarshal(data, &servers); err != nil {
		return err
	}

	m.servers = make(map[string]*MCPServer)
	for i := range servers {
		m.servers[servers[i].Name] = &servers[i]
	}

	return nil
}

func (m *MCPServerManager) saveServers() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if err := os.MkdirAll(mcpConfigDir(), 0755); err != nil {
		return err
	}

	var servers []MCPServer
	for _, s := range m.servers {
		servers = append(servers, *s)
	}

	data, err := json.MarshalIndent(servers, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(mcpServersFile(), data, 0644)
}

func (m *MCPServerManager) List() []*MCPServer {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for name, pid := range m.running {
		if srv, ok := m.servers[name]; ok {
			srv.Running = isProcessRunning(pid)
			srv.PID = pid
		}
	}

	var result []*MCPServer
	for _, srv := range m.servers {
		s := *srv
		if pid, ok := m.running[s.Name]; ok && isProcessRunning(pid) {
			s.Running = true
			s.PID = pid
		}
		result = append(result, &s)
	}
	return result
}

func (m *MCPServerManager) Get(name string) (*MCPServer, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	srv, ok := m.servers[name]
	return srv, ok
}

func (m *MCPServerManager) Add(srv *MCPServer) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.servers[srv.Name] = srv
	return m.saveServers()
}

func (m *MCPServerManager) Remove(name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.servers, name)
	delete(m.running, name)
	return m.saveServers()
}

func (m *MCPServerManager) SetRunning(name string, pid int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.running[name] = pid
}

func (m *MCPServerManager) ClearRunning(name string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.running, name)
}

func isProcessRunning(pid int) bool {
	if pid <= 0 {
		return false
	}
	process, err := os.FindProcess(pid)
	if err != nil {
		return false
	}
	err = process.Signal(os.Signal(nil))
	return err == nil
}

var mcpCmd = &cobra.Command{
	Use:   "mcp",
	Short: "Manage MCP servers",
	Long:  `List, install, and manage Model Context Protocol servers.`,
}

func init() {
	mcpManager.loadServers()

	mcpCmd.AddCommand(mcpListCmd)
	mcpCmd.AddCommand(mcpInstallCmd)
	mcpCmd.AddCommand(mcpUninstallCmd)
	mcpCmd.AddCommand(mcpStartCmd)
	mcpCmd.AddCommand(mcpStopCmd)
}

var mcpListCmd = &cobra.Command{
	Use:   "list",
	Short: "List installed MCP servers",
	RunE:  runMCPList,
}

var mcpInstallCmd = &cobra.Command{
	Use:   "install [server-name or npm-url]",
	Short: "Install an MCP server from npm",
	Args:  cobra.ExactArgs(1),
	RunE:  runMCPInstall,
}

var mcpUninstallCmd = &cobra.Command{
	Use:   "uninstall [server-name]",
	Short: "Uninstall an MCP server",
	Args:  cobra.ExactArgs(1),
	RunE:  runMCPUninstall,
}

var mcpStartCmd = &cobra.Command{
	Use:   "start [server-name]",
	Short: "Start an MCP server",
	Args:  cobra.ExactArgs(1),
	RunE:  runMCPStart,
}

var mcpStopCmd = &cobra.Command{
	Use:   "stop [server-name]",
	Short: "Stop a running MCP server",
	Args:  cobra.ExactArgs(1),
	RunE:  runMCPStop,
}

func runMCPList(cmd *cobra.Command, args []string) error {
	servers := mcpManager.List()

	if len(servers) == 0 {
		fmt.Println("No MCP servers installed.")
		fmt.Println("")
		fmt.Printf("  Install an MCP server: freecode mcp install <package>\n")
		fmt.Printf("  MCP config directory: %s\n", mcpConfigDir())
		return nil
	}

	fmt.Println("Installed MCP Servers:")
	fmt.Println("")
	fmt.Printf("  %-20s %-30s %-10s %s\n", "NAME", "PACKAGE", "VERSION", "STATUS")
	fmt.Println("  " + strings.Repeat("-", 70))

	for _, srv := range servers {
		status := "stopped"
		if srv.Running {
			status = fmt.Sprintf("running (PID %d)", srv.PID)
		}
		version := srv.Version
		if version == "" {
			version = "?"
		}
		fmt.Printf("  %-20s %-30s %-10s %s\n", srv.Name, srv.Package, version, status)
	}
	fmt.Println("")

	return nil
}

func runMCPInstall(cmd *cobra.Command, args []string) error {
	packageName := args[0]

	fmt.Printf("Installing MCP server from npm: %s\n", packageName)

	npmPath, err := exec.LookPath("npm")
	if err != nil {
		return fmt.Errorf("npm not found. Please install Node.js to install MCP servers")
	}
	fmt.Printf("Using npm at: %s\n", npmPath)

	installCmd := exec.Command(npmPath, "install", "-g", packageName)
	installCmd.Stdout = os.Stdout
	installCmd.Stderr = os.Stderr
	installCmd.Stdin = os.Stdin

	if err := installCmd.Run(); err != nil {
		return fmt.Errorf("npm install failed: %w", err)
	}

	serverName := extractServerName(packageName)

	srv := &MCPServer{
		Name:      serverName,
		Package:   packageName,
		Version:   "latest",
		Installed: true,
		Running:   false,
	}

	if err := mcpManager.Add(srv); err != nil {
		fmt.Printf("Warning: Failed to save server info: %v\n", err)
	}

	fmt.Println("")
	fmt.Printf("MCP server '%s' installed successfully!\n", serverName)
	fmt.Println("")
	fmt.Printf("  Start server: freecode mcp start %s\n", serverName)
	fmt.Printf("  List servers: freecode mcp list\n")

	return nil
}

func runMCPUninstall(cmd *cobra.Command, args []string) error {
	name := args[0]

	srv, exists := mcpManager.Get(name)
	if !exists {
		return fmt.Errorf("MCP server '%s' not found", name)
	}

	if srv.Running {
		fmt.Printf("Stopping MCP server '%s'...\n", name)
		_ = stopMCPServer(name)
	}

	if srv.Package != "" {
		fmt.Printf("Uninstalling npm package: %s\n", srv.Package)
		npmPath, err := exec.LookPath("npm")
		if err != nil {
			fmt.Printf("Warning: npm not found, skipping npm uninstall\n")
		} else {
			uninstallCmd := exec.Command(npmPath, "uninstall", "-g", srv.Package)
			uninstallCmd.Stdout = os.Stdout
			uninstallCmd.Stderr = os.Stderr
			_ = uninstallCmd.Run()
		}
	}

	if err := mcpManager.Remove(name); err != nil {
		return fmt.Errorf("failed to remove server from config: %w", err)
	}

	fmt.Printf("MCP server '%s' uninstalled successfully.\n", name)
	return nil
}

func runMCPStart(cmd *cobra.Command, args []string) error {
	name := args[0]

	srv, exists := mcpManager.Get(name)
	if !exists {
		return fmt.Errorf("MCP server '%s' not found. Install it first with: freecode mcp install <package>", name)
	}

	if srv.Running {
		return fmt.Errorf("MCP server '%s' is already running (PID %d)", name, srv.PID)
	}

	fmt.Printf("Starting MCP server '%s'...\n", name)

	binPath, err := exec.LookPath(srv.Package)
	if err != nil {
		binPath = "npx"
	}

	var runArgs []string
	if binPath == "npx" {
		runArgs = []string{"-y", srv.Package}
	}

	proc := exec.Command(binPath, runArgs...)
	proc.Stdout = os.Stdout
	proc.Stderr = os.Stderr
	proc.Stdin = os.Stdin

	if err := proc.Start(); err != nil {
		return fmt.Errorf("failed to start MCP server: %w", err)
	}

	pid := proc.Process.Pid
	mcpManager.SetRunning(name, pid)
	srv.Running = true
	srv.PID = pid

	pidDir := filepath.Join(mcpConfigDir(), "pids")
	if err := os.MkdirAll(pidDir, 0755); err == nil {
		os.WriteFile(mcpPIDFile(name), []byte(fmt.Sprintf("%d", pid)), 0644)
	}

	fmt.Printf("MCP server '%s' started (PID %d)\n", name, pid)
	return nil
}

func runMCPStop(cmd *cobra.Command, args []string) error {
	name := args[0]

	srv, exists := mcpManager.Get(name)
	if !exists {
		return fmt.Errorf("MCP server '%s' not found", name)
	}

	if !srv.Running {
		pidData, err := os.ReadFile(mcpPIDFile(name))
		if err == nil {
			var pid int
			fmt.Sscanf(string(pidData), "%d", &pid)
			if isProcessRunning(pid) {
				srv.Running = true
				srv.PID = pid
			}
		}
	}

	if !srv.Running {
		return fmt.Errorf("MCP server '%s' is not running", name)
	}

	fmt.Printf("Stopping MCP server '%s' (PID %d)...\n", name, srv.PID)

	if err := stopMCPServer(name); err != nil {
		return fmt.Errorf("failed to stop MCP server: %w", err)
	}

	srv.Running = false
	srv.PID = 0
	mcpManager.ClearRunning(name)

	os.Remove(mcpPIDFile(name))

	fmt.Printf("MCP server '%s' stopped.\n", name)
	return nil
}

func stopMCPServer(name string) error {
	srv, _ := mcpManager.Get(name)
	if srv == nil || srv.PID <= 0 {
		return fmt.Errorf("no PID found for server")
	}

	proc, err := os.FindProcess(srv.PID)
	if err != nil {
		return err
	}

	proc.Signal(os.Signal(nil))

	if isProcessRunning(srv.PID) {
		proc.Kill()
	}

	return nil
}

func extractServerName(packageOrURL string) string {
	if strings.HasPrefix(packageOrURL, "http://") || strings.HasPrefix(packageOrURL, "https://") {
		parts := strings.Split(packageOrURL, "/")
		name := parts[len(parts)-1]
		name = strings.TrimSuffix(name, ".git")
		return name
	}

	parts := strings.Split(packageOrURL, "/")
	name := parts[len(parts)-1]
	name = strings.TrimPrefix(name, "mcp-server-")
	name = strings.TrimPrefix(name, "server-")
	return name
}