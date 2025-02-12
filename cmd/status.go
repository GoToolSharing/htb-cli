package cmd

import (
	"fmt"
	"os"

	"github.com/GoToolSharing/htb-cli/config"
	"github.com/GoToolSharing/htb-cli/lib/utils"
	"github.com/GoToolSharing/htb-cli/lib/webhooks"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func coreStatusCmd() (string, error) {
	fmt.Println("---- VPN Status ----")

	fmt.Println("")
	activeMachineData, err := utils.GetInformationsFromActiveMachine()

	if err != nil {
		return "", err
	}

	ip := activeMachineData["ip"].(string)
	name := activeMachineData["name"].(string)
	os := activeMachineData["os"].(string)
	authUserInUserOwns := activeMachineData["authUserInUserOwns"].(bool)
	authUserInRootOwns := activeMachineData["authUserInRootOwns"].(bool)
	stars := activeMachineData["stars"].(float64)

	fmt.Println("---- Active Machine ----")
	fmt.Println(fmt.Sprintf("IP : %s", ip))
	fmt.Println(fmt.Sprintf("Name : %s", name))
	fmt.Println(fmt.Sprintf("OS : %s", os))
	fmt.Println(fmt.Sprintf("Stars : %v", stars))
	fmt.Println(fmt.Sprintf("User flag : %v", authUserInUserOwns))
	fmt.Println(fmt.Sprintf("Root flag : %v", authUserInRootOwns))
	return "", nil
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check HTB status",
	Run: func(cmd *cobra.Command, args []string) {
		config.GlobalConfig.Logger.Info("Status command executed")
		output, err := coreStatusCmd()
		if err != nil {
			config.GlobalConfig.Logger.Error("", zap.Error(err))
			os.Exit(1)
		}

		fmt.Println(output)

		err = webhooks.SendToDiscord("status", output)
		if err != nil {
			config.GlobalConfig.Logger.Error("", zap.Error(err))
			os.Exit(1)
		}
		config.GlobalConfig.Logger.Info("Exit status command correctly")
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
