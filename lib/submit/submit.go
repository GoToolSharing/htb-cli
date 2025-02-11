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

func SubmitFlag(url string, payload map[string]interface{}) (string, error) {
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

// coreSubmitCmd handles the submission of flags for machines or challenges, returning a status message or error.
func CoreSubmitCmd(difficultyParam int, modeType string, modeValue string) (string, int, error) {
	var payload map[string]interface{}
	var difficultyString string
	var url string
	var challengeID string
	var mID int

	if modeType == "challenge" {
		config.GlobalConfig.Logger.Info("Challenge submit requested")
		if difficultyParam != 0 {
			if difficultyParam < 1 || difficultyParam > 10 {
				return "", 0, errors.New("difficulty must be set between 1 and 10")
			}
			difficultyString = strconv.Itoa(difficultyParam * 10)
		}
		challenges, err := utils.SearchChallengeByName(modeValue)
		if err != nil {
			return "", 0, err
		}
		config.GlobalConfig.Logger.Debug(fmt.Sprintf("Challenge found: %v", challenges))

		// TODO: get this int
		challengeID = strconv.Itoa(challenges.ID)

		url = config.BaseHackTheBoxAPIURL + "/challenge/own"
		payload = map[string]interface{}{
			"difficulty":   difficultyString,
			"challenge_id": challengeID,
		}
	} else if modeType == "machine" {
		config.GlobalConfig.Logger.Info("Machine submit requested")
		machineID, err := utils.SearchItemIDByName(modeValue, "Machine")
		if err != nil {
			return "", 0, err
		}
		machineType, err := utils.GetMachineType(machineID)
		if err != nil {
			return "", 0, err
		}
		config.GlobalConfig.Logger.Debug(fmt.Sprintf("Machine Type: %s", machineType))

		if machineType == "release" {
			url = config.BaseHackTheBoxAPIURL + "/arena/own"
		} else {
			url = config.BaseHackTheBoxAPIURL + "/machine/own"

		}
		payload = map[string]interface{}{
			"id": machineID,
		}
	} else if modeType == "fortress" {
		config.GlobalConfig.Logger.Info("Fortress submit requested")
		fortressID, err := utils.SearchFortressID(modeValue)
		if err != nil {
			return "", 0, err
		}
		config.GlobalConfig.Logger.Debug(fmt.Sprintf("Fortress ID : %d", fortressID))
		url = fmt.Sprintf("%s/fortress/%d/flag", config.BaseHackTheBoxAPIURL, fortressID)
		payload = map[string]interface{}{}
	} else if modeType == "endgame" {
		config.GlobalConfig.Logger.Info("Endgame submit requested")
		endgameID, err := utils.SearchEndgameID(modeValue)
		if err != nil {
			return "", 0, err
		}
		config.GlobalConfig.Logger.Debug(fmt.Sprintf("Endgame ID : %d", endgameID))
		url = fmt.Sprintf("%s/endgame/%d/flag", config.BaseHackTheBoxAPIURL, endgameID)
		payload = map[string]interface{}{}
	} else if modeType == "prolab" {
		config.GlobalConfig.Logger.Info("Prolab submit requested")
		prolabID, err := utils.SearchProlabID(modeValue)
		if err != nil {
			return "", 0, err
		}
		config.GlobalConfig.Logger.Debug(fmt.Sprintf("Prolab ID : %d", prolabID))
		url = fmt.Sprintf("%s/prolab/%d/flag", config.BaseHackTheBoxAPIURL, prolabID)
		payload = map[string]interface{}{}
	} else if modeType == "release-arena" {
		config.GlobalConfig.Logger.Info("Release Arena submit requested")
		isConfirmed := utils.AskConfirmation("Would you like to submit a flag for the release arena ?")
		if !isConfirmed {
			return "", 0, nil
		}
		releaseID, err := utils.SearchLastReleaseArenaMachine()
		if err != nil {
			return "", 0, err
		}
		config.GlobalConfig.Logger.Debug(fmt.Sprintf("Release Arena ID : %s", releaseID))
		url = fmt.Sprintf("%s/arena/own", config.BaseHackTheBoxAPIURL)
		payload = map[string]interface{}{
			"id": releaseID,
		}
		mID = releaseID
	}

	fmt.Print("Flag : ")
	flagByte, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Println("Error reading flag")
		return "", 0, fmt.Errorf("error reading flag")
	}
	flagOriginal := string(flagByte)
	flag := strings.ReplaceAll(flagOriginal, " ", "")

	config.GlobalConfig.Logger.Debug(fmt.Sprintf("Flag: %s", flag))

	payload["flag"] = flag

	message, err := SubmitFlag(url, payload)
	if err != nil {
		return "", 0, err
	}
	return message, mID, nil
}

func GetAchievementLink(machineID int) (string, error) {
	resp, err := utils.HtbRequest(http.MethodGet, fmt.Sprintf("%s/user/info", config.BaseHackTheBoxAPIURL), nil)
	if err != nil {
		return "", err
	}
	info := utils.ParseJsonMessage(resp, "info")
	infoMap, _ := info.(map[string]interface{})
	config.GlobalConfig.Logger.Debug(fmt.Sprintf("User ID: %v", infoMap["id"]))
	config.GlobalConfig.Logger.Debug(fmt.Sprintf("Machine ID: %d", machineID))

	resp, err = utils.HtbRequest(http.MethodGet, fmt.Sprintf("%s/user/achievement/machine/%v/%d", config.BaseHackTheBoxAPIURL, infoMap["id"], machineID), nil)
	if err != nil {
		return "", err
	}
	_, ok := utils.ParseJsonMessage(resp, "message").(string)
	if !ok {
		return fmt.Sprintf("\nAchievement link: https://labs.hackthebox.com/achievement/machine/%v/%d", infoMap["id"], machineID), nil
	}
	return "", nil

}
