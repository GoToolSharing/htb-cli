package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/GoToolSharing/htb-cli/config"
	"github.com/GoToolSharing/htb-cli/lib/start"
	"github.com/GoToolSharing/htb-cli/lib/utils"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func createProfileFile(name, ip, machineID, seasonal, host string) error {
	profilePath := filepath.Join(config.BaseDirectory, "profile")

	content := []byte(fmt.Sprintf("[Default]\nname=%s\nip=%s\nmachineID=%s\nseasonal=%s\n", name, ip, machineID, seasonal))
	err := os.WriteFile(profilePath, content, 0644)
	if err != nil {
		return err
	}

	// fmt.Println("Use the following command to update your profile.")
	// fmt.Println("source ~/.local/htb-cli/profile")

	return nil
}

var playCmd = &cobra.Command{
	Use:   "play",
	Short: "Configure the pentest environment",
	Run: func(cmd *cobra.Command, args []string) {
		releaseParam, err := cmd.Flags().GetBool("release")
		if err != nil {
			config.GlobalConfig.Logger.Error("", zap.Error(err))
			os.Exit(1)
		}

		var machineID string
		var seasonal string
		host := ""
		var machineName string
		var machineIP string

		if releaseParam {
			seasonal = "true"
			machineID, err = utils.SearchLastReleaseArenaMachine()
			config.GlobalConfig.Logger.Debug(fmt.Sprintf("Machine ID : %s", machineID))
			if err != nil {
				config.GlobalConfig.Logger.Error("", zap.Error(err))
				os.Exit(1)
			}
			machineName, err = utils.GetActiveReleaseArenaMachineName()
			if err != nil {
				config.GlobalConfig.Logger.Error("", zap.Error(err))
				os.Exit(1)
			}
			output, err := start.CoreStartCmd("", machineID)
			if err != nil {
				config.GlobalConfig.Logger.Error("", zap.Error(err))
				os.Exit(1)
			}

			machineIP, err = utils.GetActiveReleaseArenaMachineIP()
			if err != nil {
				config.GlobalConfig.Logger.Error("", zap.Error(err))
				os.Exit(1)
			}

			fmt.Println(output)

		}

		err = createProfileFile(machineName, machineIP, machineID, seasonal, host)
		if err != nil {
			fmt.Println("Error creating profile file:", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(playCmd)
	playCmd.Flags().BoolP("release", "", false, "Release arena")
}
