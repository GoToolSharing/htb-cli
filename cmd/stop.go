package cmd

import (
	"fmt"
	"log"
	"net/http"

	"github.com/GoToolSharing/htb-cli/config"
	"github.com/GoToolSharing/htb-cli/lib/utils"
	"github.com/GoToolSharing/htb-cli/lib/webhooks"
	"github.com/spf13/cobra"
)

var (
	releaseAPI = fmt.Sprintf("%s/arena/stop", config.BaseHackTheBoxAPIURL)
	vipAPI     = fmt.Sprintf("%s/vm/terminate", config.BaseHackTheBoxAPIURL)
	defaultAPI = fmt.Sprintf("%s/machine/stop", config.BaseHackTheBoxAPIURL)
)

// buildMachineStopRequest constructs the URL endpoint and JSON data payload for stopping a machine based on its type and user's subscription.
func buildMachineStopRequest(machineType string, userSubscription string, machineID string) (string, []byte) {
	var apiEndpoint string
	var jsonData []byte

	if machineType == "release" {
		return releaseAPI, []byte(`{}`)
	}

	switch userSubscription {
	case "vip", "vip+":
		apiEndpoint = vipAPI
	default:
		apiEndpoint = defaultAPI
	}

	jsonData = []byte(fmt.Sprintf(`{"machine_id": "%s"}`, machineID))
	return apiEndpoint, jsonData
}

// coreStopCmd stops the currently active machine.
// It fetches machine's ID, its type, and user's subscription to determine how to stop the machine.
func coreStopCmd() (string, error) {
	// err := utils.StopVPN()
	// if err != nil {
	// 	return "", err
	// }
	machineID := utils.GetActiveMachineID()
	if machineID == "" {
		return "No machine is running", nil
	}
	log.Println("Machine ID:", machineID)

	machineType := utils.GetMachineType(machineID)
	log.Println("Machine Type:", machineType)

	userSubscription := utils.GetUserSubscription()
	log.Println("User subscription:", userSubscription)

	apiEndpoint, jsonData := buildMachineStopRequest(machineType, userSubscription, machineID)
	resp, err := utils.HtbRequest(http.MethodPost, apiEndpoint, jsonData)
	if err != nil {
		return "", err
	}

	message, ok := utils.ParseJsonMessage(resp, "message").(string)
	if !ok {
		return "", fmt.Errorf("error parsing message from response")
	}

	// err = utils.StopVPN()
	// if err != nil {
	// 	return "", err
	// }

	return message, nil
}

// stopCmd represents the "stop" command which stops the current active machine.
var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop the current machine",
	Run: func(cmd *cobra.Command, args []string) {
		output, err := coreStopCmd()
		if err != nil {
			log.Fatal(err)
		}
		if config.GlobalConf["Discord"] != "False" {
			err := webhooks.SendToDiscord("[STOP] - " + output)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
		fmt.Println(output)
	},
}

// init initializes the command by adding it to the root command.
func init() {
	rootCmd.AddCommand(stopCmd)
}
