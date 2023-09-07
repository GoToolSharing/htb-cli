package cmd

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/QU35T-code/htb-cli/utils"
	"github.com/spf13/cobra"
)

var flagCmd = &cobra.Command{
	Use:   "flag",
	Short: "Submit a flag (user and root)",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			fmt.Println("USAGE : htb-cli flag FLAG DIFFICULTY(1:10)")
			return
		}
		flag := args[0]
		difficulty, err := strconv.Atoi(args[1])
		if err != nil {
			log.Println(err)
		}
		if difficulty < 1 || difficulty > 10 {
			fmt.Println("The difficulty must be in 1 and 10")
			os.Exit(1)
		}
		machine_id := utils.GetActiveMachineID()
		machine_name := utils.GetActiveMachineName(machine_id)
		machine_id = fmt.Sprintf("%v", machine_id)
		url := "https://www.hackthebox.com/api/v4/machine/own"
		difficultyString := strconv.Itoa(difficulty * 10)
		var jsonData = []byte(`{"flag": "` + flag + `", "id": ` + machine_id.(string) + `, "difficulty": ` + difficultyString + `}`)
		resp := utils.HtbPost(url, jsonData)
		message := utils.ParseJsonMessage(resp, "message")
		fmt.Printf("Machine : %v\n\n%v", machine_name, message)
	},
}

func init() {
	rootCmd.AddCommand(flagCmd)
}
