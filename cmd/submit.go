package cmd

import (
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

func coreSubmitCmd(difficultyParam int, machineNameParam string, challengeNameParam string, flagParam string, proxyParam string) (string, error) {
	if machineNameParam == "" && challengeNameParam == "" {
		return "", errors.New("Error: either -f or -c is required")
	}
	if difficultyParam < 1 || difficultyParam > 10 {
		return "", errors.New("Error: Difficulty must set between 1 and 10")
	}

	url := ""
	var jsonData = []byte("")
	difficultyString := strconv.Itoa(difficultyParam * 10)
	if machineNameParam != "" {
		log.Println("Machine submit requested !")
		machine_id := utils.SearchItemIDByName(machineNameParam, proxyParam, "Machine")
		// TODO: Add current machine ID
		// machine_id := utils.GetActiveMachineID(proxyParam)
		// log.Println("Machine ID :", machine_id)
		url = "https://www.hackthebox.com/api/v4/machine/own"
		jsonData = []byte(`{"flag": "` + flagParam + `", "id": ` + machine_id + `, "difficulty": ` + difficultyString + `}`)
	} else if challengeNameParam != "" {
		log.Println("Challenge submit requested !")
		challenge_id := utils.SearchItemIDByName(challengeNameParam, proxyParam, "Challenge")
		url = "https://www.hackthebox.com/api/v4/challenge/own"
		jsonData = []byte(`{"flag": "` + flagParam + `", "challenge_id": ` + challenge_id + `, "difficulty": ` + difficultyString + `}`)
	}

	log.Println("Flag :", flagParam)
	log.Println("Difficulty :", difficultyString)
	resp, err := utils.HtbRequest(http.MethodPost, url, proxyParam, jsonData)
	if err != nil {
		return "", err
	}
	message := utils.ParseJsonMessage(resp, "message")
	return message.(string), nil
}

var submitCmd = &cobra.Command{
	Use:   "submit",
	Short: "Submit credentials (User and Root Flags)",
	Long:  "This command allows for the submission of user and root flags discovered on vulnerable machines.",
	Run: func(cmd *cobra.Command, args []string) {
		output, err := coreSubmitCmd(difficultyParam, machineNameParam, challengeNameParam, flagParam, proxyParam)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(output)
	},
}

func init() {
	rootCmd.AddCommand(submitCmd)
	submitCmd.Flags().StringVarP(&machineNameParam, "machine_name", "m", "", "Machine Name")
	submitCmd.Flags().StringVarP(&challengeNameParam, "challenge_name", "c", "", "Challenge Name")
	submitCmd.Flags().StringVarP(&flagParam, "flag", "f", "", "Flag")
	submitCmd.Flags().IntVarP(&difficultyParam, "difficulty", "d", 0, "Difficulty")
	submitCmd.MarkFlagRequired("difficulty")
	submitCmd.MarkFlagRequired("flag")
}
