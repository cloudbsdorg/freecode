package cli

import (
	"bytes"
	"testing"
)

func TestExecute(t *testing.T) {
	rootCmd.SetOut(&bytes.Buffer{})
	rootCmd.SetErr(&bytes.Buffer{})

	err := Execute()
	if err != nil {
		t.Errorf("Execute() error = %v", err)
	}
}

func TestRootCommandUse(t *testing.T) {
	if rootCmd.Use != "freecode" {
		t.Errorf("rootCmd.Use = %q, want %q", rootCmd.Use, "freecode")
	}
}

func TestRootCommandHasSubcommands(t *testing.T) {
	expected := []string{"run", "serve", "agent", "session", "tab", "mcp", "stats", "doctor", "upgrade", "version"}
	for _, name := range expected {
		cmd, _, err := rootCmd.Find([]string{name})
		if err != nil || cmd == rootCmd {
			t.Errorf("Subcommand %q not found", name)
		}
	}
}

func TestAgentCommand(t *testing.T) {
	cmd, _, err := rootCmd.Find([]string{"agent"})
	if err != nil || cmd == rootCmd {
		t.Error("agent command not found")
	}
	if len(cmd.Commands()) == 0 {
		t.Error("agent should have subcommands")
	}
}

func TestAgentListCommand(t *testing.T) {
	cmd, _, err := agentCmd.Find([]string{"list"})
	if err != nil || cmd == agentCmd {
		t.Error("list subcommand not found")
	}
}

func TestAgentStartCommand(t *testing.T) {
	cmd, _, err := agentCmd.Find([]string{"start"})
	if err != nil || cmd == agentCmd {
		t.Error("start subcommand not found")
	}
}

func TestAgentStopCommand(t *testing.T) {
	cmd, _, err := agentCmd.Find([]string{"stop"})
	if err != nil || cmd == agentCmd {
		t.Error("stop subcommand not found")
	}
}

func TestRunAgentStart(t *testing.T) {
	cmd := agentStartCmd
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)
	cmd.SetErr(buf)

	err := runAgentStart(cmd, []string{"test-agent"})
	if err != nil {
		t.Errorf("runAgentStart() error = %v", err)
	}
}

func TestRunAgentStop(t *testing.T) {
	cmd := agentStopCmd
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)
	cmd.SetErr(buf)

	err := runAgentStop(cmd, []string{"test-agent"})
	if err != nil {
		t.Errorf("runAgentStop() error = %v", err)
	}
}

func TestRunDoctor(t *testing.T) {
	err := runDoctor(doctorCmd, []string{})
	if err != nil {
		t.Errorf("runDoctor() error = %v", err)
	}
}

func TestRunDoctorWithAllFlag(t *testing.T) {
	doctorCmd.Flags().Set("all", "true")
	defer doctorCmd.Flags().Set("all", "false")

	err := runDoctor(doctorCmd, []string{})
	if err != nil {
		t.Errorf("runDoctor(--all) error = %v", err)
	}
}

func TestDoctorAllFlag(t *testing.T) {
	err := doctorCmd.Flags().Set("all", "true")
	if err != nil {
		t.Errorf("Failed to set all flag: %v", err)
	}
	if !doctorAll {
		t.Error("doctorAll should be true after setting flag")
	}
}

func TestStatsCommandHasSubcommands(t *testing.T) {
	if len(statsCmd.Commands()) == 0 {
		t.Error("statsCmd should have subcommands")
	}
}

func TestStatsFlags(t *testing.T) {
	tests := []struct {
		name  string
		flag  string
		value string
	}{
		{"sessions", "sessions", "true"},
		{"tools", "tools", "true"},
		{"agents", "agents", "true"},
		{"all", "all", "true"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := statsCmd.Flags().Set(tt.flag, tt.value)
			if err != nil {
				t.Errorf("Failed to set %s flag: %v", tt.name, err)
			}
		})
	}
}