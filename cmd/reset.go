package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/GoToolSharing/htb-cli/utils"
	"github.com/spf13/cobra"
)

// coreResetCmd sends a reset request for an active machine.
func coreResetCmd(proxyParam string) (string, error) {
	// Retrieve the ID of the active machine.
	machineID := utils.GetActiveMachineID(proxyParam)
	if machineID == "" {
		return "No active machine found", nil
	}
	log.Printf("Machine ID: %s", machineID)
	
	// Retrieve the type of the machine.
	machineType := utils.GetMachineType(machineID, "")
	log.Printf("Machine Type: %s", machineType)
	
	// Determine the API endpoint and construct JSON data based on the machine type.
	var endpoint string
	switch machineType {
	case "active":
		endpoint = "/vm/reset"
	default:
		endpoint = "/arena/reset"
	}
	url := baseAPIURL + endpoint

	// Construct JSON data.
	jsonData, err := json.Marshal(map[string]string{"machine_id": machineID})
	if err != nil {
		return "", fmt.Errorf("failed to create JSON data: %w", err)
	}

	// Send an HTTP request to reset the machine.
	resp, err := utils.HtbRequest(http.MethodPost, url, proxyParam, jsonData)
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
		output, err := coreResetCmd(proxyParam)
		if err != nil {
			log.Fatalf("Error: %v", err)
		}
		fmt.Println(output)
	},
}

// init adds the resetCmd to rootCmd, making it callable.
func init() {
	rootCmd.AddCommand(resetCmd)
}
