package cmd

import (
	"fmt"
	"log"
	"net/http"

	"github.com/GoToolSharing/htb-cli/utils"
	"github.com/spf13/cobra"
)

func coreStopCmd(proxyParam string) (string, error) {
	machine_id := utils.GetActiveMachineID(proxyParam)
	if machine_id == "" {
		return "No machine is running", nil
	}
	log.Println("Machine ID :", machine_id)
	machine_type := utils.GetMachineType(machine_id, "")
	log.Println("Machine Type :", machine_type)
	user_subscription := utils.GetUserSubscription(proxyParam)
	log.Println("User subscription :", user_subscription)

	if machine_type == "release" {
		url := "https://www.hackthebox.com/api/v4/arena/stop"
		resp, err := utils.HtbRequest(http.MethodPost, url, proxyParam, []byte(`{}`))
		if err != nil {
			log.Fatal(err)
		}
		message := utils.ParseJsonMessage(resp, "message")
		return message.(string), nil
	}

	url := ""
	jsonData := []byte("")
	switch user_subscription {
	case "vip":
		url = "https://www.hackthebox.com/api/v4/vm/terminate"
		jsonData = []byte(`{"machine_id": ` + machine_id + `}`)
	default:
		url = "https://www.hackthebox.com/api/v4/machine/stop"
		jsonData = []byte(`{"machine_id": ` + machine_id + `}`)
	}
	resp, err := utils.HtbRequest(http.MethodPost, url, proxyParam, jsonData)
	if err != nil {
		log.Fatal(err)
	}
	message := utils.ParseJsonMessage(resp, "message")
	return message.(string), nil
}

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop the current machine",
	Run: func(cmd *cobra.Command, args []string) {
		output, err := coreStopCmd(proxyParam)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(output)
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)
}
