package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/GoToolSharing/htb-cli/utils"
	"github.com/briandowns/spinner"
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
	case userSubscription == "vip" || userSubscription == "vip+":
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

	if strings.Contains(message, "You must stop") {
		return message, nil
	}

	ip := "Undefined"
	switch userSubscription {
	case "vip+":
		s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
		setupSignalHandler(s)
		s.Suffix = " Waiting for the machine to start in order to fetch the IP address (this might take a while)."
		s.Start()
		defer s.Stop()
		timeout := time.After(10 * time.Minute)
	Loop:
		for {
			select {
			case <-timeout:
				fmt.Println("Timeout (10 min) ! Exiting")
				s.Stop()
				return "", nil
			default:
				ip = utils.GetActiveMachineIP(proxyParam)
				if ip != "" {
					s.Stop()
					break Loop
				}
				time.Sleep(3 * time.Second)
			}
		}
	default:
		// Get IP address from active machine
		activeMachineData, err := utils.GetInformationsFromActiveMachine(proxyParam)
		if err != nil {
			return "", err
		}
		ip = activeMachineData["ip"].(string)
	}

	message = fmt.Sprintf("%s\nTarget: %s", message, ip)
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
