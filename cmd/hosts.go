package cmd

import (
	"os"

	"github.com/GoToolSharing/htb-cli/config"
	"github.com/GoToolSharing/htb-cli/lib/hosts"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var hostsCmd = &cobra.Command{
	Use:   "hosts",
	Short: "Interact with hosts file",
	Run: func(cmd *cobra.Command, args []string) {
		config.GlobalConfig.Logger.Info("Hosts command executed")
		addParam, err := cmd.Flags().GetString("add")
		if err != nil {
			config.GlobalConfig.Logger.Error("", zap.Error(err))
			os.Exit(1)
		}

		deleteParam, err := cmd.Flags().GetString("delete")
		if err != nil {
			config.GlobalConfig.Logger.Error("", zap.Error(err))
			os.Exit(1)
		}

		ipParam, err := cmd.Flags().GetString("ip")
		if err != nil {
			config.GlobalConfig.Logger.Error("", zap.Error(err))
			os.Exit(1)
		}

		if addParam != "" && deleteParam != "" {
			config.GlobalConfig.Logger.Error("Only one parameter is allowed")
			os.Exit(1)
		}

		if addParam != "" {
			err = hosts.AddEntryToHosts(ipParam, addParam)
			if err != nil {
				config.GlobalConfig.Logger.Error("", zap.Error(err))
				os.Exit(1)
			}
		}

		if deleteParam != "" {
			err = hosts.RemoveEntryFromHosts(ipParam, deleteParam)
			if err != nil {
				config.GlobalConfig.Logger.Error("", zap.Error(err))
				os.Exit(1)
			}
		}

		config.GlobalConfig.Logger.Info("Exit hosts command correctly")
	},
}

func init() {
	rootCmd.AddCommand(hostsCmd)
	hostsCmd.Flags().StringP("add", "a", "", "Add a new entry")
	hostsCmd.Flags().StringP("ip", "i", "", "IP Address")
	hostsCmd.Flags().StringP("delete", "d", "", "Delete an entry")
}
