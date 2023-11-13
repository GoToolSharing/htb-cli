package cmd

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
	"github.com/GoToolSharing/htb-cli/utils"
	"github.com/GoToolSharing/htb-cli/utils/webhooks"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var difficultyParam int
var machineNameParam string
var challengeNameParam string

// coreSubmitCmd handles the submission of flags for machines or challenges, returning a status message or error.
func coreSubmitCmd(difficultyParam int, machineNameParam string, challengeNameParam string, proxyParam string) (string, error) {
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
		challengeID, err := utils.SearchItemIDByName(challengeNameParam, proxyParam, "Challenge", batchParam)
		if err != nil {
			return "", err
		}
		url = config.BaseHackTheBoxAPIURL + "/challenge/own"
		payload["challenge_id"] = challengeID
	} else if machineNameParam != "" {
		log.Println("Machine submit requested!")
		machineID, err := utils.SearchItemIDByName(machineNameParam, proxyParam, "Machine", batchParam)
		if err != nil {
			return "", err
		}
		machineType := utils.GetMachineType(machineID, proxyParam)
		log.Printf("Machine Type: %s", machineType)

		if machineType == "release" {
			url = config.BaseHackTheBoxAPIURL + "/arena/own"
		} else {
			url = config.BaseHackTheBoxAPIURL + "/machine/own"

		}
		payload["id"] = machineID
	} else if machineNameParam == "" && challengeNameParam == "" {
		machineID := utils.GetActiveMachineID(proxyParam)
		if machineID == "" {
			return "No machine is running", nil
		}
		machineType := utils.GetMachineType(machineID, proxyParam)
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

	resp, err := utils.HtbRequest(http.MethodPost, url, proxyParam, jsonData)
	if err != nil {
		return "", err
	}

	message, ok := utils.ParseJsonMessage(resp, "message").(string)
	if !ok {
		return "", errors.New("unexpected response format")
	}
	return message, nil
}

// submitCmd defines the "submit" command for submitting flags for machines or challenges.
var submitCmd = &cobra.Command{
	Use:   "submit",
	Short: "Submit credentials (machines / challenges / arena)",
	Long:  "This command allows for the submission of user and root flags discovered on vulnerable machines / challenges",
	Run: func(cmd *cobra.Command, args []string) {
		output, err := coreSubmitCmd(difficultyParam, machineNameParam, challengeNameParam, proxyParam)
		if err != nil {
			log.Fatal(err)
		}
		if config.GlobalConf["Discord"] != "False" {
			err := webhooks.SendToDiscord("[SUBMIT] - " + output)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
		fmt.Println(output)
	},
}

// init adds the submitCmd to rootCmd and sets flags for the "submit" command.
func init() {
	rootCmd.AddCommand(submitCmd)
	submitCmd.Flags().StringVarP(&machineNameParam, "machine_name", "m", "", "Machine Name")
	submitCmd.Flags().StringVarP(&challengeNameParam, "challenge_name", "c", "", "Challenge Name")
	submitCmd.Flags().IntVarP(&difficultyParam, "difficulty", "d", 0, "Difficulty")
	err := submitCmd.MarkFlagRequired("difficulty")
	if err != nil {
		fmt.Println(err)
	}
}
