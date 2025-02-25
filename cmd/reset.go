package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/GoToolSharing/htb-cli/config"
	"github.com/GoToolSharing/htb-cli/lib/utils"
	"github.com/GoToolSharing/htb-cli/lib/webhooks"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// coreResetCmd sends a reset request for an active machine.
func coreResetCmd() (string, error) {
	// Retrieve the ID of the active machine.
	machineID, err := utils.GetActiveMachineID()
	if err != nil {
		return "", err
	}
	if machineID == 0 {
		return "No active machine found", nil
	}
	config.GlobalConfig.Logger.Debug(fmt.Sprintf("Machine ID: %d", machineID))

	// Retrieve the type of the machine.
	machineType, err := utils.GetMachineType(machineID)
	if err != nil {
		return "", err
	}
	config.GlobalConfig.Logger.Debug(fmt.Sprintf("Machine Type: %s", machineType))

	// Determine the API endpoint and construct JSON data based on the machine type.
	var endpoint string
	switch machineType {
	case "active":
		endpoint = "/vm/reset"
	default:
		endpoint = "/arena/reset"
	}
	url := config.BaseHackTheBoxAPIURL + endpoint

	// Construct JSON data.
	jsonData, err := json.Marshal(map[string]interface{}{"machine_id": machineID})
	if err != nil {
		return "", fmt.Errorf("failed to create JSON data: %w", err)
	}

	// Send an HTTP request to reset the machine.
	resp, err := utils.HtbRequest(http.MethodPost, url, jsonData)
	if err != nil {
		return "", err
	}

	// Parse and return the message from the response.
	message, ok := utils.ParseJsonMessage(resp, "message").(string)
	if !ok {
		return "", fmt.Errorf("unexpected response format")
	}
	return message, nil
}

// resetCmd defines the "reset" command, which allows the user to reset a machine.
var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset a machine",
	Long:  "Initiates a reset request for the selected machine.",
	Run: func(cmd *cobra.Command, args []string) {
		// Execute the core reset function and handle the results.
		output, err := coreResetCmd()
		if err != nil {
			config.GlobalConfig.Logger.Error("", zap.Error(err))
			os.Exit(1)
		}
		fmt.Println(output)
		err = webhooks.SendToDiscord("reset", output)
		if err != nil {
			config.GlobalConfig.Logger.Error("", zap.Error(err))
			os.Exit(1)
		}
	},
}

// init adds the resetCmd to rootCmd, making it callable.
func init() {
	rootCmd.AddCommand(resetCmd)
}
