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
	machineID, err := utils.GetActiveMachineID()
	if err != nil {
		return "", err
	}
	if machineID == "" {
		return "No machine is running", nil
	}
	config.GlobalConfig.Logger.Debug(fmt.Sprintf("Machine ID: %s", machineID))

	machineType, err := utils.GetMachineType(machineID)
	if err != nil {
		return "", err
	}
	config.GlobalConfig.Logger.Debug(fmt.Sprintf("Machine Type: %s", machineType))

	userSubscription, err := utils.GetUserSubscription()
	if err != nil {
		return "", err
	}
	config.GlobalConfig.Logger.Debug(fmt.Sprintf("User subscription: %s", userSubscription))

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

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop the current machine",
	Run: func(cmd *cobra.Command, args []string) {
		config.GlobalConfig.Logger.Info("Stop command executed")
		output, err := coreStopCmd()
		if err != nil {
			config.GlobalConfig.Logger.Error("", zap.Error(err))
			os.Exit(1)
		}
		err = webhooks.SendToDiscord(fmt.Sprintf("[STOP] - %s", output))
		if err != nil {
			config.GlobalConfig.Logger.Error("", zap.Error(err))
			os.Exit(1)
		}
		fmt.Println(output)
		config.GlobalConfig.Logger.Info("Exit stop command correctly")
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)
}
