package cmd

import (
	"fmt"
	"log"
	"net/http"

	"github.com/GoToolSharing/htb-cli/config"
	"github.com/GoToolSharing/htb-cli/utils"
	"github.com/spf13/cobra"
)

const (
	releaseAPI = "https://www.hackthebox.com/api/v4/arena/stop"
	vipAPI     = "https://www.hackthebox.com/api/v4/vm/terminate"
	defaultAPI = "https://www.hackthebox.com/api/v4/machine/stop"
)

// buildMachineStopRequest constructs the URL endpoint and JSON data payload for stopping a machine based on its type and user's subscription.
func buildMachineStopRequest(machineType, userSubscription, machineID, proxyParam string) (string, []byte) {
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
func coreStopCmd(proxyParam string) (string, error) {
	err := utils.StopVPN()
	if err != nil {
		return "", err
	}
	machineID := utils.GetActiveMachineID(proxyParam)
	if machineID == "" {
		return "No machine is running", nil
	}
	log.Println("Machine ID:", machineID)

	machineType := utils.GetMachineType(machineID, "")
	log.Println("Machine Type:", machineType)

	userSubscription := utils.GetUserSubscription(proxyParam)
	log.Println("User subscription:", userSubscription)

	apiEndpoint, jsonData := buildMachineStopRequest(machineType, userSubscription, machineID, proxyParam)
	resp, err := utils.HtbRequest(http.MethodPost, apiEndpoint, proxyParam, jsonData)
	if err != nil {
		return "", err
	}

	message, ok := utils.ParseJsonMessage(resp, "message").(string)
	if !ok {
		return "", fmt.Errorf("error parsing message from response")
	}

	err = utils.StopVPN()
	if err != nil {
		return "", err
	}

	return message, nil
}

// stopCmd represents the "stop" command which stops the current active machine.
var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop the current machine",
	Run: func(cmd *cobra.Command, args []string) {
		output, err := coreStopCmd(proxyParam)
		if err != nil {
			log.Fatal(err)
		}
		if config.GlobalConf["Discord"] != "False" {
			utils.SendDiscordWebhook("[STOP] - " + output)
		}
		fmt.Println(output)
	},
}

// init initializes the command by adding it to the root command.
func init() {
	rootCmd.AddCommand(stopCmd)
}
