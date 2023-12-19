package cmd

import (
	"fmt"

	"github.com/GoToolSharing/htb-cli/config"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Displays the current version",
	Run: func(cmd *cobra.Command, args []string) {
		config.GlobalConfig.Logger.Info("Version command executed")
		config.GlobalConfig.Logger.Debug(fmt.Sprintf("config.Version: %s", config.Version))
		if len(config.Version) == 40 {
			fmt.Println("Development version (dev branch): " + config.Version)
		} else {
			fmt.Println("Stable version (main branch): " + config.Version)
		}
		config.GlobalConfig.Logger.Info("Exit version command correctly")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
