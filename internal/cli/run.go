package cli

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/freecode/freecode/internal/args"
	"github.com/freecode/freecode/internal/ui"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run [message..]",
	Short: "Run freecode with a message",
	Long:  `Start a new freecode session or continue an existing one.`,
	RunE:  runRun,
}

var (
	continueLast bool
	sessionID    string
	forkSession  bool
	shareSession bool
	model        string
	agentName    string
	format       string
	attachFiles  []string
	sessionTitle string
	attachServer string
	password     string
	workDir      string
	serverPort   int
	modelVariant string
	showThinking bool
	skipPerms    bool
	prompt       string
)

func init() {
	runCmd.Flags().BoolVarP(&continueLast, "continue", "c", false, "Continue last session")
	runCmd.Flags().StringVarP(&sessionID, "session", "s", "", "Session ID")
	runCmd.Flags().BoolVar(&forkSession, "fork", false, "Fork session")
	runCmd.Flags().BoolVar(&shareSession, "share", false, "Share session")
	runCmd.Flags().StringVarP(&model, "model", "m", "", "Model (provider/model)")
	runCmd.Flags().StringVar(&agentName, "agent", "", "Agent to use")
	runCmd.Flags().StringVar(&format, "format", "default", "Output format")
	runCmd.Flags().StringSliceVarP(&attachFiles, "file", "f", nil, "Files to attach")
	runCmd.Flags().StringVar(&sessionTitle, "title", "", "Session title")
	runCmd.Flags().StringVar(&attachServer, "attach", "", "Attach to remote server")
	runCmd.Flags().StringVarP(&password, "password", "p", "", "Auth password")
	runCmd.Flags().StringVar(&workDir, "dir", "", "Working directory")
	runCmd.Flags().IntVar(&serverPort, "port", 0, "Local server port")
	runCmd.Flags().StringVar(&modelVariant, "variant", "", "Model variant")
	runCmd.Flags().BoolVar(&showThinking, "thinking", false, "Show thinking blocks")
	runCmd.Flags().BoolVar(&skipPerms, "dangerously-skip-permissions", false, "Skip permission checks")
	runCmd.Flags().BoolVar(&yolo, "yolo", false, "Skip all confirmations")
	runCmd.Flags().StringVar(&prompt, "prompt", "", "Prompt to execute")
}

func runRun(cmd *cobra.Command, cmdArgs []string) error {
	tuiArgs := args.Args{
		Continue:  continueLast,
		SessionID: sessionID,
		Agent:     agentName,
		Model:     model,
		Prompt:    prompt,
		Fork:      forkSession,
	}
	p := tea.NewProgram(ui.NewModel(tuiArgs), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("failed to start TUI: %w", err)
	}
	return nil
}
