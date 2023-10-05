package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/GoToolSharing/htb-cli/utils"
	"github.com/spf13/cobra"
)

var difficultyParam int
var machineNameParam string
var challengeNameParam string
var flagParam string

// coreSubmitCmd handles the submission of flags for machines or challenges, returning a status message or error.
func coreSubmitCmd(difficultyParam int, machineNameParam string, challengeNameParam string, flagParam string, proxyParam string) (string, error) {
	if machineNameParam == "" && challengeNameParam == "" {
		return "", errors.New("either machine_name (-m) or challenge_name (-c) must be provided")
	}
	if difficultyParam < 1 || difficultyParam > 10 {
		return "", errors.New("difficulty must be set between 1 and 10")
	}

	// Common payload elements
	difficultyString := strconv.Itoa(difficultyParam * 10)
	payload := map[string]string{
		"flag":       flagParam,
		"difficulty": difficultyString,
	}

	url := ""

	// Determine the API endpoint and payload based on input parameters.
	switch {
	case machineNameParam != "":
		log.Println("Machine submit requested!")
		machineID, err := utils.SearchItemIDByName(machineNameParam, proxyParam, "Machine")
		if err != nil {
			return "", err
		}
		machineType := utils.GetMachineType(machineID, proxyParam)
		log.Printf("Machine Type: %s", machineType)

		if machineType == "release" {
			url = baseAPIURL + "/arena/own"
		} else {
			url = baseAPIURL + "/machine/own"

		}
		payload["id"] = machineID
	case challengeNameParam != "":
		log.Println("Challenge submit requested!")
		challengeID, err := utils.SearchItemIDByName(challengeNameParam, proxyParam, "Challenge")
		if err != nil {
			return "", err
		}
		url = baseAPIURL + "/challenge/own"
		payload["challenge_id"] = challengeID
	}

	log.Println("Flag :", flagParam)
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
	Short: "Submit credentials (machines & challenges)",
	Long:  "This command allows for the submission of user and root flags discovered on vulnerable machines / challenges",
	Run: func(cmd *cobra.Command, args []string) {
		output, err := coreSubmitCmd(difficultyParam, machineNameParam, challengeNameParam, flagParam, proxyParam)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(output)
	},
}

// init adds the submitCmd to rootCmd and sets flags for the "submit" command.
func init() {
	rootCmd.AddCommand(submitCmd)
	submitCmd.Flags().StringVarP(&machineNameParam, "machine_name", "m", "", "Machine Name")
	submitCmd.Flags().StringVarP(&challengeNameParam, "challenge_name", "c", "", "Challenge Name")
	submitCmd.Flags().StringVarP(&flagParam, "flag", "f", "", "Flag")
	submitCmd.Flags().IntVarP(&difficultyParam, "difficulty", "d", 0, "Difficulty")
	submitCmd.MarkFlagRequired("difficulty")
	submitCmd.MarkFlagRequired("flag")
}
