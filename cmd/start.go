package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/GoToolSharing/htb-cli/utils"
	"github.com/spf13/cobra"
)

var machineChoosen string

// coreStartCmd starts a specified machine and returns a status message and any error encountered.
func coreStartCmd(machineChoosen string, proxyParam string) (string, error) {
	machineID, err := utils.SearchItemIDByName(machineChoosen, proxyParam, "Machine", batchParam)
	if err != nil {
		return "", err
	}
	log.Printf("Machine ID: %s", machineID)

	machineType := utils.GetMachineType(machineID, proxyParam)
	log.Printf("Machine Type: %s", machineType)

	userSubscription := utils.GetUserSubscription(proxyParam)
	log.Printf("User subscription: %s", userSubscription)

	var url string
	var jsonData []byte

	switch {
	case machineType == "release":
		url = baseAPIURL + "/arena/start"
		jsonData = []byte("{}")
	case userSubscription == "vip":
		url = baseAPIURL + "/vm/spawn"
		jsonData, err = json.Marshal(map[string]string{"machine_id": machineID})
		if err != nil {
			return "", fmt.Errorf("failed to create JSON data: %w", err)
		}
	default:
		url = baseAPIURL + "/machine/play/" + machineID
		jsonData = []byte("{}")
	}

	resp, err := utils.HtbRequest(http.MethodPost, url, proxyParam, jsonData)
	if err != nil {
		return "", err
	}

	message, ok := utils.ParseJsonMessage(resp, "message").(string)
	if !ok {
		return "", fmt.Errorf("unexpected response format")
	}
	return message, nil
}

// startCmd defines the "start" command which initiates the starting of a specified machine.
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start a machine",
	Long:  `Starts a Hackthebox machine specified in argument`,
	Run: func(cmd *cobra.Command, args []string) {
		output, err := coreStartCmd(machineChoosen, proxyParam)
		if err != nil {
			log.Fatalf("Error: %v", err)
		}
		fmt.Println(output)
	},
}

// init adds the startCmd to rootCmd and sets flags for the "start" command.
func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.Flags().StringVarP(&machineChoosen, "machine", "m", "", "Machine name")
	startCmd.MarkFlagRequired("machine")
}
