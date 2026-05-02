package cli

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/freecode/freecode/internal/ui"
	"github.com/spf13/cobra"
)

var (
	cfgFile string
	yolo    bool
)

var logo = `██╗     ██╗ ██████╗ ██████╗ ███████╗██████╗ ███╗   ██╗██╗███╗   ██╗ ██████╗
██║     ██║██╔════╝██╔═══██╗██╔════╝██╔══██╗████╗  ██║██║████╗  ██║██╔════╝
██║     ██║██║     ██║   ██║█████╗  ██████╔╝██╔██╗ ██║██║██╔██╗ ██║██║  ███╗
██║     ██║██║     ██║   ██║██╔══╝  ██╔══██╗██║╚██╗██║██║██║╚██╗██║██║   ██║
███████╗██║╚██████╗╚██████╔╝███████╗██║  ██║██║ ╚████║██║██║ ╚████║╚██████╔╝
╚══════╝╚═╝ ╚═════╝ ╚═════╝ ╚══════╝╚═╝  ╚═╝╚═╝  ╚═══╝╚═╝╚═╝  ╚═══╝ ╚═════╝`

var rootCmd = &cobra.Command{
	Use:   "freecode",
	Short: "Freecode - AI coding assistant",
	Long: `Freecode is a platform-independent AI coding assistant.

Built with Go for FreeBSD 16, Linux, macOS, and IllumOS.`,
	SilenceErrors: true,
	SilenceUsage:   true,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println(logo)
		fmt.Println()

		p := tea.NewProgram(ui.NewModel(), tea.WithAltScreen())
		if _, err := p.Run(); err != nil {
			return fmt.Errorf("failed to start TUI: %w", err)
		}
		return nil
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ~/.config/freecode/config.yaml)")
	rootCmd.PersistentFlags().BoolVar(&yolo, "yolo", false, "skip all confirmations")
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
}
