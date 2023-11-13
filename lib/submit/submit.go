package submit

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/GoToolSharing/htb-cli/config"
	"github.com/GoToolSharing/htb-cli/lib/utils"
	"golang.org/x/term"
)

// coreSubmitCmd handles the submission of flags for machines or challenges, returning a status message or error.
func coreSubmitCmd(difficultyParam int, machineNameParam string, challengeNameParam string) (string, error) {
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
		log.Println("Challenge submit requested!")
		challengeID, err := utils.SearchItemIDByName(challengeNameParam, "Challenge")
		if err != nil {
			return "", err
		}
		url = config.BaseHackTheBoxAPIURL + "/challenge/own"
		payload["challenge_id"] = challengeID
	} else if machineNameParam != "" {
		log.Println("Machine submit requested!")
		machineID, err := utils.SearchItemIDByName(machineNameParam, "Machine")
		if err != nil {
			return "", err
		}
		machineType := utils.GetMachineType(machineID)
		log.Printf("Machine Type: %s", machineType)

		if machineType == "release" {
			url = config.BaseHackTheBoxAPIURL + "/arena/own"
		} else {
			url = config.BaseHackTheBoxAPIURL + "/machine/own"

		}
		payload["id"] = machineID
	} else if machineNameParam == "" && challengeNameParam == "" {
		machineID := utils.GetActiveMachineID()
		if machineID == "" {
			return "No machine is running", nil
		}
		machineType := utils.GetMachineType(machineID)
		log.Printf("Machine Type: %s", machineType)

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

	log.Println("Flag :", flag)
	log.Println("Difficulty :", difficultyString)

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
