package cmd

import (
	"fmt"
	"os"

	"github.com/GoToolSharing/htb-cli/config"
	"github.com/GoToolSharing/htb-cli/lib/submit"
	"github.com/GoToolSharing/htb-cli/lib/webhooks"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var submitCmd = &cobra.Command{
	Use:   "submit",
	Short: "Submit credentials (machines / challenges / release arena)",
	Long:  "This command allows for the submission of user and root flags discovered on vulnerable machines / challenges",
	Run: func(cmd *cobra.Command, args []string) {
		config.GlobalConfig.Logger.Info("Submit command executed")
		difficultyParam, err := cmd.Flags().GetInt("difficulty")
		if err != nil {
			config.GlobalConfig.Logger.Error("", zap.Error(err))
			os.Exit(1)
		}

		machineNameParam, err := cmd.Flags().GetString("machine_name")
		if err != nil {
			config.GlobalConfig.Logger.Error("", zap.Error(err))
			os.Exit(1)
		}

		challengeNameParam, err := cmd.Flags().GetString("challenge_name")
		if err != nil {
			config.GlobalConfig.Logger.Error("", zap.Error(err))
			os.Exit(1)
		}

		fortressNameParam, err := cmd.Flags().GetString("fortress_name")
		if err != nil {
			config.GlobalConfig.Logger.Error("", zap.Error(err))
			os.Exit(1)
		}

		if challengeNameParam != "" || machineNameParam != "" {
			if difficultyParam == 0 {
				fmt.Println("required flag(s) 'difficulty' not set")
				os.Exit(1)
			}
		}

		var modeType string
		var modeValue string

		if fortressNameParam != "" {
			modeType = "fortress"
			modeValue = fortressNameParam
		} else if machineNameParam != "" {
			modeType = "machine"
			modeValue = machineNameParam
		} else if challengeNameParam != "" {
			modeType = "challenge"
			modeValue = challengeNameParam
		} else {
			modeType = "release-arena"
			modeValue = ""
		}

		config.GlobalConfig.Logger.Debug(fmt.Sprintf("Mode type: %s", modeType))
		config.GlobalConfig.Logger.Debug(fmt.Sprintf("Mode value: %s", modeValue))

		output, err := submit.CoreSubmitCmd(difficultyParam, modeType, modeValue)
		if err != nil {
			config.GlobalConfig.Logger.Error("", zap.Error(err))
			os.Exit(1)
		}

		fmt.Println(output)

		err = webhooks.SendToDiscord("submit", output)
		if err != nil {
			config.GlobalConfig.Logger.Error("", zap.Error(err))
			os.Exit(1)
		}
		config.GlobalConfig.Logger.Info("Exit submit command correctly")
	},
}

func init() {
	rootCmd.AddCommand(submitCmd)
	submitCmd.Flags().StringP("machine_name", "m", "", "Machine Name")
	submitCmd.Flags().StringP("challenge_name", "c", "", "Challenge Name")
	submitCmd.Flags().StringP("fortress_name", "f", "", "Fortress Name")
	submitCmd.Flags().IntP("difficulty", "d", 0, "Difficulty")
	// err := submitCmd.MarkFlagRequired("difficulty")
	// if err != nil {
	// 	fmt.Println(err)
	// }
}
