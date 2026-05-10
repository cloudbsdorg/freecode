package cli

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/freecode/freecode/internal/args"
	"github.com/freecode/freecode/internal/config"
	"github.com/freecode/freecode/internal/ui"
	"github.com/spf13/cobra"
)

var (
	cfgFile    string
	yolo       bool
	setupRun   bool
	continue_  bool // 'continue' is a reserved keyword
	tuiSession string
	tuiAgent   string
	tuiModel   string
	tuiPrompt  string
	tuiFork    bool
	tuiRenderer string
)

var rootCmd = &cobra.Command{
	Use:   "freecode",
	Short: "Freecode - AI coding assistant",
	Long: `Freecode is a platform-independent AI coding assistant.

Built with Go for FreeBSD, Linux, macOS, and IllumOS.`,
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE: func(cmd *cobra.Command, cmdArgs []string) error {
		tuiArgs := args.Args{
			Continue:  continue_,
			SessionID: tuiSession,
			Agent:     tuiAgent,
			Model:     tuiModel,
			Prompt:    tuiPrompt,
			Fork:      tuiFork,
			Setup:     setupRun,
			Renderer:  tuiRenderer,
		}
		p := tea.NewProgram(ui.NewModel(tuiArgs), tea.WithAltScreen())
		if _, err := p.Run(); err != nil {
			return fmt.Errorf("failed to start TUI: %w", err)
		}
		return nil
	},
}

func Execute() error {
	if err := rootCmd.ParseFlags(os.Args[1:]); err != nil {
	}

	if setupRun {
		paths := config.PathsGet()
		if err := paths.Ensure(); err != nil {
			return fmt.Errorf("failed to create directories: %w", err)
		}
	} else {
		opts := config.BootstrapOptions{}
		result, err := config.Bootstrap(opts)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Bootstrap error: %v\n", err)
		}

		if result != nil && result.WizardRan {
			fmt.Fprintf(os.Stderr, "Provider: %s, Model: %s\n", result.Provider, result.Model)
		}
	}

	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ~/.config/freecode/config.yaml)")
	rootCmd.PersistentFlags().BoolVar(&yolo, "yolo", false, "skip all confirmations")
	rootCmd.PersistentFlags().BoolVar(&setupRun, "setup", false, "run initial setup wizard")
	rootCmd.PersistentFlags().BoolVar(&continue_, "continue", false, "continue last session")
	rootCmd.PersistentFlags().StringVar(&tuiSession, "session", "", "session ID to resume")
	rootCmd.PersistentFlags().StringVar(&tuiAgent, "agent", "", "agent to use (e.g., sisyphus, oracle)")
	rootCmd.PersistentFlags().StringVar(&tuiModel, "model", "", "model to use (e.g., provider/model)")
	rootCmd.PersistentFlags().StringVar(&tuiPrompt, "prompt", "", "prompt to execute (non-interactive)")
	rootCmd.PersistentFlags().BoolVar(&tuiFork, "fork", false, "fork the session before executing")
	rootCmd.PersistentFlags().StringVar(&tuiRenderer, "renderer", "", "renderer to use: bubble, lcd, auto (default auto)")
	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(agentCmd)
	rootCmd.AddCommand(sessionCmd)
	rootCmd.AddCommand(tabCmd)
	rootCmd.AddCommand(mcpCmd)
	rootCmd.AddCommand(modelsCmd)
	rootCmd.AddCommand(providersCmd)
	rootCmd.AddCommand(statsCmd)
	rootCmd.AddCommand(doctorCmd)
	rootCmd.AddCommand(upgradeCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(accountCmd)
	rootCmd.AddCommand(webCmd)
	rootCmd.AddCommand(cmdCmd)
	rootCmd.AddCommand(plugCmd)
	rootCmd.AddCommand(generateCmd)
}
