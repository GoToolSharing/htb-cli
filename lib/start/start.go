package start

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/GoToolSharing/htb-cli/config"
	"github.com/GoToolSharing/htb-cli/lib/utils"
	"github.com/briandowns/spinner"
)

// CoreStartCmd starts a specified machine and returns a status message and any error encountered.
func CoreStartCmd(machineChoosen string, machineID string) (string, error) {
	var err error
	if machineID == "" {
		machineID, err = utils.SearchItemIDByName(machineChoosen, "Machine")
		if err != nil {
			return "", err
		}

	}
	config.GlobalConfig.Logger.Info(fmt.Sprintf("Machine ID: %s", machineID))

	machineTypeChan := make(chan string)
	machineErrChan := make(chan error)
	userSubChan := make(chan string)
	userSubErrChan := make(chan error)

	go func() {
		machineType, err := utils.GetMachineType(machineID)
		machineTypeChan <- machineType
		machineErrChan <- err
	}()

	go func() {
		userSubscription, err := utils.GetUserSubscription()
		userSubChan <- userSubscription
		userSubErrChan <- err
	}()

	machineType := <-machineTypeChan
	err = <-machineErrChan
	if err != nil {
		return "", err
	}
	config.GlobalConfig.Logger.Info(fmt.Sprintf("Machine Type: %s", machineType))

	userSubscription := <-userSubChan
	err = <-userSubErrChan
	if err != nil {
		return "", err
	}

	config.GlobalConfig.Logger.Info(fmt.Sprintf("User subscription: %s", userSubscription))

	// isActive := utils.CheckVPN()
	// if !isActive {
	// 	isConfirmed := utils.AskConfirmation("No active VPN has been detected. Would you like to start it ?", batchParam)
	// 	if isConfirmed {
	// 		utils.StartVPN(config.BaseDirectory + "/lab_QU35T3190.ovpn")
	// 	}
	// }

	var url string
	var jsonData []byte

	switch {
	case machineType == "release":
		url = config.BaseHackTheBoxAPIURL + "/arena/start"
		jsonData = []byte("{}")
	case userSubscription == "vip" || userSubscription == "vip+":
		url = config.BaseHackTheBoxAPIURL + "/vm/spawn"
		jsonData, err = json.Marshal(map[string]string{"machine_id": machineID})
		if err != nil {
			return "", fmt.Errorf("failed to create JSON data: %w", err)
		}
	default:
		url = config.BaseHackTheBoxAPIURL + "/machine/play/" + machineID
		jsonData = []byte("{}")
	}

	resp, err := utils.HtbRequest(http.MethodPost, url, jsonData)
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
	switch {
	case machineType == "release":
		s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
		utils.SetupSignalHandler(s)
		s.Suffix = " Waiting for the machine to start in order to fetch the IP address (this might take a while)."
		s.Start()
		defer s.Stop()
		timeout := time.After(10 * time.Minute)
	LoopRelease:
		for {
			select {
			case <-timeout:
				fmt.Println("Timeout (10 min) ! Exiting")
				s.Stop()
				return "", nil
			default:
				ip, err = utils.GetActiveReleaseArenaMachineIP()
				if err != nil {
					return "", err
				}
				if ip != "Undefined" {
					s.Stop()
					break LoopRelease
				}
				time.Sleep(6 * time.Second)
			}
		}
	case userSubscription == "vip+":
		s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
		utils.SetupSignalHandler(s)
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
				ip, err = utils.GetActiveMachineIP()
				if err != nil {
					return "", err
				}
				if ip != "Undefined" {
					s.Stop()
					break Loop
				}
				time.Sleep(6 * time.Second)
			}
		}
	default:
		// Get IP address from active machine
		activeMachineData, err := utils.GetInformationsFromActiveMachine()
		if err != nil {
			return "", err
		}
		ip = activeMachineData["ip"].(string)
	}

	message = fmt.Sprintf("%s\nTarget: %s", message, ip)
	return message, nil
}
