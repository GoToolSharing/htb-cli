package cmd

import (
	"fmt"

	"github.com/QU35T-code/htb-cli/utils"
	"github.com/spf13/cobra"
)

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop the current machine",
	Long:  "Stop the current machine",
	Run: func(cmd *cobra.Command, args []string) {
		machine_id := utils.GetActiveMachineID()
		machine_id_string := fmt.Sprintf("%v", machine_id)
		machine_type := utils.GetMachineType(machine_id)
		user_subscription := utils.GetUserSubscription()

		if machine_type == "release" {
			url := "https://www.hackthebox.com/api/v4/release_arena/terminate"
			var jsonData2 = []byte(`{}`)
			resp := utils.HtbPost(url, jsonData2)
			message := utils.ParseJsonMessage(resp, "message")
			fmt.Println(message)
			return
		}

		if user_subscription == "vip" {
			url := "https://www.hackthebox.com/api/v4/vm/terminate"
			var jsonData2 = []byte(`{"machine_id": ` + machine_id_string + `}`)
			resp := utils.HtbPost(url, jsonData2)
			message := utils.ParseJsonMessage(resp, "message")
			fmt.Println(message)
			return
		} else {
			url := "https://www.hackthebox.com/api/v4/machine/stop"
			var jsonData = []byte(`{"machine_id": ` + machine_id_string + `}`)
			resp := utils.HtbPost(url, jsonData)
			message := utils.ParseJsonMessage(resp, "message")
			fmt.Println(message)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)
}
