package cmd

import (
	"fmt"
	"os"

	"github.com/QU35T-code/htb-cli/utils"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Displays the active machine",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		url := "https://www.hackthebox.com/api/v4/machine/active"
		resp := utils.HtbGet(url)
		info := utils.ParseJsonMessage(resp, "info")
		if info == nil {
			url = "https://www.hackthebox.com/api/v4/release_arena/active"
			resp = utils.HtbGet(url)
			info = utils.ParseJsonMessage(resp, "info")
			if info == nil {
				fmt.Println("No machine is active")
				os.Exit(1)
			}
			infomap := info.(map[string]interface{})
			fmt.Printf("--- STATUS ---\nId : %v\nName : %v\nPlan : %v\nServer : %v\nStatus %v", infomap["id"], infomap["name"], infomap["type"], infomap["lab_server"], "Running")
			os.Exit(0)
		}
		infomap := info.(map[string]interface{})
		fmt.Printf("--- STATUS ---\nId : %v\nName : %v\nPlan : %v\nServer : %v\nStatus %v", infomap["id"], infomap["name"], infomap["type"], infomap["lab_server"], "Running")
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
