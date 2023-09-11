package cmd

import (
	"fmt"
	"log"
	"net/http"

	"github.com/GoToolSharing/htb-cli/utils"
	"github.com/spf13/cobra"
)

var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset a machine - [WIP]",
	Long:  "Initiates a reset request for the selected machine.",
	Run: func(cmd *cobra.Command, args []string) {
		machine_id := utils.GetActiveMachineID(proxyParam)
		log.Println("Machine ID :", machine_id)
		machine_type := utils.GetMachineType(machine_id, "")
		log.Println("Machine Type :", machine_type)

		url := ""
		jsonData := []byte("")
		switch machine_type {
		case "active":
			url = "https://www.hackthebox.com/api/v4/vm/reset"
			jsonData = []byte(`{"machine_id": ` + machine_id + `}`)
		default:
			url = "https://www.hackthebox.com/api/v4/arena/reset"
			jsonData = []byte(`{"machine_id": ` + machine_id + `}`)
		}
		resp, err := utils.HtbRequest(http.MethodPost, url, proxyParam, jsonData)
		if err != nil {
			log.Fatal(err)
		}
		message := utils.ParseJsonMessage(resp, "message")
		fmt.Println(message)
	},
}

func init() {
	rootCmd.AddCommand(resetCmd)
}
