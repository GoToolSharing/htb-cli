package cmd

import (
	"fmt"
	"log"
	"net/http"

	"github.com/GoToolSharing/htb-cli/utils"
	"github.com/spf13/cobra"
)

var machineChoosen string

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start a machine",
	Long:  `Starts a Hackthebox machine specified in argument`,
	Run: func(cmd *cobra.Command, args []string) {
		machine_id := utils.SearchMachineIDByName(machineChoosen, proxyParam)
		log.Println("Machine ID :", machine_id)
		machine_type := utils.GetMachineType(machine_id, proxyParam)
		log.Println("Machine Type :", machine_type)
		user_subscription := utils.GetUserSubscription(proxyParam)
		log.Println("User subscription :", user_subscription)
		if machine_type == "release" {
			url := "https://www.hackthebox.com/api/v4/arena/start"
			resp, err := utils.HtbRequest(http.MethodPost, url, proxyParam, []byte(`{}`))
			if err != nil {
				log.Fatal(err)
			}
			message := utils.ParseJsonMessage(resp, "message")
			fmt.Println(message)
			return
		}

		url := ""
		jsonData := []byte("")
		switch user_subscription {
		case "vip":
			url = "https://www.hackthebox.com/api/v4/vm/spawn"
			jsonData = []byte(`{"machine_id": ` + machine_id + `}`)
		default:
			url = "https://www.hackthebox.com/api/v4/machine/play/" + machine_id
			jsonData = []byte("{}")
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
	rootCmd.AddCommand(startCmd)
	startCmd.Flags().StringVarP(&machineChoosen, "machine", "m", "", "Machine name")
	startCmd.MarkFlagRequired("machine")
}
