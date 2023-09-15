package cmd

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/GoToolSharing/htb-cli/utils"
	"github.com/spf13/cobra"
)

var difficultyParam int
var machineNameParam string
var challengeNameParam string
var flagParam string

var submitCmd = &cobra.Command{
	Use:   "submit",
	Short: "Submit credentials (User and Root Flags) - [WIP]",
	Long:  "This command allows for the submission of user and root flags discovered on vulnerable machines.",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if verbose {
			log.SetOutput(os.Stdout)
		} else {
			log.SetOutput(io.Discard)
		}
		if machineNameParam == "" && challengeNameParam == "" {
			return fmt.Errorf("either -f or -c is required")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if difficultyParam < 1 || difficultyParam > 10 {
			fmt.Println("Difficulty must set between 1 and 10")
			os.Exit(1)
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
			log.Fatal(err)
		}
		message := utils.ParseJsonMessage(resp, "message")
		fmt.Println(message)
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
