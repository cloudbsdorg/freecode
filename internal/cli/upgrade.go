package cli

import (
	"github.com/spf13/cobra"
)

var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Upgrade freecode",
	Long:  `Check for and install freecode updates.`,
}

var (
	upgradeCheck   bool
	upgradeVersion string
	upgradeForce   bool
)

func init() {
	upgradeCmd.AddCommand(upgradeCheckCmd)
	upgradeCmd.AddCommand(upgradeInstallCmd)
}

var upgradeCheckCmd = &cobra.Command{
	Use:   "check",
	Short: "Check for updates",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var upgradeInstallCmd = &cobra.Command{
	Use:   "install [version]",
	Short: "Install a specific version",
	RunE:  runUpgradeInstall,
}

func runUpgradeInstall(cmd *cobra.Command, args []string) error {
	cmd.Help()
	return nil
}
