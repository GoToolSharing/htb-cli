package submit

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/GoToolSharing/htb-cli/config"
	"github.com/GoToolSharing/htb-cli/lib/utils"
	"golang.org/x/term"
)

// coreSubmitCmd handles the submission of flags for machines or challenges, returning a status message or error.
func CoreSubmitCmd(difficultyParam int, machineNameParam string, challengeNameParam string) (string, error) {
	if difficultyParam < 1 || difficultyParam > 10 {
		return "", errors.New("difficulty must be set between 1 and 10")
	}

	// Common payload elements
	difficultyString := strconv.Itoa(difficultyParam * 10)
	payload := map[string]string{
		"difficulty": difficultyString,
	}

	url := ""

	if challengeNameParam != "" {
		config.GlobalConfig.Logger.Info("Challenge submit requested")
		challengeID, err := utils.SearchItemIDByName(challengeNameParam, "Challenge")
		if err != nil {
			return "", err
		}
		url = config.BaseHackTheBoxAPIURL + "/challenge/own"
		payload["challenge_id"] = challengeID
	} else if machineNameParam != "" {
		config.GlobalConfig.Logger.Info("Machine submit requested")
		machineID, err := utils.SearchItemIDByName(machineNameParam, "Machine")
		if err != nil {
			return "", err
		}
		machineType, err := utils.GetMachineType(machineID)
		if err != nil {
			return "", err
		}
		config.GlobalConfig.Logger.Debug(fmt.Sprintf("Machine Type: %s", machineType))

		if machineType == "release" {
			url = config.BaseHackTheBoxAPIURL + "/arena/own"
		} else {
			url = config.BaseHackTheBoxAPIURL + "/machine/own"

		}
		payload["id"] = machineID
	} else if machineNameParam == "" && challengeNameParam == "" {
		machineID, err := utils.GetActiveMachineID()
		if err != nil {
			return "", err
		}
		if machineID == "" {
			return "No machine is running", nil
		}
		machineType, err := utils.GetMachineType(machineID)
		if err != nil {
			return "", err
		}
		config.GlobalConfig.Logger.Debug(fmt.Sprintf("Machine Type: %s", machineType))

		if machineType == "release" {
			url = config.BaseHackTheBoxAPIURL + "/arena/own"
		} else {
			url = config.BaseHackTheBoxAPIURL + "/machine/own"

		}
		payload["id"] = machineID
	}

	fmt.Print("Flag : ")
	flagByte, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Println("Error reading flag")
		return "", fmt.Errorf("error reading flag")
	}
	flagOriginal := string(flagByte)
	flag := strings.ReplaceAll(flagOriginal, " ", "")
	payload["flag"] = flag

	config.GlobalConfig.Logger.Debug(fmt.Sprintf("Flag: %s", flag))
	config.GlobalConfig.Logger.Debug(fmt.Sprintf("Difficulty: %s", difficultyString))

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to create JSON data: %w", err)
	}

	resp, err := utils.HtbRequest(http.MethodPost, url, jsonData)
	if err != nil {
		return "", err
	}

	message, ok := utils.ParseJsonMessage(resp, "message").(string)
	if !ok {
		return "", errors.New("unexpected response format")
	}
	return message, nil
}
