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

	if activeMachineData == nil {
		return "", nil
	}

	ip := activeMachineData["ip"].(string)
	name := activeMachineData["name"].(string)
	os := activeMachineData["os"].(string)
	authUserInUserOwns := activeMachineData["authUserInUserOwns"].(bool)
	authUserInRootOwns := activeMachineData["authUserInRootOwns"].(bool)
	stars := activeMachineData["stars"].(float64)

	fmt.Println("---- Active Machine ----")
	fmt.Printf("IP : %s\n", ip)
	fmt.Printf("Name : %s\n", name)
	fmt.Printf("OS : %s\n", os)
	fmt.Printf("Stars : %v\n", stars)

	if authUserInUserOwns && authUserInRootOwns {
		link, err := utils.GetAchievementLink(int(activeMachineData["id"].(float64)))
		if err != nil {
			config.GlobalConfig.Logger.Error("", zap.Error(err))
			return "", err
		}
		fmt.Println("The machine has been pwned !")
		fmt.Println(link)
		return "", nil
	}
	fmt.Printf("User flag : %v\n", authUserInUserOwns)
	fmt.Printf("Root flag : %v\n", authUserInRootOwns)
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
