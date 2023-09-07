package cmd

import (
	"fmt"

	"github.com/QU35T-code/htb-cli/utils"
	"github.com/spf13/cobra"
)

var ipCmd = &cobra.Command{
	Use:   "ip",
	Short: "Get machine IP address",
	Run: func(cmd *cobra.Command, args []string) {
		machine_id := utils.GetActiveMachineID()
		machine_id = fmt.Sprintf("%v", machine_id)
		url := "https://www.hackthebox.com/api/v4/machine/profile/" + machine_id.(string)
		resp := utils.HtbGet(url)
		info := utils.ParseJsonMessage(resp, "info")
		infomap := info.(map[string]interface{})
		fmt.Printf("Machine : %v\n\n%v", infomap["name"], infomap["ip"])
	},
}

func init() {
	rootCmd.AddCommand(ipCmd)
}
