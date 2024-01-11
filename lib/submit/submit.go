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
func CoreSubmitCmd(difficultyParam int, modeType string, modeValue string) (string, error) {
	var payload map[string]string
	var difficultyString string
	if difficultyParam != 0 {
		if difficultyParam < 1 || difficultyParam > 10 {
			return "", errors.New("difficulty must be set between 1 and 10")
		}
		difficultyString = strconv.Itoa(difficultyParam * 10)
	}

	var url string

	if modeType == "challenge" {
		config.GlobalConfig.Logger.Info("Challenge submit requested")
		challengeID, err := utils.SearchItemIDByName(modeValue, "Challenge")
		if err != nil {
			return "", err
		}
		url = config.BaseHackTheBoxAPIURL + "/challenge/own"
		payload = map[string]string{
			"difficulty":   difficultyString,
			"challenge_id": challengeID,
		}
	} else if modeType == "machine" {
		config.GlobalConfig.Logger.Info("Machine submit requested")
		machineID, err := utils.SearchItemIDByName(modeValue, "Machine")
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
		payload = map[string]string{
			"difficulty": difficultyString,
			"id":         machineID,
		}
	} else if modeType == "fortress" {
		config.GlobalConfig.Logger.Info("Fortress submit requested")
		fortressID, err := utils.SearchFortressID(modeValue)
		if err != nil {
			return "", err
		}
		config.GlobalConfig.Logger.Debug(fmt.Sprintf("Fortress ID : %d", fortressID))
		url = fmt.Sprintf("%s/fortress/%d/flag", config.BaseHackTheBoxAPIURL, fortressID)
		payload = map[string]string{}
	} else if modeType == "endgame" {
		config.GlobalConfig.Logger.Info("Endgame submit requested")
		endgameID, err := utils.SearchEndgameID(modeValue)
		if err != nil {
			return "", err
		}
		config.GlobalConfig.Logger.Debug(fmt.Sprintf("Endgame ID : %d", endgameID))
		url = fmt.Sprintf("%s/endgame/%d/flag", config.BaseHackTheBoxAPIURL, endgameID)
		payload = map[string]string{}
	} else if modeType == "release-arena" {
		config.GlobalConfig.Logger.Info("Release Arena submit requested")
		isConfirmed := utils.AskConfirmation("Would you like to submit a flag for the release arena ?")
		if !isConfirmed {
			return "", nil
		}
		releaseID, err := utils.SearchLastReleaseArenaMachine()
		if err != nil {
			return "", err
		}
		config.GlobalConfig.Logger.Debug(fmt.Sprintf("Release Arena ID : %s", releaseID))
		url = fmt.Sprintf("%s/arena/own", config.BaseHackTheBoxAPIURL)
		payload = map[string]string{
			"difficulty": difficultyString,
			"id":         releaseID,
		}
	}

	fmt.Print("Flag : ")
	flagByte, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Println("Error reading flag")
		return "", fmt.Errorf("error reading flag")
	}
	flagOriginal := string(flagByte)
	flag := strings.ReplaceAll(flagOriginal, " ", "")

	config.GlobalConfig.Logger.Debug(fmt.Sprintf("Flag: %s", flag))

	payload["flag"] = flag

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
