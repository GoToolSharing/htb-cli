package cmd

import (
	"fmt"
	"net/http"
	"os"

	"github.com/GoToolSharing/htb-cli/config"
	"github.com/GoToolSharing/htb-cli/lib/utils"
	"github.com/GoToolSharing/htb-cli/lib/webhooks"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func doAction(mode string) {
}

var pwnboxCmd = &cobra.Command{
	Use:   "pwnbox",
	Short: "Interact with the pwnbox",
	Run: func(cmd *cobra.Command, args []string) {
		modeFlag, err := cmd.Flags().GetString("mode")
		if err != nil {
			fmt.Println(err)
			return
		}

		if modeFlag == "" {
			fmt.Println("You have to choose a mode to use the pwnbox (-m)")
			return
		}

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

		// _ = startFlag

		stopFlag, err := cmd.Flags().GetBool("stop")
		if err != nil {
			fmt.Println(err)
			return
		}

		if !startFlag && !stopFlag {
			fmt.Println("Wrong action supplied : --start / --stop")
			return
		}

		// Check subscription
		config.GlobalConfig.Logger.Debug(fmt.Sprintf("Mode: %s", modeFlag))

		if startFlag {
			fmt.Println("Sorry, but HackTheBox currently uses a v3 recaptcha to start a pwnbox.")
			return
		}
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

			err = webhooks.SendToDiscord("pwnbox", message)
			if err != nil {
				config.GlobalConfig.Logger.Error("", zap.Error(err))
				os.Exit(1)
			}
			return
		}

		switch modeFlag {
		case "machines":
			fmt.Println("Machines")
			doAction("machines")
			return
		case "sp":
			fmt.Println("Starting Points")
			doAction("sp")
			return
		case "fortresses":
			fmt.Println("Fortresses")
			doAction("fortresses")
			return
		case "prolabs":
			fmt.Println("Prolabs")
			doAction("prolabs")
			return
		case "seasonals":
			fmt.Println("Seasonals")
			doAction("seasonals")
			return
		default:
			fmt.Println("Available modes : machines - sp - fortresses - prolabs - seasonals")
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(pwnboxCmd)
	pwnboxCmd.Flags().StringP("location", "l", "", "Pwnbox Location")
	pwnboxCmd.Flags().BoolP("start", "", false, "Start pwnbox")
	pwnboxCmd.Flags().BoolP("stop", "", false, "Stop active pwnbox")
	pwnboxCmd.Flags().StringP("mode", "m", "", "Select mode") // TODO: Choice list
}
