package cli

import (
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/freecode/freecode/internal/agent"
	"github.com/spf13/cobra"
)

var agentCmd = &cobra.Command{
	Use:   "agent",
	Short: "Manage agents",
	Long:  `List, start, and manage freecode agents.`,
}

var agentState = &AgentState{
	agents: make(map[string]*AgentInfo),
}

type AgentState struct {
	mu     sync.RWMutex
	agents map[string]*AgentInfo
}

type AgentInfo struct {
	Name      string
	Status    string
	StartedAt *time.Time
	Health    string
}

func (s *AgentState) SetStatus(name, status string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	info, ok := s.agents[name]
	if !ok {
		info = &AgentInfo{Name: name}
		s.agents[name] = info
	}
	info.Status = status
}

func (s *AgentState) GetInfo(name string) (*AgentInfo, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	info, ok := s.agents[name]
	return info, ok
}

func (s *AgentState) ListAgents() []*AgentInfo {
	s.mu.RLock()
	defer s.mu.RUnlock()
	agents := make([]*AgentInfo, 0, len(s.agents))
	for _, info := range s.agents {
		agents = append(agents, info)
	}
	return agents
}

func (s *AgentState) Start(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := agent.GetAgentConfig(name); !ok {
		return fmt.Errorf("agent not found: %s", name)
	}

	info, ok := s.agents[name]
	if !ok {
		info = &AgentInfo{Name: name}
		s.agents[name] = info
	}

	if info.Status == "running" {
		return fmt.Errorf("agent already running: %s", name)
	}

	info.Status = "running"
	info.Health = "healthy"
	now := time.Now()
	info.StartedAt = &now
	return nil
}

func (s *AgentState) Stop(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	info, ok := s.agents[name]
	if !ok {
		return fmt.Errorf("agent not found: %s", name)
	}

	if info.Status != "running" {
		return fmt.Errorf("agent not running: %s", name)
	}

	info.Status = "stopped"
	info.Health = "unknown"
	return nil
}

func init() {
	agentCmd.AddCommand(agentListCmd)
	agentCmd.AddCommand(agentStartCmd)
	agentCmd.AddCommand(agentStopCmd)
	agentCmd.AddCommand(agentStatusCmd)
}

var agentListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all available built-in agents",
	RunE:  runAgentList,
}

func runAgentList(cmd *cobra.Command, args []string) error {
	agents := agent.ListAgents()

	names := make([]string, 0, len(agents))
	for name := range agents {
		names = append(names, name)
	}
	sort.Strings(names)

	fmt.Println("Available Built-in Agents:")
	fmt.Println("")
	fmt.Printf("%-20s %-15s %s\n", "Name", "Mode", "Description")
	fmt.Println("----------------------------------------------------------------")

	for _, name := range names {
		cfg := agents[name]
		mode := "unknown"
		switch cfg.Mode {
		case agent.AgentModePrimary:
			mode = "primary"
		case agent.AgentModeSubagent:
			mode = "subagent"
		case agent.AgentModeAll:
			mode = "all"
		}
		fmt.Printf("%-20s %-15s %s\n", cfg.Name, mode, cfg.Description)
	}

	fmt.Println("")
	fmt.Println("Use 'freecode agent start <name>' to start an agent")
	return nil
}

var agentStartCmd = &cobra.Command{
	Use:   "start [agent-name]",
	Short: "Start a freecode agent",
	Args:  cobra.ExactArgs(1),
	RunE:  runAgentStart,
}

func runAgentStart(cmd *cobra.Command, args []string) error {
	name := args[0]

	if _, ok := agent.GetAgentConfig(name); !ok {
		return fmt.Errorf("agent not found: %s", name)
	}

	if err := agentState.Start(name); err != nil {
		return err
	}

	fmt.Printf("Agent '%s' started successfully\n", name)
	return nil
}

var agentStopCmd = &cobra.Command{
	Use:   "stop [agent-name]",
	Short: "Stop a running agent",
	Args:  cobra.ExactArgs(1),
	RunE:  runAgentStop,
}

func runAgentStop(cmd *cobra.Command, args []string) error {
	name := args[0]

	if err := agentState.Stop(name); err != nil {
		return err
	}

	fmt.Printf("Agent '%s' stopped successfully\n", name)
	return nil
}

var agentStatusCmd = &cobra.Command{
	Use:   "status [agent-name]",
	Short: "Get agent health/running status",
	Args:  cobra.ExactArgs(1),
	RunE:  runAgentStatus,
}

func runAgentStatus(cmd *cobra.Command, args []string) error {
	name := args[0]

	cfg, ok := agent.GetAgentConfig(name)
	if !ok {
		return fmt.Errorf("agent not found: %s", name)
	}

	info, exists := agentState.GetInfo(name)

	if !exists || info.Status == "stopped" || info.Status == "" {
		fmt.Printf("Agent: %s\n", name)
		fmt.Printf("Status: stopped\n")
		fmt.Printf("Description: %s\n", cfg.Description)
		return nil
	}

	fmt.Printf("Agent: %s\n", name)
	fmt.Printf("Status: %s\n", info.Status)
	fmt.Printf("Health: %s\n", info.Health)
	fmt.Printf("Description: %s\n", cfg.Description)

	if info.StartedAt != nil {
		fmt.Printf("Started At: %s\n", info.StartedAt.Format(time.RFC3339))
	}

	return nil
}