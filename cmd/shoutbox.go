package cmd

import (
	"fmt"

	"github.com/GoToolSharing/htb-cli/config"
	"github.com/spf13/cobra"
)

var shoutboxCmd = &cobra.Command{
	Use:   "shoutbox",
	Short: "Displays shoutbox information in real time",
	Run: func(cmd *cobra.Command, args []string) {
		config.GlobalConfig.Logger.Info("Shoutbox command executed")
		fmt.Println("DEPRECATED: HackTheBox no longer uses the shoutbox. A similar alternative is coming soon !")
		/*err := shoutbox.ConnectToWebSocket()
		if err != nil {
			config.GlobalConfig.Logger.Error("", zap.Error(err))
			os.Exit(1)
		}*/
		config.GlobalConfig.Logger.Info("Exit shoutbox command correctly")
	},
}

func init() {
	rootCmd.AddCommand(shoutboxCmd)
}
