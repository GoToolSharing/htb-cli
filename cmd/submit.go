package cmd

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/GoToolSharing/htb-cli/utils"
	"github.com/spf13/cobra"
)

var difficultyParam int
var flagParam string

var submitCmd = &cobra.Command{
	Use:   "submit",
	Short: "Submit credentials (User and Root Flags) - [WIP]",
	Long:  "This command allows for the submission of user and root flags discovered on vulnerable machines.",
	Run: func(cmd *cobra.Command, args []string) {
		if difficultyParam < 1 || difficultyParam > 10 {
			fmt.Println("Difficulty must set between 1 and 10")
			os.Exit(1)
		}
		machine_id := utils.GetActiveMachineID(proxyParam)
		log.Println("Machine ID :", machine_id)
		url := "https://www.hackthebox.com/api/v4/machine/own"
		difficultyString := strconv.Itoa(difficultyParam * 10)
		log.Println("Difficulty :", difficultyString)
		var jsonData = []byte(`{"flag": "` + flagParam + `", "id": ` + machine_id + `, "difficulty": ` + difficultyString + `}`)
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
	submitCmd.Flags().StringVarP(&flagParam, "flag", "f", "", "Flag")
	submitCmd.MarkFlagRequired("flag")
	submitCmd.Flags().IntVarP(&difficultyParam, "difficulty", "d", 0, "Difficulty")
	submitCmd.MarkFlagRequired("difficulty")
}
