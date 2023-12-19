package cmd

import (
	"github.com/GoToolSharing/htb-cli/config"
	"github.com/GoToolSharing/htb-cli/lib/update"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Check for updates",
	Run: func(cmd *cobra.Command, args []string) {
		config.GlobalConfig.Logger.Info("Update command executed")
		update.Check(config.Version)
		config.GlobalConfig.Logger.Info("Exit update command correctly")
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
