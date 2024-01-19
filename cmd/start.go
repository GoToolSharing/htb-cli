package cmd

import (
	"fmt"
	"os"

	"github.com/GoToolSharing/htb-cli/config"
	"github.com/GoToolSharing/htb-cli/lib/start"
	"github.com/GoToolSharing/htb-cli/lib/utils"
	"github.com/GoToolSharing/htb-cli/lib/webhooks"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// startCmd defines the "start" command which initiates the starting of a specified machine.
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start a machine",
	Long:  `Starts a Hackthebox machine specified in argument`,
	Run: func(cmd *cobra.Command, args []string) {
		config.GlobalConfig.Logger.Info("Start command executed")
		machineChoosen, err := cmd.Flags().GetString("machine")
		if err != nil {
			config.GlobalConfig.Logger.Error("", zap.Error(err))
			os.Exit(1)
		}
		var machineID string
		if machineChoosen == "" {
			config.GlobalConfig.Logger.Info("Launching the machine in release arena")
			machineID, err = utils.SearchLastReleaseArenaMachine()
			if err != nil {
				config.GlobalConfig.Logger.Error("", zap.Error(err))
				os.Exit(1)
			}
			config.GlobalConfig.Logger.Debug(fmt.Sprintf("Machine ID : %s", machineID))

		}
		output, err := start.CoreStartCmd(machineChoosen, machineID)
		if err != nil {
			config.GlobalConfig.Logger.Error("", zap.Error(err))
			os.Exit(1)
		}
		fmt.Println(output)
		err = webhooks.SendToDiscord("start", output)
		if err != nil {
			config.GlobalConfig.Logger.Error("", zap.Error(err))
			os.Exit(1)
		}
		config.GlobalConfig.Logger.Info("Exit start command correctly")
	},
}

// init adds the startCmd to rootCmd and sets flags for the "start" command.
func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.Flags().StringP("machine", "m", "", "Machine name")
}
