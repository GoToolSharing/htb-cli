package cmd

import (
	"os"

	"github.com/GoToolSharing/htb-cli/config"
	"github.com/GoToolSharing/htb-cli/lib/play"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var playCmd = &cobra.Command{
	Use:   "play",
	Short: "Configure the pentest environment",
	Run: func(cmd *cobra.Command, args []string) {
		config.GlobalConfig.Logger.Info("Play command executed")
		releaseParam, err := cmd.Flags().GetBool("release")
		if err != nil {
			config.GlobalConfig.Logger.Error("", zap.Error(err))
			os.Exit(1)
		}
		play.Configure(releaseParam)
		config.GlobalConfig.Logger.Info("Exit play command correctly")
	},
}

func init() {
	rootCmd.AddCommand(playCmd)
	playCmd.Flags().BoolP("release", "", false, "Release arena")
}
