package cmd

import (
	"fmt"
	"net/http"

	"github.com/GoToolSharing/htb-cli/config"
	"github.com/GoToolSharing/htb-cli/lib/utils"
	"github.com/spf13/cobra"
)

var pwnboxCmd = &cobra.Command{
	Use:   "pwnbox",
	Short: "Interact with the pwnbox",
	Run: func(cmd *cobra.Command, args []string) {
		locationFlag, err := cmd.Flags().GetString("location")
		if err != nil {
			fmt.Println(err)
			return
		}

		_ = locationFlag

		startFlag, err := cmd.Flags().GetBool("start")
		if err != nil {
			fmt.Println(err)
			return
		}

		_ = startFlag

		stopFlag, err := cmd.Flags().GetBool("stop")
		if err != nil {
			fmt.Println(err)
			return
		}

		// Check subscription

		// Start and stop
		if stopFlag {
			resp, err := utils.HtbRequest(http.MethodPost, config.BaseHackTheBoxAPIURL+"/pwnbox/terminate", nil)
			if err != nil {
				fmt.Println(resp)
				return
			}
			message, ok := utils.ParseJsonMessage(resp, "message").(string)
			if !ok {
				fmt.Println("unexpected response format")
				return
			}
			fmt.Println(message)
		}
	},
}

func init() {
	rootCmd.AddCommand(pwnboxCmd)
	pwnboxCmd.Flags().StringP("location", "l", "", "Pwnbox Location")
	pwnboxCmd.Flags().BoolP("start", "", false, "Start pwnbox")
	pwnboxCmd.Flags().BoolP("stop", "", false, "Stop active pwnbox")
}
