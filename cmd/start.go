package cmd

import (
	"fmt"

	"github.com/QU35T-code/htb-cli/utils"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start a machine",
	Long:  `Starts a Hackthebox machine specified in argument`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("USAGE : htb-cli start MACHINE_NAME")
			return
		}
		machine_id := utils.SearchMachineIDByName(args[0])
		machine_type := utils.GetMachineType(machine_id)
		user_subscription := utils.GetUserSubscription()
		if machine_type == "release" {
			url := "https://www.hackthebox.com/api/v4/release_arena/spawn"
			var jsonData3 = []byte(`{}`)
			resp := utils.HtbPost(url, jsonData3)
			message := utils.ParseJsonMessage(resp, "message")
			fmt.Println(message)
			return
		}
		if user_subscription == "vip" {
			url := "https://www.hackthebox.com/api/v4/vm/spawn"
			var jsonData2 = []byte(`{"machine_id": ` + machine_id + `}`)
			resp := utils.HtbPost(url, jsonData2)
			message := utils.ParseJsonMessage(resp, "message")
			fmt.Println(message)
			return
		} else {
			url := "https://www.hackthebox.com/api/v4/machine/play/" + machine_id
			var jsonData = []byte("{}")
			resp := utils.HtbPost(url, jsonData)
			message := utils.ParseJsonMessage(resp, "message")
			fmt.Println(message)
		}
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
