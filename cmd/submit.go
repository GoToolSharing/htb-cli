package cmd

import (
	"fmt"
	"log"

	"github.com/GoToolSharing/htb-cli/config"
	"github.com/GoToolSharing/htb-cli/lib/submit"
	"github.com/GoToolSharing/htb-cli/lib/webhooks"
	"github.com/spf13/cobra"
)

// submitCmd defines the "submit" command for submitting flags for machines or challenges.
var submitCmd = &cobra.Command{
	Use:   "submit",
	Short: "Submit credentials (machines / challenges / arena)",
	Long:  "This command allows for the submission of user and root flags discovered on vulnerable machines / challenges",
	Run: func(cmd *cobra.Command, args []string) {
		difficultyParam, err := cmd.Flags().GetInt("difficulty")
		if err != nil {
			fmt.Println(err)
			return
		}

		machineNameParam, err := cmd.Flags().GetString("machine_name")
		if err != nil {
			fmt.Println(err)
			return
		}

		challengeNameParam, err := cmd.Flags().GetString("challenge_name")
		if err != nil {
			fmt.Println(err)
			return
		}

		output, err := submit.CoreSubmitCmd(difficultyParam, machineNameParam, challengeNameParam)
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
	submitCmd.Flags().StringP("machine_name", "m", "", "Machine Name")
	submitCmd.Flags().StringP("challenge_name", "c", "", "Challenge Name")
	submitCmd.Flags().IntP("difficulty", "d", 0, "Difficulty")
	err := submitCmd.MarkFlagRequired("difficulty")
	if err != nil {
		fmt.Println(err)
	}
}
