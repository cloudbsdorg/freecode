package cli

import (
	"bytes"
	"os"
	"testing"

	"github.com/freecode/freecode/internal/config"
	"github.com/freecode/freecode/internal/provider"
)

func TestExecute(t *testing.T) {
	if os.Getenv("CI") != "" {
		t.Skip("skipping test that requires TTY in CI")
	}

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

func TestRunRun(t *testing.T) {
	buf := &bytes.Buffer{}
	runCmd.SetOut(buf)
	runCmd.SetErr(buf)

	err := runRun(runCmd, []string{"hello"})
	if err != nil {
		t.Errorf("runRun() error = %v", err)
	}
}

func TestRunServe(t *testing.T) {
	t.Skip("runServe blocks indefinitely waiting for SIGINT - tested manually only")

	buf := &bytes.Buffer{}
	serveCmd.SetOut(buf)
	serveCmd.SetErr(buf)

	err := runServe(serveCmd, []string{})
	if err != nil {
		t.Errorf("runServe() error = %v", err)
	}
}

func TestRunSessionRead(t *testing.T) {
	buf := &bytes.Buffer{}
	sessionCmd.SetOut(buf)
	sessionCmd.SetErr(buf)

	err := runSessionRead(sessionCmd, []string{})
	if err != nil {
		t.Errorf("runSessionRead() error = %v", err)
	}
}

func TestRunSessionSearch(t *testing.T) {
	buf := &bytes.Buffer{}
	sessionCmd.SetOut(buf)
	sessionCmd.SetErr(buf)

	err := runSessionSearch(sessionCmd, []string{"test"})
	if err != nil {
		t.Errorf("runSessionSearch() error = %v", err)
	}
}

func TestRunSessionExport(t *testing.T) {
	buf := &bytes.Buffer{}
	sessionCmd.SetOut(buf)
	sessionCmd.SetErr(buf)

	err := runSessionExport(sessionCmd, []string{})
	if err != nil {
		t.Errorf("runSessionExport() error = %v", err)
	}
}

func TestRunSessionImport(t *testing.T) {
	buf := &bytes.Buffer{}
	sessionCmd.SetOut(buf)
	sessionCmd.SetErr(buf)

	err := runSessionImport(sessionCmd, []string{})
	if err != nil {
		t.Errorf("runSessionImport() error = %v", err)
	}
}

func TestRunSessionDelete(t *testing.T) {
	buf := &bytes.Buffer{}
	sessionCmd.SetOut(buf)
	sessionCmd.SetErr(buf)

	err := runSessionDelete(sessionCmd, []string{})
	if err != nil {
		t.Errorf("runSessionDelete() error = %v", err)
	}
}

func TestRunSessionInfo(t *testing.T) {
	buf := &bytes.Buffer{}
	sessionCmd.SetOut(buf)
	sessionCmd.SetErr(buf)

	err := runSessionInfo(sessionCmd, []string{})
	if err != nil {
		t.Errorf("runSessionInfo() error = %v", err)
	}
}

func TestRunTabNew(t *testing.T) {
	buf := &bytes.Buffer{}
	tabCmd.SetOut(buf)
	tabCmd.SetErr(buf)

	err := runTabNew(tabCmd, []string{})
	if err != nil {
		t.Errorf("runTabNew() error = %v", err)
	}
}

func TestRunTabClose(t *testing.T) {
	buf := &bytes.Buffer{}
	tabCmd.SetOut(buf)
	tabCmd.SetErr(buf)

	err := runTabClose(tabCmd, []string{})
	if err != nil {
		t.Errorf("runTabClose() error = %v", err)
	}
}

func TestRunTabMove(t *testing.T) {
	buf := &bytes.Buffer{}
	tabCmd.SetOut(buf)
	tabCmd.SetErr(buf)

	err := runTabMove(tabCmd, []string{})
	if err != nil {
		t.Errorf("runTabMove() error = %v", err)
	}
}

func TestRunTabRename(t *testing.T) {
	buf := &bytes.Buffer{}
	tabCmd.SetOut(buf)
	tabCmd.SetErr(buf)

	err := runTabRename(tabCmd, []string{})
	if err != nil {
		t.Errorf("runTabRename() error = %v", err)
	}
}

func TestRunMCPInstall(t *testing.T) {
	buf := &bytes.Buffer{}
	mcpCmd.SetOut(buf)
	mcpCmd.SetErr(buf)

	err := runMCPInstall(mcpCmd, []string{})
	if err != nil {
		t.Errorf("runMCPInstall() error = %v", err)
	}
}

func TestRunMCPUninstall(t *testing.T) {
	buf := &bytes.Buffer{}
	mcpCmd.SetOut(buf)
	mcpCmd.SetErr(buf)

	err := runMCPUninstall(mcpCmd, []string{})
	if err != nil {
		t.Errorf("runMCPUninstall() error = %v", err)
	}
}

func TestRunMCPStart(t *testing.T) {
	buf := &bytes.Buffer{}
	mcpCmd.SetOut(buf)
	mcpCmd.SetErr(buf)

	err := runMCPStart(mcpCmd, []string{})
	if err != nil {
		t.Errorf("runMCPStart() error = %v", err)
	}
}

func TestRunMCPStop(t *testing.T) {
	buf := &bytes.Buffer{}
	mcpCmd.SetOut(buf)
	mcpCmd.SetErr(buf)

	err := runMCPStop(mcpCmd, []string{})
	if err != nil {
		t.Errorf("runMCPStop() error = %v", err)
	}
}

func TestRunUpgradeInstall(t *testing.T) {
	buf := &bytes.Buffer{}
	upgradeCmd.SetOut(buf)
	upgradeCmd.SetErr(buf)

	err := runUpgradeInstall(upgradeCmd, []string{})
	if err != nil {
		t.Errorf("runUpgradeInstall() error = %v", err)
	}
}

func TestRunModels(t *testing.T) {
	buf := &bytes.Buffer{}
	modelsCmd.SetOut(buf)
	modelsCmd.SetErr(buf)

	err := runModels(modelsCmd, []string{})
	if err != nil {
		t.Errorf("runModels() error = %v", err)
	}
}

func TestRunModelsWithProvider(t *testing.T) {
	buf := &bytes.Buffer{}
	modelsCmd.SetOut(buf)
	modelsCmd.SetErr(buf)
	modelsCmd.Flags().Set("provider", "openai")

	err := runModels(modelsCmd, []string{})
	if err != nil {
		t.Errorf("runModels(--provider openai) error = %v", err)
	}

	modelsCmd.Flags().Set("provider", "")
}

func TestRunModelsWithRefresh(t *testing.T) {
	buf := &bytes.Buffer{}
	modelsCmd.SetOut(buf)
	modelsCmd.SetErr(buf)
	modelsCmd.Flags().Set("refresh", "true")

	err := runModels(modelsCmd, []string{})
	if err != nil {
		t.Errorf("runModels(--refresh) error = %v", err)
	}

	modelsCmd.Flags().Set("refresh", "false")
}

func TestDiscoverConnectedProviders(t *testing.T) {
	cfg := config.DefaultConfig()
	providers := discoverConnectedProviders(cfg)
	if providers == nil {
		t.Log("discoverConnectedProviders returned nil (no API keys set)")
	}
}

func TestRefreshProviders(t *testing.T) {
	cfg := config.DefaultConfig()
	svc := provider.NewCatalogService()
	err := refreshProviders(cfg, svc)
	if err != nil {
		t.Logf("refreshProviders returned error (expected with no API keys): %v", err)
	}
}

func TestDisplayModels(t *testing.T) {
	providers := map[string]*provider.ProviderModels{}
	displayModels(providers, false)
	displayModels(providers, true)
}

func TestModelsCmdHasFlags(t *testing.T) {
	if modelsCmd.Flags().Lookup("provider") == nil {
		t.Error("modelsCmd should have --provider flag")
	}
	if modelsCmd.Flags().Lookup("refresh") == nil {
		t.Error("modelsCmd should have --refresh flag")
	}
	if modelsCmd.Flags().Lookup("list") == nil {
		t.Error("modelsCmd should have --list flag")
	}
}
