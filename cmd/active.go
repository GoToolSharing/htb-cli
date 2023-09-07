package cmd

import (
	"fmt"

	"github.com/QU35T-code/htb-cli/utils"
	"github.com/spf13/cobra"
)

var activeCmd = &cobra.Command{
	Use:   "active",
	Short: "List of active machines",
	Run: func(cmd *cobra.Command, args []string) {
		url := "https://www.hackthebox.com/api/v4/machine/list"
		resp := utils.HtbGet(url)
		info := utils.ParseJsonMessage(resp, "info")
		for _, value := range info.([]interface{}) {
			data := value.(map[string]interface{})
			fmt.Printf("Name : %v\nOS : %v\nDifficulty : %v\n\n", data["name"], data["os"], data["difficultyText"])
		}
	},
}

func init() {
	rootCmd.AddCommand(activeCmd)
}
