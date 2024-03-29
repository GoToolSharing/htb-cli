package cmd

import (
	"fmt"
	"os"

	"github.com/GoToolSharing/htb-cli/config"
	"github.com/GoToolSharing/htb-cli/lib/webhooks"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// Displays the current version of htb-cli
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Displays the current version of htb-cli",
	Run: func(cmd *cobra.Command, args []string) {
		config.GlobalConfig.Logger.Info("Version command executed")
		config.GlobalConfig.Logger.Debug(fmt.Sprintf("config.Version: %s", config.Version))
		var message string
		if config.Version == "dev" {
			message = "Development version (dev branch)"
		} else {
			message = fmt.Sprintf("Stable version (main branch): %s", config.Version)
		}

		fmt.Println(message)
		err := webhooks.SendToDiscord("version", message)
		if err != nil {
			config.GlobalConfig.Logger.Error("", zap.Error(err))
			os.Exit(1)
		}

		config.GlobalConfig.Logger.Info("Exit version command correctly")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
