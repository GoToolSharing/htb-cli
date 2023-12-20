package cmd

import (
	"os"

	"github.com/GoToolSharing/htb-cli/config"
	"github.com/GoToolSharing/htb-cli/lib/shoutbox"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var shoutboxCmd = &cobra.Command{
	Use:   "shoutbox",
	Short: "Displays shoutbox information in real time",
	Run: func(cmd *cobra.Command, args []string) {
		config.GlobalConfig.Logger.Info("Shoutbox command executed")
		err := shoutbox.ConnectToWebSocket()
		if err != nil {
			config.GlobalConfig.Logger.Error("", zap.Error(err))
			os.Exit(1)
		}
		config.GlobalConfig.Logger.Info("Exit shoutbox command correctly")
	},
}

func init() {
	rootCmd.AddCommand(shoutboxCmd)
}
