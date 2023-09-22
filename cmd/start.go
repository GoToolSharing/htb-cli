package cmd

import (
	"fmt"
	"log"
	"net/http"

	"github.com/GoToolSharing/htb-cli/utils"
	"github.com/spf13/cobra"
)

var machineChoosen string

func coreStartCmd(machineChoosen string, proxyParam string) (string, error) {
	machine_id := utils.SearchItemIDByName(machineChoosen, proxyParam, "Machine")
	log.Println("Machine ID :", machine_id)
	machine_type := utils.GetMachineType(machine_id, proxyParam)
	log.Println("Machine Type :", machine_type)
	user_subscription := utils.GetUserSubscription(proxyParam)
	log.Println("User subscription :", user_subscription)
	if machine_type == "release" {
		url := "https://www.hackthebox.com/api/v4/arena/start"
		resp, err := utils.HtbRequest(http.MethodPost, url, proxyParam, []byte(`{}`))
		if err != nil {
			return "", err
		}
		message := utils.ParseJsonMessage(resp, "message")
		return message.(string), nil
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
		return "", err
	}
	message := utils.ParseJsonMessage(resp, "message")
	return message.(string), nil
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start a machine",
	Long:  `Starts a Hackthebox machine specified in argument`,
	Run: func(cmd *cobra.Command, args []string) {
		output, err := coreStartCmd(machineChoosen, proxyParam)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(output)
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.Flags().StringVarP(&machineChoosen, "machine", "m", "", "Machine name")
	startCmd.MarkFlagRequired("machine")
}
